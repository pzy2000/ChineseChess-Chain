package sync

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

func TestGetCrossCycleTxDataCache(t *testing.T) {
	crossCycleTxMap := make(map[string]*db.CrossCycleTransaction, 0)
	crossCycleTxMap["123"] = &db.CrossCycleTransaction{
		CrossId: "123",
	}
	SetCrossCycleTxDataCache(ChainId1, 12, crossCycleTxMap)
	type args struct {
		chainId     string
		blockHeight int64
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.CrossCycleTransaction
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:     ChainId1,
				blockHeight: 12,
			},
			want: crossCycleTxMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossCycleTxDataCache(tt.args.chainId, tt.args.blockHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossCycleTxDataCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossCycleTxDataCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossTransfersCache(t *testing.T) {
	crossTransfers := []*db.CrossTransactionTransfer{
		{
			CrossId: "123",
		},
	}
	SetCrossTransfersCache(ChainId1, 12, crossTransfers)
	type args struct {
		chainId     string
		blockHeight int64
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossTransactionTransfer
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:     ChainId1,
				blockHeight: 12,
			},
			want: crossTransfers,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossTransfersCache(tt.args.chainId, tt.args.blockHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossTransfersCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossTransfersCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRealtimeDataHandle(t *testing.T) {
	blockInfo := getBlockInfoTest("1_blockInfoJsonContract.json")
	type args struct {
		blockInfo *common.BlockInfo
		hashType  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				blockInfo: blockInfo,
				hashType:  "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := RealtimeDataHandle(tt.args.blockInfo, tt.args.hashType)
			if (err != nil) != tt.wantErr {
				t.Errorf("RealtimeDataHandle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSetCrossCycleTxDataCache(t *testing.T) {
	crossCycleTxMap := make(map[string]*db.CrossCycleTransaction, 0)
	crossCycleTxMap["123"] = &db.CrossCycleTransaction{
		CrossId: "123",
	}

	type args struct {
		chainId         string
		blockHeight     int64
		crossCycleTxMap map[string]*db.CrossCycleTransaction
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:         ChainId1,
				blockHeight:     12,
				crossCycleTxMap: crossCycleTxMap,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossCycleTxDataCache(tt.args.chainId, tt.args.blockHeight, tt.args.crossCycleTxMap)
		})
	}
}

func TestSetCrossSubChainCrossCache(t *testing.T) {
	type args struct {
		chainId     string
		blockHeight int64
		dealResult  RealtimeDealResult
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:     ChainId1,
				blockHeight: 12,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossSubChainCrossCache(tt.args.chainId, tt.args.blockHeight, tt.args.dealResult)
		})
	}
}

func TestSetCrossTransfersCache(t *testing.T) {
	crossTransfers := []*db.CrossTransactionTransfer{
		{
			CrossId: "123",
		},
	}

	type args struct {
		chainId        string
		blockHeight    int64
		crossTransfers []*db.CrossTransactionTransfer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:        ChainId1,
				blockHeight:    12,
				crossTransfers: crossTransfers,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossTransfersCache(tt.args.chainId, tt.args.blockHeight, tt.args.crossTransfers)
		})
	}
}

func Test_executeDataInsertTasks(t *testing.T) {
	type args struct {
		chainId    string
		dealResult RealtimeDealResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				chainId: ChainId1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := executeDataInsertTasks(tt.args.chainId, tt.args.dealResult); (err != nil) != tt.wantErr {
				t.Errorf("executeDataInsertTasks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setDelayedUpdateCache(t *testing.T) {
	type args struct {
		chainId     string
		blockHeight int64
		dealResult  RealtimeDealResult
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:     ChainId1,
				blockHeight: 12,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setDelayedUpdateCache(tt.args.chainId, tt.args.blockHeight, tt.args.dealResult)
		})
	}
}

func TestRealtimeDataSaveToDB(t *testing.T) {
	crossCycleTx := &db.CrossCycleTransaction{
		CrossId:   "1234",
		StartTime: 123,
		EndTime:   456,
	}

	dealResult := RealtimeDealResult{
		BlockDetail: &db.Block{
			BlockHeight: 1,
			BlockHash:   "abc123",
		},
		UserList: map[string]*db.User{
			"user1": {
				UserId:   "user1",
				UserAddr: "address1",
				Role:     "admin",
			},
		},
		Transactions: map[string]*db.Transaction{
			"tx1": {
				TxId:        "tx1",
				Sender:      "user1",
				BlockHeight: 1,
				BlockHash:   "abc123",
				TxType:      "type1",
			},
		},
		CrossChainResult: &db.CrossChainResult{
			SaveCrossCycleTx: map[string]*db.CrossCycleTransaction{
				"1234": crossCycleTx,
			},
			UpdateCrossCycleTx: map[string]*db.CrossCycleTransaction{},
			InsertCrossCycleTx: []*db.CrossCycleTransaction{},
		},
	}

	type args struct {
		chainId     string
		blockHeight int64
		dealResult  RealtimeDealResult
		txTimeLog   *TxTimeLog
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:     ChainId1,
				blockHeight: 12,
				dealResult:  dealResult,
				txTimeLog:   &TxTimeLog{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//todo 1111111111111
			// if err := RealtimeDataSaveToDB(tt.args.chainId, tt.args.blockHeight, tt.args.dealResult, tt.args.txTimeLog); (err != nil) != tt.wantErr {
			// 	t.Errorf("RealtimeDataSaveToDB() error = %v, wantErr %v", err, tt.wantErr)
			// }
		})
	}
}
