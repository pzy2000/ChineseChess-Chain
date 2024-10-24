/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// InsertContract 插入合约
func InsertContract(chainId string, contract *db.Contract) error {
	if contract == nil || contract.Addr == "" {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableContract)
	err := InsertData(tableName, contract)
	if err != nil {
		return err
	}

	//设置合约缓存
	SetContractInfoCache(chainId, contract)
	//删除合约数量缓存
	DelContractCountCache(chainId)
	return nil
}

// UpdateContract 更新合约
func UpdateContract(chainId string, contract *db.Contract) error {
	if contract == nil || contract.Addr == "" {
		return nil
	}
	where := map[string]interface{}{
		"addr": contract.Addr,
	}

	params := map[string]interface{}{}
	params["contractStatus"] = contract.ContractStatus
	if contract.UpgradeTimestamp != 0 {
		params["upgradeTimestamp"] = contract.UpgradeTimestamp
	}
	if contract.Upgrader != "" {
		params["upgrader"] = contract.Upgrader
	}
	if contract.UpgradeAddr != "" {
		params["upgradeAddr"] = contract.UpgradeAddr
	}
	if contract.UpgradeOrgId != "" {
		params["upgradeOrgId"] = contract.UpgradeOrgId
	}
	if contract.Version != "" {
		params["version"] = contract.Version
	}
	if contract.ContractSymbol != "" {
		params["contractSymbol"] = contract.ContractSymbol
	}
	if contract.Decimals != 0 {
		params["decimals"] = contract.Decimals
	}

	if len(params) == 0 {
		return nil
	}

	// 获取表名
	tableName := db.GetTableName(chainId, db.TableContract)
	// 更新记录
	err := db.GormDB.Table(tableName).Model(&db.Contract{}).Where(where).Updates(params).Error
	if err != nil {
		return fmt.Errorf("update contract fail, chainId : %s, name : %s, version : %s, cause: %s",
			chainId, contract.Name, contract.Version, err.Error())
	}

	//更新合约缓存
	UpdateContractCache(chainId, contract)
	return nil
}

// UpdateContractTxNum 更新合约交易数量
func UpdateContractTxNum(chainId string, contract *db.Contract) error {
	if contract == nil || contract.Addr == "" {
		return nil
	}
	where := map[string]interface{}{
		"addr": contract.Addr,
	}

	params := map[string]interface{}{}
	if contract.TxNum > 0 {
		params["txNum"] = contract.TxNum
	}
	if contract.EventNum > 0 {
		params["eventNum"] = contract.EventNum
	}

	if len(params) == 0 {
		return nil
	}
	// 获取表名
	tableName := db.GetTableName(chainId, db.TableContract)
	// 更新记录
	err := db.GormDB.Table(tableName).Model(&db.Contract{}).Where(where).Updates(params).Error
	if err != nil {
		return fmt.Errorf("update contract fail, chainId : %s, name : %s, version : %s, cause: %s",
			chainId, contract.Name, contract.Version, err.Error())
	}

	//更新合约缓存
	UpdateContractCache(chainId, contract)
	return nil
}

// GetContractByCacheOrAddr
//
//	@Description: 通过合约地址获取合约信息，先行缓存获取，没有在从DB获取
//	@param chainId
//	@param contractAddr 合约地址
func GetContractByCacheOrAddr(chainId, contractAddr string) (*db.Contract, error) {
	//缓存获取合约信息
	contractInfo, err := GetContractCacheByNameOrAddr(chainId, contractAddr)
	if contractInfo == nil || err != nil {
		//DB获取合约信息
		contractInfo, err = GetContractByAddr(chainId, contractAddr)
	}

	if err != nil {
		return nil, err
	}

	return contractInfo, nil
}

// GetContractByCacheOrName
//
//	@Description: 通过合约地址获取合约信息，先行缓存获取，没有在从DB获取
//	@param chainId
//	@param contractName 合约名称
func GetContractByCacheOrName(chainId, contractName string) (*db.Contract, error) {
	if chainId == "" || contractName == "" {
		return nil, errors.New("contractName is null")
	}
	//缓存获取合约信息
	contractInfo, err := GetContractCacheByNameOrAddr(chainId, contractName)
	if contractInfo == nil || err != nil {
		//DB获取合约信息
		contractInfo, err = GetContractByName(chainId, contractName)
	}

	if err != nil {
		return nil, err
	}

	return contractInfo, nil
}

// GetContractByCacheOrNameAddr
//
//	@Description: 通过合约名称或地址获取合约信息，先行缓存获取，没有在从DB获取
//	@param chainId
//	@param contractName 合约名称
func GetContractByCacheOrNameAddr(chainId, contractKey string) (*db.Contract, error) {
	//缓存获取合约信息
	contractInfo, err := GetContractCacheByNameOrAddr(chainId, contractKey)
	if contractInfo == nil || err != nil {
		//DB获取合约信息
		contractInfo, err = GetContractByNameOrAddr(chainId, contractKey)
	}

	if err != nil {
		return nil, err
	}

	return contractInfo, nil
}

// GetContractByCacheOrAddrs
//
//	@Description: 根据合约地址批量获取合约
//	@param chainId
//	@param contractAdders 合约地址
//	@return map[string]*db.Contract
//	@return error
func GetContractByCacheOrAddrs(chainId string, contractAdders []string) (map[string]*db.Contract, error) {
	contracts := make([]*db.Contract, 0)
	contractMap := make(map[string]*db.Contract, 0)
	selectAddrs := make([]string, 0)
	if chainId == "" || len(contractAdders) == 0 {
		return contractMap, nil
	}

	// 将contracts转换为contractExists映射
	for _, addr := range contractAdders {
		//缓存获取合约信息
		contractInfo, err := GetContractCacheByNameOrAddr(chainId, addr)
		if err == nil && contractInfo != nil {
			contractMap[contractInfo.NameBak] = contractInfo
			contractMap[contractInfo.Addr] = contractInfo
		} else {
			selectAddrs = append(selectAddrs, addr)
		}
	}

	if len(selectAddrs) == 0 {
		return contractMap, nil
	}

	tableName := db.GetTableName(chainId, db.TableContract)
	err := db.GormDB.Table(tableName).Where("addr in ?", selectAddrs).Find(&contracts).Error
	if err != nil {
		return contractMap, err
	}

	for _, contract := range contracts {
		contractMap[contract.NameBak] = contract
		contractMap[contract.Addr] = contract
		//设置合约缓存
		SetContractInfoCache(chainId, contract)
	}
	return contractMap, nil
}

// GetContractByCacheOrAddrs
//
//	@Description: 根据合约地址批量获取合约
//	@param chainId
//	@param contractAdders 合约地址
//	@return map[string]*db.Contract
//	@return error
func GetContractByAddrs(chainId string, contractAdders []string) (map[string]*db.Contract, error) {
	contracts := make([]*db.Contract, 0)
	contractMap := make(map[string]*db.Contract, 0)
	if chainId == "" || len(contractAdders) == 0 {
		return contractMap, nil
	}

	tableName := db.GetTableName(chainId, db.TableContract)
	err := db.GormDB.Table(tableName).Where("addr in ?", contractAdders).Find(&contracts).Error
	if err != nil {
		return contractMap, err
	}

	for _, contract := range contracts {
		contractMap[contract.NameBak] = contract
		contractMap[contract.Addr] = contract
	}
	return contractMap, nil
}

// GetContractByAddersOrNames 根据合约地址获取合约
func GetContractByAddersOrNames(chainId string, nameList []string) (map[string]*db.Contract, error) {
	contractExists := make(map[string]*db.Contract, 0)
	contracts := make([]*db.Contract, 0)
	if len(nameList) == 0 {
		return contractExists, nil
	}
	tableName := db.GetTableName(chainId, db.TableContract)
	err := db.GormDB.Table(tableName).Where("nameBak IN (?) OR addr IN (?)",
		nameList, nameList).Find(&contracts).Error
	if err != nil {
		return contractExists, fmt.Errorf("GetGasByAddrInfo err, cause : %s", err.Error())
	}

	// 将contracts转换为contractExists映射
	for _, contract := range contracts {
		contractExists[contract.NameBak] = contract
		contractExists[contract.Addr] = contract
	}
	return contractExists, nil
}

// GetContractByAddr 根据合约地址获取合约
func GetContractByAddr(chainId, contractAddr string) (*db.Contract, error) {
	if chainId == "" || contractAddr == "" {
		return nil, db.ErrTableParams
	}

	contractInfo := &db.Contract{}
	where := map[string]interface{}{
		"addr": contractAddr,
	}
	tableName := db.GetTableName(chainId, db.TableContract)

	err := db.GormDB.Table(tableName).Where(where).First(&contractInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	//设置合约缓存
	SetContractInfoCache(chainId, contractInfo)
	return contractInfo, nil
}

// GetContractByName 根据合约名称获取合约
func GetContractByName(chainId, contractName string) (*db.Contract, error) {
	if chainId == "" || contractName == "" {
		return nil, db.ErrTableParams
	}

	contractInfo := &db.Contract{}
	where := map[string]interface{}{
		"nameBak": contractName,
	}
	tableName := db.GetTableName(chainId, db.TableContract)
	err := db.GormDB.Table(tableName).Where(where).First(&contractInfo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	//设置合约缓存
	SetContractInfoCache(chainId, contractInfo)

	return contractInfo, nil
}

// GetContractNum 获取合约数量
func GetContractNum(chainId string) (int64, error) {
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisOverviewContractCount, prefix, chainId)
	countRes, err := GetContractNumCache(redisKey)
	if err == nil && countRes != 0 {
		return countRes, nil
	}

	var count int64
	tableName := db.GetTableName(chainId, db.TableContract)
	err = db.GormDB.Table(tableName).Where("contractStatus > ?", SystemContractStatus).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count contract err, cause : %s", err.Error())
	}

	// 设置键值对和过期时间
	ctx := context.Background()
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, count, time.Hour).Err()
	return count, nil
}

// GetContractNumCache 获取合约数量缓存
func GetContractNumCache(redisKey string) (int64, error) {
	//获取缓存
	var count int64
	ctx := context.Background()
	result, err := cache.GlobalRedisDb.Get(ctx, redisKey).Result()
	if err != nil {
		return 0, err
	}
	count, err = strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetContractByNameOrAddr 根据合约名称获取合约，contractKey is name  or addr
func GetContractByNameOrAddr(chainId, contractKey string) (*db.Contract, error) {
	if chainId == "" || contractKey == "" {
		return nil, db.ErrTableParams
	}

	//从数据库获取
	var contract *db.Contract
	// 构建查询条件
	tableName := db.GetTableName(chainId, db.TableContract)
	err := db.GormDB.Table(tableName).Where("nameBak = ?", contractKey).
		Or("addr = ?", contractKey).First(&contract).Error
	if err != nil || contract == nil {
		return nil, err
	}

	//设置合约缓存
	SetContractInfoCache(chainId, contract)

	return contract, nil
}

// GetLatestContractList 获取最后10条合约
func GetLatestContractList(chainId string) ([]*db.Contract, error) {
	contractList := make([]*db.Contract, 0)
	tableName := db.GetTableName(chainId, db.TableContract)
	query := db.GormDB.Table(tableName).Order("timestamp desc").Limit(10)
	err := query.Find(&contractList).Error
	if err != nil {
		return contractList, err
	}
	return contractList, nil
}

// GetContractList 获取合约列表
func GetContractList(chainId string, offset, limit int, status *int32, runtimeType, contractKey string,
	creators, creatorAddrs, upgraders, upgradeAddrs []string, startTime, endTime int64) ([]*db.Contract, int64, error) {
	var count int64
	contracts := make([]*db.Contract, 0)
	where := map[string]interface{}{}
	currentWhere := map[string]interface{}{}
	whereIn := map[string]interface{}{}
	whereOr := map[string]interface{}{}
	if runtimeType != "" {
		where["runtimeType"] = runtimeType
	}
	if status != nil && *status != -1 {
		where["contractStatus"] = *status
	} else {
		currentWhere["contractStatus > ?"] = SystemContractStatus
	}
	if len(creators) > 0 {
		whereIn["createSender"] = creators
	}
	if len(creatorAddrs) > 0 {
		whereIn["creatorAddr"] = creatorAddrs
	}
	if len(upgraders) > 0 {
		whereIn["upgrader"] = upgraders
	}
	if len(upgradeAddrs) > 0 {
		whereIn["upgradeAddr"] = upgradeAddrs
	}

	if contractKey != "" {
		whereOr["nameBak"] = contractKey
		whereOr["addr"] = contractKey
	}

	selectFile := &SelectFile{
		Where:        where,
		CurrentWhere: currentWhere,
		WhereIn:      whereIn,
		WhereOr:      whereOr,
		StartTime:    startTime,
		EndTime:      endTime,
	}
	tableName := db.GetTableName(chainId, db.TableContract)
	query := BuildParamsQuery(tableName, selectFile)
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	query = query.Order("timestamp desc")
	query = query.Offset(offset * limit).Limit(limit)
	err = query.Find(&contracts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("GetContractList err, cause : %s", err.Error())
	}
	return contracts, count, nil
}

// UpdateContractNameBak 更新合约敏感词
// @desc
// @param ${param}
// @return error
func UpdateContractNameBak(chainId string, contract *db.Contract) error {
	if chainId == "" || contract == nil {
		return nil
	}

	where := map[string]interface{}{
		"addr": contract.Addr,
	}
	params := map[string]interface{}{
		"name":    contract.Name,
		"nameBak": contract.NameBak,
	}
	tableName := db.GetTableName(chainId, db.TableContract)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return fmt.Errorf("update contract name fail, chainId : %s, addr : %s , cause: %s",
			chainId, contract.Addr, err.Error())
	}
	return nil
}

// GetTotalTxNum 获取交易总量
func GetTotalTxNum(chainId string) (int64, error) {
	txTotal, err := GetTotalTxNumCache(chainId)
	if err == nil && txTotal != 0 {
		return txTotal, nil
	}

	tableName := db.GetTableName(chainId, db.TableContract)
	var totalTxNum int64 // 使用 interface{} 类型
	err = db.GormDB.Table(tableName).Select("sum(txNum)").Row().Scan(&totalTxNum)
	if err != nil {
		return 0, fmt.Errorf("select sum(txNum) err, cause : %v", err)
	}

	SetTotalTxNumCache(chainId, totalTxNum)
	return totalTxNum, nil
}

// GetContractCountByRange 获取指定时间内的交易数量
func GetContractCountByRange(chainId string, startTime, endTime int64) (int64, error) {
	var totalCount int64
	tableName := db.GetTableName(chainId, db.TableContract)
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
