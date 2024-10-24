/*
Package msgchannel comment
Copyright (C) BABEC. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package msgchannel

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"chainmaker_web/src/config"
	loggers "chainmaker_web/src/logger"
)

const ContentMsg = `{"msgtype": "text","text": {"content": %s}}`

var (
	log = loggers.GetLogger(loggers.MODULE_WEB)
)

// DingAlarmer ding
type DingAlarmer struct {
}

// SendAlarm send
func (alarmer *DingAlarmer) SendAlarm(msg string) error {
	webHook := `https://oapi.dingtalk.com/robot/send?access_token=` + config.GlobalConfig.AlarmerConf.DingAccessToken
	//content := `{"msgtype": "text",
	//		"text": {"content": "` + msg + `"}
	//	}`
	content := fmt.Sprintf(ContentMsg, msg)
	//创建一个请求
	req, err := http.NewRequest("POST", webHook, strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("create requset  err, %s", err.Error())
	}
	//nolint:gosec
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Warn("send alarm err, %s \n", err.Error())
		return fmt.Errorf("send alarm err, %s", err.Error())
	}
	//关闭请求
	resp.Body.Close()
	return err
}
