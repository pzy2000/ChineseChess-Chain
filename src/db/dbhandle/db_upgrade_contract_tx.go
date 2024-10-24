/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"fmt"
)

// GetUpgradeContractTxList
//
//	@Description: 获取版本更新交易列表
//	@param offset
//	@param limit
//	@param chainId
//	@param contractName
//	@param contractAddr
//	@param senders
//	@param runtimeType
//	@param status
//	@return []*db.UpgradeContractTransaction
//	@return int64
//	@return error
func GetUpgradeContractTxList(offset int, limit int, chainId string, contractName, contractAddr string,
	senders []string, runtimeType string, status int, startTime, endTime int64) (
	[]*db.UpgradeContractTransaction, int64, error) {
	var count int64
	txList := make([]*db.UpgradeContractTransaction, 0)
	where := map[string]interface{}{}
	whereIn := map[string]interface{}{}
	if contractName != "" {
		where["contractNameBak"] = contractName
	}
	if contractAddr != "" {
		where["contractAddr"] = contractAddr
	}
	if len(senders) > 0 {
		whereIn["sender"] = senders
	}

	if runtimeType != "" {
		where["contractRuntimeType"] = runtimeType
	}
	if status != -1 {
		where["contractResultCode"] = status
	}
	tableName := db.GetTableName(chainId, db.TableContractUpgradeTransaction)
	selectFile := &SelectFile{
		Where:     where,
		WhereIn:   whereIn,
		StartTime: startTime,
		EndTime:   endTime,
	}
	query := BuildParamsQuery(tableName, selectFile)
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("GetUpgradeContractTxList err, cause : %s", err.Error())
	}

	query = query.Order("timestamp desc")
	query = query.Offset(offset * limit).Limit(limit)
	err = query.Find(&txList).Error
	if err != nil {
		return nil, 0, fmt.Errorf("GetUpgradeContractTxList err, cause : %s", err.Error())
	}

	return txList, count, nil
}

// InsertUpgradeContractTx 新增或者更新交易
func InsertUpgradeContractTx(chainId string, transactions []*db.UpgradeContractTransaction) error {
	if len(transactions) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableContractUpgradeTransaction)
	return CreateInBatchesData(tableName, transactions)
}

// UpdateUpgradeContractName 更新合约名称
func UpdateUpgradeContractName(chainId string, contract *db.Contract) error {
	if chainId == "" || contract == nil || contract.Addr == "" {
		return nil
	}

	where := map[string]interface{}{
		"contractAddr": contract.Addr,
	}
	params := map[string]interface{}{
		"contractName":    contract.Name,
		"contractNameBak": contract.NameBak,
	}
	tableName := db.GetTableName(chainId, db.TableContractUpgradeTransaction)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}
