/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/shopspring/decimal"
)

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

func TestGetMaxBlockHeight(t *testing.T) {
	type args struct {
		heightList []int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				heightList: []int64{4, 8, 6, 45, 34},
			},
			want: 45,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMaxBlockHeight(tt.args.heightList); got != tt.want {
				t.Errorf("GetMaxBlockHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMinBlockHeight(t *testing.T) {
	type args struct {
		heightList []int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				heightList: []int64{4, 8, 6, 45, 34},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMinBlockHeight(tt.args.heightList); got != tt.want {
				t.Errorf("GetMinBlockHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInBlockHeight(t *testing.T) {
	type args struct {
		height     int64
		heightList []int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				height:     8,
				heightList: []int64{4, 8, 6, 45, 34},
			},
			want: true,
		}, {
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				height:     9,
				heightList: []int64{4, 8, 6, 45, 34},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInBlockHeight(tt.args.height, tt.args.heightList); got != tt.want {
				t.Errorf("IsInBlockHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				v: "dsadadadas",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5(tt.args.v); got == "" {
				t.Errorf("MD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveAddrPrefix(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				address: "0x123456789",
			},
			want: "123456789",
		},
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				address: "123456789",
			},
			want: "123456789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveAddrPrefix(tt.args.address); got != tt.want {
				t.Errorf("RemoveAddrPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringAmountDecimal(t *testing.T) {
	//amountDecimal := decimal.NewFromFloat()
	amountDecimal, _ := decimal.NewFromString("123456.789")

	type args struct {
		amount   string
		decimals int
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				amount:   "123456789",
				decimals: 3,
			},
			want: amountDecimal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringAmountDecimal(tt.args.amount, tt.args.decimals)
			if !got.Equal(tt.want) {
				t.Errorf("StringAmountDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_copyMap(t *testing.T) {
}

func Test_getMemberInfoKey(t *testing.T) {
	memberInfo := "-----BEGIN CERTIFICATE-----\nMIICfjCCAiSgAwIBAgIDAxXYMAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMy5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmczLmNoYWlubWFrZXIub3JnMB4XDTIzMTIwMTA4NDMxNFoXDTI4\nMTEyOTA4NDMxNFowgZcxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmczLmNoYWlubWFrZXIub3Jn\nMRIwEAYDVQQLEwljb25zZW5zdXMxLzAtBgNVBAMTJmNvbnNlbnN1czEuc2lnbi53\neC1vcmczLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\nDtzS2xYTUfOYuvQq+c3KkL36ZTY+NCHartUHS8cCW54KQXRfDQRC0e9/JR1I7kp9\n9gy9fshvmESw/gQykfKX9qNqMGgwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCDl\newGepSE/dpdcq+nGCHeC0z7QpB6fXyZs64TuzBk/UjArBgNVHSMEJDAigCDGo5Qc\nLwYIuUF03wEa03op0tOteA4YhsvUuZEovpJXGTAKBggqhkjOPQQDAgNIADBFAiBl\nu4dLNllh91R5jkLOI2IWcPd4ht1jTh/zgt8MUEZF6AIhAJzM56k7bcRfqAwDeCDB\nzGF2T1NKUHFnqJu6YbDX8D5l\n-----END CERTIFICATE-----"
	type args struct {
		chainId     string
		hashType    string
		memberType  int32
		memberBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				chainId:     "chain1",
				hashType:    "SHA256",
				memberType:  0,
				memberBytes: []byte(memberInfo),
			},
			want: "2a5b80d1d933b32a93d059eb55008e3c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getMemberInfoKey(tt.args.chainId, tt.args.hashType, tt.args.memberType, tt.args.memberBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMemberInfoKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_isZeroAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				address: "0000000000000000000000000000000000000000",
			},
			want: true,
		},
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				address: "323232323232323",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isZeroAddress(tt.args.address); got != tt.want {
				t.Errorf("isZeroAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getContractEventTest(fileName string) []*db.ContractEventData {
	contractEvents := make([]*db.ContractEventData, 0)
	// 打开 JSON 文件
	//file, err := os.Open("../testData/1_blockInfoJsonContract.json")
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return contractEvents
	}
	err = decoder.Decode(&contractEvents)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return contractEvents
	}
	return contractEvents
}

func getContractInfoMapTest(fileName string) map[string]*db.Contract {
	resultValue := make(map[string]*db.Contract, 0)
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

func getAccountMapTest(fileName string) map[string]*db.Account {
	resultValue := make(map[string]*db.Account, 0)
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

func getPositionListJsonTest(fileName string) map[string]*db.PositionData {
	resultValue := make(map[string]*db.PositionData, 0)
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

type JsonData struct {
	UpdatePositionData map[string]*db.PositionData `json:"UpdatePositionData"`
	ResultPositionData ResultPositionData          `json:"ResultPositionData"`
}

type ResultPositionData struct {
	FungiblePosition    []*db.FungiblePosition    `json:"FungiblePosition"`
	NonFungiblePosition []*db.NonFungiblePosition `json:"NonFungiblePosition"`
}

func getUpdatePositionDataTest(fileName string) (map[string]*db.PositionData, []*db.FungiblePosition,
	[]*db.NonFungiblePosition) {
	//positionData := make(map[string]*db.PositionData, 0)
	//fungiblePosition := make([]*db.FungiblePosition, 0)
	//nonFungiblePosition := make([]*db.NonFungiblePosition, 0)
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil, nil
	}

	// 解码 JSON 文件内容到 blockInfo 结构体
	var jsonData JsonData
	err = decoder.Decode(&jsonData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, nil, nil
	}
	return jsonData.UpdatePositionData, jsonData.ResultPositionData.FungiblePosition,
		jsonData.ResultPositionData.NonFungiblePosition
}

func getPositionDBJsonTest(fileName string) map[string][]*db.FungiblePosition {
	resultValue := make(map[string][]*db.FungiblePosition, 0)
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}

	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func getNonPositionDBJsonTest(fileName string) map[string][]*db.NonFungiblePosition {
	resultValue := make(map[string][]*db.NonFungiblePosition, 0)
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

func getGotTokenResultTest(fileName string) *db.TokenResult {
	resultValue := &db.TokenResult{}
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

func getChainTxRWSetTest(fileName string) *common.TxRWSet {
	var resultValue *common.TxRWSet
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

func getBlockInfoTest(fileName string) *common.BlockInfo {
	var blockInfo *common.BlockInfo
	// 打开 JSON 文件
	//file, err := os.Open("../testData/1_blockInfoJsonContract.json")

	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return blockInfo
	}

	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&blockInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return blockInfo
	}
	return blockInfo
}

func getDealResultTest(fileName string) *RealtimeDealResult {
	var dealResult *RealtimeDealResult
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return dealResult
	}
	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&dealResult)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return dealResult
	}
	return dealResult
}

func getTxInfoInfoTest(fileName string) *common.Transaction {
	var txInfo *common.Transaction
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return txInfo
	}
	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&txInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return txInfo
	}
	return txInfo
}

func getBuildTxInfoTest(fileName string) *db.Transaction {
	var txInfo *db.Transaction
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return txInfo
	}
	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&txInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return txInfo
	}
	return txInfo
}

func getUserInfoInfoTest(fileName string) *MemberAddrIdCert {
	var userInfo *MemberAddrIdCert
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return userInfo
	}
	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&userInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return userInfo
	}
	return userInfo
}

func getCrossChainInfoTest(fileName string) *tcipCommon.CrossChainInfo {
	var resultValue *tcipCommon.CrossChainInfo
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}
	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func TestParallelParseBatchWhere(t *testing.T) {
	type args struct {
		wheres    []string
		batchSize int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "空切片",
			args: args{
				wheres:    []string{},
				batchSize: 2,
			},
			want: [][]string{},
		},
		{
			name: "正常情况",
			args: args{
				wheres:    []string{"a", "b", "c", "d", "e"},
				batchSize: 2,
			},
			want: [][]string{
				{"a", "b"},
				{"c", "d"},
				{"e"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParallelParseBatchWhere(tt.args.wheres, tt.args.batchSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParallelParseBatchWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEvmAbi(t *testing.T) {
	type args struct {
		evmType string
	}
	tests := []struct {
		name string
		args args
		want *abi.ABI
	}{
		{
			name: "获取EVM-DFA合约",
			args: args{
				evmType: ContractStandardNameEVMDFA,
			},
			want: config.GlobalAbiERC20,
		},
		{
			name: "获取EVM-NFA合约",
			args: args{
				evmType: ContractStandardNameEVMNFA,
			},
			want: config.GlobalAbiERC721,
		},
		{
			name: "获取未知类型合约",
			args: args{
				evmType: "unknown",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEvmAbi(tt.args.evmType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEvmAbi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsConfigTx(t *testing.T) {
	type args struct {
		txInfo *common.Transaction
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "空交易",
			args: args{
				txInfo: nil,
			},
			want: false,
		},
		{
			name: "有效的配置交易",
			args: args{
				txInfo: &common.Transaction{
					Payload: &common.Payload{
						ContractName: syscontract.SystemContract_CHAIN_CONFIG.String(),
					},
				},
			},
			want: true,
		},
		{
			name: "无效的配置交易",
			args: args{
				txInfo: &common.Transaction{
					Payload: &common.Payload{
						ContractName: "invalid_contract",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsConfigTx(tt.args.txInfo); got != tt.want {
				t.Errorf("IsConfigTx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsContractTx(t *testing.T) {
	type args struct {
		txInfo *common.Transaction
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "空交易",
			args: args{
				txInfo: nil,
			},
			want: false,
		},
		{
			name: "有效的合约交易",
			args: args{
				txInfo: &common.Transaction{
					Payload: &common.Payload{
						ContractName: syscontract.SystemContract_CONTRACT_MANAGE.String(),
					},
					Result: &common.Result{
						ContractResult: &common.ContractResult{
							Code: 0,
						},
					},
				},
			},
			want: true,
		},
		{
			name: "无效的合约交易",
			args: args{
				txInfo: &common.Transaction{
					Payload: &common.Payload{
						ContractName: "invalid_contract",
					},
					Result: &common.Result{
						ContractResult: &common.ContractResult{
							Code: 1,
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsContractTx(tt.args.txInfo); got != tt.want {
				t.Errorf("IsContractTx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMainChainGateway(t *testing.T) {
	type args struct {
		gatewayID string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "主链网关",
			args: args{
				gatewayID: tcipCommon.MainGateway_MAIN_GATEWAY_ID.String() + "_test",
			},
			want: true,
		},
		{
			name: "非主链网关",
			args: args{
				gatewayID: "invalid_gateway",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMainChainGateway(tt.args.gatewayID); got != tt.want {
				t.Errorf("IsMainChainGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRelayCrossChainTx(t *testing.T) {
	type args struct {
		txInfo *common.Transaction
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Nil transaction",
			args: args{
				txInfo: nil,
			},
			want: false,
		},
		{
			name: "Valid cross-chain transaction",
			args: args{
				txInfo: &common.Transaction{
					Payload: &common.Payload{
						ContractName: syscontract.SystemContract_RELAY_CROSS.String(),
					},
				},
			},
			want: true,
		},
		{
			name: "Invalid cross-chain transaction",
			args: args{
				txInfo: &common.Transaction{
					Payload: &common.Payload{
						ContractName: "invalid_contract",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRelayCrossChainTx(tt.args.txInfo); got != tt.want {
				t.Errorf("IsRelayCrossChainTx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSubChainSpvContractTx(t *testing.T) {
	type args struct {
		txInfo *common.Transaction
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{
			name: "Nil transaction",
			args: args{
				txInfo: nil,
			},
			want:  false,
			want1: "",
		},
		{
			name: "Valid sub-chain SPV contract transaction",
			args: args{
				txInfo: &common.Transaction{
					Payload: &common.Payload{
						ContractName: "official_spv_test",
					},
				},
			},
			want:  true,
			want1: "official_spv_test",
		},
		{
			name: "Invalid sub-chain SPV contract transaction",
			args: args{
				txInfo: &common.Transaction{
					Payload: &common.Payload{
						ContractName: "invalid_contract",
					},
				},
			},
			want:  false,
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := IsSubChainSpvContractTx(tt.args.txInfo)
			if got != tt.want {
				t.Errorf("IsSubChainSpvContractTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("IsSubChainSpvContractTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSaveJsonFile(t *testing.T) {
	type args struct {
		fix       string
		valueJson interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Save JSON file",
			args: args{
				fix: "test",
				valueJson: map[string]string{
					"key": "value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SaveJsonFile(tt.args.fix, tt.args.valueJson)
			// You can add assertions to check if the file was created and has the correct content
		})
	}
}

func TestIsCrossEnd(t *testing.T) {
	type args struct {
		status int32
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test:case 1",
			args: args{
				status: 0,
			},
			want: false,
		},
		{
			name: "test:case 2",
			args: args{
				status: 3,
			},
			want: true,
		},
		{
			name: "test:case 4",
			args: args{
				status: 4,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCrossEnd(tt.args.status); got != tt.want {
				t.Errorf("IsCrossEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}
