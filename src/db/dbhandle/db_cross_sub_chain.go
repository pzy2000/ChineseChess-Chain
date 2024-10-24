package dbhandle

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	//SubChainStatusSuccess 子链健康状态正常
	SubChainStatusSuccess int32 = 0
	//SubChainStatusFail 子链健康状态异常
	SubChainStatusFail int32 = 1
)

// GetCrossSubChainName 获取子链名称
func GetCrossSubChainName(chainId, subChainId string) (string, error) {
	nameCache, err := GetCrossSubChainNameCache(chainId, subChainId)
	if err == nil && nameCache != "" {
		return nameCache, nil
	}
	subChainInfo, err := GetCrossSubChainInfoById(chainId, subChainId)
	if err != nil || subChainInfo == nil {
		return "", err
	}
	SetCrossSubChainNameCache(chainId, subChainId, subChainInfo.ChainName)
	return subChainInfo.ChainName, nil
}

// GetCrossSubChainNameCache 获取子链名称缓存
func GetCrossSubChainNameCache(chainId, subChainId string) (string, error) {
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossSubChainName, prefix, chainId, subChainId)
	result, err := cache.GlobalRedisDb.Get(ctx, redisKey).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

// SetCrossSubChainNameCache 缓存子链名称
func SetCrossSubChainNameCache(chainId, subChainId, subChainName string) {
	if chainId == "" || subChainId == "" || subChainName == "" {
		return
	}

	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossSubChainName, prefix, chainId, subChainId)
	// 设置键值对和过期时间(1h 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, subChainName, 24*time.Hour).Err()
}

// InsertCrossSubChain 插入子链详情
func InsertCrossSubChain(chainId string, subChainList []*db.CrossSubChainData) error {
	if len(subChainList) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	err := CreateInBatchesData(tableName, subChainList)
	if err != nil {
		return err
	}
	return nil
}

// GetCrossSubChainAll 获取所有子链
func GetCrossSubChainAll(chainId string) ([]*db.CrossSubChainData, error) {
	crossSubChains := make([]*db.CrossSubChainData, 0)
	if chainId == "" {
		return crossSubChains, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	err := db.GormDB.Table(tableName).Find(&crossSubChains).Error
	if err != nil {
		return crossSubChains, err
	}

	return crossSubChains, nil
}

// UpdateCrossSubChainStatus 更新子链健康状态
func UpdateCrossSubChainStatus(chainId, subChainId, spvContractName string, status int32) error {
	if chainId == "" {
		return db.ErrTableParams
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	where := map[string]interface{}{
		"subChainId": subChainId,
	}

	updateTime := time.Now().Unix()
	params := map[string]interface{}{
		"status":    status,
		"timestamp": updateTime,
	}
	if spvContractName != "" {
		params["spvContractName"] = spvContractName
	}
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateCrossSubChainById 更新子链信息
func UpdateCrossSubChainById(chainId string, subChainInfo *db.CrossSubChainData) error {
	if chainId == "" || subChainInfo == nil {
		return db.ErrTableParams
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	where := map[string]interface{}{
		"subChainId": subChainInfo.SubChainId,
	}
	params := map[string]interface{}{}
	if subChainInfo.GatewayId != "" {
		params["gatewayId"] = subChainInfo.GatewayId
		params["gatewayName"] = subChainInfo.GatewayName
		params["gatewayAddr"] = subChainInfo.GatewayAddr
	}

	if subChainInfo.BlockHeight > 0 {
		params["blockHeight"] = subChainInfo.BlockHeight
	}
	if subChainInfo.SpvContractName != "" {
		params["spvContractName"] = subChainInfo.SpvContractName
	}
	if subChainInfo.TxNum > 0 {
		params["txNum"] = subChainInfo.TxNum
	}
	if subChainInfo.CrossCa != "" {
		params["crossCa"] = subChainInfo.CrossCa
	}
	if subChainInfo.SdkClientCrt != "" {
		params["sdkClientCrt"] = subChainInfo.SdkClientCrt
	}
	if subChainInfo.SdkClientKey != "" {
		params["sdkClientKey"] = subChainInfo.SdkClientKey
	}
	if len(params) > 0 {
		params["timestamp"] = time.Now().Unix()
	}
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	//删除缓存
	DeleteCrossSubChainCache(chainId, subChainInfo.SubChainId)
	return nil
}

// UpdateCrossChainHeightBySpv
//
//	@Description: 根据spv合约更新子链区块高度
//	@param chainId 主链id
//	@param subChainInfo 子链信息
//	@return error
func UpdateCrossChainHeightBySpv(chainId string, subChainInfo *db.CrossSubChainData) error {
	if chainId == "" || subChainInfo == nil {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	where := map[string]interface{}{
		"spvContractName": subChainInfo.SpvContractName,
	}
	params := map[string]interface{}{}
	if subChainInfo.BlockHeight > 0 {
		params["blockHeight"] = subChainInfo.BlockHeight
	}
	if len(params) > 0 {
		params["timestamp"] = time.Now().Unix()
	}

	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// GetCrossSubChainById 根据chainId获取子链信息
func GetCrossSubChainById(chainId string, subChainIds []string) (map[string]*db.CrossSubChainData, error) {
	crossSubChainMap := make(map[string]*db.CrossSubChainData, 0)
	crossSubChains := make([]*db.CrossSubChainData, 0)
	if chainId == "" || len(subChainIds) == 0 {
		return crossSubChainMap, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	err := db.GormDB.Table(tableName).Where("subChainId in ?", subChainIds).Find(&crossSubChains).Error
	if err != nil {
		return crossSubChainMap, err
	}

	for _, subChain := range crossSubChains {
		crossSubChainMap[subChain.SubChainId] = subChain
	}
	return crossSubChainMap, nil
}

// GetCrossSubChainInfoByName 根据Name获取子链信息
func GetCrossSubChainInfoByName(chainId string, subChainName string) (*db.CrossSubChainData, error) {
	cacheResult := GetCrossSubChainInfoCache(chainId, subChainName)
	if cacheResult != nil {
		return cacheResult, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	subChainInfo := &db.CrossSubChainData{}
	where := map[string]interface{}{
		"chainName": subChainName,
	}
	err := db.GormDB.Table(tableName).Where(where).First(subChainInfo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	SetCrossSubChainInfoCache(chainId, subChainName, subChainInfo)
	return subChainInfo, nil
}

// GetCrossSubChainInfoById 根据Id获取子链信息
func GetCrossSubChainInfoById(chainId, subChainId string) (*db.CrossSubChainData, error) {
	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	subChainInfo := &db.CrossSubChainData{}
	if subChainId == "" {
		return subChainInfo, nil
	}

	where := map[string]interface{}{
		"subChainId": subChainId,
	}
	err := db.GormDB.Table(tableName).Where(where).First(subChainInfo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	//设置子链缓存
	SetCrossSubChainInfoCache(chainId, subChainId, subChainInfo)
	return subChainInfo, nil
}

// GetAllSubChainBlockHeight 获取所有子链高度
func GetAllSubChainBlockHeight(chainId string) (int64, error) {
	if chainId == "" {
		return 0, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	var totalBlockHeight int64
	err := db.GormDB.Table(tableName).Select("sum(blockHeight)").Scan(&totalBlockHeight).Error
	if err != nil {
		return 0, nil
	}

	return totalBlockHeight, nil
}

// GetCrossSubChainAllCount 获取所有子链数量
func GetCrossSubChainAllCount(chainId string) (int64, error) {
	var count int64
	if chainId == "" {
		return count, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	err := db.GormDB.Table(tableName).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetCrossSubChainInfoCache 获取子链信息缓存
func GetCrossSubChainInfoCache(chainId, subChainId string) *db.CrossSubChainData {
	ctx := context.Background()
	cacheResult := &db.CrossSubChainData{}
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossSubChainData, prefix, chainId, subChainId)
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

// SetCrossSubChainInfoCache 缓存子链信息
func SetCrossSubChainInfoCache(chainId, subChainId string, subChainInfo *db.CrossSubChainData) {
	if subChainInfo == nil {
		return
	}

	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossSubChainData, prefix, chainId, subChainId)
	retJson, err := json.Marshal(subChainInfo)
	if err != nil {
		log.Errorf("【Redis】set cache failed, key:%v, result:%v", redisKey, retJson)
		return
	}
	// 设置键值对和过期时间(12h 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), 12*time.Hour).Err()
}

// DeleteCrossSubChainCache 删除子链缓存信息
func DeleteCrossSubChainCache(chainId, subChainId string) {
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossSubChainData, prefix, chainId, subChainId)
	_ = cache.GlobalRedisDb.Del(ctx, redisKey).Err()
}

// GetCrossLatestSubChainList 获取最新10条子链
func GetCrossLatestSubChainList(chainId string) ([]*db.CrossSubChainData, error) {
	crossSubChains := make([]*db.CrossSubChainData, 0)
	if chainId == "" {
		return crossSubChains, db.ErrTableParams
	}

	cacheResult := GetCrossSubChainListCache(chainId)
	if cacheResult != nil {
		return cacheResult, nil
	}
	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	err := db.GormDB.Table(tableName).Order("timestamp DESC").Limit(10).Find(&crossSubChains).Error
	if err != nil {
		return crossSubChains, err
	}

	SetCrossSubChainListCache(chainId, crossSubChains)
	return crossSubChains, nil
}

// GetCrossSubChainListCache 获取子链信息缓存
func GetCrossSubChainListCache(chainId string) []*db.CrossSubChainData {
	ctx := context.Background()
	cacheResult := make([]*db.CrossSubChainData, 0)
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossLatestSubChainList, prefix, chainId)
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

// SetCrossSubChainListCache 缓存子链信息
func SetCrossSubChainListCache(chainId string, subChainList []*db.CrossSubChainData) {
	if len(subChainList) == 0 {
		return
	}

	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossLatestSubChainList, prefix, chainId)
	retJson, err := json.Marshal(subChainList)
	if err != nil {
		log.Errorf("【Redis】set cache failed, key:%v, result:%v", redisKey, retJson)
		return
	}
	// 设置键值对和过期时间(1h 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), time.Minute).Err()
}

// GetCrossSubChainList 获取子链列表
func GetCrossSubChainList(offset, limit int, chainId, subChainId, chainName string) (
	[]*db.CrossSubChainData, int64, error) {
	var count int64
	subChainList := make([]*db.CrossSubChainData, 0)
	tableName := db.GetTableName(chainId, db.TableCrossSubChainData)
	where := map[string]interface{}{}
	if subChainId != "" {
		where["subChainId"] = subChainId
	}
	if chainName != "" {
		where["chainName"] = chainName
	}
	query := db.GormDB.Table(tableName).Where(where)
	err := query.Count(&count).Error
	if err != nil {
		return subChainList, count, err
	}
	query = query.Order("timestamp desc").Offset(offset * limit).Limit(limit)
	err = query.Find(&subChainList).Error
	return subChainList, count, err
}
