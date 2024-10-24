/*
Package alarms comment
Copyright (C) BABEC. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package alarms

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"fmt"
	"strconv"
	"time"
)

// 单用户短时间大批量交易探测
// 一段时间内（可配置，当前配置为10min）单用户交易量大于MaxTxNum（可配置，当前为10000），触发预警
func monitorTx() {
	for _, chainConfig := range config.SubscribeChains {
		singleTx(chainConfig.ChainId)
	}

}

func singleTx(chainId string) {
	limit := config.GlobalConfig.MonitorConf.MonitorTxConf.TxLimit
	offset := 0
	var users []*db.User
	for {
		user, total, err := dbhandle.GetUserList(offset, limit, chainId, "", nil, nil)
		offset++
		if err != nil {
			log.Errorf("GetUserList err : %s", err.Error())
		}
		users = append(users, user...)
		if offset*limit+1 > int(total) {
			break
		}
	}

	for _, user := range users {
		nowTime := time.Now().Unix()
		startTime := time.Now().Add(-time.Minute * time.Duration(config.GlobalConfig.MonitorConf.Interval)).Unix()
		num, err := dbhandle.GetTransactionNumByRange(chainId, user.UserAddr, startTime, nowTime)
		if err != nil {
			log.Errorf("GetTxList err : %s", err.Error())
			return
		}
		if int(num) > config.GlobalConfig.MonitorConf.MonitorTxConf.MaxTxNum {
			msg := fmt.Sprintf("【%s告警】用户[%s] %d 分钟内发起交易数异常：%d",
				config.GlobalConfig.AlarmerConf.Prefix,
				user.UserId, config.GlobalConfig.MonitorConf.Interval, num)
			log.Warn(msg)
			singUserAbnormal.WithLabelValues(chainId, user.UserId, strconv.FormatInt(num, 10)).Inc()
			sendMsg(msg)
		}
	}
}
