/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/sync/saveTasks"
	"context"
	"sync"
)

// GetRealtimeCacheData 区块同步缓存数据，供异步更新计算使用
type GetRealtimeCacheData struct {
	TxList         map[string]*db.Transaction
	ContractAddrs  map[string]string
	GasRecords     []*db.GasRecord
	ContractEvents []*db.ContractEvent
	CrossTransfers []*db.CrossTransactionTransfer
	UserInfoMap    map[string]*db.User
}

// DelayedUpdateData 异步计算存储数据
type DelayedUpdateData struct {
	InsertSubChainCross []*db.CrossSubChainCrossChain
	UpdateSubChainCross []*db.CrossSubChainCrossChain
	UpdateSubChainData  []*db.CrossSubChainData
	InsertGasList       []*db.Gas
	UpdateGasList       []*db.Gas
	UpdateTxBlack       *db.UpdateTxBlack
	ContractResult      *db.GetContractResult
	FungibleTransfer    []*db.FungibleTransfer
	NonFungibleTransfer []*db.NonFungibleTransfer
	BlockPosition       *db.BlockPosition
	UpdateAccountResult *db.UpdateAccountResult
	TokenResult         *db.TokenResult
	ContractMap         map[string]*db.Contract
	IDAInsertAssetsData *db.IDAAssetsDataDB
	IDAUpdateAssetsData *db.IDAAssetsUpdateDB
}

// GetDBResult 需要用到的数据库数据
type GetDBResult struct {
	GasList                []*db.Gas
	PositionMapList        map[string][]*db.FungiblePosition
	NonPositionMapList     map[string][]*db.NonFungiblePosition
	FungibleContractMap    map[string]*db.FungibleContract
	NonFungibleContractMap map[string]*db.NonFungibleContract
	AddBlackTxList         []*db.Transaction
	DeleteBlackTxList      []*db.BlackTransaction
	CrossSubChainCross     []*db.CrossSubChainCrossChain
	CrossSubChainMap       map[string]*db.CrossSubChainData
	AccountBNSList         []*db.Account
	AccountDIDList         []*db.Account
	AccountDBMap           map[string]*db.Account
	IDAContractMap         map[string]*db.IDAContract
	IDAAssetDetailMap      map[string]*db.IDAAssetDetail
}

type BatchDelayedUpdateLog struct {
	GetRealtimeCacheTime      int64
	DelayedUpdateDataTime     int64
	UpdateDataToDBTime        int64
	UpdateBlockStatusToDBTime int64
}

// BatchDelayedUpdate
//
//	@Description:批量延迟更新
//	@param chainId 链id
//	@param blockHeights 批量处理区块高度
//	@return error
func BatchDelayedUpdate(chainId string, blockHeights []int64) error {
	if len(blockHeights) == 0 {
		return nil
	}

	//获取缓存数据,缓存缺失从数据库查询（同步插入的交易列表，合约类别等计算需要用到的数据）
	delayedUpdateNeedCache, err := GetRealtimeDataCache(chainId, blockHeights)
	if err != nil {
		return err
	}

	//计算所有需要更新的数据
	delayedUpdateData, err := BuildDelayedUpdateData(chainId, blockHeights, delayedUpdateNeedCache)
	if err != nil {
		return err
	}

	//并发插入，更新数据库
	err = ParallelParseUpdateDataToDB(chainId, delayedUpdateData)
	if err != nil {
		return err
	}

	//最后更新区块状态，所有数据更新结束
	err = saveTasks.UpdateBlockStatusToDB(chainId, blockHeights)
	if err != nil {
		return err
	}

	//更新首页合约缓存交易量
	UpdateLatestContractCache(chainId, delayedUpdateData.ContractResult.UpdateContractTxEventNum)
	return nil
}

// GetRealtimeDataCache
//
//	@Description: 根据缓存获取需要异步更新的数据
//	@param chainId
//	@param blockHeights
//	@return *GetRealtimeCacheData 缓存数据
//	@return error
func GetRealtimeDataCache(chainId string, blockHeights []int64) (*GetRealtimeCacheData, error) {
	//获取同步插入缓存数据
	delayedUpdateData := &GetRealtimeCacheData{
		TxList:         make(map[string]*db.Transaction),
		ContractAddrs:  make(map[string]string, 0),
		UserInfoMap:    make(map[string]*db.User, 0),
		GasRecords:     make([]*db.GasRecord, 0),
		ContractEvents: make([]*db.ContractEvent, 0),
		CrossTransfers: make([]*db.CrossTransactionTransfer, 0),
	}

	//缓存缺失的height
	heightDB := make([]int64, 0)
	//缓存缺失的主链height
	crossHeightDB := make([]int64, 0)
	for _, height := range blockHeights {
		//获取缓存数据
		cacheResult := GetDelayedUpdateCache(chainId, height)
		if cacheResult == nil {
			heightDB = append(heightDB, height)
		} else {
			// 合并缓存数据到 delayedUpdateData
			for k, v := range cacheResult.TxList {
				delayedUpdateData.TxList[k] = v
			}
			// 合并缓存数据到 delayedUpdateData
			for userAddr, user := range cacheResult.UserInfoMap {
				delayedUpdateData.UserInfoMap[userAddr] = user
			}

			for _, addr := range cacheResult.ContractAddrs {
				delayedUpdateData.ContractAddrs[addr] = addr
			}

			delayedUpdateData.GasRecords = append(delayedUpdateData.GasRecords, cacheResult.GasRecords...)
			delayedUpdateData.ContractEvents = append(delayedUpdateData.ContractEvents, cacheResult.ContractEvents...)
		}

		//获取主子链缓存
		crossCycleTransfers, err := GetCrossTransfersCache(chainId, height)
		if err != nil {
			crossHeightDB = append(crossHeightDB, height)
		} else {
			delayedUpdateData.CrossTransfers = append(delayedUpdateData.CrossTransfers, crossCycleTransfers...)
		}
	}

	//缓存没有，从数据库获取
	if len(heightDB) > 0 {
		//缓存没有，数据库取数据
		err := GetDelayedUpdateByDB(chainId, heightDB, delayedUpdateData)
		if err != nil {
			return delayedUpdateData, err
		}
	}

	//主子链缓存没有，从数据库获取
	if len(crossHeightDB) > 0 {
		//缓存没有，数据库取数据
		crossCycleTransfers, err := dbhandle.GetCrossCycleTransferByHeight(chainId, crossHeightDB)
		if err != nil {
			return delayedUpdateData, err
		}
		delayedUpdateData.CrossTransfers = append(delayedUpdateData.CrossTransfers, crossCycleTransfers...)
	}

	return delayedUpdateData, nil
}

// BuildDelayedUpdateData
//
//	@Description: 计算所有需要更新的数据
//	@param chainId
//	@param blockHeights 批量处理的区块高度列表
//	@param delayedUpdateCache 同步插入的缓存数据
//	@return DelayedUpdateData 需要更新数据库的结构化数据
func BuildDelayedUpdateData(chainId string, blockHeights []int64, delayedUpdateCache *GetRealtimeCacheData) (
	*DelayedUpdateData, error) {
	//本次批量处理的最新的区块高度，用于确定异常情况重复更新问题
	minHeight := GetMinBlockHeight(blockHeights)

	//获取本次涉及到的合约信息
	contractMap, err := GetContractMapByAddrs(chainId, delayedUpdateCache.ContractAddrs)
	if err != nil {
		return nil, err
	}

	//解析合约event
	topicEventResult := DealTopicEventData(delayedUpdateCache.ContractEvents, contractMap, delayedUpdateCache.TxList)
	//解析出的事件列表
	eventDataList := topicEventResult.ContractEventData

	//主子链解析transfer
	crossCycleTransfers := delayedUpdateCache.CrossTransfers
	//计算主子链跨链次数
	crossSubChainIdMap := ParseCrossCycleTxTransfer(crossCycleTransfers)

	//DB并发获取合约，交易，持仓等数据库数据
	delayGetDBResult, err := DelayParallelParseGetDB(chainId, delayedUpdateCache, contractMap, topicEventResult,
		crossSubChainIdMap)
	if err != nil {
		return nil, err
	}

	//主子链计算跨链交易数
	crossSubChainCrossDB := delayGetDBResult.CrossSubChainCross
	insertSubChainCross, updateSubChainCross, err := DealSubChainCrossChainNum(chainId, crossSubChainIdMap,
		crossSubChainCrossDB, minHeight)
	if err != nil {
		return nil, err
	}

	//主子链计算子链交易数
	saveSubChainTxNum := DealCrossSubChainTxNum(crossSubChainIdMap, delayGetDBResult.CrossSubChainMap)

	// 获取新增，更新账户信息
	accountTxNum, accountNFTNum := DealAccountTxNFTNum(delayedUpdateCache.TxList, eventDataList)
	insertAccountMap, updateAccountMap, accountMap, err := BuildAccountInsertOrUpdate(chainId, minHeight, delayGetDBResult,
		topicEventResult, accountTxNum, accountNFTNum)
	if err != nil {
		return nil, err
	}
	updateAccountResult := &db.UpdateAccountResult{
		InsertAccount: insertAccountMap,
		UpdateAccount: updateAccountMap,
	}

	//计算新增，更新gas数据
	insertGasList, updateGasList := buildGasInfo(delayedUpdateCache.GasRecords, delayGetDBResult.GasList, minHeight)

	//统计transfer流转记录
	fungibleTransfer, nonFungibleTransfer := dealTransferList(eventDataList, contractMap, delayedUpdateCache.TxList)

	//统计token列表
	tokenResult := DealNonFungibleToken(chainId, eventDataList, contractMap, accountMap)

	//统计新增持仓数据
	positionList := BuildPositionList(eventDataList, contractMap, accountMap)
	//计算持仓数据
	positionDBMap := delayGetDBResult.PositionMapList
	nonPositionDBMap := delayGetDBResult.NonPositionMapList
	positionOperates := BuildUpdatePositionData(minHeight, positionList, positionDBMap, nonPositionDBMap)

	//交易黑名单
	updateTxBlack := &db.UpdateTxBlack{
		AddTxBlack:    make([]*db.BlackTransaction, 0),
		DeleteTxBlack: make([]*db.Transaction, 0),
	}
	for _, txInfo := range delayGetDBResult.AddBlackTxList {
		//添加黑名单
		updateTxBlack.AddTxBlack = append(updateTxBlack.AddTxBlack, (*db.BlackTransaction)(txInfo))
	}
	for _, txInfo := range delayGetDBResult.DeleteBlackTxList {
		//删除黑名单
		updateTxBlack.DeleteTxBlack = append(updateTxBlack.DeleteTxBlack, (*db.Transaction)(txInfo))
	}

	//计算合约交易量
	updateContractTxEventNum := UpdateContractTxAndEventNum(minHeight, contractMap, delayedUpdateCache.TxList,
		delayedUpdateCache.ContractEvents)

	//计算持有量和发行总
	//持有人数
	holdCountMap := DealContractHoldCount(positionOperates)
	//发行总量
	totalSupplyMap := DealContractTotalSupply(eventDataList, contractMap)
	fungibleMap := delayGetDBResult.FungibleContractMap
	nonFungibleMap := delayGetDBResult.NonFungibleContractMap

	//计算同质化合约持有人数，发行量最终数据
	updateFTContractMap := DealFungibleContractUpdateData(holdCountMap, totalSupplyMap, fungibleMap, minHeight)
	//计算FT合约交易流转数量
	ftContractTransferMap := FTContractTransferNum(fungibleTransfer, fungibleMap, minHeight)
	mergedFTContract := MergeFTContractMaps(minHeight, updateFTContractMap, ftContractTransferMap)

	//计算非同质化合约持有人数，发行量最终数据
	updateNFTContractMap := DealNonFungibleContractUpdateData(holdCountMap, totalSupplyMap, nonFungibleMap, minHeight)
	//计算NFT合约交易流转数量
	nftContractTransferMap := NFTContractTransferNum(nonFungibleTransfer, nonFungibleMap, minHeight)
	mergedNFTContract := MergeNFTContractMaps(minHeight, updateNFTContractMap, nftContractTransferMap)

	//数据要素计算方式
	//计算IDA合约数据
	idaContractMap := delayGetDBResult.IDAContractMap
	//更新IDA合约数据
	updateIDAContractData := DealIDAContractUpdateData(idaContractMap, topicEventResult.IDAEventData, minHeight)
	//插入IDA数据资产
	dealInsertIDAAssets := DealInsertIDAAssetsData(idaContractMap, topicEventResult.IDAEventData)
	//更新IDA合约资产
	dealUpdateIDAAssets := DealUpdateIDAAssetsData(topicEventResult.IDAEventData, delayGetDBResult.IDAAssetDetailMap)

	//合约更新数据
	contractResult := &db.GetContractResult{
		UpdateContractTxEventNum: updateContractTxEventNum,
		IdentityContract:         topicEventResult.IdentityContract,
		UpdateFungibleContract:   mergedFTContract,
		UpdateNonFungible:        mergedNFTContract,
		UpdateIdaContract:        updateIDAContractData,
	}

	buildDelayedUpdateData := &DelayedUpdateData{
		InsertSubChainCross: insertSubChainCross,
		UpdateSubChainCross: updateSubChainCross,
		UpdateSubChainData:  saveSubChainTxNum,
		InsertGasList:       insertGasList,
		UpdateGasList:       updateGasList,
		UpdateTxBlack:       updateTxBlack,
		ContractResult:      contractResult,
		FungibleTransfer:    fungibleTransfer,
		NonFungibleTransfer: nonFungibleTransfer,
		BlockPosition:       positionOperates,
		UpdateAccountResult: updateAccountResult,
		TokenResult:         tokenResult,
		ContractMap:         contractMap,
		IDAInsertAssetsData: dealInsertIDAAssets,
		IDAUpdateAssetsData: dealUpdateIDAAssets,
	}
	return buildDelayedUpdateData, nil
}

// GetDelayedUpdateByDB
//
//	@Description: 缓存数据如果没有，根据区块高度从数据库查
//	@param chainId
//	@param heightDB
//	@param delayedUpdateData 异步更新需要用到的数据
//	@return error
func GetDelayedUpdateByDB(chainId string, heightDB []int64, delayedUpdateData *GetRealtimeCacheData) error {
	//获取交易列表
	txInfoList, err := dbhandle.GetTxInfoByBlockHeight(chainId, heightDB)
	if err != nil {
		return err
	}

	//解析交易id列表，合约名称列表
	txIds, contractNameMap, txInfoMap := ExtractTxIdsAndContractNames(txInfoList)
	if len(txIds) == 0 {
		return nil
	}

	// 合并缓存数据到 delayedUpdateData
	for k, v := range txInfoMap {
		delayedUpdateData.TxList[k] = v
	}
	for name := range contractNameMap {
		delayedUpdateData.ContractAddrs[name] = name
	}

	errCh := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		//获取gas记录
		gasRecords, err := GetGasRecord(chainId, txIds)
		if err != nil {
			errCh <- err
			return
		}
		if len(gasRecords) > 0 {
			delayedUpdateData.GasRecords = append(delayedUpdateData.GasRecords, gasRecords...)
		}
	}()

	go func() {
		defer wg.Done()
		//获取合约event
		contractEvents, err := GetContractEvents(chainId, txIds)
		if err != nil {
			errCh <- err
			return
		}
		if len(contractEvents) > 0 {
			delayedUpdateData.ContractEvents = append(delayedUpdateData.ContractEvents, contractEvents...)
		}
	}()

	wg.Wait()
	close(errCh)
	for errDB := range errCh {
		if errDB != nil {
			// 重试多次仍未成功，停掉链，重新订阅
			log.Errorf("Error: %v", errDB)
			return errDB
		}
	}

	return nil
}

// ParallelParseUpdateDataToDB
//
//	@Description: 并发更新所有的表数据，失败会进行重试
//	@param chainId
//	@param delayedUpdateData 处理好的更新数据
//	@return error
func ParallelParseUpdateDataToDB(chainId string, delayedUpdateData *DelayedUpdateData) error {
	var err error
	// 数据插入
	// 初始化重试计数映射
	retryCountMap := &sync.Map{}
	// 创建一个可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 创建任务列表
	tasksList := createTasksDelayedUpdate(chainId, delayedUpdateData)
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

// createTasksDelayedUpdate
//
//	@Description: 数据插入任务列表
//	@param chainId
//	@param delayedUpdate 需要插入，更新的数据
//	@return []saveTasks.Task 任务列表
func createTasksDelayedUpdate(chainId string, delayedUpdate *DelayedUpdateData) []saveTasks.Task {
	// 定义任务列表
	tasksList := []saveTasks.Task{
		{
			Name:     "TaskUpdateContractResult",
			Function: saveTasks.TaskUpdateContractResult,
			Args:     []interface{}{chainId, delayedUpdate.ContractResult},
		},
		{
			Name:     "TaskInsertFungibleTransferToDB",
			Function: saveTasks.TaskInsertFungibleTransferToDB,
			Args:     []interface{}{chainId, delayedUpdate.FungibleTransfer},
		},
		{
			Name:     "TaskInsertNonFungibleTransferToDB",
			Function: saveTasks.TaskInsertNonFungibleTransferToDB,
			Args:     []interface{}{chainId, delayedUpdate.NonFungibleTransfer},
		},
		{
			Name:     "TaskSaveAccountListToDB",
			Function: saveTasks.TaskSaveAccountListToDB,
			Args:     []interface{}{chainId, delayedUpdate.UpdateAccountResult},
		},
		{
			Name:     "TaskSaveTokenResultToDB",
			Function: saveTasks.TaskSaveTokenResultToDB,
			Args:     []interface{}{chainId, delayedUpdate.TokenResult},
		},
		{
			Name:     "TaskUpdateTxBlackToDB",
			Function: saveTasks.TaskUpdateTxBlackToDB,
			Args:     []interface{}{chainId, delayedUpdate.UpdateTxBlack},
		},
		{
			Name:     "TaskSaveGasToDB",
			Function: saveTasks.TaskSaveGasToDB,
			Args:     []interface{}{chainId, delayedUpdate.InsertGasList, delayedUpdate.UpdateGasList},
		},
		{
			Name:     "TaskSaveFungibleContractResult",
			Function: saveTasks.TaskSaveFungibleContractResult,
			Args:     []interface{}{chainId, delayedUpdate.ContractResult},
		},
		{
			Name:     "TaskSavePositionToDB",
			Function: saveTasks.TaskSavePositionToDB,
			Args:     []interface{}{chainId, delayedUpdate.BlockPosition},
		},
		{
			Name:     "TaskCrossSubChainCrossToDB",
			Function: saveTasks.TaskCrossSubChainCrossToDB,
			Args:     []interface{}{chainId, delayedUpdate.InsertSubChainCross, delayedUpdate.UpdateSubChainCross},
		},
		{
			Name:     "TaskCrossUpdateSubChainTxNumToDB",
			Function: saveTasks.TaskCrossUpdateSubChainTxNumToDB,
			Args:     []interface{}{chainId, delayedUpdate.UpdateSubChainData},
		},
		{
			Name:     "TaskSaveIDAAssetDataToDB",
			Function: saveTasks.TaskSaveIDAAssetDataToDB,
			Args:     []interface{}{chainId, delayedUpdate.IDAInsertAssetsData},
		},
		{
			Name:     "TaskUpdateIDAAssetDataToDB",
			Function: saveTasks.TaskUpdateIDAAssetDataToDB,
			Args:     []interface{}{chainId, delayedUpdate.IDAUpdateAssetsData},
		},
	}

	return tasksList
}
