package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func insertGasTest() ([]*db.Gas, error) {
	insertList := []*db.Gas{
		{
			Address:    ownerAddr1,
			GasBalance: 100,
			GasTotal:   1000,
			GasUsed:    900,
		},
	}
	err := InsertBatchGas(ChainID, insertList)
	return insertList, err
}

func TestGetGasByAddrInfo(t *testing.T) {
	_, err := insertGasTest()
	if err != nil {
		return
	}

	type args struct {
		chainId  string
		addrList []string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				addrList: []string{
					ownerAddr1,
				},
			},
			want:    100,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGasByAddrInfo(tt.args.chainId, tt.args.addrList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGasByAddrInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGasByAddrInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGasInfoByAddr(t *testing.T) {
	insertList, err := insertGasTest()
	if err != nil {
		return
	}
	type args struct {
		chainId  string
		addrList []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Gas
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				addrList: []string{
					ownerAddr1,
				},
			},
			want:    insertList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGasInfoByAddr(tt.args.chainId, tt.args.addrList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGasInfoByAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Gas{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetGasInfoByAddr() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetGasList(t *testing.T) {
	insertList, err := insertGasTest()
	if err != nil {
		return
	}

	type args struct {
		offset   int
		limit    int
		chainId  string
		addrList []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Gas
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				offset:  0,
				limit:   10,
				addrList: []string{
					ownerAddr1,
				},
			},
			want:    insertList,
			want1:   int64(len(insertList)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetGasList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.addrList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGasList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Gas{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetGasList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
			if got1 != tt.want1 {
				t.Errorf("GetGasList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInsertBatchGas(t *testing.T) {
	insertList := []*db.Gas{
		{
			Address:    ownerAddr1,
			GasBalance: 100,
			GasTotal:   1000,
			GasUsed:    900,
		},
	}

	type args struct {
		chainId string
		gasList []*db.Gas
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				gasList: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertBatchGas(tt.args.chainId, tt.args.gasList); (err != nil) != tt.wantErr {
				t.Errorf("InsertBatchGas() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateGas(t *testing.T) {
	insertList, err := insertGasTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		gasInfo *db.Gas
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				gasInfo: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateGas(tt.args.chainId, tt.args.gasInfo); (err != nil) != tt.wantErr {
				t.Errorf("UpdateGas() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
