/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"encoding/json"
	"reflect"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/google/go-cmp/cmp"
)

var TxInfoEventJson = "{\"payload\":{\"chain_id\":\"chain1\",\"tx_id\":\"17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5\",\"timestamp\":1702461023,\"contract_name\":\"goErc20_1\",\"method\":\"Mint\",\"parameters\":[{\"key\":\"account\",\"value\":\"MThmYzRlNzQyOWFmODQxOWQ1YmIzMDdlMzRkYjM5OGI5YTIzMzFjNg==\"},{\"key\":\"amount\",\"value\":\"MTAwMDAwMDAwMDA=\"}],\"limit\":{\"gas_limit\":13000}},\"sender\":{\"signer\":{\"org_id\":\"wx-org1.chainmaker.org\",\"member_type\":1,\"member_info\":\"LK4/KplYsQcFU2All0UxorspVALdt/tgHuZ8QxiME2M=\"},\"signature\":\"MEQCIAK4XuZoU0XB+ya2PRNtebY/BACX8BQOBRMQBMidbfXpAiA6EURjarfxU/qbwCrkqptOdXav7orDeVR38aHXzynb5g==\"},\"result\":{\"contract_result\":{\"result\":\"b2s=\",\"message\":\"Success\",\"gas_used\":156,\"contract_event\":[{\"topic\":\"mint\",\"tx_id\":\"17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5\",\"contract_name\":\"goErc20_1\",\"event_data\":[\"18fc4e7429af8419d5bb307e34db398b9a2331c6\",\"10000000000\"]}]},\"rw_set_hash\":\"BOf54ycn6MjSiUL06tEU8WN2cDvTXJWkShVNowYYAK4=\"}}"
var ContractEventsJson = "[{\"txId\":\"17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5\",\"eventIndex\":1,\"topic\":\"mint\",\"topicBak\":\"\",\"contractName\":\"goErc20_1\",\"contractNameBak\":\"goErc20_1\",\"contractAddr\":\"\",\"contractVersion\":\"\",\"eventData\":[\"18fc4e7429af8419d5bb307e34db398b9a2331c6\",\"10000000000\"],\"eventDataBak\":\"\",\"timestamp\":1702461023,\"createdAt\":\"0001-01-01T00:00:00Z\",\"updatedAt\":\"0001-01-01T00:00:00Z\"}]"

var ContractInfoMap = map[string]*db.Contract{
	"aba31ce4cd49f08073d2f115eb12610544242ff9": {
		Name:         "goErc20_1",
		NameBak:      "goErc20_1",
		Addr:         "aba31ce4cd49f08073d2f115eb12610544242ff9",
		ContractType: "CMDFA",
		TxNum:        100,
	},
	"goErc20_1": {
		Name:         "goErc20_1",
		NameBak:      "goErc20_1",
		Addr:         "aba31ce4cd49f08073d2f115eb12610544242ff9",
		ContractType: "CMDFA",
		TxNum:        100,
	},
}

//var txInfoList = []*db.Transaction{
//	{
//		TxId:               "17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5",
//		Sender:             "client1.sign.wx-org1.chainmaker.org",
//		SenderOrgId:        "wx-org1.chainmaker.org",
//		BlockHeight:        40,
//		BlockHash:          "d3b2b488033c2faa100949667572b1875d82f7a32bd35bccf8232f5d3eef6545",
//		TxType:             "INVOKE_CONTRACT",
//		Timestamp:          1702461023,
//		TxIndex:            1,
//		TxStatusCode:       "SUCCESS",
//		RwSetHash:          "04e7f9e32727e8c8d28942f4ead114f16376703bd35c95a44a154da3061800ae",
//		ContractResultCode: 0,
//		ContractName:       "goErc20_1",
//		ContractNameBak:    "goErc20_1",
//		ContractAddr:       "aba31ce4cd49f08073d2f115eb12610544242ff9",
//		ContractType:       "CMDFA",
//		UserAddr:           "171262347a59fded92021a32421a5dad05424e03",
//	},
//}

func TestDealContractEvents(t *testing.T) {
	transactionInfo := &common.Transaction{}
	err := json.Unmarshal([]byte(txInfoJson), transactionInfo)
	if err != nil {
		return
	}

	transactionEvent := &common.Transaction{}
	err = json.Unmarshal([]byte(TxInfoEventJson), transactionEvent)
	if err != nil {
		return
	}

	contractEvents := make([]*db.ContractEvent, 0)
	err = json.Unmarshal([]byte(ContractEventsJson), &contractEvents)
	if err != nil {
		return
	}
	type args struct {
		txInfo *common.Transaction
	}
	tests := []struct {
		name string
		args args
		want []*db.ContractEvent
	}{
		{
			name: "Test case 1",
			args: args{
				txInfo: transactionInfo,
			},
			want: []*db.ContractEvent{},
		},
		{
			name: "Test case 1",
			args: args{
				txInfo: transactionEvent,
			},
			want: contractEvents,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealContractEvents(tt.args.txInfo)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealContractEvents() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestDealContractTxNum(t *testing.T) {
	type args struct {
		minHeight     int64
		contractMap   map[string]*db.Contract
		txList        map[string]*db.Transaction
		contractEvent []*db.ContractEvent
	}
	txInfoListMap := map[string]*db.Transaction{}
	txInfoListMap["17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5"] = &db.Transaction{
		TxId:               "17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5",
		Sender:             "client1.sign.wx-org1.chainmaker.org",
		SenderOrgId:        "wx-org1.chainmaker.org",
		BlockHeight:        40,
		BlockHash:          "d3b2b488033c2faa100949667572b1875d82f7a32bd35bccf8232f5d3eef6545",
		TxType:             "INVOKE_CONTRACT",
		Timestamp:          1702461023,
		TxIndex:            1,
		TxStatusCode:       "SUCCESS",
		RwSetHash:          "04e7f9e32727e8c8d28942f4ead114f16376703bd35c95a44a154da3061800ae",
		ContractResultCode: 0,
		ContractName:       "goErc20_1",
		ContractNameBak:    "goErc20_1",
		ContractAddr:       "aba31ce4cd49f08073d2f115eb12610544242ff9",
		ContractType:       "CMDFA",
		UserAddr:           "171262347a59fded92021a32421a5dad05424e03",
	}

	tests := []struct {
		name      string
		args      args
		wantTxNum int64
	}{
		{
			name: "Test case 1",
			args: args{
				minHeight:   1000,
				contractMap: ContractInfoMap,
				txList:      txInfoListMap,
			},
			wantTxNum: 101,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UpdateContractTxAndEventNum(tt.args.minHeight, tt.args.contractMap, tt.args.txList, tt.args.contractEvent)
			if !reflect.DeepEqual(got[0].TxNum, tt.wantTxNum) {
				t.Errorf("DealContractTxNum() = %v, want %v", got[0], tt.wantTxNum)
			}
		})
	}
}

func TestDealTopicEventData(t *testing.T) {
	contractEvents := make([]*db.ContractEvent, 0)
	err := json.Unmarshal([]byte(ContractEventsJson), &contractEvents)
	if err != nil {
		return
	}
	transactionEvent := &common.Transaction{}
	err = json.Unmarshal([]byte(TxInfoEventJson), transactionEvent)
	if err != nil {
		return
	}
	txInfoMap := map[string]*db.Transaction{
		"17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5": {
			TxId:               "17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5",
			Sender:             "client1.sign.wx-org1.chainmaker.org",
			SenderOrgId:        "wx-org1.chainmaker.org",
			BlockHeight:        40,
			BlockHash:          "d3b2b488033c2faa100949667572b1875d82f7a32bd35bccf8232f5d3eef6545",
			TxType:             "INVOKE_CONTRACT",
			Timestamp:          1702461023,
			TxIndex:            1,
			TxStatusCode:       "SUCCESS",
			RwSetHash:          "04e7f9e32727e8c8d28942f4ead114f16376703bd35c95a44a154da3061800ae",
			ContractResultCode: 0,
			ContractName:       "goErc20_1",
			ContractNameBak:    "goErc20_1",
			ContractAddr:       "aba31ce4cd49f08073d2f115eb12610544242ff9",
			ContractType:       "CMDFA",
			UserAddr:           "171262347a59fded92021a32421a5dad05424e03",
		},
	}
	contractInfoMap := map[string]*db.Contract{
		"goErc20_1": {
			Name:         "goErc20_1",
			NameBak:      "goErc20_1",
			Addr:         "aba31ce4cd49f08073d2f115eb12610544242ff9",
			ContractType: "CMDFA",
		},
	}

	type args struct {
		contractEvent   []*db.ContractEvent
		contractInfoMap map[string]*db.Contract
		txInfoMap       map[string]*db.Transaction
	}
	tests := []struct {
		name string
		args args
		want *TopicEventResult
	}{
		{
			name: "Test case 1",
			args: args{
				contractEvent:   contractEvents,
				contractInfoMap: contractInfoMap,
				txInfoMap:       txInfoMap,
			},
			want: &TopicEventResult{
				AddBlack:         []string{},
				DeleteBlack:      []string{},
				IdentityContract: []*db.IdentityContract{},
				ContractEventData: []*db.ContractEventData{
					{
						Index:        1,
						Topic:        "mint",
						TxId:         "17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5",
						ContractName: "goErc20_1",
						EventData: &db.TransferTopicEventData{
							FromAddress: "",
							ToAddress:   "18fc4e7429af8419d5bb307e34db398b9a2331c6",
							Amount:      "10000000000",
						},
						Timestamp: 1702461023,
					},
				},
				OwnerAdders: []string{
					"18fc4e7429af8419d5bb307e34db398b9a2331c6",
				},
				DIDAccount:       map[string][]string{},
				BNSBindEventData: []*db.BNSTopicEventData{},
				BNSUnBindDomain:  []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealTopicEventData(tt.args.contractEvent, tt.args.contractInfoMap, tt.args.txInfoMap)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealTopicEventData got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestDealEvidence(t *testing.T) {
	txInfo := getTxInfoInfoTest("3_txInfoJson_evidenceList.json")
	txInfo1 := getTxInfoInfoTest("3_txInfoJson_evidenceList.json")
	userInfo := getUserInfoInfoTest("3_userInfoJson_evidence.json")
	if txInfo == nil || userInfo == nil {
		return
	}

	txInfo.Payload.Method = PayloadMethodEvidence
	txInfo1.Payload.Method = PayloadMethodEvidenceBatch
	type args struct {
		blockHeight int64
		txInfo      *common.Transaction
		userInfo    *MemberAddrIdCert
	}
	tests := []struct {
		name          string
		args          args
		wantEvidences []*db.EvidenceContract
		wantErr       bool
	}{
		{
			name: "Test case 1",
			args: args{
				blockHeight: 2,
				txInfo:      txInfo,
				userInfo:    userInfo,
			},
			wantErr: false,
		},
		{
			name: "Test case 1",
			args: args{
				blockHeight: 2,
				txInfo:      txInfo,
				userInfo:    userInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEvidences, err := DealEvidence(tt.args.blockHeight, tt.args.txInfo, tt.args.userInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("DealEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotEvidences) == 0 {
				t.Errorf("DealEvidence() gotEvidences =  %v", gotEvidences)
			}
		})
	}
}

func TestGetContractEvents(t *testing.T) {
	type args struct {
		chainId string
		txIds   []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.ContractEvent
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				chainId: ChainId1,
				txIds: []string{
					"1234",
					"12344",
				},
			},
			want:    make([]*db.ContractEvent, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractEvents(tt.args.chainId, tt.args.txIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetContractEvents() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestBuildTransferEventData(t *testing.T) {
	type args struct {
		topicEventResult *TopicEventResult
		ownerAddrMap     map[string]string
		contractInfoMap  map[string]*db.Contract
		event            *db.ContractEvent
		senderUser       string
		eventData        []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Test case 1",
			args: args{
				topicEventResult: &TopicEventResult{},
				ownerAddrMap:     map[string]string{},
				contractInfoMap: map[string]*db.Contract{
					"ContractName": {
						ContractType: "CMDFA",
					},
				},
				event: &db.ContractEvent{
					TxId:            "1231212313",
					Topic:           "mint",
					ContractName:    "ContractName",
					ContractNameBak: "ContractName",
					ContractAddr:    "1234",
				},
				eventData: []string{
					"12345",
					"123",
				},
			},
			want: map[string]string{
				"12345": "12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildTransferEventData(tt.args.topicEventResult, tt.args.ownerAddrMap, tt.args.contractInfoMap, tt.args.event, tt.args.senderUser, tt.args.eventData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildTransferEventData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildIdentityEventData(t *testing.T) {
	eventInfo := &db.ContractEvent{
		TxId:            "1231212313",
		Topic:           "setIdentity",
		ContractName:    "ContractName",
		ContractNameBak: "ContractName",
		ContractAddr:    "1234",
	}
	type args struct {
		topicEventResult *TopicEventResult
		contractInfoMap  map[string]*db.Contract
		event            *db.ContractEvent
		eventData        []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1",
			args: args{
				topicEventResult: &TopicEventResult{},
				contractInfoMap: map[string]*db.Contract{
					"ContractName": {
						ContractType: "CMID",
					},
				},
				event: eventInfo,
				eventData: []string{
					"12345666",
					"12345777",
					"123455",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BuildIdentityEventData(tt.args.topicEventResult, tt.args.contractInfoMap, tt.args.event, tt.args.eventData)
		})
	}
}
