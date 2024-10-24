/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/chain"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	loggers "chainmaker_web/src/logger"
	alarms "chainmaker_web/src/monitor"
	"chainmaker_web/src/monitor_prometheus"
	"chainmaker_web/src/router"
	"chainmaker_web/src/sensitive"
	"chainmaker_web/src/sync"
	"chainmaker_web/src/utils"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

const (
	configParam = "config"
)

func main() {
	// 解析命令行参数：config
	confYml, envConfig := configYml()
	//初始化配置
	conf := config.InitConfig(confYml, envConfig)

	if conf.PProf != nil && conf.PProf.IsOpen {
		// 初始化pprof HTTP服务器
		go pprofServer(conf.PProf.Port)
	}

	// 初始化缓存
	cache.InitRedis(conf.RedisDB)
	// 初始化日志配置
	loggers.SetLogConfig(conf.LogConf)
	// 初始化数据库配置
	db.InitDbConn(conf.DBConf)
	// 初始化敏感词配置
	sensitive.SetSensitiveConfig(conf.SensitiveConf)
	// 初试化链配置信息
	chain.InitChainConfig()
	// 初始化prometheus监控
	mServer := monitor_prometheus.NewMonitorServer(utils.MonitorNameSpace, conf.WebConf.MonitorPort)
	_ = mServer.Start()
	if conf.SubscribeConfig == nil || conf.SubscribeConfig.Enable {
		// 启动同步
		go sync.StartSync(config.SubscribeChains)
	}
	// 启动告警监控
	go alarms.Start(conf.MonitorConf)
	// http-server启动
	router.HttpServe(conf.WebConf)

	log.Println("Program has finished running")
}

func configYml() (string, string) {
	configPath := flag.String(configParam, "configs", "config.yml's path")
	env := flag.String("env", "", "yml file name")

	flag.Parse()
	return *configPath, *env
}

func pprofServer(port string) {
	// 创建一个自定义的 http.Server 实例
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 使用 ListenAndServe 方法启动服务器
	err := server.ListenAndServe()
	if err != nil {
		log.Println("pprof ListenAndServe err:", err)
	}
}
