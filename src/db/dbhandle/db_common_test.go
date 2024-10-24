package dbhandle

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func TestBuildParamsQuery(t *testing.T) {
	type args struct {
		tableName  string
		selectFile *SelectFile
	}
	tests := []struct {
		name string
		args args
		want *gorm.DB
	}{
		{
			name: "Test BuildParamsQuery",
			args: args{
				tableName: "test_table",
				selectFile: &SelectFile{
					Where: map[string]interface{}{
						"id": 1,
					},
					NotNull: []string{"name"},
				},
			},
			want: nil, // 注意：这里的want应该是一个*gorm.DB类型的实例，但是由于我们无法预先知道其内部状态，所以暂时设置为nil
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildParamsQuery(tt.args.tableName, tt.args.selectFile)
			if got == nil {
				t.Errorf("BuildParamsQuery() = %v, want not nil", got)
			}
		})
	}
}

func TestBuildParamsQueryNew(t *testing.T) {
	type args struct {
		tableName  string
		selectFile *SelectFile
	}
	tests := []struct {
		name string
		args args
		want *gorm.DB
	}{
		{
			name: "Test BuildParamsQuery",
			args: args{
				tableName: "test_table",
				selectFile: &SelectFile{
					Where: map[string]interface{}{
						"id": 1,
					},
					NotNull: []string{"name"},
				},
			},
			want: nil, // 注意：这里的want应该是一个*gorm.DB类型的实例，但是由于我们无法预先知道其内部状态，所以暂时设置为nil
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildParamsQuery(tt.args.tableName, tt.args.selectFile)
			if got == nil {
				t.Errorf("BuildParamsQuery() = %v, want not nil", got)
			}
		})
	}
}

func TestCreateInBatchesData(t *testing.T) {
	chainList := []*db.Chain{
		{
			ChainId: ChainId1,
			Version: "1.0",
		},
		{
			ChainId: ChainId2,
			Version: "2.0",
		},
	}

	type args struct {
		tableName string
		data      interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test CreateInBatchesData",
			args: args{
				tableName: "test_chain",
				data:      chainList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateInBatchesData(tt.args.tableName, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("CreateInBatchesData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertData(t *testing.T) {
	type args struct {
		tableName string
		data      interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test InsertData",
			args: args{
				tableName: "test_chain",
				data:      map[string]interface{}{"chainId": "chain1", "version": "1.0"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertData(tt.args.tableName, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("InsertData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestToInterfaceSlice(t *testing.T) {
	type args struct {
		slice interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "Test ToInterfaceSlice",
			args: args{
				slice: []int{1, 2, 3},
			},
			want: []interface{}{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInterfaceSlice(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToInterfaceSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isDuplicateKeyError(t *testing.T) {
	mysqlErr1 := &mysql.MySQLError{}
	mysqlErr1.Number = 1062
	mysqlErr2 := &mysql.MySQLError{}
	mysqlErr2.Number = 1054

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test isDuplicateKeyError",
			args: args{
				err: mysqlErr1,
			},
			want: true,
		},
		{
			name: "Test isDuplicateKeyError",
			args: args{
				err: mysqlErr2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDuplicateKeyError(tt.args.err); got != tt.want {
				t.Errorf("isDuplicateKeyError() = %v, want %v", got, tt.want)
			}
		})
	}
}
