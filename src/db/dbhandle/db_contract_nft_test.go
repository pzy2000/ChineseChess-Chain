package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	NFTContractName1 = "NFTContractName1"
	NFTContractAddr1 = "123456789"
)

func insertNonFungibleContractTest1() ([]*db.NonFungibleContract, error) {
	insertList := []*db.NonFungibleContract{
		{
			ContractName:    NFTContractName1,
			ContractNameBak: NFTContractName1,
			ContractAddr:    NFTContractAddr1,
			//TotalSupply:     "12",
			HolderCount: 2,
			Timestamp:   123456,
		},
	}
	err := InsertNonFungibleContract(ChainID, insertList)
	return insertList, err
}

func TestGetNonFungibleContractByAddr(t *testing.T) {
	insertList, err := insertNonFungibleContractTest1()
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
		want    *db.NonFungibleContract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractAddr: NFTContractAddr1,
			},
			want:    insertList[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNonFungibleContractByAddr(tt.args.chainId, tt.args.contractAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNonFungibleContractByAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.NonFungibleContract{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetFungibleContractByAddr() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetNonFungibleContractList(t *testing.T) {
	_, err := insertNonFungibleContractTest1()
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
		want    []*db.NonFungibleContractWithTxNum
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
			_, got1, err := GetNonFungibleContractList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.contractKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNonFungibleContractList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 != tt.want1 {
				t.Errorf("GetNonFungibleContractList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInsertNonFungibleContract(t *testing.T) {
	insertList := []*db.NonFungibleContract{
		{
			ContractName:    NFTContractName1,
			ContractNameBak: NFTContractName1,
			ContractAddr:    NFTContractAddr1,
			//TotalSupply:     "12",
			HolderCount: 2,
			Timestamp:   123456,
		},
	}

	type args struct {
		chainId   string
		contracts []*db.NonFungibleContract
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
			if err := InsertNonFungibleContract(tt.args.chainId, tt.args.contracts); (err != nil) != tt.wantErr {
				t.Errorf("InsertNonFungibleContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueryNonFungibleContractExists(t *testing.T) {
	insertList, err := insertNonFungibleContractTest1()
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
		want    map[string]*db.NonFungibleContract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				addrList: []string{
					NFTContractAddr1,
				},
			},
			want: map[string]*db.NonFungibleContract{
				insertList[0].ContractAddr: insertList[0],
			},
			wantErr: false,
		},
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				addrList: []string{
					"123",
				},
			},
			want:    map[string]*db.NonFungibleContract{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryNonFungibleContractExists(tt.args.chainId, tt.args.addrList)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryNonFungibleContractExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.NonFungibleContract{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("QueryNonFungibleContractExists() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestUpdateNonFungibleContract(t *testing.T) {
	insertList, err := insertNonFungibleContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId  string
		contract *db.NonFungibleContract
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
			if err := UpdateNonFungibleContract(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateNonFungibleContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateNonFungibleContractName(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}

	_, err = insertNonFungibleContractTest1()
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
			if err := UpdateNonFungibleContractName(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateNonFungibleContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
