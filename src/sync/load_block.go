/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

// BlockWaitUpdate
// @Description: 异步更新channel数据
type BlockWaitUpdate struct {
	ChainId     string
	BlockHeight int64
}

// DataSaveToDB
// @Description: 处理完的订阅数据,放入存储channel,顺序插入数据库
type DataSaveToDB struct {
	ChainId     string
	BlockHeight int64
	DealResult  RealtimeDealResult
	TxTimeLog   *TxTimeLog
}

// PeriodicGetSubscribeLock
//
//	@Description:  10分钟请求一次,获取订阅锁,获取到锁才进行订阅
//	@param sdkClient
func PeriodicGetSubscribeLock(sdkClient *SdkClient) {
	ctx := sdkClient.Ctx
	chainId := sdkClient.ChainId
	prefix := config.GlobalConfig.RedisDB.Prefix
	lockKey := fmt.Sprintf(cache.RedisSubscribeLockKey, prefix, chainId)
	// 尝试获取分布式锁（第一次尝试）
	lock := cache.GlobalRedisDb.SetNX(ctx, lockKey, chainId, 3*time.Minute)
	if lock.Val() {
		log.Infof("【load】Periodic Get Subscribe Lock (first attempt)【true】, LockKey:%s", lockKey)
		// 获取到锁,说明其他节点订阅失败了,启动订阅
		err := blockListen(sdkClient)
		if err != nil {
			//重启链
			ReStartChain(sdkClient)
		}
	} else {
		log.Infof("【load】Periodic Get Subscribe Lock (first attempt)【false】, LockKey:%s", lockKey)
		// 如果第一次尝试失败,启动定时器
		startSubscribeLockTicker(ctx, sdkClient, lockKey)
	}
}

// startSubscribeLockTicker
//
//	@Description: 当获取锁失败,有其他机器订阅时,启动订阅定时器,知道订阅成功,停止定时器
//	@param sdkClient
//	@param LockKey
func startSubscribeLockTicker(ctx context.Context, sdkClient *SdkClient, lockKey string) {
	//1分钟定时器,获取订阅锁
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop() // 在函数返回时停止定时器

	for {
		select {
		case <-ticker.C:
			//链订阅已经停止,停止定时器
			if sdkClient.Status == STOP {
				return
			}

			// 尝试获取分布式锁
			lock := cache.GlobalRedisDb.SetNX(ctx, lockKey, sdkClient.ChainId, 3*time.Minute)
			log.Infof("【load】Periodic Get Subscribe Lock【%v】, LockKey:%s", lock.Val(), lockKey)
			if lock.Val() {
				// 获取到锁,说明其他节点订阅失败了,启动订阅
				// 订阅链
				err := blockListen(sdkClient)
				if err != nil {
					//重启链
					ReStartChain(sdkClient)
					return
				}
			}
		case <-ctx.Done():
			// 当接收到上下文取消的通知时,返回函数,停止定时器
			return
		}
	}
}

// blockListen
//
//	@Description: 订阅区块数据
//	@param sdkClient 链连接
//	@return error
func blockListen(sdkClient *SdkClient) error {
	log.Infof("【load】 begin to subscribe [chain:%s] ", sdkClient.ChainId)
	//判断连接池还在不在,不在的话不在重启链
	poolSdkClient := GetSdkClient(sdkClient.ChainId)
	if poolSdkClient == nil {
		log.Infof("【ReStartChain】poolSdkClient is null, chain is cancel,chainId:%v", sdkClient.ChainId)
		return nil
	}

	var (
		//err error
		err error
		//chainId chainId
		chainId = sdkClient.ChainId
		// 使用sdkClient的Context
		ctx = sdkClient.Ctx
		//blockInfoCh 创建一个缓冲通道来存储格式化后的订阅数据 blockInfo
		blockInfoCh = make(chan *common.BlockInfo, config.BlockInsertWorkerCount)
		//blockWaitUpdateCh 创建一个容量为 20 的通道来存储处理完成,等待异步更新的区块数据
		blockWaitUpdateCh = make(chan *BlockWaitUpdate, config.BlockWaitUpdateWorkerCount)
		//dataSaveCh 创建一个无缓冲的通道,用来顺序执行插入操作 ,因为需要确保区块是顺序插入的,所以采用无缓存通道。
		dataSaveCh = make(chan *DataSaveToDB)
		//blockListenErrCh 创建一个错误通道来接收 子线程 的错误
		blockListenErrCh = make(chan error)
	)

	defer sdkClient.Cancel() // 在函数退出时调用Cancel
	defer close(blockInfoCh)
	defer close(blockWaitUpdateCh)
	defer close(dataSaveCh)
	//defer close(blockListenErrCh)

	//消费blockWaitUpdateCh队列,
	//启动异步更新协成,计算交易数量,持仓信息等区块数据
	go DelayUpdateOperation(ctx, blockWaitUpdateCh, blockListenErrCh)

	//写入blockWaitUpdateCh队列,
	//将上一次未更新数据写入通道,继续更新操作
	err = waitUpdateChFailedData(ctx, chainId, blockWaitUpdateCh)
	if err != nil {
		return err
	}

	//消费dataSaveCh队列,写入blockWaitUpdateCh队列
	//按顺序插入DB,区块数据,需要确保区块是顺序插入的,写入异步处理队列blockWaitUpdateCh
	go OrderedSaveBlockData(ctx, dataSaveCh, blockWaitUpdateCh, blockListenErrCh)

	//写入blockInfoCh队列
	//订阅数据,将订阅数据处理成结构化数据
	go SubscribeBlockSetToBlockInfoCh(ctx, sdkClient, blockInfoCh, blockListenErrCh)

	hash := sdkClient.GetChainHashType()
	//消费blockInfoCh队列,写入dataSaveCh队列
	//将格式话的区块数据,处理成可以插入DB的数据格式
	go RealtimeInsertOperation(ctx, hash, blockInfoCh, dataSaveCh, blockListenErrCh)

	// 使用 select 语句等待错误或上下文取消
	select {
	case errCh := <-blockListenErrCh:
		// 接收区块错误,无法处理成结构化数据,重启链
		log.Errorf("【sync block】Subscribe block failed, err:%v", errCh)
		return errCh
	case <-ctx.Done():
		log.Errorf("【sync block】Subscribe block failed, context cancel, err:%v", ctx.Err())
		// 上下文已取消,停止监听
		return ctx.Err()
	}
}

// OrderedSaveBlockData
//
//	@Description:  将订阅处理完的数据,按顺序插入DB,采用channel确保区块是顺序插入的
//	@param ctx
//	@param dataSaveCh 按照blockHeight顺序写入dataSaveCh,保证插入顺序
//	@param blockWaitUpdateCh 插入完成后写入blockWaitUpdateCh通道,等待更新数据
//	@param errCh
func OrderedSaveBlockData(ctx context.Context, dataSaveCh chan *DataSaveToDB, blockWaitUpdateCh chan *BlockWaitUpdate,
	errCh chan<- error) {
	for data := range dataSaveCh {
		chainId := data.ChainId
		blockHeight := data.BlockHeight
		log.Infof("【Realtime insert】start block-%s[%d]", chainId, blockHeight)
		startTime := time.Now()
		err := RealtimeDataSaveToDB(chainId, blockHeight, data.DealResult, data.TxTimeLog)
		log.Infof("【Realtime insert】end block-%s[%d], duration_time(ms):%v", chainId, blockHeight,
			time.Since(startTime).Milliseconds())
		if err != nil {
			errCh <- fmt.Errorf("【Realtime insert】err block-%s[%d] failed, err:%v", chainId, blockHeight, err)
			// 如果处理失败,取消上下文,会重启链
			return
		}

		// 将处理完成的结果写入 blockWaitUpdateCh
		resultData := &BlockWaitUpdate{
			ChainId:     chainId,
			BlockHeight: blockHeight,
		}
		select {
		case <-ctx.Done():
			// 上下文已取消,不要发送数据到通道
			return
		case blockWaitUpdateCh <- resultData:
		}
	}
}

// waitUpdateChFailedData
//
//	@Description: 程序异常重启后,先从数据库获取未更新的数据
//	@param ctx
//	@param chainId
//	@param blockWaitUpdateCh
//	@return error
func waitUpdateChFailedData(ctx context.Context, chainId string, blockWaitUpdateCh chan *BlockWaitUpdate) error {
	blockList, err := dbhandle.GetBlockByStatus(chainId, dbhandle.DelayUpdateFail)
	if err != nil {
		return err
	}

	//将未更新数据写入异步更新队列blockWaitUpdateCh
	for _, block := range blockList {
		// 将处理完成的结果写入 blockWaitUpdateCh
		resultData := &BlockWaitUpdate{
			ChainId:     chainId,
			BlockHeight: block.BlockHeight,
		}
		select {
		case <-ctx.Done():
			// 上下文已取消,不要发送数据到通道
			return nil
		case blockWaitUpdateCh <- resultData:
			//blockList的长度可能是blockWaitUpdateCh长度的2倍
		}
	}

	return nil
}

// RealtimeInsertOperation
//
//	@Description: BlockInsertWorkerCount个线程并发处理解析区块数据,存储到dataSaveCh通道,等待入库
//	@param ctx
//	@param hash hash值
//	@param blockInfoCh 订阅区块通道
//	@param dataSaveCh 保存区块数据通道
//	@param errCh
func RealtimeInsertOperation(ctx context.Context, hash string, blockInfoCh chan *common.BlockInfo,
	dataSaveCh chan *DataSaveToDB, errCh chan<- error) {
	workerCount := config.BlockInsertWorkerCount
	// 使用 sync.WaitGroup 来等待所有 worker 协程完成
	var wg sync.WaitGroup
	wg.Add(workerCount)
	// 启动 worker 协程,订阅blockInfoCh队列,并发解析区块数据,写入dataSaveCh通道
	for i := 0; i < workerCount; i++ {
		go ParallelParseBlockWork(ctx, &wg, hash, blockInfoCh, dataSaveCh, errCh)
	}
	wg.Wait()
}

// ParallelParseBlockWork
//
//	@Description: 消费blockInfoCh通道数据,解析成格式化的DB数据,存储到dataSaveCh,等待存储DB
//	@param ctx
//	@param wg
//	@param hashType
//	@param blockInfoCh 订阅区块通道
//	@param dataSaveCh 保存区块数据通道
//	@param errCh
func ParallelParseBlockWork(ctx context.Context, wg *sync.WaitGroup, hashType string,
	blockInfoCh chan *common.BlockInfo, dataSaveCh chan *DataSaveToDB, errCh chan<- error) {
	defer wg.Done()
	//blockInfoCh 阻塞持续等待blockInfoCh
	for blockInfo := range blockInfoCh {
		if blockInfo == nil {
			log.Errorf("blockInfoCh blockInfo failed.\n")
			continue
		}

		chainId := blockInfo.Block.Header.ChainId
		blockHeight := int64(blockInfo.Block.Header.BlockHeight)
		startTime := time.Now()
		// 处理区块数据
		log.Infof("【Realtime deal】start block-%s[%d]", chainId, blockHeight)
		dealResult, txTimeLog, err := RealtimeDataHandle(blockInfo, hashType)
		log.Infof("【Realtime deal】end block-%s[%d] duration_time(ms):%v",
			chainId, blockHeight, time.Since(startTime).Milliseconds())
		if err != nil {
			errCh <- fmt.Errorf("【Realtime deal】err block-%s[%d] failed, err:%v", chainId, blockHeight, err)
			// 如果处理失败,取消上下文,会重启链
			return
		}

		dataToSave := &DataSaveToDB{
			ChainId:     chainId,
			BlockHeight: blockHeight,
			DealResult:  *dealResult,
			TxTimeLog:   txTimeLog,
		}

		// 等待当前区块高度等于最大高度时,将数据写入 dataCh,否则需要等待一下
		done := false
		for !done {
			select {
			case <-ctx.Done():
				// 上下文已取消,跳出循环
				return
			default:
				if blockHeight == GetMaxHeight(chainId) {
					dataSaveCh <- dataToSave
					//只有写入dataSaveCh队列后才会加1
					setMaxHeight(chainId, blockHeight+1)
					done = true
				} else {
					//不sleep的话,高并发会占满cpu
					time.Sleep(20 * time.Millisecond)
				}
			}
		}
	}
}

// DelayUpdateOperation
//
//	@Description: 异步数据计算更新,订阅blockWaitUpdateCh通道数据
//	@param blockWaitUpdateCh 等到异步计算的通道
//	@param errCh
func DelayUpdateOperation(ctx context.Context, blockWaitUpdateCh chan *BlockWaitUpdate, errCh chan<- error) {
	for {
		var (
			blockWaitUpdates []*BlockWaitUpdate
			maxCount         = config.BlockUpdateWorkerCount
			blockHeightList  []int64
			chainId          string
			heightStr        string
		)

		// 从 blockWaitUpdateCh 中读取一个区块数据
		select {
		case blockInfo, ok := <-blockWaitUpdateCh:
			if !ok {
				// blockWaitUpdateCh 通道已关闭,退出函数
				log.Errorf("Child process terminated due to channel closed.\n")
				return
			}
			blockWaitUpdates = append(blockWaitUpdates, blockInfo)
		case <-ctx.Done():
			log.Infof("Context cancelled, exiting function.\n")
			return
		}

		// 继续读取,直到达到最大数量（例如,5个）或通道为空
		for i := 0; i < maxCount; i++ {
			select {
			case updateInfo := <-blockWaitUpdateCh:
				blockWaitUpdates = append(blockWaitUpdates, updateInfo)
			case <-ctx.Done():
				log.Infof("Context cancelled, exiting function.\n")
				return
			default:
				// 通道为空,停止读取
				continue
			}
		}

		for _, blockData := range blockWaitUpdates {
			chainId = blockData.ChainId
			// 处理已完成的 BlockHeight
			blockHeightList = append(blockHeightList, blockData.BlockHeight)
			heightStr += strconv.FormatInt(blockData.BlockHeight, 10) + ","
		}

		if chainId == "" {
			break
		}

		//开始异步计算
		log.Infof("【Delay update】start block-%s[%s]", chainId, heightStr)
		//更新数据库
		startTime := time.Now()
		err := BatchDelayedUpdate(chainId, blockHeightList)
		//异步计算结束
		durationTime := time.Since(startTime).Milliseconds()
		log.Infof("【Delay update】end block-%s[%s], duration_time:%vms", chainId, heightStr, durationTime)
		if err != nil {
			errCh <- fmt.Errorf("【Delay update】BatchDelayedUpdate failed, err: %v", err)
			return
		}
	}
}

// SubscribeBlockSetToBlockInfoCh
//
//	@Description: 将订阅的区块数据Json解析成common.BlockInfo结构,写入blockInfoCh通道
//	@param ctx
//	@param sdkClient 链连接
//	@param blockInfoCh 区块处理通道
//	@param errCh
func SubscribeBlockSetToBlockInfoCh(ctx context.Context, sdkClient *SdkClient, blockInfoCh chan *common.BlockInfo,
	errCh chan<- error) {
	chainId := sdkClient.ChainId
	chainClient := sdkClient.ChainClient

	// get max block height for this chain
	maxBlockHeight := dbhandle.GetMaxBlockHeight(chainId)
	log.Infof("【sync load】 begin to subscribe block-%s[%d] ", chainId, maxBlockHeight)
	if maxBlockHeight > 0 {
		setMaxHeight(sdkClient.ChainId, maxBlockHeight+1)
	} else {
		setMaxHeight(sdkClient.ChainId, 0)
	}

	//订阅区块
	c, err := chainClient.SubscribeBlock(ctx, maxBlockHeight, -1, true, false)
	if err != nil {
		subscribeFail.WithLabelValues(chainId).Inc()
		errCh <- fmt.Errorf("【Sync Block】 Get Block By SDK failed:, err: %v", err)
		return
	}

	// 创建一个定时器用于刷新锁的过期时间
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		//链订阅已经停止,停止定时器
		if sdkClient.Status == STOP {
			return
		}

		select {
		case block, ok := <-c:
			if !ok {
				subscribeFail.WithLabelValues(chainId).Inc()
				errCh <- fmt.Errorf("【Sync Block】 SubscribeBlock- Chan Is Closed, chainId:%v, ok:%v",
					chainId, ok)
				return
			}

			blockInfo, ok := block.(*common.BlockInfo)
			if !ok {
				subscribeFail.WithLabelValues(chainId).Inc()
				errCh <- fmt.Errorf("【Sync Block】 SubscribeBlock- The Data Type Error, chainId:%v", chainId)
				return
			}

			//根据区块高度获取区块信息
			height := int64(blockInfo.Block.Header.BlockHeight)
			blockDB, _ := dbhandle.GetBlockByHeight(chainId, height)
			if blockDB != nil && blockDB.BlockHash != "" {
				//数据库已经存在
				setMaxHeight(chainId, height+1)
				log.Infof("【Sync Block】block is existed, chainId:%v, block height:%v \n", chainId, height)
			} else {
				select {
				case <-ctx.Done():
					// 上下文已取消,不要发送数据到通道
					return
				case blockInfoCh <- blockInfo:
					// 成功发送数据到通道
				}
			}
		case <-ticker.C:
			// 定期刷新分布式锁的过期时间
			prefix := config.GlobalConfig.RedisDB.Prefix
			lockKey := fmt.Sprintf(cache.RedisSubscribeLockKey, prefix, chainId)
			cache.GlobalRedisDb.Expire(ctx, lockKey, 3*time.Minute)
			log.Infof("【Sync Block】Redis Set Lock, LockKey:%s", lockKey)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// BuildChainInfo
//
//	@Description: 构造链信息
//	@param subscribeChain 链订阅信息
//	@return *config.ChainInfo 链配置结构
func BuildChainInfo(subscribeChain *db.Subscribe) *config.ChainInfo {
	chainInfo := &config.ChainInfo{
		ChainId:  subscribeChain.ChainId,
		AuthType: subscribeChain.AuthType,
		OrgId:    subscribeChain.OrgId,
		HashType: subscribeChain.HashType,
		UserInfo: &config.UserInfo{
			UserKey:  subscribeChain.UserKey,
			UserCert: subscribeChain.UserCert,
		},
	}
	var nodeList []*config.NodeInfo
	_ = json.Unmarshal([]byte(subscribeChain.NodeList), &nodeList)
	chainInfo.NodesList = nodeList
	return chainInfo
}

//type BlockDataStore struct {
//	data      map[int64]*DataSaveToDB
//	nextBlock int64
//	lock      sync.Mutex
//	cond      *sync.Cond
//}
//
//func NewBlockDataStore() *BlockDataStore {
//	bds := &BlockDataStore{
//		data:      make(map[int64]*DataSaveToDB),
//		nextBlock: 0,
//	}
//	bds.cond = sync.NewCond(&bds.lock)
//	return bds
//}

//func (bds *BlockDataStore) AddData(data *DataSaveToDB) {
//	bds.lock.Lock()
//	defer bds.lock.Unlock()
//
//	bds.data[data.BlockHeight] = data
//	bds.cond.Broadcast()
//}
//
//func (bds *BlockDataStore) GetNextData() *DataSaveToDB {
//	bds.lock.Lock()
//	defer bds.lock.Unlock()
//
//	for bds.data[bds.nextBlock] == nil {
//		bds.cond.Wait()
//	}
//	data := bds.data[bds.nextBlock]
//	delete(bds.data, bds.nextBlock)
//	bds.nextBlock++
//	return data
//}
//
//func RealtimeInsertOperation(ctx context.Context, hash string, blockInfoCh chan *common.BlockInfo,
//	errCh chan<- error) {
//	workerCount := config.BlockInsertWorkerCount
//
//	// 创建一个 BlockDataStore 实例
//	blockDataStore := NewBlockDataStore()
//
//	// 使用 sync.WaitGroup 来等待所有 worker 协程完成
//	var wg sync.WaitGroup
//	wg.Add(workerCount + 1) // 注意这里需要增加一个协程
//
//	// 启动一个新的协程,按照 blockHeight 的顺序插入数据库
//	go func() {
//		defer wg.Done()
//
//		for {
//			select {
//			case <-ctx.Done():
//				return
//			default:
//				data := blockDataStore.GetNextData()
//				// 在这里将 data 插入数据库
//				fmt.Printf("Inserting block %d\n", data.BlockHeight)
//			}
//		}
//	}()
//
//	// 启动 worker 协程,订阅blockInfoCh队列,并发解析区块数据,写入 BlockDataStore 实例
//	for i := 0; i < workerCount; i++ {
//		go ParallelParseBlockWork(ctx, &wg, hash, blockInfoCh, blockDataStore, errCh)
//	}
//	wg.Wait()
//}
