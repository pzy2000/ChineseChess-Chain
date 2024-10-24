/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/db"
	"context"
	"encoding/json"
	"testing"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/google/go-cmp/cmp"
)

func SetMemberInfoCache(chainId, hashType string, member *accesscontrol.Member) error {
	//缓存key
	memberKey, memberKeyErr := getMemberInfoKey(chainId, hashType, int32(member.MemberType), member.MemberInfo)
	if memberKeyErr != nil {
		return memberKeyErr
	}
	userInfoJson := getUserInfoInfoTest("0_userInfoJson.json")
	userInfoJsonBytes, _ := json.Marshal(userInfoJson)
	//缓存数据
	// 设置键值对和过期时间
	ctx := context.Background()
	err := cache.GlobalRedisDb.Set(ctx, memberKey, string(userInfoJsonBytes), time.Hour).Err()
	return err
}

func TestParallelParseContract(t *testing.T) {
	blockInfo := getBlockInfoTest("1_blockInfoJsonContract.json")
	txInfo := getTxInfoInfoTest("1_txInfoJsonContract.json")
	dealResult := getDealResultTest("1_dealResultJsonContract.json")
	chainId := blockInfo.Block.Header.ChainId
	err := SetMemberInfoCache(chainId, "SHA256", txInfo.Sender.Signer)
	if err != nil {
		return
	}

	type args struct {
		blockInfo  *common.BlockInfo
		hashType   string
		dealResult *RealtimeDealResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				blockInfo:  blockInfo,
				hashType:   "SHA256",
				dealResult: dealResult,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParallelParseContract(tt.args.blockInfo, tt.args.hashType, tt.args.dealResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParallelParseContract() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.args.dealResult.ContractWriteSetData) == 0 {
				t.Errorf("ParallelParseContract() error = %v, ContractWriteSetData is nil", err)
			}
		})
	}
}

func Test_containsAllFunctions(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJsonContractERC20.json")
	if blockInfo == nil {
		return
	}
	rwSetList := blockInfo.RwsetList[0]
	contractWriteSet, err := GetContractByWriteSet(rwSetList.TxWrites)
	if err != nil {
		return
	}
	if contractWriteSet.ContractResult == nil {
		return
	}

	signatures := ExtractFunctionSignatures(contractWriteSet.ByteCode)
	type args struct {
		evmType       string
		signatures    [][]byte
		functionNames map[string]bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1",
			args: args{
				evmType:       ContractStandardNameEVMDFA,
				signatures:    signatures,
				functionNames: copyMap(ERC20Functions),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = containsAllFunctions(tt.args.evmType, tt.args.signatures, tt.args.functionNames)
			//if got != tt.want {
			//	t.Errorf("containsAllFunctions() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestBuildContractInfo(t *testing.T) {
	blockInfo := getBlockInfoTest("1_blockInfoJsonContract.json")
	txInfo := getTxInfoInfoTest("1_txInfoJsonContract.json")
	userInfo := getUserInfoInfoTest("1_userInfoJsonContract.json")

	if blockInfo == nil || txInfo == nil || userInfo == nil {
		return
	}
	type args struct {
		i         int
		blockInfo *common.BlockInfo
		txInfo    *common.Transaction
		userInfo  *MemberAddrIdCert
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
				txInfo:    txInfo,
				userInfo:  userInfo,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildContractInfo(tt.args.i, tt.args.blockInfo, tt.args.txInfo, tt.args.userInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildContractInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetContractByWriteSet(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJsonContractERC20.json")
	if blockInfo == nil {
		return
	}
	rwSetList := blockInfo.RwsetList[0]
	type args struct {
		txWriteList []*common.TxWrite
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				txWriteList: rwSetList.TxWrites,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractByWriteSet(tt.args.txWriteList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByWriteSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ContractResult == nil || len(got.ByteCode) == 0 {
				t.Errorf("GetContractByWriteSet() got = %v", got)
			}
		})
	}
}

func Test_containsAllFunctions1(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJsonContractERC20.json")
	if blockInfo == nil {
		return
	}
	rwSetList := blockInfo.RwsetList[0]
	contractWriteSet, err := GetContractByWriteSet(rwSetList.TxWrites)
	if err != nil {
		return
	}
	if contractWriteSet.ContractResult == nil {
		return
	}

	signatures := ExtractFunctionSignatures(contractWriteSet.ByteCode)
	type args struct {
		evmType       string
		signatures    [][]byte
		functionNames map[string]bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1",
			args: args{
				evmType:       ContractStandardNameEVMDFA,
				signatures:    signatures,
				functionNames: copyMap(ERC20Functions),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = containsAllFunctions(tt.args.evmType, tt.args.signatures, tt.args.functionNames)
			//if got != tt.want {
			//	t.Errorf("containsAllFunctions() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestBuildContractInfo1(t *testing.T) {
	blockInfo := getBlockInfoTest("1_blockInfoJsonContract.json")
	txInfo := getTxInfoInfoTest("1_txInfoJsonContract.json")
	userInfo := getUserInfoInfoTest("1_userInfoJsonContract.json")

	if blockInfo == nil || txInfo == nil || userInfo == nil {
		return
	}
	type args struct {
		i         int
		blockInfo *common.BlockInfo
		txInfo    *common.Transaction
		userInfo  *MemberAddrIdCert
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
				txInfo:    txInfo,
				userInfo:  userInfo,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildContractInfo(tt.args.i, tt.args.blockInfo, tt.args.txInfo, tt.args.userInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildContractInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGenesisBlockSystemContract(t *testing.T) {
	type args struct {
		blockInfo  *common.BlockInfo
		dealResult *RealtimeDealResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "非创世区块",
			args: args{
				blockInfo: &common.BlockInfo{
					Block: &common.Block{
						Header: &common.BlockHeader{
							BlockHeight: 0,
						},
						Txs: []*common.Transaction{
							{
								Payload: &common.Payload{
									TxId: "1234",
								},
							},
						},
					},
					RwsetList: []*common.TxRWSet{
						{
							TxId: "1234",
						},
					},
				},
				dealResult: &RealtimeDealResult{},
			},
			wantErr: false,
		},
		// 添加其他测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenesisBlockSystemContract(tt.args.blockInfo, tt.args.dealResult); (err != nil) != tt.wantErr {
				t.Errorf("GenesisBlockSystemContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetContractSDKData(t *testing.T) {
	type args struct {
		chainId      string
		contractInfo *ContractWriteSetData
		byteCode     []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:      ChainId,
				contractInfo: &ContractWriteSetData{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetContractSDKData(tt.args.chainId, tt.args.contractInfo, tt.args.byteCode); (err != nil) != tt.wantErr {
				t.Errorf("GetContractSDKData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessContractInsertOrUpdate(t *testing.T) {
	writeSetData := &ContractWriteSetData{
		ContractName:    "chainName",
		ContractNameBak: "chainName",
		ContractSymbol:  "chainName",
		ContractAddr:    "123456789",
		ContractType:    "ContractType",
		Version:         "string",
		SenderTxId:      "123",
	}

	contractWriteSetMap := make(map[string]*ContractWriteSetData, 0)
	contractWriteSetMap["123"] = writeSetData
	dealResult := RealtimeDealResult{
		ContractWriteSetData: contractWriteSetMap,
	}

	type args struct {
		chainId    string
		dealResult RealtimeDealResult
	}
	tests := []struct {
		name    string
		args    args
		want    RealtimeDealResult
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				chainId:    ChainId1,
				dealResult: dealResult,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ProcessContractInsertOrUpdate(tt.args.chainId, tt.args.dealResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessContractInsertOrUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_dealStandardContract(t *testing.T) {
	contract1 := &db.Contract{
		ContractType: "ERC20",
		Name:         "ContractName",
		NameBak:      "ContractName",
		Addr:         "12345",
	}
	contract2 := &db.Contract{
		ContractType: "ERC721",
		Name:         "ContractName",
		NameBak:      "ContractName",
		Addr:         "12345",
	}
	want1 := &db.FungibleContract{
		ContractType:    "ERC20",
		ContractName:    "ContractName",
		ContractNameBak: "ContractName",
		ContractAddr:    "12345",
	}
	want2 := &db.NonFungibleContract{
		ContractType:    "ERC721",
		ContractName:    "ContractName",
		ContractNameBak: "ContractName",
		ContractAddr:    "12345",
	}

	type args struct {
		contract *db.Contract
	}
	tests := []struct {
		name  string
		args  args
		want  *db.FungibleContract
		want1 *db.NonFungibleContract
	}{
		{
			name: "Test case 1",
			args: args{
				contract: contract1,
			},
			want: want1,
		},
		{
			name: "Test case 2",
			args: args{
				contract: contract2,
			},
			want1: want2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := dealStandardContract(tt.args.contract)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("dealStandardContract() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
			if !cmp.Equal(got1, tt.want1) {
				t.Errorf("dealStandardContract() got1 = %v, want1 %v\ndiff: %s", got1, tt.want1, cmp.Diff(got1, tt.want1))
			}
		})
	}
}

func TestUpdateLatestContractCache(t *testing.T) {
	contract := &db.Contract{
		Name:    "nameBak",
		NameBak: "NameBak",
		Addr:    "123456",
	}

	contracts := []*db.Contract{
		contract,
	}

	SetLatestContractListCache(ChainId1, 12, contracts, contracts)
	type args struct {
		chainId         string
		updateContracts []*db.Contract
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1",
			args: args{
				updateContracts: []*db.Contract{
					contract,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateLatestContractCache(tt.args.chainId, tt.args.updateContracts)
		})
	}
}
