/*
Package config comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/hokaccha/go-prettyjson"

	"github.com/spf13/viper"
)

const (
	DefaultMaxByteSize      = 256
	DefaultMaxPoolCount     = 50
	DefaultPositionRankTime = 60
)

var (
	// GlobalConfig global
	GlobalConfig *Config
	//gConfPath gConfPath
	gConfPath       string
	GlobalAbiERC20  *abi.ABI
	GlobalAbiERC721 *abi.ABI
	//MaxDBByteSize 数据库批量插入最大字节,单位byte
	MaxDBByteSize int
	MaxDBPoolSize int
)

// Chain 链基本数据
type Chain struct {
	ChainId  string
	AuthType string
}

// SubscribeChains 订阅的所有链数据
var SubscribeChains = make([]*ChainInfo, 0)

// MonitorChains 监控链
var MonitorChains = make([]*ChainInfo, 0)

// printLog 输出日志
func (c *Config) printLog(env string) {
	if env == "" {
		return
	}

	json, err := prettyjson.Marshal(c)
	if err != nil {
		log.Fatalf("marshal alarm config failed, %s", err.Error())
	}
	fmt.Println(string(json))
}

// InitConfig init
func InitConfig(confPath, env string) *Config {
	gConfPath = confPath
	if gConfPath == "" {
		gConfPath = GetConfigDirPath()
	}
	webViper, err := initCMViper(gConfPath, env)
	if err != nil {
		fmt.Println("can not load config.yml, exit")
		panic(err)
	}
	browserConfig := &Config{}
	if err = webViper.Unmarshal(&browserConfig); err != nil {
		log.Fatal("Unmarshal config failed, ", err)
	}
	browserConfig.printLog(env)

	//数据库批量插入最大字节，kb转换成byte1
	MaxDBByteSize = DefaultMaxByteSize * 1024
	MaxDBPoolSize = DefaultMaxPoolCount
	if browserConfig.DBConf != nil {
		if browserConfig.DBConf.MaxByteSize > DefaultMaxByteSize {
			MaxDBByteSize = browserConfig.DBConf.MaxByteSize * 1024
		}
		if browserConfig.DBConf.MaxPoolSize < DefaultMaxPoolCount {
			MaxDBPoolSize = browserConfig.DBConf.MaxPoolSize
		}
	}

	if browserConfig.RedisDB != nil &&
		browserConfig.RedisDB.PositionRankTime <= 0 {
		browserConfig.RedisDB.PositionRankTime = DefaultPositionRankTime
	}

	//配置全局变量
	GlobalConfig = browserConfig

	GlobalAbiERC20 = getERC20AbiJson()
	GlobalAbiERC721 = getERC721AbiJson()

	//获取配置中所有订阅数据
	SubscribeChains = GetChainsInfoAll()
	//获取监控链
	MonitorChains = getMonitorChainsConfig()
	return browserConfig
}

// getERC20AbiJson 解析EVM abi
func getERC20AbiJson() *abi.ABI {
	abiPath := GlobalConfig.SubscribeConfig.ERC20AbiFile
	if abiPath == "" {
		return nil
	}
	ercAbiJson, err := os.Open(abiPath)
	if err != nil {
		return nil
	}
	ercAbi, err := abi.JSON(ercAbiJson)
	if err != nil {
		log.Fatalf("readFile erc20_abi failed: %v", err)
		return nil
	}
	return &ercAbi
}

// getERC20AbiJson 解析EVM abi
func getERC721AbiJson() *abi.ABI {
	abiPath := GlobalConfig.SubscribeConfig.ERC721AbiFile
	if abiPath == "" {
		return nil
	}
	ercAbiJson, err := os.Open(abiPath)
	if err != nil {
		//log.Fatalf("readFile erc721_abi failed: %v", err)
		return nil
	}
	ercAbi, err := abi.JSON(ercAbiJson)
	if err != nil {
		log.Fatalf("readFile erc721_abi failed: %v", err)
		return nil
	}

	return &ercAbi
}

// GetConfigDirPath 绝对路径
func GetConfigDirPath() string {
	_, currentFilePath, _, _ := runtime.Caller(0)
	configDir := filepath.Join(filepath.Dir(currentFilePath), "../../", "configs")
	return configDir
}

// initCMViper
func initCMViper(gConfPath, env string) (*viper.Viper, error) {

	cmViper := viper.New()
	// 使用 env 参数构建配置文件名
	configFilePath := GetConfigFilePath(gConfPath, env)
	cmViper.SetConfigFile(configFilePath)
	if err := cmViper.ReadInConfig(); err != nil {
		return nil, err
	}
	return cmViper, nil
}

// GetConfigFilePath 配置文件路径
func GetConfigFilePath(gConfPath, env string) string {
	configFile := "config.yml"
	if env != "" {
		configFile = fmt.Sprintf("config.%s.yml", env)
	}
	return gConfPath + "/" + configFile
}

// GetChainsInfoAll 获取订阅链信息
func GetChainsInfoAll() []*ChainInfo {
	chainList := make([]*ChainInfo, 0)
	//如果没有订阅信息使用配置信息
	if len(GlobalConfig.ChainsConfig) == 0 {
		return chainList
	}

	for _, chain := range GlobalConfig.ChainsConfig {
		chainInfo := buildChainInfo(chain)
		chainList = append(chainList, chainInfo)
	}

	return chainList
}

// getMonitorChainsConfig 获取监控链配置信息
func getMonitorChainsConfig() []*ChainInfo {
	chainList := make([]*ChainInfo, 0)
	if GlobalConfig == nil || GlobalConfig.MonitorConf == nil ||
		GlobalConfig.MonitorConf.ChainsConfig == nil {
		return chainList
	}

	if !GlobalConfig.MonitorConf.Enable {
		return chainList
	}

	//处理配置监控链数据
	for _, chain := range GlobalConfig.MonitorConf.ChainsConfig {
		chainInfo := buildChainInfo(chain)
		chainList = append(chainList, chainInfo)
	}
	return chainList
}

// buildChainInfo buildChainInfo
func buildChainInfo(chain *ChainConfig) *ChainInfo {
	if chain == nil {
		return &ChainInfo{}
	}

	chainInfo := &ChainInfo{
		ChainId:  chain.ChainId,
		AuthType: chain.AuthType,
		OrgId:    chain.OrgId,
		HashType: chain.HashType,
	}
	nodesList := make([]*NodeInfo, 0)
	userInfo := &UserInfo{
		UserKey:  "",
		UserCert: "",
	}

	for _, node := range chain.NodesConfig {
		orgCA := ReadAbsPathFile(node.CaPaths)
		nodesList = append(nodesList, &NodeInfo{
			Addr:        node.Remotes,
			OrgCA:       orgCA,
			TLSHostName: node.TlsHost,
			Tls:         node.Tls,
		})
	}
	chainInfo.NodesList = nodesList

	//根据路径读取密钥文件
	userInfo.UserKey = ReadAbsPathFile(chain.UserConf.PrivKeyFile)
	userInfo.UserCert = ReadAbsPathFile(chain.UserConf.CertFile)
	chainInfo.UserInfo = userInfo
	return chainInfo
}

func ReadAbsPathFile(path string) string {
	if path == "" {
		return ""
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		// 处理错误
		fmt.Println("Error getting absolute path:", err)
		return ""
	}
	pathBytes, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Println("Error ReadFile path: err:", absPath, err)
		return ""
	}

	return string(pathBytes)
}

// ToClickHouseUrl to
func (dbConfig *DBConf) ToClickHouseUrl(useDataBase bool) string {
	database := "default"
	if useDataBase && dbConfig.Database != "" {
		database = dbConfig.Database
	}
	connStr := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, database)
	return connStr
}

// ToMysqlUrl to
func (dbConfig *DBConf) ToMysqlUrl(useDataBase bool) string {
	var url string
	if useDataBase {
		url = fmt.Sprintf("tcp(%s:%s)/%s", dbConfig.Host, dbConfig.Port, dbConfig.Database)
	} else {
		url = fmt.Sprintf("tcp(%s:%s)/", dbConfig.Host, dbConfig.Port)
	}
	return dbConfig.Username + ":" + dbConfig.Password + "@" + url + MysqlDefaultConf
}

// ToPgsqlUrl to
func (dbConfig *DBConf) ToPgsqlUrl(useDataBase bool) string {
	var url string
	if useDataBase {
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable client_encoding=UTF8",
			dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Database)
	} else {
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable client_encoding=UTF8",
			dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, "template1")
	}
	return url
}

// ToUrl to
func (webConfig *WebConf) ToUrl() string {
	return webConfig.Address + ":" + strconv.Itoa(webConfig.Port)
}
