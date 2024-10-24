/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDealBackListEventData(t *testing.T) {
	eventData := []string{
		"chainID",
		"123456789",
		"223456789",
		"323456789",
	}
	want := []string{
		"123456789",
		"223456789",
		"323456789",
	}
	type args struct {
		contractName string
		topic        string
		eventData    []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 []string
	}{
		{
			name: "Test case 1",
			args: args{
				contractName: "TRANSACTION_MANAGER",
				topic:        "100",
				eventData:    eventData,
			},
			want:  want,
			want1: []string{},
		},
		{
			name: "Test case 1",
			args: args{
				contractName: "TRANSACTION_MANAGER",
				topic:        "101",
				eventData:    eventData,
			},
			want:  []string{},
			want1: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := DealBackListEventData(tt.args.contractName, tt.args.topic, tt.args.eventData)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealBackListEventData() got = %v, want %v", got, tt.want)
			}
			if !cmp.Equal(got1, tt.want1) {
				t.Errorf("DealBackListEventData() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDealDockerDFAEventData(t *testing.T) {
	eventData := []string{
		"123456789",
		"223456789",
		"323456789",
	}
	eventData1 := []string{
		"123456789",
		"223456789",
		"323456789",
		"323456789",
		"323456789",
	}
	want := &db.TransferTopicEventData{
		ToAddress: "123456789",
		Amount:    "223456789",
	}

	want1 := &db.TransferTopicEventData{
		FromAddress: "123456789",
		ToAddress:   "223456789",
		Amount:      "323456789",
	}

	type args struct {
		topic     string
		eventData []string
	}
	tests := []struct {
		name string
		args args
		want *db.TransferTopicEventData
	}{
		{
			name: "Test case 1",
			args: args{
				topic:     "101",
				eventData: eventData,
			},
			want: nil,
		},
		{
			name: "Test case 1",
			args: args{
				topic:     "mint",
				eventData: eventData,
			},
			want: want,
		},
		{
			name: "Test case 1",
			args: args{
				topic:     "transfer",
				eventData: eventData1,
			},
			want: want1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealDockerDFAEventData(tt.args.topic, tt.args.eventData)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealDockerDFAEventData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDealDockerNFAEventData(t *testing.T) {
	eventData := []string{
		"123456789",
		"223456789",
		"323456789",
	}
	want := &db.TransferTopicEventData{
		FromAddress: "123456789",
		ToAddress:   "223456789",
		TokenId:     "323456789",
	}

	type args struct {
		topic      string
		eventData  []string
		senderUser string
	}
	tests := []struct {
		name string
		args args
		want *db.TransferTopicEventData
	}{
		{
			name: "Test case 1",
			args: args{
				topic:      "transfer",
				eventData:  eventData,
				senderUser: "123456789",
			},
			want: nil,
		},
		{
			name: "Test case 1",
			args: args{
				topic:      "TransferFrom",
				eventData:  eventData,
				senderUser: "123456789",
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealDockerNFAEventData(tt.args.topic, tt.args.eventData, tt.args.senderUser)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealDockerNFAEventData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDealEVMDFAEventData(t *testing.T) {
	eventData := []string{
		"123456789",
		"223456789",
		"323456789",
	}
	eventData1 := []string{
		"000000000000000000000000123456789",
		"000000000000000000000000223456789",
		"000000000000000000000000323456789",
	}
	want := &db.TransferTopicEventData{
		FromAddress: "123456789",
		ToAddress:   "223456789",
		Amount:      "13476652937",
	}
	type args struct {
		topic     string
		eventData []string
	}
	tests := []struct {
		name string
		args args
		want *db.TransferTopicEventData
	}{
		{
			name: "Test case 1",
			args: args{
				topic:     "12344",
				eventData: eventData,
			},
			want: nil,
		},
		{
			name: "Test case 1",
			args: args{
				topic:     EVMEventTopicTransfer,
				eventData: eventData,
			},
			want: nil,
		},
		{
			name: "Test case 1",
			args: args{
				topic:     EVMEventTopicTransfer,
				eventData: eventData1,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealEVMDFAEventData(tt.args.topic, tt.args.eventData)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealEVMDFAEventData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDealEVMNFAEventData(t *testing.T) {
	eventData := []string{
		"123456789",
		"223456789",
		"323456789",
	}
	eventData1 := []string{
		"000000000000000000000000123456789",
		"000000000000000000000000223456789",
		"000000000000000000000000323456789",
	}
	want := &db.TransferTopicEventData{
		FromAddress: "123456789",
		ToAddress:   "223456789",
		TokenId:     "13476652937",
	}
	type args struct {
		topic     string
		eventData []string
	}
	tests := []struct {
		name string
		args args
		want *db.TransferTopicEventData
	}{
		{
			name: "Test case 1",
			args: args{
				topic:     "12344",
				eventData: eventData,
			},
			want: nil,
		},
		{
			name: "Test case 1",
			args: args{
				topic:     EVMEventTopicTransfer,
				eventData: eventData,
			},
			want: nil,
		},
		{
			name: "Test case 1",
			args: args{
				topic:     EVMEventTopicTransfer,
				eventData: eventData1,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealEVMNFAEventData(tt.args.topic, tt.args.eventData)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealEVMNFAEventData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDealIdentityEventData(t *testing.T) {
	eventData := []string{
		"123456789",
		"22",
		"323456789",
	}
	want := &db.IdentityEventData{
		UserAddr: "123456789",
		Level:    "22",
		PkPem:    "323456789",
	}
	type args struct {
		topic     string
		eventData []string
	}
	tests := []struct {
		name string
		args args
		want *db.IdentityEventData
	}{
		{
			name: "Test case 1",
			args: args{
				topic:     "121212",
				eventData: eventData,
			},
			want: nil,
		},
		{
			name: "Test case 1",
			args: args{
				topic:     "setIdentity",
				eventData: eventData,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealIdentityEventData(tt.args.topic, tt.args.eventData)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealIdentityEventData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDealUserBNSEventData(t *testing.T) {
	eventData := []string{
		"BNS:123456789.com",
		"123456789",
		"1",
	}
	want := &db.BNSTopicEventData{
		Domain:       "BNS:123456789.com",
		Value:        "123456789",
		ResourceType: "1",
	}
	type args struct {
		contractName string
		topic        string
		eventData    []string
	}
	tests := []struct {
		name  string
		args  args
		want  *db.BNSTopicEventData
		want1 string
	}{
		{
			name: "Test case 1",
			args: args{
				topic:        "121212",
				eventData:    eventData,
				contractName: "official_bns",
			},
			want:  nil,
			want1: "",
		},
		{
			name: "Test case 1",
			args: args{
				topic:        BNSBindEvent,
				eventData:    eventData,
				contractName: "official_bns",
			},
			want:  want,
			want1: "",
		},
		{
			name: "Test case 1",
			args: args{
				topic:        BNSUnBindEvent,
				eventData:    eventData,
				contractName: "official_bns",
			},
			want:  nil,
			want1: "BNS:123456789.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := DealUserBNSEventData(tt.args.contractName, tt.args.topic, tt.args.eventData)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealUserBNSEventData() = %v, want %v", got, tt.want)
			}
			if !cmp.Equal(got1, tt.want1) {
				t.Errorf("DealUserBNSEventData() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDealUserDIDEventData(t *testing.T) {
	didDocument := db.DidDocument{
		Id: "DID:123456789",
		VerificationMethod: []db.VerificationMethod{
			{
				Id:           "DID:123456789:Key1",
				Type:         "12121",
				Controller:   "2121212",
				PublicKeyPem: "121212",
				Address:      "123456789",
			},
			{
				Id:           "DID:123456789:Key2",
				Type:         "12121",
				Controller:   "2121212",
				PublicKeyPem: "121212",
				Address:      "223456789",
			},
		},
	}
	didDocumentJson, _ := json.Marshal(didDocument)
	eventData := []string{
		"DID:123456789",
		string(didDocumentJson),
	}
	want := &db.DIDTopicEventData{
		Did: "DID:123456789",
		VerificationMethod: []db.VerificationMethod{
			{
				Id:           "DID:123456789:Key1",
				Type:         "12121",
				Controller:   "2121212",
				PublicKeyPem: "121212",
				Address:      "123456789",
			},
			{
				Id:           "DID:123456789:Key2",
				Type:         "12121",
				Controller:   "2121212",
				PublicKeyPem: "121212",
				Address:      "223456789",
			},
		},
	}

	type args struct {
		contractName string
		topic        string
		eventData    []string
	}
	tests := []struct {
		name string
		args args
		want *db.DIDTopicEventData
	}{
		{
			name: "Test case 1",
			args: args{
				topic:        "121212",
				eventData:    eventData,
				contractName: "DID",
			},
			want: nil,
		},
		{
			name: "Test case 1",
			args: args{
				topic:        DIDSetDidDocument,
				eventData:    eventData,
				contractName: "DID",
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealUserDIDEventData(tt.args.contractName, tt.args.topic, tt.args.eventData)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealUserDIDEventData() = %v, want %v", got, tt.want)
			}
		})
	}
}
