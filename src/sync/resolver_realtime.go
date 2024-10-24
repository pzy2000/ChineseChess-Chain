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
	"chainmaker_web/src/sync/saveTasks"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"
	"github.com/redis/go-redis/v9"
)

// TxTimeLog 处理时间日志
type TxTimeLog struct {
	TimeStart              int64
	TimeContractEvents     int64
	TimeRealtimeDataHandle int64
	TimeRealtimeSaveToDB   int64
	TimeParseTransactions  int64
	TimeParseContracts     int64
	TimeDealBlock          int64
}

// RealtimeDealResult 区块数据格式化后结果
type RealtimeDealResult struct {
	//BlockDetail 区块数据
	BlockDetail *db.Block
	//UserList 用户数据
	UserList map[string]*db.User
	//Transactions 交易数据
	Transactions map[string]*db.Transaction
	//UpgradeContractTransaction 合约版本交易
	UpgradeContractTx []*db.UpgradeContractTransaction
	//ChainConfigList 链配置数据
	ChainConfigList []*pbConfig.ChainConfig
	//InsertContracts 新增合约
	ContractWriteSetData map[string]*ContractWriteSetData
	InsertContracts      []*db.Contract
	//UpdateContracts 修改合约
	UpdateContracts []*db.Contract
	//FungibleContract
	FungibleContract []*db.FungibleContract
	//NonFungibleContract
	NonFungibleContract []*db.NonFungibleContract
	InsertIDAContracts  []*db.IDAContract
	//EvidenceList 存证合约
	EvidenceList []*db.EvidenceContract
	//ContractEvents 合约事件
	ContractEvents []*db.ContractEvent
	//GasRecordList gas消耗列表
	GasRecordList []*db.GasRecord
	//CrossChainResult 主子链相关数据
	CrossChainResult *db.CrossChainResult
}

// CrossChainSaveDB 跨链主子链数据
type CrossChainSaveDB struct {
	//InsertSubChainList 新增子链
	InsertSubChainList []*db.CrossSubChainData
	//UpdateSubChainList 更新子链
	UpdateSubChainList []*db.CrossSubChainData
	//SubChainBlockHeight 需要更新的子链高度列表
	SubChainBlockHeight map[string]int64
}

// RealtimeDataHandle
//
//	@Description: 订阅解析区块数据
//	@param blockInfo 订阅区块
//	@param hashType
//	@return *RealtimeDealResult 解析格式化数据
//	@return *TxTimeLog 耗时日志
//	@return error
func RealtimeDataHandle(blockInfo *common.BlockInfo, hashType string) (*RealtimeDealResult, *TxTimeLog, error) {
	var err error
	txTimeLog := &TxTimeLog{}
	timeStart := time.Now()
	dealResult := &RealtimeDealResult{
		UserList:             map[string]*db.User{},
		Transactions:         map[string]*db.Transaction{},
		ContractWriteSetData: map[string]*ContractWriteSetData{},
	}
	txTimeLog.TimeStart = timeStart.UnixMilli()
	timestamp := blockInfo.Block.Header.BlockTimestamp

	err = GenesisBlockSystemContract(blockInfo, dealResult)
	if err != nil {
		return dealResult, txTimeLog, err
	}

	errCh := make(chan error, 4)
	var wg sync.WaitGroup
	wg.Add(4)

	//根据交易数据解析，交易，event，GasRecord的数据
	go func() {
		defer wg.Done()
		task := func() error {
			startTime := time.Now()
			//并发处理transactions
			dealResult, err = ParallelParseTransactions(blockInfo, hashType, dealResult)
			if err != nil {
				return err
			}
			txTimeLog.TimeParseTransactions = time.Since(startTime).Milliseconds()
			return nil
		}
		//失败会一直重试
		saveTasks.WithRetry(task, "ParallelParseTransactions", errCh)
	}()

	//解析合约数据
	go func() {
		defer wg.Done()
		task := func() error {
			startTime := time.Now()
			//并发解析所有合约数据
			err = ParallelParseContract(blockInfo, hashType, dealResult)
			if err != nil {
				return err
			}
			txTimeLog.TimeParseContracts = time.Since(startTime).Milliseconds()
			return nil
		}
		//失败会一直重试
		saveTasks.WithRetry(task, "ParallelParseContract", errCh)
	}()

	//解析区块数据
	go func() {
		defer wg.Done()
		task := func() error {
			startTime := time.Now()
			//并发处理Block
			modBlock, errB := DealBlockInfo(blockInfo, hashType)
			if errB != nil {
				return errB
			}
			dealResult.BlockDetail = modBlock
			txTimeLog.TimeDealBlock = time.Since(startTime).Milliseconds()
			return nil
		}
		//失败会一直重试
		saveTasks.WithRetry(task, "DealBlockInfo", errCh)
	}()

	//根据读写集解析链配置，主子链数据
	go func() {
		defer wg.Done()
		task := func() error {
			//并发解析读写集数据
			err = ParallelParseWriteSetData(blockInfo, dealResult)
			if err != nil {
				return err
			}
			return nil
		}
		//失败会一直重试
		saveTasks.WithRetry(task, "ParallelParseWriteSetData", errCh)
	}()

	wg.Wait()
	//-------敏感词过滤-----
	startTime := time.Now()
	_ = filterTxAndEvent(dealResult.Transactions, dealResult.ContractEvents)
	txTimeLog.TimeContractEvents = time.Since(startTime).Milliseconds()

	// ----8----接收错误通道中的错误
	close(errCh)
	for err := range errCh {
		if err != nil {
			// 重试多次仍未成功，停掉链，重新订阅
			log.Errorf("Error: %v", err)
			return dealResult, txTimeLog, err
		}
	}

	//根据主子链网关获取子链详情
	if dealResult.CrossChainResult != nil {
		gateWayIds := dealResult.CrossChainResult.GateWayIds
		err = BuildCrossSubChainData(gateWayIds, dealResult, timestamp)
		if err != nil {
			return dealResult, txTimeLog, err
		}
	}

	txTimeLog.TimeRealtimeDataHandle = time.Since(timeStart).Milliseconds()
	return dealResult, txTimeLog, nil
}

// SetTransactionContract
//
//	@Description: 将合约数据写入交易表
//	@param chainId
//	@param transactionMap
func SetTransactionContract(chainId string, dealResult RealtimeDealResult) {
	transactionMap := dealResult.Transactions
	contractWriteSetData := dealResult.ContractWriteSetData
	//交易信息写入合约数据
	for _, transaction := range transactionMap {
		contractInfo, err := dbhandle.GetContractByCacheOrNameAddr(chainId, transaction.ContractNameBak)
		if contractInfo != nil && err == nil {
			transaction.ContractName = contractInfo.Name
			transaction.ContractNameBak = contractInfo.NameBak
			transaction.ContractAddr = contractInfo.Addr
			transaction.ContractRuntimeType = contractInfo.RuntimeType
			transaction.ContractType = contractInfo.ContractType
		}
	}

	for _, contractTx := range dealResult.UpgradeContractTx {
		if writeSetData, ok := contractWriteSetData[contractTx.TxId]; ok {
			contractTx.ContractName = writeSetData.ContractName
			contractTx.ContractNameBak = writeSetData.ContractNameBak
			contractTx.ContractAddr = writeSetData.ContractAddr
			contractTx.ContractRuntimeType = writeSetData.RuntimeType
			contractTx.ContractVersion = writeSetData.Version
			contractTx.ContractType = writeSetData.ContractType
			contractTx.ContractByteCode = writeSetData.ContractByteCode
		} else if txInfo, ok := transactionMap[contractTx.TxId]; ok {
			contractTx.ContractName = txInfo.ContractName
			contractTx.ContractNameBak = txInfo.ContractNameBak
			contractTx.ContractAddr = txInfo.ContractAddr
			contractTx.ContractRuntimeType = txInfo.ContractRuntimeType
			contractTx.ContractVersion = txInfo.ContractVersion
			contractTx.ContractType = txInfo.ContractType
		}
	}
}

// RealtimeDataSaveToDB
//
//	@Description:  顺序插入同步处理数据
//	@param chainId 链ID
//	@param blockHeight 区块高度
//	@param dealResult 格式化后待存储DB的数据
//	@param txTimeLog 耗时日志
//	@return error
func RealtimeDataSaveToDB(chainId string, blockHeight int64, dealResult RealtimeDealResult,
	txTimeLog *TxTimeLog) error {
	timeStart := time.Now()
	var err error
	// 检查是否启用了 gas
	dealResult, err = checkGasEnabled(chainId, dealResult)
	if err != nil {
		return err
	}

	//处理合约数据，判断是新增合约还是更新合约
	dealResult, err = ProcessContractInsertOrUpdate(chainId, dealResult)
	if err != nil {
		return err
	}

	//主子链-确认transfer是否存在，因为是并发处理，处理过程中无法确认是否存在
	if dealResult.CrossChainResult != nil && len(dealResult.CrossChainResult.SaveCrossCycleTx) > 0 {
		err = InsertOrUpdateCrossCycleTx(chainId, dealResult)
		if err != nil {
			return err
		}
	}

	//将合约数据写入交易表
	SetTransactionContract(chainId, dealResult)

	// 执行数据插入任务
	err = executeDataInsertTasks(chainId, dealResult)
	if err != nil {
		return err
	}

	//最后保存block数据
	err = dbhandle.InsertBlock(chainId, dealResult.BlockDetail)
	if err != nil {
		return err
	}

	//异步更新使用
	//设置缓存数据，
	setDelayedUpdateCache(chainId, blockHeight, dealResult)
	//缓存主子链流转数据
	SetCrossSubChainCrossCache(chainId, blockHeight, dealResult)

	//浏览器首页使用
	//最新交易缓存
	BuildLatestTxListCache(chainId, dealResult.Transactions)
	//缓存首页交易总量
	BuildOverviewTxTotalCache(chainId, dealResult.Transactions)
	//缓存最新区块高度
	BuildOverviewMaxBlockHeightCache(chainId, dealResult.BlockDetail)
	//最新区块缓存
	BuildLatestBlockListCache(chainId, dealResult.BlockDetail)
	//最新合约缓存
	SetLatestContractListCache(chainId, blockHeight, dealResult.InsertContracts, dealResult.UpdateContracts)

	txTimeLog.TimeRealtimeSaveToDB = time.Since(timeStart).Milliseconds()
	return nil
}

// executeDataInsertTasks
//
//	@Description: 执行数据并发插入任务
//	@param chainId
//	@param dealResult
//	@return error
func executeDataInsertTasks(chainId string, dealResult RealtimeDealResult) error {
	var err error
	// 初始化重试计数映射
	retryCountMap := &sync.Map{}
	// 创建一个可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 创建任务列表
	tasksList := createTasks(chainId, dealResult)
	// 创建一个错误通道
	errCh := make(chan error, len(tasksList))

	// 并发执行无依赖任务
	var wg sync.WaitGroup
	wg.Add(len(tasksList))
	for _, task := range tasksList {
		go saveTasks.ExecuteTaskWithRetry(ctx, &wg, task, retryCountMap, errCh)
	}

	wg.Wait()
	close(errCh)
	if len(errCh) > 0 {
		err = <-errCh
		if err != nil {
			// 取消其他任务
			cancel()
			// 处理错误
			log.Errorf("Error: %v", err)
			return err
		}
	}

	return nil
}

// createTasks
//
//	@Description:  数据插入任务列表
//	@param chainId 链ID
//	@param dealResult 待插入数据
//	@return []saveTasks.Task 任务列表
func createTasks(chainId string, dealResult RealtimeDealResult) []saveTasks.Task {
	// 定义任务列表
	tasksList := []saveTasks.Task{
		{
			Name:     "TaskSaveTransactionsToDB",
			Function: saveTasks.TaskSaveTransactionsToDB,
			Args:     []interface{}{chainId, dealResult.Transactions, dealResult.UpgradeContractTx},
		},
		{
			Name:     "TaskSaveUserToDB",
			Function: saveTasks.TaskSaveUserToDB,
			Args:     []interface{}{chainId, dealResult.UserList},
		},
		{
			Name:     "TaskSaveContractToDB",
			Function: saveTasks.TaskSaveContractToDB,
			Args:     []interface{}{chainId, dealResult.InsertContracts, dealResult.UpdateContracts},
		},
		{
			Name:     "TaskSaveStandardContractToDB",
			Function: saveTasks.TaskSaveStandardContractToDB,
			Args: []interface{}{chainId, dealResult.FungibleContract, dealResult.NonFungibleContract,
				dealResult.InsertIDAContracts},
		},
		{
			Name:     "TaskEvidenceContractToDB",
			Function: saveTasks.TaskEvidenceContractToDB,
			Args:     []interface{}{chainId, dealResult.EvidenceList},
		},
		{
			Name:     "TaskContractEventsToDB",
			Function: saveTasks.TaskContractEventsToDB,
			Args:     []interface{}{chainId, dealResult.ContractEvents},
		},
		{
			Name:     "TaskGasRecordToDB",
			Function: saveTasks.TaskGasRecordToDB,
			Args:     []interface{}{chainId, dealResult.GasRecordList},
		},
		{
			Name:     "TaskSaveChainConfig",
			Function: saveTasks.TaskSaveChainConfig,
			Args:     []interface{}{chainId, dealResult.ChainConfigList},
		},
		{
			Name:     "TaskSaveRelayCrossChainToDB",
			Function: saveTasks.TaskSaveRelayCrossChainToDB,
			Args:     []interface{}{chainId, dealResult.CrossChainResult},
		},
	}

	return tasksList
}

// setDelayedUpdateCache
//
//	@Description: 设置数据缓存，异步计算使用
//	@param chainId
//	@param blockHeight
//	@param dealResult
func setDelayedUpdateCache(chainId string, blockHeight int64, dealResult RealtimeDealResult) {
	contractAddrs := make(map[string]string, 0)
	delayedUpdateData := &GetRealtimeCacheData{
		TxList:         make(map[string]*db.Transaction),
		ContractAddrs:  make(map[string]string, 0),
		GasRecords:     make([]*db.GasRecord, 0),
		ContractEvents: make([]*db.ContractEvent, 0),
		UserInfoMap:    make(map[string]*db.User, 0),
	}

	for _, txInfo := range dealResult.Transactions {
		if txInfo.ContractAddr == "" {
			continue
		}
		contractAddrs[txInfo.ContractAddr] = txInfo.ContractAddr
	}

	if len(dealResult.Transactions) == 0 {
		return
	}

	delayedUpdateData.TxList = dealResult.Transactions
	delayedUpdateData.GasRecords = dealResult.GasRecordList
	delayedUpdateData.ContractEvents = dealResult.ContractEvents
	delayedUpdateData.UserInfoMap = dealResult.UserList
	delayedUpdateData.ContractAddrs = contractAddrs

	prefix := config.GlobalConfig.RedisDB.Prefix
	heightStr := strconv.FormatInt(blockHeight, 10)
	redisKey := fmt.Sprintf(cache.RedisDelayedUpdateData, prefix, chainId, heightStr)
	retJson, _ := json.Marshal(delayedUpdateData)
	// 设置键值对和过期时间
	ctx := context.Background()
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), 15*time.Minute).Err()
}

// GetDelayedUpdateCache
//
//	@Description: 获取数据缓存，异步计算使用
//	@param chainId
//	@param blockHeight
//	@return *GetRealtimeCacheData 缓存数据
func GetDelayedUpdateCache(chainId string, blockHeight int64) *GetRealtimeCacheData {
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	heightStr := strconv.FormatInt(blockHeight, 10)
	redisKey := fmt.Sprintf(cache.RedisDelayedUpdateData, prefix, chainId, heightStr)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes == nil || redisRes.Val() == "" {
		return nil
	}

	cacheResult := &GetRealtimeCacheData{}
	err := json.Unmarshal([]byte(redisRes.Val()), &cacheResult)
	if err != nil {
		return nil
	}
	return cacheResult
}

// SetCrossSubChainCrossCache
//
//	@Description: 缓存跨链信息，异步更新使用
//	@param chainId
//	@param blockHeight
//	@param dealResult
func SetCrossSubChainCrossCache(chainId string, blockHeight int64, dealResult RealtimeDealResult) {
	if dealResult.CrossChainResult == nil {
		return
	}
	crossTransferMap := dealResult.CrossChainResult.CrossTransfer
	if len(crossTransferMap) == 0 {
		return
	}

	crossTransfer := make([]*db.CrossTransactionTransfer, 0)
	for _, transfer := range crossTransferMap {
		crossTransfer = append(crossTransfer, transfer)
	}
	SetCrossTransfersCache(chainId, blockHeight, crossTransfer)
	//SetCrossCycleTxDataCache(chainId, blockHeight, saveCrossCycleTx)
}

// SetCrossTransfersCache
//
//	@Description: 缓存子链信息
//	@param chainId
//	@param blockHeight
//	@param crossTransfers 跨链流转数据
func SetCrossTransfersCache(chainId string, blockHeight int64, crossTransfers []*db.CrossTransactionTransfer) {
	if len(crossTransfers) == 0 {
		return
	}

	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	heightStr := strconv.FormatInt(blockHeight, 10)
	redisKey := fmt.Sprintf(cache.RedisCrossTxTransfers, prefix, chainId, heightStr)
	retJson, err := json.Marshal(crossTransfers)
	if err != nil {
		log.Errorf("【Redis】set cache failed, key:%v, result:%v", redisKey, retJson)
		return
	}
	// 设置键值对和过期时间(1h 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), time.Hour).Err()
}

// GetCrossTransfersCache
//
//	@Description: 获取子链信息缓存
//	@param chainId
//	@param blockHeight
//	@return []*db.CrossTransactionTransfer 跨链流转信息
//	@return error
func GetCrossTransfersCache(chainId string, blockHeight int64) ([]*db.CrossTransactionTransfer, error) {
	ctx := context.Background()
	var cacheResult []*db.CrossTransactionTransfer
	prefix := config.GlobalConfig.RedisDB.Prefix
	heightStr := strconv.FormatInt(blockHeight, 10)
	redisKey := fmt.Sprintf(cache.RedisCrossTxTransfers, prefix, chainId, heightStr)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	// 检查 Redis 错误
	if err := redisRes.Err(); err != nil {
		if err == redis.Nil {
			// 没有发生错误，缓存不存在
			return nil, nil
		}
		// 发生了错误
		log.Errorf("【Redis】get cache failed, key:%v, error:%v", redisKey, err)
		return nil, err
	}

	err := json.Unmarshal([]byte(redisRes.Val()), &cacheResult)
	if err != nil {
		log.Errorf("【Redis】get cache failed, key:%v, result:%v", redisKey, redisRes)
		return nil, err
	}
	return cacheResult, nil
}

// SetCrossCycleTxDataCache
//
//	@Description: 缓存子链跨链交易信息
//	@param chainId
//	@param blockHeight
//	@param crossCycleTxMap 跨链交易数据
func SetCrossCycleTxDataCache(chainId string, blockHeight int64, crossCycleTxMap map[string]*db.CrossCycleTransaction) {
	if len(crossCycleTxMap) == 0 {
		return
	}

	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	heightStr := strconv.FormatInt(blockHeight, 10)
	redisKey := fmt.Sprintf(cache.RedisCrossCycleTxData, prefix, chainId, heightStr)
	retJson, err := json.Marshal(crossCycleTxMap)
	if err != nil {
		log.Errorf("【Redis】set cache failed, key:%v, result:%v", redisKey, retJson)
		return
	}
	// 设置键值对和过期时间(1h 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), time.Hour).Err()
}

// GetCrossCycleTxDataCache
//
//	@Description: 获取子链信息缓存
//	@param chainId
//	@param blockHeight
//	@return map[string]*db.CrossCycleTransaction 跨链交易数据
//	@return error
func GetCrossCycleTxDataCache(chainId string, blockHeight int64) (map[string]*db.CrossCycleTransaction, error) {
	ctx := context.Background()
	cacheResult := make(map[string]*db.CrossCycleTransaction, 0)
	prefix := config.GlobalConfig.RedisDB.Prefix
	heightStr := strconv.FormatInt(blockHeight, 10)
	redisKey := fmt.Sprintf(cache.RedisCrossCycleTxData, prefix, chainId, heightStr)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	// 检查 Redis 错误
	if err := redisRes.Err(); err != nil {
		if err == redis.Nil {
			// 没有发生错误，缓存不存在
			return nil, nil
		}
		// 发生了错误
		log.Errorf("【Redis】get cache failed, key:%v, error:%v", redisKey, err)
		return nil, err
	}

	err := json.Unmarshal([]byte(redisRes.Val()), &cacheResult)
	if err != nil {
		log.Errorf("【Redis】get cache failed, key:%v, result:%v", redisKey, redisRes)
		return nil, err
	}
	return cacheResult, nil
}
