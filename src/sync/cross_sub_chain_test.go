package sync

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func convertToMap(data []*db.CrossSubChainData) map[string]int64 {
	result := make(map[string]int64)
	for _, item := range data {
		result[item.ChainId] = item.TxNum
	}
	return result
}

func TestBuildExecutionTransaction(t *testing.T) {
	crossChainInfo := getCrossChainInfoTest("cross_0_crossChainInfo.json")
	type args struct {
		txContent *tcipCommon.TxContent
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1",
			args: args{
				txContent: crossChainInfo.CrossChainTxContent[0].TxContent,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildExecutionTransaction(tt.args.txContent)
			if got == nil || got.TxId == "" {
				t.Errorf("BuildExecutionTransaction() got = %v", got)
			}
		})
	}
}

func TestDealCrossSubChainTxNum(t *testing.T) {
	type args struct {
		subChainIdMap  map[string]map[string]int64
		subChainDataDB map[string]*db.CrossSubChainData
	}

	subChainIdMap := map[string]map[string]int64{
		"chain1": {
			"subChain1": 1,
			"subChain2": 1,
		},
		"subChain1": {
			"chain1":    1,
			"subChain2": 1,
			"subChain3": 1,
		},
		"subChain2": {
			"chain1":    1,
			"subChain1": 1,
		},
		"subChain3": {
			"subChain1": 1,
		},
	}

	subChainDataDB := map[string]*db.CrossSubChainData{
		"subChain1": {
			ChainId: "subChain1",
			TxNum:   2,
		},
		"subChain2": {
			ChainId: "subChain2",
			TxNum:   2,
		},
		"subChain3": {
			ChainId: "subChain3",
			TxNum:   3,
		},
	}
	want := []*db.CrossSubChainData{
		{
			ChainId: "subChain1",
			TxNum:   5,
		},
		{
			ChainId: "subChain2",
			TxNum:   4,
		},
		{
			ChainId: "subChain3",
			TxNum:   4,
		},
	}
	tests := []struct {
		name string
		args args
		want []*db.CrossSubChainData
	}{
		{
			name: "Test case 1",
			args: args{
				subChainIdMap:  subChainIdMap,
				subChainDataDB: subChainDataDB,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealCrossSubChainTxNum(tt.args.subChainIdMap, tt.args.subChainDataDB)
			gotMap := convertToMap(got)
			wantMap := convertToMap(tt.want)

			if !reflect.DeepEqual(gotMap, wantMap) {
				t.Errorf("DealCrossSubChainTxNum() = %v, want %v", gotMap, wantMap)
			}
		})
	}
}

func TestDealSubChainCrossChainNum(t *testing.T) {
	subChainDataDB := []*db.CrossSubChainCrossChain{
		{
			SubChainId:  "subChain1",
			ChainId:     "chain1",
			ChainName:   "chain1",
			TxNum:       12,
			BlockHeight: 12,
		},
		{
			SubChainId:  "subChain1",
			ChainId:     "subChain2",
			ChainName:   "subChain2",
			TxNum:       10,
			BlockHeight: 10,
		},
		{
			SubChainId:  "subChain2",
			ChainId:     "subChain1",
			ChainName:   "subChain1",
			TxNum:       10,
			BlockHeight: 10,
		},
	}

	subChainIdMap := map[string]map[string]int64{
		"chain1": {
			"subChain1": 1,
		},
		"subChain1": {
			"chain1":    1,
			"subChain3": 1,
		},
		"subChain3": {
			"subChain1": 1,
		},
	}

	type args struct {
		chainId         string
		subChainIdMap   map[string]map[string]int64
		subChainCrossDB []*db.CrossSubChainCrossChain
		minHeight       int64
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossSubChainCrossChain
		want1   []*db.CrossSubChainCrossChain
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				subChainIdMap:   subChainIdMap,
				subChainCrossDB: subChainDataDB,
				minHeight:       30,
			},
			want: []*db.CrossSubChainCrossChain{
				{
					SubChainId:  "chain1",
					ChainId:     "subChain1",
					ChainName:   "",
					TxNum:       1,
					BlockHeight: 30,
				},
				{
					SubChainId:  "subChain1",
					ChainId:     "subChain3",
					ChainName:   "",
					TxNum:       1,
					BlockHeight: 30,
				},
				{
					SubChainId:  "subChain3",
					ChainId:     "subChain1",
					ChainName:   "",
					TxNum:       1,
					BlockHeight: 30,
				},
			},
			want1: []*db.CrossSubChainCrossChain{
				{
					SubChainId:  "subChain1",
					ChainId:     "chain1",
					ChainName:   "chain1",
					TxNum:       13,
					BlockHeight: 30,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DealSubChainCrossChainNum(tt.args.chainId, tt.args.subChainIdMap, tt.args.subChainCrossDB, tt.args.minHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("DealSubChainCrossChainNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 将实际结果和预期结果转换为映射
			gotMap := crossSubChainCrossChainSliceToMap(got)
			wantMap := crossSubChainCrossChainSliceToMap(tt.want)
			got1Map := crossSubChainCrossChainSliceToMap(got1)
			want1Map := crossSubChainCrossChainSliceToMap(tt.want1)

			if !reflect.DeepEqual(gotMap, wantMap) {
				t.Errorf("DealSubChainCrossChainNum() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1Map, want1Map) {
				t.Errorf("DealSubChainCrossChainNum() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func crossSubChainCrossChainSliceToMap(slice []*db.CrossSubChainCrossChain) map[string]map[string]int64 {
	result := make(map[string]map[string]int64)
	for _, item := range slice {
		if _, ok := result[item.SubChainId]; !ok {
			result[item.SubChainId] = make(map[string]int64)
		}
		result[item.SubChainId][item.ChainId] = item.TxNum
	}
	return result
}

func TestGetBusinessTransaction(t *testing.T) {
	type args struct {
		chainId        string
		crossChainInfo *tcipCommon.CrossChainInfo
	}
	crossChainInfo := getCrossChainInfoTest("cross_0_crossChainInfo.json")

	tests := []struct {
		name    string
		args    args
		wantLen int
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:        "chain1",
				crossChainInfo: crossChainInfo,
			},
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetBusinessTransaction(tt.args.chainId, tt.args.crossChainInfo)
			if len(got) != tt.wantLen {
				t.Errorf("GetBusinessTransaction() = %v, wantLen %v", got, tt.wantLen)
			}
		})
	}
}

func TestGetCrossCycleTransaction(t *testing.T) {
	crossChainInfo := getCrossChainInfoTest("cross_0_crossChainInfo.json")

	type args struct {
		crossChainInfo *tcipCommon.CrossChainInfo
		blockHeight    int64
		timestamp      int64
	}
	tests := []struct {
		name string
		args args
		want *db.CrossCycleTransaction
	}{
		{
			name: "Test case 1",
			args: args{
				timestamp:      123456789,
				blockHeight:    12,
				crossChainInfo: crossChainInfo,
			},
			want: &db.CrossCycleTransaction{
				CrossId:     "0",
				Status:      3,
				StartTime:   123456789,
				EndTime:     123456789,
				BlockHeight: 12,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCrossCycleInsertTx(tt.args.crossChainInfo, tt.args.blockHeight, tt.args.timestamp)
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.CrossCycleTransaction{}, "CreatedAt", "BlockHeight", "UpdatedAt", "ID")) {
				t.Errorf("GetCrossCycleTransaction() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetCrossTxTransfer1(t *testing.T) {
	crossChainInfo := getCrossChainInfoTest("cross_1_ChainInfoJson.json")

	type args struct {
		chainId        string
		blockHeight    int64
		crossChainInfo *tcipCommon.CrossChainInfo
	}
	tests := []struct {
		name string
		args args
		want []*db.CrossTransactionTransfer
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:        "chain1",
				blockHeight:    19,
				crossChainInfo: crossChainInfo,
			},
			want: []*db.CrossTransactionTransfer{
				{
					CrossId:         "0",
					FromGatewayId:   "MAIN_GATEWAY_ID-relay-1",
					FromChainId:     "chain1",
					FromIsMainChain: true,
					ToGatewayId:     "0",
					ToIsMainChain:   false,
					ToChainId:       "chainmaker001",
					BlockHeight:     19,
					ContractName:    "crossChainSaveQuery",
					ContractMethod:  "save",
					Parameter:       "{\"key\":\"main_son_save_key\",\"value\":\"main_son_save_value\"}",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCrossTxTransfer(tt.args.chainId, tt.args.blockHeight, tt.args.crossChainInfo)
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.CrossTransactionTransfer{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetCrossTxTransfer() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetMainCrossTransaction1(t *testing.T) {
	txInfo := getTxInfoInfoTest("cross_1_txInfoJson.json")
	crossChainInfo := getCrossChainInfoTest("cross_1_ChainInfoJson.json")

	type args struct {
		crossChainInfo *tcipCommon.CrossChainInfo
		txInfo         *common.Transaction
	}
	tests := []struct {
		name string
		args args
		want *db.CrossMainTransaction
	}{
		{
			name: "Test case 1",
			args: args{
				txInfo:         txInfo,
				crossChainInfo: crossChainInfo,
			},
			want: &db.CrossMainTransaction{
				TxId:      "AC1052E93C62A79B5ECF942AC1BF37ED8CD9F9593873DBC81807DF8CBE5E451F",
				CrossId:   "0",
				ChainMsg:  "[{\"gateway_id\":\"0\",\"chain_rid\":\"chainmaker001\",\"contract_name\":\"crossChainSaveQuery\",\"method\":\"save\",\"parameter\":\"{\\\"key\\\":\\\"main_son_save_key\\\",\\\"value\\\":\\\"main_son_save_value\\\"}\",\"confirm_info\":{\"chain_rid\":\"chainmaker001\",\"contract_name\":\"crossChainSaveQuery\",\"method\":\"saveState\",\"parameter\":\"{\\\"key\\\":\\\"main_son_save_key\\\",\\\"state\\\":\\\"confirmEnd\\\"}\"},\"cancel_info\":{\"chain_rid\":\"chainmaker001\",\"contract_name\":\"crossChainSaveQuery\",\"method\":\"saveState\",\"parameter\":\"{\\\"key\\\":\\\"main_son_save_key\\\",\\\"state\\\":\\\"cancelEnd\\\"}\"}}]",
				Status:    1,
				CrossType: 1,
				Timestamp: 1709177844,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMainCrossTransaction(tt.args.crossChainInfo, tt.args.txInfo)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMainCrossTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCrossCycleTxTransfer1(t *testing.T) {
	type args struct {
		transfers []*db.CrossTransactionTransfer
	}
	transfers := []*db.CrossTransactionTransfer{
		{
			CrossId:         "1",
			FromChainId:     "chain1",
			FromIsMainChain: true,
			ToChainId:       "subChain1",
			ToIsMainChain:   false,
		},
		{
			CrossId:         "1",
			FromChainId:     "chain1",
			FromIsMainChain: true,
			ToChainId:       "subChain2",
			ToIsMainChain:   false,
		},
		{
			CrossId:         "2",
			FromChainId:     "subChain1",
			FromIsMainChain: false,
			ToChainId:       "subChain2",
			ToIsMainChain:   false,
		},
		{
			CrossId:         "2",
			FromChainId:     "subChain1",
			FromIsMainChain: false,
			ToChainId:       "subChain3",
			ToIsMainChain:   false,
		},
	}

	want := map[string]map[string]int64{
		"chain1": {
			"subChain1": 1,
			"subChain2": 1,
		},
		"subChain1": {
			"chain1":    1,
			"subChain2": 1,
			"subChain3": 1,
		},
		"subChain2": {
			"chain1":    1,
			"subChain1": 1,
		},
		"subChain3": {
			"subChain1": 1,
		},
	}

	tests := []struct {
		name string
		args args
		want map[string]map[string]int64
	}{
		{
			name: "Test case 1",
			args: args{
				transfers: transfers,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCrossCycleTxTransfer(tt.args.transfers)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCrossCycleTxTransfer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossTxTransfer(t *testing.T) {
	crossChainInfo := getCrossChainInfoTest("cross_1_ChainInfoJson.json")

	type args struct {
		chainId        string
		blockHeight    int64
		crossChainInfo *tcipCommon.CrossChainInfo
	}
	tests := []struct {
		name string
		args args
		want []*db.CrossTransactionTransfer
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:        "chain1",
				blockHeight:    19,
				crossChainInfo: crossChainInfo,
			},
			want: []*db.CrossTransactionTransfer{
				{
					CrossId:         "0",
					FromGatewayId:   "MAIN_GATEWAY_ID-relay-1",
					FromChainId:     "chain1",
					FromIsMainChain: true,
					ToGatewayId:     "0",
					ToIsMainChain:   false,
					ToChainId:       "chainmaker001",
					BlockHeight:     19,
					ContractName:    "crossChainSaveQuery",
					ContractMethod:  "save",
					Parameter:       "{\"key\":\"main_son_save_key\",\"value\":\"main_son_save_value\"}",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCrossTxTransfer(tt.args.chainId, tt.args.blockHeight, tt.args.crossChainInfo)
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.CrossTransactionTransfer{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetCrossTxTransfer() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetMainCrossTransaction(t *testing.T) {
	txInfo := getTxInfoInfoTest("cross_1_txInfoJson.json")
	crossChainInfo := getCrossChainInfoTest("cross_1_ChainInfoJson.json")

	type args struct {
		crossChainInfo *tcipCommon.CrossChainInfo
		txInfo         *common.Transaction
	}
	tests := []struct {
		name string
		args args
		want *db.CrossMainTransaction
	}{
		{
			name: "Test case 1",
			args: args{
				txInfo:         txInfo,
				crossChainInfo: crossChainInfo,
			},
			want: &db.CrossMainTransaction{
				TxId:      "AC1052E93C62A79B5ECF942AC1BF37ED8CD9F9593873DBC81807DF8CBE5E451F",
				CrossId:   "0",
				ChainMsg:  "[{\"gateway_id\":\"0\",\"chain_rid\":\"chainmaker001\",\"contract_name\":\"crossChainSaveQuery\",\"method\":\"save\",\"parameter\":\"{\\\"key\\\":\\\"main_son_save_key\\\",\\\"value\\\":\\\"main_son_save_value\\\"}\",\"confirm_info\":{\"chain_rid\":\"chainmaker001\",\"contract_name\":\"crossChainSaveQuery\",\"method\":\"saveState\",\"parameter\":\"{\\\"key\\\":\\\"main_son_save_key\\\",\\\"state\\\":\\\"confirmEnd\\\"}\"},\"cancel_info\":{\"chain_rid\":\"chainmaker001\",\"contract_name\":\"crossChainSaveQuery\",\"method\":\"saveState\",\"parameter\":\"{\\\"key\\\":\\\"main_son_save_key\\\",\\\"state\\\":\\\"cancelEnd\\\"}\"}}]",
				Status:    1,
				CrossType: 1,
				Timestamp: 1709177844,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMainCrossTransaction(tt.args.crossChainInfo, tt.args.txInfo)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMainCrossTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCrossCycleTxTransfer(t *testing.T) {
	type args struct {
		transfers []*db.CrossTransactionTransfer
	}
	transfers := []*db.CrossTransactionTransfer{
		{
			CrossId:         "1",
			FromChainId:     "chain1",
			FromIsMainChain: true,
			ToChainId:       "subChain1",
			ToIsMainChain:   false,
		},
		{
			CrossId:         "1",
			FromChainId:     "chain1",
			FromIsMainChain: true,
			ToChainId:       "subChain2",
			ToIsMainChain:   false,
		},
		{
			CrossId:         "2",
			FromChainId:     "subChain1",
			FromIsMainChain: false,
			ToChainId:       "subChain2",
			ToIsMainChain:   false,
		},
		{
			CrossId:         "2",
			FromChainId:     "subChain1",
			FromIsMainChain: false,
			ToChainId:       "subChain3",
			ToIsMainChain:   false,
		},
	}

	want := map[string]map[string]int64{
		"chain1": {
			"subChain1": 1,
			"subChain2": 1,
		},
		"subChain1": {
			"chain1":    1,
			"subChain2": 1,
			"subChain3": 1,
		},
		"subChain2": {
			"chain1":    1,
			"subChain1": 1,
		},
		"subChain3": {
			"subChain1": 1,
		},
	}

	tests := []struct {
		name string
		args args
		want map[string]map[string]int64
	}{
		{
			name: "Test case 1",
			args: args{
				transfers: transfers,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCrossCycleTxTransfer(tt.args.transfers)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCrossCycleTxTransfer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildCrossSubChainData(t *testing.T) {
	type args struct {
		gateWayIds []int64
		dealResult *RealtimeDealResult
		timestamp  int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				gateWayIds: []int64{
					1, 2, 3,
				},
				dealResult: &RealtimeDealResult{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BuildCrossSubChainData(tt.args.gateWayIds, tt.args.dealResult, tt.args.timestamp); (err != nil) != tt.wantErr {
				t.Errorf("BuildCrossSubChainData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fetchCrossSubChainData(t *testing.T) {
	type args struct {
		gateWayId int64
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want []*db.CrossSubChainData
	}{
		{
			name: "Test case 1",
			args: args{
				gateWayId: 1,
				timestamp: 1234,
			},
			want: []*db.CrossSubChainData{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fetchCrossSubChainData(tt.args.gateWayId, tt.args.timestamp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchCrossSubChainData() = %v, want %v", got, tt.want)
			}
		})
	}
}
