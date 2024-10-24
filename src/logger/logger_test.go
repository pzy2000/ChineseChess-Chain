/*
Package loggers comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package loggers

import (
	"chainmaker_web/src/config"
	"flag"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

const (
	configParam = "config"
)

func configYml() string {
	configPath := flag.String(configParam, "../../configs", "config.yml's path")
	flag.Parse()
	return *configPath
}

func TestDefaultLogConfig(t *testing.T) {
	tests := []struct {
		name string
		want *config.LogConf
	}{
		{
			name: "Test case 1:",
			want: &config.LogConf{
				LogLevelDefault: "INFO",
				FilePath:        "../log/web.log",
				MaxAge:          365,
				RotationTime:    6,
				LogInConsole:    true,
				ShowColor:       false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultLogConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultLogConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefaultLogNodeConfig(t *testing.T) {
	tests := []struct {
		name string
		want config.LogConf
	}{
		{
			name: "Test case 1:",
			want: config.LogConf{
				LogLevelDefault: "INFO",
				FilePath:        "../log/web.log",
				MaxAge:          365,
				RotationTime:    6,
				LogInConsole:    true,
				ShowColor:       true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDefaultLogNodeConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDefaultLogNodeConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLoggerByChain(t *testing.T) {
	type args struct {
		name    string
		chainId string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1:",
			args: args{
				name:    "name1",
				chainId: "chain1",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = GetLoggerByChain(tt.args.name, tt.args.chainId)
			logHeader := tt.args.name + tt.args.chainId
			if _, ok := loggers[logHeader]; ok != tt.want {
				t.Errorf("GetLoggerByChain() = %v, want %v", ok, tt.want)
			}
		})
	}
}

func TestRefreshLogConfig(t *testing.T) {
	type args struct {
		config *config.LogConf
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1:",
			args: args{
				config: &config.LogConf{
					LogLevels: map[string]string{
						"name1": "debug",
					},
					LogLevelDefault: "info",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 初始化 loggerLevels
			loggerLevels["name1"] = make(map[string]zap.AtomicLevel)
			loggerLevels["name1"]["name1chain1"] = zap.NewAtomicLevelAt(zap.InfoLevel)

			// 调用 RefreshLogConfig
			RefreshLogConfig(tt.args.config)

			// 检查 loggerLevels 是否已更新为预期的值
			expectedLevel := zap.DebugLevel
			actualLevel := loggerLevels["name1"]["name1chain1"].Level()

			if actualLevel != expectedLevel {
				t.Errorf("RefreshLogConfig() = %v, want %v", actualLevel, expectedLevel)
			}
		})
	}
}

func TestSetLogConfig(t *testing.T) {
	type args struct {
		config *config.LogConf
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1:",
			args: args{
				config: &config.LogConf{
					LogLevels: map[string]string{
						"name1": "debug",
					},
					LogLevelDefault: "info",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLogConfig(tt.args.config)
		})
	}
}
