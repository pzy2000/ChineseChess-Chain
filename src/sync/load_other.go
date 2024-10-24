/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/utils"
	"encoding/json"
	"encoding/pem"
	"errors"
	"strings"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
	commonutils "chainmaker.org/chainmaker/utils/v2"

	"chainmaker_web/src/config"
	"chainmaker_web/src/db/dbhandle"

	"chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/pb-go/v2/discovery"
)

// PeriodicCheckSubChainStatus
//
//	@Description: 检查子链健康状态， 1小时检查一次
//	@param sdkClient 链连接
func PeriodicCheckSubChainStatus(sdkClient *SdkClient) {
	chainId := sdkClient.ChainId

	//1小时定时器
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		//链订阅已经停止,停止定时器
		if sdkClient.Status == STOP {
			return
		}

		//查询子链列表
		crossSubChains, err := dbhandle.GetCrossSubChainAll(chainId)
		if err != nil {
			log.Errorf("[load_other] get cross sub chain failed:%v", err)
		}

		for _, subChain := range crossSubChains {
			//Grpc-获取子链健康状态
			chainOk, errGrpc := CheckSubChainStatus(subChain)
			if errGrpc != nil {
				subChainJson, _ := json.Marshal(subChain)
				log.Errorf("[load_other] CheckSubChainStatus failed, err:%v, subChainJson:%v",
					errGrpc, string(subChainJson))
			}

			status := dbhandle.SubChainStatusSuccess
			if !chainOk {
				status = dbhandle.SubChainStatusFail
			}

			//为获取同步区块高度合约
			var spvContractName string
			if subChain.SpvContractName == "" {
				subChainInfo, errRpc := utils.GetCrossSubChainInfo(subChain.SubChainId)
				if errRpc != nil || subChainInfo == nil {
					log.Errorf("【load_other】http get sub chain failed, err:%v, SubChainId:%v",
						err, subChain.SubChainId)
				} else {
					spvContractName = subChainInfo.SpvContractName
				}
			}

			//健康状态变更，更新数据库
			if status != subChain.Status || spvContractName != "" {
				err = dbhandle.UpdateCrossSubChainStatus(chainId, subChain.SubChainId, spvContractName, status)
				if err != nil {
					log.Errorf("[load_other] update cross sub chain status failed, err:%v", err)
				}
			}
		}

	}
}

// PeriodicLoadStart
//
//	@Description:  1小时请求一次，检查是否有新增节点
//	@param sdkClient
func PeriodicLoadStart(sdkClient *SdkClient) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if sdkClient.Status == STOP {
			return
		}
		loadNodeInfo(sdkClient)
	}
}

// loadChainRefInfos
//
//	@Description: load chain and other information
//	@param sdkClient
//	@return error
func loadChainRefInfos(sdkClient *SdkClient) error {
	var err error
	//处理节点数据
	loadNodeInfo(sdkClient)
	if sdkClient.ChainInfo.AuthType == config.PUBLIC {
		//处理user数据
		err = loadChainUser(sdkClient)
	} else {
		//处理组织数据
		err = loadOrgInfo(sdkClient)
	}

	if err != nil {
		return err
	}

	return nil
}

// loadChainUser loadChainUser
//
// loadChainUser
//
//	@Description: 更新user数据
//	@param sdkClient
//	@return error
func loadChainUser(sdkClient *SdkClient) error {
	//链配置信息
	chainConfig := sdkClient.ChainConfig
	chainId := chainConfig.ChainId
	userList := GetAdminUserByConfig(chainConfig)
	return dbhandle.BatchInsertUser(chainId, userList)
}

// GetAdminUserByConfig
//
//	@Description: 根据链配置获取user信息
//	@param chainConfig 链配置
//	@return []*db.User user列表
func GetAdminUserByConfig(chainConfig *pbconfig.ChainConfig) []*db.User {
	userList := make([]*db.User, 0)
	if chainConfig.ChainId == "" || len(chainConfig.TrustRoots) <= 0 {
		return userList
	}

	hashType := chainConfig.Crypto.Hash
	for _, root := range chainConfig.TrustRoots[0].Root {
		publicKey, err := asym.PublicKeyFromPEM([]byte(root))
		if err != nil {
			log.Error("[SDK] get publicKey by PK err : " + err.Error())
			continue
		}
		addr, err := commonutils.PkToAddrStr(publicKey, pbconfig.AddrType_ETHEREUM, crypto.HashAlgoMap[hashType])
		if err != nil {
			log.Error("[SDK] get addr by PK err : " + err.Error())
			continue
		}
		//userId, err := helper.CreateLibp2pPeerIdWithPublicKey(publicKey)
		//if err != nil {
		//	continue
		//}

		user := &db.User{
			UserId:    addr,
			UserAddr:  addr,
			Role:      "admin",
			OrgId:     config.PUBLIC,
			Timestamp: time.Now().Unix(),
		}
		userList = append(userList, user)
	}
	return userList
}

// loadOrgInfo
//
//	@Description: 加载组织数据
//	@param sdkClient
//	@return error
func loadOrgInfo(sdkClient *SdkClient) error {
	chainConfig := sdkClient.ChainConfig
	return dbhandle.SaveOrgByConfig(chainConfig)
}

// loadNodeInfo
//
//	@Description: 加载节点数据
//	@param sdkClient
func loadNodeInfo(sdkClient *SdkClient) {
	sdkClient.lock.Lock()
	defer sdkClient.lock.Unlock()
	chainId := sdkClient.ChainId
	if chainId == "" {
		return
	}

	//获取链上节点和需要删除的节点
	nodeList, deleteNodeIds := GetChainNodeData(sdkClient)
	//删除节点数据
	if len(deleteNodeIds) > 0 {
		//deErr := dbhandle.UpdateNode(chainId, deleteNodeIds, config.StatusDeleted)
		deErr := dbhandle.DeleteNodeById(chainId, deleteNodeIds)
		if deErr != nil {
			log.Error("[DB] Delete Node Info Failed : " + deErr.Error())
			return
		}
	}

	var nodes []*db.Node
	//public模式
	authType := sdkClient.ChainInfo.AuthType
	if authType == config.PUBLIC {
		chainConfig := sdkClient.ChainConfig
		consensusNodes := chainConfig.GetConsensus().Nodes
		consensusNodeIds := make(map[string]int)
		if len(consensusNodes) > 0 {
			for _, nodeId := range consensusNodes[0].NodeId {
				consensusNodeIds[nodeId] = 0
			}
		}
		nodes = parsePublicNodeList(nodeList, consensusNodeIds)
	} else {
		nodes = parseNodeList(nodeList)
	}

	err := dbhandle.BatchInsertNode(chainId, nodes)
	if err != nil {
		log.Error("[DB] Update Node Info Failed : " + err.Error())
	}
}

// GetChainNodeData
//
//	@Description:  获取链上节点和需要删除的节点
//	@param sdkClient
//	@return []*discovery.Node
//	@return []string
func GetChainNodeData(sdkClient *SdkClient) ([]*discovery.Node, []string) {
	deleteNodeIds := make([]string, 0)
	chainNodeIdMap := make(map[string]string)
	chainId := sdkClient.ChainId
	//链配置
	chainClient := sdkClient.ChainClient
	//获取链上节点数据
	var nodeList []*discovery.Node
	chainInfo, err := chainClient.GetChainInfo()
	if err != nil || chainInfo == nil || len(chainInfo.NodeList) == 0 {
		log.Errorf("[SDK] Get Chain Info Failed : %v", err)
		return nodeList, deleteNodeIds
	}

	log.Infof("【loadNodeInfo】chainClient GetChainInfo, chainInfo:%v", chainInfo)
	nodeList = chainInfo.NodeList
	//if len(chainInfo.NodeList) == 0 {
	//	//获取默认节点
	//	nodeList, err = DealChainConfigError(sdkClient)
	//	if err != nil || len(nodeList) == 0 {
	//		log.Infof("Get chain node Failed : %v", err)
	//		return nodeList, deleteNodeIds
	//	}
	//}

	//链上节点
	for _, node := range nodeList {
		chainNodeIdMap[node.NodeId] = node.NodeId
	}

	//数据库节点
	nodeIds, err := dbhandle.GetNodesRef(chainId)
	if err != nil {
		log.Errorf("Get nodeIds fail, err:%v", err)
		return nodeList, deleteNodeIds
	}

	//链上没有的节点在数据库中删除
	for _, dbNodeId := range nodeIds {
		if _, ok := chainNodeIdMap[dbNodeId]; !ok {
			deleteNodeIds = append(deleteNodeIds, dbNodeId)
		}
	}
	return nodeList, deleteNodeIds
}

// parseNodeList parse node information
func parsePublicNodeList(nodeList []*discovery.Node, consensusNodeIds map[string]int) []*db.Node {
	nodes := make([]*db.Node, 0)
	for _, v := range nodeList {
		node := db.Node{
			NodeId:  v.GetNodeId(),
			Address: v.GetNodeAddress(),
		}
		if _, ok := consensusNodeIds[v.GetNodeId()]; ok {
			node.Role = "consensus"
		} else {
			node.Role = "common"
		}
		nodes = append(nodes, &node)
	}
	return nodes
}

// parseNodeList parse node information
func parseNodeList(nodeList []*discovery.Node) []*db.Node {
	nodes := make([]*db.Node, 0)
	for _, v := range nodeList {
		node := db.Node{}
		//节点地址
		node.Address = v.GetNodeAddress()
		node.NodeId = v.GetNodeId()
		nodeName, orgIds, roles, err := parseNodeInfo(v)
		if err != nil {
			continue
		}
		orgId, role := strings.Join(orgIds, ","), strings.Join(roles, ",")
		node.NodeName = nodeName
		node.OrgId = orgId
		node.Role = role
		nodes = append(nodes, &node)
	}
	return nodes
}

// parseNodeInfo parse node information
func parseNodeInfo(node *discovery.Node) (string, []string, []string, error) {
	// return OrgId/Role
	_, rest := pem.Decode(node.GetNodeTlsCert())
	if rest == nil {
		log.Error("can not decode tls cert")
		return "", nil, nil, errors.New("can not decode tls cert")
	}
	cert, err := x509.ParseCertificate(rest)
	if err != nil {
		return "", nil, nil, err
	}
	return cert.Subject.CommonName, cert.Subject.Organization, cert.Subject.OrganizationalUnit, nil
}
