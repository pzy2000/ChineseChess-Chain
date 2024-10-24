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

	"github.com/redis/go-redis/v9"
)

/*-------同质化持仓-----------*/

// InsertFungiblePosition 插入同质化持仓
func InsertFungiblePosition(chainId string, positions []*db.FungiblePosition) error {
	if len(positions) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableFungiblePosition)
	err := CreateInBatchesData(tableName, positions)
	if err != nil {
		return err
	}

	// 更新缓存
	for _, position := range positions {
		ftPositionCache := GetFTPositionDataCache(chainId, position.OwnerAddr)
		if len(ftPositionCache) > 0 {
			exists := false
			for _, v := range ftPositionCache {
				if v.ContractAddr == position.ContractAddr {
					exists = true
					break
				}
			}

			if !exists {
				ftPositionCache = append(ftPositionCache, position)
				SetFTPositionDataCache(chainId, position.OwnerAddr, ftPositionCache)
			}
		}
	}
	return nil
}

// UpdateFungiblePosition 更新同质化持仓
func UpdateFungiblePosition(chainId string, positions []*db.FungiblePosition) error {
	if len(positions) == 0 {
		return nil
	}

	ownerAdders := make([]string, 0)
	for _, position := range positions {
		ownerAdders = append(ownerAdders, position.OwnerAddr)
	}
	positionMapList, err := GetFungiblePositionByOwners(chainId, ownerAdders)
	if err != nil {
		return err
	}

	tableName := db.GetTableName(chainId, db.TableFungiblePosition)
	for _, position := range positions {
		where := map[string]interface{}{
			"ownerAddr":    position.OwnerAddr,
			"contractAddr": position.ContractAddr,
		}

		params := map[string]interface{}{
			"amount":      position.Amount,
			"blockHeight": position.BlockHeight,
		}
		err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
		if err != nil {
			return err
		}

		// 更新缓存
		if cachedPositionList, ok := positionMapList[position.OwnerAddr]; ok {
			for _, cachedPosition := range cachedPositionList {
				if cachedPosition.ContractAddr == position.ContractAddr {
					cachedPosition.Amount = position.Amount
					cachedPosition.BlockHeight = position.BlockHeight
				}
			}
			SetFTPositionDataCache(chainId, position.OwnerAddr, cachedPositionList)
		}
	}
	return nil
}

// DeleteFungiblePosition 删除同质化持仓
func DeleteFungiblePosition(chainId string, positions []*db.FungiblePosition) error {
	if len(positions) == 0 {
		return nil
	}

	ownerAdders := make([]string, 0)
	for _, position := range positions {
		ownerAdders = append(ownerAdders, position.OwnerAddr)
	}
	positionMapList, err := GetFungiblePositionByOwners(chainId, ownerAdders)
	if err != nil {
		return err
	}

	tableName := db.GetTableName(chainId, db.TableFungiblePosition)
	for _, position := range positions {
		where := map[string]interface{}{
			"ownerAddr":    position.OwnerAddr,
			"contractAddr": position.ContractAddr,
		}
		err := db.GormDB.Table(tableName).Where(where).Delete(&db.FungiblePosition{}).Error
		if err != nil {
			return err
		}

		// 更新缓存
		if cachedPositionList, ok := positionMapList[position.OwnerAddr]; ok {
			updatedPositions := make([]*db.FungiblePosition, 0)
			for _, cachedPosition := range cachedPositionList {
				if cachedPosition.ContractAddr != position.ContractAddr {
					updatedPositions = append(updatedPositions, cachedPosition)
				}
			}
			SetFTPositionDataCache(chainId, position.OwnerAddr, updatedPositions)
		}
	}
	return nil
}

// UpdatePositionContractName 更新合约名称
func UpdatePositionContractName(chainId, contractName, contractAddr string) error {
	if chainId == "" || contractName == "" || contractAddr == "" {
		return nil
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	params := map[string]interface{}{
		"contractName": contractName,
	}
	tableName := db.GetTableName(chainId, db.TableFungiblePosition)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// GetFTPositionCountByAddr 获取同质化合约持仓
func GetFTPositionCountByAddr(chainId, ownerAddr string) (int64, error) {
	var count int64
	if chainId == "" || ownerAddr == "" {
		return count, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableFungiblePosition)
	query := db.GormDB.Table(tableName).Where("ownerAddr = ?", ownerAddr)
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetFTPositionListByAddr 获取同质化合约持仓
//
//nolint:goconst
func GetFTPositionListByAddr(offset, limit int, chainId, ownerAddr string) ([]*db.FTPositionJoinContract, error) {
	positionList := make([]*db.FTPositionJoinContract, 0)
	if chainId == "" || ownerAddr == "" {
		return positionList, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableFungiblePosition)
	contractTableName := db.GetTableName(chainId, db.TableFungibleContract)

	// 查询数据：在FungiblePosition和FungibleContract表上执行连表查询
	dataQuery := db.GormDB.Table(tableName).Select(tableName+".*, "+contractTableName+".contractType").
		Joins("LEFT JOIN "+contractTableName+" ON "+
			tableName+".contractAddr = "+contractTableName+".contractAddr").
		Where(tableName+".ownerAddr = ?", ownerAddr)

	// 对 Amount 进行降序排序
	dataQuery = dataQuery.Order(tableName + ".amount desc").Offset(offset * limit).Limit(limit)
	err := dataQuery.Find(&positionList).Error
	if err != nil {
		return positionList, err
	}
	return positionList, nil
}

// GetFungiblePositionByOwnerAddr 获取同质化合约持仓
func GetFungiblePositionByOwnerAddr(chainId, ownerAddr string) ([]*db.FungiblePosition, error) {
	positionList := make([]*db.FungiblePosition, 0)
	if len(ownerAddr) == 0 {
		return positionList, nil
	}

	ftPositionCache := GetFTPositionDataCache(chainId, ownerAddr)
	if len(ftPositionCache) > 0 {
		return ftPositionCache, nil
	}

	tableName := db.GetTableName(chainId, db.TableFungiblePosition)
	err := db.GormDB.Table(tableName).Where("ownerAddr = ?", ownerAddr).Find(&positionList).Error
	if err != nil {
		return positionList, err
	}

	SetFTPositionDataCache(chainId, ownerAddr, positionList)
	return positionList, nil
}

// GetFungiblePositionByOwners 获取同质化合约持仓
func GetFungiblePositionByOwners(chainId string, ownerAddr []string) (map[string][]*db.FungiblePosition, error) {
	positionExists := make(map[string][]*db.FungiblePosition, 0)
	if len(ownerAddr) == 0 {
		return positionExists, nil
	}
	missingAddrs := make([]string, 0)
	for _, addr := range ownerAddr {
		ftPosition := GetFTPositionDataCache(chainId, addr)
		if len(ftPosition) > 0 {
			positionExists[addr] = ftPosition
		} else {
			missingAddrs = append(missingAddrs, addr)
		}
	}

	if len(missingAddrs) == 0 {
		return positionExists, nil
	}

	positionList := make([]*db.FungiblePosition, 0)
	tableName := db.GetTableName(chainId, db.TableFungiblePosition)
	err := db.GormDB.Table(tableName).Where("ownerAddr in ?", missingAddrs).Find(&positionList).Error
	if err != nil {
		return positionExists, err
	}

	for _, position := range positionList {
		positionExists[position.OwnerAddr] = append(positionExists[position.OwnerAddr], position)
		SetFTPositionDataCache(chainId, position.OwnerAddr, positionExists[position.OwnerAddr])
	}
	return positionExists, nil
}

// GetFTPositionDataCache
//
//	@Description: 获取持仓地址-持仓数据缓存
//	@param chainId
//	@param address
//	@return []*db.FungiblePosition
func GetFTPositionDataCache(chainId, address string) []*db.FungiblePosition {
	result := make([]*db.FungiblePosition, 0)
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisFTContractPositionData, prefix, chainId, address)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes != nil && redisRes.Val() != "" {
		err := json.Unmarshal([]byte(redisRes.Val()), &result)
		if err == nil {
			return result
		}
		log.Errorf("GetFTPositionDataCache json Unmarshal err : %s, redisRes:%v", err.Error(), redisRes)
	}

	return result
}

// SetFTPositionDataCache
//
//	@Description: 根据持仓地址-设置持仓缓存
//	@param chainId
//	@param address
//	@param ftPosition
func SetFTPositionDataCache(chainId, address string, ftPosition []*db.FungiblePosition) {
	if len(ftPosition) == 0 {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisFTContractPositionData, prefix, chainId, address)
	retJson, err := json.Marshal(ftPosition)
	if err == nil {
		// 设置键值对和过期时间
		ctx := context.Background()
		_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), 30*time.Minute).Err()
	}
}

// SetFTPositionListCache
//
//	@Description: 根据合约地址-设置持仓缓存
//	@param chainId
//	@param address
//	@param ftPosition
func SetFTPositionListCache(chainId, contractAddr string, ftPosition []*db.FungiblePosition) {
	if len(ftPosition) == 0 {
		return
	}

	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisContractPositionOwnerList, prefix, chainId, contractAddr)
	// 创建一个 Redis 批量操作
	pipe := cache.GlobalRedisDb.Pipeline()

	// 将数据写入Redis缓存
	for _, position := range ftPosition {
		// 将字符串转换为 decimal.Decimal 值
		//amountDecimal, _ := decimal.NewFromString(position.Amount)
		amountFloat, _ := position.Amount.Float64()
		pipe.ZAdd(ctx, redisKey, redis.Z{
			Score:  amountFloat,
			Member: position.OwnerAddr,
		})
	}

	// 执行批量操作
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Error("SetFTPositionListCache Error executing pipeline: " + err.Error())
	}

	// 设置缓存过期时间
	expirationTime := time.Duration(config.GlobalConfig.RedisDB.PositionRankTime) * time.Second
	cache.GlobalRedisDb.Expire(ctx, redisKey, expirationTime)
}

func SetNFTPositionListCache(chainId, contractAddr string, nftPosition []*db.NonFungiblePosition) {
	if len(nftPosition) == 0 {
		return
	}

	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisContractPositionOwnerList, prefix, chainId, contractAddr)
	// 创建一个 Redis 批量操作
	pipe := cache.GlobalRedisDb.Pipeline()

	// 将数据写入Redis缓存
	for _, position := range nftPosition {
		// 将字符串转换为 decimal.Decimal 值
		//amountDecimal, _ := decimal.NewFromString(position.Amount)
		amountFloat, _ := position.Amount.Float64()
		pipe.ZAdd(ctx, redisKey, redis.Z{
			Score:  amountFloat,
			Member: position.OwnerAddr,
		})
	}

	// 执行批量操作
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Error("SetNFTPositionListCache Error executing pipeline: " + err.Error())
	}

	// 设置缓存过期时间
	expirationTime := time.Duration(config.GlobalConfig.RedisDB.PositionRankTime) * time.Second
	cache.GlobalRedisDb.Expire(ctx, redisKey, expirationTime)
}

// GetFungiblePositionList
//
//	@Description: 获取持仓列表
//	@param offset
//	@param limit
//	@param chainId
//	@param contractAddr
//	@param ownerAddr
//	@return []*db.FungiblePosition
//	@return error
func GetFungiblePositionList(offset, limit int, chainId, contractAddr, ownerAddr string) (
	[]*db.FungiblePosition, error) {
	positionList := make([]*db.FungiblePosition, 0)
	if chainId == "" || contractAddr == "" {
		return positionList, nil
	}

	tableName := db.GetTableName(chainId, db.TableFungiblePosition)
	// 查询总数
	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	if ownerAddr != "" {
		where["ownerAddr"] = ownerAddr
	}
	query := db.GormDB.Table(tableName).Where(where)
	query = query.Order("amount desc").Offset(offset * limit).Limit(limit)
	err := query.Find(&positionList).Error
	if err != nil {
		return positionList, err
	}
	return positionList, nil
}

// GetFTPositionJoinAccount
//
//	@Description:  根据合约地址，用户地址连表account，查询持仓列表
//	@param offset
//	@param limit
//	@param chainId
//	@param contractAddr
//	@param ownerAddr
//	@return []*db.ContractPositionAccount
//	@return error
//
//nolint:goconst
func GetFTPositionJoinAccount(offset, limit int, chainId, contractAddr, ownerAddr string) (
	[]*db.ContractPositionAccount, error) {
	positionList := make([]*db.ContractPositionAccount, 0)
	if chainId == "" || contractAddr == "" {
		return positionList, nil
	}

	positionTableName := db.GetTableName(chainId, db.TableFungiblePosition)
	accountTableName := db.GetTableName(chainId, db.TableAccount)
	// 查询条件
	where := map[string]interface{}{
		positionTableName + ".contractAddr": contractAddr,
	}

	if ownerAddr != "" {
		where[positionTableName+".ownerAddr"] = ownerAddr
	}

	query := db.GormDB.Table(positionTableName).Where(where)
	// 使用 Joins 方法关联 Account 表
	query = query.Joins("JOIN " + accountTableName + " ON " +
		positionTableName + ".ownerAddr = " + accountTableName + ".address")

	// 选择 FungiblePosition 的 amount 和 Account 的 BNS 字段
	query = query.Select(positionTableName + ".id, " + positionTableName + ".ownerAddr, " +
		positionTableName + ".amount, " + accountTableName + ".bns, " + accountTableName + ".addrType")

	query = query.Order("amount desc")
	query = query.Offset(offset * limit).Limit(limit)
	err := query.Find(&positionList).Error
	if err != nil {
		return positionList, err
	}
	return positionList, nil
}

// GetFTPositionByAddrJoinAccount
//
//	@Description: 根据合约地址，用户地址连表account，查询持仓列表
//	@param chainId
//	@param contractAddr
//	@param ownerAddr
//	@return []*db.ContractPositionAccount
//	@return error
//
//nolint:goconst
func GetFTPositionByAddrJoinAccount(chainId, contractAddr string, ownerAddr []string) (
	[]*db.ContractPositionAccount, error) {
	positionList := make([]*db.ContractPositionAccount, 0)
	if chainId == "" || contractAddr == "" || len(ownerAddr) == 0 {
		return positionList, nil
	}

	positionTableName := db.GetTableName(chainId, db.TableFungiblePosition)
	accountTableName := db.GetTableName(chainId, db.TableAccount)

	query := db.GormDB.Table(positionTableName)
	// 选择 FungiblePosition 的 amount 和 Account 的 BNS 字段
	query = query.Select(positionTableName + ".id, " + positionTableName + ".ownerAddr, " +
		positionTableName + ".amount, " + accountTableName + ".bns, " + accountTableName + ".addrType")

	// 使用 Joins 方法关联 Account 表
	query = query.Joins("JOIN " + accountTableName + " ON " +
		positionTableName + ".ownerAddr = " + accountTableName + ".address")
	query = query.Where("contractAddr = ?", contractAddr).Where("ownerAddr in ?", ownerAddr)
	err := query.Find(&positionList).Error
	if err != nil {
		return positionList, err
	}
	return positionList, nil
}
