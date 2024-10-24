package saveTasks

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"context"
	"fmt"
	"sync"
	"time"

	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"
)

// TaskFunc 是一个类型别名，表示任务函数的类型。
type TaskFunc func(...interface{}) error

// Task 任务
type Task struct {
	Name     string
	Function TaskFunc
	Args     []interface{}
}

// RetrySleepTime 重试等待时间
func RetrySleepTime(retryCount int) int {
	if retryCount < 3 {
		return 0
	}

	return retryCount - 3
}

// WithRetry 执行任务，失败后重试
func WithRetry(task func() error, logFuncName string, errCh chan<- error) {
	retryCount := 0
	for {
		err := task()
		if err == nil {
			break
		}
		retryCount++
		log.Errorf("dealBlockData %v-[%d] failed, err:%v", logFuncName, retryCount, err)
		if retryCount > config.MaxRetryCount {
			errCh <- fmt.Errorf("dealBlockData Error: %v, Retry count: %d", err, retryCount)
			break
		}
		//重试总失败就先等待一下
		sleepTime := time.Duration(RetrySleepTime(retryCount))
		time.Sleep(time.Second * sleepTime)
	}
}

// ExecuteTaskWithRetry 会执行给定的任务，并在出错时重试。
// 任务将继续重试，直到成功或达到 maxRetryCount 为止。
func ExecuteTaskWithRetry(ctx context.Context, wg *sync.WaitGroup, task Task, retryCountMap *sync.Map,
	errCh chan<- error) {
	defer wg.Done()
	for {
		err := task.Function(task.Args...)
		if err == nil {
			return
		}
		retryCount, _ := retryCountMap.LoadOrStore(task.Name, 0)
		retryCountMap.Store(task.Name, retryCount.(int)+1)
		if retryCount.(int) >= config.MaxRetryCount {
			// 将错误发送到错误通道
			errCh <- fmt.Errorf("ExecuteTaskWithRetry task %s failed err: %v", task.Name, err)
			return
		}

		select {
		case <-ctx.Done():
			// 上下文已取消
			return
		default:
			// 继续重试
			log.Errorf("ExecuteTaskWithRetry task[%v] run failed err:%v", task.Name, err)
			//重试总失败就先等待一下
			sleepTime := time.Duration(RetrySleepTime(retryCount.(int)))
			time.Sleep(time.Second * sleepTime)
		}
	}
}

// TaskSaveTransactionsToDB 保存交易数据
func TaskSaveTransactionsToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveTransactionsToDB: expected string for args[0], got %T", args[0])
	}
	transactions, ok := args[1].(map[string]*db.Transaction)
	if !ok {
		return fmt.Errorf("TaskSaveTransactionsToDB: expected string for args[1], got %T", args[1])
	}
	upgradeContractTxs, ok := args[2].([]*db.UpgradeContractTransaction)
	if !ok {
		return fmt.Errorf("TaskSaveTransactionsToDB: expected string for args[2], got %T", args[2])
	}
	err := SaveTransactionsToDB(chainId, transactions)
	if err != nil {
		return err
	}
	err = SaveUpgradeContractTxToDB(chainId, upgradeContractTxs)
	return err
}

// TaskSaveUserToDB 保存user数据
func TaskSaveUserToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveUserToDB: expected string for args[0], got %T", args[0])
	}

	users, ok := args[1].(map[string]*db.User)
	if !ok {
		return fmt.Errorf("TaskSaveUserToDB: expected string for args[1], got %T", args[1])
	}

	return SaveUserToDB(chainId, users)
}

// TaskSaveContractToDB 保存合约数据
func TaskSaveContractToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[0], got %T", args[0])
	}
	insertContracts, ok := args[1].([]*db.Contract)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[1], got %T", args[1])
	}
	updateContracts, ok := args[2].([]*db.Contract)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[1], got %T", args[1])
	}

	return SaveContractToDB(chainId, insertContracts, updateContracts)
}

// TaskSaveStandardContractToDB 保存合约数据
func TaskSaveStandardContractToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[0], got %T", args[0])
	}
	fungibleContracts, ok := args[1].([]*db.FungibleContract)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[1], got %T", args[1])
	}
	nonFungibleContracts, ok := args[2].([]*db.NonFungibleContract)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[1], got %T", args[1])
	}
	idaContracts, ok := args[3].([]*db.IDAContract)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[1], got %T", args[1])
	}

	return SaveStandardContractToDB(chainId, fungibleContracts, nonFungibleContracts, idaContracts)
}

// TaskEvidenceContractToDB 保存存证合约数据
func TaskEvidenceContractToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[0], got %T", args[0])
	}
	contractList, ok := args[1].([]*db.EvidenceContract)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[2], got %T", args[2])
	}

	return SaveEvidenceContractToDB(chainId, contractList)
}

// TaskContractEventsToDB 保存topic事件
func TaskContractEventsToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[0], got %T", args[0])
	}
	dbContractEvent, ok := args[1].([]*db.ContractEvent)
	if !ok {
		return fmt.Errorf("TaskContractEventsToDB: expected string for args[1], got %T", args[1])
	}

	return SaveContractEventsToDB(chainId, dbContractEvent)
}

// TaskGasRecordToDB 保存gas数据
func TaskGasRecordToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[0], got %T", args[0])
	}
	gasRecords, ok := args[1].([]*db.GasRecord)
	if !ok {
		return fmt.Errorf("TaskGasRecordToDB: expected string for args[1], got %T", args[1])
	}
	return SaveGasRecordToDB(chainId, gasRecords)
}

// TaskSaveChainConfig 保存链配置信息
func TaskSaveChainConfig(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[0], got %T", args[0])
	}
	chainConfigList, ok := args[1].([]*pbConfig.ChainConfig)
	if !ok {
		return fmt.Errorf("TaskSaveChainConfig: expected string for args[1], got %T", args[1])
	}

	return SaveChainConfig(chainId, chainConfigList)
}

func TaskSaveRelayCrossChainToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveRelayCrossChainToDB: expected string for args[0], got %T", args[0])
	}
	crossChainResult, ok := args[1].(*db.CrossChainResult)
	if !ok {
		return fmt.Errorf("TaskSaveRelayCrossChainToDB: expected string for args[1], got %T", args[1])
	}

	return SaveRelayCrossChainToDB(chainId, crossChainResult)
}

//-------异步更新--------

// TaskUpdateTxBlackToDB 更新交易黑名单数据
func TaskUpdateTxBlackToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[0], got %T", args[0])
	}
	updateTxBlack, ok := args[1].(*db.UpdateTxBlack)
	if !ok {
		return fmt.Errorf("TaskUpdateTxBlackToDB: expected string for args[1], got %T", args[1])
	}

	return UpdateTxBlackToDB(chainId, updateTxBlack)
}

// TaskUpdateContractResult 更新合约数据
func TaskUpdateContractResult(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskUpdateContractResult: expected string for args[0], got %T", args[0])
	}
	contractResult, ok := args[1].(*db.GetContractResult)
	if !ok {
		return fmt.Errorf("TaskUpdateContractResult: expected string for args[1], got %T", args[1])
	}

	//更新合约交易量
	err := UpdateContractTxNum(chainId, contractResult.UpdateContractTxEventNum)
	if err != nil {
		return err
	}

	//更新IDA合约资产数量
	err = UpdateIDAContract(chainId, contractResult.UpdateIdaContract)
	if err != nil {
		return err
	}

	return UpdateContract(chainId, contractResult.IdentityContract)
}

// TaskInsertFungibleTransferToDB 保留流转数据
func TaskInsertFungibleTransferToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskInsertFungibleTransferToDB: expected string for args[0], got %T", args[0])
	}
	transferList, ok := args[1].([]*db.FungibleTransfer)
	if !ok {
		return fmt.Errorf("TaskInsertFungibleTransferToDB: expected string for args[1], got %T", args[1])
	}

	return InsertFungibleTransferToDB(chainId, transferList)
}

// TaskInsertNonFungibleTransferToDB 保留流转数据
func TaskInsertNonFungibleTransferToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskInsertNonFungibleTransferToDB: expected string for args[0], got %T", args[0])
	}
	transferList, ok := args[1].([]*db.NonFungibleTransfer)
	if !ok {
		return fmt.Errorf("TaskInsertNonFungibleTransferToDB: expected string for args[1], got %T", args[1])
	}
	return InsertNonFungibleTransferToDB(chainId, transferList)
}

// TaskSaveAccountListToDB 保存账户数据
func TaskSaveAccountListToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveContractResult: expected string for args[0], got %T", args[0])
	}
	updateAccountResult, ok := args[1].(*db.UpdateAccountResult)
	if !ok {
		return fmt.Errorf("TaskSaveAccountListToDB: expected string for args[1], got %T", args[1])
	}

	return SaveAccountToDB(chainId, updateAccountResult)
}

// TaskSaveTokenResultToDB 保存非同质化token
func TaskSaveTokenResultToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveContractResult: expected string for args[0], got %T", args[0])
	}
	tokenResult, ok := args[1].(*db.TokenResult)
	if !ok {
		return fmt.Errorf("TaskSaveTokenResultToDB: expected string for args[1], got %T", args[1])
	}

	return SaveNonFungibleToken(chainId, tokenResult)
}

func TaskSaveGasToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveGasToDB: expected string for args[0], got %T", args[0])
	}
	insertGasList, ok := args[1].([]*db.Gas)
	if !ok {
		return fmt.Errorf("TaskSaveGasToDB: expected string for args[1], got %T", args[1])
	}
	updateGasList, ok := args[2].([]*db.Gas)
	if !ok {
		return fmt.Errorf("TaskSaveGasToDB: expected string for args[1], got %T", args[1])
	}

	err := InsertGasToDB(chainId, insertGasList)
	if err != nil {
		return err
	}
	err = UpdateGasToDB(chainId, updateGasList)
	if err != nil {
		return err
	}
	return nil
}

func TaskSaveFungibleContractResult(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveFungibleContractResult: expected string for args[0], got %T", args[0])
	}
	contractResult, ok := args[1].(*db.GetContractResult)
	if !ok {
		return fmt.Errorf("TaskSaveFungibleContractResult: expected string for args[1], got %T", args[1])
	}

	return SaveFungibleContractResult(chainId, contractResult)
}

func TaskSavePositionToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSavePositionToDB: expected string for args[0], got %T", args[0])
	}
	blockPosition, ok := args[1].(*db.BlockPosition)
	if !ok {
		return fmt.Errorf("TaskSavePositionToDB: expected string for args[1], got %T", args[1])
	}

	return SavePositionToDB(chainId, blockPosition)
}

func TaskCrossSubChainCrossToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSavePositionToDB: expected string for args[0], got %T", args[0])
	}
	insertList, ok := args[1].([]*db.CrossSubChainCrossChain)
	if !ok {
		return fmt.Errorf("TaskCrossSubChainCrossToDB: expected string for args[1], got %T", args[1])
	}
	updateList, ok := args[2].([]*db.CrossSubChainCrossChain)
	if !ok {
		return fmt.Errorf("TaskCrossSubChainCrossToDB: expected string for args[2], got %T", args[2])
	}

	return SaveCrossSubChainCrossToDB(chainId, insertList, updateList)
}

func TaskCrossUpdateSubChainTxNumToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskCrossUpdateSubChainTxNumToDB: expected string for args[0], got %T", args[0])
	}
	updateList, ok := args[1].([]*db.CrossSubChainData)
	if !ok {
		return fmt.Errorf("TaskCrossUpdateSubChainTxNumToDB: expected string for args[1], got %T", args[1])
	}

	return UpdateCrossSubChainData(chainId, updateList)
}

func TaskSaveIDAAssetDataToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveIDAAssetDataToDB: expected string for args[0], got %T", args[0])
	}
	idaAssetsDataDB, ok := args[1].(*db.IDAAssetsDataDB)
	if !ok {
		return fmt.Errorf("TaskSaveIDAAssetDataToDB: expected string for args[1], got %T", args[1])
	}

	return SaveIDAAssetDataToDB(chainId, idaAssetsDataDB)
}

func TaskUpdateIDAAssetDataToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskUpdateIDAAssetDataToDB: expected string for args[0], got %T", args[0])
	}
	idaAssetsDataDB, ok := args[1].(*db.IDAAssetsUpdateDB)
	if !ok {
		return fmt.Errorf("TaskUpdateIDAAssetDataToDB: expected string for args[1], got %T", args[1])
	}

	return UpdateIDAAssetDataToDB(chainId, idaAssetsDataDB)
}
