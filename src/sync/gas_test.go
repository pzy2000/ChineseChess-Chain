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

	"chainmaker.org/chainmaker/pb-go/v2/syscontract"

	"github.com/google/go-cmp/cmp/cmpopts"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/google/go-cmp/cmp"
)

var txInfoJson = "{\"payload\":{\"chain_id\":\"chain1\",\"tx_id\":\"17a059bf5f096c6bca7e84373d3eb43e4265ac491b7947fd93db8c267879827d\",\"timestamp\":1702459337,\"contract_name\":\"ACCOUNT_MANAGER\",\"method\":\"RECHARGE_GAS\",\"parameters\":[{\"key\":\"batch_recharge\",\"value\":\"CjMKKDE3MTI2MjM0N2E1OWZkZWQ5MjAyMWEzMjQyMWE1ZGFkMDU0MjRlMDMQgICE/qbe4RE=\"}]},\"sender\":{\"signer\":{\"org_id\":\"wx-org1.chainmaker.org\",\"member_info\":\"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNlRENDQWg2Z0F3SUJBZ0lERFdIZk1Bb0dDQ3FHU000OUJBTUNNSUdLTVFzd0NRWURWUVFHRXdKRFRqRVEKTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVmTUIwR0ExVUVDaE1XZDNndApiM0puTVM1amFHRnBibTFoYTJWeUxtOXlaekVTTUJBR0ExVUVDeE1KY205dmRDMWpaWEowTVNJd0lBWURWUVFECkV4bGpZUzUzZUMxdmNtY3hMbU5vWVdsdWJXRnJaWEl1YjNKbk1CNFhEVEl6TVRJd01UQTRORE14TkZvWERUSTQKTVRFeU9UQTRORE14TkZvd2daRXhDekFKQmdOVkJBWVRBa05PTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBdwpEZ1lEVlFRSEV3ZENaV2xxYVc1bk1SOHdIUVlEVlFRS0V4WjNlQzF2Y21jeExtTm9ZV2x1YldGclpYSXViM0puCk1ROHdEUVlEVlFRTEV3WmpiR2xsYm5ReExEQXFCZ05WQkFNVEkyTnNhV1Z1ZERFdWMybG5iaTUzZUMxdmNtY3gKTG1Ob1lXbHViV0ZyWlhJdWIzSm5NRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUVadTcwb0xRWQp2UEptNnZlUFdMbENWeTJHOHpqYUQxL2tpYUpuMnNyRnc3WVR2cDNYV2d5OVJVM1ZRNnJsa004VExaYWY5Z2NRCmdScWFoSEJnaTZ0T1FLTnFNR2d3RGdZRFZSMFBBUUgvQkFRREFnYkFNQ2tHQTFVZERnUWlCQ0FvNGxwakRiNXAKSmdCc2JBc3U5aXEwQlM1V3p3N0IvMy9kelM0anpadEdTakFyQmdOVkhTTUVKREFpZ0NBUFJxKy8xd1FQajhBawplVkl5bDhENmkwZGdxdnh5NWV1QytERjVXVnVVTnpBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlBakdlZ0pndWQ1CnZPU0plVktENzdyUzFwOWE5TytQQU1UM3ptbWd6MlJZWndJaEFPNDE4Z3V2NUlhckFJMmt1MXlGbTVQK2FmYWQKeW1lNnp2c1RVbEdhOHhLZgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\"},\"signature\":\"MEYCIQC1/uFuWZ5FiJyNgMcJcoIPLKuIx32MBrvrvBpxklLmewIhAP0wKfeKRV7QUZhsH889vHVwTmviHKm9YhLF8e8fODdl\"},\"result\":{\"contract_result\":{\"result\":\"c3VjY2Vzcw==\",\"message\":\"OK\"},\"rw_set_hash\":\"USmT6IWvaESBjvsiR5nE8InKSJ6e53AfukSQwgj9/wE=\"}}"
var userResultJson = "{\"SenderUserId\":\"client1.sign.wx-org1.chainmaker.org\",\"SenderUserAddr\":\"171262347a59fded92021a32421a5dad05424e03\",\"SenderOrgId\":\"wx-org1.chainmaker.org\",\"SenderRole\":\"client\",\"PayerUserId\":\"\",\"PayerUserAddr\":\"\"}"
var gasRecordsJson = "[{\"gasIndex\":1,\"txId\":\"17a059bf5f096c6bca7e84373d3eb43e4265ac491b7947fd93db8c267879827d\",\"address\":\"171262347a59fded92021a32421a5dad05424e03\",\"payerAddress\":\"\",\"gasAmount\":10000000000000000,\"businessType\":1,\"timestamp\":1702459337,\"createdAt\":\"0001-01-01T00:00:00Z\",\"updatedAt\":\"0001-01-01T00:00:00Z\"}]"

func TestGetGasRecord(t *testing.T) {
	type args struct {
		chainId string
		txIds   []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.GasRecord
		wantErr bool
	}{
		{
			name: "Test case 1: Valid blockInfo and hashType",
			args: args{
				chainId: "chain1",
				txIds:   []string{"12345", "45456565"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetGasRecord(tt.args.chainId, tt.args.txIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGasRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_buildGasInfo(t *testing.T) {
	gasRecords := []*db.GasRecord{
		{
			GasIndex:     0,
			TxId:         "1234566777",
			Address:      "123456789",
			BusinessType: 1,
			GasAmount:    10000,
		},
		{
			GasIndex:     1,
			TxId:         "1234566777",
			Address:      "123456789",
			BusinessType: 2,
			GasAmount:    2000,
		},
		{
			GasIndex:     0,
			TxId:         "1234566777",
			Address:      "1234567890",
			BusinessType: 1,
			GasAmount:    2000,
		},
	}
	gasInfoList := []*db.Gas{
		{
			BlockHeight: 10,
			Address:     "123456789",
			GasBalance:  10000,
			GasTotal:    20000,
			GasUsed:     10000,
		},
	}

	want := []*db.Gas{
		{
			BlockHeight: 12,
			Address:     "1234567890",
			GasBalance:  2000,
			GasTotal:    2000,
			GasUsed:     0,
		},
	}
	want1 := []*db.Gas{
		{
			BlockHeight: 12,
			Address:     "123456789",
			GasBalance:  18000,
			GasTotal:    30000,
			GasUsed:     12000,
		},
	}

	type args struct {
		gasRecords  []*db.GasRecord
		gasInfoList []*db.Gas
		minHeight   int64
	}
	tests := []struct {
		name  string
		args  args
		want  []*db.Gas
		want1 []*db.Gas
	}{
		{
			name: "Test case 1: Valid blockInfo and hashType",
			args: args{
				gasRecords:  gasRecords,
				gasInfoList: gasInfoList,
				minHeight:   12,
			},
			want:  want,
			want1: want1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := buildGasInfo(tt.args.gasRecords, tt.args.gasInfoList, tt.args.minHeight)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("buildGasInfo() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
			if !cmp.Equal(got1, tt.want1) {
				t.Errorf("buildGasInfo() got1 = %v, want %v\ndiff: %s", got1, tt.want1, cmp.Diff(got1, tt.want1))
			}
		})
	}
}

func Test_buildGasRecord(t *testing.T) {
	transactionInfo := &common.Transaction{}
	transactionInfo1 := &common.Transaction{}
	err := json.Unmarshal([]byte(txInfoJson), transactionInfo)
	err = json.Unmarshal([]byte(txInfoJson), transactionInfo1)
	if err != nil {
		return
	}
	userResult := &db.SenderPayerUser{}
	err = json.Unmarshal([]byte(userResultJson), userResult)
	if err != nil {
		return
	}
	gasRecords := make([]*db.GasRecord, 0)
	err = json.Unmarshal([]byte(gasRecordsJson), &gasRecords)
	if err != nil {
		return
	}

	transactionInfo1.Payload.Method = syscontract.GasAccountFunction_RECHARGE_GAS.String()
	transactionInfo1.Result.Code = common.TxStatusCode_SUCCESS
	type args struct {
		txInfo     *common.Transaction
		userResult *db.SenderPayerUser
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.GasRecord
		wantErr bool
	}{
		{
			name: "Test case 1: Valid blockInfo and hashType",
			args: args{
				txInfo:     transactionInfo,
				userResult: userResult,
			},
			want: gasRecords,
		},
		{
			name: "Test case 1: Valid blockInfo and hashType",
			args: args{
				txInfo:     transactionInfo1,
				userResult: userResult,
			},
			want: gasRecords,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildGasRecord(tt.args.txInfo, tt.args.userResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildGasRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.GasRecord{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("buildGasRecord() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}
