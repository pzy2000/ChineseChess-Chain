package alarms

import (
	"chainmaker_web/src/monitor_prometheus"
	"chainmaker_web/src/utils"
)

var (
	nodeBlockHeightGauge = monitor_prometheus.NewGaugeVec(utils.MonitorNameSpace, "node_block_height",
		"chain node block height", "chainId", "remotes") // 当前节点的区块链高度
	//txNotOnChain = monitor_prometheus.NewCounterVec(utils.MonitorNameSpace, "put_tx_failed",
	//"put tx failed in 1 minute", "chainId") // 1分钟内未成功上链
	explorerAfterChain = monitor_prometheus.NewCounterVec(utils.MonitorNameSpace, "sync_blk_large",
		"explorer sync block height diff chain node block height exceeds the gate",
		"chainId", "dbHeight", "chainHeight")
	probeChainNodeFailed = monitor_prometheus.NewCounterVec(utils.MonitorNameSpace, "probe_node_failed",
		"explorer probe chainNode failed", "chainId")
	singUserAbnormal = monitor_prometheus.NewCounterVec(utils.MonitorNameSpace, "single_user_abnormal",
		"single user too much tx in short time ", "chainId", "userId", "txNum")
)
