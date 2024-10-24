package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func insertGasRecordTest() ([]*db.GasRecord, error) {
	insertList := []*db.GasRecord{
		{
			TxId:         txId1,
			Address:      AccountAddr1,
			GasAmount:    1000,
			BusinessType: 0,
		},
	}
	err := InsertGasRecord(ChainID, insertList)
	return insertList, err
}

func TestGetGasRecordByTxIds(t *testing.T) {
	insertList, err := insertGasRecordTest()
	if err != nil {
		return
	}

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
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				txIds: []string{
					txId1,
					txId2,
				},
			},
			want:    insertList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGasRecordByTxIds(tt.args.chainId, tt.args.txIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGasRecordByTxIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.GasRecord{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetGasRecordByTxIds() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetGasRecordList(t *testing.T) {
	insertList, err := insertGasRecordTest()
	if err != nil {
		return
	}

	type args struct {
		offset       int
		limit        int
		chainId      string
		addrList     []string
		startTime    int64
		endTime      int64
		businessType int
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.GasRecord
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				offset:  0,
				limit:   10,
			},
			want:    insertList,
			want1:   int64(len(insertList)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetGasRecordList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.addrList, tt.args.startTime, tt.args.endTime, tt.args.businessType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGasRecordList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.GasRecord{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetGasRecordList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
			if got1 != tt.want1 {
				t.Errorf("GetGasRecordList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInsertGasRecord(t *testing.T) {
	insertList := []*db.GasRecord{
		{
			TxId:         txId1,
			Address:      AccountAddr1,
			GasAmount:    1000,
			BusinessType: 0,
		},
	}

	type args struct {
		chainId    string
		gasRecords []*db.GasRecord
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				gasRecords: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertGasRecord(tt.args.chainId, tt.args.gasRecords); (err != nil) != tt.wantErr {
				t.Errorf("InsertGasRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
