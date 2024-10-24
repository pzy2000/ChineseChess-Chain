/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

// nolint
import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	loggers "chainmaker_web/src/logger"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
)

var (
	log              = loggers.GetLogger(loggers.MODULE_SYNC)
	sdkClientPool    = NewSdkClientPool()
	StartSyncCancels sync.Map // 使用 sync.Map 代替锁
)

// Start
//
//	@Description: 开始链订阅
//	@param chainList 链列表
func StartSync(chainList []*config.ChainInfo) {
	for _, chainInfo := range chainList {
		// 订阅链
		go func(chainInfo *config.ChainInfo) {
			// 创建新的 context 和取消函数
			ctx, cancel := context.WithCancel(context.Background())
			StartSyncCancels.Store(chainInfo.ChainId, cancel)

			for {
				select {
				case <-ctx.Done():
					log.Infof("【sync】区块链【%v】订阅任务已取消", chainInfo.ChainId)
					return
				default:
					err := SubscribeChain(chainInfo)
					if err == nil {
						log.Infof("【sync】区块链【%v】订阅成功", chainInfo.ChainId)
						return
					}
					log.Errorf("【sync】区块链【%v】订阅失败, 正在尝试重新订阅", chainInfo.ChainId)
					time.Sleep(time.Second * 10) // 添加一个短暂延迟，避免频繁重试
				}
			}
		}(chainInfo)
	}
}

func SubscribeChain(chainInfo *config.ChainInfo) error {
	chainId := chainInfo.ChainId
	//判断连接池中是否存在，已经存在的话需要先先停止之前的订阅后在重新订阅
	poolSdkClient := GetSdkClient(chainId)
	if poolSdkClient != nil {
		//停止订阅
		StopChain(poolSdkClient)
	}

	// 创建连接
	chainConfig, subscribeErr := subscribe(chainInfo)
	subscribeStatus := db.SubscribeOK
	if subscribeErr != nil {
		log.Errorf("【sync】区块链【%v】连接失败: %v, 尝试重新订阅...", chainId, subscribeErr)
		subscribeStatus = db.SubscribeFailed
	} else {
		log.Infof("【sync】区块链【%v】连接成功, 开启订阅", chainInfo.ChainId)
		//开启订阅
		sdkClientPool.LoadChains(chainId)
	}

	//更新订阅状态
	errDB := SaveSubscribeToDB(chainInfo, chainConfig, subscribeStatus)
	if errDB != nil {
		log.Errorf("【sync】 SaveSubscribeToDB failed, err:%v", errDB)
	}

	if subscribeErr != nil || errDB != nil {
		return fmt.Errorf("【sync】SubscribeChain fialed")
	}

	return nil
}

// subscribe
//
//	@Description:  在这里执行订阅操作，如果订阅成功，返回nil，否则返回错误
//	@param chainInfo
//	@return error
func subscribe(chainInfo *config.ChainInfo) (*pbconfig.ChainConfig, error) {
	chainInfoJson, _ := json.Marshal(chainInfo)
	log.Infof("【Sync】 chainId[%v] init sdk clients Start, chainInfoJson:%v", chainInfo.ChainId, string(chainInfoJson))

	chainClient, err := CreateChainClient(chainInfo)
	sdkClient := NewSdkClient(chainInfo, chainClient)
	if err != nil {
		log.Errorf("【Sync】创建chain Client失败: err:%v, chainInfo:%v",
			err.Error(), string(chainInfoJson))
		return nil, err
	}
	//判断节点是否存活
	chainConfig, err := sdkClient.ChainClient.GetChainConfig()
	sdkClient.ChainConfig = chainConfig
	if err != nil {
		_ = chainClient.Stop()
		log.Errorf("【Sync】try to connect chain failed:err:%v , chainInfo:%v",
			err.Error(), string(chainInfoJson))
		return nil, err
	}

	// 加入到连接池
	clientPool := NewSingleSdkClientPool(chainInfo, sdkClient, chainClient)
	sdkClientPool.addSdkClientPool(clientPool)

	log.Infof("【Sync】 chainId[%v] init sdk clients success", chainInfo.ChainId)
	return chainConfig, nil
}

// SaveSubscribeToDB
//
//	@Description: 将订阅信息存储数据库
//	@param chainInfo
//	@param chainConfig
//	@return error
func SaveSubscribeToDB(chainInfo *config.ChainInfo, chainConfig *pbconfig.ChainConfig, subscribeStatus int) error {
	//更新订阅状态
	err := dbhandle.InsertOrUpdateSubscribe(chainInfo, subscribeStatus)
	log.Infof("【sync】 save Subscribe finished, err:%v, chain:%v, subscribeStatus:%v",
		err, chainInfo, subscribeStatus)
	if err != nil {
		return err
	}

	//处理链数据
	chain := &db.Chain{}
	if chainConfig == nil {
		chain.ChainId = chainInfo.ChainId
		chain.AuthType = chainInfo.AuthType
		chain.Timestamp = time.Now().Unix()
	} else {
		chain = &db.Chain{
			ChainId:           chainInfo.ChainId,
			Version:           chainConfig.Version,
			EnableGas:         chainConfig.AccountConfig.EnableGas,
			BlockInterval:     int(chainConfig.Block.BlockInterval),
			BlockSize:         int(chainConfig.Block.BlockSize),
			BlockTxCapacity:   int(chainConfig.Block.BlockTxCapacity),
			TxTimestampVerify: chainConfig.Block.TxTimestampVerify,
			TxTimeout:         int(chainConfig.Block.TxTimeout),
			Consensus:         chainConfig.Consensus.Type.String(),
			HashType:          chainConfig.Crypto.Hash,
			AuthType:          chainInfo.AuthType,
			Timestamp:         time.Now().Unix(),
		}
	}

	//插入，更新链信息
	err = dbhandle.InsertUpdateChainInfo(chain, subscribeStatus)
	log.Infof("【sync】 save chaininfo finished, err:%v, chain:%v", err, chain)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSubscribeToDB
//
//	@Description: 将订阅信息存储数据库
//	@param chainInfo
//	@param chainConfig
//	@return error
// func UpdateSubscribeToDB(chainInfo *config.ChainInfo, chainConfig *pbconfig.ChainConfig, subscribeStatus int) error {
// 	//更新订阅状态
// 	subscribeInfo := dbhandle.BuildSubscribeInfo(chainInfo, subscribeStatus)
// 	err := dbhandle.UpdateSubscribe(subscribeInfo)
// 	log.Infof("【sync】 save Subscribe finished, err:%v, chain:%v, subscribeStatus:%v",
// 		err, chainInfo, subscribeStatus)
// 	if err != nil {
// 		return err
// 	}

// 	if chainConfig == nil {
// 		return nil
// 	}

// 	//处理链数据
// 	chain := &db.Chain{
// 		ChainId:           chainInfo.ChainId,
// 		Version:           chainConfig.Version,
// 		EnableGas:         chainConfig.AccountConfig.EnableGas,
// 		BlockInterval:     int(chainConfig.Block.BlockInterval),
// 		BlockSize:         int(chainConfig.Block.BlockSize),
// 		BlockTxCapacity:   int(chainConfig.Block.BlockTxCapacity),
// 		TxTimestampVerify: chainConfig.Block.TxTimestampVerify,
// 		TxTimeout:         int(chainConfig.Block.TxTimeout),
// 		Consensus:         chainConfig.Consensus.Type.String(),
// 		HashType:          chainConfig.Crypto.Hash,
// 		AuthType:          chainInfo.AuthType,
// 	}

// 	//插入，更新链信息
// 	err = dbhandle.UpdateChainInfo(chain)
// 	log.Infof("【sync】 save chaininfo finished, err:%v, chain:%v", err, chain)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
