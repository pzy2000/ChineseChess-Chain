package msgchannel

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"chainmaker_web/src/config"
)

// WechatAlarmer wechat
type WechatAlarmer struct {
}

// SendAlarm send
func (alarmer *WechatAlarmer) SendAlarm(msg string) error {

	webHook := `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=` + config.GlobalConfig.AlarmerConf.WechatAccessToken
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
