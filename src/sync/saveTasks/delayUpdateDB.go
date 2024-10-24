/*
Package saveTasks comment： resolver delay update DB
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package saveTasks

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"sync"

	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

// -------异步更新--------

// InsertFungibleTransferToDB 保存流转信息
func InsertFungibleTransferToDB(chainId string, transferList []*db.FungibleTransfer) error {
	if len(transferList) == 0 {
		return nil
	}

	var (
		goRoutinePool *ants.Pool
		err           error
		wg            sync.WaitGroup
		errChan       = make(chan error, config.MaxDBPoolSize)
	)

	// 将交易分割为大小为10的批次
	batches := batchFungibleTransfers(transferList)
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	for _, batch := range batches {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(insertList []*db.FungibleTransfer) func() {
			return func() {
				defer wg.Done()
				//插入数据
				err = dbhandle.InsertFungibleTransfer(chainId, insertList)
				if err != nil {
					errChan <- err
				}

			}
		}(batch))
		if errSub != nil {
			log.Error("InsertFungibleTransfer submit Failed : " + err.Error())
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	return nil
}

// InsertNonFungibleTransferToDB 保存流转信息
func InsertNonFungibleTransferToDB(chainId string, transferList []*db.NonFungibleTransfer) error {
	if len(transferList) == 0 {
		return nil
	}

	var (
		goRoutinePool *ants.Pool
		err           error
		wg            sync.WaitGroup
		errChan       = make(chan error, config.MaxDBPoolSize)
	)

	// 将交易分割为大小为10的批次
	batches := batchNonFungibleTransfers(transferList)
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	for _, batch := range batches {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(insertList []*db.NonFungibleTransfer) func() {
			return func() {
				defer wg.Done()
				//插入数据
				err = dbhandle.InsertNonFungibleTransfer(chainId, insertList)
				if err != nil {
					errChan <- err
				}

			}
		}(batch))
		if errSub != nil {
			log.Error("InsertNonFungibleTransfer submit Failed : " + err.Error())
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	return nil
}

// SaveNonFungibleToken 保存Token
func SaveNonFungibleToken(chainId string, tokenResult *db.TokenResult) error {
	if len(tokenResult.InsertUpdateToken) > 0 {
		var tokenIds, contractAddrs []string
		for _, token := range tokenResult.InsertUpdateToken {
			tokenIds = append(tokenIds, token.TokenId)
			contractAddrs = append(contractAddrs, token.ContractAddr)
		}
		tokenMap, err := dbhandle.SelectTokenByID(chainId, tokenIds, contractAddrs)
		if err != nil {
			return err
		}

		insertList := make([]*db.NonFungibleToken, 0)
		updateList := make([]*db.NonFungibleToken, 0)
		for _, token := range tokenResult.InsertUpdateToken {
			key := token.TokenId + "_" + token.ContractAddr
			if value, ok := tokenMap[key]; ok {
				//value.AddrType = token.AddrType
				value.OwnerAddr = token.OwnerAddr
				updateList = append(updateList, value)
			} else {
				insertList = append(insertList, token)
			}
		}
		err = InsertNonFungibleTokenConcurrent(chainId, insertList)
		if err != nil {
			return err
		}
		err = UpdateNonFungibleToken(chainId, updateList)
		if err != nil {
			return err
		}
	}

	if len(tokenResult.DeleteToken) > 0 {
		err := dbhandle.DeleteNonFungibleToken(chainId, tokenResult.DeleteToken)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertNonFungibleTokenConcurrent(chainId string, tokenList []*db.NonFungibleToken) error {
	var (
		goRoutinePool *ants.Pool
		err           error
	)
	if len(tokenList) == 0 {
		return nil
	}
	// 将交易分割为大小为10的批次
	batches := batchNonFungibleToken(tokenList)
	errChan := make(chan error, config.MaxDBPoolSize)
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	var wg sync.WaitGroup
	for _, batch := range batches {
		wg.Add(1)
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(insertList []*db.NonFungibleToken) func() {
			return func() {
				defer wg.Done()
				//插入数据
				err = dbhandle.InsertNonFungibleToken(chainId, insertList)
				if err != nil {
					errChan <- err
				}

			}
		}(batch))
		if errSub != nil {
			log.Error("InsertNonFungibleTokenConcurrent submit Failed : " + err.Error())
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	return nil
}

// UpdateNonFungibleToken 更新token
func UpdateNonFungibleToken(chainId string, tokenList []*db.NonFungibleToken) error {
	var (
		goRoutinePool *ants.Pool
		err           error
	)
	if len(tokenList) == 0 {
		return nil
	}

	errChan := make(chan error, config.MaxDBPoolSize)
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	var wg sync.WaitGroup
	for _, token := range tokenList {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(token *db.NonFungibleToken) func() {
			return func() {
				defer wg.Done()
				//更新数据
				err = dbhandle.UpdateNonFungibleToken(chainId, token)
				if err != nil {
					errChan <- err
				}

			}
		}(token))
		if errSub != nil {
			log.Error("UpdateNonFungibleToken submit Failed : " + err.Error())
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	return nil
}

// SavePositionToDB 保存持仓数据
func SavePositionToDB(chainId string, dbPositionOperates *db.BlockPosition) error {
	insertFungiblePosition := dbPositionOperates.InsertFungiblePosition
	updateFungiblePosition := dbPositionOperates.UpdateFungiblePosition
	deleteFungiblePosition := dbPositionOperates.DeleteFungiblePosition
	insertNonFungible := dbPositionOperates.InsertNonFungible
	updateNonFungible := dbPositionOperates.UpdateNonFungible
	deleteNonFungible := dbPositionOperates.DeleteNonFungible

	if len(insertFungiblePosition) > 0 {
		err := dbhandle.InsertFungiblePosition(chainId, insertFungiblePosition)
		if err != nil {
			log.Errorf("savePositionDB InsertFungiblePosition failed , ", zap.Any("position", insertFungiblePosition))
			return err
		}
	}
	if len(updateFungiblePosition) > 0 {
		err := dbhandle.UpdateFungiblePosition(chainId, updateFungiblePosition)
		if err != nil {
			log.Errorf("savePositionDB UpdateFungiblePosition failed , ",
				zap.Any("position", updateFungiblePosition))
			return err
		}
	}
	if len(deleteFungiblePosition) > 0 {
		err := dbhandle.DeleteFungiblePosition(chainId, deleteFungiblePosition)
		if err != nil {
			log.Errorf("savePositionDB DeleteFungiblePosition failed , ",
				zap.Any("position", updateFungiblePosition))
			return err
		}
	}
	if len(insertNonFungible) > 0 {
		err := dbhandle.InsertNonFungiblePosition(chainId, insertNonFungible)
		if err != nil {
			log.Errorf("savePositionDB insertNonFungiblePosition failed , ",
				zap.Any("position", insertNonFungible))
			return err
		}
	}
	if len(updateNonFungible) > 0 {
		err := dbhandle.UpdateNonFungiblePosition(chainId, updateNonFungible)
		if err != nil {
			log.Errorf("savePositionDB UpdateNonFungiblePosition failed , ",
				zap.Any("position", updateNonFungible))
			return err
		}
	}
	if len(deleteNonFungible) > 0 {
		err := dbhandle.DeleteNonFungiblePosition(chainId, deleteNonFungible)
		if err != nil {
			log.Errorf("savePositionDB UpdateNonFungiblePosition failed , ",
				zap.Any("position", updateNonFungible))
			return err
		}
	}

	return nil
}

// InsertGasToDB 保存gas
func InsertGasToDB(chainId string, insertGas []*db.Gas) error {
	if len(insertGas) < 0 {
		return nil
	}

	err := dbhandle.InsertBatchGas(chainId, insertGas)
	if err != nil {
		log.Errorf("SaveGasToDB InsertBatchGas failed , ", zap.Any("InsertBatchGas", insertGas))
		return err
	}
	return nil
}

// UpdateGasToDB 更新gas
func UpdateGasToDB(chainId string, updateGas []*db.Gas) error {
	if len(updateGas) < 0 {
		return nil
	}

	for _, gasInfo := range updateGas {
		err := dbhandle.UpdateGas(chainId, gasInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateTxBlackToDB 更新交易黑名单
func UpdateTxBlackToDB(chainId string, txBlockList *db.UpdateTxBlack) error {
	//需要添加黑名单
	err := dbhandle.InsertBlackTransactions(chainId, txBlockList.AddTxBlack)
	if err != nil {
		log.Errorf("InsertBlackTransactions chainId-%s failed, err:%v, AddBlack:%v ",
			chainId, err, txBlockList.AddTxBlack)
		return err
	}

	//删除黑名单
	err = dbhandle.DeleteBlackTransaction(chainId, txBlockList.DeleteTxBlack)
	if err != nil {
		log.Errorf("DeleteBlackTransactions chainId-%s failed, err:%v, AddBlack:%v ",
			chainId, err, txBlockList.DeleteTxBlack)
		return err
	}
	return nil
}

// UpdateContract 更新合约数据
func UpdateContract(chainId string, identityContract []*db.IdentityContract) error {
	err := dbhandle.InsertIdentityContract(chainId, identityContract)
	if err != nil {
		log.Errorf("SaveContractResult InsertIdentityContract err:%v contract:%v", err, identityContract)
		return err
	}
	return nil
}

// UpdateContractTxNum 更新合约交易量数据
func UpdateContractTxNum(chainId string, updateContracts []*db.Contract) error {
	var (
		goRoutinePool *ants.Pool
		err           error
	)
	if len(updateContracts) == 0 {
		return nil
	}
	errChan := make(chan error, config.MaxDBPoolSize)
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	var wg sync.WaitGroup
	for _, contract := range updateContracts {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		errSub := goRoutinePool.Submit(func(contract *db.Contract) func() {
			return func() {
				defer wg.Done()
				updateErr := dbhandle.UpdateContractTxNum(chainId, contract)
				if updateErr != nil {
					errChan <- updateErr
				}

			}
		}(contract))
		if errSub != nil {
			log.Error(" submit Failed : " + err.Error())
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	return nil
}

// SaveFungibleContractResult 更新合约数据
func SaveFungibleContractResult(chainId string, contractResult *db.GetContractResult) error {
	var err error
	updateFungibleContract := contractResult.UpdateFungibleContract
	updateNonFungible := contractResult.UpdateNonFungible

	for _, contract := range updateFungibleContract {
		//更新数据
		err = dbhandle.UpdateFungibleContract(chainId, contract)
		if err != nil {
			return err
		}
	}

	for _, contract := range updateNonFungible {
		//更新数据
		err = dbhandle.UpdateNonFungibleContract(chainId, contract)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateBlockStatusToDB 更新区块状态
func UpdateBlockStatusToDB(chainId string, blockHeightList []int64) error {
	if len(blockHeightList) == 0 {
		return nil
	}
	for _, blockHeight := range blockHeightList {
		err := dbhandle.UpdateBlockUpdateStatus(chainId, blockHeight, dbhandle.DelayUpdateSuccess)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateIDAContract 更新IDA合约
func UpdateIDAContract(chainId string, updateIdaContract map[string]*db.IDAContract) error {
	if len(updateIdaContract) == 0 {
		return nil
	}
	for contractAddr, idaContract := range updateIdaContract {
		err := dbhandle.UpdateIDAContractByAddr(chainId, contractAddr, idaContract)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveIDAAssetDataToDB
func SaveIDAAssetDataToDB(chainId string, idaAssetsData *db.IDAAssetsDataDB) error {
	if idaAssetsData == nil {
		return nil
	}

	idaDetails := idaAssetsData.IDAAssetDetail
	idaAttachments := idaAssetsData.IDAAssetAttachment
	idaApis := idaAssetsData.IDAAssetApi
	idaDatas := idaAssetsData.IDAAssetData
	err := dbhandle.InsertIDADetail(chainId, idaDetails)
	if err != nil {
		return err
	}

	err = dbhandle.InsertIDAAttachments(chainId, idaAttachments)
	if err != nil {
		return err
	}
	err = dbhandle.InsertIDAAssetApi(chainId, idaApis)
	if err != nil {
		return err
	}
	err = dbhandle.InsertIDAAssetData(chainId, idaDatas)
	if err != nil {
		return err
	}

	return nil
}

// UpdateIDAAssetDataToDB
func UpdateIDAAssetDataToDB(chainId string, updateAssetsData *db.IDAAssetsUpdateDB) error {
	if updateAssetsData == nil {
		return nil
	}

	updateDetails := updateAssetsData.UpdateAssetDetails
	for _, detail := range updateDetails {
		err := dbhandle.UpdateIDADetailByCode(chainId, detail)
		if err != nil {
			return err
		}
	}

	insertAttachments := updateAssetsData.InsertAttachment
	deleteAttachments := updateAssetsData.DeleteAttachmentCodes
	//需要先删除在插入
	err := dbhandle.DeleteIDAAttachments(chainId, deleteAttachments)
	if err != nil {
		return err
	}
	err = dbhandle.InsertIDAAttachments(chainId, insertAttachments)
	if err != nil {
		return err
	}

	insertApis := updateAssetsData.InsertIDAAssetApi
	deleteApis := updateAssetsData.DeleteAssetApiCodes
	//需要先删除在插入
	err = dbhandle.DeleteIDAApis(chainId, deleteApis)
	if err != nil {
		return err
	}
	err = dbhandle.InsertIDAAssetApi(chainId, insertApis)
	if err != nil {
		return err
	}

	insertDatas := updateAssetsData.InsertIDAAssetData
	deleteDatas := updateAssetsData.DeleteAssetApiCodes
	//需要先删除在插入
	err = dbhandle.DeleteIDADatas(chainId, deleteDatas)
	if err != nil {
		return err
	}
	err = dbhandle.InsertIDAAssetData(chainId, insertDatas)
	if err != nil {
		return err
	}

	return nil
}
