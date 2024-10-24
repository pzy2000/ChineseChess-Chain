package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	tokenId1 = "TokenId1"
)

func insertTokenTest() ([]*db.NonFungibleToken, error) {
	insertList := []*db.NonFungibleToken{
		{
			OwnerAddr:    ownerAddr1,
			ContractAddr: ContractAddr1,
			ContractName: ContractName1,
			TokenId:      tokenId1,
		},
	}
	err := InsertNonFungibleToken(ChainID, insertList)
	return insertList, err
}

func TestDeleteNonFungibleToken(t *testing.T) {
	insertList, err := insertTokenTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		tokenList []*db.NonFungibleToken
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
				tokenList: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteNonFungibleToken(tt.args.chainId, tt.args.tokenList); (err != nil) != tt.wantErr {
				t.Errorf("DeleteNonFungibleToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetNonFungibleTokenDetail(t *testing.T) {
	insertList, err := insertTokenTest()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		tokenId      string
		contractAddr string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.NonFungibleToken
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				tokenId:      tokenId1,
				contractAddr: ContractAddr1,
			},
			want:    insertList[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNonFungibleTokenDetail(tt.args.chainId, tt.args.tokenId, tt.args.contractAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNonFungibleTokenDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.NonFungibleToken{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetNonFungibleTokenDetail() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetNonFungibleTokenList(t *testing.T) {
	insertList, err := insertTokenTest()
	if err != nil {
		return
	}

	type args struct {
		offset      int
		limit       int
		chainId     string
		tokenId     string
		contractKey string
		ownerAddrs  []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.NonFungibleToken
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				tokenId: tokenId1,
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
			got, err := GetNonFungibleTokenList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.tokenId, tt.args.contractKey, tt.args.ownerAddrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNonFungibleTokenList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.NonFungibleToken{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetNonFungibleTokenDetail() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestInsertNonFungibleToken(t *testing.T) {
	insertList := []*db.NonFungibleToken{
		{
			OwnerAddr:    ownerAddr1,
			ContractAddr: ContractAddr1,
			ContractName: ContractName1,
			TokenId:      tokenId1,
		},
	}

	type args struct {
		chainId   string
		tokenList []*db.NonFungibleToken
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
				tokenList: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertNonFungibleToken(tt.args.chainId, tt.args.tokenList); (err != nil) != tt.wantErr {
				t.Errorf("InsertNonFungibleToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSelectTokenByID(t *testing.T) {
	insertList, err := insertTokenTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string

		tokenIds []string
		addrs    []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.NonFungibleToken
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				tokenIds: []string{
					tokenId1,
				},
				addrs: []string{
					ContractAddr1,
				},
			},
			want: map[string]*db.NonFungibleToken{
				tokenId1 + "_" + ContractAddr1: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectTokenByID(tt.args.chainId, tt.args.tokenIds, tt.args.addrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectTokenByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.NonFungibleToken{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("SelectTokenByID() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestUpdateNonFungibleToken(t *testing.T) {
	insertList, err := insertTokenTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		tokenInfo *db.NonFungibleToken
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
				tokenInfo: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateNonFungibleToken(tt.args.chainId, tt.args.tokenInfo); (err != nil) != tt.wantErr {
				t.Errorf("UpdateNonFungibleToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateNonFungibleTokenBak(t *testing.T) {
	insertList, err := insertTokenTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		tokenData *db.NonFungibleToken
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
				tokenData: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateNonFungibleTokenBak(tt.args.chainId, tt.args.tokenData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateNonFungibleTokenBak() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateTokenContractName(t *testing.T) {
	_, err := insertTokenTest()
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
				contractName: contractName1,
				contractAddr: ContractAddr1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateTokenContractName(tt.args.chainId, tt.args.contractName, tt.args.contractAddr); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTokenContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
