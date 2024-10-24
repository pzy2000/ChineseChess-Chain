/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

const TxSuccess = "SUCCESS"

// InsertTransactions
//
//	@Description: 添加交易数据
//	@param chainId
//	@param transactions
//	@return error
func InsertTransactions(chainId string, transactions []*db.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableTransaction)
	return CreateInBatchesData(tableName, transactions)
}

// GetTxInfoByBlockHeight
//
//	@Description:根据区块高度获取交易数据
//	@param chainId
//	@param blockHeight
//	@return []*db.Transaction
//	@return error
func GetTxInfoByBlockHeight(chainId string, blockHeight []int64) ([]*db.Transaction, error) {
	transactions := make([]*db.Transaction, 0)
	tableName := db.GetTableName(chainId, db.TableTransaction)
	err := db.GormDB.Table(tableName).Where("blockHeight in ?", blockHeight).Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

// InsertBlackTransactions
//
//	@Description: 插入黑名单交易
//	@param chainId
//	@param transactions
//	@return error
func InsertBlackTransactions(chainId string, transactions []*db.BlackTransaction) error {
	if len(transactions) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableBlackTransaction)
	//添加黑名单交易
	err := CreateInBatchesData(tableName, transactions)
	if err != nil {
		return fmt.Errorf("InsertBlackTransactions InsertBatchNew err, cause : %v", err.Error())
	}

	//删除交易数据
	var txIds []string
	for _, txInfo := range transactions {
		txIds = append(txIds, txInfo.TxId)
	}
	err = DeleteTransactionByTxId(chainId, txIds)
	if err != nil {
		return fmt.Errorf("InsertBlackTransactions DeleteTransaction err, cause : %v", err.Error())
	}
	return nil
}

// DeleteTransactionByTxId delete交易
func DeleteTransactionByTxId(chainId string, txIds []string) error {
	if len(txIds) == 0 {
		return nil
	}
	tableName := db.GetTableName(chainId, db.TableTransaction)
	return db.GormDB.Table(tableName).Where("txId in ?", txIds).Delete(&db.Transaction{}).Error
}

// DeleteBlackTransaction delete黑名单交易
func DeleteBlackTransaction(chainId string, transactions []*db.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	//恢复黑名单交易
	tableName := db.GetTableName(chainId, db.TableTransaction)
	err := CreateInBatchesData(tableName, transactions)
	if err != nil {
		return err
	}

	//删除黑名单数据
	var txIds []string
	for _, txInfo := range transactions {
		txIds = append(txIds, txInfo.TxId)
	}
	tableNameBlack := db.GetTableName(chainId, db.TableBlackTransaction)
	return db.GormDB.Table(tableNameBlack).Where("txId in ?", txIds).Delete(&db.Transaction{}).Error
}

// GetTransactionIDList
//
//	@Description: 获取交易id列表
//	@param chainId
//	@param offset
//	@param limit
//	@param startTime
//	@param endTime
//	@param txId
//	@param blockHash
//	@param contractName
//	@param contractAddr
//	@param method
//	@param senders
//	@param userAddrs
//	@return []*db.Transaction
//	@return error
func GetTransactionIDList(chainId, contractName, blockHash string, offset, limit int, startTime, endTime int64,
	txId string, txStatus int, senders, userAddrs []string) ([]string, error) {
	txIdList := make([]string, 0) // 修改这里
	//构造查询条件
	selectFile := GetTransactionSelectFile(txId, contractName, blockHash, startTime, endTime, txStatus, senders, userAddrs)
	tableName := db.GetTableName(chainId, db.TableTransaction)
	query := db.GormDB.Table(tableName).Select("txId")
	query = BuildParamsQueryNew(query, selectFile)
	query = query.Order("blockHeight desc, timestamp desc") // 修改这里
	query = query.Offset(offset * limit).Limit(limit)
	err := query.Find(&txIdList).Error
	if err != nil {
		return nil, fmt.Errorf("GetTransactionList err, cause : %s", err.Error())
	}

	return txIdList, nil
}

// GetBlockTransactionList
//
//	@Description:
//	@param chainId
//	@param txIds
//	@return []*db.BlockTxListResult
//	@return error
func GetBlockTransactionList(chainId string, txIds []string) ([]*db.BlockTxListResult, error) {
	txList := make([]*db.BlockTxListResult, 0)
	if chainId == "" {
		return txList, db.ErrTableParams
	}
	tableName := db.GetTableName(chainId, db.TableTransaction)
	query := db.GormDB.Table(tableName).
		Select("txId, sender, senderOrgId, userAddr, txStatusCode, contractName, contractAddr, timestamp").
		Where("txId in ?", txIds)
	err := query.Find(&txList).Error
	if err != nil {
		return nil, fmt.Errorf("GetTransactionList err, cause : %s", err.Error())
	}

	return txList, nil
}

func GetContractTransactionList(chainId string, txIds []string) ([]*db.ContractTxListResult, error) {
	txList := make([]*db.ContractTxListResult, 0)
	if chainId == "" {
		return txList, db.ErrTableParams
	}
	tableName := db.GetTableName(chainId, db.TableTransaction)
	query := db.GormDB.Table(tableName).
		Select("txId, sender, userAddr, txStatusCode, contractName, contractAddr, "+
			"contractMethod, blockHeight, timestamp").
		Where("txId in ?", txIds)
	err := query.Find(&txList).Error
	if err != nil {
		return nil, fmt.Errorf("GetTransactionList err, cause : %s", err.Error())
	}

	return txList, nil
}

func GetBlockTxIDList(chainId, blockHash string, offset, limit int) ([]string, error) {
	txIdList := make([]string, 0)
	if chainId == "" || blockHash == "" {
		return txIdList, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableTransaction)
	query := db.GormDB.Table(tableName).Select("txId").Where("blockHash = ?", blockHash)
	query = query.Offset(offset * limit).Limit(limit)
	err := query.Find(&txIdList).Error
	if err != nil {
		return nil, fmt.Errorf("GetTransactionList err, cause : %s", err.Error())
	}

	return txIdList, nil
}

func GetContractTxIDList(offset, limit int, chainId, contractAddr, contractMethod string,
	userAddrs []string, txStatus int) ([]string, error) {
	txIdList := make([]string, 0)
	if chainId == "" || contractAddr == "" {
		return txIdList, db.ErrTableParams
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	if contractMethod != "" {
		where["contractMethod"] = contractMethod
	}

	tableName := db.GetTableName(chainId, db.TableTransaction)
	query := db.GormDB.Table(tableName).Select("txId").Where(where)
	if len(userAddrs) > 0 {
		query = query.Where("userAddr in ?", userAddrs)
	}
	if txStatus == 0 {
		query = query.Where("txStatusCode = ?", "SUCCESS")
	} else if txStatus == 1 {
		query = query.Where("txStatusCode != ?", "SUCCESS")
	}

	query = query.Order("blockHeight desc, timestamp desc")
	query = query.Offset(offset * limit).Limit(limit)
	err := query.Find(&txIdList).Error
	if err != nil {
		return nil, fmt.Errorf("GetTransactionList err, cause : %s", err.Error())
	}

	return txIdList, nil
}

func GetContractTxCount(chainId, contractAddr, contractMethod string, userAddrs []string, txStatus int) (int64, error) {
	var txCount int64
	if chainId == "" || contractAddr == "" {
		return 0, db.ErrTableParams
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	if contractMethod != "" {
		where["contractMethod"] = contractMethod
	}

	tableName := db.GetTableName(chainId, db.TableTransaction)
	query := db.GormDB.Table(tableName).Select("txId").Where(where)
	if len(userAddrs) > 0 {
		query = query.Where("userAddr in ?", userAddrs)
	}
	if txStatus == 0 {
		query = query.Where("txStatusCode = ?", TxSuccess)
	} else if txStatus == 1 {
		query = query.Where("txStatusCode != ?", TxSuccess)
	}

	err := query.Count(&txCount).Error
	if err != nil {
		return 0, fmt.Errorf("GetContractTxCount err, cause : %s", err.Error())
	}

	return txCount, nil
}

// BatchQueryTxList
//
//	@Description: 根据txIds查询交易
//	@param chainId
//	@param txIds
//	@return []*db.Transaction
//	@return error
func BatchQueryTxList(chainId string, txIds []string) ([]*db.Transaction, error) {
	transaction := make([]*db.Transaction, 0)
	if len(txIds) == 0 {
		return transaction, nil
	}
	tableName := db.GetTableName(chainId, db.TableTransaction)
	err := db.GormDB.Table(tableName).Where("txId in ?", txIds).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

// GetUserTxIDList
//
//	@Description: 根据账户地址获取交易id列表
//	@param chainId
//	@param userAddr
//	@param offset
//	@param limit
//	@return []string
//	@return error
func GetUserTxIDList(chainId string, userAddrs []string, offset, limit int) ([]string, error) {
	txIdList := make([]string, 0)
	if chainId == "" || len(userAddrs) == 0 {
		return txIdList, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableTransaction)
	query := db.GormDB.Table(tableName).Select("txId").Where("userAddr in ?", userAddrs)
	query = query.Order("blockHeight desc, timestamp desc")
	query = query.Offset(offset * limit).Limit(limit)
	err := query.Find(&txIdList).Error
	if err != nil {
		return nil, fmt.Errorf("GetTransactionList err, cause : %s", err.Error())
	}

	return txIdList, nil
}

// GetTransactionListCount
//
//	@Description: 获取交易列表交易总数
//	@param chainId
//	@param txId
//	@param blockHash
//	@param contractName
//	@param contractAddr
//	@param method
//	@param startTime
//	@param endTime
//	@param senders
//	@param userAddrs
//	@return int64
//	@return error
func GetTransactionListCount(chainId, txId, contractName, blockHash string, startTime, endTime int64, txStatus int,
	senders, userAddrs []string) (
	int64, error) {
	//构造查询条件
	selectFile := GetTransactionSelectFile(txId, contractName, blockHash, startTime, endTime, txStatus, senders, userAddrs)
	//获取count缓存数据
	count, err := GetTransactionListCountCache(chainId, selectFile)
	if err == nil && count != 0 {
		return count, nil
	}

	tableName := db.GetTableName(chainId, db.TableTransaction)
	query := BuildParamsQuery(tableName, selectFile)
	err = query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	//设置缓存数据
	SetTransactionListCountCache(chainId, selectFile, count)
	return count, nil
}

// GetTransactionSelectFile
//
//	@Description: 构造交易列表查询条件
//	@param chainId
//	@param txId
//	@param blockHash
//	@param contractName
//	@param contractAddr
//	@param method
//	@param startTime
//	@param endTime
//	@param senders
//	@param userAddrs
//	@return *SelectFile
func GetTransactionSelectFile(txId, contractName, blockHash string, startTime, endTime int64, txStatus int,
	senders, userAddrs []string) *SelectFile {
	where := map[string]interface{}{}
	notWhere := map[string]interface{}{}
	whereIn := map[string]interface{}{}
	if txId != "" {
		where["txId"] = txId
	}
	if contractName != "" {
		where["contractName"] = contractName
	}
	if blockHash != "" {
		where["blockHash"] = blockHash
	}
	if txStatus == 0 {
		where["txStatusCode"] = TxSuccess
	} else if txStatus == 1 {
		notWhere["txStatusCode"] = TxSuccess
	}

	if len(senders) == 1 {
		where["sender"] = senders[0]
	} else if len(senders) > 1 {
		whereIn["sender"] = senders
	}
	if len(userAddrs) == 1 {
		where["userAddr"] = userAddrs[0]
	} else if len(userAddrs) > 1 {
		whereIn["userAddr"] = userAddrs
	}

	selectFile := &SelectFile{
		Where:     where,
		WhereIn:   whereIn,
		NotWhere:  notWhere,
		StartTime: startTime,
		EndTime:   endTime,
	}

	return selectFile
}

// GetSafeWordTransactionList
//
//	@Description: 获取敏感词交易
//	@param chainId
//	@param startTime
//	@param endTime
//	@return []*db.Transaction
//	@return error
func GetSafeWordTransactionList(chainId string, startTime int64, endTime int64) ([]*db.Transaction, error) {
	txList := make([]*db.Transaction, 0)
	tableName := db.GetTableName(chainId, db.TableTransaction)
	condition := gorm.Expr("ContractParametersBak <> ? OR contractResultBak <> ?", "", "")
	// 更新满足条件的记录
	query := db.GormDB.Table(tableName).Where(condition)
	// 添加时间范围条件
	if startTime > 0 && endTime > 0 {
		query = query.Where("timestamp BETWEEN ? AND ?", startTime, endTime)
	}
	query = query.Limit(10)
	err := query.Find(&txList).Error
	if err != nil {
		return txList, err
	}
	return txList, nil
}

// GetTransactionNumByRange 获取交易数量
func GetTransactionNumByRange(chainId string, userAddr string, startTime int64, endTime int64) (int64, error) {
	if userAddr == "" || chainId == "" {
		return 0, nil
	}

	var count int64
	where := map[string]interface{}{
		"userAddr": userAddr,
	}
	selectFile := &SelectFile{
		Where:     where,
		StartTime: startTime,
		EndTime:   endTime,
	}
	tableName := db.GetTableName(chainId, db.TableTransaction)
	query := BuildParamsQuery(tableName, selectFile)
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTxListNumByRange 获取指定时间的交易列表
func GetTxListNumByRange(chainId string, startTime, endTime, interval int64) (map[int64]int64, error) {
	tableName := db.GetTableName(chainId, db.TableTransaction)
	// 创建一个包含所有时间段的初始映射
	txMap := make(map[int64]int64)
	for t := endTime; t > startTime; t -= interval {
		// 将时间戳转换为整数
		txKey := t / interval * interval
		txMap[txKey] = 0
	}

	type Result struct {
		Day      string
		DocCount int64
	}
	var results []Result
	err := db.GormDB.Table(tableName).
		Select("FLOOR(timestamp / ?) * ? as day, COUNT(*) as doc_count", interval, interval).
		Where("timestamp >= ? AND timestamp <= ?", startTime, endTime).
		Group("day").
		Find(&results).Error

	if err != nil {
		return txMap, fmt.Errorf("count day txs err, cause : %s", err.Error())
	}

	for _, result := range results {
		dayFloat, err := strconv.ParseFloat(result.Day, 64)
		if err != nil {
			return txMap, fmt.Errorf("count day txs err, cause : %s", err.Error())
		}
		txMap[int64(dayFloat)] = result.DocCount
	}
	return txMap, nil
}

// GetTxCountByRange 获取指定时间内的交易数量
func GetTxCountByRange(chainId string, startTime, endTime int64) (int64, error) {
	var totalCount int64
	tableName := db.GetTableName(chainId, db.TableTransaction)
	query := db.GormDB.Table(tableName)
	// 添加时间范围条件
	if startTime > 0 && endTime > 0 {
		query = query.Where("timestamp BETWEEN ? AND ?", startTime, endTime)
	}
	err := query.Count(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

// GetTxNumByContractName 根据合约名称获取交易数量
func GetTxNumByContractName(chainId, contractName string) (int64, error) {
	if chainId == "" || contractName == "" {
		return 0, nil
	}
	var count int64
	where := map[string]interface{}{
		"contractNameBak": contractName,
	}
	tableName := db.GetTableName(chainId, db.TableTransaction)
	err := db.GormDB.Table(tableName).Where(where).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count tx err, cause : %s", err.Error())
	}
	return count, nil
}

// GetTransactionByTxId 根据交易id获取交易信息
func GetTransactionByTxId(txId, chainId string) (*db.Transaction, error) {
	if txId == "" || chainId == "" {
		return nil, db.ErrTableParams
	}

	transaction := &db.Transaction{}
	where := map[string]interface{}{
		"txId": txId,
	}
	tableName := db.GetTableName(chainId, db.TableTransaction)
	err := db.GormDB.Table(tableName).Where(where).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return transaction, nil
}

// GetBlackTxInfoByTxId 获取黑名单交易
func GetBlackTxInfoByTxId(chainId, txId string) (*db.BlackTransaction, error) {
	if txId == "" || chainId == "" {
		return nil, db.ErrTableParams
	}

	transaction := &db.BlackTransaction{}
	where := map[string]interface{}{
		"txId": txId,
	}
	tableName := db.GetTableName(chainId, db.TableBlackTransaction)
	err := db.GormDB.Table(tableName).Where(where).First(&transaction).Error
	if err != nil || transaction == nil {
		return nil, err
	}

	return transaction, nil
}

// BatchQueryBlackTxList 根据txIds查询黑名单交易
func BatchQueryBlackTxList(chainId string, txIds []string) ([]*db.BlackTransaction, error) {
	transaction := make([]*db.BlackTransaction, 0)
	if len(txIds) == 0 {
		return transaction, nil
	}

	tableName := db.GetTableName(chainId, db.TableBlackTransaction)
	err := db.GormDB.Table(tableName).Where("txId in ?", txIds).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

// GetLatestTxList
//
//	@Description: 获取最后10个交易列表
//	@param chainId
//	@return []*db.Transaction 交易列表
//	@return error
func GetLatestTxList(chainId string) ([]*db.Transaction, error) {
	transactions := make([]*db.Transaction, 0)
	tableName := db.GetTableName(chainId, db.TableTransaction)
	selectFile := &SelectFile{}
	query := BuildParamsQuery(tableName, selectFile)
	query = query.Order("blockHeight desc, timestamp desc")
	query = query.Limit(10)
	err := query.Find(&transactions).Error
	if err != nil {
		return nil, fmt.Errorf("GetLatestTxList By chainId err, cause : %s", err.Error())
	}

	// 缓存交易信息
	SetLatestTxListCache(chainId, transactions)
	return transactions, nil
}

// UpdateTransactionBak
//
//	@Description: 更新交易敏感词
//	@param chainId
//	@param transaction 交易信息
//	@return error
func UpdateTransactionBak(chainId string, transaction *db.Transaction) error {
	if chainId == "" || transaction == nil {
		return nil
	}

	where := map[string]interface{}{
		"txId": transaction.TxId,
	}

	params := map[string]interface{}{
		"contractResult":        transaction.ContractResult,
		"contractResultBak":     transaction.ContractResultBak,
		"contractMessage":       transaction.ContractMessage,
		"contractMessageBak":    transaction.ContractMessageBak,
		"contractParameters":    transaction.ContractParameters,
		"contractParametersBak": transaction.ContractParametersBak,
		"readSet":               transaction.ReadSet,
		"readSetBak":            transaction.ReadSetBak,
		"writeSet":              transaction.WriteSet,
		"writeSetBak":           transaction.WriteSetBak,
	}

	tableName := db.GetTableName(chainId, db.TableTransaction)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateTransactionContractName
//
//	@Description: 更新交易合约名称
//	@param chainId
//	@param contract 合约信息
//	@return error
func UpdateTransactionContractName(chainId string, contract *db.Contract) error {
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
	tableName := db.GetTableName(chainId, db.TableTransaction)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}
