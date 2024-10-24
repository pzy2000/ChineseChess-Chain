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
	"strings"
	"time"
)

/*-------非同质化持仓-----------*/

// InsertNonFungiblePosition 插入非同质化持仓
func InsertNonFungiblePosition(chainId string, positions []*db.NonFungiblePosition) error {
	if len(positions) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableNonFungiblePosition)
	err := CreateInBatchesData(tableName, positions)
	if err != nil {
		return err
	}

	// 更新缓存
	for _, position := range positions {
		ftPositionCache := GetNFTPositionDataCache(chainId, position.OwnerAddr)
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
				SetNFTPositionDataCache(chainId, position.OwnerAddr, ftPositionCache)
			}
		}
	}
	return nil
}

// UpdateNonFungiblePosition 更新非同质化持仓
func UpdateNonFungiblePosition(chainId string, positions []*db.NonFungiblePosition) error {
	if len(positions) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableNonFungiblePosition)
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
		cachedPositions := GetNFTPositionDataCache(chainId, position.OwnerAddr)
		for _, cachedPosition := range cachedPositions {
			if cachedPosition.ContractAddr == position.ContractAddr {
				cachedPosition.Amount = position.Amount
				cachedPosition.BlockHeight = position.BlockHeight
			}
		}
		SetNFTPositionDataCache(chainId, position.OwnerAddr, cachedPositions)
	}
	return nil
}

// DeleteNonFungiblePosition 删除非同质化持仓
func DeleteNonFungiblePosition(chainId string, positions []*db.NonFungiblePosition) error {
	if len(positions) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableNonFungiblePosition)
	for _, position := range positions {
		where := map[string]interface{}{
			"ownerAddr":    position.OwnerAddr,
			"contractAddr": position.ContractAddr,
		}
		err := db.GormDB.Table(tableName).Where(where).Delete(&db.NonFungiblePosition{}).Error
		if err != nil {
			return err
		}

		// 更新缓存
		cachedPositions := GetNFTPositionDataCache(chainId, position.OwnerAddr)
		updatedPositions := make([]*db.NonFungiblePosition, 0)
		for _, cachedPosition := range cachedPositions {
			if cachedPosition.ContractAddr != position.ContractAddr {
				updatedPositions = append(updatedPositions, cachedPosition)
			}
		}
		SetNFTPositionDataCache(chainId, position.OwnerAddr, updatedPositions)
	}
	return nil
}

// GetNonFungiblePositionListWithRank
//
//	@Description:
//	@param offset
//	@param limit
//	@param chainId
//	@param contractAddr
//	@param ownerAddr
//	@return []*db.PositionWithRank
//	@return int64
//	@return error
func GetNonFungiblePositionListWithRank(offset, limit int, chainId, contractAddr, ownerAddr string) (
	[]*db.PositionWithRank, int64, error) {
	var totalCount int64
	positionList := make([]*db.PositionWithRank, 0)
	if chainId == "" {
		return positionList, totalCount, nil
	}

	// 查询总数
	where := map[string]interface{}{}
	if contractAddr != "" {
		where["contractAddr"] = contractAddr
	}
	if ownerAddr != "" {
		where["ownerAddr"] = ownerAddr
	}
	tableName := db.GetTableName(chainId, db.TableNonFungiblePosition)
	query := db.GormDB.Table(tableName).Where(where)
	err := query.Count(&totalCount).Error
	if err != nil {
		return positionList, 0, err
	}

	// 查询持仓列表和排名
	rankQuery, args := BuildPositionRankSql(tableName, contractAddr, ownerAddr, limit, offset)
	err = db.GormDB.Raw(rankQuery, args...).Scan(&positionList).Error

	if err != nil {
		return nil, 0, err
	}

	return positionList, totalCount, nil
}

//nolint:goconst
func BuildPositionRankSql(tableName, contractAddr, ownerAddr string, limit, offset int) (string, []interface{}) {
	// 查询持仓列表和排名
	sonSql := fmt.Sprintf("SELECT *, RANK() OVER (ORDER BY CAST(amount AS DECIMAL(65, 30)) DESC) as holdRank "+
		"FROM %s", tableName)
	var args []interface{}
	var whereClauses []string
	//contractAddr和ownerAddr计算排名只能用其中一个值
	//同一个合约间排名，同一个持有账户排名
	if contractAddr != "" {
		sonSql += " WHERE contractAddr = ? "
		args = append(args, contractAddr)
	} else if ownerAddr != "" {
		sonSql += " WHERE ownerAddr = ? "
		args = append(args, ownerAddr)
	}
	rankQuery := fmt.Sprintf("SELECT * FROM (%s) as subquery", sonSql)

	if contractAddr != "" {
		whereClauses = append(whereClauses, "contractAddr = ?")
		args = append(args, contractAddr)
	}

	if ownerAddr != "" {
		whereClauses = append(whereClauses, "ownerAddr = ?")
		args = append(args, ownerAddr)
	}

	if len(whereClauses) > 0 {
		rankQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rankQuery += " ORDER BY holdRank ASC LIMIT ? OFFSET ?" // 添加了"ORDER BY"子句
	args = append(args, limit, offset)

	return rankQuery, args
}

// GetNFTPositionList 获取非同质化合约持仓
func GetNFTPositionList(offset int, limit int, chainId, contractAddr, ownerAddr string) (
	[]*db.NonFungiblePosition, error) {
	positionList := make([]*db.NonFungiblePosition, 0)
	if chainId == "" || contractAddr == "" {
		return positionList, db.ErrTableParams
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	if ownerAddr != "" {
		where["ownerAddr"] = ownerAddr
	}

	tableName := db.GetTableName(chainId, db.TableNonFungiblePosition)
	query := db.GormDB.Table(tableName).Where(where)
	// 对 Amount 进行降序排序
	query = query.Order("amount desc").Offset(offset * limit).Limit(limit)
	err := query.Find(&positionList).Error
	return positionList, err
}

// GetNFTPositionHoldRankByAmount 获取排名
func GetNFTPositionHoldRankByAmount(chainId, contractAddr string, amount string) (int64, error) {
	var count int64
	if chainId == "" || contractAddr == "" {
		return count, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableNonFungiblePosition)
	err := db.GormDB.Table(tableName).Where("contractAddr = ? AND amount > ?",
		contractAddr, amount).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateNonPositionContractName 更新合约名称
func UpdateNonPositionContractName(chainId, contractName, contractAddr string) error {
	if chainId == "" || contractName == "" || contractAddr == "" {
		return nil
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	params := map[string]interface{}{
		"contractName": contractName,
	}
	tableName := db.GetTableName(chainId, db.TableNonFungiblePosition)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// GetNonFungiblePositionByOwner 获取非同质化合约持仓
func GetNonFungiblePositionByOwner(chainId string, ownerAddr []string) (map[string][]*db.NonFungiblePosition, error) {
	positionExists := make(map[string][]*db.NonFungiblePosition, 0)
	if len(ownerAddr) == 0 {
		return positionExists, nil
	}

	missingAddrs := make([]string, 0)
	for _, addr := range ownerAddr {
		nftPosition := GetNFTPositionDataCache(chainId, addr)
		if len(nftPosition) > 0 {
			positionExists[addr] = nftPosition
		} else {
			missingAddrs = append(missingAddrs, addr)
		}
	}

	if len(missingAddrs) == 0 {
		return positionExists, nil
	}

	positionList := make([]*db.NonFungiblePosition, 0)
	tableName := db.GetTableName(chainId, db.TableNonFungiblePosition)
	err := db.GormDB.Table(tableName).Where("ownerAddr in ?", missingAddrs).Find(&positionList).Error
	if err != nil {
		return positionExists, err
	}

	for _, position := range positionList {
		positionExists[position.OwnerAddr] = append(positionExists[position.OwnerAddr], position)
		SetNFTPositionDataCache(chainId, position.OwnerAddr, positionExists[position.OwnerAddr])
	}
	return positionExists, nil
}

// GetNFTPositionDataCache
//
//	@Description: 获取持仓数据缓存
//	@param chainId
//	@param address
//	@return []*db.FungiblePosition
func GetNFTPositionDataCache(chainId, address string) []*db.NonFungiblePosition {
	result := make([]*db.NonFungiblePosition, 0)
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisNFTContractPositionData, prefix, chainId, address)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes != nil && redisRes.Val() != "" {
		err := json.Unmarshal([]byte(redisRes.Val()), &result)
		if err == nil {
			return result
		}
		log.Errorf("GetNFTPositionDataCache json Unmarshal err : %s, redisRes:%v", err.Error(), redisRes)
	}

	return nil
}

// SetNFTPositionDataCache
//
//	@Description: 设置持仓缓存
//	@param chainId
//	@param address
//	@param ftPosition
func SetNFTPositionDataCache(chainId, address string, nftPosition []*db.NonFungiblePosition) {
	if len(nftPosition) == 0 {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisNFTContractPositionData, prefix, chainId, address)
	retJson, err := json.Marshal(nftPosition)
	if err == nil {
		// 设置键值对和过期时间
		ctx := context.Background()
		_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), 30*time.Minute).Err()
	}
}

//nolint:goconst
func GetNFTPositionJoinAccount(offset, limit int, chainId, contractAddr, ownerAddr string) (
	[]*db.ContractPositionAccount, error) {
	positionList := make([]*db.ContractPositionAccount, 0)
	if chainId == "" || contractAddr == "" {
		return positionList, nil
	}

	positionTableName := db.GetTableName(chainId, db.TableNonFungiblePosition)
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

//nolint:goconst
func GetNFTPositionByAddrJoinAccount(chainId, contractAddr string, ownerAddr []string) (
	[]*db.ContractPositionAccount, error) {
	positionList := make([]*db.ContractPositionAccount, 0)
	if chainId == "" || contractAddr == "" || len(ownerAddr) == 0 {
		return positionList, nil
	}

	positionTableName := db.GetTableName(chainId, db.TableNonFungiblePosition)
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
