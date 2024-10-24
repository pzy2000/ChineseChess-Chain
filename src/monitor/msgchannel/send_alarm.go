package msgchannel

import (
	"chainmaker_web/src/config"
)

// SendMsg SendMsg
func SendMsg(msg string) {
	if config.GlobalConfig.AlarmerConf.DingEnable {
		ding := &DingAlarmer{}
		err := ding.SendAlarm(msg)
		if err != nil {
			log.Warn("send dingding alarm err , %s \n", err.Error())
		}
	}

	if config.GlobalConfig.AlarmerConf.WechatEnable {
		w := &WechatAlarmer{}
		err := w.SendAlarm(msg)
		if err != nil {
			log.Warn("send wechat alarm err , %s \n", err.Error())
		}
	}
}
