// Package chain provides chain Methods
package chain

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db/dbhandle"
	loggers "chainmaker_web/src/logger"
	"encoding/json"
)

var (
	log = loggers.GetLogger(loggers.MODULE_SYNC)
)

// InitChainConfig
//
//	@Description: 初始化链订阅数据
func InitChainConfig() {
	//从数据库获取订阅数据
	subscribeChains, err := GetSubscribeChains()
	if len(subscribeChains) > 0 && err == nil {
		//将数据库订阅和配置订阅数据合并
		mergedChains := mergeChainInfo(config.SubscribeChains, subscribeChains)
		config.SubscribeChains = mergedChains
	}
}

// mergeChainInfo
//
//	@Description: 将DB订阅数据和配置额订阅数据合并， 合并配置数据和DB数据,相同数据用DB数据
//	@param configChains 配置的订阅数据
//	@param dbChains 数据库的订阅数据
//	@return []*config.ChainInfo 合并后的订阅数据
func mergeChainInfo(configChains, dbChains []*config.ChainInfo) []*config.ChainInfo {
	mergedChains := make([]*config.ChainInfo, 0)

	// 创建一个映射，用于存储 chains2 中的 ChainInfo，以便于根据 ChainId 进行查找
	dbChainsMap := make(map[string]*config.ChainInfo)
	for _, chain := range dbChains {
		dbChainsMap[chain.ChainId] = chain
	}

	// 遍历 configChains，如果存在相同的 ChainId，则使用 dbChains 中的数据
	for _, chain := range configChains {
		if chain2, ok := dbChainsMap[chain.ChainId]; ok {
			mergedChains = append(mergedChains, chain2)
			// 从 chains2Map 中删除已合并的 ChainInfo，以便于后续处理 chains2 中剩余的数据
			delete(dbChainsMap, chain.ChainId)
		} else {
			mergedChains = append(mergedChains, chain)
		}
	}

	// 将 dbChains 中剩余的 ChainInfo 添加到 mergedChains 中
	for _, chain := range dbChainsMap {
		mergedChains = append(mergedChains, chain)
	}

	return mergedChains
}

// GetSubscribeChains 获取订阅信息，没有就用配置文件
func GetSubscribeChains() ([]*config.ChainInfo, error) {
	var err error
	chainConfigs := make([]*config.ChainInfo, 0)
	//数据库获取订阅数据
	subscribeChains, err := dbhandle.GetDBSubscribeChains()
	if err != nil {
		log.Errorf("Init GetSubscribeChains err:%v", err)
		return nil, err
	}
	if len(subscribeChains) == 0 {
		return chainConfigs, nil
	}
	//数据库链数据
	for _, chainInfo := range subscribeChains {
		var nodeList []*config.NodeInfo
		if chainInfo.NodeList != "" {
			err = json.Unmarshal([]byte(chainInfo.NodeList), &nodeList)
			if err != nil {
				log.Errorf("chain node list json Unmarshal failed, err:%v", err)
				continue
			}
		}

		chain := &config.ChainInfo{
			ChainId:   chainInfo.ChainId,
			AuthType:  chainInfo.AuthType,
			OrgId:     chainInfo.OrgId,
			HashType:  chainInfo.HashType,
			NodesList: nodeList,
			UserInfo: &config.UserInfo{
				UserKey:  chainInfo.UserKey,
				UserCert: chainInfo.UserCert,
			},
		}
		chainConfigs = append(chainConfigs, chain)
	}
	return chainConfigs, nil
}

// GetSubscribeByChainId 获取订阅信息，没有就用配置文件
func GetSubscribeByChainId(chainId string) (*config.ChainInfo, error) {
	var err error
	//数据库获取订阅数据
	chainInfoDB, err := dbhandle.GetSubscribeByChainId(chainId)
	if err != nil || chainInfoDB == nil {
		log.Errorf("Init GetSubscribeChains err:%v", err)
		return nil, err
	}

	var nodeList []*config.NodeInfo
	if chainInfoDB.NodeList != "" {
		err = json.Unmarshal([]byte(chainInfoDB.NodeList), &nodeList)
		if err != nil {
			log.Errorf("chain node list json Unmarshal failed, err:%v", err)
			return nil, err
		}
	}

	chainConfig := &config.ChainInfo{
		ChainId:   chainInfoDB.ChainId,
		AuthType:  chainInfoDB.AuthType,
		OrgId:     chainInfoDB.OrgId,
		HashType:  chainInfoDB.HashType,
		NodesList: nodeList,
		UserInfo: &config.UserInfo{
			UserKey:  chainInfoDB.UserKey,
			UserCert: chainInfoDB.UserCert,
		},
	}
	return chainConfig, nil
}

// GetConfigShow - 获取链配置是否显示
func GetConfigShow() bool {
	if config.GlobalConfig == nil || config.GlobalConfig.ChainConf == nil {
		return false
	}

	return config.GlobalConfig.ChainConf.ShowConfig
}

// GetIsMainChain - 是否是主链
func GetIsMainChain() bool {
	if config.GlobalConfig == nil || config.GlobalConfig.ChainConf == nil {
		return false
	}

	return config.GlobalConfig.ChainConf.IsMainChain
}
