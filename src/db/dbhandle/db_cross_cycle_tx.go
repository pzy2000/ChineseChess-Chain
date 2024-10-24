package dbhandle

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

// InsertCrossCycleTx 插入跨链交易周期数据
func InsertCrossCycleTx(chainId string, crossCycleTxs []*db.CrossCycleTransaction) error {
	if len(crossCycleTxs) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossCycleTransaction)
	err := CreateInBatchesData(tableName, crossCycleTxs)
	if err != nil {
		return err
	}

	// 删除指定redisKey的所有数据
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossLatestTransactions, prefix, chainId)
	cache.GlobalRedisDb.Del(ctx, redisKey)
	return nil
}

// UpdateCrossCycleTx 更新交易状态
func UpdateCrossCycleTx(chainId string, crossCycleTx *db.CrossCycleTransaction) error {
	if chainId == "" || crossCycleTx == nil {
		return db.ErrTableParams
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossCycleTransaction)
	where := map[string]interface{}{
		"crossId": crossCycleTx.CrossId,
	}
	params := map[string]interface{}{}
	if crossCycleTx.Status > 0 {
		params["status"] = crossCycleTx.Status
	}
	if crossCycleTx.EndTime > 0 {
		params["endTime"] = crossCycleTx.EndTime
	}
	if crossCycleTx.Duration > 0 {
		params["duration"] = crossCycleTx.Duration
	}
	if crossCycleTx.BlockHeight > 0 {
		params["blockHeight"] = crossCycleTx.BlockHeight
	}
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	// 删除指定redisKey的所有数据
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossLatestTransactions, prefix, chainId)
	cache.GlobalRedisDb.Del(ctx, redisKey)
	return nil
}

// GetCrossCycleTransactionById 根据Id获取周期交易
func GetCrossCycleTransactionById(chainId string, crossIds []string) (map[string]*db.CrossCycleTransaction, error) {
	crossCycleTxMap := make(map[string]*db.CrossCycleTransaction, 0)
	crossCycleTxs := make([]*db.CrossCycleTransaction, 0)
	if chainId == "" || len(crossIds) == 0 {
		return crossCycleTxMap, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableCrossCycleTransaction)
	err := db.GormDB.Table(tableName).Where("crossId in ?", crossIds).Find(&crossCycleTxs).Error
	if err != nil {
		return crossCycleTxMap, err
	}

	for _, cycleTx := range crossCycleTxs {
		crossCycleTxMap[cycleTx.CrossId] = cycleTx
	}
	return crossCycleTxMap, nil
}

// GetCrossCycleById 根据Id获取周期交易
func GetCrossCycleById(chainId, crossId string) (*db.CrossCycleTransaction, error) {
	tableName := db.GetTableName(chainId, db.TableCrossCycleTransaction)
	cycleTxInfo := &db.CrossCycleTransaction{}
	where := map[string]interface{}{
		"crossId": crossId,
	}
	err := db.GormDB.Table(tableName).Where(where).First(cycleTxInfo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return cycleTxInfo, nil
}

// GetCycleShortestTime 最短时间
func GetCycleShortestTime(chainId string, startTime, endTime int64) (int64, error) {
	shortestTimeCache, err := GetCrossCycleTimeCache(chainId, "ShortestTime")
	if err == nil {
		return shortestTimeCache, nil
	}

	var shortestTime int64
	tableName := db.GetTableName(chainId, db.TableCrossCycleTransaction)
	err = db.GormDB.Table(tableName).Where("startTime >= ? AND endTime <= ? AND duration > 0", startTime, endTime).
		Select("MIN(duration)").Scan(&shortestTime).Error
	if err != nil {
		return 0, err
	}
	SetCrossCycleTimeCache(chainId, "ShortestTime", shortestTime)
	return shortestTime, nil
}

// GetCycleLongestTime 最长时间
func GetCycleLongestTime(chainId string, startTime, endTime int64) (int64, error) {
	longestTimeCache, err := GetCrossCycleTimeCache(chainId, "LongestTime")
	if err == nil {
		return longestTimeCache, nil
	}

	var longestTime int64
	tableName := db.GetTableName(chainId, db.TableCrossCycleTransaction)
	err = db.GormDB.Table(tableName).Where("startTime >= ? AND endTime <= ? AND duration > 0", startTime, endTime).
		Select("MAX(duration)").Scan(&longestTime).Error
	if err != nil {
		return 0, err
	}

	SetCrossCycleTimeCache(chainId, "LongestTime", longestTime)
	return longestTime, nil
}

// GetCycleAverageTime
//
//	@Description: 跨链交易完整周期完成-平均完成时间
//	@param chainId
//	@param startTime 统计开始时间
//	@param endTime 统计结束时间
//	@return int64 平均时间
//	@return error
func GetCycleAverageTime(chainId string, startTime, endTime int64) (int64, error) {
	var result *float64 // 使用 *float64 类型
	//获取缓存数据
	averageTimeCache, err := GetCrossCycleTimeCache(chainId, "AverageTime")
	if err == nil {
		return averageTimeCache, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossCycleTransaction)
	err = db.GormDB.Table(tableName).Where("startTime >= ? AND endTime <= ? AND duration > 0", startTime, endTime).
		Select("AVG(duration)").Scan(&result).Error
	if err != nil {
		return 0, err
	}

	if result == nil || math.IsNaN(*result) {
		return 0, nil
	}
	averageTime := int64(math.Round(*result))
	if averageTime < 0 {
		log.Errorf("GetCycleAverageTime averageTime < 0")
		return 0, nil
	}

	SetCrossCycleTimeCache(chainId, "AverageTime", averageTime)
	return averageTime, nil
}

// GetCrossCycleTxAllCount 获取所有跨链交易数量
func GetCrossCycleTxAllCount(chainId string) (int64, error) {
	var count int64
	if chainId == "" {
		return count, nil
	}

	// 尝试从缓存中获取数据
	cacheKey := fmt.Sprintf(cache.CacheCrossTxCount, chainId)
	countCache, found := cache.GlobalCacheInstance.Get(cacheKey)
	if found {
		return countCache.(int64), nil
	}

	// 如果缓存中没有数据，则从数据库中查询
	tableName := db.GetTableName(chainId, db.TableCrossCycleTransaction)
	err := db.GormDB.Table(tableName).Count(&count).Error
	if err != nil {
		return 0, err
	}

	// 将查询结果存入缓存，并在失效前刷新
	cache.GlobalCacheInstance.Add(cacheKey, count)
	time.AfterFunc(5*time.Minute, func() {
		refreshCache(chainId)
	})

	return count, nil
}

// refreshCache 刷新指定chainId的缓存
func refreshCache(chainId string) {
	var count int64
	tableName := db.GetTableName(chainId, db.TableCrossCycleTransaction)
	err := db.GormDB.Table(tableName).Count(&count).Error
	if err != nil {
		return
	}
	// 更新缓存
	cacheKey := fmt.Sprintf(cache.CacheCrossTxCount, chainId)
	cache.GlobalCacheInstance.Add(cacheKey, count)

	// 重新设置失效时间
	time.AfterFunc(5*time.Minute, func() {
		refreshCache(chainId)
	})
}

func getCycleQuery(chainId string) *gorm.DB {
	tableNameCycle := db.GetTableName(chainId, db.TableCrossCycleTransaction)
	tableNameTransfer := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	query := db.GormDB.Table(tableNameCycle + " AS cct").
		Select("cct.*, ctt.*, cct.crossId as crossId").
		Joins(leftJoinKeyword + tableNameTransfer + " AS ctt ON cct.crossId = ctt.crossId")
	return query
}

// GetCrossLatestCycleTxList 获取最后10个交易列表
func GetCrossLatestCycleTxList(chainId string) ([]*db.CycleJoinTransferResult, error) {
	cacheResults := GetCrossLatestCycleTxCache(chainId)
	if cacheResults != nil {
		return cacheResults, nil
	}

	var joinResults []*db.CycleJoinTransferResult
	query := getCycleQuery(chainId)
	err := query.Order("cct.startTime DESC").Limit(10).Scan(&joinResults).Error

	if err != nil {
		return nil, err
	}

	SetCrossLatestCycleTxCache(chainId, joinResults)
	return joinResults, nil
}

// GetCrossLatestCycleTxCache 获取最后10个交易列表
func GetCrossLatestCycleTxCache(chainId string) []*db.CycleJoinTransferResult {
	ctx := context.Background()
	cacheResult := make([]*db.CycleJoinTransferResult, 0)
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossLatestTransactions, prefix, chainId)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes == nil || redisRes.Val() == "" {
		return nil
	}

	err := json.Unmarshal([]byte(redisRes.Val()), &cacheResult)
	if err != nil {
		log.Errorf("【Redis】get cache failed, key:%v, result:%v", redisKey, redisRes)
		return nil
	}
	return cacheResult
}

// SetCrossLatestCycleTxCache 缓存最后10个交易列表
func SetCrossLatestCycleTxCache(chainId string, crossCycleTxs []*db.CycleJoinTransferResult) {
	if len(crossCycleTxs) == 0 {
		return
	}

	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossLatestTransactions, prefix, chainId)
	retJson, err := json.Marshal(crossCycleTxs)
	if err != nil {
		log.Errorf("【Redis】set cache failed, key:%v, result:%v", redisKey, retJson)
		return
	}
	// 设置键值对和过期时间(1h 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), time.Hour).Err()
}

// GetCrossCycleTxList 获取交易列表
func GetCrossCycleTxList(offset, limit int, startTime, endTime int64, chainId, crossId, subChainId,
	fromChainId, toChainId string) ([]*db.CycleJoinTransferResult, int64, error) {
	var count int64
	txList := make([]*db.CycleJoinTransferResult, 0)

	query := getCycleQuery(chainId)
	if crossId != "" {
		query = query.Where("cct.crossId = ?", crossId)
	}
	if startTime != 0 && endTime != 0 {
		query = query.Where("cct.startTime BETWEEN ? AND ?", startTime, endTime)
	}

	if subChainId == "" {
		if fromChainId != "" {
			query = query.Where("ctt.fromChainId = ?", fromChainId)
		}
		if toChainId != "" {
			query = query.Where("ctt.toChainId = ?", toChainId)
		}
	} else {
		if fromChainId == "" && toChainId == "" {
			query = query.Where("ctt.fromChainId = ? OR ctt.toChainId = ?", subChainId, subChainId)
		} else if fromChainId != "" && toChainId != "" {
			query = query.Where("ctt.fromChainId = ?", fromChainId)
			query = query.Where("ctt.toChainId = ?", toChainId)
		} else {
			if fromChainId != "" {
				if fromChainId == subChainId {
					query = query.Where("ctt.fromChainId = ?", subChainId)
				} else {
					query = query.Where("ctt.fromChainId = ?", fromChainId)
					query = query.Where("ctt.toChainId = ?", subChainId)
				}
			} else if toChainId != "" {
				if toChainId == subChainId {
					query = query.Where("ctt.toChainId = ?", subChainId)
				} else {
					query = query.Where("ctt.fromChainId = ?", subChainId)
					query = query.Where("ctt.toChainId = ?", toChainId)
				}
			}
		}
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	query = query.Order("cct.startTime DESC").
		Offset(offset * limit).Limit(limit)

	err = query.Scan(&txList).Error
	if err != nil {
		return nil, 0, fmt.Errorf("GetCrossCycleTxList err, cause : %s", err.Error())
	}

	return txList, count, nil
}

// GetCrossCycleTxByHeight 根据height获取交易流转
func GetCrossCycleTxByHeight(chainId string, blockHeight []int64) ([]*db.CrossCycleTransaction, error) {
	tableName := db.GetTableName(chainId, db.TableCrossCycleTransaction)
	cycleTxs := make([]*db.CrossCycleTransaction, 0)
	err := db.GormDB.Table(tableName).Where("blockHeight IN ?", blockHeight).Find(&cycleTxs).Error
	if err != nil {
		return cycleTxs, err
	}

	return cycleTxs, nil
}

// GetCrossSubChainTxCount 获取交易列表
func GetCrossSubChainTxCount(startTime, endTime int64, chainId, crossId, subChainId, fromChainId, toChainId string) (
	int64, error) {
	var count int64

	query := getCycleQuery(chainId)
	if crossId != "" {
		query = query.Where("cct.crossId = ?", crossId)
	}
	if startTime != 0 && endTime != 0 {
		query = query.Where("cct.startTime BETWEEN ? AND ?", startTime, endTime)
	}

	if fromChainId == "" && toChainId == "" && subChainId != "" {
		query = query.Where("ctt.fromChainId = ? OR ctt.toChainId = ?", subChainId, subChainId)
	} else {
		if fromChainId != "" {
			query = query.Where("ctt.fromChainId = ?", fromChainId)
		}
		if toChainId != "" {
			query = query.Where("ctt.toChainId = ?", toChainId)
		}
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// // GetCrossCycleTxCount 获取交易列表
// func GetCrossCycleTxCount(startTime, endTime int64, chainId, crossId, fromChainId, toChainId string) (int64, error) {
// 	var count int64
// 	query := getCycleQuery(chainId)
// 	if crossId != "" {
// 		query = query.Where("cct.crossId = ?", crossId)
// 	}
// 	if startTime != 0 && endTime != 0 {
// 		query = query.Where("cct.startTime BETWEEN ? AND ?", startTime, endTime)
// 	}
// 	if fromChainId != "" {
// 		query = query.Where("ctt.fromChainId = ?", fromChainId)
// 	}
// 	if toChainId != "" {
// 		query = query.Where("ctt.toChainId = ?", toChainId)
// 	}

// 	err := query.Count(&count).Error
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// GetCrossSubChainTxList 获取交易列表
func GetCrossSubChainTxList(offset, limit int, startTime, endTime int64, chainId, crossId, subChainId,
	fromChainId, toChainId string) ([]*db.CycleJoinTransferResult, error) {
	txList := make([]*db.CycleJoinTransferResult, 0)
	query := getCycleQuery(chainId)
	if crossId != "" {
		query = query.Where("cct.crossId = ?", crossId)
	}
	if startTime != 0 && endTime != 0 {
		query = query.Where("cct.startTime BETWEEN ? AND ?", startTime, endTime)
	}

	if fromChainId == "" && toChainId == "" && subChainId != "" {
		query = query.Where("ctt.fromChainId = ? OR ctt.toChainId = ?", subChainId, subChainId)
	} else {
		if fromChainId != "" {
			query = query.Where("ctt.fromChainId = ?", fromChainId)
		}
		if toChainId != "" {
			query = query.Where("ctt.toChainId = ?", toChainId)
		}
	}

	query = query.Order("cct.startTime DESC").Offset(offset * limit).Limit(limit)
	err := query.Scan(&txList).Error
	if err != nil {
		return nil, fmt.Errorf("GetCrossSubChainTxList err, cause : %s", err.Error())
	}

	return txList, nil
}
