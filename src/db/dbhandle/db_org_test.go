package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
)

const (
	orgId1 = "123"
	orgId2 = "223"
)

func insertOrgTest() ([]*db.Org, error) {
	insertList := []*db.Org{
		{
			OrgId: orgId1,
		},
		{
			OrgId: orgId2,
		},
	}
	tableName := db.GetTableName(ChainID, db.TableOrg)
	err := CreateInBatchesData(tableName, insertList)
	return insertList, err
}

func TestGetOrgList(t *testing.T) {
	insertList, err := insertOrgTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		orgId   string
		offset  int
		limit   int
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Org
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
			want:    insertList,
			want1:   int64(len(insertList)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetOrgList(tt.args.chainId, tt.args.orgId, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrgList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Org{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetOrgList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
			if got1 != tt.want1 {
				t.Errorf("GetOrgList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetOrgNum(t *testing.T) {
	insertList, err := insertOrgTest()
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
			want:    int64(len(insertList)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOrgNum(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrgNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetOrgNum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveOrgByConfig(t *testing.T) {
	type args struct {
		chainConfig *pbconfig.ChainConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test SaveOrgByConfig with valid input",
			args: args{
				chainConfig: &pbconfig.ChainConfig{
					ChainId: ChainID,
					TrustRoots: []*pbconfig.TrustRootConfig{
						{OrgId: "org1"},
						{OrgId: "org2"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test SaveOrgByConfig with no chain ID",
			args: args{
				chainConfig: &pbconfig.ChainConfig{
					ChainId: "",
					TrustRoots: []*pbconfig.TrustRootConfig{
						{OrgId: "org1"},
						{OrgId: "org2"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test SaveOrgByConfig with no trust roots",
			args: args{
				chainConfig: &pbconfig.ChainConfig{
					ChainId:    ChainID,
					TrustRoots: []*pbconfig.TrustRootConfig{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveOrgByConfig(tt.args.chainConfig); (err != nil) != tt.wantErr {
				t.Errorf("SaveOrgByConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
