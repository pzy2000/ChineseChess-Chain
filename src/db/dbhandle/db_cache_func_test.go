package dbhandle

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"
)

func TestGetTotalTxNumCache(t *testing.T) {
	SetTotalTxNumCache(ChainID, 123)
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
			want:    123,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTotalTxNumCache(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTotalTxNumCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTotalTxNumCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMaxBlockHeightCache(t *testing.T) {
	SetMaxBlockHeightCache(ChainID, 12)
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
			want:    12,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMaxBlockHeightCache(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMaxBlockHeightCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMaxBlockHeightCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLatestTxListCache(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Transaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    make([]*db.Transaction, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestTxListCache(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestTxListCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLatestTxListCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLatestBlockListCache(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Block
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    make([]*db.Block, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetLatestBlockListCache(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestBlockListCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
