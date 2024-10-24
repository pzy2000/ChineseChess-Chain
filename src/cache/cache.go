// Package cache provides cache Methods
package cache

import (
	"chainmaker_web/src/config"
	"context"
	"fmt"

	lru "github.com/hashicorp/golang-lru"
	"github.com/redis/go-redis/v9"
)

var (
	//GlobalRedisDb redis连接
	GlobalRedisDb redis.Cmdable
	// 定义全局缓存实例
	GlobalCacheInstance *lru.Cache
)

const (
	//RedisSubscribeLockKey 分布式订阅锁
	RedisSubscribeLockKey = "%s_%s_redis_subscribe_lock_key"
	//RedisDelayedUpdateData 异步更新数据缓存
	RedisDelayedUpdateData = "%s_%s_delay_update_data_%s"
	//RedisDbChainConfig  链配置缓存
	RedisDbChainConfig = "%s_%s_db_chain_config"
	//RedisContractInfoByNameAddr 合约缓存
	RedisContractInfoByNameAddr = "%s_%s_contract_info_by_name_addr_%s"
	//RedisDBAccountData 账户信息缓存
	RedisDBAccountData = "%s_%s_db_account_data_%s"
	//RedisUserMemberInfoKey user缓存
	RedisUserMemberInfoKey = "%s_%s_user_member_info_key_%s_%d"
)

const (
	//RedisLatestTransactions 最新交易缓存
	RedisLatestTransactions = "%s_%s_latest_transaction_list"
	//RedisLatestBlockList 最新区块缓存
	RedisLatestBlockList = "%s_%s_latest_block_list"
	//RedisLatestContractList 最新合约列表缓存
	RedisLatestContractList = "%s_%s_latest_contract_list"
	//RedisHomeOverviewData 首页统计数据缓存
	RedisHomeOverviewData = "%s_%s_home_overview_data"
	//RedisOverviewContractCount 首页合约数量缓存
	RedisOverviewContractCount = "%s_%s_overview_contract_count"
	//RedisOverviewNodeCount 首页节点数据缓存
	RedisOverviewNodeCount = "%s_%s_overview_node_count_%s"
	//RedisOverviewTxTotal 首页交易总量缓存
	RedisOverviewTxTotal = "%s_%s_overview_tx_total"
	//RedisOverviewMaxBlockHeight 首页最大区块高度缓存
	RedisOverviewMaxBlockHeight = "%s_%s_overview_max_block_height"
	//RedisOverviewTxNumTime 首页24小时交易量缓存
	RedisOverviewTxNumTime = "%s_%s_overview_tx_num_time_%s_%s_%d"
)

const (
	//RedisCrossSubChainName 子链
	RedisCrossSubChainName = "%s_%s_cross_sub_chain_name_%s"
	//RedisCrossOverviewData 主子链首页统计数据缓存
	RedisCrossOverviewData = "%s_%s_cross_home_overview_data"
	//RedisCrossCycleTxDurationMonth 跨链交易月度耗时统计缓存
	RedisCrossCycleTxDurationMonth = "%s_%s_cross_cycle_tx_duration_month_%s_%s"
	//RedisCrossLatestTransactions 最新跨链交易缓存
	RedisCrossLatestTransactions = "%s_%s_cross_latest_tx_list"
	//RedisCrossLatestSubChainList 最新子链列表缓存
	RedisCrossLatestSubChainList = "%s_%s_cross_latest_sub_chain_list"
	//RedisCrossSubChainData 跨链子链缓存
	RedisCrossSubChainData = "%s_%s_cross_sub_chain_data_%s"
	//RedisCrossContractCount 跨链子链缓存
	RedisCrossContractCount = "%s_%s_cross_contract_count_%s"
	//RedisCrossTxTransfers 新增跨链流转缓存，异步更新交易数据使用
	RedisCrossTxTransfers = "%s_%s_cross_tx_transfers_%s"
	//RedisCrossCycleTxData 跨链交易状态
	RedisCrossCycleTxData = "%s_%s_cross_cycle_tx_data_%s"
	//RedisCrossSubChainCrossChain 子链跨链数据
	RedisCrossSubChainCrossChain = "%s_%s_cross_sub_chain_cross_chain_%s"
)

const (
	//RedisTransactionListCount 交易列表交易数量
	RedisTransactionListCount = "%s_%s_transaction_list_count_%s"
	//RedisContractEventListCount 合约事件数
	RedisContractEventListCount = "%s_%s_contract_event_list_count_%s"
)

const (
	//RedisUserAddressGasInfo gas数
	RedisUserAddressGasInfo = "%s_%s_user_address_gas_Info_%s"
	//RedisFTContractPositionData 同质化合约持仓数据
	RedisFTContractPositionData = "%s_%s_ft_contract_position_data_%s"
	//RedisNFTContractPositionData 非同质化合约持仓数据
	RedisNFTContractPositionData = "%s_%s_nft_contract_position_data_%s"
	//RedisFTContractData 同质化合约数据
	RedisFTContractData = "%s_%s_ft_contract_data_%s"
	//RedisNFTContractData 非同质化合数据
	RedisNFTContractData = "%s_%s_nft_contract_data_%s"
	//RedisContractPositionOwnerList 合约持仓用户列表数据
	RedisContractPositionOwnerList = "%s_%s_contract_position_owner_list_%s"
)

const (
	CacheCrossTxCount = "cache_cross_tx_count_%s"
)

// InitRedis redis初始化
func InitRedis(cfg *config.RedisConfig) {
	var cmd *redis.Cmd
	if cfg.Type == "cluster" {
		//集群模式
		clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    cfg.Host,
			Password: cfg.Password, // 如果您的集群需要密码，请在此设置
			Username: cfg.Username,
		})
		// 使用 Do 方法发送 PING 命令并获取原始响应
		cmd = clusterClient.Do(context.Background(), "PING")
		GlobalRedisDb = clusterClient
	} else {
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.Host[0],
			Password: cfg.Password,
			Username: cfg.Username,
		})
		// 使用 Do 方法发送 PING 命令并获取原始响应
		cmd = client.Do(context.Background(), "PING")
		GlobalRedisDb = client
	}
	if cmd == nil {
		panic("InitRedis failed cmd is nil")
	}

	_, err := cmd.Result()
	if err != nil {
		panic(fmt.Sprintf("InitRedis failed %s ", err.Error()))
	}

	//设置内存缓存
	GlobalCacheInstance, _ = lru.New(1024) // 缓存大小设置为1024
}
