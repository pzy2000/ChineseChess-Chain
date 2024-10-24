package dbhandle

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	contractName1 = "contract1"
	contractName2 = "contract2"
	//contractName3  = "contract3"
	contractAdder1 = "12345678"
	contractAdder2 = "22345678"
	//contractAdder3 = "32345678"
)

func insertContractTest1() (*db.Contract, error) {
	contractInfo := &db.Contract{
		Name:         contractName1,
		NameBak:      contractName1,
		Addr:         contractAdder1,
		RuntimeType:  "DOCKER_GO",
		ContractType: "CMDFA",
		TxNum:        12,
		Timestamp:    12345,
	}
	err := InsertContract(ChainID, contractInfo)
	return contractInfo, err
}

func insertContractTest2() (*db.Contract, error) {
	contractInfo := &db.Contract{
		Name:           contractName2,
		NameBak:        contractName2,
		Addr:           contractAdder2,
		RuntimeType:    "EVM",
		ContractType:   "ERC20",
		ContractStatus: 0,
		Timestamp:      123456,
	}
	err := InsertContract(ChainID, contractInfo)
	return contractInfo, err
}

func TestGetContractByAdders(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}
	contractInfo2, err := insertContractTest2()
	if err != nil {
		return
	}

	contractAdders := []string{
		contractAdder1,
		contractAdder2,
	}
	type args struct {
		chainId        string
		contractAdders []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:        ChainID,
				contractAdders: contractAdders,
			},
			want: map[string]*db.Contract{
				contractAdder2: contractInfo2,
				contractAdder1: contractInfo1,
				contractName1:  contractInfo1,
				contractName2:  contractInfo2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractByCacheOrAddrs(tt.args.chainId, tt.args.contractAdders)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByCacheOrAddrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Contract{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetContractByCacheOrAddrs() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetContractByAddersOrNames(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}
	contractInfo2, err := insertContractTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId  string
		nameList []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				nameList: []string{
					contractName1,
					contractAdder2,
				},
			},
			want: map[string]*db.Contract{
				contractAdder2: contractInfo2,
				contractName2:  contractInfo2,
				contractAdder1: contractInfo1,
				contractName1:  contractInfo1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractByAddersOrNames(tt.args.chainId, tt.args.nameList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByAddersOrNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Contract{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetContractByAddersOrNames() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetContractByAddr(t *testing.T) {
	contractInfo1, err := insertContractTest1()
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
		want    *db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractAddr: contractAdder1,
			},
			want:    contractInfo1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractByAddr(tt.args.chainId, tt.args.contractAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Contract{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetContractByAddersOrNames() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetContractByName(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		contractName string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractName: contractName1,
			},
			want:    contractInfo1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractByName(tt.args.chainId, tt.args.contractName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Contract{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetContractByAddersOrNames() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetContractByNameOrAddr(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId     string
		contractKey string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:     ChainID,
				contractKey: contractName1,
			},
			want:    contractInfo1,
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:     ChainID,
				contractKey: contractAdder1,
			},
			want:    contractInfo1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractByNameOrAddr(tt.args.chainId, tt.args.contractKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByNameOrAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Contract{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetContractByNameOrAddr() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetContractList(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}
	contractInfo2, err := insertContractTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		offset       int
		limit        int
		status       *int32
		runtimeType  string
		contractKey  string
		creators     []string
		creatorAddrs []string
		upgrades     []string
		upgradeAddrs []string
		startTime    int64
		endTime      int64
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Contract
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
			want: []*db.Contract{
				contractInfo2,
				contractInfo1,
			},
			want1:   2,
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:     ChainID,
				offset:      0,
				limit:       10,
				contractKey: contractAdder1,
			},
			want: []*db.Contract{
				contractInfo1,
			},
			want1:   1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetContractList(tt.args.chainId, tt.args.offset, tt.args.limit, tt.args.status, tt.args.runtimeType, tt.args.contractKey,
				tt.args.creators, tt.args.creatorAddrs, tt.args.upgrades, tt.args.upgradeAddrs, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Contract{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetContractByNameOrAddr() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
			if got1 != tt.want1 {
				t.Errorf("GetContractList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetContractNum(t *testing.T) {
	_, err := insertContractTest1()
	if err != nil {
		return
	}
	_, err = insertContractTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId string
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
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractNum(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetContractNum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetContractNumCache(t *testing.T) {
	_, err := insertContractTest1()
	if err != nil {
		return
	}

	_, err = insertContractTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId string
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
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractNum(tt.args.chainId)
			if err != nil {
				return
			}

			prefix := config.GlobalConfig.RedisDB.Prefix
			redisKey := fmt.Sprintf(cache.RedisOverviewContractCount, prefix, tt.args.chainId)
			got, err = GetContractNumCache(redisKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractNumCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetContractNumCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLatestContractList(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}
	contractInfo2, err := insertContractTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want: []*db.Contract{
				contractInfo2,
				contractInfo1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestContractList(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestContractList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Contract{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetContractByNameOrAddr() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestInsertContract(t *testing.T) {
	if db.GormDB == nil {
		return
	}
	_, err := insertContractTest1()
	if err != nil {
		t.Errorf("InsertContract() error = %v", err)
	}
}

func TestUpdateContract(t *testing.T) {
	contractInfo2, err := insertContractTest2()
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
				chainId: ChainID,
				contract: &db.Contract{
					Addr:           contractInfo2.Addr,
					ContractStatus: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateContract(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateContractNameBak(t *testing.T) {
	contractInfo2, err := insertContractTest2()
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
				chainId: ChainID,
				contract: &db.Contract{
					Addr:    contractInfo2.Addr,
					Name:    "123",
					NameBak: "123",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateContractNameBak(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateContractNameBak() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateContractTxNum(t *testing.T) {
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
				chainId: ChainID,
				contract: &db.Contract{
					Addr:  contractInfo1.Addr,
					TxNum: 23,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateContractTxNum(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateContractTxNum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTotalTxNum(t *testing.T) {
	_, err := insertContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId string
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
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTotalTxNum(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTotalTxNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == 0 {
				t.Errorf("GetTotalTxNum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetContractByCacheOrAddr(t *testing.T) {
	type args struct {
		chainId      string
		contractAddr string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractAddr: contractAdder1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetContractByCacheOrAddr(tt.args.chainId, tt.args.contractAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByCacheOrAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetContractByCacheOrName(t *testing.T) {
	type args struct {
		chainId      string
		contractName string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractName: "name",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractByCacheOrName(tt.args.chainId, tt.args.contractName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByCacheOrName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContractByCacheOrName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetContractByCacheOrNameAddr(t *testing.T) {
	type args struct {
		chainId     string
		contractKey string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:     ChainID,
				contractKey: contractAdder1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetContractByCacheOrNameAddr(tt.args.chainId, tt.args.contractKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByCacheOrNameAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetContractCountByRange(t *testing.T) {
	type args struct {
		chainId   string
		startTime int64
		endTime   int64
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
				startTime: 123455667,
				endTime:   123366666,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetContractCountByRange(tt.args.chainId, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractCountByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
