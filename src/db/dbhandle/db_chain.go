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
	"errors"
	"fmt"
	"time"

	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"
	"gorm.io/gorm"
)

// UpdateChainInfoByConfig 根据config更新链配置
func UpdateChainInfoByConfig(chainId string, chainConfig *pbConfig.ChainConfig) (err error) {
	where := map[string]interface{}{
		"chainId": chainId,
	}
	params := map[string]interface{}{
		"version": chainConfig.Version,
	}
	if chainConfig.Block != nil {
		params["blockInterval"] = int(chainConfig.Block.BlockInterval)
		params["blockSize"] = int(chainConfig.Block.BlockSize)
		params["blockTxCapacity"] = int(chainConfig.Block.BlockTxCapacity)
		params["txTimeout"] = int(chainConfig.Block.TxTimeout)
		params["txTimestampVerify"] = chainConfig.Block.TxTimestampVerify
	}

	if chainConfig.AccountConfig != nil {
		params["enableGas"] = chainConfig.AccountConfig.EnableGas
	}
	err = db.GormDB.Table(db.TableChain).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	// 删除指定redisKey的所有数据
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisDbChainConfig, prefix, chainId)
	cache.GlobalRedisDb.Del(ctx, redisKey)
	return nil
}

// GetChainInfoById 获取链配置
func GetChainInfoById(chainId string) (*db.Chain, error) {
	if chainId == "" {
		return nil, db.ErrTableParams
	}
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisDbChainConfig, prefix, chainId)
	chainInfo, err := GetChainInfoCache(redisKey)
	if err == nil && chainInfo != nil {
		return chainInfo, nil
	}

	chainInfo = &db.Chain{}
	where := map[string]interface{}{
		"chainId": chainId,
	}

	err = db.GormDB.Table(db.TableChain).Where(where).First(chainInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("GetChainInfoById err, cause : %s", err.Error())
	}

	//缓存数据
	retJson, err := json.Marshal(chainInfo)
	if err != nil {
		log.Errorf("GetChainInfoById json marshal err: %v，chainInfo：%v", err, chainInfo)
	} else {
		// 设置键值对和过期时间
		ctx := context.Background()
		_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), time.Hour).Err()
	}

	return chainInfo, nil
}

// GetChainInfoCache get Cache
func GetChainInfoCache(redisKey string) (*db.Chain, error) {
	//获取缓存
	var result *db.Chain
	ctx := context.Background()
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes != nil && redisRes.Val() != "" {
		err := json.Unmarshal([]byte(redisRes.Val()), &result)
		if err != nil {
			log.Errorf("GetChainInfoCache json Unmarshal err : %s, redisRes:%v", err.Error(), redisRes)
			return nil, err
		}
	}

	return result, nil
}

// GetChainListByPage get
// @desc
// @param ${param}
// @return []*ChainWithStatus
// @return int64
// @return error
func GetChainListByPage(offset, limit int, chainId string) ([]*db.Chain, int64, error) {
	var count int64
	chains := make([]*db.Chain, 0)
	where := map[string]interface{}{}
	if chainId != "" {
		where["chainId"] = chainId
	}
	err := db.GormDB.Table(db.TableChain).Where(where).Count(&count).Error
	if err != nil {
		return chains, 0, err
	}

	query := db.GormDB.Table(db.TableChain).Where(where).Order("timestamp desc")
	err = query.Offset(offset * limit).Limit(limit).Find(&chains).Error
	if err != nil {
		return chains, 0, err
	}

	return chains, count, err
}

// InsertChainInfo insert
func InsertChainInfo(chainData *db.Chain) error {
	if chainData == nil {
		return fmt.Errorf("chain is nil")
	}

	err := db.GormDB.Table(db.TableChain).Create(&chainData).Error
	if err != nil {
		log.Errorf("InsertChainInfo fail, err:%v", err)
		return err
	}
	return err
}

// InsertUpdateChainInfo insert
func InsertUpdateChainInfo(chainData *db.Chain, subscribeStatus int) error {
	if chainData == nil {
		return fmt.Errorf("chain is nil")
	}

	chainInfo, err := GetChainInfoById(chainData.ChainId)
	if err != nil {
		return err
	}

	//不存在，插入数据
	if chainInfo == nil {
		err = db.GormDB.Table(db.TableChain).Create(&chainData).Error
		if err != nil {
			log.Errorf("InsertChainInfo fail, err:%v", err)
			return err
		}
	}

	//失败的订阅不用更新chain
	if subscribeStatus != db.SubscribeOK {
		return nil
	}

	//更新链数据
	err = UpdateChainInfo(chainData)
	return err
}

// UpdateChainInfo
//
//	@Description: 更新链数据
//	@param chainInfo
//	@return error
func UpdateChainInfo(chainInfo *db.Chain) error {
	where := map[string]interface{}{
		"chainId": chainInfo.ChainId,
	}

	params := map[string]interface{}{
		"version":           chainInfo.Version,
		"enableGas":         chainInfo.EnableGas,
		"blockInterval":     chainInfo.BlockInterval,
		"authType":          chainInfo.AuthType,
		"hashType":          chainInfo.HashType,
		"blockSize":         chainInfo.BlockSize,
		"txTimestampVerify": chainInfo.TxTimestampVerify,
		"txTimeout":         chainInfo.TxTimeout,
		"consensus":         chainInfo.Consensus,
		"blockTxCapacity":   chainInfo.BlockTxCapacity,
	}

	return db.GormDB.Table(db.TableChain).Where(where).Updates(params).Error
}

// DeleteChain delete
// @desc
// @param ${param}
// @return error
func DeleteChain(chainId string) error {
	chainInfo := &db.Chain{}
	where := map[string]interface{}{
		"chainId": chainId,
	}
	err := db.GormDB.Table(db.TableChain).Where(where).Delete(chainInfo).Error
	if err != nil {
		log.Error("[DB] Delete Chain Failed: " + err.Error())
	}
	return err
}
