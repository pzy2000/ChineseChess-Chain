package chain

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
)

func insertSubscribeTest() (*db.Subscribe, error) {
	subscribeInfo := &db.Subscribe{
		ChainId:  "chain1",
		UserKey:  "1234",
		UserCert: "1234",
	}

	err := dbhandle.InsertSubscribe(subscribeInfo)
	return subscribeInfo, err
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

func TestMain(m *testing.M) {
	// 在运行测试之前，更改当前工作目录
	config.SubscribeChains = getChainListConfigTest("0_chainListConfig.json")
	// 初始化数据库配置
	dbCfg, err := db.InitMySQLContainer()
	if err != nil || dbCfg.Host == "" {
		return
	}
	// 运行其他测试
	os.Exit(m.Run())
}

func TestGetConfigShow(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Test case 1:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = GetConfigShow()
		})
	}
}

func TestGetSubscribeChains(t *testing.T) {
	tests := []struct {
		name    string
		want    []*config.ChainInfo
		wantErr bool
	}{
		{
			name:    "Test case 1:",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetSubscribeChains()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubscribeChains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInitChainConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test case 1:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitChainConfig()
		})
	}
}

func Test_mergeChainInfo(t *testing.T) {
	type args struct {
		configChains []*config.ChainInfo
		dbChains     []*config.ChainInfo
	}
	tests := []struct {
		name string
		args args
		want []*config.ChainInfo
	}{
		{
			name: "Test case 1:",
			args: args{
				configChains: []*config.ChainInfo{
					{
						ChainId:  "chain3",
						AuthType: "123",
						OrgId:    "123",
						HashType: "123",
					},
					{
						ChainId:  "chain1",
						AuthType: "1234",
						OrgId:    "1234",
						HashType: "1234",
					},
				},
				dbChains: []*config.ChainInfo{
					{
						ChainId:  "chain1",
						AuthType: "123",
						OrgId:    "123",
						HashType: "123",
					},
					{
						ChainId:  "chain2",
						AuthType: "1234",
						OrgId:    "1234",
						HashType: "1234",
					},
				},
			},
			want: []*config.ChainInfo{
				{
					ChainId:  "chain3",
					AuthType: "123",
					OrgId:    "123",
					HashType: "123",
				},
				{
					ChainId:  "chain1",
					AuthType: "123",
					OrgId:    "123",
					HashType: "123",
				},
				{
					ChainId:  "chain2",
					AuthType: "1234",
					OrgId:    "1234",
					HashType: "1234",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeChainInfo(tt.args.configChains, tt.args.dbChains)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("mergeChainInfo() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetIsMainChain(t *testing.T) {
	config.GlobalConfig.ChainConf.IsMainChain = true
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "test: case 1",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIsMainChain(); got != tt.want {
				t.Errorf("GetIsMainChain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSubscribeByChainId(t *testing.T) {
	subscribeInfo, err := insertSubscribeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    *config.ChainInfo
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: "chain1",
			},
			want: &config.ChainInfo{
				ChainId:  subscribeInfo.ChainId,
				AuthType: subscribeInfo.AuthType,
				OrgId:    subscribeInfo.OrgId,
				UserInfo: &config.UserInfo{
					UserKey:  subscribeInfo.UserKey,
					UserCert: subscribeInfo.UserCert,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSubscribeByChainId(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubscribeByChainId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Subscribe{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("DealContractEvents() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetSubscribeChains1(t *testing.T) {
	_, err := insertSubscribeTest()
	if err != nil {
		return
	}

	tests := []struct {
		name    string
		want    []*config.ChainInfo
		wantErr bool
	}{
		{
			name: "test: case 1",
			want: []*config.ChainInfo{
				{
					ChainId: "chain1",
					UserInfo: &config.UserInfo{
						UserKey:  "1234",
						UserCert: "1234",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSubscribeChains()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubscribeChains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetSubscribeChains() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_mergeChainInfo1(t *testing.T) {
	type args struct {
		configChains []*config.ChainInfo
		dbChains     []*config.ChainInfo
	}
	tests := []struct {
		name string
		args args
		want []*config.ChainInfo
	}{
		{
			name: "Test case 1:",
			args: args{
				configChains: []*config.ChainInfo{
					{
						ChainId:  "chain3",
						AuthType: "123",
						OrgId:    "123",
						HashType: "123",
					},
					{
						ChainId:  "chain1",
						AuthType: "1234",
						OrgId:    "1234",
						HashType: "1234",
					},
				},
				dbChains: []*config.ChainInfo{
					{
						ChainId:  "chain1",
						AuthType: "123",
						OrgId:    "123",
						HashType: "123",
					},
					{
						ChainId:  "chain2",
						AuthType: "1234",
						OrgId:    "1234",
						HashType: "1234",
					},
				},
			},
			want: []*config.ChainInfo{
				{
					ChainId:  "chain3",
					AuthType: "123",
					OrgId:    "123",
					HashType: "123",
				},
				{
					ChainId:  "chain1",
					AuthType: "123",
					OrgId:    "123",
					HashType: "123",
				},
				{
					ChainId:  "chain2",
					AuthType: "1234",
					OrgId:    "1234",
					HashType: "1234",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeChainInfo(tt.args.configChains, tt.args.dbChains)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("mergeChainInfo() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetConfigShow2(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Test case 1",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = GetConfigShow()
		})
	}
}

func TestGetSubscribeByChainId1(t *testing.T) {
	subscribeInfo, err := insertSubscribeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    *config.ChainInfo
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				chainId: subscribeInfo.ChainId,
			},
			want: &config.ChainInfo{
				ChainId:  subscribeInfo.ChainId,
				AuthType: subscribeInfo.AuthType,
				UserInfo: &config.UserInfo{
					UserKey:  subscribeInfo.UserKey,
					UserCert: subscribeInfo.UserCert,
				},
			},
			wantErr: false,
		},
		{
			name: "Test case 2",
			args: args{
				chainId: subscribeInfo.ChainId,
			},
			want: &config.ChainInfo{
				ChainId:  subscribeInfo.ChainId,
				AuthType: subscribeInfo.AuthType,
				UserInfo: &config.UserInfo{
					UserKey:  subscribeInfo.UserKey,
					UserCert: subscribeInfo.UserCert,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSubscribeByChainId(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubscribeByChainId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetSubscribeByChainId() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetSubscribeChains2(t *testing.T) {
	subscribeInfo, err := insertSubscribeTest()
	if err != nil {
		return
	}

	tests := []struct {
		name    string
		want    []*config.ChainInfo
		wantErr bool
	}{
		{
			name: "Test GetSubscribeChains",
			want: []*config.ChainInfo{
				{
					ChainId:  subscribeInfo.ChainId,
					AuthType: subscribeInfo.AuthType,
					UserInfo: &config.UserInfo{
						UserKey:  subscribeInfo.UserKey,
						UserCert: subscribeInfo.UserCert,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSubscribeChains()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubscribeChains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetSubscribeChains() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestInitChainConfig1(t *testing.T) {
	_, err := insertSubscribeTest()
	if err != nil {
		return
	}

	tests := []struct {
		name string
		want []*config.ChainInfo
	}{
		{
			name: "test: case 1",
			want: []*config.ChainInfo{
				{
					ChainId: "chain1",
					UserInfo: &config.UserInfo{
						UserKey:  "1234",
						UserCert: "1234",
					},
				},
			},
		},
		{
			name: "test: case 2",
			want: []*config.ChainInfo{
				{
					ChainId: "chain1",
					UserInfo: &config.UserInfo{
						UserKey:  "1234",
						UserCert: "1234",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitChainConfig()
			got := config.SubscribeChains
			if !cmp.Equal(got, tt.want) {
				t.Errorf("InitChainConfig() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_mergeChainInfo2(t *testing.T) {
	type args struct {
		configChains []*config.ChainInfo
		dbChains     []*config.ChainInfo
	}
	tests := []struct {
		name string
		args args
		want []*config.ChainInfo
	}{
		{
			name: "Test mergeChainInfo",
			args: args{
				configChains: []*config.ChainInfo{
					{ChainId: "chain1", AuthType: "auth1"},
					{ChainId: "chain2", AuthType: "auth2"},
				},
				dbChains: []*config.ChainInfo{
					{ChainId: "chain2", AuthType: "auth2"},
					{ChainId: "chain3", AuthType: "auth3"},
				},
			},
			want: []*config.ChainInfo{
				{ChainId: "chain1", AuthType: "auth1"},
				{ChainId: "chain2", AuthType: "auth2"},
				{ChainId: "chain3", AuthType: "auth3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeChainInfo(tt.args.configChains, tt.args.dbChains)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeChainInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
