package sync

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"testing"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

func TestBuildChainInfo(t *testing.T) {
	nodeList := []*config.NodeInfo{
		{
			Addr:        "123456",
			OrgCA:       "12345678",
			TLSHostName: "1234569",
			Tls:         true,
		},
	}

	nodeJson, _ := json.Marshal(nodeList)
	subscribeChain := &db.Subscribe{
		ChainId:    "chain1",
		OrgId:      "123",
		UserKey:    "345",
		UserCert:   "567",
		AuthType:   "567",
		HashType:   "5678",
		NodeList:   string(nodeJson),
		Status:     0,
		NodeCACert: "8888",
	}

	type args struct {
		subscribeChain *db.Subscribe
	}
	tests := []struct {
		name string
		args args
		want *config.ChainInfo
	}{
		{
			name: "test case 1",
			args: args{
				subscribeChain: subscribeChain,
			},
			want: &config.ChainInfo{
				ChainId:   "chain1",
				AuthType:  "567",
				OrgId:     "123",
				HashType:  "5678",
				NodesList: nodeList,
				UserInfo: &config.UserInfo{
					UserKey:  "345",
					UserCert: "567",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildChainInfo(tt.args.subscribeChain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildChainInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelayUpdateOperation(t *testing.T) {
	ctx := context.Background()
	//blockWaitUpdateCh 创建一个容量为 20 的通道来存储处理完成，等待异步更新的区块数据
	blockWaitUpdateCh := make(chan *BlockWaitUpdate, config.BlockWaitUpdateWorkerCount)
	//blockListenErrCh 创建一个错误通道来接收 子线程 的错误
	blockListenErrCh := make(chan error)
	// 将处理完成的结果写入 blockWaitUpdateCh
	resultData := &BlockWaitUpdate{
		ChainId:     "chain1",
		BlockHeight: 10,
	}
	blockWaitUpdateCh <- resultData

	type args struct {
		blockWaitUpdateCh chan *BlockWaitUpdate
		errCh             chan<- error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test case 1",
			args: args{
				blockWaitUpdateCh: blockWaitUpdateCh,
				errCh:             blockListenErrCh,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go DelayUpdateOperation(ctx, tt.args.blockWaitUpdateCh, tt.args.errCh)
			// 使用 select 语句等待错误或上下文取消
			select {
			case errCh := <-blockListenErrCh:
				t.Errorf("DelayUpdateOperation() error = %v, wantErr %v", errCh, errCh)
				close(blockWaitUpdateCh)
				close(blockListenErrCh)
			case <-time.After(2 * time.Second):
				// 等待5秒后关闭通道
				close(blockWaitUpdateCh)
				close(blockListenErrCh)
			}
		})
	}
}

func TestParallelParseBlockWork(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJson_1704945589.json")
	if blockInfo == nil || len(blockInfo.RwsetList) == 0 {
		return
	}
	chainId := blockInfo.Block.Header.ChainId
	blockHeight := int64(blockInfo.Block.Header.BlockHeight)
	setMaxHeight(chainId, blockHeight)

	hashType := "your_hash_type"

	// 创建通道和 WaitGroup
	blockInfoCh := make(chan *common.BlockInfo, 1)
	dataSaveCh := make(chan *DataSaveToDB, 1)
	errCh := make(chan error, 1)
	wg := &sync.WaitGroup{}

	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	type args struct {
		ctx         context.Context
		wg          *sync.WaitGroup
		hashType    string
		blockInfoCh chan *common.BlockInfo
		dataSaveCh  chan *DataSaveToDB
		errCh       chan<- error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test case 1",
			args: args{
				ctx:         ctx,
				wg:          wg,
				hashType:    hashType,
				blockInfoCh: blockInfoCh,
				dataSaveCh:  dataSaveCh,
				errCh:       errCh,
			},
		},
	}
	go func() {
		time.Sleep(10 * time.Second)
		cancel()
		close(blockInfoCh)
		close(dataSaveCh)
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 在一个新的 Goroutine 中运行 ParallelParseBlockWork
			wg.Add(1)
			go ParallelParseBlockWork(tt.args.ctx, tt.args.wg, tt.args.hashType, tt.args.blockInfoCh, tt.args.dataSaveCh, tt.args.errCh)
			// 发送测试数据到 blockInfoCh
			blockInfoCh <- blockInfo
			// 检查结果
			select {
			case err := <-errCh:
				t.Errorf("ParallelParseBlockWork() error = %v", err)
			case dataSave := <-dataSaveCh:
				if dataSave == nil {
					// 验证 dataSave 的内容
					t.Errorf("ParallelParseBlockWork() dataSave = %v", dataSave)
				}
			}
		})
	}
}

func Test_startSubscribeLockTicker(t *testing.T) {
	type args struct {
		sdkClient *SdkClient
		lockKey   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test case 1",
			args: args{
				sdkClient: &SdkClient{
					ChainId: ChainId1,
				},
				lockKey: "1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			go startSubscribeLockTicker(ctx, tt.args.sdkClient, tt.args.lockKey)
			time.Sleep(1000)
			// 在测试结束时取消上下文
			defer cancel()
		})
	}
}
