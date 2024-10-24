package dbhandle

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/db"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
)

const (
	ChainID      = "chain1"
	AccountAddr1 = "123456789"
	AccountAddr2 = "223456789"
	AccountAddr3 = "323456789"
	AccountAddr4 = "423456789"
	BNSAddr1     = "bns.com"
	BNSAddr2     = "bns.2com"
	BNSAddr3     = "bns.com2"
	BNSAddr4     = "bns.com3"
	DIDAddr1     = "did:12345"
	DIDAddr2     = "did:22345"
	DIDAddr3     = "did:42345"
	DIDAddr4     = "did:42345"
)

func insertAccountTest() ([]*db.Account, error) {
	accountList := []*db.Account{
		{
			AddrType: 0,
			Address:  AccountAddr1,
			DID:      DIDAddr1,
			BNS:      BNSAddr1,
		},
		{
			AddrType: 0,
			Address:  AccountAddr2,
			DID:      DIDAddr2,
			BNS:      BNSAddr2,
		},
	}
	err := InsertAccount(ChainID, accountList)
	return accountList, err
}

func insertAccountTest2() ([]*db.Account, error) {
	accountList := []*db.Account{
		{
			AddrType: 1,
			Address:  AccountAddr3,
			DID:      DIDAddr3,
			BNS:      BNSAddr3,
		},
		{
			AddrType: 1,
			Address:  AccountAddr4,
			DID:      DIDAddr4,
			BNS:      BNSAddr4,
		},
	}
	err := InsertAccount(ChainID, accountList)
	return accountList, err
}

func TestMain(m *testing.M) {
	//初始化配置
	//_ = config.InitConfig("", "")
	// 初始化数据库配置
	dbCfg, err := db.InitMySQLContainer()
	if err != nil || dbCfg == nil {
		return
	}

	redisCfg, err := db.InitRedisContainer()
	if err != nil || redisCfg == nil {
		return
	}
	cache.InitRedis(redisCfg)
	// 运行其他测试
	os.Exit(m.Run())
}

func TestGetAccountByAddr(t *testing.T) {
	_, err := insertAccountTest()
	if err != nil {
		return
	}

	gotWant := &db.Account{
		AddrType: 0,
		Address:  AccountAddr1,
		DID:      DIDAddr1,
		BNS:      BNSAddr1,
	}
	type args struct {
		chainId string
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				address: AccountAddr1,
			},
			wantErr: false,
			want:    gotWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccountByAddr(tt.args.chainId, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("DealContractEvents() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")))
			}
		})
	}
}

func TestGetAccountByBNS(t *testing.T) {
	_, err := insertAccountTest()
	if err != nil {
		return
	}

	gotWant := &db.Account{
		AddrType: 0,
		Address:  AccountAddr1,
		DID:      "did:12345",
		BNS:      "bns.com",
	}

	type args struct {
		chainId string
		bns     string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				bns:     "bns.com",
			},
			wantErr: false,
			want:    gotWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccountByBNS(tt.args.chainId, tt.args.bns)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByBNS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("DealContractEvents() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")))
			}
		})
	}
}

func TestGetAccountByDID(t *testing.T) {
	_, err := insertAccountTest()
	if err != nil {
		return
	}

	gotWant := []*db.Account{
		{
			AddrType: 0,
			Address:  AccountAddr1,
			DID:      "did:12345",
			BNS:      "bns.com",
		},
	}

	type args struct {
		chainId string
		did     string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				did:     "did:12345",
			},
			wantErr: false,
			want:    gotWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccountByDID(tt.args.chainId, tt.args.did)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByDID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("DealContractEvents() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")))
			}
		})
	}
}

func TestGetAccountDetail(t *testing.T) {
	_, err := insertAccountTest()
	if err != nil {
		return
	}

	gotWant := &db.Account{
		AddrType: 0,
		Address:  AccountAddr1,
		DID:      "did:12345",
		BNS:      "bns.com",
	}

	type args struct {
		chainId string
		address string
		bns     string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				address: AccountAddr1,
				bns:     BNSAddr1,
			},
			wantErr: false,
			want:    gotWant,
		},
		{
			name: "test: case 2",
			args: args{
				chainId: ChainID,
				address: AccountAddr1,
				bns:     "bns.com23232",
			},
			wantErr: false,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccountDetail(tt.args.chainId, tt.args.address, tt.args.bns)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("DealContractEvents() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")))
			}
		})
	}
}

func TestGetAccountList(t *testing.T) {
	accountList, err := insertAccountTest()
	if err != nil {
		return
	}

	type args struct {
		offset   int
		limit    int
		chainId  string
		addrType *int
	}
	addrType := 0
	tests := []struct {
		name    string
		args    args
		want    []*db.Account
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 2",
			args: args{
				chainId:  ChainID,
				offset:   0,
				limit:    10,
				addrType: &addrType,
			},
			wantErr: false,
			want:    accountList,
			want1:   int64(len(accountList)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetAccountList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.addrType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("DealContractEvents() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")))
			}

			if got1 != tt.want1 {
				t.Errorf("GetAccountList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInsertAccount(t *testing.T) {
	accountList := []*db.Account{
		{
			AddrType: 1,
			Address:  "12345678933333",
			DID:      "did:12345333",
			BNS:      "bns.com33",
		},
		{
			AddrType: 1,
			Address:  "2234567893333",
			DID:      "did:22345333",
			BNS:      "bns.2com333",
		},
	}
	type args struct {
		chainId     string
		accountList []*db.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:     "chain1",
				accountList: accountList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertAccount(tt.args.chainId, tt.args.accountList); (err != nil) != tt.wantErr {
				t.Errorf("InsertAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueryAccountExists(t *testing.T) {
	addrListDB, err := insertAccountTest()
	if err != nil {
		return
	}

	wantMap := map[string]*db.Account{}
	for _, value := range addrListDB {
		wantMap[value.Address] = value
	}

	type args struct {
		chainId  string
		addrList []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: "chain1",
				addrList: []string{
					AccountAddr1,
					AccountAddr2,
					"343234345",
				},
			},
			wantErr: false,
			want:    wantMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryAccountExists(tt.args.chainId, tt.args.addrList)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryAccountExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("DealContractEvents() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")))
			}
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	_, err := insertAccountTest2()
	if err != nil {
		return
	}

	accountInfo := &db.Account{
		AddrType: 1,
		Address:  AccountAddr3,
		DID:      "did:12345",
		BNS:      "bns.com",
	}

	type args struct {
		chainId     string
		accountInfo *db.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:     ChainID,
				accountInfo: accountInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateAccount(tt.args.chainId, tt.args.accountInfo); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAccountByBNSList(t *testing.T) {
	type args struct {
		chainId string
		bnsList []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				bnsList: []string{
					"BNS:123",
				},
			},
			want:    make([]*db.Account, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccountByBNSList(tt.args.chainId, tt.args.bnsList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByBNSList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetAccountByBNSList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")))
			}
		})
	}
}

func TestGetAccountByDIDList(t *testing.T) {
	type args struct {
		chainId string
		didList []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				didList: []string{
					"BNS:123",
				},
			},
			want:    make([]*db.Account, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccountByDIDList(tt.args.chainId, tt.args.didList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByDIDList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetAccountByDIDList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, cmpopts.IgnoreFields(db.Account{}, "CreatedAt", "UpdatedAt")))
			}
		})
	}
}
