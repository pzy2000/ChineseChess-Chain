/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
)

func TestBuildPositionList(t *testing.T) {
	contractEvents := getContractEventTest("0_contractEventData.json")
	contractInfoMap := getContractInfoMapTest("0_contractInfoMap.json")
	accountMap := getAccountMapTest("0_accountMap.json")
	positionList := getPositionListJsonTest("0_positionListJson.json")

	type args struct {
		contractEvents  []*db.ContractEventData
		contractInfoMap map[string]*db.Contract
		accountMap      map[string]*db.Account
	}
	tests := []struct {
		name string
		args args
		want map[string]*db.PositionData
	}{
		{
			name: "Test case 1",
			args: args{
				contractEvents:  contractEvents,
				contractInfoMap: contractInfoMap,
				accountMap:      accountMap,
			},
			want: positionList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildPositionList(tt.args.contractEvents, tt.args.contractInfoMap, tt.args.accountMap)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("BuildPositionList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestBuildUpdatePositionData(t *testing.T) {
	positionDataList, updateFungiblePosition, updateNonFungible := getUpdatePositionDataTest(
		"4_updatePositionData.json")
	positionDBMap := getPositionDBJsonTest("4_positionDBData.json")
	nonPositionDBMap := getNonPositionDBJsonTest("4_nonPositionDBData.json")

	wantResult := &db.BlockPosition{
		UpdateFungiblePosition: updateFungiblePosition,
		UpdateNonFungible:      updateNonFungible,
	}

	type args struct {
		minHeight        int64
		positionList     map[string]*db.PositionData
		positionDBMap    map[string][]*db.FungiblePosition
		nonPositionDBMap map[string][]*db.NonFungiblePosition
	}
	tests := []struct {
		name string
		args args
		want *db.BlockPosition
	}{
		{
			name: "Test case 1",
			args: args{
				minHeight:        20,
				positionList:     positionDataList,
				positionDBMap:    positionDBMap,
				nonPositionDBMap: nonPositionDBMap,
			},
			want: wantResult,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildUpdatePositionData(tt.args.minHeight, tt.args.positionList, tt.args.positionDBMap,
				tt.args.nonPositionDBMap)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("BuildUpdatePositionData() got = %v, want %v\ndiff: %s", got,
					tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestDealNonFungibleToken(t *testing.T) {
	contractEvents := getContractEventTest("5_contractEventDataToken.json")
	contractInfoMap := getContractInfoMapTest("0_contractInfoMap.json")
	accountMap := getAccountMapTest("5_accountMapToken.json")
	gotTokenResult := getGotTokenResultTest("5_gotTokenResult.json")
	type args struct {
		contractEvents  []*db.ContractEventData
		contractInfoMap map[string]*db.Contract
		accountMap      map[string]*db.Account
	}
	tests := []struct {
		name string
		args args
		want *db.TokenResult
	}{
		{
			name: "Test case 1",
			args: args{
				contractEvents:  contractEvents,
				contractInfoMap: contractInfoMap,
				accountMap:      accountMap,
			},
			want: gotTokenResult,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealNonFungibleToken(ChainId, tt.args.contractEvents, tt.args.contractInfoMap, tt.args.accountMap)
			if len(got.InsertUpdateToken) != 2 {
				t.Errorf("DealNonFungibleToken got= %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestPositionToAddressData(t *testing.T) {
	type args struct {
		positionList map[string]*db.PositionData
		address      string
		amount       decimal.Decimal
		accountMap   map[string]*db.Account
		contract     *db.Contract
	}
	tests := []struct {
		name string
		args args
		want map[string]*db.PositionData
	}{
		{
			name: "测试 PositionToAddressData",
			args: args{
				positionList: make(map[string]*db.PositionData),
				address:      "0x123",
				amount:       decimal.NewFromInt(10),
				accountMap: map[string]*db.Account{
					"0x123": {
						AddrType: AddrTypeUser,
					},
				},
				contract: &db.Contract{
					Addr:           "0x456",
					Name:           "TestContract",
					ContractSymbol: "TEST",
					Decimals:       18,
				},
			},
			want: map[string]*db.PositionData{
				"0x123_0x456": {
					AddrType:     AddrTypeUser,
					OwnerAddr:    "0x123",
					ContractAddr: "0x456",
					ContractName: "TestContract",
					Symbol:       "TEST",
					Amount:       decimal.NewFromInt(10),
					Decimals:     18,
				},
			},
		},
		// 添加其他测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PositionToAddressData(tt.args.positionList, tt.args.address, tt.args.amount, tt.args.accountMap, tt.args.contract)
			if !cmp.Equal(tt.args.positionList, tt.want) {
				t.Errorf("PositionToAddressData() = %v, want %v", tt.args.positionList, tt.want)
			}
		})
	}
}

func TestPositionFromAddressData(t *testing.T) {
	type args struct {
		positionList map[string]*db.PositionData
		address      string
		amount       decimal.Decimal
		accountMap   map[string]*db.Account
		contract     *db.Contract
	}
	tests := []struct {
		name string
		args args
		want map[string]*db.PositionData
	}{
		{
			name: "测试 PositionFromAddressData",
			args: args{
				positionList: make(map[string]*db.PositionData),
				address:      "0x123",
				amount:       decimal.NewFromInt(10),
				accountMap: map[string]*db.Account{
					"0x123": {
						AddrType: AddrTypeUser,
					},
				},
				contract: &db.Contract{
					Addr:           "0x456",
					Name:           "TestContract",
					ContractSymbol: "TEST",
					Decimals:       18,
				},
			},
			want: map[string]*db.PositionData{
				"0x123_0x456": {
					AddrType:     AddrTypeUser,
					OwnerAddr:    "0x123",
					ContractAddr: "0x456",
					ContractName: "TestContract",
					Symbol:       "TEST",
					Amount:       decimal.NewFromInt(-10),
					Decimals:     18,
				},
			},
		},
		// 添加其他测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PositionFromAddressData(tt.args.positionList, tt.args.address, tt.args.amount, tt.args.accountMap, tt.args.contract)
			if !cmp.Equal(tt.args.positionList, tt.want) {
				t.Errorf("PositionFromAddressData() = %v, want %v", tt.args.positionList, tt.want)
			}
		})
	}
}

func Test_dealFungiblePosition(t *testing.T) {
	positionDBMap := getPositionDBJsonTest("4_positionDBData.json")

	amountDecimal := decimal.NewFromInt(10000000)
	position := &db.PositionData{
		AddrType:     0,
		OwnerAddr:    "171262347a59fded92021a32421a5dad05424e03",
		ContractAddr: "ea9c7f588e2bce761ae33ac5bf31092abefb1aae",
		ContractName: "EVM_ERC202",
		Amount:       amountDecimal,
		Decimals:     0,
		ContractType: "EVMDFA",
	}

	want := []*db.FungiblePosition{
		{
			//AddrType:     0,
			OwnerAddr:    "171262347a59fded92021a32421a5dad05424e03",
			ContractAddr: "ea9c7f588e2bce761ae33ac5bf31092abefb1aae",
			ContractName: "EVM_ERC202",
			Amount:       GetNumberDecimal("10000000"),
			BlockHeight:  34,
		},
	}
	type args struct {
		minHeight     int64
		position      *db.PositionData
		positionDBMap map[string][]*db.FungiblePosition
		positionOp    *db.BlockPosition
	}
	tests := []struct {
		name string
		args args
		want []*db.FungiblePosition
	}{
		{
			name: "Test case 1",
			args: args{
				minHeight:     34,
				position:      position,
				positionDBMap: positionDBMap,
				positionOp:    &db.BlockPosition{},
			},
			want: want,
		},
		{
			name: "Test case 2",
			args: args{
				minHeight:     34,
				position:      position,
				positionDBMap: positionDBMap,
				positionOp:    &db.BlockPosition{},
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dealFungiblePosition(tt.args.minHeight, tt.args.position, tt.args.positionDBMap, tt.args.positionOp)
			positionRes := tt.args.positionOp.InsertFungiblePosition
			if !cmp.Equal(positionRes, tt.want, cmpopts.IgnoreFields(db.FungiblePosition{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("dealFungiblePosition() got = %v, want %v\ndiff: %s", positionRes, tt.want, cmp.Diff(positionRes, tt.want))
			}
		})
	}
}

func Test_dealNonFungiblePosition(t *testing.T) {
	amountDecimal := decimal.NewFromInt(10)
	positionDBMap := getNonPositionDBJsonTest("4_nonPositionDBData.json")
	position := &db.PositionData{
		AddrType:     0,
		OwnerAddr:    "18fc4e7429af8419d5bb307e34db398b9a2331c6",
		ContractAddr: "aba31ce4cd49f08073d2f115eb12610544242456",
		ContractName: "goErc721_1",
		Amount:       amountDecimal,
		Decimals:     0,
		ContractType: "EVMDFA",
	}
	want := []*db.NonFungiblePosition{
		{
			//AddrType:     0,
			OwnerAddr:    "18fc4e7429af8419d5bb307e34db398b9a2331c6",
			ContractAddr: "aba31ce4cd49f08073d2f115eb12610544242456",
			ContractName: "goErc721_1",
			Amount:       GetNumberDecimal("44"),
			BlockHeight:  34,
		},
	}

	type args struct {
		minHeight     int64
		position      *db.PositionData
		positionDBMap map[string][]*db.NonFungiblePosition
		positionOp    *db.BlockPosition
	}
	tests := []struct {
		name string
		args args
		want []*db.NonFungiblePosition
	}{
		{
			name: "Test case 1",
			args: args{
				minHeight:     34,
				position:      position,
				positionDBMap: positionDBMap,
				positionOp:    &db.BlockPosition{},
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dealNonFungiblePosition(tt.args.minHeight, tt.args.position, tt.args.positionDBMap, tt.args.positionOp)
			positionRes := tt.args.positionOp.UpdateNonFungible
			if !cmp.Equal(positionRes, tt.want) {
				t.Errorf("dealNonFungiblePosition got= %v, want %v\ndiff: %s", positionRes, tt.want,
					cmp.Diff(positionRes, tt.want))
			}
		})
	}
}
