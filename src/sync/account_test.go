/*
 * @Author: dongxuliang dongxuliang@tencent.com
 * @Date: 2024-07-15 14:22:41
 * @LastEditors: dongxuliang dongxuliang@tencent.com
 * @LastEditTime: 2024-07-15 15:16:54
 * @FilePath: /chainmaker-explorer-backend/src/sync/account_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package sync

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildAccountInsertOrUpdate(t *testing.T) {
	getDBResult := &GetDBResult{
		PositionMapList:        make(map[string][]*db.FungiblePosition, 0),
		NonPositionMapList:     make(map[string][]*db.NonFungiblePosition, 0),
		FungibleContractMap:    make(map[string]*db.FungibleContract, 0),
		NonFungibleContractMap: make(map[string]*db.NonFungibleContract, 0),
		AddBlackTxList:         make([]*db.Transaction, 0),
		DeleteBlackTxList:      make([]*db.BlackTransaction, 0),
		CrossSubChainCross:     make([]*db.CrossSubChainCrossChain, 0),
		CrossSubChainMap:       make(map[string]*db.CrossSubChainData, 0),
		AccountBNSList:         make([]*db.Account, 0),
		AccountDIDList:         make([]*db.Account, 0),
		AccountDBMap:           make(map[string]*db.Account, 0),
	}
	topicEventResult := &TopicEventResult{
		AddBlack:          make([]string, 0),
		DeleteBlack:       make([]string, 0),
		IdentityContract:  make([]*db.IdentityContract, 0),
		ContractEventData: make([]*db.ContractEventData, 0),
		OwnerAdders: []string{
			"1234",
			"2234",
			"3234",
		},
		DIDAccount:       make(map[string][]string, 0),
		DIDUnBindList:    make([]string, 0),
		BNSBindEventData: make([]*db.BNSTopicEventData, 0),
		BNSUnBindDomain:  make([]string, 0),
	}

	type args struct {
		chainId          string
		delayGetDBResult *GetDBResult
		topicEventResult *TopicEventResult
		txList           map[string]*db.Transaction
		minHeight        int64
		accountTx        map[string]int64
		accountNFT       map[string]int64
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Account
		want1   []*db.Account
		want2   map[string]*db.Account
		wantErr bool
	}{
		{
			name: "测试构建账户插入或更新",
			args: args{
				chainId:          ChainId1,
				delayGetDBResult: getDBResult,
				topicEventResult: topicEventResult,
				minHeight:        13,
				txList:           make(map[string]*db.Transaction, 0),
			},
			want: []*db.Account{
				{
					Address:     "1234",
					BlockHeight: 13,
				},
				{
					Address:     "2234",
					BlockHeight: 13,
				},
				{
					Address:     "3234",
					BlockHeight: 13,
				},
			},
			want1: nil,
			want2: map[string]*db.Account{
				"1234": {Address: "1234", BlockHeight: 13},
				"2234": {Address: "2234", BlockHeight: 13},
				"3234": {Address: "3234", BlockHeight: 13},
			},
			wantErr: false,
		},
		// 添加其他测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := BuildAccountInsertOrUpdate(tt.args.chainId, tt.args.minHeight, tt.args.delayGetDBResult, tt.args.topicEventResult, tt.args.accountTx, tt.args.accountNFT)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildAccountInsertOrUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sortAccounts(got)
			sortAccounts(tt.want)

			if !cmp.Equal(got, tt.want) {
				t.Errorf("buildGasInfo() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
			if !cmp.Equal(got1, tt.want1) {
				t.Errorf("buildGasInfo() got1 = %v, want1 %v\ndiff: %s", got1, tt.want1, cmp.Diff(got1, tt.want1))
			}
			if !cmp.Equal(got2, tt.want2) {
				t.Errorf("buildGasInfo() got2 = %v, want2 %v\ndiff: %s", got2, tt.want2, cmp.Diff(got2, tt.want2))
			}
		})
	}
}

func Test_processBNSAccounts(t *testing.T) {
	type args struct {
		bnsBindEventData []*db.BNSTopicEventData
		unBindBNSs       []*db.Account
		accountInsertMap map[string]*db.Account
		accountUpdateMap map[string]*db.Account
		accountMap       map[string]*db.Account
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "测试处理BNS帐户",
			args: args{
				bnsBindEventData: []*db.BNSTopicEventData{
					{
						Domain: "1234",
						Value:  "12345",
					},
				},
				unBindBNSs: []*db.Account{
					{
						Address: "123456",
						BNS:     "1234",
					},
				},
				accountInsertMap: map[string]*db.Account{
					"12345": {
						Address: "123456",
						BNS:     "1234",
						DID:     "did:32345",
					},
				},
				accountUpdateMap: map[string]*db.Account{
					"67890": {
						Address: "67890",
						BNS:     "1234",
						DID:     "did:32345",
					},
				},
				accountMap: map[string]*db.Account{
					"12345": {
						Address: "123456",
						BNS:     "1234",
						DID:     "did:32345",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processBNSAccounts(tt.args.bnsBindEventData, tt.args.unBindBNSs, tt.args.accountInsertMap, tt.args.accountUpdateMap, tt.args.accountMap)
		})
	}
}

func Test_processDIDAccounts(t *testing.T) {
	type args struct {
		didAccount       map[string][]string
		unBindDIDs       []*db.Account
		accountInsertMap map[string]*db.Account
		accountUpdateMap map[string]*db.Account
		accountMap       map[string]*db.Account
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "测试处理BNS帐户",
			args: args{
				didAccount: map[string][]string{
					"did:12345": {
						"12345",
						"22345",
					},
					"did:22345": {
						"67890",
						"6789000",
					},
				},
				unBindDIDs: []*db.Account{
					{
						Address: "123456",
						BNS:     "1234",
						DID:     "did:32345",
					},
				},
				accountInsertMap: map[string]*db.Account{
					"12345": {
						Address: "123456",
						BNS:     "1234",
						DID:     "did:32345",
					},
				},
				accountUpdateMap: map[string]*db.Account{
					"67890": {
						Address: "67890",
						BNS:     "1234",
						DID:     "did:32345",
					},
				},
				accountMap: map[string]*db.Account{
					"12345": {
						Address: "123456",
						BNS:     "1234",
						DID:     "did:32345",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processDIDAccounts(tt.args.didAccount, tt.args.unBindDIDs, tt.args.accountInsertMap, tt.args.accountUpdateMap, tt.args.accountMap)
		})
	}
}

func TestGetAccountType2(t *testing.T) {
	type args struct {
		chainId string
		address string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test case 1",
			args: args{
				chainId: "chain1",
				address: "123456789",
			},
			want: AddrTypeUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAccountType(tt.args.chainId, tt.args.address); got != tt.want {
				t.Errorf("GetAccountType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dealAccountNFTNum(t *testing.T) {
	type args struct {
		chainId          string
		minHeight        int64
		accountNFT       map[string]int64
		accountInsertMap map[string]*db.Account
		accountUpdateMap map[string]*db.Account
		accountMap       map[string]*db.Account
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test case 1",
			args: args{
				chainId: ChainId1,
				accountNFT: map[string]int64{
					"123456789": 12,
				},
				accountInsertMap: map[string]*db.Account{},
				accountUpdateMap: map[string]*db.Account{},
				accountMap:       map[string]*db.Account{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dealAccountNFTNum(tt.args.chainId, tt.args.minHeight, tt.args.accountNFT, tt.args.accountInsertMap, tt.args.accountUpdateMap, tt.args.accountMap)
		})
	}
}

func Test_dealAccountTxNum(t *testing.T) {
	type args struct {
		chainId          string
		minHeight        int64
		accountTx        map[string]int64
		accountInsertMap map[string]*db.Account
		accountUpdateMap map[string]*db.Account
		accountMap       map[string]*db.Account
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test case 1",
			args: args{
				chainId: ChainId1,
				accountTx: map[string]int64{
					"123456789": 12,
				},
				accountInsertMap: map[string]*db.Account{},
				accountUpdateMap: map[string]*db.Account{},
				accountMap:       map[string]*db.Account{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dealAccountTxNum(tt.args.chainId, tt.args.minHeight, tt.args.accountTx, tt.args.accountInsertMap, tt.args.accountUpdateMap, tt.args.accountMap)
		})
	}
}
