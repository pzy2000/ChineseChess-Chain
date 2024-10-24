/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/shopspring/decimal"

	"chainmaker.org/chainmaker/contract-utils/standard"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/gogo/protobuf/proto"
	"github.com/panjf2000/ants/v2"
	"github.com/redis/go-redis/v9"
)

// ContractWriteSetData 读写集解析合约数据
type ContractWriteSetData struct {
	ContractName     string
	ContractNameBak  string
	ContractSymbol   string
	ContractAddr     string
	ContractType     string
	ContractByteCode []byte
	Version          string
	RuntimeType      string
	ContractStatus   int32
	BlockHeight      int64
	OrgId            string
	SenderTxId       string
	Sender           string
	SenderAddr       string
	Timestamp        int64
	Decimals         int
}

// ParallelParseContract
//
//	@Description: 并发处理合约数据
//	@param blockInfo 区块数据
//	@param hashType
//	@param dealResult 结果集
//	@return error
func ParallelParseContract(blockInfo *common.BlockInfo, hashType string, dealResult *RealtimeDealResult) error {
	var (
		goRoutinePool *ants.Pool
		mutx          sync.Mutex
		errPool       error
	)

	chainId := blockInfo.Block.Header.ChainId
	blockHeight := int64(blockInfo.Block.Header.BlockHeight)
	errChan := make(chan error, 10)
	if goRoutinePool, errPool = ants.NewPool(10, ants.WithPreAlloc(false)); errPool != nil {
		log.Error(GoRoutinePoolErr + errPool.Error())
		return errPool
	}
	defer goRoutinePool.Release()
	var wg sync.WaitGroup
	for i, tx := range blockInfo.Block.Txs {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		errSub := goRoutinePool.Submit(func(i int, blockInfo *common.BlockInfo, txInfo *common.Transaction) func() {
			return func() {
				defer wg.Done()
				var err error
				payload := txInfo.Payload
				//所有上链的数据都是invoke数据
				if payload.TxType != common.TxType_QUERY_CONTRACT &&
					payload.TxType != common.TxType_ARCHIVE &&
					payload.TxType != common.TxType_SUBSCRIBE &&
					payload.TxType != common.TxType_INVOKE_CONTRACT {
					return
				}

				//根据sender计算Addr,Id,cert
				userInfo := &MemberAddrIdCert{}
				if txInfo.Sender != nil && txInfo.Sender.Signer != nil {
					userInfo, err = getMemberIdAddrAndCert(chainId, hashType, txInfo.Sender.Signer)
					if err != nil {
						log.Error("BuildContractInfo getMemberIdAddrAndCert err: " + err.Error())
					}
				}

				//存证合约
				evidenceList, err := DealEvidence(blockHeight, txInfo, userInfo)
				if err != nil {
					errChan <- err
					return
				}

				//处理通用合约数据
				contractInfo, err := BuildContractInfo(i, blockInfo, txInfo, userInfo)
				if err != nil {
					errChan <- fmt.Errorf("Contract: Build ContractInfo Failed: err:%v", err)
					return
				}

				mutx.Lock()         // 锁定互斥锁
				defer mutx.Unlock() // 使用 defer 确保互斥锁被解锁
				if contractInfo != nil {
					//合约修改事件】
					if dealResult.ContractWriteSetData == nil {
						dealResult.ContractWriteSetData = make(map[string]*ContractWriteSetData, 0)
					}
					dealResult.ContractWriteSetData[contractInfo.SenderTxId] = contractInfo
				}

				//存证合约
				if len(evidenceList) > 0 {
					dealResult.EvidenceList = append(dealResult.EvidenceList, evidenceList...)
				}
			}
		}(i, blockInfo, tx))
		if errSub != nil {
			log.Error("ParallelParseContract submit Failed : " + errSub.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		errLog := <-errChan
		return errLog
	}
	return nil
}

// GenesisBlockSystemContract
//
//	@Description: 创世区块解析系统合约
//	@param blockInfo
//	@param dealResult
//	@return error
func GenesisBlockSystemContract(blockInfo *common.BlockInfo, dealResult *RealtimeDealResult) error {
	blockHeight := blockInfo.Block.Header.BlockHeight
	timestamp := blockInfo.Block.Header.BlockTimestamp
	//创世区块
	if blockHeight != 0 {
		return nil
	}

	for i, txInfo := range blockInfo.Block.Txs {
		rwSetList := blockInfo.RwsetList[i]
		if rwSetList == nil {
			continue
		}

		//解析读写集
		systemContractMap, err := GenesisBlockGetContractByWriteSet(rwSetList.TxWrites)
		if err != nil {
			return err
		}

		if len(systemContractMap) == 0 {
			continue
		}

		for _, contract := range systemContractMap {
			runtimeType := common.RuntimeType_name[int32(contract.RuntimeType)]
			contractInfo := &db.Contract{
				Name:           contract.Name,
				NameBak:        contract.Name,
				Version:        contract.Version,
				RuntimeType:    runtimeType,
				CreateTxId:     txInfo.Payload.TxId,
				BlockHeight:    int64(blockHeight),
				Addr:           contract.Address,
				ContractStatus: dbhandle.SystemContractStatus,
				ContractType:   ContractStandardNameOTHER,
				Timestamp:      timestamp,
			}
			dealResult.InsertContracts = append(dealResult.InsertContracts, contractInfo)
		}
		//合约修改事件
		//dealResult.ContractTxList = append(dealResult.ContractTxList, txInfo.Payload.TxId)
	}

	return nil
}

// BuildContractInfo
//
//	@Description: 构造合约数据
//	@param i
//	@param blockInfo
//	@param txInfo
//	@param userInfo
//	@return *db.Contract  合约数据
//	@return string  新增或者更新合约
//	@return error
func BuildContractInfo(i int, blockInfo *common.BlockInfo, txInfo *common.Transaction, userInfo *MemberAddrIdCert) (
	*ContractWriteSetData, error) {
	if blockInfo == nil || txInfo == nil || userInfo == nil {
		return nil, nil
	}

	isContractTx := IsContractTx(txInfo)
	//非合约类交易不用处理合约数据
	if !isContractTx {
		return nil, nil
	}

	blockHeight := blockInfo.Block.Header.BlockHeight
	chainId := blockInfo.Block.Header.ChainId
	rwSetList := blockInfo.RwsetList[i]
	payload := txInfo.Payload
	if rwSetList == nil {
		return nil, nil
	}

	//解析读写集
	contractWriteSet, err := GetContractByWriteSet(rwSetList.TxWrites)
	if err != nil || contractWriteSet.ContractResult == nil {
		return nil, nil
	}

	contractResult := contractWriteSet.ContractResult
	decimal, _ := strconv.Atoi(contractWriteSet.Decimal)
	runtimeType := contractResult.RuntimeType.String()
	if runtimeType == RuntimeTypeGo {
		runtimeType = RuntimeTypeDockerGo
	}

	contractInfo := &ContractWriteSetData{
		ContractName:     contractResult.Name,
		ContractNameBak:  contractResult.Name,
		ContractAddr:     contractResult.Address,
		ContractSymbol:   contractWriteSet.Symbol,
		ContractByteCode: contractWriteSet.ByteCode,
		Version:          contractResult.Version,
		RuntimeType:      runtimeType,
		ContractStatus:   int32(contractResult.Status),
		BlockHeight:      int64(blockHeight),
		OrgId:            contractResult.Creator.OrgId,
		SenderTxId:       payload.TxId,
		Sender:           userInfo.UserId,
		SenderAddr:       userInfo.UserAddr,
		Timestamp:        payload.Timestamp,
		Decimals:         decimal,
	}

	// 计算合约类型
	err = GetContractSDKData(chainId, contractInfo, contractWriteSet.ByteCode)
	if err != nil {
		return nil, err
	}

	//敏感词过滤
	_, flag := FilteringSensitive(contractInfo.ContractName)
	if flag {
		contractInfo.ContractName = config.ContractWarnMsg
	}

	return contractInfo, nil
}

// GetContractSDKData
//
//	@Description: 从SDK获取合约类型，简称，小数
//	@param chainId
//	@param contractInfo 合约数据
//	@return string 合约类型
//	@return string 合约检查
//	@return int 合约小数
//	@return error
func GetContractSDKData(chainId string, contractInfo *ContractWriteSetData, byteCode []byte) error {
	//计算合约类型
	contractType, err := GetContractType(chainId, contractInfo.ContractName, contractInfo.RuntimeType, byteCode)
	if err != nil {
		log.Error("BuildContractInfo get contractType err: " + err.Error())
		return err
	}
	contractInfo.ContractType = contractType

	//只有同质化合约才有合约简称和小数
	if contractType == ContractStandardNameCMDFA ||
		contractType == ContractStandardNameEVMDFA {
		//获取合约简称
		if contractInfo.ContractSymbol == "" {
			symbol, _ := GetContractSymbol(contractType, chainId, contractInfo.ContractAddr)
			contractInfo.ContractSymbol = symbol
		}

		//获取合约合约小数
		if contractInfo.Decimals == 0 {
			decimals, _ := GetContractDecimals(contractType, chainId, contractInfo.ContractAddr)
			contractInfo.Decimals = decimals
		}
	}
	return nil
}

// ProcessContractInsertOrUpdate
//
//	@Description: 并发处理合约数据
//	@param chainId
//	@param dealResult
//	@return error
func ProcessContractInsertOrUpdate(chainId string, dealResult RealtimeDealResult) (RealtimeDealResult, error) {
	if len(dealResult.ContractWriteSetData) == 0 {
		return dealResult, nil
	}

	var wg sync.WaitGroup
	var mutx sync.Mutex

	errCh := make(chan error, len(dealResult.ContractWriteSetData))
	for _, contractData := range dealResult.ContractWriteSetData {
		wg.Add(1)
		go func(cd *ContractWriteSetData) {
			defer wg.Done()
			var contractInfo *db.Contract
			// 缓存判断合约是否存在
			contractDB, err := dbhandle.GetContractByCacheOrAddr(chainId, cd.ContractAddr)
			if err != nil {
				errCh <- err
				return
			}

			if contractDB == nil {
				txNum, _ := dbhandle.GetTxNumByContractName(chainId, cd.ContractNameBak)
				contractInfo = &db.Contract{
					Name:             cd.ContractName,
					NameBak:          cd.ContractNameBak,
					Addr:             cd.ContractAddr,
					Version:          cd.Version,
					RuntimeType:      cd.RuntimeType,
					ContractStatus:   cd.ContractStatus,
					ContractType:     cd.ContractType,
					ContractSymbol:   cd.ContractSymbol,
					Decimals:         cd.Decimals,
					TxNum:            txNum,
					OrgId:            cd.OrgId,
					CreateTxId:       cd.SenderTxId,
					CreateSender:     cd.Sender,
					CreatorAddr:      cd.SenderAddr,
					Timestamp:        cd.Timestamp,
					Upgrader:         cd.Sender,
					UpgradeAddr:      cd.SenderAddr,
					UpgradeOrgId:     cd.OrgId,
					UpgradeTimestamp: cd.Timestamp,
				}

				//获取标准化合约数据
				dbFungibleContract, dbNonFungibleContract := dealStandardContract(contractInfo)

				//构造IDA合约数据
				dbIDAContract := dealIDAContractData(contractInfo)

				// 使用互斥锁保护切片操作
				mutx.Lock()
				dealResult.InsertContracts = append(dealResult.InsertContracts, contractInfo)
				if dbFungibleContract != nil {
					dealResult.FungibleContract = append(dealResult.FungibleContract, dbFungibleContract)
				}
				if dbNonFungibleContract != nil {
					dealResult.NonFungibleContract = append(dealResult.NonFungibleContract, dbNonFungibleContract)
				}
				if dbIDAContract != nil {
					dealResult.InsertIDAContracts = append(dealResult.InsertIDAContracts, dbIDAContract)
				}
				mutx.Unlock()
			} else {
				contractInfo = contractDB
				contractInfo.ContractStatus = cd.ContractStatus
				contractInfo.Upgrader = cd.SenderAddr
				contractInfo.UpgradeAddr = cd.SenderAddr
				contractInfo.UpgradeOrgId = cd.OrgId
				contractInfo.UpgradeTimestamp = cd.Timestamp
				contractInfo.Version = cd.Version

				// 使用互斥锁保护切片操作
				mutx.Lock()
				dealResult.UpdateContracts = append(dealResult.UpdateContracts, contractInfo)
				mutx.Unlock()
			}
			//更新合约缓存
			dbhandle.UpdateContractCache(chainId, contractInfo)
		}(contractData)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return dealResult, <-errCh
	}

	return dealResult, nil
}

// dealStandardContract
//
//	@Description: 将合约数据处理成同质化，非同质化合约数据
//	@param contract
//	@return *db.FungibleContract 同质化合约
//	@return *db.NonFungibleContract 非同质化合约
func dealStandardContract(contract *db.Contract) (*db.FungibleContract, *db.NonFungibleContract) {
	switch contract.ContractType {
	case ContractStandardNameCMDFA:
		fallthrough
	case ContractStandardNameEVMDFA:
		//同质化合约
		fungibleContract := &db.FungibleContract{
			ContractName:    contract.Name,
			ContractNameBak: contract.NameBak,
			Symbol:          contract.ContractSymbol,
			ContractAddr:    contract.Addr,
			ContractType:    contract.ContractType,
			TotalSupply:     decimal.Zero,
			Timestamp:       contract.Timestamp,
		}
		return fungibleContract, nil
	case ContractStandardNameCMNFA:
		fallthrough
	case ContractStandardNameEVMNFA:
		//非同质化合约
		nonFungibleContract := &db.NonFungibleContract{
			ContractName:    contract.Name,
			ContractNameBak: contract.NameBak,
			ContractAddr:    contract.Addr,
			ContractType:    contract.ContractType,
			TotalSupply:     decimal.Zero,
			Timestamp:       contract.Timestamp,
		}
		return nil, nonFungibleContract
	}

	return nil, nil
}

// dealIDAContractData
//
//	@Description: 将合约数据处理成同质化，非同质化合约数据
//	@param contract
//	@return *db.FungibleContract 同质化合约
//	@return *db.NonFungibleContract 非同质化合约
func dealIDAContractData(contract *db.Contract) *db.IDAContract {
	if contract.ContractType != standard.ContractStandardNameCMIDA {
		return nil
	}

	idaContract := &db.IDAContract{
		ContractName:    contract.Name,
		ContractNameBak: contract.NameBak,
		ContractAddr:    contract.Addr,
		ContractType:    contract.ContractType,
		Timestamp:       contract.Timestamp,
	}
	return idaContract
}

// GetContractType
//
//	@Description: 获取合约类型
//	@param chainId
//	@param contractName 合约名称
//	@param runtimeType
//	@param bytecode 创建合约bytecode
//	@return string
//	@return error
func GetContractType(chainId, contractName, runtimeType string, bytecode []byte) (string, error) {
	var err error
	contractType := ContractStandardNameOTHER
	if runtimeType == RuntimeTypeDockerGo {
		contractType, err = DockerGetContractType(chainId, contractName)
		if err != nil {
			log.Errorf("【sdk】GetContractType docker go err :%v", err)
			//失败重试一次
			contractType, err = DockerGetContractType(chainId, contractName)
			if err != nil {
				log.Errorf("【sdk】GetContractType docker go err :%v", err)
			}
		}
		return contractType, nil
	} else if runtimeType == RuntimeTypeEVM {
		if len(bytecode) == 0 {
			log.Errorf("【sdk】GetContractType EVM err, bytecode is nil")
			return ContractStandardNameOTHER, fmt.Errorf("GetContractType bytecode is nil ")
		}

		//获取4字节列表
		signatures := ExtractFunctionSignatures(bytecode)
		// 检查是否包含所有ERC20函数
		if containsAllFunctions(ContractStandardNameEVMDFA, signatures, copyMap(ERC20Functions)) {
			return ContractStandardNameEVMDFA, nil
		}

		// 检查是否包含所有ERC721函数
		if containsAllFunctions(ContractStandardNameEVMNFA, signatures, copyMap(ERC721Functions)) {
			return ContractStandardNameEVMNFA, nil
		}
	}

	log.Infof("GetContractType Not a standard contract，contractName：%v, runtimeType:%v",
		contractName, runtimeType)
	return contractType, nil
}

// containsAllFunctions
//
//	@Description: 判断EVM合约方法是否在标准合约方法中
//	@param evmType evm合约类型
//	@param signatures 合约方法的byte值
//	@param functionNames 标准合约方法名称
//	@return bool
func containsAllFunctions(evmType string, signatures [][]byte, functionNames map[string]bool) bool {
	// 创建一个通道用于接收找到的函数名
	foundChan := make(chan string, len(signatures))
	ercAbi := GetEvmAbi(evmType)
	if ercAbi == nil {
		log.Errorf("containsAllFunctions unmarshal ercAbi failed, ercAbi is null")
		return false
	}

	// 使用一个 WaitGroup 来等待所有的 goroutine 完成
	var wg sync.WaitGroup

	// 遍历签名并调用 EVMGetMethodName
	var allNameList []string
	for _, sig := range signatures {
		wg.Add(1)
		go func(sig []byte) {
			defer wg.Done()
			name, _ := EVMGetMethodName(ercAbi, sig)
			foundChan <- name
			allNameList = append(allNameList, name)
		}(sig)
	}

	// 等待所有的 goroutine 完成
	go func() {
		wg.Wait()
		close(foundChan)
	}()

	// 从通道中读取找到的函数名
	for range signatures {
		name := <-foundChan
		if _, found := functionNames[name]; found {
			// 如果找到匹配的函数名，则从映射中删除
			delete(functionNames, name)

			// 如果映射为空，则已找到所有函数名
			if len(functionNames) == 0 {
				return true
			}
		}
	}

	if len(functionNames) > 0 {
		allNameListJson, _ := json.Marshal(allNameList)
		functionNamesJson, _ := json.Marshal(functionNames)
		log.Infof("【sdk】EVM ContractType containsAllFunctions allNameList:%v, not have name :%v",
			string(allNameListJson), string(functionNamesJson))
	}
	// 如果映射不为空，则没有找到所有函数名
	return false
}

// GetContractSymbol
//
//	@Description: 获取合约简称
//	@param chainId
//	@param contractType 合约类型
//	@param contractAddr 合约地址
//	@return string 简称
//	@return error
func GetContractSymbol(chainId, contractType, contractAddr string) (string, error) {
	var symbolName string
	var err error
	if contractType == ContractStandardNameCMDFA ||
		contractType == ContractStandardNameCMNFA {
		symbolName, err = DockerGetContractSymbol(chainId, contractAddr)
	}

	if contractType == ContractStandardNameEVMDFA ||
		contractType == ContractStandardNameEVMNFA {
		symbolName, err = EVMGetContractSymbol(chainId, contractAddr, contractType)
	}
	return symbolName, err
}

// GetTotalSupply
//
//	@Description: 获取合约总发行量
//	@param contractType
//	@param chainId
//	@param contractName
//	@return string
//	@return error
func GetTotalSupply(contractType, chainId, contractName string) (string, error) {
	totalSupply := "0"
	var err error
	if contractType == ContractStandardNameCMDFA {
		totalSupply, err = DockerGetTotalSupply(chainId, contractName)
	}

	if contractType == ContractStandardNameEVMDFA {
		totalSupply, err = EvmGetTotalSupply(contractType, chainId, contractName)
	}

	return totalSupply, err
}

// GetContractDecimals
//
//	@Description: 获取合约小数位数
//	@param chainId
//	@param contractType 合约类型
//	@param contractName 合约名称，地址
//	@return int
//	@return error
func GetContractDecimals(chainId, contractType, contractName string) (int, error) {
	var decimals int
	var err error
	if contractType == ContractStandardNameCMDFA {
		decimals, err = DockerGetDecimals(chainId, contractName)
	}

	if contractType == ContractStandardNameEVMDFA {
		decimals, err = EvmGetDecimals(contractType, chainId, contractName)
	}

	return decimals, err
}

// SetLatestContractListCache
//
//	@Description: 缓存最新合约列表
//	@param chainId
//	@param blockHeight
//	@param insertContracts
//	@param updateContracts
func SetLatestContractListCache(chainId string, blockHeight int64, insertContracts, updateContracts []*db.Contract) {
	if len(insertContracts) == 0 && len(updateContracts) == 0 {
		return
	}

	ctx := context.Background()
	//添加缓存合约信息
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisLatestContractList, prefix, chainId)
	//新增合约缓存
	for i, contract := range insertContracts {
		contractJson, err := json.Marshal(contract)
		if err != nil {
			log.Errorf("Error Marshal contract err: %v，redisKey：:%v", err, redisKey)
		}
		cache.GlobalRedisDb.ZAdd(ctx, redisKey, redis.Z{
			Score:  float64(blockHeight*10000 + int64(i)),
			Member: string(contractJson),
		})
	}

	// 保留最新的 10 条区块数据
	cache.GlobalRedisDb.ZRemRangeByRank(ctx, redisKey, 0, -11)
	//更新合约版本，交易数
	UpdateLatestContractCache(chainId, updateContracts)
}

// UpdateLatestContractCache 最新合约列表
func UpdateLatestContractCache(chainId string, updateContracts []*db.Contract) {
	if len(updateContracts) == 0 {
		return
	}

	//如果缓存内的合约版本更新了，也需要实时更新
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisLatestContractList, prefix, chainId)
	// 获取缓存中的合约列表及其 Score
	contractList, err := cache.GlobalRedisDb.ZRangeWithScores(ctx, redisKey, 0, -1).Result()
	if err != nil {
		log.Errorf("Error ZRangeWithScores contract err: %v，redisKey：:%v", err, redisKey)
		return
	}

	updatedContractMap := make(map[string]*db.Contract, 0)
	for _, contract := range updateContracts {
		updatedContractMap[contract.Addr] = contract
	}

	// 遍历合约列表，找到需要更新的合约
	for _, contractWithScore := range contractList {
		var contract db.Contract
		err = json.Unmarshal([]byte(contractWithScore.Member.(string)), &contract)
		if err != nil {
			log.Errorf("Error Unmarshal contract err: %v，redisKey：:%v, redisRes:%v", err, redisKey,
				contractWithScore.Member)
			return
		}

		// 检查是否是需要更新的合约
		if _, ok := updatedContractMap[contract.Addr]; !ok {
			continue
		}

		// 从缓存中删除旧的合约数据
		cache.GlobalRedisDb.ZRem(ctx, redisKey, contractWithScore.Member)

		// 将更新后的合约数据添加到缓存中，使用原来的 Score
		updateContract := updatedContractMap[contract.Addr]
		redisContract := contract
		redisContract.Version = updateContract.Version
		redisContract.ContractStatus = updateContract.ContractStatus
		redisContract.UpgradeAddr = updateContract.UpgradeAddr
		redisContract.UpgradeTimestamp = updateContract.UpgradeTimestamp
		if updateContract.TxNum > 0 {
			redisContract.TxNum = updateContract.TxNum
		}
		if updateContract.EventNum > 0 {
			redisContract.EventNum = updateContract.EventNum
		}
		updatedContractJson, err := json.Marshal(redisContract)
		if err != nil {
			log.Errorf("Error Marshal contract err: %v，redisContract：:%v", err, redisContract)
			return
		}

		cache.GlobalRedisDb.ZAdd(ctx, redisKey, redis.Z{
			Score:  contractWithScore.Score,
			Member: string(updatedContractJson),
		})
	}
}

// GetContractByWriteSet
//
//	@Description: 根据读写接，解析合约数据
//	@param txWriteList
//	@return *db.GetContractWriteSet
//	@return error
func GetContractByWriteSet(txWriteList []*common.TxWrite) (*db.GetContractWriteSet, error) {
	contractWriteSet := &db.GetContractWriteSet{}
	var contractResult common.Contract
	for _, write := range txWriteList {
		if strings.HasPrefix(string(write.Key), "Contract:") {
			err := proto.Unmarshal(write.Value, &contractResult)
			if err != nil {
				return contractWriteSet, err
			}
			contractWriteSet.ContractResult = &contractResult
		} else if strings.HasPrefix(string(write.Key), "ContractByteCode:") {
			contractWriteSet.ByteCode = write.Value
		} else if string(write.Key) == "decimal" {
			contractWriteSet.Decimal = string(write.Value)
		} else if string(write.Key) == "symbol" {
			contractWriteSet.Symbol = string(write.Value)
		}
	}
	return contractWriteSet, nil
}

// GenesisBlockGetContractByWriteSet
//
//	@Description: 根据读写接，解析合约数据
//	@param txWriteList
//	@return *db.GetContractWriteSet
//	@return error
func GenesisBlockGetContractByWriteSet(txWriteList []*common.TxWrite) (map[string]common.Contract, error) {
	systemContractList := make(map[string]common.Contract, 0)
	for _, write := range txWriteList {
		var contractResult common.Contract
		if strings.HasPrefix(string(write.Key), "Contract:") {
			err := proto.Unmarshal(write.Value, &contractResult)
			if err != nil {
				return systemContractList, err
			}

			if contractResult.Address != "" {
				systemContractList[contractResult.Address] = contractResult
			}
		}
	}
	return systemContractList, nil
}

// GetContractMapByAddrs
//
//	@Description: 根据合约地址获取合约数据
//	@param chainId
//	@param contractAddrMap
//	@return map[string]*db.Contract
//	@return error
func GetContractMapByAddrs(chainId string, contractAddrMap map[string]string) (map[string]*db.Contract, error) {
	//获取本次涉及到的合约信息
	contractAddrs := make([]string, 0)
	for _, addr := range contractAddrMap {
		contractAddrs = append(contractAddrs, addr)
	}
	//contractMap, err := dbhandle.GetContractByCacheOrAddrs(chainId, contractAddrs)
	contractMap, err := dbhandle.GetContractByAddrs(chainId, contractAddrs)
	if err != nil {
		return contractMap, err
	}
	return contractMap, err
}
