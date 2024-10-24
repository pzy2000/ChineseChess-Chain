/*
Package alarms comment
Copyright (C) BABEC. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package alarms

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/sync"
	"fmt"
)

// 节点探活
// 连接节点，查询链配置，成功则OK，失败3次则告警。
func monitorNode(sdkClients []*sync.SdkClient, index int) {
	log.Infof("[Monitor] monitor node start,sdk clients num:%d", len(sdkClients))
	//nodeStatus := make(map[string]bool)
	prefix := config.GlobalConfig.AlarmerConf.Prefix
	for _, c := range sdkClients {
		tryNum := config.GlobalConfig.MonitorConf.TryConnNum
		if len(c.ChainInfo.NodesList) == 0 {
			continue
		}
		nodeInfo := c.ChainInfo.NodesList[0]
		nodeStatus := false
		chainClient := c.ChainClient
		for j := 0; j < tryNum; j++ {
			_, err := chainClient.GetChainConfig()
			if err != nil {
				log.Errorf("[SDK] %dst monitor node Get Chain Config Failed : %s", j+1, err.Error())
				continue
			}
			nodeStatus = true
			break
		}

		if !nodeStatus {
			msg := fmt.Sprintf("【%s告警】底链查询链配置信息失败,节点【%s】失活", prefix, nodeInfo.Addr)
			log.Warn(msg)
			probeChainNodeFailed.WithLabelValues(config.GlobalConfig.MonitorConf.ChainsConfig[index].ChainId).Inc()
			sendMsg(msg)
		}
	}
}
