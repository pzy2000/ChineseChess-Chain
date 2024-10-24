package dbhandle

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"
)

func TestGetAllSubChainBlockHeight(t *testing.T) {
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
			got, err := GetAllSubChainBlockHeight(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllSubChainBlockHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAllSubChainBlockHeight() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossLatestSubChainList(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossSubChainData
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    make([]*db.CrossSubChainData, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossLatestSubChainList(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossLatestSubChainList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossLatestSubChainList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossSubChainAll(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossSubChainData
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    make([]*db.CrossSubChainData, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossSubChainAll(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossSubChainAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossSubChainAllCount(t *testing.T) {
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
			got, err := GetCrossSubChainAllCount(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainAllCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCrossSubChainAllCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossSubChainById(t *testing.T) {
	type args struct {
		chainId     string
		subChainIds []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.CrossSubChainData
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				subChainIds: []string{
					"123",
				},
			},
			want:    make(map[string]*db.CrossSubChainData, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossSubChainById(tt.args.chainId, tt.args.subChainIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossSubChainById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossSubChainInfoById(t *testing.T) {
	type args struct {
		chainId    string
		subChainId string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.CrossSubChainData
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				subChainId: "123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossSubChainInfoById(tt.args.chainId, tt.args.subChainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainInfoById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossSubChainInfoById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossSubChainInfoByName(t *testing.T) {
	type args struct {
		chainId      string
		subChainName string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.CrossSubChainData
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				subChainName: "123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossSubChainInfoByName(tt.args.chainId, tt.args.subChainName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainInfoByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossSubChainInfoByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossSubChainInfoCache(t *testing.T) {
	type args struct {
		chainId    string
		subChainId string
	}
	tests := []struct {
		name string
		args args
		want *db.CrossSubChainData
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				subChainId: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCrossSubChainInfoCache(tt.args.chainId, tt.args.subChainId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossSubChainInfoCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossSubChainList(t *testing.T) {
	type args struct {
		offset     int
		limit      int
		chainId    string
		subChainId string
		chainName  string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossSubChainData
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				offset:     0,
				limit:      10,
				subChainId: "123",
				chainName:  "123",
			},
			want: make([]*db.CrossSubChainData, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetCrossSubChainList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.subChainId, tt.args.chainName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossSubChainList() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetCrossSubChainList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetCrossSubChainListCache(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name string
		args args
		want []*db.CrossSubChainData
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCrossSubChainListCache(tt.args.chainId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossSubChainListCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossSubChainName(t *testing.T) {
	type args struct {
		chainId    string
		subChainId string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				subChainId: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossSubChainName(tt.args.chainId, tt.args.subChainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCrossSubChainName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossSubChainNameCache(t *testing.T) {
	subName := "123"
	type args struct {
		chainId    string
		subChainId string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				subChainId: "123",
			},
			want: subName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossSubChainNameCache(tt.args.chainId, tt.args.subChainId, subName)
			got, err := GetCrossSubChainNameCache(tt.args.chainId, tt.args.subChainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainNameCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCrossSubChainNameCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertCrossSubChain(t *testing.T) {
	insertList := []*db.CrossSubChainData{
		{
			SubChainId: "6666",
			TxNum:      23,
			ChainId:    "6666",
			ChainName:  "6666",
		},
	}
	type args struct {
		chainId      string
		subChainList []*db.CrossSubChainData
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
				subChainList: insertList,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertCrossSubChain(tt.args.chainId, tt.args.subChainList); (err != nil) != tt.wantErr {
				t.Errorf("InsertCrossSubChain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetCrossSubChainInfoCache(t *testing.T) {
	type args struct {
		chainId      string
		subChainId   string
		subChainInfo *db.CrossSubChainData
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				subChainId: "123",
				subChainInfo: &db.CrossSubChainData{
					ChainId:    "123",
					SubChainId: "123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossSubChainInfoCache(tt.args.chainId, tt.args.subChainId, tt.args.subChainInfo)
		})
	}
}

func TestSetCrossSubChainListCache(t *testing.T) {
	type args struct {
		chainId      string
		subChainList []*db.CrossSubChainData
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				subChainList: []*db.CrossSubChainData{
					{
						ChainId:    "123",
						SubChainId: "123",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossSubChainListCache(tt.args.chainId, tt.args.subChainList)
		})
	}
}

func TestSetCrossSubChainNameCache(t *testing.T) {
	type args struct {
		chainId      string
		subChainId   string
		subChainName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				subChainId:   "123",
				subChainName: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossSubChainNameCache(tt.args.chainId, tt.args.subChainId, tt.args.subChainName)
		})
	}
}

func TestUpdateCrossSubChainById(t *testing.T) {
	type args struct {
		chainId      string
		subChainInfo *db.CrossSubChainData
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
				subChainInfo: &db.CrossSubChainData{
					ChainId:    "123",
					SubChainId: "123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateCrossSubChainById(tt.args.chainId, tt.args.subChainInfo); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCrossSubChainById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateCrossSubChainStatus(t *testing.T) {
	type args struct {
		chainId    string
		subChainId string
		status     int32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				subChainId: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateCrossSubChainStatus(tt.args.chainId, tt.args.subChainId, "", tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCrossSubChainStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
