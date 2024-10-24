/*
Package alarms comment
Copyright (C) BABEC. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package alarms

import (
	"bytes"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"fmt"
	"time"
)

type safeWord struct {
	txID               string
	contractParameters string
	contractResult     string
}

// 敏感词历史
// 查询一段时间（可配置，当前为10min）内的敏感词列表，直接告警
func monitorSafeWord() {
	log.Infof("[Monitor] monitorSafeWord start")
	var flag bool
	problem := make(map[string][]safeWord)
	nowTime := time.Now().Unix()
	startTime := time.Now().Add(-time.Minute * time.Duration(config.GlobalConfig.MonitorConf.Interval)).Unix()
	txLists := make([]*db.Transaction, 0)
	for _, chainConfig := range config.GlobalConfig.ChainsConfig {
		txs, err := dbhandle.GetSafeWordTransactionList(chainConfig.ChainId, startTime, nowTime)
		if err != nil {
			log.Errorf("GetTxList err : %s", err.Error())
			return
		}
		txLists = append(txLists, txs...)
	}

	flag = false
	for i, tx := range txLists {
		if i > config.GlobalConfig.MonitorConf.SafeWordLimit {
			break
		}
		if tx.ContractParametersBak != "" ||
			bytes.Equal(tx.ContractResultBak, config.ContractResultMsg) {
			flag = true
			sender := tx.Sender
			s := safeWord{txID: tx.TxId, contractParameters: tx.ContractParametersBak,
				contractResult: string(tx.ContractResultBak)}
			problem[sender] = append(problem[sender], s)
		}
	}
	if flag {
		msg := fmt.Sprintf("【%s告警】敏感词列表：%v", config.GlobalConfig.AlarmerConf.Prefix, problem)
		log.Warn(msg)
		sendMsg(msg)
	}
}
