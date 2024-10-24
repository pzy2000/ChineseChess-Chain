package dbhandle

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	txId1      = "123456"
	txId2      = "223456"
	timestamp1 = 12345000
	timestamp2 = 12445000
)

func insertTxTest() ([]*db.Transaction, error) {
	insertList := []*db.Transaction{
		{
			TxId:         txId1,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp1,
			ContractName: ContractName1,
		},
		{
			TxId:         txId2,
			Sender:       "123",
			UserAddr:     ContractTxUser2,
			BlockHeight:  12,
			Timestamp:    timestamp2,
			ContractName: ContractName2,
		},
	}
	err := InsertTransactions(ChainID, insertList)
	return insertList, err
}

func insertBlackTransactionsTest() ([]*db.BlackTransaction, error) {
	insertList := []*db.BlackTransaction{
		{
			TxId:         txId1,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp1,
			ContractName: ContractName1,
		},
		{
			TxId:         txId2,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp2,
			ContractName: ContractName2,
		},
	}
	err := InsertBlackTransactions(ChainID, insertList)
	return insertList, err
}

func TestBatchQueryBlackTxList(t *testing.T) {
	txList, err := insertBlackTransactionsTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		txIds   []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.BlackTransaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				txIds: []string{
					txId1,
					txId2,
				},
			},
			want:    txList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BatchQueryBlackTxList(tt.args.chainId, tt.args.txIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchQueryBlackTxList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.BlackTransaction{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetSubscribeByChainId() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestBatchQueryTxList(t *testing.T) {
	txList, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		txIds   []string
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
				txIds: []string{
					txId1,
					txId2,
				},
			},
			want:    txList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BatchQueryTxList(tt.args.chainId, tt.args.txIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchQueryTxList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Transaction{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("BatchQueryTxList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestDeleteBlackTransaction(t *testing.T) {
	_, err := insertBlackTransactionsTest()
	if err != nil {
		return
	}

	insertList := []*db.Transaction{
		{
			TxId:        txId1,
			Sender:      "123",
			UserAddr:    ContractTxUser1,
			BlockHeight: 12,
		},
		{
			TxId:        txId2,
			Sender:      "123",
			UserAddr:    ContractTxUser1,
			BlockHeight: 12,
		},
	}

	type args struct {
		chainId      string
		transactions []*db.Transaction
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
				transactions: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteBlackTransaction(tt.args.chainId, tt.args.transactions); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBlackTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteTransactionByTxId(t *testing.T) {
	_, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		txIds   []string
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
				txIds: []string{
					txId1,
					txId2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteTransactionByTxId(tt.args.chainId, tt.args.txIds); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTransactionByTxId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetBlackTxInfoByTxId(t *testing.T) {
	insertList, err := insertBlackTransactionsTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		txId    string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.BlackTransaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				txId:    txId1,
			},
			want:    insertList[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBlackTxInfoByTxId(tt.args.chainId, tt.args.txId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlackTxInfoByTxId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.BlackTransaction{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("BatchQueryTxList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetLatestTxList(t *testing.T) {
	insertList, err := insertTxTest()
	if err != nil {
		return
	}

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
			want: []*db.Transaction{
				insertList[1],
				insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestTxList(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestTxList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Transaction{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("BatchQueryTxList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetSafeWordTransactionList(t *testing.T) {
	_, err := insertTxTest()
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
		want    []*db.Transaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    []*db.Transaction{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSafeWordTransactionList(tt.args.chainId, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSafeWordTransactionList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSafeWordTransactionList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTransactionByTxId(t *testing.T) {
	insertList, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		txId    string
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Transaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				txId:    txId1,
			},
			want:    insertList[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTransactionByTxId(tt.args.txId, tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionByTxId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Transaction{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetTransactionByTxId() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetTransactionList(t *testing.T) {
	insertList, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		offset       int
		limit        int
		txStatus     int
		contractName string
		blockHash    string
		startTime    int64
		endTime      int64
		txId         string
		senders      []string
		userAddrs    []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Transaction
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
			want: []*db.Transaction{
				insertList[1],
				insertList[0],
			},
			want1:   int64(len(insertList)),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetTransactionListCount(tt.args.chainId, tt.args.txId, tt.args.contractName, tt.args.blockHash,
				tt.args.startTime, tt.args.endTime, tt.args.txStatus, tt.args.senders, tt.args.userAddrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetTransactionNumByRange(t *testing.T) {
	_, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		userAddr  string
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
				chainId:  ChainID,
				userAddr: ContractTxUser2,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTransactionNumByRange(tt.args.chainId, tt.args.userAddr, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionNumByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTransactionNumByRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTxInfoByBlockHeight(t *testing.T) {
	insertList, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId     string
		blockHeight []int64
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
				chainId:     ChainID,
				blockHeight: []int64{12, 13},
			},
			want:    insertList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTxInfoByBlockHeight(tt.args.chainId, tt.args.blockHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTxInfoByBlockHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Transaction{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetTxInfoByBlockHeight() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetTxListNumByRange(t *testing.T) {
	_, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		startTime int64
		endTime   int64
		interval  int64
	}
	tests := []struct {
		name    string
		args    args
		want    map[int64]int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:   ChainID,
				startTime: timestamp1 - 3600,
				endTime:   timestamp2 + 3600,
				interval:  3600,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTxListNumByRange(tt.args.chainId, tt.args.startTime, tt.args.endTime, tt.args.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTxListNumByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("GetTxListNumByRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTxNumByContractName(t *testing.T) {
	_, err := insertTxTest()
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
		want    int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractName: contractName1,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTxNumByContractName(tt.args.chainId, tt.args.contractName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTxNumByContractName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTxNumByContractName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertBlackTransactions(t *testing.T) {
	insertList := []*db.BlackTransaction{
		{
			TxId:         txId1,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp1,
			ContractName: ContractName1,
		},
		{
			TxId:         txId2,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp2,
			ContractName: ContractName2,
		},
	}

	type args struct {
		chainId      string
		transactions []*db.BlackTransaction
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
				transactions: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertBlackTransactions(tt.args.chainId, tt.args.transactions); (err != nil) != tt.wantErr {
				t.Errorf("InsertBlackTransactions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertTransactions(t *testing.T) {
	insertList := []*db.Transaction{
		{
			TxId:         txId1,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp1,
			ContractName: ContractName1,
		},
		{
			TxId:         txId2,
			Sender:       "123",
			UserAddr:     ContractTxUser2,
			BlockHeight:  12,
			Timestamp:    timestamp2,
			ContractName: ContractName2,
		},
	}

	type args struct {
		chainId      string
		transactions []*db.Transaction
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
				transactions: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertTransactions(tt.args.chainId, tt.args.transactions); (err != nil) != tt.wantErr {
				t.Errorf("InsertTransactions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateTransactionBak(t *testing.T) {
	insertList, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId     string
		transaction *db.Transaction
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
				transaction: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateTransactionBak(tt.args.chainId, tt.args.transaction); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTransactionBak() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateTransactionContractName(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}

	_, err = insertTxTest()
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
				chainId:  ChainID,
				contract: contractInfo1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateTransactionContractName(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTransactionContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
