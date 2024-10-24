/*
Package saveTasks comment： resolver realtime inset DB
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package saveTasks

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	loggers "chainmaker_web/src/logger"
	"sync"

	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"
	"github.com/panjf2000/ants/v2"
)

var (
	log = loggers.GetLogger(loggers.MODULE_WEB)
)

const (
	GoRoutinePoolErr = "new ants pool error: "
)

// SaveTransactionsToDB
//
//	@Description: 存储交易数据
//	@param chainId
//	@param transactions 交易列表
//	@return error
func SaveTransactionsToDB(chainId string, transactions map[string]*db.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}
	var (
		goRoutinePool *ants.Pool
		err           error
	)

	// 将交易分割为大小为50的批次
	batches := batchTransactions(transactions)
	errChan := make(chan error, config.MaxDBPoolSize)
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	var wg sync.WaitGroup
	for _, batch := range batches {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(insertList []*db.Transaction) func() {
			return func() {
				defer wg.Done()
				//插入数据
				err = dbhandle.InsertTransactions(chainId, insertList)
				if err != nil {
					errChan <- err
				}

			}
		}(batch))
		if errSub != nil {
			log.Error("InsertTransactions submit Failed : " + err.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	return nil
}

// SaveUpgradeContractTxToDB
//
//	@Description: 保存合约交易信息
//	@param chainId
//	@param transactions 交易map
//	@param contractTxId 合约升级交易id
//	@return error
func SaveUpgradeContractTxToDB(chainId string, upgradeContractTxs []*db.UpgradeContractTransaction) error {
	if len(upgradeContractTxs) == 0 {
		return nil
	}

	//插入数据
	err := dbhandle.InsertUpgradeContractTx(chainId, upgradeContractTxs)
	if err != nil {
		return err
	}

	return nil
}

// SaveUserToDB 保存用户信息
func SaveUserToDB(chainId string, users map[string]*db.User) error {
	if len(users) == 0 {
		return nil
	}
	// 将分割为大小为100的批次
	batches := batchUsers(users)
	for _, batch := range batches {
		//插入数据
		err := dbhandle.BatchInsertUser(chainId, batch)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveContractToDB 保存合约
func SaveContractToDB(chainId string, insertContracts, updateContracts []*db.Contract) error {
	for _, contract := range insertContracts {
		//insert
		err := dbhandle.InsertContract(chainId, contract)
		if err != nil {
			log.Errorf("saveContractToDB InsertContract failed, contract:%v", contract)
			return err
		}
	}

	for _, contract := range updateContracts {
		//update
		err := dbhandle.UpdateContract(chainId, contract)
		if err != nil {
			log.Errorf("saveContractToDB UpdateContract failed, contract:%v", contract)
			return err
		}
	}

	return nil
}

// SaveEvidenceContractToDB 保存存证合约
func SaveEvidenceContractToDB(chainId string, contractList []*db.EvidenceContract) error {
	if len(contractList) == 0 {
		return nil
	}

	hashList := make([]string, 0)
	insertList := contractList
	for _, evidence := range contractList {
		hashList = append(hashList, evidence.Hash)
	}

	//数据是否已经存在
	evidenceList, err := dbhandle.GetEvidenceContractByHashLit(chainId, hashList)
	if err != nil {
		return err
	}

	if len(evidenceList) > 0 {
		insertList = make([]*db.EvidenceContract, 0)
		for _, evidence := range contractList {
			if _, ok := evidenceList[evidence.Hash]; !ok {
				insertList = append(insertList, evidence)
			}
		}
	}

	//插入数据
	err = dbhandle.InsertEvidenceContract(chainId, insertList)
	if err != nil {
		log.Errorf("SaveEvidenceContractToDB failed, contract:%v", insertList)
		return err
	}
	return nil
}

// SaveContractEventsToDB
//
//	@Description: 存储合约事件数据
//	@param chainId
//	@param contractEvents 合约事件
//	@return error
func SaveContractEventsToDB(chainId string, contractEvents []*db.ContractEvent) error {
	if len(contractEvents) == 0 {
		return nil
	}
	var (
		goRoutinePool *ants.Pool
		err           error
		wg            sync.WaitGroup
		errChan       = make(chan error, config.MaxDBPoolSize)
	)

	// 将交易分割为大小为SaveBatchSize的批次
	batches := batchContractEvents(chainId, contractEvents)
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	for _, batch := range batches {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(insertList []*db.ContractEvent) func() {
			return func() {
				defer wg.Done()
				//插入数据
				err = dbhandle.InsertContractEvent(chainId, insertList)
				if err != nil {
					errChan <- err
				}

			}
		}(batch))
		if errSub != nil {
			log.Error("InsertContractEvent submit Failed : " + err.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	return nil
}

// SaveGasRecordToDB 保存gasR
func SaveGasRecordToDB(chainId string, gasRecords []*db.GasRecord) error {
	if len(gasRecords) == 0 {
		return nil
	}
	var (
		goRoutinePool *ants.Pool
		err           error
		wg            sync.WaitGroup
		errChan       = make(chan error, config.MaxDBPoolSize)
	)

	// 将交易分割为大小为10的批次
	batches := batchGasRecords(gasRecords)
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	for _, batch := range batches {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(insertList []*db.GasRecord) func() {
			return func() {
				defer wg.Done()
				//插入数据
				err = dbhandle.InsertGasRecord(chainId, insertList)
				if err != nil {
					errChan <- err
				}

			}
		}(batch))
		if errSub != nil {
			log.Error("InsertGasRecord submit Failed : " + err.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	return nil
}

// SaveChainConfig 更新链配置
func SaveChainConfig(chainId string, chainConfigList []*pbConfig.ChainConfig) error {
	if len(chainConfigList) == 0 {
		return nil
	}

	for _, chainConfig := range chainConfigList {
		if chainConfig == nil || chainConfig.ChainId == "" {
			continue
		}

		err := dbhandle.UpdateChainInfoByConfig(chainId, chainConfig)
		if err != nil {
			log.Errorf("SaveChainConfig UpdateChainInfoByConfig failed chainConfig:%v ", chainConfig)
			return err
		}
	}

	return nil
}

func SaveStandardContractToDB(chainId string, fungibleContract []*db.FungibleContract,
	nonFungibleContract []*db.NonFungibleContract, idaContracts []*db.IDAContract) error {
	var err error
	err = dbhandle.InsertFungibleContract(chainId, fungibleContract)
	if err != nil {
		log.Errorf("SaveStandardContractToDB err:%v contract:%v", err, fungibleContract)
		return err
	}

	err = dbhandle.InsertNonFungibleContract(chainId, nonFungibleContract)
	if err != nil {
		log.Errorf("SaveStandardContractToDB err:%v contract:%v", err, nonFungibleContract)
		return err
	}

	err = dbhandle.InsertIDAContract(chainId, idaContracts)
	if err != nil {
		log.Errorf("SaveStandardContractToDB err:%v contract:%v", err, idaContracts)
		return err
	}

	return nil
}

// SaveAccountToDB SaveAccount
func SaveAccountToDB(chainId string, accountResult *db.UpdateAccountResult) error {
	insertAccounts := accountResult.InsertAccount
	updateAccounts := accountResult.UpdateAccount
	//插入账户
	err := dbhandle.InsertAccount(chainId, insertAccounts)
	if err != nil {
		return err
	}

	//更新账户
	for _, account := range updateAccounts {
		err = dbhandle.UpdateAccount(chainId, account)
		if err != nil {
			return err
		}
	}
	return err
}
