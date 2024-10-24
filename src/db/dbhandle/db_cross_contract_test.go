package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

const (
	subChainId1      = "subChainId1"
	subContractName1 = "subContractName1"
)

func insertCrossContractTest1() ([]*db.CrossChainContract, error) {
	newUUID := uuid.New().String()
	contractList := []*db.CrossChainContract{
		{
			ID:           newUUID,
			SubChainId:   subChainId1,
			ContractName: subContractName1},
	}
	err := InsertCrossContract(ChainID, subChainId1, contractList)
	return contractList, err
}

func TestGetCrossContractByName(t *testing.T) {
	insertList, err := insertCrossContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId    string
		subChainId string
		nameList   []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossChainContract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				subChainId: subChainId1,
				nameList: []string{
					subContractName1,
					"123",
				},
			},
			want:    insertList,
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:    ChainID,
				subChainId: subChainId1,
			},
			want:    make([]*db.CrossChainContract, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossContractByName(tt.args.chainId, tt.args.subChainId, tt.args.nameList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossContractByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.CrossChainContract{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetCrossContractByName() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetCrossContractCount(t *testing.T) {
	insertList, err := insertCrossContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId    string
		subChainId string
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
				chainId:    ChainID,
				subChainId: subChainId1,
			},
			want:    int64(len(insertList)),
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:    ChainID,
				subChainId: "1234",
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossContractCount(tt.args.chainId, tt.args.subChainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossContractCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCrossContractCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossContractCountCache(t *testing.T) {
	insertList, err := insertCrossContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId    string
		subChainId string
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
				chainId:    ChainID,
				subChainId: subChainId1,
			},
			want:    int64(len(insertList)),
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:    ChainID,
				subChainId: "1234",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = GetCrossContractCount(tt.args.chainId, tt.args.subChainId)
			got, err := GetCrossContractCountCache(tt.args.chainId, tt.args.subChainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossContractCountCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCrossContractCountCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertCrossContract(t *testing.T) {
	newUUID := uuid.New().String()
	contractList := []*db.CrossChainContract{
		{
			ID:           newUUID,
			SubChainId:   subChainId1,
			ContractName: subContractName1},
	}

	type args struct {
		chainId    string
		subChainId string
		insertList []*db.CrossChainContract
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
				subChainId: subChainId1,
				insertList: contractList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertCrossContract(tt.args.chainId, tt.args.subChainId, tt.args.insertList); (err != nil) != tt.wantErr {
				t.Errorf("InsertCrossContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetCrossContractCountCache(t *testing.T) {
	type args struct {
		chainId       string
		subChainId    string
		contractCount int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:       ChainID,
				subChainId:    "123556",
				contractCount: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossContractCountCache(tt.args.chainId, tt.args.subChainId, tt.args.contractCount)
		})
	}
}
