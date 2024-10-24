package db

import (
	"reflect"
	"testing"
)

func TestGetTableName(t *testing.T) {
	type args struct {
		chainId string
		table   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "tase case 1",
			args: args{
				chainId: "chain1",
				table:   "test_table",
			},
			want: "chain1_test_table",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTableName(tt.args.chainId, tt.args.table); got != tt.want {
				t.Errorf("GetTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetClickhouseTableOptions(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "获取Clickhouse表选项",
			want: GetClickhouseTableOptions(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetClickhouseTableOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClickhouseTableOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
