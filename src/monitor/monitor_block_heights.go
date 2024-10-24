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

// 区块高度不一致检测
// 连接4个共识 4个同步节点，获取区块高度，若区块高度差大于10则告警
func monitorBlockHeights(sdkClients []*sync.SdkClient, index int) {
	bhs := make(map[string]int64, 0)
	for _, c := range sdkClients {
		if len(c.ChainInfo.NodesList) == 0 {
			continue
		}

		chainClient := c.ChainClient
		nodeInfo := c.ChainInfo.NodesList[0]
		bh, err := chainClient.GetCurrentBlockHeight()
		if err != nil {
			log.Error("[SDK] Get chain block height Failed : " + err.Error())
		}
		bhs[nodeInfo.Addr] = int64(bh)
		nodeBlockHeightGauge.WithLabelValues(c.ChainId, nodeInfo.Addr).Set(float64(bh))
	}

	num := 0
	problem := make(map[string]int64, 0)
	maxHeight, maxNode := maxNum(bhs)
	heightDiff := config.GlobalConfig.MonitorConf.MaximumHeightDiff
	for nodeIp, height := range bhs {
		if height+heightDiff < maxHeight {
			problem[nodeIp] = height
			num++
		}
	}

	if num > 0 {
		prefix := config.GlobalConfig.AlarmerConf.Prefix
		msg := fmt.Sprintf("【%s告警】底链区块高度不一致：当前订阅节点最高区块高度为：%s[%d]， 落后区块节点有： %v",
			prefix, maxNode, maxHeight, problem)
		log.Warn(msg)
		sendMsg(msg)
	}
}

// maxNum maxNum
func maxNum(arr map[string]int64) (max int64, maxNode string) {
	for ip, height := range arr {
		if max < height {
			max = height
			maxNode = ip
		}
	}
	return max, maxNode
}
