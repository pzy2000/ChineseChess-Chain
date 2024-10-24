package sync

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/utils"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/gogo/protobuf/proto"
)

// BuildCrossSubChainData
//
//	@Description: 根据网关id获取子链信息
//	@param gateWayIds 网关id
//	@param dealResult
//	@param timestamp
//	@return error
func BuildCrossSubChainData(gateWayIds []int64, dealResult *RealtimeDealResult, timestamp int64) error {
	crossSubChainList := make([]*db.CrossSubChainData, 0)
	//主链新建，更新子链网关
	for _, gateWayId := range gateWayIds {
		//根据gateWayId获取子链数据
		crossSubChains := fetchCrossSubChainData(gateWayId, timestamp)
		if len(crossSubChains) > 0 {
			crossSubChainList = append(crossSubChainList, crossSubChains...)
		}
	}

	if len(crossSubChainList) == 0 {
		return nil
	}

	dealResult.CrossChainResult.SaveSubChainList = crossSubChainList
	return nil
}

// fetchCrossSubChainData 根据gateWayId获取子链数据
func fetchCrossSubChainData(gateWayId, timestamp int64) []*db.CrossSubChainData {
	crossSubChainList := make([]*db.CrossSubChainData, 0)
	gatewayInfo, err := utils.GetCrossGatewayInfo(gateWayId)
	if err != nil || gatewayInfo == nil || len(gatewayInfo.SubChainInfo) == 0 {
		log.Errorf("【Realtime deal】http get sub chain failed, err:%v, gateWayId:%d", err, gateWayId)
		return crossSubChainList
	}

	for _, subChain := range gatewayInfo.SubChainInfo {
		chainDB := &db.CrossSubChainData{
			SubChainId:      subChain.SubChainId,
			ChainId:         subChain.ChainId,
			ChainName:       subChain.ChainName,
			GatewayId:       gatewayInfo.GatewayId,
			GatewayName:     gatewayInfo.GatewayName,
			GatewayAddr:     gatewayInfo.Address,
			ChainType:       subChain.ChainType,
			TxVerifyType:    gatewayInfo.TxVerifyType,
			Enable:          gatewayInfo.Enable,
			CrossCa:         gatewayInfo.CrossCa,
			SdkClientCrt:    gatewayInfo.SdkClientCrt,
			SdkClientKey:    gatewayInfo.SdkClientKey,
			SpvContractName: subChain.SpvContractName,
			Introduction:    subChain.Introduction,
			ExplorerAddr:    subChain.ExplorerAddress,
			ExplorerTxAddr:  subChain.ExplorerTxAddress,
			Timestamp:       timestamp,
		}

		//Grpc-获取子链健康状态
		chainOk, errGrpc := CheckSubChainStatus(chainDB)
		if errGrpc != nil {
			log.Errorf("[load_other] CheckSubChainStatus GetCredentialsByCA failed, err:%v", errGrpc)
		}

		if chainOk {
			chainDB.Status = dbhandle.SubChainStatusSuccess
		} else {
			chainDB.Status = dbhandle.SubChainStatusFail
		}
		crossSubChainList = append(crossSubChainList, chainDB)
	}

	return crossSubChainList
}

// GetMainCrossTransaction 获取主链跨链交易
func GetMainCrossTransaction(crossChainInfo *tcipCommon.CrossChainInfo,
	txInfo *common.Transaction) *db.CrossMainTransaction {
	if crossChainInfo == nil {
		return nil
	}
	crossChainMsg, _ := json.Marshal(crossChainInfo.CrossChainMsg)
	//跨链交易
	crossTransaction := &db.CrossMainTransaction{
		TxId:      txInfo.Payload.TxId,
		CrossId:   crossChainInfo.CrossChainId,
		ChainMsg:  string(crossChainMsg),
		CrossType: int32(crossChainInfo.CrossType),
		Status:    int32(crossChainInfo.State),
		Timestamp: txInfo.Payload.Timestamp,
	}
	return crossTransaction
}

// GetBusinessTransaction 跨链-具体执行的交易
func GetBusinessTransaction(chainId string,
	crossChainInfo *tcipCommon.CrossChainInfo) map[string]*db.CrossBusinessTransaction {
	executionTxMap := make(map[string]*db.CrossBusinessTransaction, 0)
	if crossChainInfo == nil {
		return executionTxMap
	}

	//跨链周期结束在解析业务交易数据，防止重复解析
	if crossChainInfo.State != tcipCommon.CrossChainStateValue_CONFIRM_END &&
		crossChainInfo.State != tcipCommon.CrossChainStateValue_CANCEL_END {
		return executionTxMap
	}

	//子链交易
	crossId := crossChainInfo.CrossChainId
	if crossChainInfo.FirstTxContent != nil && crossChainInfo.FirstTxContent.TxContent != nil {
		txContent := crossChainInfo.FirstTxContent.TxContent
		//解析交易数据
		executionTx := BuildExecutionTransaction(txContent)
		isMainChain := IsMainChainGateway(txContent.GatewayId)
		subChainId := txContent.ChainRid
		if isMainChain {
			subChainId = chainId
		}
		executionTx.IsMainChain = isMainChain
		executionTx.CrossId = crossId
		executionTx.SubChainId = subChainId
		executionTx.GatewayId = txContent.GatewayId
		executionTx.TxId = txContent.TxId
		executionTx.TxStatus = int32(txContent.TxResult)
		mapKey := crossId + "_" + subChainId
		if _, ok := executionTxMap[mapKey]; !ok {
			executionTxMap[executionTx.TxId] = executionTx
		}
	}
	if len(crossChainInfo.CrossChainTxContent) > 0 {
		for _, txContent := range crossChainInfo.CrossChainTxContent {
			if txContent.TxContent == nil {
				continue
			}

			//解析交易数据
			executionTx := BuildExecutionTransaction(txContent.TxContent)
			isMainChain := IsMainChainGateway(txContent.TxContent.GatewayId)
			subChainId := txContent.TxContent.ChainRid
			if isMainChain {
				subChainId = chainId
			}
			executionTx.IsMainChain = isMainChain
			executionTx.CrossId = crossId
			executionTx.SubChainId = subChainId
			executionTx.GatewayId = txContent.TxContent.GatewayId

			contractRes, _ := json.Marshal(txContent.TryResult)
			executionTx.CrossContractResult = string(contractRes)
			executionTx.TxId = txContent.TxContent.TxId
			executionTx.TxStatus = int32(txContent.TxContent.TxResult)
			mapKey := crossId + "_" + subChainId
			if _, ok := executionTxMap[mapKey]; !ok {
				executionTxMap[mapKey] = executionTx
			}
		}
	}

	return executionTxMap
}

// GetCrossTxTransfer
//
//	@Description: 跨链-交易流转方向
//	@param chainId 主链id
//	@param blockHeight 主链高度
//	@param crossChainInfo 跨链详情
//	@return []*db.CrossTransactionTransfer 跨链交易流转方向
func GetCrossTxTransfer(chainId string, blockHeight int64,
	crossChainInfo *tcipCommon.CrossChainInfo) []*db.CrossTransactionTransfer {
	//跨链流转交易
	transferList := make([]*db.CrossTransactionTransfer, 0)
	if crossChainInfo.FirstTxContent == nil || crossChainInfo.FirstTxContent.TxContent == nil {
		return transferList
	}

	//只有初始化的时候解析一次就行了（可能执行多次CrossChainStateValue_WAIT_EXECUTE会有多个），后面都是重复数据
	if crossChainInfo.State != tcipCommon.CrossChainStateValue_NEW &&
		crossChainInfo.State != tcipCommon.CrossChainStateValue_WAIT_EXECUTE {
		return transferList
	}

	fromTxContent := crossChainInfo.FirstTxContent.TxContent
	isMainChain := IsMainChainGateway(fromTxContent.GatewayId)
	fromSubChainId := fromTxContent.ChainRid
	if isMainChain {
		fromSubChainId = chainId
	}

	newUUID := uuid.New().String()
	transfer := &db.CrossTransactionTransfer{
		ID:              newUUID,
		BlockHeight:     blockHeight,
		CrossId:         crossChainInfo.CrossChainId,
		FromChainId:     fromSubChainId,
		FromIsMainChain: isMainChain,
		FromGatewayId:   fromTxContent.GatewayId,
	}

	for i := 0; i < len(crossChainInfo.CrossChainMsg); i++ {
		chainMsg := crossChainInfo.CrossChainMsg[i]
		isMainChain = IsMainChainGateway(chainMsg.GatewayId)
		toSubChainId := chainMsg.ChainRid
		if isMainChain {
			toSubChainId = chainId
		}
		transfer.ToChainId = toSubChainId
		transfer.ToIsMainChain = isMainChain
		transfer.ToGatewayId = chainMsg.GatewayId
		transfer.ContractName = chainMsg.ContractName
		transfer.ContractMethod = chainMsg.Method
		transfer.Parameter = chainMsg.Parameter

		transferList = append(transferList, transfer)
	}

	return transferList
}

// GetCrossCycleTransaction 跨链-交易周期数据
func GetCrossCycleTransaction(crossChainInfo *tcipCommon.CrossChainInfo, crossTxInfo *db.CrossCycleTransaction,
	blockHeight, timestamp int64) *db.CrossCycleTransaction {
	if crossChainInfo == nil ||
		crossChainInfo.FirstTxContent == nil ||
		crossChainInfo.FirstTxContent.TxContent == nil {
		return nil
	}

	cycleTransaction := &db.CrossCycleTransaction{
		CrossId:     crossChainInfo.CrossChainId,
		Status:      int32(crossChainInfo.State),
		BlockHeight: blockHeight,
	}
	//结束状态，更新结束时间
	if crossChainInfo.State == tcipCommon.CrossChainStateValue_CONFIRM_END ||
		crossChainInfo.State == tcipCommon.CrossChainStateValue_CANCEL_END {
		//跨链交易提交，回滚，周期结束
		cycleTransaction.EndTime = timestamp
		//计算周期时长
		if crossTxInfo.StartTime > 0 {
			cycleTransaction.Duration = cycleTransaction.EndTime - crossTxInfo.StartTime
		}
	}

	return cycleTransaction
}

// GetCrossCycleInsertTx 跨链-交易周期数据
func GetCrossCycleInsertTx(crossChainInfo *tcipCommon.CrossChainInfo,
	blockHeight, timestamp int64) *db.CrossCycleTransaction {
	if crossChainInfo == nil ||
		crossChainInfo.FirstTxContent == nil {
		return nil
	}

	//跨链交易周期
	newUUID := uuid.New().String()
	cycleTransaction := &db.CrossCycleTransaction{
		ID:          newUUID,
		CrossId:     crossChainInfo.CrossChainId,
		Status:      int32(crossChainInfo.State),
		BlockHeight: blockHeight,
	}
	//解析手笔交易
	startTime := timestamp
	if crossChainInfo.FirstTxContent.TxContent != nil {
		txContent := crossChainInfo.FirstTxContent.TxContent
		businessTx := BuildExecutionTransaction(txContent)
		//新建跨链交易，周期开始
		if businessTx.Timestamp > 0 {
			startTime = businessTx.Timestamp
		}
	}
	cycleTransaction.StartTime = startTime

	//结束状态，更新结束时间
	if crossChainInfo.State == tcipCommon.CrossChainStateValue_CONFIRM_END ||
		crossChainInfo.State == tcipCommon.CrossChainStateValue_CANCEL_END {
		//跨链交易提交，回滚，周期结束
		cycleTransaction.EndTime = timestamp
		//计算周期时长
		cycleTransaction.Duration = timestamp - startTime
	}
	return cycleTransaction
}

// BuildExecutionTransaction
//
//	@Description: 将跨链交易内容解析成交易信息
//	@param txContent
//	@return *db.CrossBusinessTransaction
func BuildExecutionTransaction(txContent *tcipCommon.TxContent) *db.CrossBusinessTransaction {
	var txInfo common.Transaction
	newUUID := uuid.New().String()
	executionTx := &db.CrossBusinessTransaction{
		ID: newUUID,
	}
	if txContent == nil || len(txContent.Tx) == 0 {
		return executionTx
	}

	err := proto.Unmarshal(txContent.Tx, &txInfo)
	payload := txInfo.Payload
	if err != nil || payload == nil {
		log.Errorf("BuildExecutionTransaction txContent json Unmarshal failed, err:%v", err)
		return executionTx
	}

	//构造交易数据
	executionTx = &db.CrossBusinessTransaction{
		ID:                 newUUID,
		TxId:               payload.TxId,
		TxType:             payload.TxType.String(),
		ContractMessage:    txInfo.Result.ContractResult.Message,
		GasUsed:            txInfo.Result.ContractResult.GasUsed,
		ContractResult:     txInfo.Result.ContractResult.Result,
		ContractResultCode: txInfo.Result.ContractResult.Code,
		RwSetHash:          hex.EncodeToString(txInfo.Result.RwSetHash),
		Timestamp:          payload.Timestamp,
		TxStatusCode:       txInfo.Result.Code.String(),
		ContractMethod:     payload.Method,
		ContractName:       payload.ContractName,
	}
	parametersBytes, err := json.Marshal(payload.Parameters)
	if err == nil {
		executionTx.ContractParameters = string(parametersBytes)
	}

	return executionTx
}

// ParseCrossCycleTxTransfer 主子链跨链次数
func ParseCrossCycleTxTransfer(transfers []*db.CrossTransactionTransfer) map[string]map[string]int64 {
	subChainIdMap := make(map[string]map[string]int64, 0)
	if len(transfers) == 0 {
		return subChainIdMap
	}

	for _, transfer := range transfers {
		if _, ok := subChainIdMap[transfer.FromChainId]; !ok {
			subChainIdMap[transfer.FromChainId] = make(map[string]int64, 0)
		}
		subChainIdMap[transfer.FromChainId][transfer.ToChainId]++

		if _, ok := subChainIdMap[transfer.ToChainId]; !ok {
			subChainIdMap[transfer.ToChainId] = make(map[string]int64, 0)
		}
		subChainIdMap[transfer.ToChainId][transfer.FromChainId]++
	}

	return subChainIdMap
}

// DealSubChainCrossChainNum
//
//	@Description: 根据子链跨链流转计算子链跨链交易数量明细
//	@param chainId
//	@param subChainIdMap 本次子链跨链数据
//	@param subChainCrossDB 数据库子链跨链数据
//	@param minHeight 本次批量处理最低区块高度
//	@return []*db.CrossSubChainCrossChain 新增跨链数据
//	@return []*db.CrossSubChainCrossChain  更新跨链数据
//	@return error
func DealSubChainCrossChainNum(chainId string, subChainIdMap map[string]map[string]int64,
	subChainCrossDB []*db.CrossSubChainCrossChain, minHeight int64) ([]*db.CrossSubChainCrossChain,
	[]*db.CrossSubChainCrossChain, error) {
	insertSubChainCross := make([]*db.CrossSubChainCrossChain, 0)
	updateSubChainCross := make([]*db.CrossSubChainCrossChain, 0)
	if len(subChainIdMap) == 0 {
		return insertSubChainCross, updateSubChainCross, nil
	}

	crossSubChainDBMap := make(map[string]map[string]*db.CrossSubChainCrossChain, 0)
	for _, subChain := range subChainCrossDB {
		if _, ok := crossSubChainDBMap[subChain.SubChainId]; !ok {
			crossSubChainDBMap[subChain.SubChainId] = make(map[string]*db.CrossSubChainCrossChain, 0)
		}
		crossSubChainDBMap[subChain.SubChainId][subChain.ChainId] = subChain
	}

	mainChainId := config.GlobalConfig.ChainConf.MainChainId
	mainChainName := config.GlobalConfig.ChainConf.MainChainName
	for subChainId, relatedChains := range subChainIdMap {
		for relatedChainId, txNum := range relatedChains {
			subDBMap, subChainExists := crossSubChainDBMap[subChainId]
			subDB, relatedChainExists := subDBMap[relatedChainId]
			if !subChainExists || !relatedChainExists {
				// 数据库不存在，新增数据
				chainName := mainChainName
				//是否是主链
				if relatedChainId != mainChainId {
					chainName, _ = dbhandle.GetCrossSubChainName(chainId, relatedChainId)
				}

				newUUID := uuid.New().String()
				crossSubChain := &db.CrossSubChainCrossChain{
					ID:          newUUID,
					SubChainId:  subChainId,
					ChainId:     relatedChainId,
					ChainName:   chainName,
					TxNum:       txNum,
					BlockHeight: minHeight,
				}
				insertSubChainCross = append(insertSubChainCross, crossSubChain)
			} else if subDB.BlockHeight < minHeight {
				//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
				// 更新数据
				subDB.TxNum += txNum
				subDB.BlockHeight = minHeight
				updateSubChainCross = append(updateSubChainCross, subDB)
			}
		}
	}

	return insertSubChainCross, updateSubChainCross, nil
}

// DealCrossSubChainTxNum
//
//	@Description: 根据本次跨链流转交易计算子链交易数
//	@param subChainIdMap 本次跨链交易数
//	@param subChainDataDB 数据库子链交易数
//	@return []*db.CrossSubChainData 子链信息
func DealCrossSubChainTxNum(subChainIdMap map[string]map[string]int64,
	subChainDataDB map[string]*db.CrossSubChainData) []*db.CrossSubChainData {
	saveSubChainTxNum := make([]*db.CrossSubChainData, 0)
	for subChainId, relatedChains := range subChainIdMap {
		if _, ok := subChainDataDB[subChainId]; ok {
			var countTx int64
			for _, txNum := range relatedChains {
				countTx += txNum
			}
			subChainDataDB[subChainId].TxNum += countTx
			saveSubChainTxNum = append(saveSubChainTxNum, subChainDataDB[subChainId])
		}
	}
	return saveSubChainTxNum
}

// InsertOrUpdateCrossCycleTx 判断是新增子链周期交易还是更新交易时间
func InsertOrUpdateCrossCycleTx(chainId string, dealResult RealtimeDealResult) error {
	var (
		crossIds       []string
		insertCycleTxs []*db.CrossCycleTransaction
	)
	saveCycleTxs := dealResult.CrossChainResult.SaveCrossCycleTx
	updateCycleTxs := dealResult.CrossChainResult.UpdateCrossCycleTx
	for _, transfer := range saveCycleTxs {
		crossIds = append(crossIds, transfer.CrossId)
	}
	crossTxMap, errDB := dbhandle.GetCrossCycleTransactionById(chainId, crossIds)
	if errDB != nil {
		return errDB
	}

	for _, transfer := range saveCycleTxs {
		if transferDB, ok := crossTxMap[transfer.CrossId]; ok {
			//数据库存在，就更新
			if _, okUp := updateCycleTxs[transfer.CrossId]; !okUp {
				isEnd := IsCrossEnd(transfer.Status)
				if isEnd {
					transfer.Duration = transfer.EndTime - transferDB.StartTime
				}
				dealResult.CrossChainResult.UpdateCrossCycleTx[transfer.CrossId] = transfer
			}
		} else {
			//数据库不存在，就插入
			insertCycleTxs = append(insertCycleTxs, transfer)
		}
	}
	dealResult.CrossChainResult.InsertCrossCycleTx = insertCycleTxs
	return nil
}
