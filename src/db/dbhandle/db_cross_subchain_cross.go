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

// InsertCrossSubChainCross 插入子链跨链详情
func InsertCrossSubChainCross(chainId string, insertList []*db.CrossSubChainCrossChain) error {
	if len(insertList) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossSubChainCrossChain)
	err := CreateInBatchesData(tableName, insertList)
	if err != nil {
		return err
	}

	// 更新缓存
	for _, v := range insertList {
		crossCache := GetCrossSubChainCrossCache(chainId, v.SubChainId)
		if len(crossCache) > 0 {
			exists := false
			for _, cross := range crossCache {
				if cross.ChainId == v.ChainId {
					exists = true
					break
				}
			}

			if !exists {
				crossCache = append(crossCache, v)
				SetCrossSubChainCrossCache(chainId, v.SubChainId, crossCache)
			}
		}
	}
	return nil
}

// UpdateCrossSubChainCross 更新子链跨链数据
func UpdateCrossSubChainCross(chainId string, subChainCross *db.CrossSubChainCrossChain) error {
	if chainId == "" || subChainCross == nil {
		return db.ErrTableParams
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossSubChainCrossChain)
	where := map[string]interface{}{
		"subChainId": subChainCross.SubChainId,
		"chainId":    subChainCross.ChainId,
	}

	params := map[string]interface{}{}
	if subChainCross.TxNum > 0 {
		params["txNum"] = subChainCross.TxNum
		params["blockHeight"] = subChainCross.BlockHeight
	}
	if len(params) == 0 {
		return nil
	}

	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	// 更新缓存
	cachedSubChainCrossList := GetCrossSubChainCrossCache(chainId, subChainCross.SubChainId)
	for _, v := range cachedSubChainCrossList {
		if v.ChainId == subChainCross.ChainId {
			v.TxNum = subChainCross.TxNum
		}
	}
	SetCrossSubChainCrossCache(chainId, subChainCross.SubChainId, cachedSubChainCrossList)
	return nil
}

// GetSubChainCrossChainList 获取子链跨链列表
func GetSubChainCrossChainList(chainId, subChainId string) ([]*db.CrossSubChainCrossChain, error) {
	subChainCrossList := make([]*db.CrossSubChainCrossChain, 0)
	if chainId == "" || subChainId == "" {
		return subChainCrossList, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossSubChainCrossChain)
	where := map[string]interface{}{
		"subChainId": subChainId,
	}
	err := db.GormDB.Table(tableName).Where(where).Order("txNum desc").Find(&subChainCrossList).Error
	return subChainCrossList, err
}

func GetCrossSubChainCrossNum(chainId string, subChainIds []string) ([]*db.CrossSubChainCrossChain, error) {
	crossSubChains := make([]*db.CrossSubChainCrossChain, 0)
	if chainId == "" || len(subChainIds) == 0 {
		return crossSubChains, nil
	}

	missingChainIds := make([]string, 0)
	for _, subChainId := range subChainIds {
		crossSubChainInfo := GetCrossSubChainCrossCache(chainId, subChainId)
		if len(crossSubChainInfo) > 0 {
			crossSubChains = append(crossSubChains, crossSubChainInfo...)
		} else {
			missingChainIds = append(missingChainIds, subChainId)
		}
	}

	if len(missingChainIds) == 0 {
		return crossSubChains, nil
	}

	findCrossSubChains := make([]*db.CrossSubChainCrossChain, 0)
	tableName := db.GetTableName(chainId, db.TableCrossSubChainCrossChain)
	err := db.GormDB.Table(tableName).Where("subChainId in ?", missingChainIds).Find(&findCrossSubChains).Error
	if err != nil {
		return crossSubChains, err
	}

	crossSubChains = append(crossSubChains, findCrossSubChains...)
	subChainMap := make(map[string][]*db.CrossSubChainCrossChain, 0)
	for _, v := range findCrossSubChains {
		subChainMap[v.SubChainId] = append(subChainMap[v.SubChainId], v)
	}

	for k, v := range subChainMap {
		SetCrossSubChainCrossCache(chainId, k, v)
	}
	return crossSubChains, nil
}

func GetCrossSubChainCrossCache(chainId, subChainId string) []*db.CrossSubChainCrossChain {
	result := make([]*db.CrossSubChainCrossChain, 0)
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossSubChainCrossChain, prefix, chainId, subChainId)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes != nil && redisRes.Val() != "" {
		err := json.Unmarshal([]byte(redisRes.Val()), &result)
		if err == nil {
			return result
		}
		log.Errorf("GetCrossSubChainCrossCache json Unmarshal err : %s, redisRes:%v", err.Error(), redisRes)
	}

	return result
}

func SetCrossSubChainCrossCache(chainId, subChainId string, crossChainList []*db.CrossSubChainCrossChain) {
	if len(crossChainList) == 0 {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossSubChainCrossChain, prefix, chainId, subChainId)
	retJson, err := json.Marshal(crossChainList)
	if err == nil {
		// 设置键值对和过期时间
		ctx := context.Background()
		_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), 30*time.Minute).Err()
	}
}
