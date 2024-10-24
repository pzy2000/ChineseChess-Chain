/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"sync"

	"github.com/google/uuid"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/panjf2000/ants/v2"
)

const (
	//BusinessTypeRecharge gas充值
	BusinessTypeRecharge = 1
	//BusinessTypeConsume gas消费
	BusinessTypeConsume = 2
)

// buildGasRecord
//
//	@Description: 根据交易数据构造gas消耗列表
//	@param txInfo 交易列表
//	@param userResult 交易用户
//	@return []*db.GasRecord gas消耗
//	@return error
func buildGasRecord(txInfo *common.Transaction, userResult *db.SenderPayerUser) ([]*db.GasRecord, error) {
	gasRecords := make([]*db.GasRecord, 0)
	payload := txInfo.Payload
	txId := payload.TxId

	//gas充值
	if payload.Method == syscontract.GasAccountFunction_RECHARGE_GAS.String() {
		if txInfo.Result.Code != common.TxStatusCode_SUCCESS {
			return gasRecords, nil
		}

		req := &syscontract.RechargeGasReq{}
		for _, parameter := range payload.Parameters {
			if parameter.Key == utils.KeyGasBatchRecharge {
				err := req.Unmarshal(parameter.Value)
				if err != nil {
					return gasRecords, err
				}
				break
			}
		}

		if len(req.BatchRechargeGas) > 0 {
			for i, gasReq := range req.BatchRechargeGas {
				newUUID := uuid.New().String()
				gasInfo := &db.GasRecord{
					ID:           newUUID,
					Address:      gasReq.Address,
					GasAmount:    gasReq.GasAmount,
					BusinessType: BusinessTypeRecharge,
					Timestamp:    payload.Timestamp,
					TxId:         txId,
					GasIndex:     i + 1,
				}
				gasRecords = append(gasRecords, gasInfo)
			}
		}
	} else {
		//gas消费
		newUUID := uuid.New().String()
		gasInfo := &db.GasRecord{
			ID:           newUUID,
			TxId:         txId,
			GasIndex:     1,
			GasAmount:    int64(txInfo.Result.ContractResult.GasUsed),
			BusinessType: BusinessTypeConsume,
			Timestamp:    payload.Timestamp,
		}

		if userResult != nil {
			if userResult.PayerUserAddr != "" {
				//PayerUserAddr代付
				gasInfo.Address = userResult.PayerUserAddr
			} else {
				gasInfo.Address = userResult.SenderUserAddr
			}
		}
		if gasInfo.GasAmount == 0 || gasInfo.Address == "" {
			return gasRecords, nil
		}
		gasRecords = append(gasRecords, gasInfo)
	}

	return gasRecords, nil
}

// buildGasInfo
//
//	@Description: 根据gas消耗列表，和DB中gas余额计算新的gas余额
//	@param gasRecords gas消耗
//	@param gasInfoList DB中gas余额
//	@param minHeight 批量处理最小高度，用作版本号，避免重复计算
//	@return []*db.Gas 新增gas用户
//	@return []*db.Gas 更新gas用户
func buildGasInfo(gasRecords []*db.GasRecord, gasInfoList []*db.Gas, minHeight int64) ([]*db.Gas, []*db.Gas) {
	var (
		gasUseAmount   = make(map[string]int64)
		gasTotalAmount = make(map[string]int64)
		addrMap        = make(map[string]string, 0)
		insertGas      = make([]*db.Gas, 0)
		updateGas      = make([]*db.Gas, 0)
	)
	for _, gasInfo := range gasRecords {
		addr := gasInfo.Address
		if _, ok := addrMap[addr]; !ok {
			addrMap[addr] = addr
		}

		if gasInfo.BusinessType == BusinessTypeRecharge {
			//gas充值
			if amount, ok := gasTotalAmount[addr]; ok {
				gasTotalAmount[addr] = amount + gasInfo.GasAmount
			} else {
				gasTotalAmount[addr] = gasInfo.GasAmount
			}
		} else {
			//gas消耗
			if amount, ok := gasUseAmount[addr]; ok {
				gasUseAmount[addr] = amount + gasInfo.GasAmount
			} else {
				gasUseAmount[addr] = gasInfo.GasAmount
			}
		}
	}

	if len(addrMap) == 0 {
		return insertGas, updateGas
	}

	// 创建一个映射来存储gasInfoList中的地址
	gasInfoMap := make(map[string]*db.Gas)
	for _, gas := range gasInfoList {
		gasInfoMap[gas.Address] = gas
	}

	for _, addr := range addrMap {
		gas, okMap := gasInfoMap[addr]
		if okMap {
			//数据库存在
			if gas.BlockHeight >= minHeight {
				//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
				continue
			}

			if amount, ok := gasUseAmount[gas.Address]; ok {
				gas.GasUsed = gas.GasUsed + amount
			}
			if amount, ok := gasTotalAmount[gas.Address]; ok {
				gas.GasTotal = gas.GasTotal + amount
			}
			gas.GasBalance = gas.GasTotal - gas.GasUsed
			gas.BlockHeight = minHeight
			updateGas = append(updateGas, gas)
		} else {
			// 将不存在的地址添加到InsertGas
			gas = &db.Gas{
				Address:     addr,
				BlockHeight: minHeight,
			}
			if amount, ok := gasUseAmount[gas.Address]; ok {
				gas.GasUsed = amount
			}
			if amount, ok := gasTotalAmount[gas.Address]; ok {
				gas.GasTotal = amount
			}
			gas.GasBalance = gas.GasTotal - gas.GasUsed
			insertGas = append(insertGas, gas)
		}
	}

	return insertGas, updateGas
}

// GetGasRecord
//
//	@Description: 数据库获取gas记录
//	@param chainId
//	@param txIds 交易ID列表
//	@return []*db.GasRecord  gas消耗
//	@return error
func GetGasRecord(chainId string, txIds []string) ([]*db.GasRecord, error) {
	gasRecords := make([]*db.GasRecord, 0)
	if len(txIds) == 0 {
		return gasRecords, nil
	}
	var (
		goRoutinePool *ants.Pool
		mutx          sync.Mutex
		err           error
	)
	// 将交易分割为大小为10的批次
	batches := ParallelParseBatchWhere(txIds, 100)
	errChan := make(chan error, 10)
	if goRoutinePool, err = ants.NewPool(10, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return gasRecords, err
	}

	defer goRoutinePool.Release()
	var wg sync.WaitGroup
	for _, batch := range batches {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(batch []string) func() {
			return func() {
				defer wg.Done()
				//查询数据
				gasList, eventErr := dbhandle.GetGasRecordByTxIds(chainId, txIds)
				if eventErr != nil {
					errChan <- eventErr
				}
				mutx.Lock()         // 锁定互斥锁
				defer mutx.Unlock() // 使用 defer 确保互斥锁被解锁
				if len(gasList) > 0 {
					gasRecords = append(gasRecords, gasList...)
				}
			}
		}(batch))
		if errSub != nil {
			log.Error("GetGasRecord submit Failed : " + err.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return gasRecords, err
	}

	return gasRecords, nil

}

// buildGasAddrList
//
//	@Description: 计算需要更新的GAS数据
//	@param gasRecords 本次gas消耗列表
//	@return []string 计算后的gas余额
func buildGasAddrList(gasRecords []*db.GasRecord) []string {
	//获取gas余额
	addrMap := make(map[string]string, 0)
	addrList := make([]string, 0)
	for _, gasInfo := range gasRecords {
		addr := gasInfo.Address
		if _, ok := addrMap[addr]; !ok {
			addrMap[addr] = addr
		}
	}
	for _, addr := range addrMap {
		addrList = append(addrList, addr)
	}

	return addrList
}
