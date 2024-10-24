//nolint
package db

import (
	"chainmaker_web/src/config"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var nodesList = []*config.NodeInfo{
	{
		Addr:        "9.135.180.61:12301",
		OrgCA:       "-----BEGIN CERTIFICATE-----\nMIICnjCCAkSgAwIBAgIDBM78MAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTIzMTIwMTA4NDMxNFoXDTMz\nMTEyODA4NDMxNFowgYoxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMRIwEAYDVQQLEwlyb290LWNlcnQxIjAgBgNVBAMTGWNhLnd4LW9yZzEuY2hhaW5t\nYWtlci5vcmcwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARIrtmDyFGTcWqIJyM8\nweHuk+9GqTkllI9E59P4h3Ms/jP8xBaa815Zkh1y5WPqFxqyN5rfrRhRMp8LqoeU\nuF+To4GWMIGTMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MCkGA1Ud\nDgQiBCAPRq+/1wQPj8AkeVIyl8D6i0dgqvxy5euC+DF5WVuUNzBFBgNVHREEPjA8\ngg5jaGFpbm1ha2VyLm9yZ4IJbG9jYWxob3N0ghljYS53eC1vcmcxLmNoYWlubWFr\nZXIub3JnhwR/AAABMAoGCCqGSM49BAMCA0gAMEUCIQCSFT8YV2rsga4TyT/qs0Qp\nAv0aRTURq7XqEmnuX3fDGQIgC93pXvi6GY0T6beC80HR3ib/TmTQ8YvVsIt1p/Tk\nc6E=\n-----END CERTIFICATE-----\n",
		TLSHostName: "chainmaker.org",
		Tls:         true,
	},
}
var ChainListConfigTest = &config.ChainInfo{
	ChainId:   "chain1",
	AuthType:  "permissionedwithcert",
	OrgId:     "wx-org1.chainmaker.org",
	HashType:  "SHA256",
	NodesList: nodesList,
	UserInfo: &config.UserInfo{
		UserKey:  "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEILN+eElgD7gwq5t/Z/ZQ4JifW/8RC/YaVW2unaMko2AXoAoGCCqGSM49\nAwEHoUQDQgAEhw6eLytgitYparmL/ALv7k7GnCHHR8937bvQeihews0Df+QrsLAD\nwrfwfE8V8AeI72E2yHAX6LvFg7t2JFSlug==\n-----END EC PRIVATE KEY-----\n",
		UserCert: "-----BEGIN CERTIFICATE-----\nMIICdTCCAhygAwIBAgIDCW/zMAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTIzMTIwMTA4NDMxNFoXDTI4\nMTEyOTA4NDMxNFowgY8xCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMQ4wDAYDVQQLEwVhZG1pbjErMCkGA1UEAxMiYWRtaW4xLnNpZ24ud3gtb3JnMS5j\naGFpbm1ha2VyLm9yZzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABIcOni8rYIrW\nKWq5i/wC7+5Oxpwhx0fPd+270HooXsLNA3/kK7CwA8K38HxPFfAHiO9hNshwF+i7\nxYO7diRUpbqjajBoMA4GA1UdDwEB/wQEAwIGwDApBgNVHQ4EIgQgfs51BSgHvk4M\niDcCJPghxY3TGPUGIWuVWoZDyC37YtswKwYDVR0jBCQwIoAgD0avv9cED4/AJHlS\nMpfA+otHYKr8cuXrgvgxeVlblDcwCgYIKoZIzj0EAwIDRwAwRAIgE10Ns1fwCn1U\nfnDJnf0STF5Ipm36yscpvbBl0LPouQECIBs4UPGsU/gXlvzCGjYnCFGTfj4ny+F2\nuspvLrkFozei\n-----END CERTIFICATE-----\n",
	},
}

var (
	onceMySQL      sync.Once
	onceClickHouse sync.Once
	onceRedis      sync.Once
	mysqlCfg       *config.DBConf
	clickHouseCfg  *config.DBConf
	redisCfg       *config.RedisConfig
)

func InitMySQLContainer() (*config.DBConf, error) {
	var err error
	onceMySQL.Do(func() {
		//初始化配置
		_ = config.InitConfig("", "")
		//config.SubscribeChains = getChainListConfigTest("0_chainListConfig.json")
		config.SubscribeChains = []*config.ChainInfo{
			ChainListConfigTest,
		}
		mysqlCfg, err = CreateMySQLContainer()
		if mysqlCfg == nil {
			return
		}
		log.Infof("=========mysqlCfg:%v=", mysqlCfg)
		config.GlobalConfig.DBConf = mysqlCfg
		InitDbConn(mysqlCfg)
	})
	config.GlobalConfig.DBConf = mysqlCfg
	return mysqlCfg, err
}

func CreateMySQLContainer() (*config.DBConf, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL"),
	}
	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := mysqlC.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := mysqlC.MappedPort(ctx, "3306")
	if err != nil {
		return nil, err
	}

	dbCfg := &config.DBConf{
		Host:       ip,
		Port:       port.Port(),
		Database:   "chainmaker_explorer_dev",
		Username:   "root",
		Password:   "password",
		Prefix:     "test_",
		DbProvider: "Mysql",
	}

	//time.Sleep(time.Second * 10)
	return dbCfg, nil
}

func InitClickHouseContainer() (*config.DBConf, error) {
	var err error
	onceClickHouse.Do(func() {
		//初始化配置
		_ = config.InitConfig("", "")
		clickHouseCfg, err = CreateClickHouseContainer()
		if clickHouseCfg == nil {
			return
		}
		config.GlobalConfig.DBConf = clickHouseCfg
		InitDbConn(clickHouseCfg)
	})
	config.GlobalConfig.DBConf = clickHouseCfg
	return clickHouseCfg, err
}

func CreateClickHouseContainer() (*config.DBConf, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "yandex/clickhouse-server:latest",
		ExposedPorts: []string{"8123/tcp"},
		WaitingFor:   wait.ForLog("Listening for connections with native protocol").WithStartupTimeout(2 * time.Minute),
	}
	clickHouseC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := clickHouseC.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := clickHouseC.MappedPort(ctx, "8123")
	if err != nil {
		return nil, err
	}

	dbCfg := &config.DBConf{
		Host:       ip,
		Port:       port.Port(),
		Database:   "chainmaker_explorer_dev",
		Username:   "root",
		Password:   "password",
		Prefix:     "test",
		DbProvider: "ClickHouse",
	}

	//time.Sleep(time.Second * 10)
	return dbCfg, nil
}

func InitRedisContainer() (*config.RedisConfig, error) {
	var err error
	onceRedis.Do(func() {
		//初始化配置
		_ = config.InitConfig("", "")
		redisCfg, err = CreateRedisContainer()
	})
	config.GlobalConfig.RedisDB = redisCfg
	log.Infof("=========redisCfg:%v=", redisCfg)
	return redisCfg, err
}

func CreateRedisContainer() (*config.RedisConfig, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		//Env: map[string]string{
		//	"REDIS_PASSWORD": "password11",
		//},
		WaitingFor: wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := redisC.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		return nil, err
	}

	redisCfg = &config.RedisConfig{
		Type:     "node",
		Host:     []string{ip + ":" + port.Port()},
		Password: "",
		Username: "",
	}
	//time.Sleep(time.Second * 10)
	return redisCfg, nil
}

func getChainListConfigTest(fileName string) []*config.ChainInfo {
	resultValue := make([]*config.ChainInfo, 0)
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}
	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func geFileData(fileName string) (*json.Decoder, error) {
	file, err := os.Open("../testData/" + fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}

	// 解码 JSON 文件内容到 blockInfo 结构体
	decoder := json.NewDecoder(file)
	return decoder, err
}
