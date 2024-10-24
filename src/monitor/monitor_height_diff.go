/*
Package alarms comment
Copyright (C) BABEC. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package alarms

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/sync"
	"fmt"
	"strconv"
	"time"
)

var dbCurrentHeight int64

// 浏览器区块高度与链高度差距过大告警
// 查询链高度1min后再查询浏览器高度，浏览器高度应大于等于链高度
func monitorHeightDiff(sdkClients []*sync.SdkClient) {
	for _, sdkClient := range sdkClients {
		msg := singleHandle(sdkClient)
		if msg != "" {
			log.Warn(msg)
			sendMsg(msg)
		}
	}
}

// 区块高度监控
func singleHandle(sdkClient *sync.SdkClient) string {
	chainClient := sdkClient.ChainClient
	chainId := sdkClient.ChainId
	nodeRemotes := ""
	if len(sdkClient.ChainInfo.NodesList) > 0 {
		nodeRemotes = sdkClient.ChainInfo.NodesList[0].Addr
	}
	block, err := chainClient.GetLastBlock(false)
	if err != nil {
		msg := GetLastBlockMsg(chainId, nodeRemotes, err)
		return msg
	}
	bcHeight := block.Block.Header.BlockHeight
	timestamp := block.Block.Header.BlockTimestamp

	<-time.After(time.Second * 60)
	dbHeight := dbhandle.GetMaxBlockHeight(chainId)
	if dbHeight == 0 {
		msg := GetMaxBlockHeightMsg(chainId, nodeRemotes)
		return msg
	}

	heightDiffConf := config.GlobalConfig.MonitorConf.MaximumHeightDiff
	heightDiff := int64(bcHeight) - dbHeight
	if heightDiff > heightDiffConf {
		if dbCurrentHeight == dbHeight {
			msg := GetBackwardMsg(chainId, nodeRemotes)
			log.Warn(msg)
		}
		msg := GetBackwardDiffMsg(chainId, nodeRemotes, heightDiff, int64(bcHeight), dbHeight, timestamp)
		explorerAfterChain.WithLabelValues(sdkClient.ChainId,
			strconv.FormatInt(dbHeight, 10), strconv.FormatInt(int64(bcHeight), 10)).Inc()
		dbCurrentHeight = dbHeight
		return msg
	}

	return ""
}

// GetLastBlockMsg GetLastBlockMsg
func GetLastBlockMsg(chainId, nodeRemotes string, blockErr error) string {
	prefix := config.GlobalConfig.AlarmerConf.Prefix
	msg := fmt.Sprintf("【%s告警】chainId[%s] nodeIP[%s] 获取最新区块失败, err:%v", prefix, chainId,
		nodeRemotes, blockErr)
	return msg
}

// GetMaxBlockHeightMsg GetMaxBlockHeightMsg
func GetMaxBlockHeightMsg(chainId, nodeRemotes string) string {
	prefix := config.GlobalConfig.AlarmerConf.Prefix
	msg := fmt.Sprintf("【%s告警】chainId[%s] nodeIP[%s] 还未开始同步区块信息, 区块高度0", prefix, chainId, nodeRemotes)
	return msg
}

// GetBackwardMsg GetBackwardMsg
func GetBackwardMsg(chainId, nodeRemotes string) string {
	prefix := config.GlobalConfig.AlarmerConf.Prefix
	msg := fmt.Sprintf("【%s告警】chainId[%s] nodeIP[%s] 落后很久了，长时间未同步区块了", prefix, chainId, nodeRemotes)
	return msg
}

// GetBackwardDiffMsg GetBackwardDiffMsg
func GetBackwardDiffMsg(chainId, nodeRemotes string, diffHeight, bcHeight, dbHeight, timestamp int64) string {
	prefix := config.GlobalConfig.AlarmerConf.Prefix
	msg := fmt.Sprintf("【%s告警】chainId[%s] nodeIP[%s] 浏览器落后底链高度%d，底链高度为：%d 浏览器高度为：%d 区块时间:%d 当前时间:%d",
		prefix, chainId, nodeRemotes, diffHeight, bcHeight, dbHeight, timestamp, time.Now().Unix())
	return msg
}
