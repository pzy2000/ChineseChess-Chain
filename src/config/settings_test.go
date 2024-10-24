package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/viper"
)

func TestConfig_printLog(t *testing.T) {
	type fields struct {
		WebConf         *WebConf
		SubscribeConfig *SubscribeConfig
		ChainsConfig    []*ChainConfig
		DBConf          *DBConf
		LogConf         *LogConf
		ChainConf       *ChainConf
		AlarmerConf     *AlarmerConf
		MonitorConf     *MonitorConf
		SensitiveConf   *SensitiveConf
		RedisDB         *RedisConfig
	}
	type args struct {
		env string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test case 1:",
			fields: fields{
				WebConf: nil,
				DBConf: &DBConf{
					Host: "127.0.0.1",
					Port: "8080",
				},
			},
			args: args{
				env: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				WebConf:         tt.fields.WebConf,
				SubscribeConfig: tt.fields.SubscribeConfig,
				ChainsConfig:    tt.fields.ChainsConfig,
				DBConf:          tt.fields.DBConf,
				LogConf:         tt.fields.LogConf,
				ChainConf:       tt.fields.ChainConf,
				AlarmerConf:     tt.fields.AlarmerConf,
				MonitorConf:     tt.fields.MonitorConf,
				SensitiveConf:   tt.fields.SensitiveConf,
				RedisDB:         tt.fields.RedisDB,
			}
			c.printLog(tt.args.env)
		})
	}
}

func TestDBConf_ToClickHouseUrl(t *testing.T) {
	type fields struct {
		Host       string
		Port       string
		Database   string
		Username   string
		Password   string
		Prefix     string
		DbProvider string
	}
	type args struct {
		useDataBase bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test case 1:",
			fields: fields{
				Host:     "127.0.0.1",
				Port:     "8080",
				Database: "test1",
				Username: "default",
				Password: "123456",
			},
			args: args{
				useDataBase: true,
			},
			want: "clickhouse://default:123456@127.0.0.1:8080/test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbConfig := &DBConf{
				Host:       tt.fields.Host,
				Port:       tt.fields.Port,
				Database:   tt.fields.Database,
				Username:   tt.fields.Username,
				Password:   tt.fields.Password,
				Prefix:     tt.fields.Prefix,
				DbProvider: tt.fields.DbProvider,
			}
			if got := dbConfig.ToClickHouseUrl(tt.args.useDataBase); got != tt.want {
				t.Errorf("ToClickHouseUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBConf_ToMysqlUrl(t *testing.T) {
	type fields struct {
		Host       string
		Port       string
		Database   string
		Username   string
		Password   string
		Prefix     string
		DbProvider string
	}
	type args struct {
		useDataBase bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test case 1:",
			fields: fields{
				Host:     "127.0.0.1",
				Port:     "8080",
				Database: "test1",
				Username: "default",
				Password: "123456",
			},
			args: args{
				useDataBase: true,
			},
			want: "default:123456@tcp(127.0.0.1:8080)/test1?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbConfig := &DBConf{
				Host:       tt.fields.Host,
				Port:       tt.fields.Port,
				Database:   tt.fields.Database,
				Username:   tt.fields.Username,
				Password:   tt.fields.Password,
				Prefix:     tt.fields.Prefix,
				DbProvider: tt.fields.DbProvider,
			}
			if got := dbConfig.ToMysqlUrl(tt.args.useDataBase); got != tt.want {
				t.Errorf("ToMysqlUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChainsInfoAll(t *testing.T) {
	tests := []struct {
		name string
		want []*ChainInfo
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChainsInfoAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChainsInfoAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadAbsPathFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test case 1:",
			args: args{
				path: "../../configs/config.yml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadAbsPathFile(tt.args.path); got == "" {
				t.Errorf("ReadAbsPathFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebConf_ToUrl(t *testing.T) {
	type fields struct {
		Address             string
		Port                int
		CrossDomain         bool
		ThirdApplyUrl       string
		RelayCrossChainUrl  string
		TestnetUrl          string
		OpennetUrl          string
		MonitorPort         int
		ManageBackendApiKey string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test case 1:",
			fields: fields{
				Address: "0.0.0.0",
				Port:    79999,
			},
			want: "0.0.0.0:79999",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webConfig := &WebConf{
				Address:             tt.fields.Address,
				Port:                tt.fields.Port,
				CrossDomain:         tt.fields.CrossDomain,
				ThirdApplyUrl:       tt.fields.ThirdApplyUrl,
				RelayCrossChainUrl:  tt.fields.RelayCrossChainUrl,
				TestnetUrl:          tt.fields.TestnetUrl,
				OpennetUrl:          tt.fields.OpennetUrl,
				MonitorPort:         tt.fields.MonitorPort,
				ManageBackendApiKey: tt.fields.ManageBackendApiKey,
			}
			if got := webConfig.ToUrl(); got != tt.want {
				t.Errorf("ToUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildChainInfo(t *testing.T) {
	// 创建输入参数：ChainConfig 结构体实例
	sampleChainConfig := &ChainConfig{
		ChainId:  "test-chain-id",
		AuthType: "test-auth-type",
		OrgId:    "test-org-id",
		HashType: "test-hash-type",
		NodesConfig: []*NodeConf{
			{
				Remotes: "test-remote",
				CaPaths: "test-ca-path",
				TlsHost: "test-tls-host",
				Tls:     true,
			},
		},
		UserConf: &CertConf{
			PrivKeyFile: "test-privkey-file",
			CertFile:    "test-cert-file",
		},
	}

	// 创建期望的输出结果：ChainInfo 结构体实例
	expectedChainInfo := &ChainInfo{
		ChainId:  "test-chain-id",
		AuthType: "test-auth-type",
		OrgId:    "test-org-id",
		HashType: "test-hash-type",
		NodesList: []*NodeInfo{
			{
				Addr:        "test-remote",
				OrgCA:       "", // 假设 ReadAbsPathFile 返回空字符串
				TLSHostName: "test-tls-host",
				Tls:         true,
			},
		},
		UserInfo: &UserInfo{
			UserKey:  "", // 假设 ReadAbsPathFile 返回空字符串
			UserCert: "", // 假设 ReadAbsPathFile 返回空字符串
		},
	}

	type args struct {
		chain *ChainConfig
	}
	tests := []struct {
		name string
		args args
		want *ChainInfo
	}{
		{
			name: "Test_buildChainInfo",
			args: args{
				chain: sampleChainConfig,
			},
			want: expectedChainInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildChainInfo(tt.args.chain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildChainInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getERC20AbiJson(t *testing.T) {
	// 创建一个临时文件
	tmpFile, err := ioutil.TempFile("", "test-erc20-abi-*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // 在测试结束时删除临时文件

	// 写入一个有效的 ERC20 ABI JSON 到临时文件
	_, err = tmpFile.WriteString(`[{"type": "function", "name": "balanceOf", "inputs": [{"name": "owner", "type": "address"}], "outputs": [{"name": "balance", "type": "uint256"}]}]`)
	if err != nil {
		t.Fatalf("Failed to write ERC20 ABI JSON to temporary file: %v", err)
	}

	// 设置 GlobalConfig.SubscribeConfig.ERC20AbiFile 为临时文件的路径
	GlobalConfig = &Config{
		SubscribeConfig: &SubscribeConfig{
			ERC20AbiFile: tmpFile.Name(),
		},
	}

	tests := []struct {
		name string
		want *abi.ABI
	}{
		{
			name: "Test_getERC20AbiJson",
			want: nil, // 将在下面设置期望的 abi.ABI 结构体实例
		},
	}

	// 创建一个期望的 abi.ABI 结构体实例
	expectedAbi, err := abi.JSON(strings.NewReader(`[{"type": "function", "name": "balanceOf", "inputs": [{"name": "owner", "type": "address"}], "outputs": [{"name": "balance", "type": "uint256"}]}]`))
	if err != nil {
		t.Fatalf("Failed to create expected ABI instance: %v", err)
	}
	tests[0].want = &expectedAbi

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getERC20AbiJson(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getERC20AbiJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getERC721AbiJson(t *testing.T) {
	// 创建一个临时文件
	tmpFile, err := ioutil.TempFile("", "test-erc20-abi-*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // 在测试结束时删除临时文件

	// 写入一个有效的 ERC20 ABI JSON 到临时文件
	_, err = tmpFile.WriteString(`[{"type": "function", "name": "balanceOf", "inputs": [{"name": "owner", "type": "address"}], "outputs": [{"name": "balance", "type": "uint256"}]}]`)
	if err != nil {
		t.Fatalf("Failed to write ERC20 ABI JSON to temporary file: %v", err)
	}

	// 设置 GlobalConfig.SubscribeConfig.ERC20AbiFile 为临时文件的路径
	GlobalConfig = &Config{
		SubscribeConfig: &SubscribeConfig{
			ERC721AbiFile: tmpFile.Name(),
		},
	}

	tests := []struct {
		name string
		want *abi.ABI
	}{
		{
			name: "getERC721AbiJson",
			want: nil, // 将在下面设置期望的 abi.ABI 结构体实例
		},
	}

	// 创建一个期望的 abi.ABI 结构体实例
	expectedAbi, err := abi.JSON(strings.NewReader(`[{"type": "function", "name": "balanceOf", "inputs": [{"name": "owner", "type": "address"}], "outputs": [{"name": "balance", "type": "uint256"}]}]`))
	if err != nil {
		t.Fatalf("Failed to create expected ABI instance: %v", err)
	}
	tests[0].want = &expectedAbi

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getERC721AbiJson(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getERC721AbiJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfigDirPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Test case 1:",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = GetConfigDirPath()
		})
	}
}

func Test_initCMViper1(t *testing.T) {
	type args struct {
		env       string
		gConfPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *viper.Viper
		wantErr bool
	}{
		{
			name: "Test case 1:",
			args: args{
				gConfPath: "../../configs",
				env:       "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := initCMViper(tt.args.gConfPath, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("initCMViper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("InitConfig() = %v", got)
			}
		})
	}
}

func TestGetConfigFilePath(t *testing.T) {
	type args struct {
		gConfPath string
		env       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test case 1:",
			args: args{
				gConfPath: "../../configs",
				env:       "",
			},
			want: "../../configs/config.yml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfigFilePath(tt.args.gConfPath, tt.args.env); got != tt.want {
				t.Errorf("GetConfigFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChainsInfoAll1(t *testing.T) {
	GlobalConfig = &Config{
		ChainsConfig: []*ChainConfig{
			{
				ChainId:  "123",
				AuthType: "123",
				OrgId:    "123",
				HashType: "123",
				NodesConfig: []*NodeConf{
					{
						TlsHost: "456",
					},
				},
				UserConf: &CertConf{
					PrivKeyFile: "",
					CertFile:    "",
				},
			},
		},
	}

	want := []*ChainInfo{
		{
			ChainId:  "123",
			AuthType: "123",
			OrgId:    "123",
			HashType: "123",
			NodesList: []*NodeInfo{
				{
					TLSHostName: "456",
				},
			},
			UserInfo: &UserInfo{},
		},
	}
	tests := []struct {
		name string
		want []*ChainInfo
	}{
		{
			name: "Test case 1:",
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetChainsInfoAll()
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetChainsInfoAll() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestReadAbsPathFile1(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test case 1:",
			args: args{
				path: "../../configs/config.yml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadAbsPathFile(tt.args.path); got == "" {
				t.Errorf("ReadAbsPathFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBConf_ToClickHouseUrl1(t *testing.T) {
	type fields struct {
		Host       string
		Port       string
		Database   string
		Username   string
		Password   string
		Prefix     string
		DbProvider string
	}
	type args struct {
		useDataBase bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbConfig := &DBConf{
				Host:       tt.fields.Host,
				Port:       tt.fields.Port,
				Database:   tt.fields.Database,
				Username:   tt.fields.Username,
				Password:   tt.fields.Password,
				Prefix:     tt.fields.Prefix,
				DbProvider: tt.fields.DbProvider,
			}
			if got := dbConfig.ToClickHouseUrl(tt.args.useDataBase); got != tt.want {
				t.Errorf("ToClickHouseUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	c := Config{
		DBConf: &DBConf{
			Host:     "127.0.0.1",
			Port:     "3306",
			Database: "por",
			Username: "root",
			Password: "chainmaker",
		},
	}
	fmt.Println(c.DBConf.ToMysqlUrl(true))
}

func Test_initCMViper(t *testing.T) {
	type args struct {
		env       string
		gConfPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *viper.Viper
		wantErr bool
	}{
		{
			name: "Test case 1:",
			args: args{
				gConfPath: "../../configs",
				env:       "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := initCMViper(tt.args.gConfPath, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("initCMViper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("InitConfig() = %v", got)
			}
		})
	}
}

func TestDBConf_ToPgsqlUrl(t *testing.T) {
	type fields struct {
		Host        string
		Port        string
		Database    string
		Username    string
		Password    string
		Prefix      string
		DbProvider  string
		MaxByteSize int
		MaxPoolSize int
	}
	type args struct {
		useDataBase bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test case 1:",
			fields: fields{
				Host:     "127.0.0.1",
				Port:     "8090",
				Database: "string",
				Username: "Username",
				Password: "Password",
			},
			args: args{
				useDataBase: true,
			},
			want: "host=127.0.0.1 port=8090 user=Username password=Password dbname=string sslmode=disable client_encoding=UTF8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbConfig := &DBConf{
				Host:        tt.fields.Host,
				Port:        tt.fields.Port,
				Database:    tt.fields.Database,
				Username:    tt.fields.Username,
				Password:    tt.fields.Password,
				Prefix:      tt.fields.Prefix,
				DbProvider:  tt.fields.DbProvider,
				MaxByteSize: tt.fields.MaxByteSize,
				MaxPoolSize: tt.fields.MaxPoolSize,
			}
			if got := dbConfig.ToPgsqlUrl(tt.args.useDataBase); got != tt.want {
				t.Errorf("DBConf.ToPgsqlUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitConfig(t *testing.T) {
	type args struct {
		confPath string
		env      string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "Test case 1:",
			args: args{
				confPath: "",
				env:      "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = InitConfig(tt.args.confPath, tt.args.env)
		})
	}
}
