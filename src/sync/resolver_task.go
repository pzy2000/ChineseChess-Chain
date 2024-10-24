package sync

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/sync/saveTasks"
	"sync"

	"github.com/panjf2000/ants/v2"
)

type TaskSync interface {
	Run() error
}

type TaskFunc func() error

func (f TaskFunc) Run() error {
	return f()
}

// GetContractEvents 数据库获取event数据
func GetContractEvents(chainId string, txIds []string) ([]*db.ContractEvent, error) {
	contractEvents := make([]*db.ContractEvent, 0)
	if len(txIds) == 0 {
		return contractEvents, nil
	}
	var (
		goRoutinePool *ants.Pool
		mutx          sync.Mutex
		err           error
	)
	// 将交易分割为大小为10的批次
	batches := ParallelParseBatchWhere(txIds, 100)
	errChan := make(chan error, 10)
	if goRoutinePool, err = ants.NewPool(10, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return contractEvents, err
	}

	defer goRoutinePool.Release()
	var wg sync.WaitGroup
	for _, batch := range batches {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(batch []string) func() {
			return func() {
				defer wg.Done()
				//查询数据
				eventList, eventErr := dbhandle.GetEventDataByTxIds(chainId, txIds)
				if eventErr != nil {
					errChan <- eventErr
				}
				mutx.Lock()         // 锁定互斥锁
				defer mutx.Unlock() // 使用 defer 确保互斥锁被解锁
				if len(eventList) > 0 {
					contractEvents = append(contractEvents, eventList...)
				}
			}
		}(batch))
		if errSub != nil {
			log.Errorf("GetContractEvents submit Failed, err:%v", err)
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return contractEvents, err
	}

	return contractEvents, nil
}

// DelayParallelParseGetDB
//
//	@Description:并发获取合约，交易等数据库数据
//	@param chainId
//	@param delayedUpdateCache 同步插入的缓存数据，用于获取where条件
//	@param contractMap 本次涉及到的合约DB数据
//	@param topicEventResult 处理好的event事件数据
//	@param crossSubChainIdMap 本次涉及的主子链中的子链信息
//	@return *GetDBResult DB数据
//	@return error
func DelayParallelParseGetDB(chainId string, delayedUpdateCache *GetRealtimeCacheData,
	contractMap map[string]*db.Contract, topicEventResult *TopicEventResult,
	crossSubChainIdMap map[string]map[string]int64) (*GetDBResult, error) {

	getDBResult := &GetDBResult{
		PositionMapList:        make(map[string][]*db.FungiblePosition, 0),
		NonPositionMapList:     make(map[string][]*db.NonFungiblePosition, 0),
		FungibleContractMap:    make(map[string]*db.FungibleContract, 0),
		NonFungibleContractMap: make(map[string]*db.NonFungibleContract, 0),
		AddBlackTxList:         make([]*db.Transaction, 0),
		DeleteBlackTxList:      make([]*db.BlackTransaction, 0),
		CrossSubChainCross:     make([]*db.CrossSubChainCrossChain, 0),
		CrossSubChainMap:       make(map[string]*db.CrossSubChainData, 0),
		AccountBNSList:         make([]*db.Account, 0),
		AccountDIDList:         make([]*db.Account, 0),
		AccountDBMap:           make(map[string]*db.Account, 0),
		IDAContractMap:         make(map[string]*db.IDAContract, 0),
	}

	// 定义一个任务列表
	tasks := []TaskSync{
		getGasListTask(chainId, delayedUpdateCache, getDBResult),
		getPositionMapTask(chainId, getDBResult, topicEventResult.OwnerAdders),
		getNonPositionMapTask(chainId, getDBResult, topicEventResult.OwnerAdders),
		getFungibleContractTask(chainId, getDBResult, contractMap),
		getNonFungibleContractTask(chainId, getDBResult, contractMap),
		getAddBlackTxListTask(chainId, getDBResult, topicEventResult.AddBlack),
		getDeleteBlackTxListTask(chainId, getDBResult, topicEventResult.DeleteBlack),
		getSubChainCrossListTask(chainId, getDBResult, crossSubChainIdMap),
		getSubChainDBMapTask(chainId, getDBResult, crossSubChainIdMap),
		getAccountByBNSListTask(chainId, getDBResult, topicEventResult.BNSUnBindDomain),
		getAccountByDIDListTask(chainId, getDBResult, topicEventResult.DIDUnBindList),
		getAccountMapTask(chainId, getDBResult, topicEventResult, delayedUpdateCache.UserInfoMap),
		getIDAContractTask(chainId, getDBResult, topicEventResult.IDAEventData),
		getIDAAssetDetailsTask(chainId, getDBResult, topicEventResult.IDAEventData),
	}

	errCh := make(chan error, len(tasks))
	var wg sync.WaitGroup
	// 使用for循环启动所有任务
	for _, task := range tasks {
		wg.Add(1)
		go func(t TaskSync) {
			defer wg.Done()
			saveTasks.WithRetry(t.Run, "DelayParallelParseGetDB", errCh)
		}(task)
	}

	wg.Wait()
	close(errCh)
	for errDB := range errCh {
		if errDB != nil {
			// 重试多次仍未成功，停掉链，重新订阅
			log.Errorf("Error: %v", errDB)
			return nil, errDB
		}
	}

	return getDBResult, nil
}

// getGasListTask
//
//	@Description: 任务：获取gas列表
//	@param chainId
//	@param delayedUpdateCache 同步订阅缓存数据
//	@param getDBResult DB结果数据
//	@return TaskSync 任务方法
func getGasListTask(chainId string, delayedUpdateCache *GetRealtimeCacheData, getDBResult *GetDBResult) TaskSync {
	return TaskFunc(func() error {
		addrList := buildGasAddrList(delayedUpdateCache.GasRecords)
		getGasList, err := dbhandle.GetGasInfoByAddr(chainId, addrList)
		if err != nil {
			return err
		}
		getDBResult.GasList = getGasList
		return nil
	})
}

// getPositionMapTask
//
//	@Description: 任务：根据ownerAdders获取同质化持仓数据
//	@param chainId
//	@param getDBResult 数据库数据
//	@param ownerAdders 账户地址
//	@return TaskSync
func getPositionMapTask(chainId string, getDBResult *GetDBResult, ownerAdders []string) TaskSync {
	return TaskFunc(func() error {
		//根据owner获取合约信息
		positionMapList, err := dbhandle.GetFungiblePositionByOwners(chainId, ownerAdders)
		if err != nil {
			return err
		}
		getDBResult.PositionMapList = positionMapList
		return nil
	})
}

// getNonPositionMapTask
//
//	@Description: 任务：根据ownerAdders获取非同质化持仓数据
//	@param chainId
//	@param getDBResult 数据库数据
//	@param ownerAdders 账户地址
//	@return TaskSync
func getNonPositionMapTask(chainId string, getDBResult *GetDBResult, ownerAdders []string) TaskSync {
	return TaskFunc(func() error {
		//根据owner获取合约信息
		positionMapList, err := dbhandle.GetNonFungiblePositionByOwner(chainId, ownerAdders)
		if err != nil {
			return err
		}
		getDBResult.NonPositionMapList = positionMapList
		return nil
	})
}

// getFungibleContractTask
//
//	@Description: 任务：获取同质化合约数据
//	@param chainId
//	@param getDBResult
//	@param contractMap 本次涉及全部合约
//	@return TaskSync
func getFungibleContractTask(chainId string, getDBResult *GetDBResult, contractMap map[string]*db.Contract) TaskSync {
	return TaskFunc(func() error {
		fungibleAddr := make([]string, 0)
		for _, contract := range contractMap {
			if contract.ContractType == ContractStandardNameCMDFA ||
				contract.ContractType == ContractStandardNameEVMDFA {
				fungibleAddr = append(fungibleAddr, contract.Addr)
			}
		}

		if len(fungibleAddr) == 0 {
			return nil
		}

		//同质化合约
		fungibleContract, err := dbhandle.QueryFungibleContractExists(chainId, fungibleAddr)
		if err != nil {
			return err
		}

		getDBResult.FungibleContractMap = fungibleContract
		return nil
	})
}

// getNonFungibleContractTask
//
//	@Description: 任务：非同质化合约数据
//	@param chainId
//	@param getDBResult
//	@param contractMap 本次涉及全部合约
//	@return TaskSync
func getNonFungibleContractTask(chainId string, getDBResult *GetDBResult,
	contractMap map[string]*db.Contract) TaskSync {
	return TaskFunc(func() error {
		nonFungibleAddr := make([]string, 0)
		for _, contract := range contractMap {
			if contract.ContractType == ContractStandardNameCMNFA ||
				contract.ContractType == ContractStandardNameEVMNFA {
				nonFungibleAddr = append(nonFungibleAddr, contract.Addr)
			}
		}

		if len(nonFungibleAddr) == 0 {
			return nil
		}

		//非同质化合约
		nonFungibleContract, err := dbhandle.QueryNonFungibleContractExists(chainId, nonFungibleAddr)
		if err != nil {
			return err
		}
		getDBResult.NonFungibleContractMap = nonFungibleContract
		return nil
	})
}

// getAddBlackTxListTask
//
//	@Description: 任务：查询需要条加黑名单的交易数据
//	@param chainId
//	@param getDBResult
//	@param addBlackTxIds 黑名单交易列表
//	@return TaskSync
func getAddBlackTxListTask(chainId string, getDBResult *GetDBResult, addBlackTxIds []string) TaskSync {
	return TaskFunc(func() error {
		//添加黑名单交易
		txList, err := dbhandle.BatchQueryTxList(chainId, addBlackTxIds)
		if err != nil {
			return err
		}
		getDBResult.AddBlackTxList = txList
		return nil
	})
}

// getDeleteBlackTxListTask
//
//	@Description: 任务：获取需要删除黑名单交易的交易数据
//	@param chainId
//	@param getDBResult
//	@param deleteBlackTxIds 删除交易黑名单的交易列表
//	@return TaskSync
func getDeleteBlackTxListTask(chainId string, getDBResult *GetDBResult, deleteBlackTxIds []string) TaskSync {
	return TaskFunc(func() error {
		//删除黑名单交易
		blackTxList, err := dbhandle.BatchQueryBlackTxList(chainId, deleteBlackTxIds)
		if err != nil {
			return err
		}
		getDBResult.DeleteBlackTxList = blackTxList
		return nil
	})
}

// getSubChainCrossListTask
//
//	@Description: 任务：获取子链每条跨链交易数量
//	@param chainId
//	@param getDBResult
//	@param crossSubChainIdMap 子链列表
//	@return TaskSync
func getSubChainCrossListTask(chainId string, getDBResult *GetDBResult,
	crossSubChainIdMap map[string]map[string]int64) TaskSync {
	return TaskFunc(func() error {
		var subChainIds []string
		for subChainId := range crossSubChainIdMap {
			subChainIds = append(subChainIds, subChainId)
		}
		//主子链-获取子链跨链列表交易数据
		crossSubChainCrossList, err := dbhandle.GetCrossSubChainCrossNum(chainId, subChainIds)
		if err != nil {
			return err
		}
		getDBResult.CrossSubChainCross = crossSubChainCrossList
		return nil
	})
}

// getSubChainDBMapTask
//
//	@Description: 任务：获取子链详情
//	@param chainId
//	@param getDBResult
//	@param crossSubChainIdMap 子链列表
//	@return TaskSync
func getSubChainDBMapTask(chainId string, getDBResult *GetDBResult,
	crossSubChainIdMap map[string]map[string]int64) TaskSync {
	return TaskFunc(func() error {
		var subChainIds []string
		for subChainId := range crossSubChainIdMap {
			subChainIds = append(subChainIds, subChainId)
		}

		//主子链-获取子链信息
		crossSubChainDBMap, err := dbhandle.GetCrossSubChainById(chainId, subChainIds)
		if err != nil {
			return err
		}
		getDBResult.CrossSubChainMap = crossSubChainDBMap
		return nil
	})
}

// getAccountByBNSListTask
//
//	@Description: 任务：获取解绑BNS的账户信息
//	@param chainId
//	@param getDBResult
//	@param crossSubChainIdMap 子链列表
//	@return TaskSync
func getAccountByBNSListTask(chainId string, getDBResult *GetDBResult, bnsUnBindDomain []string) TaskSync {
	return TaskFunc(func() error {
		accountBNSList, err := dbhandle.GetAccountByBNSList(chainId, bnsUnBindDomain)
		if err != nil {
			return err
		}
		getDBResult.AccountBNSList = accountBNSList
		return nil
	})
}

// getAccountByDIDListTask
//
//	@Description: 任务：获取解绑DID的账户信息详情
//	@param chainId
//	@param getDBResult
//	@param crossSubChainIdMap 子链列表
//	@return TaskSync
func getAccountByDIDListTask(chainId string, getDBResult *GetDBResult, didUnBindList []string) TaskSync {
	return TaskFunc(func() error {
		accountDIDList, err := dbhandle.GetAccountByDIDList(chainId, didUnBindList)
		if err != nil {
			return err
		}
		getDBResult.AccountDIDList = accountDIDList
		return nil
	})
}

// getAccountMapTask
//
//	@Description: 根据转账地址，BNS，DID绑定解绑地址，user地址获取账户新增
//	@param chainId
//	@param getDBResult 存储DB数据
//	@param topicEventResult 合约事件解析转账数据，BNS，DID
//	@param userInfoMap user信息
//	@return TaskSync
func getAccountMapTask(chainId string, getDBResult *GetDBResult, topicEventResult *TopicEventResult,
	userInfoMap map[string]*db.User) TaskSync {
	return TaskFunc(func() error {
		//BNS账户
		for _, event := range topicEventResult.BNSBindEventData {
			topicEventResult.OwnerAdders = append(topicEventResult.OwnerAdders, event.Value)
		}
		//did账户
		for _, didAccounts := range topicEventResult.DIDAccount {
			topicEventResult.OwnerAdders = append(topicEventResult.OwnerAdders, didAccounts...)
		}
		//user账户
		for _, user := range userInfoMap {
			topicEventResult.OwnerAdders = append(topicEventResult.OwnerAdders, user.UserAddr)
		}

		//获取账户信息,走缓存
		accountMap, err := dbhandle.QueryAccountExists(chainId, topicEventResult.OwnerAdders)
		if err != nil {
			return err
		}

		//获取解绑用户
		unBindBNSList, err := dbhandle.GetAccountByBNSList(chainId, topicEventResult.BNSUnBindDomain)
		if err != nil {
			return err
		}

		//获取解绑用户
		unBindDIDList, err := dbhandle.GetAccountByDIDList(chainId, topicEventResult.DIDUnBindList)
		if err != nil {
			return err
		}

		for _, account := range unBindBNSList {
			accountMap[account.Address] = account
		}
		for _, account := range unBindDIDList {
			accountMap[account.Address] = account
		}

		getDBResult.AccountDBMap = accountMap
		return nil
	})
}

// getIDAContractTask
//
//	@Description: 任务：非同质化合约数据
//	@param chainId
//	@param getDBResult
//	@param contractMap 本次涉及全部合约
//	@return TaskSync
func getIDAContractTask(chainId string, getDBResult *GetDBResult, idaEventData *IDAEventData) TaskSync {
	return TaskFunc(func() error {
		if idaEventData == nil {
			return nil
		}

		idaContractAddr := make([]string, 0)
		for address := range idaEventData.IDACreatedMap {
			idaContractAddr = append(idaContractAddr, address)
		}

		for address := range idaEventData.IDADeletedCodeMap {
			idaContractAddr = append(idaContractAddr, address)
		}

		if len(idaContractAddr) == 0 {
			return nil
		}

		//非同质化合约
		idaContractMap, err := dbhandle.GetIDAContractMapByAddrs(chainId, idaContractAddr)
		if err != nil {
			return err
		}
		getDBResult.IDAContractMap = idaContractMap
		return nil
	})
}

func getIDAAssetDetailsTask(chainId string, getDBResult *GetDBResult, idaEventData *IDAEventData) TaskSync {
	return TaskFunc(func() error {
		if idaEventData == nil {
			return nil
		}

		assetCodes := make([]string, 0)
		for assetCode := range idaEventData.IDAUpdatedMap {
			assetCodes = append(assetCodes, assetCode)
		}

		for _, codes := range idaEventData.IDADeletedCodeMap {
			assetCodes = append(assetCodes, codes...)
		}

		if len(assetCodes) == 0 {
			return nil
		}

		//非同质化合约
		idaAssetMap, err := dbhandle.GetIDAAssetDetailMapByCodes(chainId, assetCodes)
		if err != nil {
			return err
		}

		getDBResult.IDAAssetDetailMap = idaAssetMap
		return nil
	})
}
