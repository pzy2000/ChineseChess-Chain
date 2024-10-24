package dbhandle

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"context"
	"fmt"
	"strconv"
	"time"
)

// InsertCrossContract 插入跨链合约
func InsertCrossContract(chainId, subChainId string, insertList []*db.CrossChainContract) error {
	if chainId == "" || subChainId == "" || len(insertList) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossChainContract)
	err := CreateInBatchesData(tableName, insertList)
	if err != nil {
		return err
	}

	//删除缓存
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossContractCount, prefix, chainId, subChainId)
	_ = cache.GlobalRedisDb.Del(ctx, redisKey).Err()
	return nil
}

func GetCrossContractByName(chainId, subChainId string, nameList []string) ([]*db.CrossChainContract, error) {
	crossContracts := make([]*db.CrossChainContract, 0)
	if chainId == "" || subChainId == "" || len(nameList) == 0 {
		return crossContracts, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossChainContract)

	where := map[string]interface{}{
		"subChainId": subChainId,
	}
	whereIn := map[string]interface{}{
		"contractName": nameList,
	}

	selectFile := &SelectFile{
		Where:   where,
		WhereIn: whereIn,
	}
	query := BuildParamsQuery(tableName, selectFile)
	err := query.Find(&crossContracts).Error
	return crossContracts, err
}

// GetCrossContractCount 根据子链id获取合约数
func GetCrossContractCount(chainId, subChainId string) (int64, error) {
	contractCount, err := GetCrossContractCountCache(chainId, subChainId)
	if err == nil && contractCount > 0 {
		return contractCount, nil
	}

	if chainId == "" || subChainId == "" {
		return contractCount, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossChainContract)
	where := map[string]interface{}{
		"subChainId": subChainId,
	}
	err = db.GormDB.Table(tableName).Where(where).Count(&contractCount).Error
	if err != nil {
		return contractCount, err
	}

	SetCrossContractCountCache(chainId, subChainId, contractCount)
	return contractCount, err
}

// GetCrossContractCountCache
//
//	@Description: 获取跨链合约数量
//	@param chainId 主链id
//	@param subChainId 子链id
//	@return int64 跨链合约数量
//	@return error
func GetCrossContractCountCache(chainId, subChainId string) (int64, error) {
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossContractCount, prefix, chainId, subChainId)
	redisRes, err := cache.GlobalRedisDb.Get(ctx, redisKey).Result()
	if err != nil {
		return 0, err
	}
	count, err := strconv.ParseInt(redisRes, 10, 64)
	if err != nil {
		log.Errorf("【Redis】get cache failed, key:%v, result:%v", redisKey, redisRes)
		return 0, err
	}

	return count, nil
}

// SetCrossContractCountCache
//
//	@Description: 缓存子链跨链合约数量
//	@param chainId 主链id
//	@param subChainId 子链ID
//	@param contractCount 跨链合约数量
func SetCrossContractCountCache(chainId, subChainId string, contractCount int64) {
	if contractCount == 0 {
		return
	}

	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossContractCount, prefix, chainId, subChainId)
	// 设置键值对和过期时间(24h 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, contractCount, 24*time.Hour).Err()
}
