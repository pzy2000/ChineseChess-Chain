package sync

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

const ChainId = "testChainId"

func TestBuildLatestTxListCache(t *testing.T) {
	// 准备测试数据
	chainId := ChainId
	txMap := map[string]*db.Transaction{
		"tx1": {
			TxId:        "tx1",
			BlockHeight: 5,
			TxIndex:     2,
		},
		"tx2": {
			TxId:        "tx2",
			BlockHeight: 4,
			TxIndex:     1,
		},
		"tx3": {
			TxId:        "tx3",
			BlockHeight: 6,
			TxIndex:     3,
		},
	}

	// 调用 BuildLatestTxListCache 函数
	BuildLatestTxListCache(chainId, txMap)

	// 从缓存中获取交易列表
	txList, err := dbhandle.GetLatestTxListCache(chainId)
	if err != nil {
		return
	}

	// 检查交易列表长度是否为 3
	if len(txList) != 3 {
		t.Errorf("Expected txList length to be 3, got %d", len(txList))
	}

	// 检查交易列表是否按照预期排序
	expectedOrder := []string{"tx3", "tx1", "tx2"}
	for i, tx := range txList {
		if tx.TxId != expectedOrder[i] {
			t.Errorf("Expected txList order to be %v, got %v", expectedOrder, txList)
		}
	}
}

func TestGetLatestTxListCache(t *testing.T) {
	// 准备测试数据
	chainId := ChainId
	txMap := map[string]*db.Transaction{
		"tx1": {
			TxId:        "tx1",
			BlockHeight: 5,
			TxIndex:     2,
		},
		"tx2": {
			TxId:        "tx2",
			BlockHeight: 4,
			TxIndex:     1,
		},
		"tx3": {
			TxId:        "tx3",
			BlockHeight: 6,
			TxIndex:     3,
		},
	}

	// 调用 BuildLatestTxListCache 函数
	BuildLatestTxListCache(chainId, txMap)

	// 从缓存中获取交易列表
	txList, err := dbhandle.GetLatestTxListCache(chainId)
	if err != nil {
		return
	}

	// 检查交易列表长度是否为 3
	if len(txList) != 3 {
		t.Errorf("Expected txList length to be 3, got %d", len(txList))
	}

	// 检查交易列表是否按照预期排序
	expectedOrder := []string{"tx3", "tx1", "tx2"}
	for i, tx := range txList {
		if tx.TxId != expectedOrder[i] {
			t.Errorf("Expected txList order to be %v, got %v", expectedOrder, txList)
		}
	}
}

func TestParallelParseTransactions(t *testing.T) {
	dealResult := &RealtimeDealResult{
		UserList:     map[string]*db.User{},
		Transactions: map[string]*db.Transaction{},
	}
	blockInfo := getBlockInfoTest("2_blockInfoJsonContractERC20.json")
	if blockInfo == nil {
		return
	}

	type args struct {
		blockInfo  *common.BlockInfo
		hashType   string
		dealResult *RealtimeDealResult
	}
	tests := []struct {
		name    string
		args    args
		want    *RealtimeDealResult
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				blockInfo:  blockInfo,
				hashType:   "1222222",
				dealResult: dealResult,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParallelParseTransactions(tt.args.blockInfo, tt.args.hashType, tt.args.dealResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParallelParseTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got.Transactions) == 0 {
				t.Errorf("ParallelParseTransactions() got = %v, len Transactions == 0", got)
			}
			if len(got.ContractEvents) == 0 {
				t.Errorf("ParallelParseTransactions() got = %v, len ContractEvents == 0", got)
			}
		})
	}
}

func TestSetLatestTxListCache(t *testing.T) {
	// 准备测试数据
	chainId := ChainId
	txList := []*db.Transaction{
		{
			TxId:        "tx1",
			BlockHeight: 5,
			TxIndex:     2,
		},
		{
			TxId:        "tx2",
			BlockHeight: 4,
			TxIndex:     1,
		},
		{
			TxId:        "tx3",
			BlockHeight: 6,
			TxIndex:     3,
		},
	}

	// 调用 BuildLatestTxListCache 函数
	dbhandle.SetLatestTxListCache(chainId, txList)

	// 从缓存中获取交易列表
	txListCache, err := dbhandle.GetLatestTxListCache(chainId)
	if err != nil {
		return
	}

	// 检查交易列表长度是否为 3
	if len(txListCache) != 3 {
		t.Errorf("Expected txList length to be 3, got %d", len(txList))
	}

	// 检查交易列表是否按照预期排序
	expectedOrder := []string{"tx3", "tx1", "tx2"}
	for i, tx := range txListCache {
		if tx.TxId != expectedOrder[i] {
			t.Errorf("Expected txList order to be %v, got %v", expectedOrder, txList)
		}
	}
}

func Test_buildReadWriteSet(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJson_1704945589.json")
	if blockInfo == nil || len(blockInfo.RwsetList) == 0 {
		return
	}
	type args struct {
		rwsetList *common.TxRWSet
	}

	type testStruct struct {
		name string
		args args
	}

	var tests []testStruct
	temp := testStruct{
		name: "test case 1",
		args: args{
			rwsetList: blockInfo.RwsetList[0],
		},
	}
	tests = append(tests, temp)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildReadWriteSet(tt.args.rwsetList)
		})
	}
}

func Test_buildTransaction(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJsonContractERC20.json")
	if blockInfo == nil {
		return
	}
	txInfo := getTxInfoInfoTest("2_txInfoJsonContractERC20.json")
	if txInfo == nil {
		return
	}
	buildTxResult := getBuildTxInfoTest("2_buildTxResult.json")
	if buildTxResult == nil {
		return
	}
	type args struct {
		i          int
		blockInfo  *common.BlockInfo
		txInfo     *common.Transaction
		userResult *db.SenderPayerUser
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Transaction
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				i:          0,
				blockInfo:  blockInfo,
				txInfo:     txInfo,
				userResult: nil,
			},
			want:    buildTxResult,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildTransaction(tt.args.i, tt.args.blockInfo, tt.args.txInfo, tt.args.userResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("buildTransaction() got = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestGetLatestTxListCache1(t *testing.T) {
	// 准备测试数据
	chainId := ChainId
	txMap := map[string]*db.Transaction{
		"tx1": {
			TxId:        "tx1",
			BlockHeight: 5,
			TxIndex:     2,
		},
		"tx2": {
			TxId:        "tx2",
			BlockHeight: 4,
			TxIndex:     1,
		},
		"tx3": {
			TxId:        "tx3",
			BlockHeight: 6,
			TxIndex:     3,
		},
	}

	// 调用 BuildLatestTxListCache 函数
	BuildLatestTxListCache(chainId, txMap)

	// 从缓存中获取交易列表
	txList, err := dbhandle.GetLatestTxListCache(chainId)
	if err != nil {
		return
	}

	// 检查交易列表长度是否为 3
	if len(txList) != 3 {
		t.Errorf("Expected txList length to be 3, got %d", len(txList))
	}

	// 检查交易列表是否按照预期排序
	expectedOrder := []string{"tx3", "tx1", "tx2"}
	for i, tx := range txList {
		if tx.TxId != expectedOrder[i] {
			t.Errorf("Expected txList order to be %v, got %v", expectedOrder, txList)
		}
	}
}
