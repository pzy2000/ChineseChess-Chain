/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/monitor_prometheus"
	"chainmaker_web/src/utils"
)

var (
	sensitiveCallTotal = monitor_prometheus.NewCounterVec(utils.MonitorNameSpace, "sensitive_call",
		"sensitive call count", "success", "hit") //敏感词调用总次数
	callSensitiveSuccess = "1"
	callSensitiveFailed  = "0"
	callSensitiveHit     = "1"
	callSensitiveNotHit  = "0"
	syncHeightGauge      = monitor_prometheus.NewGaugeVec(utils.MonitorNameSpace, "sync_block_height",
		"has synced block height", "chainId") // 当前数据库已同步块高
	//processBlockDataHisto = monitor_prometheus.NewHistogramVec(utils.MonitorNameSpace, "process_block",
	//	"process block data",
	//	[]float64{0.01, 0.1, 0.3, 0.5, 0.7, 2.0, 3.0, 5.0, 7.0, 10, 15, 20, 30},
	//	"chainId")
	//insertblockDataHisto = monitor_prometheus.NewHistogramVec(utils.MonitorNameSpace, "insert_block",
	//	"process block data",
	//	[]float64{0.01, 0.1, 0.3, 0.5, 0.7, 2.0, 3.0, 5.0, 7.0, 10, 15, 20, 30},
	//	"chainId")
	//syncBlockDataHisto = monitor_prometheus.NewHistogramVec(utils.MonitorNameSpace, "sync_block",
	//	"process block data",
	//	[]float64{0.01, 0.1, 0.3, 0.5, 0.7, 2.0, 3.0, 5.0, 7.0, 10, 15, 20, 30},
	//	"chainId")
	//subscribeSuccess = monitor_prometheus.NewCounterVec(utils.MonitorNameSpace, "subscribe_success_counter",
	//	"success subscribe count", "chainId") //成功订阅到区块数据
	subscribeFail = monitor_prometheus.NewCounterVec(utils.MonitorNameSpace, "subscribe_fail_counter",
		"fail subscribe count", "chainId") //失败订阅到区块数据
	////subscribeRetry = monitor_prometheus.NewCounterVec(utils.MonitorNameSpace, "subscribe_retry_counter",
	////	"retry subscribe count", "chainId") //失败订阅到区块数据
	//
	//==insertBlockDataIntoDbFailedTotal = monitor_prometheus.NewCounterVec(utils.MonitorNameSpace,
	//	"insert_blk_failed", "insert block data into es failed", "chainId", "method") //敏感词调用总次数

)
