package sensitive

import (
	"chainmaker_web/src/config"
	"testing"
)

func TestGetSensitiveEnable(t *testing.T) {
	tests := []struct {
		name        string
		setValue    bool
		wantEnabled bool
	}{
		{
			name:        "启用敏感词",
			setValue:    true,
			wantEnabled: true,
		},
		{
			name:        "禁用敏感词",
			setValue:    false,
			wantEnabled: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetSensitiveConfig(&config.SensitiveConf{Enable: tt.setValue})
			if got := GetSensitiveEnable(); got != tt.wantEnabled {
				t.Errorf("GetSensitiveEnable() = %v, want %v", got, tt.wantEnabled)
			}
		})
	}
}

func TestSetSensitiveConfig(t *testing.T) {
	type args struct {
		conf *config.SensitiveConf
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "设置敏感词配置",
			args: args{
				conf: &config.SensitiveConf{Enable: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetSensitiveConfig(tt.args.conf)
			// 添加断言以检查敏感词配置是否已正确设置
			if got := GetSensitiveEnable(); got != tt.args.conf.Enable {
				t.Errorf("SetSensitiveConfig() didn't set the value correctly, got = %v, want %v", got, tt.args.conf.Enable)
			}
		})
	}
}
