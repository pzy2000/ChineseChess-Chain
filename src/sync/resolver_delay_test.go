package sync

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/saveTasks"
	"reflect"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	ChainId1 = "chain1"
)

func TestBatchDelayedUpdate(t *testing.T) {
	type args struct {
		chainId      string
		blockHeights []int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试批量延迟更新",
			args: args{
				chainId:      ChainId1,
				blockHeights: []int64{1, 2, 3},
			},
			wantErr: false,
		},
		// 添加其他测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BatchDelayedUpdate(tt.args.chainId, tt.args.blockHeights); (err != nil) != tt.wantErr {
				t.Errorf("BatchDelayedUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// 定义一个排序函数，按照 Address 排序
func sortAccounts(accounts []*db.Account) {
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].Address < accounts[j].Address
	})
}

func TestBuildDelayedUpdateData(t *testing.T) {
	buildDelayedUpdateData := &DelayedUpdateData{
		InsertSubChainCross: make([]*db.CrossSubChainCrossChain, 0),
		UpdateSubChainCross: make([]*db.CrossSubChainCrossChain, 0),
		UpdateSubChainData:  make([]*db.CrossSubChainData, 0),
		InsertGasList:       make([]*db.Gas, 0),
		UpdateGasList:       make([]*db.Gas, 0),
		UpdateTxBlack: &db.UpdateTxBlack{
			AddTxBlack:    make([]*db.BlackTransaction, 0),
			DeleteTxBlack: make([]*db.Transaction, 0),
		},
		ContractResult: &db.GetContractResult{
			UpdateContractTxEventNum: make([]*db.Contract, 0),
			IdentityContract:         make([]*db.IdentityContract, 0),
			UpdateFungibleContract:   make([]*db.FungibleContract, 0),
			UpdateNonFungible:        make([]*db.NonFungibleContract, 0),
		},
		FungibleTransfer:    make([]*db.FungibleTransfer, 0),
		NonFungibleTransfer: make([]*db.NonFungibleTransfer, 0),
		BlockPosition:       &db.BlockPosition{},
		UpdateAccountResult: &db.UpdateAccountResult{},
		TokenResult: &db.TokenResult{
			InsertUpdateToken: make([]*db.NonFungibleToken, 0),
			DeleteToken:       make([]*db.NonFungibleToken, 0),
		},
		ContractMap: make(map[string]*db.Contract),
	}

	type args struct {
		chainId            string
		blockHeights       []int64
		delayedUpdateCache *GetRealtimeCacheData
	}
	tests := []struct {
		name    string
		args    args
		want    *DelayedUpdateData
		wantErr bool
	}{
		{
			name: "测试构建延迟更新数据",
			args: args{
				chainId:            ChainId1,
				blockHeights:       []int64{1, 2, 3},
				delayedUpdateCache: &GetRealtimeCacheData{},
			},
			want:    buildDelayedUpdateData,
			wantErr: false,
		},
		// 添加其他测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildDelayedUpdateData(tt.args.chainId, tt.args.blockHeights, tt.args.delayedUpdateCache)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildDelayedUpdateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestExtractTxIdsAndContractNames(t *testing.T) {
	type args struct {
		txInfoList []*db.Transaction
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 map[string]string
		want2 map[string]*db.Transaction
	}{
		{
			name: "测试提取交易ID和合约名称",
			args: args{
				txInfoList: []*db.Transaction{
					{
						TxId:        "1234",
						Sender:      "1234",
						BlockHash:   "1234",
						SenderOrgId: "1234",
						BlockHeight: 12,
					},
				},
			},
			want: []string{
				"1234",
			}, // 填充期望的结果
			want1: map[string]string{}, // 填充期望的结果
			want2: map[string]*db.Transaction{
				"1234": {
					TxId:        "1234",
					Sender:      "1234",
					BlockHash:   "1234",
					SenderOrgId: "1234",
					BlockHeight: 12,
				},
			},
		},
		// 添加其他测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := ExtractTxIdsAndContractNames(tt.args.txInfoList)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractTxIdsAndContractNames() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ExtractTxIdsAndContractNames() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("ExtractTxIdsAndContractNames() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestGetDelayedUpdateByDB(t *testing.T) {
	type args struct {
		chainId           string
		heightDB          []int64
		delayedUpdateData *GetRealtimeCacheData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试从数据库获取延迟更新",
			args: args{
				chainId:           ChainId1,
				heightDB:          []int64{1, 2, 3},
				delayedUpdateData: &GetRealtimeCacheData{
					// 填充测试数据
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetDelayedUpdateByDB(tt.args.chainId, tt.args.heightDB, tt.args.delayedUpdateData); (err != nil) != tt.wantErr {
				t.Errorf("GetDelayedUpdateByDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRealtimeDataCache(t *testing.T) {
	delayedUpdateData := &GetRealtimeCacheData{
		TxList:         make(map[string]*db.Transaction),
		ContractAddrs:  make(map[string]string, 0),
		UserInfoMap:    make(map[string]*db.User, 0),
		GasRecords:     make([]*db.GasRecord, 0),
		ContractEvents: make([]*db.ContractEvent, 0),
		CrossTransfers: make([]*db.CrossTransactionTransfer, 0),
	}
	type args struct {
		chainId      string
		blockHeights []int64
	}
	tests := []struct {
		name    string
		args    args
		want    *GetRealtimeCacheData
		wantErr bool
	}{
		{
			name: "测试获取实时数据缓存",
			args: args{
				chainId:      ChainId1,
				blockHeights: []int64{1, 2, 3},
			},
			want:    delayedUpdateData,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRealtimeDataCache(tt.args.chainId, tt.args.blockHeights)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRealtimeDataCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetRealtimeDataCache() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestParallelParseUpdateDataToDB(t *testing.T) {
	buildDelayedUpdateData := &DelayedUpdateData{
		InsertSubChainCross: make([]*db.CrossSubChainCrossChain, 0),
		UpdateSubChainCross: make([]*db.CrossSubChainCrossChain, 0),
		UpdateSubChainData:  make([]*db.CrossSubChainData, 0),
		InsertGasList:       make([]*db.Gas, 0),
		UpdateGasList:       make([]*db.Gas, 0),
		UpdateTxBlack: &db.UpdateTxBlack{
			AddTxBlack:    make([]*db.BlackTransaction, 0),
			DeleteTxBlack: make([]*db.Transaction, 0),
		},
		ContractResult: &db.GetContractResult{
			UpdateContractTxEventNum: make([]*db.Contract, 0),
			IdentityContract:         make([]*db.IdentityContract, 0),
		},
		FungibleTransfer:    make([]*db.FungibleTransfer, 0),
		NonFungibleTransfer: make([]*db.NonFungibleTransfer, 0),
		BlockPosition:       &db.BlockPosition{},
		UpdateAccountResult: &db.UpdateAccountResult{},
		TokenResult: &db.TokenResult{
			InsertUpdateToken: make([]*db.NonFungibleToken, 0),
			DeleteToken:       make([]*db.NonFungibleToken, 0),
		},
		ContractMap: make(map[string]*db.Contract),
	}

	type args struct {
		chainId           string
		delayedUpdateData *DelayedUpdateData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试并行解析更新数据到数据库",
			args: args{
				chainId:           ChainId1,
				delayedUpdateData: buildDelayedUpdateData,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParallelParseUpdateDataToDB(tt.args.chainId, tt.args.delayedUpdateData); (err != nil) != tt.wantErr {
				t.Errorf("ParallelParseUpdateDataToDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createTasksDelayedUpdate(t *testing.T) {
	buildDelayedUpdateData := &DelayedUpdateData{
		InsertSubChainCross: make([]*db.CrossSubChainCrossChain, 0),
		UpdateSubChainCross: make([]*db.CrossSubChainCrossChain, 0),
		UpdateSubChainData:  make([]*db.CrossSubChainData, 0),
		InsertGasList:       make([]*db.Gas, 0),
		UpdateGasList:       make([]*db.Gas, 0),
		UpdateTxBlack: &db.UpdateTxBlack{
			AddTxBlack:    make([]*db.BlackTransaction, 0),
			DeleteTxBlack: make([]*db.Transaction, 0),
		},
		ContractResult: &db.GetContractResult{
			UpdateContractTxEventNum: make([]*db.Contract, 0),
			IdentityContract:         make([]*db.IdentityContract, 0),
		},
		FungibleTransfer:    make([]*db.FungibleTransfer, 0),
		NonFungibleTransfer: make([]*db.NonFungibleTransfer, 0),
		BlockPosition:       &db.BlockPosition{},
		UpdateAccountResult: &db.UpdateAccountResult{},
		TokenResult: &db.TokenResult{
			InsertUpdateToken: make([]*db.NonFungibleToken, 0),
			DeleteToken:       make([]*db.NonFungibleToken, 0),
		},
		ContractMap: make(map[string]*db.Contract),
	}

	type args struct {
		chainId       string
		delayedUpdate *DelayedUpdateData
	}
	tests := []struct {
		name string
		args args
		want []saveTasks.Task
	}{
		{
			name: "测试创建延迟更新任务",
			args: args{
				chainId:       ChainId1,
				delayedUpdate: buildDelayedUpdateData,
			},
			want: []saveTasks.Task{
				// 填充期望的结果
			},
		},
		// 添加其他测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createTasksDelayedUpdate(tt.args.chainId, tt.args.delayedUpdate)
			if len(got) == 0 {
				t.Errorf("createTasksDelayedUpdate() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}
