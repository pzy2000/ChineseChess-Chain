package service

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"chainmaker_web/src/sync"
	"encoding/json"
	"fmt"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

const RestartSubscribe = 0
const StopSubscribe = 1

// GetChainListHandler get
type GetChainListHandler struct {
}

// Handle GetChainListHandler 区块列表
func (getChainListHandler *GetChainListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetChainListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//根据分页获取区块
	chainList, count, err := dbhandle.GetChainListByPage(params.Offset, params.Limit, params.ChainId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	var chainIds []string
	for _, v := range chainList {
		chainIds = append(chainIds, v.ChainId)
	}

	//获取订阅信息
	subscribeList, err := dbhandle.GetSubscribeByChainIds(chainIds)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// 创建一个映射，将 chainId 映射到对应的 Subscribe 状态
	subscribeMap := make(map[string]*db.Subscribe)
	for _, subscribe := range subscribeList {
		subscribeMap[subscribe.ChainId] = subscribe
	}

	// view
	chainListView := arraylist.New()
	for _, chain := range chainList {
		status := 1
		if _, ok := subscribeMap[chain.ChainId]; ok {
			status = subscribeMap[chain.ChainId].Status
		}
		chainView := &entity.ChainListView{
			ChainId:      chain.ChainId,
			ChainVersion: chain.Version,
			Consensus:    chain.Consensus,
			Status:       status,
			Timestamp:    chain.CreatedAt.Unix(),
			AuthType:     chain.AuthType,
		}
		chainListView.Add(chainView)
	}

	ConvergeListResponse(ctx, chainListView.Values(), count, nil)
}

// SubscribeChainHandler sub
type SubscribeChainHandler struct {
}

// Handle SubscribeChainHandler 新增订阅
func (handler *SubscribeChainHandler) Handle(ctx *gin.Context) {
	params := entity.BindSubscribeChainHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取区块链订阅信息
	subscribeInfo, err := dbhandle.GetSubscribeByChainId(params.ChainId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}
	if subscribeInfo != nil {
		newError := entity.NewError(entity.ErrorSubscribe, "chain id already exists")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//非国密
	hashType := crypto.CRYPTO_ALGO_SHA256
	if params.HashType == config.SM2 {
		//国密
		hashType = crypto.CRYPTO_ALGO_SM3
	}

	nodes, err := json.Marshal(params.NodeList)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	sdkConfig := &db.Subscribe{
		ChainId:  params.ChainId,
		OrgId:    params.OrgId,
		UserKey:  params.UserKey,
		UserCert: params.UserCert,
		NodeList: string(nodes),
		Status:   db.SubscribeOK,
		AuthType: params.AuthType,
		HashType: hashType,
	}

	chainInfo := sync.BuildChainInfo(sdkConfig)
	chainList := make([]*config.ChainInfo, 0)
	chainList = append(chainList, chainInfo)
	//数据表初始化
	dbCfg := config.GlobalConfig.DBConf
	db.InitDBTable(dbCfg, chainList)

	//订阅链
	err = sync.SubscribeChain(chainInfo)
	if err != nil {
		errFmt := fmt.Errorf("订阅链失败，请检查配置数据")
		ConvergeHandleFailureResponse(ctx, errFmt)
		return
	}

	//加入全局变量，所有订阅的链数据
	// 遍历SubscribeChains，查找相同的ChainId
	found := false
	for _, existingChainInfo := range config.SubscribeChains {
		if existingChainInfo.ChainId == chainInfo.ChainId {
			// 如果找到相同的ChainId，更新数据
			*existingChainInfo = *chainInfo
			found = true
			break
		}
	}
	// 如果没有找到相同的ChainId，将chainInfo添加到SubscribeChains
	if !found {
		config.SubscribeChains = append(config.SubscribeChains, chainInfo)
	}

	ConvergeDataResponse(ctx, sdkConfig, nil)
}

// CancelSubscribeHandler cancel
type CancelSubscribeHandler struct {
}

// Handle deal，暂停，开始订阅
func (handler *CancelSubscribeHandler) Handle(ctx *gin.Context) {
	//var err error
	params := entity.BindCancelSubscribeHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//查看订阅连接池是否存在
	sdkClient := sync.GetSdkClient(params.ChainId)
	if params.Status == RestartSubscribe {
		if sdkClient != nil {
			//订阅已经存在，不需要重新订阅
			log.Infof("Subscribe chain already exists")
			ConvergeDataResponse(ctx, "OK", nil)
			return
		}

		//查询数据库订阅信息
		subscribeInfo, err := dbhandle.GetSubscribeByChainId(params.ChainId)
		if err != nil {
			ConvergeHandleFailureResponse(ctx, err)
			return
		}

		//重新订阅
		chainInfo := sync.BuildChainInfo(subscribeInfo)
		//go sync.Start(chainList)
		err = sync.SubscribeChain(chainInfo)
		if err != nil {
			ConvergeHandleFailureResponse(ctx, err)
			return
		}
	} else if params.Status == StopSubscribe {
		if sdkClient == nil {
			//订阅已经不存在，不需要暂停订阅
			log.Infof("Subscribe chain already exists")
			ConvergeDataResponse(ctx, "OK", nil)
			return
		}
		//暂停订阅
		sync.StopChain(sdkClient)
		//设置订阅状态
		err := dbhandle.SetSubscribeStatus(params.ChainId, db.SubscribeCanceled)
		if err != nil {
			log.Errorf("remove %s subscribe failed", params.ChainId)
			ConvergeHandleFailureResponse(ctx, err)
			return
		}
	}
	ConvergeDataResponse(ctx, "OK", nil)
}

// ModifySubscribeHandler modify
type ModifySubscribeHandler struct {
}

// Handle ModifySubscribeHandler 修改订阅
func (handler *ModifySubscribeHandler) Handle(ctx *gin.Context) {
	params := entity.BindModifySubscribeHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取区块链订阅信息
	sdkConfig, err := dbhandle.GetSubscribeByChainId(params.ChainId)
	if err != nil {
		log.Debugf("get %s subscribe failed", params.ChainId)
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	hashType := crypto.CRYPTO_ALGO_SHA256
	if params.HashType == config.SM2 {
		hashType = crypto.CRYPTO_ALGO_SM3
	}

	// 更新数据
	sdkConfig.OrgId = params.OrgId
	sdkConfig.UserCert = params.UserCert
	sdkConfig.UserKey = params.UserKey
	sdkConfig.AuthType = params.AuthType
	sdkConfig.HashType = hashType
	nodes, err := json.Marshal(params.NodeList)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}
	sdkConfig.NodeList = string(nodes)

	chainInfo := sync.BuildChainInfo(sdkConfig)
	//订阅链
	// 新的连接会覆盖掉连接池内旧的连接。
	//go sync.Start(chainList)
	//订阅链
	err = sync.SubscribeChain(chainInfo)
	if err != nil {
		errFmt := fmt.Errorf("订阅链失败，请检查配置数据")
		ConvergeHandleFailureResponse(ctx, errFmt)
		return
	}

	//修改全局变量，所有订阅的链数据
	for i, chain := range config.SubscribeChains {
		if chain.ChainId == chainInfo.ChainId {
			config.SubscribeChains[i] = chainInfo
			break
		}
	}

	ConvergeDataResponse(ctx, "OK", nil)
}

// DeleteSubscribeHandler delete
type DeleteSubscribeHandler struct {
}

// Handle DeleteSubscribeHandler删除订阅数据
func (handler *DeleteSubscribeHandler) Handle(ctx *gin.Context) {
	params := entity.BindDeleteSubscribeHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// 如果已经存在对应的取消函数，先调用取消
	if cancel, ok := sync.StartSyncCancels.Load(params.ChainId); ok {
		cancel.(context.CancelFunc)()
	}

	//获取连接池信息
	sdkClient := sync.GetSdkClient(params.ChainId)
	if sdkClient != nil {
		//订阅信息还存在，暂停订阅
		sync.StopChain(sdkClient)
	}

	//删除ChainId所有相关表
	err := db.DeleteTablesByChainID(params.ChainId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//删除链信息
	err = dbhandle.DeleteChain(params.ChainId)
	if err != nil {
		log.Errorf("remove chain info failed, chainId:%v", params.ChainId)
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//删除订阅信息
	err = dbhandle.DeleteSubscribe(params.ChainId)
	if err != nil {
		log.Errorf("remove %s subscribe failed", params.ChainId)
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//删除缓存数据
	go DeleteRedisCache(params.ChainId)

	ConvergeDataResponse(ctx, "OK", nil)
}

// DeleteRedisCache 删除缓存数据
func DeleteRedisCache(chainId string) {
	ctx := context.Background()
	// 构造匹配前缀
	prefix := config.GlobalConfig.RedisDB.Prefix
	pattern := fmt.Sprintf("%s_%s*", prefix, chainId)
	log.Infof("DeleteRedisCache redis keys with pattern: %s", pattern)

	// 使用 SCAN 命令找到所有匹配的键
	var cursor uint64
	var lengthTotal int
	for {
		var keys []string
		var err error
		// 这里的100表示每次扫描操作返回的键的最大数量
		keys, cursor, err = cache.GlobalRedisDb.Scan(ctx, cursor, pattern, 50).Result()
		lengthTotal += len(keys)
		if err != nil {
			log.Errorf("DeleteRedisCache Failed to scan keys: %v", err)
		}

		if len(keys) > 0 {
			// 删除匹配的键
			err := cache.GlobalRedisDb.Del(ctx, keys...).Err()
			if err != nil {
				log.Errorf("DeleteRedisCache Failed to delete keys: %v", err)
			}
		}

		// 如果 cursor 为 0，则表示遍历完毕
		if cursor == 0 {
			break
		}
	}

	log.Infof("DeleteRedisCache redis keys length:%v", lengthTotal)
}
