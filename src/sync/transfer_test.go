package sync

import (
	"chainmaker_web/src/db"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
)

func Test_dealTransferList(t *testing.T) {
	dbFungibleTransfer := make([]*db.FungibleTransfer, 0)
	dbNonFungibleTransfer := make([]*db.NonFungibleTransfer, 0)
	type args struct {
		eventDataList   []*db.ContractEventData
		contractInfoMap map[string]*db.Contract
		txInfoMap       map[string]*db.Transaction
	}
	contractEvents := getContractEventTest("5_contractEventDataToken.json")
	contractInfoMap := getContractInfoMapTest("0_contractInfoMap.json")

	gotResultJson := "[{\"txId\":\"17a09dcaf9a2484dcad8e7a648c1923b7dbe731c43574a6e9f0d6d7372048619\",\"eventIndex\":0,\"contractName\":\"goErc20_1\",\"contractAddr\":\"aba31ce4cd49f08073d2f115eb12610544242123\",\"contractMethod\":\"\",\"blockTime\":0,\"topic\":\"transfer\",\"fromAddr\":\"\",\"toAddr\":\"18fc4e7429af8419d5bb307e34db398b9a2331c6\",\"tokenId\":\"10000\",\"timestamp\":1702534154,\"ID\":0,\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null},{\"txId\":\"17a09dcaf9a2484dcad8e7a648c1923b7dbe731c43574a6e9f0d6d7372048619\",\"eventIndex\":2,\"contractName\":\"goErc20_1\",\"contractAddr\":\"aba31ce4cd49f08073d2f115eb12610544242123\",\"contractMethod\":\"\",\"blockTime\":0,\"topic\":\"transfer\",\"fromAddr\":\"18fc4e7429af8419d5bb307e34db398b9a2331c6\",\"toAddr\":\"18fc4e7429af8419d5bb307e34db398b9a233112\",\"tokenId\":\"10000\",\"timestamp\":1702534154,\"ID\":0,\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null},{\"txId\":\"17a09dcaf9a2484dcad8e7a648c1923b7dbe731c43574a6e9f0d6d7372048600\",\"eventIndex\":2,\"contractName\":\"goErc20_1\",\"contractAddr\":\"aba31ce4cd49f08073d2f115eb12610544242123\",\"contractMethod\":\"\",\"blockTime\":0,\"topic\":\"transfer\",\"fromAddr\":\"171262347a59fded92021a32421a5dad05424e03\",\"toAddr\":\"18fc4e7429af8419d5bb307e34db398b9a2331c6\",\"tokenId\":\"12345\",\"timestamp\":1702534154,\"ID\":0,\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null},{\"txId\":\"17a09dcaf9a2484dcad8e7a648c1923b7dbe731c43574a6e9f0d6d7372048600\",\"eventIndex\":3,\"contractName\":\"goErc20_1\",\"contractAddr\":\"aba31ce4cd49f08073d2f115eb12610544242123\",\"contractMethod\":\"\",\"blockTime\":0,\"topic\":\"transfer\",\"fromAddr\":\"171262347a59fded92021a32421a5dad05424e03\",\"toAddr\":\"\",\"tokenId\":\"20000\",\"timestamp\":1702534154,\"ID\":0,\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null}]"
	err := json.Unmarshal([]byte(gotResultJson), &dbNonFungibleTransfer)
	if err != nil {
		return
	}

	contractInfo := &db.Contract{
		Addr:         "1234",
		Name:         "ContractName",
		NameBak:      "ContractName",
		ContractType: "CMNFA",
	}

	tests := []struct {
		name  string
		args  args
		want  []*db.FungibleTransfer
		want1 []*db.NonFungibleTransfer
	}{
		{
			name: "test case 1",
			args: args{
				eventDataList:   contractEvents,
				contractInfoMap: contractInfoMap,
				txInfoMap:       nil,
			},
			want:  dbFungibleTransfer,
			want1: dbNonFungibleTransfer,
		},
		{
			name: "test case 2",
			args: args{
				eventDataList: []*db.ContractEventData{
					{
						Topic:        "Mint",
						TxId:         "21121212",
						ContractName: "ContractName",
						EventData:    &db.TransferTopicEventData{},
					},
				},
				contractInfoMap: map[string]*db.Contract{
					"ContractName": contractInfo,
				},
				txInfoMap: nil,
			},
			want:  dbFungibleTransfer,
			want1: dbNonFungibleTransfer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := dealTransferList(tt.args.eventDataList, tt.args.contractInfoMap, tt.args.txInfoMap)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dealTransferList() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("dealTransferList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_dealTransferList2(t *testing.T) {
	type args struct {
		eventDataList   []*db.ContractEventData
		contractInfoMap map[string]*db.Contract
		txInfoMap       map[string]*db.Transaction
	}
	contractEvents := getContractEventTest("5_contractEventDataAmount.json")
	contractInfoMap := getContractInfoMapTest("0_contractInfoMap.json")

	dbFungibleTransferMap := []*db.FungibleTransfer{
		{
			TxId:         "17a09dcaf9a2484dcad8e7a648c1923b7dbe731c43574a6e9f0d6d7372048619",
			EventIndex:   0,
			ContractName: "goErc20_1",
			ContractAddr: "aba31ce4cd49f08073d2f115eb12610544242ff9",
			Topic:        "mint",
			FromAddr:     "",
			ToAddr:       "18fc4e7429af8419d5bb307e34db398b9a2331c6",
			Amount:       GetNumberDecimal("1230000000000123.456789123456789123"),
			Timestamp:    1702534154,
		},
		{
			TxId:         "17a09dcaf9a2484dcad8e7a648c1923b7dbe731c43574a6e9f0d6d7372048618",
			EventIndex:   2,
			ContractName: "goErc20_1",
			ContractAddr: "aba31ce4cd49f08073d2f115eb12610544242ff9",
			Topic:        "transfer",
			FromAddr:     "18fc4e7429af8419d5bb307e34db398b9a2331c6",
			ToAddr:       "18fc4e7429af8419d5bb307e34db398b9a233112",
			Amount:       GetNumberDecimal("0.00000000000000001"),
			Timestamp:    1702534154,
		},
		{
			TxId:         "17a09dcaf9a2484dcad8e7a648c1923b7dbe731c43574a6e9f0d6d7372048607",
			EventIndex:   2,
			ContractName: "goErc20_1",
			ContractAddr: "aba31ce4cd49f08073d2f115eb12610544242ff9",
			Topic:        "transfer",
			FromAddr:     "171262347a59fded92021a32421a5dad05424e03",
			ToAddr:       "18fc4e7429af8419d5bb307e34db398b9a2331c6",
			Amount:       GetNumberDecimal("0.000000000000000009"),
			Timestamp:    1702534154,
		},
		{
			TxId:         "17a09dcaf9a2484dcad8e7a648c1923b7dbe731c43574a6e9f0d6d7372048606",
			EventIndex:   3,
			ContractName: "goErc20_1",
			ContractAddr: "aba31ce4cd49f08073d2f115eb12610544242ff9",
			Topic:        "burn",
			FromAddr:     "171262347a59fded92021a32421a5dad05424e03",
			ToAddr:       "",
			Amount:       GetNumberDecimal("0.00000000000001"),
			Timestamp:    1702534154},
	}

	tests := []struct {
		name  string
		args  args
		want  []*db.FungibleTransfer
		want1 []*db.NonFungibleTransfer
	}{
		{
			name: "test case 1",
			args: args{
				eventDataList:   contractEvents,
				contractInfoMap: contractInfoMap,
				txInfoMap:       nil,
			},
			want:  dbFungibleTransferMap,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := dealTransferList(tt.args.eventDataList, tt.args.contractInfoMap, tt.args.txInfoMap)
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.FungibleTransfer{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("mergeChainInfo() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func convertTransferToMap(data []*db.FungibleTransfer) map[string]*db.FungibleTransfer {
	result := make(map[string]*db.FungibleTransfer)
	for _, item := range data {
		result[item.TxId] = item
	}
	return result
}
