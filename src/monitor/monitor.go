/*
Package alarms comment
Copyright (C) BABEC. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package alarms

import (
	"chainmaker_web/src/config"
	loggers "chainmaker_web/src/logger"
	"chainmaker_web/src/monitor/msgchannel"
	"chainmaker_web/src/sync"
	"time"
)

var (
	log = loggers.GetLogger(loggers.MODULE_WEB)
)

// Start 启动监控
func Start(monitorConfig *config.MonitorConf) {
	log.Infof("[Monitor] isenable: %v", monitorConfig.Enable)
	if !monitorConfig.Enable {
		return
	}

	log.Infof("[Monitor] Init sdk clients start")
	monitorChains := config.SubscribeChains
	//如果配置链监控的链则值监控配置的链数据，否则监控所有订阅链
	if len(config.MonitorChains) > 0 {
		monitorChains = config.MonitorChains
	}

	interval := config.GlobalConfig.MonitorConf.Interval
	if interval < 10 {
		interval = 10
	}

	//建立链节点连接
	monitorClients := GetMonitorSdkClients(monitorChains)
	allSdkClients := GetMonitorSdkClients(config.SubscribeChains)
	//interval触发一次
	ticker := time.NewTicker(time.Minute * time.Duration(interval))
	for range ticker.C {
		//所有节点都监控
		MonitorChain(allSdkClients)
		for i, chainSdkClients := range monitorClients {
			if len(chainSdkClients) > 0 {
				//只监控配置的节点
				Monitor(chainSdkClients, i)
			}
		}
	}
}

// GetMonitorSdkClients get monitor client 所有节点的client
func GetMonitorSdkClients(monitorChains []*config.ChainInfo) [][]*sync.SdkClient {
	clients := make([][]*sync.SdkClient, 0)
	for _, chainInfo := range monitorChains {
		sdkClients := make([]*sync.SdkClient, 0)
		tempChain := chainInfo
		for _, node := range chainInfo.NodesList {
			tempNode := make([]*config.NodeInfo, 0)
			tempNode = append(tempNode, node)
			tempChain.NodesList = tempNode

			chainClient, err := sync.CreateChainClient(tempChain)
			sdkClient := sync.NewSdkClient(tempChain, chainClient)
			if err != nil {
				log.Error("创建chain Client失败: ", err.Error())
				continue
			}
			sdkClients = append(sdkClients, sdkClient)
		}
		log.Infof("[Monitor] init sdk clients success.nodes num : %d", len(sdkClients))
		clients = append(clients, sdkClients)
	}

	return clients
}

// Monitor 监控链指定节点
func Monitor(sdkClients []*sync.SdkClient, index int) {
	//初始化配置
	//节点高度差检测
	go monitorBlockHeights(sdkClients, index)
	//节点探活
	go monitorNode(sdkClients, index)
	//链出块检测
	//go monitorBlockPut(sdkClients, index)
}

// MonitorChain 对所有链的监控
func MonitorChain(allSdkClients [][]*sync.SdkClient) {
	//敏感词监控
	go monitorSafeWord()
	//单用户短时间大批量交易探测
	go monitorTx()

	//sdkClients := GetMonitorSdkClients(config.SubscribeChains)
	if len(allSdkClients) == 0 {
		return
	}

	for _, chainSdkClients := range allSdkClients {
		if len(chainSdkClients) == 0 {
			continue
		}

		//区块高度监控
		go monitorHeightDiff(chainSdkClients)
	}
}

func sendMsg(msg string) {
	msgchannel.SendMsg(msg)
}
