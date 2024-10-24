package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"
)

func insertFungibleTransferTest() ([]*db.FungibleTransfer, error) {
	insertList := []*db.FungibleTransfer{
		{
			TxId:         ContractTxId1,
			ContractName: ContractName1,
			ContractAddr: ContractAddr1,
			FromAddr:     "123",
			ToAddr:       "456",
			Amount:       GetAmountDecimal("2"),
			Timestamp:    1234,
		},
	}
	err := InsertFungibleTransfer(ChainID, insertList)
	return insertList, err
}

func insertNonFungibleTransferTest() ([]*db.NonFungibleTransfer, error) {
	insertList := []*db.NonFungibleTransfer{
		{
			TxId:         ContractTxId1,
			ContractName: ContractName1,
			ContractAddr: ContractAddr1,
			FromAddr:     "123",
			ToAddr:       "456",
			TokenId:      "2",
			Timestamp:    1234,
		},
	}
	err := InsertNonFungibleTransfer(ChainID, insertList)
	return insertList, err
}

func TestGetNonFungibleTransferList(t *testing.T) {
	_, err := insertNonFungibleTransferTest()
	if err != nil {
		return
	}

	type args struct {
		offset       int
		limit        int
		chainId      string
		tokenId      string
		contractAddr string
		userAddr     string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.NonFungibleTransfer
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractAddr: ContractAddr1,
				offset:       0,
				limit:        10,
			},
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetNFTTransferList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.contractAddr, tt.args.userAddr, tt.args.tokenId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNonFungibleTransferList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInsertFungibleTransfer(t *testing.T) {
	insertList := []*db.FungibleTransfer{
		{
			TxId:         ContractTxId1,
			ContractName: ContractName1,
			ContractAddr: ContractAddr1,
			FromAddr:     "123",
			ToAddr:       "456",
			Amount:       GetAmountDecimal("2"),
			Timestamp:    1234,
		},
	}

	type args struct {
		chainId   string
		transfers []*db.FungibleTransfer
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
				transfers: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertFungibleTransfer(tt.args.chainId, tt.args.transfers); (err != nil) != tt.wantErr {
				t.Errorf("InsertFungibleTransfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertNonFungibleTransfer(t *testing.T) {
	insertList := []*db.NonFungibleTransfer{
		{
			TxId:         ContractTxId1,
			ContractName: ContractName1,
			ContractAddr: ContractAddr1,
			FromAddr:     "123",
			ToAddr:       "456",
			TokenId:      "2",
			Timestamp:    1234,
		},
	}

	type args struct {
		chainId   string
		transfers []*db.NonFungibleTransfer
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
				transfers: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertNonFungibleTransfer(tt.args.chainId, tt.args.transfers); (err != nil) != tt.wantErr {
				t.Errorf("InsertNonFungibleTransfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateNonTransferContractName(t *testing.T) {
	_, err := insertNonFungibleTransferTest()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		contractName string
		contractAddr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractName: ContractName2,
				contractAddr: ContractAddr1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateNonTransferContractName(tt.args.chainId, tt.args.contractName, tt.args.contractAddr); (err != nil) != tt.wantErr {
				t.Errorf("UpdateNonTransferContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateTransferContractName(t *testing.T) {
	_, err := insertFungibleTransferTest()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		contractName string
		contractAddr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractName: ContractName2,
				contractAddr: ContractAddr1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateTransferContractName(tt.args.chainId, tt.args.contractName, tt.args.contractAddr); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTransferContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFTTransferList(t *testing.T) {
	type args struct {
		offset       int
		limit        int
		chainId      string
		contractAddr string
		userAddr     string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.FungibleTransfer
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				userAddr:     "1234",
				contractAddr: ContractAddr1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetFTTransferList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.contractAddr, tt.args.userAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFTTransferList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetFTTransferCount(t *testing.T) {
	type args struct {
		chainId      string
		contractAddr string
		userAddr     string
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
				chainId:      ChainID,
				userAddr:     "1234",
				contractAddr: ContractAddr1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFTTransferCount(tt.args.chainId, tt.args.contractAddr, tt.args.userAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFTTransferCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFTTransferCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNFTTransferCount(t *testing.T) {
	type args struct {
		chainId      string
		contractAddr string
		userAddr     string
		tokenId      string
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
				chainId:      ChainID,
				userAddr:     "1234",
				contractAddr: ContractAddr1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNFTTransferCount(tt.args.chainId, tt.args.contractAddr, tt.args.userAddr, tt.args.tokenId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNFTTransferCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNFTTransferCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNFTTransferList(t *testing.T) {
	type args struct {
		offset       int
		limit        int
		chainId      string
		contractAddr string
		userAddr     string
		tokenId      string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.NonFungibleTransfer
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				userAddr:     "1234",
				contractAddr: ContractAddr1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetNFTTransferList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.contractAddr, tt.args.userAddr, tt.args.tokenId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNFTTransferList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
