/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
	"sync"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"github.com/panjf2000/ants/v2"
)

const MD5Str = "md5:"

// ParallelParseTransactions
//
//	@Description: 并发解析所有交易数据
//	@param blockInfo
//	@param hashType
//	@param dealResult
//	@return var
//	@return err
func ParallelParseTransactions(blockInfo *common.BlockInfo, hashType string, dealResult *RealtimeDealResult) (
	*RealtimeDealResult, error) {
	var (
		goRoutinePool *ants.Pool
		mutx          sync.Mutex
		err           error
	)
	errChan := make(chan error, 10)
	if goRoutinePool, err = ants.NewPool(10, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return dealResult, err
	}
	chainId := blockInfo.Block.Header.ChainId
	defer goRoutinePool.Release()
	var wg sync.WaitGroup
	for i, tx := range blockInfo.Block.Txs {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		errSub := goRoutinePool.Submit(func(i int, blockInfo *common.BlockInfo, txInfo *common.Transaction) func() {
			return func() {
				defer wg.Done()

				//构造ContractEvent
				tempContractEvents := DealContractEvents(txInfo)

				// 计算账户sender
				userResult, senderErr := GetSenderAndPayerUser(chainId, hashType, txInfo)
				if senderErr != nil {
					log.Errorf("【Realtime deal】ParallelParseTransactions get User err:%v", senderErr)
				}

				//构造GasRecord
				tempGasRecords, gasErr := buildGasRecord(txInfo, userResult)
				if gasErr != nil {
					errChan <- gasErr
					return
				}

				//构造Transaction数据
				transaction, tranErr := buildTransaction(i, blockInfo, txInfo, userResult)
				if tranErr != nil {
					errChan <- tranErr
					return
				}

				//构造合约升级数据
				upgradeTx := buildUpgradeContractTransaction(txInfo, transaction)

				mutx.Lock()         // 锁定互斥锁
				defer mutx.Unlock() // 使用 defer 确保互斥锁被解锁
				//交易记录
				if transaction.TxId != "" {
					dealResult.Transactions[transaction.TxId] = transaction
				}
				if upgradeTx != nil {
					dealResult.UpgradeContractTx = append(dealResult.UpgradeContractTx, upgradeTx)
				}

				//合约event
				if len(tempContractEvents) > 0 {
					dealResult.ContractEvents = append(dealResult.ContractEvents, tempContractEvents...)
				}
				//gas记录
				if len(tempGasRecords) > 0 {
					dealResult.GasRecordList = append(dealResult.GasRecordList, tempGasRecords...)
				}
				//user列表
				if userResult.SenderUserAddr != "" {
					// 检查userMap中是否已经存在具有相同UserAddr的userInfo
					if _, ok := dealResult.UserList[userResult.SenderUserAddr]; !ok {
						userInfo := &db.User{
							UserId:    userResult.SenderUserId,
							UserAddr:  userResult.SenderUserAddr,
							Role:      userResult.SenderRole,
							OrgId:     userResult.SenderOrgId,
							Timestamp: txInfo.Payload.Timestamp,
						}
						dealResult.UserList[userResult.SenderUserAddr] = userInfo
					}
				}

				if userResult.PayerUserAddr != "" {
					// 检查userMap中是否已经存在具有相同UserAddr的userInfo
					if _, ok := dealResult.UserList[userResult.PayerUserAddr]; !ok {
						userInfo := &db.User{
							UserId:    userResult.PayerUserId,
							UserAddr:  userResult.PayerUserAddr,
							Timestamp: txInfo.Payload.Timestamp,
						}
						dealResult.UserList[userResult.PayerUserAddr] = userInfo
					}
				}
			}
		}(i, blockInfo, tx))
		if errSub != nil {
			log.Errorf("ParallelParseTransactions submit Failed : %v", err)
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return dealResult, err
	}
	return dealResult, nil
}

func buildUpgradeContractTransaction(txInfo *common.Transaction,
	transaction *db.Transaction) *db.UpgradeContractTransaction {
	if txInfo == nil || transaction == nil || transaction.TxId == "" {
		return nil
	}

	//非合约升级数据不需要处理
	isContractTx := IsContractTxByName(txInfo.Payload.ContractName, txInfo.Payload.Method)
	if !isContractTx {
		return nil
	}

	//构造合约升级交易数据
	upgradeContractTransaction := &db.UpgradeContractTransaction{
		TxId:                transaction.TxId,
		SenderOrgId:         transaction.SenderOrgId,
		Sender:              transaction.Sender,
		UserAddr:            transaction.UserAddr,
		BlockHeight:         transaction.BlockHeight,
		BlockHash:           transaction.BlockHash,
		Timestamp:           transaction.Timestamp,
		TxStatusCode:        transaction.TxStatusCode,
		ContractResultCode:  transaction.ContractResultCode,
		ContractRuntimeType: transaction.ContractRuntimeType,
		ContractName:        transaction.ContractName,
		ContractNameBak:     transaction.ContractNameBak,
		ContractAddr:        transaction.ContractAddr,
		ContractVersion:     transaction.ContractVersion,
		ContractType:        transaction.ContractType,
	}

	return upgradeContractTransaction
}

func buildTransaction(i int, blockInfo *common.BlockInfo, txInfo *common.Transaction, userResult *db.SenderPayerUser) (
	*db.Transaction, error) {
	payload := txInfo.Payload
	contractNameAddr := payload.ContractName

	//构造交易数据
	transaction := &db.Transaction{
		TxId:               payload.TxId,
		TxIndex:            i + 1,
		TxType:             payload.TxType.String(),
		BlockHeight:        int64(blockInfo.Block.Header.BlockHeight),
		BlockHash:          hex.EncodeToString(blockInfo.Block.Header.BlockHash),
		ContractMessage:    txInfo.Result.ContractResult.Message,
		GasUsed:            txInfo.Result.ContractResult.GasUsed,
		Sequence:           payload.Sequence,
		ContractResult:     txInfo.Result.ContractResult.Result,
		ContractResultCode: txInfo.Result.ContractResult.Code,
		ExpirationTime:     payload.ExpirationTime,
		RwSetHash:          hex.EncodeToString(txInfo.Result.RwSetHash),
		Timestamp:          payload.Timestamp,
		TxStatusCode:       txInfo.Result.Code.String(),
		ContractMethod:     payload.Method,
	}

	for _, parameter := range payload.Parameters {
		switch parameter.Key {
		case syscontract.InitContract_CONTRACT_NAME.String():
			contractNameAddr = string(parameter.Value)
		case syscontract.InitContract_CONTRACT_VERSION.String():
			transaction.ContractVersion = string(parameter.Value)
		case syscontract.InitContract_CONTRACT_RUNTIME_TYPE.String():
			transaction.ContractRuntimeType = string(parameter.Value)
		case syscontract.InitContract_CONTRACT_BYTECODE.String():
			parameter.Value = []byte(MD5Str + MD5(string(parameter.Value)))
		case syscontract.UpgradeContract_CONTRACT_BYTECODE.String():
			parameter.Value = []byte(MD5Str + MD5(string(parameter.Value)))
		}
	}

	transaction.ContractName = contractNameAddr
	transaction.ContractNameBak = contractNameAddr
	parametersBytes, err := json.Marshal(payload.Parameters)
	if err == nil {
		transaction.ContractParameters = string(parametersBytes)
	}

	//解析读写集
	transaction.ReadSet, transaction.WriteSet = buildReadWriteSet(blockInfo.RwsetList[i])
	if userResult != nil {
		transaction.Sender = userResult.SenderUserId
		transaction.SenderOrgId = userResult.SenderOrgId
		transaction.UserAddr = userResult.SenderUserAddr
		transaction.PayerAddr = userResult.PayerUserAddr
	}

	if len(txInfo.Endorsers) > 0 {
		endorsementBytes, _ := json.Marshal(txInfo.Endorsers)
		transaction.Endorsement = string(endorsementBytes)
	}

	if len(txInfo.Result.ContractResult.ContractEvent) > 0 {
		eventList := make([]config.RwSet, 0)
		for k, event := range txInfo.Result.ContractResult.ContractEvent {
			eventList = append(eventList, config.RwSet{
				Index:        k,
				ContractName: event.ContractName,
				Key:          event.Topic,
				Value:        strings.Join(event.EventData, ","),
			})
		}
		eventByte, _ := json.Marshal(eventList)
		transaction.Event = string(eventByte)
	}

	return transaction, nil
}

// BuildLatestTxListCache 缓存交易信息
func BuildLatestTxListCache(chainId string, txMap map[string]*db.Transaction) {
	// 从缓存中获取交易列表
	txList, err := dbhandle.GetLatestTxListCache(chainId)
	if len(txList) == 0 || err != nil {
		//缓存可能丢失
		txList, _ = dbhandle.GetLatestTxList(chainId)
		if len(txList) == 0 {
			for _, txInfo := range txMap {
				txList = append(txList, txInfo)
			}
		}
	} else {
		//缓存存在,缓存数据加入新数据
		for _, txInfo := range txMap {
			txList = append(txList, txInfo)
		}
	}

	// 根据 blockHeight 和 txInfo.TxIndex 排序交易列表
	sort.Slice(txList, func(i, j int) bool {
		if txList[i].BlockHeight == txList[j].BlockHeight {
			return txList[i].TxIndex < txList[j].TxIndex
		}
		return txList[i].BlockHeight > txList[j].BlockHeight
	})

	// 保留最新的 10 条交易数据
	if len(txList) > 10 {
		txList = txList[:10]
	}

	// 缓存交易信息
	dbhandle.SetLatestTxListCache(chainId, txList)
}

// BuildOverviewTxTotalCache
//
//	@Description: 缓存首页交易总量
//	@param chainId
//	@param transactions
func BuildOverviewTxTotalCache(chainId string, transactions map[string]*db.Transaction) {
	if len(transactions) == 0 {
		return
	}

	txCount, err := dbhandle.GetTotalTxNum(chainId)
	if err != nil {
		log.Errorf("BuildOverviewTxTotalCache GetTotalTxNum err :%v", err)
		return
	}

	txTotal := txCount + int64(len(transactions))
	dbhandle.SetTotalTxNumCache(chainId, txTotal)
}

// buildReadWriteSet 解析读写集，数据过长也不在进行截断
func buildReadWriteSet(rwSetList *common.TxRWSet) (string, string) {
	readList := make([]config.RwSet, 0)
	for j, read := range rwSetList.TxReads {
		value := make([]byte, len(read.Value))
		copy(value, read.Value)
		valueStr := string(value)
		if strings.HasPrefix(string(read.Key), "ContractByteCode:") {
			valueStr = MD5Str + MD5(string(read.Value))
		}
		readList = append(readList, config.RwSet{
			Index:        j,
			Key:          string(read.Key),
			Value:        valueStr,
			ContractName: read.ContractName,
		})
	}

	writeList := make([]config.RwSet, 0)
	for j, write := range rwSetList.TxWrites {
		value := make([]byte, len(write.Value))
		copy(value, write.Value)
		valueStr := string(value)
		if strings.HasPrefix(string(write.Key), "ContractByteCode:") {
			valueStr = MD5Str + MD5(string(write.Value))
		}

		writeList = append(writeList, config.RwSet{
			Index:        j,
			Key:          string(write.Key),
			Value:        valueStr,
			ContractName: write.ContractName,
		})
	}
	readByte, _ := json.Marshal(readList)
	writeByte, _ := json.Marshal(writeList)
	return string(readByte), string(writeByte)
}
