package dbhandle

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shopspring/decimal"
)

const (
	ownerAddr1 = "123456"
	ownerAddr2 = "223456"
)

func GetAmountDecimal(amount string) decimal.Decimal {
	// 将字符串转换为 decimal.Decimal 值
	amountDecimal, _ := decimal.NewFromString(amount)
	return amountDecimal
}

func insertPositionTest() ([]*db.FungiblePosition, error) {
	insertList := []*db.FungiblePosition{
		{
			OwnerAddr:    ownerAddr1,
			ContractAddr: ContractAddr1,
			ContractName: ContractName1,
			Amount:       GetAmountDecimal("1234"),
		},
	}
	err := InsertFungiblePosition(ChainID, insertList)
	return insertList, err
}

func insertNonPositionTest() ([]*db.NonFungiblePosition, error) {
	insertList := []*db.NonFungiblePosition{
		{
			OwnerAddr:    ownerAddr1,
			ContractAddr: ContractAddr1,
			ContractName: ContractName1,
			Amount:       GetAmountDecimal("1234"),
		},
	}
	err := InsertNonFungiblePosition(ChainID, insertList)
	return insertList, err
}

func TestBuildPositionRankSql(t *testing.T) {
	type args struct {
		tableName    string
		contractAddr string
		ownerAddr    string
		limit        int
		offset       int
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []interface{}
	}{
		{
			name: "Test BuildPositionRankSql with contractAddr",
			args: args{
				tableName:    "testTable",
				contractAddr: "testContract",
				ownerAddr:    "",
				limit:        10,
				offset:       0,
			},
			want:  "SELECT * FROM (SELECT *, RANK() OVER (ORDER BY CAST(amount AS DECIMAL(65, 30)) DESC) as holdRank FROM testTable WHERE contractAddr = ? ) as subquery WHERE contractAddr = ? ORDER BY holdRank ASC LIMIT ? OFFSET ?",
			want1: []interface{}{"testContract", "testContract", 10, 0},
		},
		{
			name: "Test BuildPositionRankSql with ownerAddr",
			args: args{
				tableName:    "testTable",
				contractAddr: "",
				ownerAddr:    "testOwner",
				limit:        10,
				offset:       0,
			},
			want:  "SELECT * FROM (SELECT *, RANK() OVER (ORDER BY CAST(amount AS DECIMAL(65, 30)) DESC) as holdRank FROM testTable WHERE ownerAddr = ? ) as subquery WHERE ownerAddr = ? ORDER BY holdRank ASC LIMIT ? OFFSET ?",
			want1: []interface{}{"testOwner", "testOwner", 10, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := BuildPositionRankSql(tt.args.tableName, tt.args.contractAddr, tt.args.ownerAddr, tt.args.limit, tt.args.offset)
			if got != tt.want {
				t.Errorf("BuildPositionRankSql() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("BuildPositionRankSql() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeleteFungiblePosition(t *testing.T) {
	insetList, err := insertPositionTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		positions []*db.FungiblePosition
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
				positions: insetList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteFungiblePosition(tt.args.chainId, tt.args.positions); (err != nil) != tt.wantErr {
				t.Errorf("DeleteFungiblePosition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteNonFungiblePosition(t *testing.T) {
	insetList, err := insertNonPositionTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		positions []*db.NonFungiblePosition
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
				positions: insetList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteNonFungiblePosition(tt.args.chainId, tt.args.positions); (err != nil) != tt.wantErr {
				t.Errorf("DeleteNonFungiblePosition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFungiblePositionByOwners(t *testing.T) {
	insetList, err := insertPositionTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		ownerAddr []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]*db.FungiblePosition
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				ownerAddr: []string{
					ownerAddr1,
					ownerAddr2,
				},
			},
			want: map[string][]*db.FungiblePosition{
				ownerAddr1: insetList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFungiblePositionByOwners(tt.args.chainId, tt.args.ownerAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFungiblePositionByOwners() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.FungiblePosition{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetFungiblePositionByOwners() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetFungiblePositionList(t *testing.T) {
	insetList, err := insertPositionTest()
	if err != nil {
		return
	}

	type args struct {
		offset       int
		limit        int
		chainId      string
		contractAddr string
		ownerAddr    string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.FungiblePosition
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				offset:       0,
				limit:        10,
				contractAddr: ContractAddr1,
			},
			want:    insetList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFungiblePositionList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.contractAddr, tt.args.ownerAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFungiblePositionList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.FungiblePosition{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetFungiblePositionList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetNFTPositionHoldRankByAmount(t *testing.T) {
	insetList, err := insertNonPositionTest()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		contractAddr string
		amount       string
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
				contractAddr: ContractAddr1,
			},
			want:    int64(len(insetList)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNFTPositionHoldRankByAmount(tt.args.chainId, tt.args.contractAddr, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNFTPositionHoldRankByAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNFTPositionHoldRankByAmount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNonFungiblePositionByOwner(t *testing.T) {
	insetList, err := insertNonPositionTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		ownerAddr []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]*db.NonFungiblePosition
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				ownerAddr: []string{
					ownerAddr1,
					ownerAddr2,
				},
			},
			want: map[string][]*db.NonFungiblePosition{
				ownerAddr1: insetList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNonFungiblePositionByOwner(tt.args.chainId, tt.args.ownerAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNonFungiblePositionByOwner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.NonFungiblePosition{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetNonFungiblePositionByOwner() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetNonFungiblePositionList(t *testing.T) {
	insetList, err := insertNonPositionTest()
	if err != nil {
		return
	}

	type args struct {
		offset       int
		limit        int
		chainId      string
		contractAddr string
		ownerAddr    string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.NonFungiblePosition
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				offset:       0,
				limit:        10,
				contractAddr: ContractAddr1,
			},
			want:    insetList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNFTPositionList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.contractAddr, tt.args.ownerAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNonFungiblePositionList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.NonFungiblePosition{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetNonFungiblePositionList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetNonFungiblePositionListWithRank(t *testing.T) {
	insertList, err := insertNonPositionTest()
	if err != nil {
		return
	}

	type args struct {
		offset       int
		limit        int
		chainId      string
		contractAddr string
		ownerAddr    string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.PositionWithRank
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				offset:       0,
				limit:        10,
				contractAddr: ContractAddr1,
			},
			want: []*db.PositionWithRank{
				{
					FungiblePosition: db.FungiblePosition{
						OwnerAddr:    insertList[0].OwnerAddr,
						ContractAddr: insertList[0].ContractAddr,
						ContractName: insertList[0].ContractName,
						Amount:       insertList[0].Amount,
					},
					HoldRank: 1,
				},
			},
			want1:   int64(len(insertList)),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetNonFungiblePositionListWithRank(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.contractAddr, tt.args.ownerAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNonFungiblePositionListWithRank() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.PositionWithRank{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetNonFungiblePositionListWithRank() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
			if got1 != tt.want1 {
				t.Errorf("GetNonFungiblePositionListWithRank() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInsertFungiblePosition(t *testing.T) {
	insertList := []*db.FungiblePosition{
		{
			OwnerAddr:    ownerAddr1,
			ContractAddr: ContractAddr1,
			ContractName: ContractName1,
			Amount:       GetAmountDecimal("1234"),
		},
	}

	type args struct {
		chainId   string
		positions []*db.FungiblePosition
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
				positions: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertFungiblePosition(tt.args.chainId, tt.args.positions); (err != nil) != tt.wantErr {
				t.Errorf("InsertFungiblePosition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertNonFungiblePosition(t *testing.T) {
	insertList := []*db.NonFungiblePosition{
		{
			OwnerAddr:    ownerAddr1,
			ContractAddr: ContractAddr1,
			ContractName: ContractName1,
			Amount:       GetAmountDecimal("1234"),
		},
	}

	type args struct {
		chainId   string
		positions []*db.NonFungiblePosition
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
				positions: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertNonFungiblePosition(tt.args.chainId, tt.args.positions); (err != nil) != tt.wantErr {
				t.Errorf("InsertNonFungiblePosition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateFungiblePosition(t *testing.T) {
	insertList, err := insertPositionTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		positions []*db.FungiblePosition
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
				positions: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateFungiblePosition(tt.args.chainId, tt.args.positions); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFungiblePosition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateNonFungiblePosition(t *testing.T) {
	insertList, err := insertNonPositionTest()
	if err != nil {
		return
	}
	type args struct {
		chainId   string
		positions []*db.NonFungiblePosition
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
				positions: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateNonFungiblePosition(tt.args.chainId, tt.args.positions); (err != nil) != tt.wantErr {
				t.Errorf("UpdateNonFungiblePosition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateNonPositionContractName(t *testing.T) {
	_, err := insertNonPositionTest()
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
			if err := UpdateNonPositionContractName(tt.args.chainId, tt.args.contractName, tt.args.contractAddr); (err != nil) != tt.wantErr {
				t.Errorf("UpdateNonPositionContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdatePositionContractName(t *testing.T) {
	_, err := insertPositionTest()
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
			if err := UpdatePositionContractName(tt.args.chainId, tt.args.contractName, tt.args.contractAddr); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePositionContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFTPositionByAddrJoinAccount(t *testing.T) {
	type args struct {
		chainId      string
		contractAddr string
		ownerAddr    []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.ContractPositionAccount
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				ownerAddr:    []string{"123", "456"},
				contractAddr: ContractAddr1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetFTPositionByAddrJoinAccount(tt.args.chainId, tt.args.contractAddr, tt.args.ownerAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFTPositionByAddrJoinAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetFTPositionJoinAccount(t *testing.T) {
	type args struct {
		offset       int
		limit        int
		chainId      string
		contractAddr string
		ownerAddr    string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.ContractPositionAccount
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				ownerAddr:    "123",
				contractAddr: ContractAddr1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetFTPositionJoinAccount(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.contractAddr, tt.args.ownerAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFTPositionJoinAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSetNFTPositionListCache(t *testing.T) {
	nftPosition := []*db.NonFungiblePosition{
		{
			ContractAddr: ContractAddr1,
			ContractName: "123",
			OwnerAddr:    "123",
		},
	}
	type args struct {
		chainId      string
		contractAddr string
		nftPosition  []*db.NonFungiblePosition
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractAddr: ContractAddr1,
				nftPosition:  nftPosition,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetNFTPositionListCache(tt.args.chainId, tt.args.contractAddr, tt.args.nftPosition)
		})
	}
}

func TestSetFTPositionListCache(t *testing.T) {
	ftPosition := []*db.FungiblePosition{
		{
			ContractAddr: ContractAddr1,
			ContractName: "123",
			OwnerAddr:    "123",
		},
	}
	type args struct {
		chainId      string
		contractAddr string
		ftPosition   []*db.FungiblePosition
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractAddr: ContractAddr1,
				ftPosition:   ftPosition,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetFTPositionListCache(tt.args.chainId, tt.args.contractAddr, tt.args.ftPosition)
		})
	}
}

func TestGetFungiblePositionByOwnerAddr(t *testing.T) {
	type args struct {
		chainId   string
		ownerAddr string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.FungiblePosition
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:   ChainID,
				ownerAddr: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetFungiblePositionByOwnerAddr(tt.args.chainId, tt.args.ownerAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFungiblePositionByOwnerAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetFTPositionListByAddr(t *testing.T) {
	type args struct {
		offset    int
		limit     int
		chainId   string
		ownerAddr string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.FTPositionJoinContract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:   ChainID,
				ownerAddr: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetFTPositionListByAddr(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.ownerAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFTPositionListByAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetFTPositionCountByAddr(t *testing.T) {
	type args struct {
		chainId   string
		ownerAddr string
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
				chainId:   ChainID,
				ownerAddr: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFTPositionCountByAddr(tt.args.chainId, tt.args.ownerAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFTPositionCountByAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFTPositionCountByAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
