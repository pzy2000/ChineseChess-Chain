package sync

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"encoding/base64"
	"encoding/json"
	"strings"
	"sync"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"
	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/gogo/protobuf/proto"
	"github.com/panjf2000/ants/v2"
)

// CrossContract 跨链合约
type CrossContract struct {
	SubChainId   string
	ContractName string
}

// ParallelParseWriteSetData
//
//	@Description: 并发解析读写集数据，获取链配置，主子链数据
//	@param blockInfo
//	@param dealResult
//	@return error
//
// nolint:gocyclo
func ParallelParseWriteSetData(blockInfo *common.BlockInfo, dealResult *RealtimeDealResult) error {
	var (
		crossResult = &db.CrossChainResult{
			SubChainBlockHeight:   make(map[string]int64, 0),
			GateWayIds:            make([]int64, 0),
			CrossMainTransaction:  make([]*db.CrossMainTransaction, 0),
			CrossTransfer:         make(map[string]*db.CrossTransactionTransfer, 0),
			SaveCrossCycleTx:      make(map[string]*db.CrossCycleTransaction, 0),
			UpdateCrossCycleTx:    make(map[string]*db.CrossCycleTransaction, 0),
			BusinessTxMap:         make(map[string]*db.CrossBusinessTransaction, 0),
			CrossChainContractMap: make(map[string]map[string]string, 0),
		}
		goRoutinePool *ants.Pool
		mutx          sync.Mutex
		err           error
	)
	errChan := make(chan error, 10)
	if goRoutinePool, err = ants.NewPool(10, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	var wg sync.WaitGroup
	chainId := blockInfo.Block.Header.ChainId
	blockHeight := blockInfo.Block.Header.BlockHeight
	for i, tx := range blockInfo.Block.Txs {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		errSub := goRoutinePool.Submit(func(i int, blockInfo *common.BlockInfo, txInfo *common.Transaction) func() {
			return func() {
				defer wg.Done()
				var (
					mainTransaction     *db.CrossMainTransaction
					saveCycleTx         *db.CrossCycleTransaction
					updateCycleTx       *db.CrossCycleTransaction
					crossChainContracts []*CrossContract
				)
				crossTxTransferList := make([]*db.CrossTransactionTransfer, 0)
				businessTxMap := make(map[string]*db.CrossBusinessTransaction, 0)

				//根据写集解析链配置数据
				rwSetList := blockInfo.RwsetList[i]
				chainConfig, chainErr := GetChainConfigByWriteSet(rwSetList, txInfo)
				if chainErr != nil {
					errChan <- chainErr
					return
				}

				//根据写集解析跨链交易数据
				crossChainInfo, crossErr := GetCrossChainInfoByWriteSet(rwSetList, txInfo)
				if crossErr != nil {
					errChan <- crossErr
					return
				}

				//存在跨链交易
				if crossChainInfo != nil && crossChainInfo.CrossChainId != "" {
					//解析跨链数据
					mainTransaction, crossTxTransferList, businessTxMap, saveCycleTx, updateCycleTx,
						crossChainContracts = processCrossChainInfo(chainId, int64(blockHeight), crossChainInfo, txInfo)
				}

				//根据写集解析跨链同步子链区块高度
				subBlockHeight, spvContractName := GetSubChainBlockHeightByWriteSet(rwSetList, txInfo)

				//根据写集解析跨链gateWayId(主链新建，更新子链gateWayId)
				gateWayId := GetSubChainGatewayIdByWriteSet(rwSetList, txInfo)

				mutx.Lock()         // 锁定互斥锁
				defer mutx.Unlock() // 使用 defer 确保互斥锁被解锁
				//修改链配置
				if chainConfig != nil && chainConfig.ChainId != "" {
					dealResult.ChainConfigList = append(dealResult.ChainConfigList, chainConfig)
				}

				//跨链交易
				if mainTransaction != nil {
					crossResult.CrossMainTransaction = append(crossResult.CrossMainTransaction, mainTransaction)
				}

				//跨链-子链网关
				if gateWayId != nil {
					crossResult.GateWayIds = append(crossResult.GateWayIds, *gateWayId)
				}

				//跨链交易流转
				if len(crossTxTransferList) > 0 {
					for _, transfer := range crossTxTransferList {
						mapKey := transfer.CrossId + "_" + transfer.FromChainId + "_" + transfer.ToChainId
						if _, ok := crossResult.CrossTransfer[mapKey]; !ok {
							crossResult.CrossTransfer[mapKey] = transfer
						}
					}
				}

				//更新跨链交易状态
				if saveCycleTx != nil {
					crossId := saveCycleTx.CrossId
					if cycleTx, ok := crossResult.SaveCrossCycleTx[crossId]; ok {
						//是否是更新的状态，如果是换成最新数据
						if saveCycleTx.Status > cycleTx.Status {
							saveCycleTx.StartTime = cycleTx.StartTime
						}
						isEnd := IsCrossEnd(saveCycleTx.Status)
						if isEnd {
							saveCycleTx.Duration = saveCycleTx.EndTime - saveCycleTx.StartTime
						}
					}
					crossResult.SaveCrossCycleTx[crossId] = saveCycleTx
				}
				//更新跨链交易状态
				if updateCycleTx != nil {
					crossId := updateCycleTx.CrossId
					if cycleTx, ok := crossResult.UpdateCrossCycleTx[crossId]; ok {
						//是否是更新的状态，如果是换成最新数据
						if updateCycleTx.Status > cycleTx.Status {
							crossResult.UpdateCrossCycleTx[crossId] = updateCycleTx
						}
					} else {
						crossResult.UpdateCrossCycleTx[crossId] = updateCycleTx
					}
				}

				//跨链-具体执行的交易
				if len(businessTxMap) > 0 {
					for mapKey, executionTx := range businessTxMap {
						if _, ok := crossResult.BusinessTxMap[mapKey]; !ok {
							crossResult.BusinessTxMap[mapKey] = executionTx
						}
					}
				}

				//跨链-子链高度
				if subBlockHeight > 0 && spvContractName != "" {
					if height, ok := crossResult.SubChainBlockHeight[spvContractName]; ok {
						if subBlockHeight > height {
							crossResult.SubChainBlockHeight[spvContractName] = subBlockHeight
						}
					} else {
						crossResult.SubChainBlockHeight[spvContractName] = subBlockHeight
					}
				}

				//跨链合约
				if len(crossChainContracts) > 0 {
					for _, crossContract := range crossChainContracts {
						subChain := crossContract.SubChainId
						subConName := crossContract.ContractName
						if crossResult.CrossChainContractMap[subChain] == nil {
							crossResult.CrossChainContractMap[subChain] = make(map[string]string)
						}
						crossResult.CrossChainContractMap[subChain][subConName] = subConName
					}
				}
			}
		}(i, blockInfo, tx))
		if errSub != nil {
			log.Error("ParallelParseWriteSetData submit Failed : " + err.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	dealResult.CrossChainResult = crossResult
	return nil
}

// processCrossChainInfo
//
//	@Description: 根据区块跨链数据crossChainInfo，解析主子链数据
//	@param chainId 链ID
//	@param blockHeight 链高度
//	@param crossChainInfo 主子链区块数据
//	@param txInfo 交易数据
//	@return *db.CrossMainTransaction z主子链:主链交易
//	@return []*db.CrossTransactionTransfer 主子链交易流转数据
//	@return map[string]*db.CrossBusinessTransaction 主子链具体任务交易
//	@return *db.CrossCycleTransaction 主子链跨链周期数据
//	@return []*CrossContract 主子链跨链合约数据
func processCrossChainInfo(chainId string, blockHeight int64, crossChainInfo *tcipCommon.CrossChainInfo,
	txInfo *common.Transaction) (*db.CrossMainTransaction, []*db.CrossTransactionTransfer,
	map[string]*db.CrossBusinessTransaction, *db.CrossCycleTransaction, *db.CrossCycleTransaction, []*CrossContract) {
	var (
		mainTransaction     *db.CrossMainTransaction
		insertCycleTx       *db.CrossCycleTransaction
		updateCycleTx       *db.CrossCycleTransaction
		crossChainContracts []*CrossContract
		crossTxTransferList []*db.CrossTransactionTransfer
		businessTxMap       map[string]*db.CrossBusinessTransaction
	)

	//存在跨链交易
	if crossChainInfo == nil || crossChainInfo.CrossChainId == "" {
		return mainTransaction, crossTxTransferList, businessTxMap, insertCycleTx, updateCycleTx, crossChainContracts
	}

	timestamp := txInfo.Payload.Timestamp
	crossId := crossChainInfo.CrossChainId
	//判断是否已经存在，存在就不在重复解析
	crossTxInfo, _ := dbhandle.GetCrossCycleById(chainId, crossId)
	if crossTxInfo == nil {
		//跨链-主子链交易流转信息
		crossTxTransferList = GetCrossTxTransfer(chainId, blockHeight, crossChainInfo)
		//跨链合约
		for _, chainMsg := range crossChainInfo.CrossChainMsg {
			crossChainContracts = append(crossChainContracts, &CrossContract{
				SubChainId:   chainMsg.ChainRid,
				ContractName: chainMsg.ContractName,
			})
		}
		//开始，更新跨链周期数据
		insertCycleTx = GetCrossCycleInsertTx(crossChainInfo, blockHeight, timestamp)
	} else {
		//开始，更新跨链周期数据
		updateCycleTx = GetCrossCycleTransaction(crossChainInfo, crossTxInfo, blockHeight, timestamp)
	}

	//跨链-主链交易
	mainTransaction = GetMainCrossTransaction(crossChainInfo, txInfo)

	//跨链-具体执行的交易数据
	businessTxMap = GetBusinessTransaction(chainId, crossChainInfo)
	return mainTransaction, crossTxTransferList, businessTxMap, insertCycleTx, updateCycleTx, crossChainContracts
}

// GetChainConfigByWriteSet
//
//	@Description: 通过读写集解析链配置
//	@param txRWSet
//	@param txInfo
//	@return *pbConfig.ChainConfig
//	@return error
func GetChainConfigByWriteSet(txRWSet *common.TxRWSet, txInfo *common.Transaction) (*pbConfig.ChainConfig, error) {
	chainConfig := &pbConfig.ChainConfig{}
	//是否配置类交易
	isConfigTx := IsConfigTx(txInfo)
	if !isConfigTx || txRWSet == nil {
		return nil, nil
	}

	for _, write := range txRWSet.TxWrites {
		if string(write.Key) == TxReadWriteKeyChainConfig {
			err := proto.Unmarshal(write.Value, chainConfig)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	if chainConfig.ChainId != "" {
		return chainConfig, nil
	}
	return nil, nil
}

// GetCrossChainInfoByWriteSet
//
//	@Description: 通过读写集解析跨链数据
//	@param txRWSet
//	@param txInfo
func GetCrossChainInfoByWriteSet(txRWSet *common.TxRWSet, txInfo *common.Transaction) (
	*tcipCommon.CrossChainInfo, error) {
	crossChainInfo := &tcipCommon.CrossChainInfo{}
	//是否跨链类交易
	isCrossChainTx := IsRelayCrossChainTx(txInfo)
	if !isCrossChainTx || txRWSet == nil {
		return nil, nil
	}

	for _, write := range txRWSet.TxWrites {
		//以c开头的key可以解析成CrossChainInfo
		if strings.HasPrefix(string(write.Key), "c") {
			// Base64 解码
			decoded, err := base64.StdEncoding.DecodeString(string(write.Value))
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(decoded, crossChainInfo)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	if crossChainInfo.CrossChainId != "" {
		return crossChainInfo, nil
	}
	return nil, nil
}

// GetSubChainBlockHeightByWriteSet
//
//	@Description:  通过读写集解析跨链子链区块高度，bh开头的key表示区块高度
//	@param txRWSet
//	@param txInfo
//	@return int64
//	@return string
func GetSubChainBlockHeightByWriteSet(txRWSet *common.TxRWSet, txInfo *common.Transaction) (int64, string) {
	//是否跨链类交易
	isSubChainSpvTx, contractName := IsSubChainSpvContractTx(txInfo)
	if !isSubChainSpvTx || txRWSet == nil {
		return 0, contractName
	}

	var maxBlockHeight int64
	for _, write := range txRWSet.TxWrites {
		//以c开头的key可以解析成CrossChainInfo
		if strings.HasPrefix(string(write.Key), "bh") {
			blockHeight := BlockHeightExtractNumber(string(write.Key))
			if blockHeight > maxBlockHeight {
				maxBlockHeight = blockHeight
			}
		}
	}

	return maxBlockHeight, contractName
}

// GetSubChainGatewayIdByWriteSet
//
//	@Description: 读写集解析子链gateWay，g开头的key
//	@param txRWSet
//	@param txInfo
//	@return *int64
func GetSubChainGatewayIdByWriteSet(txRWSet *common.TxRWSet, txInfo *common.Transaction) *int64 {
	//是否跨链类交易
	isCrossChainTx := IsRelayCrossChainTx(txInfo)
	if !isCrossChainTx || txRWSet == nil {
		return nil
	}

	for _, write := range txRWSet.TxWrites {
		//以g开头的key可以解析成GatewayId
		if strings.HasPrefix(string(write.Key), "g") {
			gatewayId := GatewayIdExtractNumber(string(write.Key))
			return &gatewayId
		}
	}

	return nil
}
