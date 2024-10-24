package sync

import (
	"chainmaker_web/src/db"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"
	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
)

func TestGetChainConfigByWriteSet(t *testing.T) {
	type args struct {
		txRWSet *common.TxRWSet
		txInfo  *common.Transaction
	}
	txInfo := getTxInfoInfoTest("1_txInfoJsonContract.json")
	tests := []struct {
		name    string
		args    args
		want    *pbConfig.ChainConfig
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				txRWSet: nil,
				txInfo:  txInfo,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetChainConfigByWriteSet(tt.args.txRWSet, tt.args.txInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChainConfigByWriteSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChainConfigByWriteSet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossChainInfoByWriteSet(t *testing.T) {
	txInfo := getTxInfoInfoTest("cross_1_txInfoJson.json")
	txRWSet := getChainTxRWSetTest("cross_1_txRWSet.json")
	crossChainInfo := getCrossChainInfoTest("cross_1_ChainInfoJson.json")

	type args struct {
		txRWSet *common.TxRWSet
		txInfo  *common.Transaction
	}
	tests := []struct {
		name    string
		args    args
		want    *tcipCommon.CrossChainInfo
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				txRWSet: txRWSet,
				txInfo:  txInfo,
			},
			want:    crossChainInfo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossChainInfoByWriteSet(tt.args.txRWSet, tt.args.txInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossChainInfoByWriteSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossChainInfoByWriteSet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSubChainBlockHeightByWriteSet(t *testing.T) {
	txInfo := getTxInfoInfoTest("cross_2_txInfoJson_1709178475.json")
	txRWSet := getChainTxRWSetTest("cross_2_txRWSetJson_1709178475.json")
	type args struct {
		txRWSet *common.TxRWSet
		txInfo  *common.Transaction
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 string
	}{
		{
			name: "test case 1",
			args: args{
				txRWSet: txRWSet,
				txInfo:  txInfo,
			},
			want:  14,
			want1: "official_spv0chainmaker001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetSubChainBlockHeightByWriteSet(tt.args.txRWSet, tt.args.txInfo)
			if got != tt.want {
				t.Errorf("GetSubChainBlockHeightByWriteSet() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetSubChainBlockHeightByWriteSet() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetSubChainGatewayIdByWriteSet(t *testing.T) {
	txInfo := getTxInfoInfoTest("cross_3_txInfoJson_1709189120.json")
	txRWSet := getChainTxRWSetTest("cross_3_txRWSetJson_1709189120.json")
	type args struct {
		txRWSet *common.TxRWSet
		txInfo  *common.Transaction
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "test case 1",
			args: args{
				txRWSet: txRWSet,
				txInfo:  txInfo,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSubChainGatewayIdByWriteSet(tt.args.txRWSet, tt.args.txInfo)
			if got == nil || !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("GetSubChainGatewayIdByWriteSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processCrossChainInfo(t *testing.T) {
	txInfo := getTxInfoInfoTest("cross_1_txInfoJson.json")
	crossChainInfo := getCrossChainInfoTest("cross_1_ChainInfoJson.json")

	var (
		wantGot  *db.CrossMainTransaction
		wantGot1 []*db.CrossTransactionTransfer
		wantGot2 = make(map[string]*db.CrossBusinessTransaction, 0)
	)
	wantGotJson := "{\"txId\":\"AC1052E93C62A79B5ECF942AC1BF37ED8CD9F9593873DBC81807DF8CBE5E451F\",\"crossId\":\"0\",\"chainMsg\":\"[{\\\"gateway_id\\\":\\\"0\\\",\\\"chain_rid\\\":\\\"chainmaker001\\\",\\\"contract_name\\\":\\\"crossChainSaveQuery\\\",\\\"method\\\":\\\"save\\\",\\\"parameter\\\":\\\"{\\\\\\\"key\\\\\\\":\\\\\\\"main_son_save_key\\\\\\\",\\\\\\\"value\\\\\\\":\\\\\\\"main_son_save_value\\\\\\\"}\\\",\\\"confirm_info\\\":{\\\"chain_rid\\\":\\\"chainmaker001\\\",\\\"contract_name\\\":\\\"crossChainSaveQuery\\\",\\\"method\\\":\\\"saveState\\\",\\\"parameter\\\":\\\"{\\\\\\\"key\\\\\\\":\\\\\\\"main_son_save_key\\\\\\\",\\\\\\\"state\\\\\\\":\\\\\\\"confirmEnd\\\\\\\"}\\\"},\\\"cancel_info\\\":{\\\"chain_rid\\\":\\\"chainmaker001\\\",\\\"contract_name\\\":\\\"crossChainSaveQuery\\\",\\\"method\\\":\\\"saveState\\\",\\\"parameter\\\":\\\"{\\\\\\\"key\\\\\\\":\\\\\\\"main_son_save_key\\\\\\\",\\\\\\\"state\\\\\\\":\\\\\\\"cancelEnd\\\\\\\"}\\\"}}]\",\"status\":1,\"crossType\":1,\"timestamp\":1709177844,\"ID\":0,\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null}"
	wantGot1Json := "[{\"crossId\":\"0\",\"fromGatewayId\":\"MAIN_GATEWAY_ID-relay-1\",\"fromChainId\":\"chain1\",\"fromIsMainChain\":true,\"toGatewayId\":\"0\",\"toChainId\":\"chainmaker001\",\"toIsMainChain\":false,\"blockHeight\":12,\"contractName\":\"crossChainSaveQuery\",\"contractMethod\":\"save\",\"parameter\":\"{\\\"key\\\":\\\"main_son_save_key\\\",\\\"value\\\":\\\"main_son_save_value\\\"}\",\"ID\":0,\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null}]"
	_ = json.Unmarshal([]byte(wantGotJson), &wantGot)
	_ = json.Unmarshal([]byte(wantGot1Json), &wantGot1)

	type args struct {
		chainId        string
		blockHeight    int64
		crossChainInfo *tcipCommon.CrossChainInfo
		txInfo         *common.Transaction
	}
	tests := []struct {
		name  string
		args  args
		want  *db.CrossMainTransaction
		want1 []*db.CrossTransactionTransfer
		want2 map[string]*db.CrossBusinessTransaction
	}{
		{
			name: "test case 1",
			args: args{
				chainId:        "chain1",
				blockHeight:    12,
				crossChainInfo: crossChainInfo,
				txInfo:         txInfo,
			},
			want:  wantGot,
			want1: wantGot1,
			want2: wantGot2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, _, _, _ := processCrossChainInfo(tt.args.chainId, tt.args.blockHeight, tt.args.crossChainInfo, tt.args.txInfo)
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.CrossMainTransaction{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("processCrossChainInfo() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}

			if !cmp.Equal(got1, tt.want1, cmpopts.IgnoreFields(db.CrossTransactionTransfer{}, "CreatedAt", "UpdatedAt", "BlockHeight", "ID")) {
				t.Errorf("processCrossChainInfo() got1 = %v, want1 %v\ndiff: %s", got1, tt.want1, cmp.Diff(got1, tt.want1))
			}

			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("processCrossChainInfo() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
