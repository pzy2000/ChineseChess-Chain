/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/chain"
	"chainmaker_web/src/config"
	"context"
	"errors"
	"sync"
	"time"

	commonlog "chainmaker.org/chainmaker/common/v2/log"
	"go.uber.org/zap"

	"chainmaker.org/chainmaker/common/v2/crypto"
	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

const (
	// STOP 1
	STOP = 1
	// ConnCount count
	ConnCount = 5
)

// ErrorRestart chain restart
var ErrorRestart = errors.New("chainId need restart")

// SdkClient sdk
type SdkClient struct {
	// 定时任务每次只启动一个
	lock        sync.Mutex
	Ctx         context.Context
	Cancel      func()
	ChainId     string
	ChainInfo   *config.ChainInfo
	ChainClient *sdk.ChainClient
	ChainConfig *pbconfig.ChainConfig
	Status      int
}

// CreateChainClient CreateSdkClientWithChainId return Sdk chain client with chain-id
// @desc
// @param ${param}
// @return *sdk.ChainClient
// @return error
func CreateChainClient(chainInfo *config.ChainInfo) (*sdk.ChainClient, error) {
	nodeList := chainInfo.NodesList
	nodeOptions := make([]sdk.ChainClientOption, 0)
	nodeOptions = append(nodeOptions, sdk.WithChainClientChainId(chainInfo.ChainId))
	nodeOptions = append(nodeOptions, sdk.WithAuthType(chainInfo.AuthType))

	// 公钥模式
	if chainInfo.AuthType == config.PUBLIC {
		cryptoConfig := sdk.NewCryptoConfig(sdk.WithHashAlgo(chainInfo.HashType))
		nodeOptions = append(nodeOptions, sdk.WithCryptoConfig(cryptoConfig))
		nodeOptions = append(nodeOptions, sdk.WithUserSignKeyBytes([]byte(chainInfo.UserInfo.UserKey)))
	} else {
		nodeOptions = append(nodeOptions, sdk.WithChainClientOrgId(chainInfo.OrgId))
		nodeOptions = append(nodeOptions, sdk.WithUserKeyBytes([]byte(chainInfo.UserInfo.UserKey)))
		nodeOptions = append(nodeOptions, sdk.WithUserCrtBytes([]byte(chainInfo.UserInfo.UserCert)))
		nodeOptions = append(nodeOptions, sdk.WithUserSignKeyBytes([]byte(chainInfo.UserInfo.UserKey)))
		nodeOptions = append(nodeOptions, sdk.WithUserSignCrtBytes([]byte(chainInfo.UserInfo.UserCert)))
	}

	rpcClient := sdk.NewRPCClientConfig(
		sdk.WithRPCClientMaxReceiveMessageSize(1024),
	)
	nodeOptions = append(nodeOptions, sdk.WithRPCClientConfig(rpcClient))

	conf := &sdk.Pkcs11Config{Enabled: false}
	nodeOptions = append(nodeOptions, sdk.WithPkcs11Config(conf))

	nodeOptions = append(nodeOptions, sdk.WithRPCClientConfig(sdk.NewRPCClientConfig(
		sdk.WithRPCClientMaxReceiveMessageSize(1024))))

	for _, nodeInfo := range nodeList {
		node := sdk.NewNodeConfig(
			// 节点地址，格式：127.0.0.1:12301
			sdk.WithNodeAddr(nodeInfo.Addr),
			// 节点连接数
			sdk.WithNodeConnCnt(ConnCount),
			// 节点是否启用TLS认证
			sdk.WithNodeUseTLS(nodeInfo.Tls),
			// 根证书路径，支持多个
			sdk.WithNodeCACerts([]string{nodeInfo.OrgCA}),
			// TLS Hostname
			sdk.WithNodeTLSHostName(nodeInfo.TLSHostName),
		)
		nodeOptions = append(nodeOptions, sdk.AddChainClientNodeConfig(node))
	}
	nodeOptions = append(nodeOptions, sdk.WithChainClientLogger(getDefaultLogger()))
	chainClient, err := sdk.NewChainClient(nodeOptions...)
	if err != nil {
		return nil, err
	}

	return chainClient, nil
}

// NewSdkClient NewSdkClient
func NewSdkClient(chainInfo *config.ChainInfo, client *sdk.ChainClient) *SdkClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &SdkClient{
		ChainId:     chainInfo.ChainId,
		ChainClient: client,
		ChainInfo:   chainInfo,
		Ctx:         ctx,
		Cancel:      cancel,
	}
}

// GetChainClient get client
func (sdkClient *SdkClient) GetChainClient() *sdk.ChainClient {
	return sdkClient.ChainClient
}

// SdkClientPool pool
type SdkClientPool struct {
	sdkClients sync.Map
}

// SingleSdkClientPool pool
type SingleSdkClientPool struct {
	chainInfo       *config.ChainInfo
	systemSdkClient *SdkClient
	queryClient     *sdk.ChainClient
	sdkClients      sync.Map
}

// NewSingleSdkClientPool new一个pool
func NewSingleSdkClientPool(chainInfo *config.ChainInfo, systemSdkClient *SdkClient,
	queryClient *sdk.ChainClient) *SingleSdkClientPool {
	pool := &SingleSdkClientPool{
		chainInfo:       chainInfo,
		systemSdkClient: systemSdkClient,
		queryClient:     queryClient,
		sdkClients:      sync.Map{},
	}
	pool.addSdkClient(systemSdkClient) // 添加sdkClient到pool中
	return pool
}

// NewSdkClientPool new一个pool
func NewSdkClientPool() *SdkClientPool {
	return &SdkClientPool{
		sdkClients: sync.Map{},
	}
}

// GetChainClient 获取指定客户端
func GetChainClient(chainId string) *sdk.ChainClient {
	val1, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return nil
	}
	singleSdkClient, ok := val1.(*SingleSdkClientPool)
	if !ok {
		return nil
	}
	if val2, ok := singleSdkClient.sdkClients.Load(chainId); ok {
		return val2.(*SdkClient).ChainClient
	}
	return nil
}

// GetSdkClient 获取指定客户端
func GetSdkClient(chainId string) *SdkClient {
	val1, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return nil
	}
	singleSdkClient, ok := val1.(*SingleSdkClientPool)
	if !ok {
		return nil
	}
	if val2, ok := singleSdkClient.sdkClients.Load(chainId); ok {
		return val2.(*SdkClient)
	}
	return nil
}

// GetAllSdkClient 获取指定客户端
func GetAllSdkClient(chainList []*config.ChainInfo) []*SdkClient {
	clients := make([]*SdkClient, 0)
	for _, chainInfo := range chainList {
		sdkClients := GetSdkClient(chainInfo.ChainId)
		if sdkClients == nil {
			continue
		}

		clients = append(clients, sdkClients)
	}
	return clients
}

// addSdkClient add SDKClient
func (pool *SingleSdkClientPool) addSdkClient(sdkClient *SdkClient) {
	pool.sdkClients.Store(sdkClient.ChainId, sdkClient)
}

// RemoveSdkClient addSdkClient add SDKClient
// @desc
// @param ${param}
// @return error
func (pool *SdkClientPool) RemoveSdkClient(chainId string) {
	pool.sdkClients.Delete(chainId)
}

func (pool *SdkClientPool) addSdkClientPool(singleSdkClientPool *SingleSdkClientPool) {
	pool.sdkClients.Delete(singleSdkClientPool.chainInfo.ChainId)
	pool.sdkClients.Store(singleSdkClientPool.chainInfo.ChainId, singleSdkClientPool)
}

// UpdateChainConfig updateSdkClient add SDKClient
// @desc
// @param ${param}
// @return error
func (pool *SdkClientPool) UpdateChainConfig(chainId string, newChainConfig *pbconfig.ChainConfig) {
	sdkClient := GetSdkClient(chainId)
	if sdkClient != nil {
		sdkClient.ChainConfig = newChainConfig
	}
}

// LoadChains 加载链
func (pool *SdkClientPool) LoadChains(chainId string) {
	value, ok := pool.sdkClients.Load(chainId)
	if ok {
		value.(*SingleSdkClientPool).sdkClients.Range(func(key, value interface{}) bool {
			value.(*SdkClient).Load()
			return true
		})
	} else {
		log.Errorf("【sync】LoadChains failed, chainId[%v]", chainId)
	}
}

// Load resource init
func (sdkClient *SdkClient) Load() {
	log.Infof("[WEB] begin to load chain's information, [chain:%s] ", sdkClient.ChainId)
	//处理节点，组织数据
	err := loadChainRefInfos(sdkClient)
	if err != nil {
		//重启链
		ReStartChain(sdkClient)
	}

	//订阅区块数据
	go PeriodicGetSubscribeLock(sdkClient)

	//定时处理
	//定期处理节点数据
	go PeriodicLoadStart(sdkClient)
	//定期检查子链健康状态
	go PeriodicCheckSubChainStatus(sdkClient)
}

// GetChainHashType 获取hash
func (sdkClient *SdkClient) GetChainHashType() string {
	hash := sdkClient.ChainConfig.Crypto.Hash
	if hash == "" {
		log.Error("[SDK] Get Chain Config Failed : ")
		return crypto.CRYPTO_ALGO_SM3
	}

	return hash
}

// ReStartChain 重启订阅
func ReStartChain(sdkClient *SdkClient) {
	//停掉这个连接
	sdkClient.Status = STOP
	sdkClient.Cancel()

	//判断连接池还在不在，不在的话不在重启
	poolSdkClient := GetSdkClient(sdkClient.ChainId)
	if poolSdkClient == nil {
		log.Infof("【ReStartChain】poolSdkClient is null, chain is cancel，chainId:%v", sdkClient.ChainId)
		return
	}

	//重启这条链的订阅
	chainConfig, err := chain.GetSubscribeByChainId(sdkClient.ChainId)
	if err == nil && chainConfig != nil {
		time.Sleep(time.Second * 10)
		StartSync([]*config.ChainInfo{chainConfig})
		//ReStart(chainConfig)
	}
}

// StopChain 停止订阅
func StopChain(sdkClient *SdkClient) {
	//停掉这个连接
	sdkClient.Status = STOP
	sdkClient.Cancel()

	// 移除订阅连接
	RemoveSubscribeChain(sdkClient.ChainId)
}

// RemoveSubscribeChain remove
// @desc 删除订阅
// @param ${param}
// @return error
func RemoveSubscribeChain(chainId string) {
	sdkClientPool = NewSdkClientPool()
	sdkClientPool.RemoveSdkClient(chainId)
}

func getDefaultLogger() *zap.SugaredLogger {
	configInfo := commonlog.LogConfig{
		Module:       "[SDK]",
		LogPath:      "../log/sdk.log",
		LogLevel:     commonlog.LEVEL_INFO,
		MaxAge:       30,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: false,
	}

	logger, _ := commonlog.InitSugarLogger(&configInfo)
	return logger
}
