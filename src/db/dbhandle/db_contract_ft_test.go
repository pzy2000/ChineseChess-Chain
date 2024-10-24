package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	FTContractName1   = "FTContractName1"
	FTContractSymbol1 = "FTContractSymbol1"
	FTContractAddr1   = "123456"
)

func insertFungibleContractTest1() ([]*db.FungibleContract, error) {
	insertList := []*db.FungibleContract{
		{
			Symbol:          FTContractSymbol1,
			ContractName:    FTContractName1,
			ContractNameBak: FTContractName1,
			ContractAddr:    FTContractAddr1,
			//TotalSupply:     "12",
			HolderCount: 2,
			Timestamp:   123456,
		},
	}
	err := InsertFungibleContract(ChainID, insertList)
	return insertList, err
}

func TestGetFungibleContractByAddr(t *testing.T) {
	insertList, err := insertFungibleContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		contractAddr string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.FungibleContract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractAddr: FTContractAddr1,
			},
			want:    insertList[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFungibleContractByAddr(tt.args.chainId, tt.args.contractAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFungibleContractByAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.FungibleContract{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetFungibleContractByAddr() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetFungibleContractList(t *testing.T) {
	_, err := insertFungibleContractTest1()
	if err != nil {
		return
	}
	type args struct {
		offset      int
		limit       int
		chainId     string
		contractKey string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.FungibleContractWithTxNum
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				offset:  0,
				limit:   10,
				chainId: ChainID,
			},
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1, err := GetFungibleContractList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.contractKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFungibleContractList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 != tt.want1 {
				t.Errorf("GetFungibleContractList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInsertFungibleContract(t *testing.T) {
	insertList := []*db.FungibleContract{
		{
			Symbol:          FTContractSymbol1,
			ContractName:    FTContractName1,
			ContractNameBak: FTContractName1,
			ContractAddr:    FTContractAddr1,
			//TotalSupply:     "12",
			HolderCount: 2,
			Timestamp:   123456,
		},
	}

	type args struct {
		chainId   string
		contracts []*db.FungibleContract
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:   ChainID,
				contracts: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertFungibleContract(tt.args.chainId, tt.args.contracts); (err != nil) != tt.wantErr {
				t.Errorf("InsertFungibleContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueryFungibleContractExists(t *testing.T) {
	insertList, err := insertFungibleContractTest1()
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
		want    map[string]*db.FungibleContract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				addrList: []string{
					FTContractAddr1,
				},
			},
			want: map[string]*db.FungibleContract{
				insertList[0].ContractAddr: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryFungibleContractExists(tt.args.chainId, tt.args.addrList)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryFungibleContractExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.FungibleContract{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("QueryFungibleContractExists() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestUpdateFungibleContract(t *testing.T) {
	insertList, err := insertFungibleContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId  string
		contract *db.FungibleContract
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:  ChainID,
				contract: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateFungibleContract(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFungibleContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateFungibleContractName(t *testing.T) {
	_, err := insertFungibleContractTest1()
	if err != nil {
		return
	}

	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId  string
		contract *db.Contract
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:  ChainID,
				contract: contractInfo1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateFungibleContractName(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFungibleContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
