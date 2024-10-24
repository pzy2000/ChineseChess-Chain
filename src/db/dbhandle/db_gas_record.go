/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
)

const (
	//BusinessTypeRecharge gas充值
	BusinessTypeRecharge = 1
	//BusinessTypeConsume gas消费
	BusinessTypeConsume = 2
)

// InsertGasRecord 批量保存gasRecords
func InsertGasRecord(chainId string, gasRecords []*db.GasRecord) error {
	if len(gasRecords) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableGasRecord)
	return CreateInBatchesData(tableName, gasRecords)
}

// GetGasRecordByTxIds 根据ID获取gas
func GetGasRecordByTxIds(chainId string, txIds []string) ([]*db.GasRecord, error) {
	gasRecords := make([]*db.GasRecord, 0)
	if len(txIds) == 0 {
		return gasRecords, nil
	}
	tableName := db.GetTableName(chainId, db.TableGasRecord)
	err := db.GormDB.Table(tableName).Where("txId in ?", txIds).Find(&gasRecords).Error
	return gasRecords, err
}

// GetGasRecordList GetGasRecordList
func GetGasRecordList(offset int, limit int, chainId string, addrList []string, startTime, endTime int64,
	businessType int) ([]*db.GasRecord, int64, error) {
	var count int64
	gasRecordList := make([]*db.GasRecord, 0)
	if chainId == "" {
		return gasRecordList, 0, db.ErrTableParams
	}
	where := map[string]interface{}{}
	whereIn := map[string]interface{}{}
	if businessType == BusinessTypeRecharge ||
		businessType == BusinessTypeConsume {
		where["businessType"] = businessType
	}
	if len(addrList) > 0 {
		whereIn["address"] = addrList
	}

	selectFile := &SelectFile{
		Where:     where,
		WhereIn:   whereIn,
		StartTime: startTime,
		EndTime:   endTime,
	}

	tableName := db.GetTableName(chainId, db.TableGasRecord)
	query := BuildParamsQuery(tableName, selectFile)
	err := query.Count(&count).Error
	if err != nil {
		return gasRecordList, 0, err
	}
	err = query.Order("timestamp desc").Offset(offset * limit).Limit(limit).Find(&gasRecordList).Error
	if err != nil {
		return gasRecordList, 0, err
	}

	return gasRecordList, count, nil
}
