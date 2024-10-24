package dbhandle

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

const (
	crossId1        = "CrossId1"
	crossId2        = "CrossId2"
	crossStartTime1 = 123456
	crossEndTime1   = 223456
	crossDuration1  = 100000
)

func insertCrossCycleTxTest1() ([]*db.CrossCycleTransaction, error) {
	newUUID := uuid.New().String()
	insertList := []*db.CrossCycleTransaction{
		{
			ID:        newUUID,
			CrossId:   crossId1,
			StartTime: crossStartTime1,
			EndTime:   crossEndTime1,
			Duration:  crossDuration1,
		},
	}
	err := InsertCrossCycleTx(ChainID, insertList)
	return insertList, err
}

func TestGetCrossCycleById(t *testing.T) {
	insertList, err := insertCrossCycleTxTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		crossId string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.CrossCycleTransaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				crossId: crossId1,
			},
			want:    insertList[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossCycleById(tt.args.chainId, tt.args.crossId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossCycleById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.CrossCycleTransaction{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetCrossCycleById() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetCrossCycleTimeCache(t *testing.T) {
	typeStr := "AverageTime"
	SetCrossCycleTimeCache(ChainId1, typeStr, crossDuration1)
	type args struct {
		chainId string
		typeStr string
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
				typeStr: typeStr,
			},
			want:    crossDuration1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossCycleTimeCache(tt.args.chainId, tt.args.typeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossCycleTimeCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCrossCycleTimeCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossCycleTransactionById(t *testing.T) {
	insertList, err := insertCrossCycleTxTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId  string
		crossIds []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.CrossCycleTransaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				crossIds: []string{
					crossId1,
					"123444",
				},
			},
			want: map[string]*db.CrossCycleTransaction{
				crossId1: insertList[0],
			},
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId: ChainID,
				crossIds: []string{
					"23423423",
					"123444",
				},
			},
			want:    map[string]*db.CrossCycleTransaction{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossCycleTransactionById(tt.args.chainId, tt.args.crossIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossCycleTransactionById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.CrossCycleTransaction{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetCrossCycleTransactionById() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetCrossCycleTxAllCount(t *testing.T) {
	insertList, err := insertCrossCycleTxTest1()
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
			got, err := GetCrossCycleTxAllCount(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossCycleTxAllCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCrossCycleTxAllCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossCycleTxList(t *testing.T) {
	insertList, err := insertCrossCycleTxTest1()
	if err != nil {
		return
	}

	type args struct {
		offset      int
		limit       int
		startTime   int64
		endTime     int64
		chainId     string
		crossId     string
		subChainId  string
		fromChainId string
		toChainId   string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CycleJoinTransferResult
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
			want1:   int64(len(insertList)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1, err := GetCrossCycleTxList(tt.args.offset, tt.args.limit, tt.args.startTime, tt.args.endTime, tt.args.chainId, tt.args.crossId, tt.args.subChainId, tt.args.fromChainId, tt.args.toChainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossCycleTxList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 != tt.want1 {
				t.Errorf("GetCrossCycleTxList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetCrossLatestCycleTxCache(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name string
		args args
		want []*db.CycleJoinTransferResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCrossLatestCycleTxCache(tt.args.chainId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossLatestCycleTxCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossLatestCycleTxList(t *testing.T) {
	_, err := insertCrossCycleTxTest1()
	if err != nil {
		return
	}
	_, err = insertCrossTxTransfersTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CycleJoinTransferResult
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want: []*db.CycleJoinTransferResult{
				{
					CrossId: crossId1,
					CrossCycleTransaction: db.CrossCycleTransaction{
						ID:        "354f07a9-2897-4a3d-903b-7ba65d4abd0d",
						CrossId:   "CrossId1",
						StartTime: 123456,
						EndTime:   223456,
						Duration:  crossDuration1,
					},
					CrossTransactionTransfer: db.CrossTransactionTransfer{
						ID:           "b413a95c-3aec-44be-a21d-94dfc283f9ea",
						CrossId:      "CrossId1",
						FromChainId:  "CrossId1",
						ToChainId:    "CrossId2",
						BlockHeight:  12,
						ContractName: "subContractName1",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossLatestCycleTxList(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossLatestCycleTxList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			ignoreFields := cmpopts.IgnoreFields(db.CycleJoinTransferResult{},
				"CrossCycleTransaction.ID",
				"CrossCycleTransaction.CreatedAt",
				"CrossCycleTransaction.UpdatedAt",
				"CrossTransactionTransfer.ID",
				"CrossTransactionTransfer.CreatedAt",
				"CrossTransactionTransfer.UpdatedAt",
			)
			if !cmp.Equal(got, tt.want, ignoreFields) {
				t.Errorf("GetCrossLatestCycleTxList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetCycleAverageTime(t *testing.T) {
	_, err := insertCrossCycleTxTest1()
	if err != nil {
		return
	}

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
				startTime: crossStartTime1 - 1,
				endTime:   crossEndTime1 + 1,
			},
			want:    crossDuration1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCycleAverageTime(tt.args.chainId, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCycleAverageTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCycleAverageTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCycleLongestTime(t *testing.T) {
	_, err := insertCrossCycleTxTest1()
	if err != nil {
		return
	}

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
				startTime: crossStartTime1 - 1,
				endTime:   crossEndTime1 + 1,
			},
			want:    crossDuration1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCycleLongestTime(tt.args.chainId, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCycleLongestTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCycleLongestTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCycleShortestTime(t *testing.T) {
	_, err := insertCrossCycleTxTest1()
	if err != nil {
		return
	}

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
				startTime: crossStartTime1 - 1,
				endTime:   crossEndTime1 + 1,
			},
			want:    crossDuration1,
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:   ChainId1,
				startTime: 100,
				endTime:   200,
			},
			want:    crossDuration1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCycleShortestTime(tt.args.chainId, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCycleShortestTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCycleShortestTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertCrossCycleTx(t *testing.T) {
	newUUID := uuid.New().String()
	insertList := []*db.CrossCycleTransaction{
		{
			ID:        newUUID,
			CrossId:   crossId1,
			StartTime: crossStartTime1,
			EndTime:   crossEndTime1,
			Duration:  crossDuration1,
		},
	}

	type args struct {
		chainId       string
		crossCycleTxs []*db.CrossCycleTransaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:       ChainID,
				crossCycleTxs: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertCrossCycleTx(tt.args.chainId, tt.args.crossCycleTxs); (err != nil) != tt.wantErr {
				t.Errorf("InsertCrossCycleTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetCrossCycleTimeCache(t *testing.T) {
	type args struct {
		chainId  string
		typeStr  string
		duration int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:  ChainID,
				typeStr:  "123",
				duration: 1000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossCycleTimeCache(tt.args.chainId, tt.args.typeStr, tt.args.duration)
		})
	}
}

func TestSetCrossLatestCycleTxCache(t *testing.T) {

	crossCycleTxs := []*db.CycleJoinTransferResult{
		{
			CrossId: crossId1,
			CrossCycleTransaction: db.CrossCycleTransaction{
				ID:        "354f07a9-2897-4a3d-903b-7ba65d4abd0d",
				CrossId:   "CrossId1",
				StartTime: 123456,
				EndTime:   223456,
				Duration:  crossDuration1,
			},
			CrossTransactionTransfer: db.CrossTransactionTransfer{
				ID:           "b413a95c-3aec-44be-a21d-94dfc283f9ea",
				CrossId:      "CrossId1",
				FromChainId:  "CrossId1",
				ToChainId:    "CrossId2",
				ContractName: "subContractName1",
			},
		},
	}

	type args struct {
		chainId       string
		crossCycleTxs []*db.CycleJoinTransferResult
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:       ChainID,
				crossCycleTxs: crossCycleTxs,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossLatestCycleTxCache(tt.args.chainId, tt.args.crossCycleTxs)
		})
	}
}

func TestUpdateCrossCycleTx(t *testing.T) {
	insertList, err := insertCrossCycleTxTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		crossCycleTx *db.CrossCycleTransaction
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
				crossCycleTx: insertList[0],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateCrossCycleTx(tt.args.chainId, tt.args.crossCycleTx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCrossCycleTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getCycleQuery(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name string
		args args
		want *gorm.DB
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
			_ = getCycleQuery(tt.args.chainId)
		})
	}
}
