package dbhandle

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"fmt"
	"testing"

	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	ChainId1 = "chain1"
	ChainId2 = "chain2"
	ChainId3 = "chain3"
)

func insertChainTest1() (*db.Chain, error) {
	chainInfo := &db.Chain{
		ChainId:   ChainId1,
		Version:   "12345",
		ChainName: ChainId1,
		Timestamp: 12345,
	}
	err := InsertUpdateChainInfo(chainInfo, db.SubscribeOK)
	return chainInfo, err
}

func insertChainTest2() (*db.Chain, error) {
	chainInfo := &db.Chain{
		ChainId:   ChainId2,
		Version:   "12345",
		ChainName: ChainId2,
		Timestamp: 123456,
	}
	err := InsertUpdateChainInfo(chainInfo, db.SubscribeOK)
	return chainInfo, err
}

func TestDeleteChain(t *testing.T) {
	_, err := insertChainTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteChain(tt.args.chainId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteChain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetChainInfoById(t *testing.T) {
	chainInfo, err := insertChainTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Chain
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    chainInfo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetChainInfoById(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChainInfoById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Chain{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetBlockByHash() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")))
			}
		})
	}
}

func TestGetChainInfoCache(t *testing.T) {
	chainInfo, err := insertChainTest1()
	if err != nil {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Chain
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    chainInfo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = GetChainInfoById(tt.args.chainId)
			redisKey := fmt.Sprintf(cache.RedisDbChainConfig, prefix, tt.args.chainId)
			got, err := GetChainInfoCache(redisKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChainInfoCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Chain{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetChainInfoCache() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")))
			}
		})
	}
}

func TestGetChainListByPage(t *testing.T) {
	chainInfo, err := insertChainTest1()
	if err != nil {
		return
	}

	type args struct {
		offset  int
		limit   int
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Chain
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				offset:  0,
				limit:   10,
				chainId: ChainId1,
			},
			want: []*db.Chain{
				chainInfo,
			},
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetChainListByPage(tt.args.offset, tt.args.limit, tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChainListByPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Chain{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetChainListByPage() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")))
			}
			if got1 != tt.want1 {
				t.Errorf("GetChainListByPage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInsertUpdateChainInfo(t *testing.T) {
	chainInfo := &db.Chain{
		ChainId:   ChainId1,
		Version:   "12345",
		ChainName: ChainId1,
		Timestamp: 12345,
	}

	type args struct {
		chainData       *db.Chain
		subscribeStatus int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainData: chainInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertUpdateChainInfo(tt.args.chainData, tt.args.subscribeStatus); (err != nil) != tt.wantErr {
				t.Errorf("InsertUpdateChainInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateChainInfo(t *testing.T) {
	chainInfo, err := insertChainTest2()
	if err != nil {
		return
	}

	type args struct {
		chainInfo *db.Chain
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainInfo: chainInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateChainInfo(tt.args.chainInfo); (err != nil) != tt.wantErr {
				t.Errorf("UpdateChainInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateChainInfoByConfig(t *testing.T) {
	type args struct {
		chainId     string
		chainConfig *pbConfig.ChainConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainId1,
				chainConfig: &pbConfig.ChainConfig{
					Version: "1.1",
					Block: &pbConfig.BlockConfig{
						BlockInterval:     123,
						BlockSize:         123,
						BlockTxCapacity:   123,
						TxTimeout:         123,
						TxTimestampVerify: false,
					},
					AccountConfig: &pbConfig.GasAccountConfig{
						EnableGas: true,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateChainInfoByConfig(tt.args.chainId, tt.args.chainConfig); (err != nil) != tt.wantErr {
				t.Errorf("UpdateChainInfoByConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
