/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
)

// InsertContractEvent 批量保存event
func InsertContractEvent(chainId string, contractEvents []*db.ContractEvent) error {
	if len(contractEvents) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableContractEvent)
	return CreateInBatchesData(tableName, contractEvents)
}

// GetEventDataByTxIds 根据ID获取event
func GetEventDataByTxIds(chainId string, txIds []string) ([]*db.ContractEvent, error) {
	eventList := make([]*db.ContractEvent, 0)
	if len(txIds) == 0 {
		return eventList, nil
	}
	tableName := db.GetTableName(chainId, db.TableContractEvent)
	err := db.GormDB.Table(tableName).Where("txId in ?", txIds).Find(&eventList).Error
	return eventList, err
}

// GetEventDataByIds
//
//	@Description: 根据主键id获取数据
//	@param chainId
//	@param ids 主键id
//	@return []*db.ContractEvent
//	@return error
func GetEventDataByIds(chainId string, ids []string) ([]*db.ContractEvent, error) {
	eventList := make([]*db.ContractEvent, 0)
	if len(ids) == 0 {
		return eventList, nil
	}
	tableName := db.GetTableName(chainId, db.TableContractEvent)
	err := db.GormDB.Table(tableName).Where("id in ?", ids).Find(&eventList).Error
	return eventList, err
}

// GetEventListCount
//
//	@Description: 获取合约事件总数量
//	@param offset
//	@param limit
//	@param chainId
//	@param contractName
//	@param contractAddr
//	@param txId
//	@return []*db.ContractEvent
//	@return int64
//	@return error
func GetEventListCount(chainId, contractName, contractAddr, txId string) (int64, error) {
	where := map[string]interface{}{}
	if contractName != "" {
		where["contractNameBak"] = contractName
	}
	if contractAddr != "" {
		where["contractAddr"] = contractAddr
	}
	if txId != "" {
		where["txId"] = txId
	}

	//获取缓存数据
	totalCount, err := GetContractEventCountCache(chainId, where)
	if err == nil && totalCount != 0 {
		return totalCount, nil
	}

	tableName := db.GetTableName(chainId, db.TableContractEvent)
	query := db.GormDB.Table(tableName).Where(where)
	err = query.Count(&totalCount).Error
	if err != nil {
		return 0, err
	}

	//设置缓存
	SetContractEventCountCache(chainId, where, totalCount)
	return totalCount, nil
}

// GetEventIdList
//
//	@Description: 获取event主键列表
//	@param offset
//	@param limit
//	@param chainId
//	@param contractName
//	@param contractAddr
//	@param txId
//	@return []*db.ContractEvent
//	@return int64
//	@return error
func GetEventIdList(offset, limit int, chainId, contractName, contractAddr, txId string) ([]string, error) {
	contractEventIds := make([]string, 0)
	where := map[string]interface{}{}
	if contractName != "" {
		where["contractNameBak"] = contractName
	}
	if contractAddr != "" {
		where["contractAddr"] = contractAddr
	}
	if txId != "" {
		where["txId"] = txId
	}
	tableName := db.GetTableName(chainId, db.TableContractEvent)
	query := db.GormDB.Table(tableName).Select("id").Where(where)
	query = query.Order("timestamp desc").Offset(offset * limit).Limit(limit)
	err := query.Find(&contractEventIds).Error
	if err != nil {
		return contractEventIds, err
	}
	return contractEventIds, nil
}

//// GetEventList 获取事件信息
//func GetEventList(offset, limit int, chainId, contractName, contractAddr, txId string) ([]*db.ContractEvent,
//	int64, error) {
//	var count int64
//	contractEvents := make([]*db.ContractEvent, 0)
//	if chainId == "" {
//		return contractEvents, 0, db.ErrTableParams
//	}
//	where := map[string]interface{}{}
//	if contractName != "" {
//		where["contractNameBak"] = contractName
//	}
//	if contractAddr != "" {
//		where["contractAddr"] = contractAddr
//	}
//	if txId != "" {
//		where["txId"] = txId
//	}
//	tableName := db.GetTableName(chainId, db.TableContractEvent)
//	query := db.GormDB.Table(tableName).Where(where)
//	err := query.Count(&count).Error
//	if err != nil {
//		return contractEvents, 0, err
//	}
//	query = query.Order("timestamp desc").Offset(offset * limit).Limit(limit)
//	err = query.Find(&contractEvents).Error
//	if err != nil {
//		return contractEvents, 0, err
//	}
//	return contractEvents, count, nil
//}

// UpdateContractEventBak 更新事件敏感词
func UpdateContractEventBak(chainId string, contractEvent *db.ContractEvent) error {
	if chainId == "" || contractEvent == nil {
		return nil
	}

	where := map[string]interface{}{
		"txId":       contractEvent.TxId,
		"eventIndex": contractEvent.EventIndex,
	}
	params := map[string]interface{}{
		"topic":        contractEvent.Topic,
		"topicBak":     contractEvent.TopicBak,
		"eventData":    contractEvent.EventData,
		"eventDataBak": contractEvent.EventDataBak,
	}
	tableName := db.GetTableName(chainId, db.TableContractEvent)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateContractEventSensitiveWord 更新合约名称
func UpdateContractEventSensitiveWord(chainId string, contract *db.Contract) error {
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
	tableName := db.GetTableName(chainId, db.TableContractEvent)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}
