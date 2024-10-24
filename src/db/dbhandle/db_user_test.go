package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	UserId1   = "123456"
	UserAddr1 = "123456789"
)

func insertUserTest1() ([]*db.User, error) {
	insertList := []*db.User{
		{
			UserId:    UserId1,
			UserAddr:  UserAddr1,
			Timestamp: 123456,
		},
	}
	err := BatchInsertUser(ChainID, insertList)
	return insertList, err
}

func TestBatchInsertUser(t *testing.T) {
	insertList := []*db.User{
		{
			UserId:    UserId1,
			UserAddr:  UserAddr1,
			Timestamp: 123456,
		},
	}

	type args struct {
		chainId  string
		userList []*db.User
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
				userList: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BatchInsertUser(tt.args.chainId, tt.args.userList); (err != nil) != tt.wantErr {
				t.Errorf("BatchInsertUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserList(t *testing.T) {
	_, err := insertUserTest1()
	if err != nil {
		return
	}
	type args struct {
		offset    int
		limit     int
		chainId   string
		orgId     string
		userIds   []string
		userAddrs []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.User
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
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1, err := GetUserList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.orgId, tt.args.userIds, tt.args.userAddrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 != tt.want1 {
				t.Errorf("GetUserList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetUserListByAdder(t *testing.T) {
	insertList, err := insertUserTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		adders  []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.User
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				adders: []string{
					UserAddr1,
				},
			},
			want: map[string]*db.User{
				insertList[0].UserAddr: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserListByAdder(tt.args.chainId, tt.args.adders)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserListByAdder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.User{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetUserListByAdder() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetUserNum(t *testing.T) {
	type args struct {
		chainId string
		orgId   string
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
				orgId:   "1",
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserNum(tt.args.chainId, tt.args.orgId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserNum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateUserStatus(t *testing.T) {
	_, err := insertUserTest1()
	if err != nil {
		return
	}

	type args struct {
		address string
		chainId string
		status  int
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
				address: UserAddr1,
				status:  1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateUserStatus(tt.args.address, tt.args.chainId, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
