/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"strconv"

	"gorm.io/gorm"
)

// GetBlockByHeight 根据blockHeight获取block详情
func GetBlockByHeight(chainId string, blockHeight int64) (*db.Block, error) {
	if chainId == "" {
		return nil, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableBlock)
	blockList := make([]*db.Block, 0)
	where := map[string]interface{}{
		"blockHeight": blockHeight,
	}
	//err := db.GormDB.Table(tableName).Where(where).First(blockInfo).Error
	err := db.GormDB.Table(tableName).Where(where).Find(&blockList).Error
	if err != nil {
		return nil, err
	} else if len(blockList) == 0 {
		return nil, nil
	}

	return blockList[0], nil
}

// InsertBlock 插入区块信息
func InsertBlock(chainId string, blockInfo *db.Block) error {
	if blockInfo == nil || chainId == "" {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableBlock)
	return InsertData(tableName, blockInfo)
}

// GetBlockListByRange 获取指定时间内的block列表
func GetBlockListByRange(chainId string, startTime, endTime int64) (int64, error) {
	if startTime == 0 && endTime == 0 {
		maxBlockHeight := GetMaxBlockHeight(chainId)
		return maxBlockHeight, nil
	}

	var totalCount int64
	tableName := db.GetTableName(chainId, db.TableBlock)
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

// GetMaxBlockHeight get max block height in this chain{chainId}
func GetMaxBlockHeight(chainId string) int64 {
	maxBlockHeight, err := GetMaxBlockHeightCache(chainId)
	if err == nil && maxBlockHeight != 0 {
		return maxBlockHeight
	}

	var blockInfo *db.Block
	tableName := db.GetTableName(chainId, db.TableBlock)
	err = db.GormDB.Table(tableName).Order("blockHeight desc").First(&blockInfo).Error
	if err != nil || blockInfo == nil {
		log.Info("GetMaxBlockHeight failed, err : %s", err)
		return 0
	}

	SetMaxBlockHeightCache(chainId, blockInfo.BlockHeight)
	return blockInfo.BlockHeight
}

// GetBlockList 根据chainId,blockKey获取block列表
func GetBlockList(offset int, limit int, chainId string, blockKey string) ([]*db.Block, error) {
	blockList := make([]*db.Block, 0)
	if chainId == "" {
		return blockList, db.ErrTableParams
	}

	where := map[string]interface{}{}
	if blockKey != "" {
		intValue, err := strconv.ParseInt(blockKey, 10, 64)
		if err == nil {
			//blockHeight 数字
			where["blockHeight"] = intValue
		} else {
			//blockHash 不是数字
			where["blockHash"] = blockKey
		}
	}

	tableName := db.GetTableName(chainId, db.TableBlock)
	query := db.GormDB.Table(tableName).Where(where).Order("blockHeight desc").
		Offset(offset * limit).Limit(limit)
	err := query.Find(&blockList).Error
	if err != nil {
		return blockList, err
	}
	return blockList, nil
}

// GetBlockListCount 根据chainId,blockKey获取block列表
func GetBlockListCount(chainId string, blockKey string) (int64, error) {
	var blockInfo *db.Block
	if chainId == "" {
		return 0, db.ErrTableParams
	}

	where := map[string]interface{}{}
	tableName := db.GetTableName(chainId, db.TableBlock)
	if blockKey != "" {
		intValue, err := strconv.ParseInt(blockKey, 10, 64)
		if err == nil {
			//blockHeight 数字
			where["blockHeight"] = intValue
		} else {
			//blockHash 不是数字
			where["blockHash"] = blockKey
		}
		err = db.GormDB.Table(tableName).Where(where).First(&blockInfo).Error
		if err != nil {
			return 0, nil
		}
		if blockInfo != nil {
			return 1, nil
		}
	} else {
		maxBlockHeight := GetMaxBlockHeight(chainId)
		return maxBlockHeight + 1, nil
	}

	return 0, nil
}

// GetLatestBlockList 获取最后10条block
func GetLatestBlockList(chainId string) ([]*db.Block, error) {
	blockList := make([]*db.Block, 0)
	tableName := db.GetTableName(chainId, db.TableBlock)
	query := db.GormDB.Table(tableName).Order("blockHeight desc").Limit(10)
	err := query.Find(&blockList).Error
	if err != nil {
		return blockList, err
	}

	SetLatestBlockListCache(chainId, blockList)
	return blockList, nil
}

// GetBlockByHash 根据blockHash获取block详情
func GetBlockByHash(blockHash string, chainId string) (*db.Block, error) {
	if blockHash == "" || chainId == "" {
		return nil, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableBlock)
	blockInfo := &db.Block{}
	where := map[string]interface{}{
		"blockHash": blockHash,
	}
	err := db.GormDB.Table(tableName).Where(where).First(blockInfo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return blockInfo, nil
}

// GetBlockByStatus 根据status获取block详情
func GetBlockByStatus(chainId string, status int) ([]*db.Block, error) {
	if chainId == "" {
		return nil, db.ErrTableParams
	}

	blocks := make([]*db.Block, 0)
	tableName := db.GetTableName(chainId, db.TableBlock)
	where := map[string]interface{}{
		"delayUpdateStatus": status,
	}
	selectFile := &SelectFile{
		Where: where,
	}
	query := BuildParamsQuery(tableName, selectFile)
	query = query.Order("blockHeight asc")
	err := query.Find(&blocks).Error
	return blocks, err
}

// UpdateBlockUpdateStatus 更新更新状态
func UpdateBlockUpdateStatus(chainId string, blockHeight int64, updateStatus int) error {
	tableName := db.GetTableName(chainId, db.TableBlock)
	where := map[string]interface{}{
		"blockHeight": blockHeight,
	}
	params := map[string]interface{}{
		"delayUpdateStatus": updateStatus,
	}
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	return err
}
