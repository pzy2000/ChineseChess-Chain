package sensitive

import (
	"chainmaker_web/src/config"
)

var (
	sensitiveEnable bool
)

// SetSensitiveConfig - 设置敏感词配置
func SetSensitiveConfig(conf *config.SensitiveConf) {
	sensitiveEnable = conf.Enable
}

// GetSensitiveEnable - 获取是否使用敏感词
func GetSensitiveEnable() bool {
	return sensitiveEnable
}
