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
	"encoding/json"
	"fmt"
	"time"
)

// InsertBatchGas 批量插入gas
func InsertBatchGas(chainId string, gasList []*db.Gas) error {
	if len(gasList) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableGas)
	return CreateInBatchesData(tableName, gasList)
}

// UpdateGas update
func UpdateGas(chainId string, gasInfo *db.Gas) error {
	if gasInfo == nil || gasInfo.Address == "" {
		return nil
	}

	where := map[string]interface{}{
		"address": gasInfo.Address,
	}
	params := map[string]interface{}{
		"gasBalance":  gasInfo.GasBalance,
		"gasTotal":    gasInfo.GasTotal,
		"gasUsed":     gasInfo.GasUsed,
		"blockHeight": gasInfo.BlockHeight,
	}

	tableName := db.GetTableName(chainId, db.TableGas)
	err := db.GormDB.Table(tableName).Model(&db.Gas{}).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	//设置缓存
	SetAddressGasInfoCache(chainId, gasInfo)
	return nil
}

// GetGasList GetGasList
func GetGasList(offset int, limit int, chainId string, addrList []string) ([]*db.Gas, int64, error) {
	var count int64
	gasList := make([]*db.Gas, 0)
	if chainId == "" {
		return gasList, 0, db.ErrTableParams
	}

	whereIn := map[string]interface{}{}
	if len(addrList) > 0 {
		whereIn["address"] = addrList
	}
	selectFile := &SelectFile{
		WhereIn: whereIn,
	}
	tableName := db.GetTableName(chainId, db.TableGas)
	query := BuildParamsQuery(tableName, selectFile)
	err := query.Count(&count).Error
	if err != nil {
		return gasList, 0, err
	}
	err = query.Order("gasBalance desc").Offset(offset * limit).Limit(limit).Find(&gasList).Error
	if err != nil {
		return gasList, 0, err
	}

	return gasList, count, nil
}

// GetGasByAddrInfo 根据多个addr获取Gas余额
func GetGasByAddrInfo(chainId string, addrList []string) (int64, error) {
	if len(addrList) == 0 || chainId == "" {
		return 0, nil
	}

	gasList := make([]*db.Gas, 0)
	whereIn := map[string]interface{}{}
	if len(addrList) > 0 {
		whereIn["address"] = addrList
	}
	selectFile := &SelectFile{
		WhereIn: whereIn,
	}
	tableName := db.GetTableName(chainId, db.TableGas)
	query := BuildParamsQuery(tableName, selectFile)
	err := query.Find(&gasList).Error
	if err != nil {
		return 0, fmt.Errorf("GetGasByAddrInfo err, cause : %s", err.Error())
	}

	var gasBalance int64
	for _, datum := range gasList {
		gasBalance = gasBalance + datum.GasBalance
	}
	return gasBalance, nil
}

// GetGasInfoByAddr 根据Addr获取gas详情
func GetGasInfoByAddr(chainId string, addrList []string) ([]*db.Gas, error) {
	gasList := make([]*db.Gas, 0)
	if len(addrList) == 0 || chainId == "" {
		return gasList, nil
	}

	// 尝试从缓存中获取gas信息
	missingAddrs := make([]string, 0)
	for _, addr := range addrList {
		gasInfo := GetAddressGasInfoCache(chainId, addr)
		if gasInfo != nil {
			gasList = append(gasList, gasInfo)
		} else {
			missingAddrs = append(missingAddrs, addr)
		}
	}

	if len(missingAddrs) == 0 {
		return gasList, nil
	}

	tableName := db.GetTableName(chainId, db.TableGas)
	query := db.GormDB.Table(tableName).Where("address in ?", missingAddrs)
	missingGasList := make([]*db.Gas, 0)
	err := query.Find(&missingGasList).Error
	if err != nil {
		return gasList, err
	}

	// 将查询结果添加到gasList中，并将结果写入缓存
	for _, gasInfo := range missingGasList {
		gasList = append(gasList, gasInfo)
		SetAddressGasInfoCache(chainId, gasInfo)
	}
	return gasList, nil
}

// GetAddressGasInfoCache
//
//	@Description: 获取账户gas信息
//	@param chainId
//	@param address
//	@return *db.Gas
func GetAddressGasInfoCache(chainId, address string) *db.Gas {
	var result *db.Gas
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisUserAddressGasInfo, prefix, chainId, address)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes != nil && redisRes.Val() != "" {
		err := json.Unmarshal([]byte(redisRes.Val()), &result)
		if err == nil {
			return result
		}
		log.Errorf("GetAddressGasInfoCache json Unmarshal err : %s, redisRes:%v", err.Error(), redisRes)
	}

	return nil
}

// SetAddressGasInfoCache
//
//	@Description: 设置账户gas
//	@param chainId
//	@param gasInfo
func SetAddressGasInfoCache(chainId string, gasInfo *db.Gas) {
	if gasInfo == nil {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisUserAddressGasInfo, prefix, chainId, gasInfo.Address)
	retJson, err := json.Marshal(gasInfo)
	if err == nil {
		// 设置键值对和过期时间
		ctx := context.Background()
		_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), 30*time.Minute).Err()
	}
}
