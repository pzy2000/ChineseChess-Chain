package dbhandle

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"
)

func TestGetEvidenceContract(t *testing.T) {
	type args struct {
		offset       int
		limit        int
		chainId      string
		contractName string
		txId         string
		hashList     []string
		senderAddrs  []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.EvidenceContract
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    make([]*db.EvidenceContract, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetEvidenceContract(tt.args.offset, tt.args.limit, 0, tt.args.chainId, tt.args.contractName, tt.args.txId, "", tt.args.hashList, tt.args.senderAddrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvidenceContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEvidenceContract() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetEvidenceContract() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetEvidenceContractByHash(t *testing.T) {
	type args struct {
		chainId string
		hash    string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.EvidenceContract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				hash:    "123",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEvidenceContractByHash(tt.args.chainId, tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvidenceContractByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEvidenceContractByHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEvidenceContractByHashLit(t *testing.T) {
	type args struct {
		chainId  string
		hashList []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.EvidenceContract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				hashList: []string{
					"123",
				},
			},
			want:    make(map[string]*db.EvidenceContract, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEvidenceContractByHashLit(tt.args.chainId, tt.args.hashList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvidenceContractByHashLit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEvidenceContractByHashLit() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIdentityContract(t *testing.T) {
	type args struct {
		offset       int
		limit        int
		chainId      string
		contractAddr string
		userAddrs    []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.IdentityContract
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				offset:       0,
				limit:        0,
				contractAddr: "123",
				userAddrs: []string{
					"123",
				},
			},
			want:    make([]*db.IdentityContract, 0),
			want1:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetIdentityContract(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.contractAddr, tt.args.userAddrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIdentityContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIdentityContract() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetIdentityContract() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInsertEvidenceContract(t *testing.T) {
	insertList := []*db.EvidenceContract{
		{
			ID:           "31234567",
			TxId:         "31234567",
			SenderAddr:   "31234567",
			ContractName: "31234567",
			Hash:         "31234567",
			Timestamp:    123456,
		},
	}
	type args struct {
		chainId   string
		contracts []*db.EvidenceContract
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
			if err := InsertEvidenceContract(tt.args.chainId, tt.args.contracts); (err != nil) != tt.wantErr {
				t.Errorf("InsertEvidenceContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertIdentityContract(t *testing.T) {
	insertList := []*db.IdentityContract{
		{
			ID:           "31234567",
			TxId:         "31234567",
			ContractName: "31234567",
		},
	}
	type args struct {
		chainId   string
		contracts []*db.IdentityContract
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
			if err := InsertIdentityContract(tt.args.chainId, tt.args.contracts); (err != nil) != tt.wantErr {
				t.Errorf("InsertIdentityContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateEvidenceBak(t *testing.T) {
	type args struct {
		chainId  string
		evidence *db.EvidenceContract
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
				evidence: &db.EvidenceContract{
					Hash:        "1234",
					MetaData:    "1234",
					MetaDataBak: "1234",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateEvidenceBak(tt.args.chainId, tt.args.evidence); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEvidenceBak() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateEvidenceContractName(t *testing.T) {
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
				contractName: "1234",
				contractAddr: "1234",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateEvidenceContractName(tt.args.chainId, tt.args.contractName, tt.args.contractAddr); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEvidenceContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateIdentityContractName(t *testing.T) {
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
				contractName: "1234",
				contractAddr: "1234",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateIdentityContractName(tt.args.chainId, tt.args.contractName, tt.args.contractAddr); (err != nil) != tt.wantErr {
				t.Errorf("UpdateIdentityContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
