/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"bytes"
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/utils"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
)

const (
	// SUCCESS sus
	SUCCESS = "SUCCESS"
	// SUCCESS_STATUS c
	SUCCESS_STATUS = "成功"
	// FAIL_STATUS fail
	FAIL_STATUS = "失败"

	//ShowStatusTrue 展示交易信息
	ShowStatusTrue = 0
	//ShowStatusFalse 隐藏交易信息
	ShowStatusFalse = 1

	//SortTypeDesc GetTxNumByTime 时间排序-倒序
	SortTypeDesc = 0
	//SortTypeAsc GetTxNumByTime 时间排序-正序
	SortTypeAsc = 1
)

// GetTxDetailHandler get
type GetTxDetailHandler struct {
}

// Handle deal
func (getTxDetailHandler *GetTxDetailHandler) Handle(ctx *gin.Context) {
	var (
		transaction *db.Transaction
		err         error
	)
	params := entity.BindGetTxDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetTxDetail param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//根据交易id获取交易数据
	transaction, err = dbhandle.GetTransactionByTxId(params.TxId, params.ChainId)
	if err != nil || transaction == nil {
		log.Errorf("GetTxDetail err : %s, TxId:%v", err, params.TxId)
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	var endorsement string
	if transaction.Endorsement != "" && transaction.Endorsement != "null" {
		var endorsers []common.EndorsementEntry
		err := json.Unmarshal([]byte(transaction.Endorsement), &endorsers)
		if err != nil {
			log.Error("transaction.Endorsement json unmarshal err " + err.Error())
		}
		for _, endorser := range endorsers {
			endorsement += endorser.Signer.OrgId + ","
		}
		endorsement = strings.Trim(endorsement, ",")
	}

	showStatus := ShowStatusTrue
	if bytes.Equal(transaction.ContractResultBak, config.ContractResultMsg) ||
		transaction.ContractMessageBak != "" ||
		transaction.ContractParametersBak != "" ||
		transaction.ReadSetBak != "" ||
		transaction.WriteSetBak != "" {
		showStatus = ShowStatusFalse
	}

	//获取账户BNS地址
	var accountAddrs []string
	var addrBns string
	var payerBns string
	if transaction.UserAddr != "" {
		accountAddrs = append(accountAddrs, transaction.UserAddr)
	}
	if transaction.PayerAddr != "" {
		accountAddrs = append(accountAddrs, transaction.PayerAddr)
	}

	accountMap, _ := dbhandle.QueryAccountExists(params.ChainId, accountAddrs)
	if account, ok := accountMap[transaction.UserAddr]; ok {
		addrBns = account.BNS
	}
	if account, ok := accountMap[transaction.PayerAddr]; ok {
		payerBns = account.BNS
	}

	var txStatus int
	if transaction.TxStatusCode != "SUCCESS" {
		txStatus = 1
	}

	txView := &entity.TxDetailView{
		TxId:               transaction.TxId,
		BlockHeight:        transaction.BlockHeight,
		BlockHash:          transaction.BlockHash,
		Sender:             transaction.Sender,
		SenderOrgId:        transaction.SenderOrgId,
		ContractName:       transaction.ContractName,
		ContractNameBak:    transaction.ContractName,
		ContractAddr:       transaction.ContractAddr,
		ContractVersion:    transaction.ContractVersion,
		ContractMessage:    transaction.ContractMessage,
		TxStatusCode:       transaction.TxStatusCode,
		TxStatus:           txStatus,
		ContractResultCode: transaction.ContractResultCode,
		ContractResult:     transaction.ContractResult,
		RwSetHash:          transaction.RwSetHash,
		ContractMethod:     transaction.ContractMethod,
		ContractParameters: transaction.ContractParameters,
		Endorsement:        endorsement,
		TxType:             transaction.TxType,
		Timestamp:          transaction.Timestamp,
		UserAddr:           transaction.UserAddr,
		UserAddrBns:        addrBns,
		ContractRead:       transaction.ReadSet,
		ContractWrite:      transaction.WriteSet,
		GasUsed:            transaction.GasUsed,
		Payer:              transaction.PayerAddr,
		PayerBns:           payerBns,
		Event:              transaction.Event,
		RuntimeType:        transaction.ContractRuntimeType,
		ShowStatus:         showStatus,
	}
	ConvergeDataResponse(ctx, txView, nil)
}

// GetLatestTxListHandler get
type GetLatestTxListHandler struct {
}

// Handle deal
func (getLatestTxListHandler *GetLatestTxListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetTxLatestListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetLatestTxList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//从缓存获取最新的交易数据
	txList, err := dbhandle.GetLatestTxListCache(params.ChainId)
	if err != nil {
		log.Errorf("getTxListFromRedis get redis fail err:%v", err)
	}

	if len(txList) == 0 {
		txList, err = dbhandle.GetLatestTxList(params.ChainId)
		if err != nil {
			ConvergeHandleFailureResponse(ctx, err)
			return
		}
	}

	//获取账户BNS地址
	var accountAddrs []string
	for _, txInfo := range txList {
		accountAddrs = append(accountAddrs, txInfo.UserAddr)
	}
	accountMap, _ := dbhandle.QueryAccountExists(params.ChainId, accountAddrs)

	txViews := arraylist.New()
	for i, txInfo := range txList {
		var status string
		var senderBns string
		if txInfo.TxStatusCode == SUCCESS {
			status = SUCCESS_STATUS
		} else {
			status = FAIL_STATUS
		}

		if account, ok := accountMap[txInfo.UserAddr]; ok {
			senderBns = account.BNS
		}

		latestBlockListView := &entity.LatestTxListView{
			Id:              i + 1,
			TxId:            txInfo.TxId,
			BlockHeight:     txInfo.BlockHeight,
			BlockHash:       txInfo.BlockHash,
			Status:          status,
			ContractName:    txInfo.ContractName,
			ContractNameBak: txInfo.ContractName,
			ContractAddr:    txInfo.ContractAddr,
			Sender:          txInfo.Sender,
			Timestamp:       txInfo.Timestamp,
			UserAddr:        txInfo.UserAddr,
			UserAddrBns:     senderBns,
			GasUsed:         txInfo.GasUsed,
		}
		txViews.Add(latestBlockListView)
	}
	ConvergeListResponse(ctx, txViews.Values(), int64(len(txList)), nil)
}

// GetContractTxListHandler get
type GetContractTxListHandler struct {
}

// Handle deal
func (GetContractTxListHandler *GetContractTxListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetContractTxListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetContractTxList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var (
		err          error
		contractInfo *db.Contract
		txCount      int64
	)

	chainId := params.ChainId
	//获取合约详情
	if params.ContractAddr != "" {
		contractInfo, err = dbhandle.GetContractByCacheOrAddr(chainId, params.ContractAddr)
	} else {
		contractInfo, err = dbhandle.GetContractByCacheOrName(chainId, params.ContractName)
	}
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//合约没有交易
	if contractInfo == nil || contractInfo.TxNum == 0 {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	var userAddrs []string
	if params.UserAddrs != "" {
		userAddrs = strings.Split(params.UserAddrs, ",")
	}

	//有搜索条件，需要计算交易数量1
	if len(userAddrs) > 0 || params.ContractMethod != "" || *params.TxStatus >= 0 {
		txCount, err = dbhandle.GetContractTxCount(chainId, contractInfo.Addr, params.ContractMethod,
			userAddrs, *params.TxStatus)
	} else {
		txCount = contractInfo.TxNum
	}

	if err != nil {
		log.Errorf("GetContractTxCount err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//获取交易列表
	contractTxIDList, err := dbhandle.GetContractTxIDList(params.Offset, params.Limit, chainId, contractInfo.Addr,
		params.ContractMethod, userAddrs, *params.TxStatus)
	if err != nil {
		log.Errorf("GetContractTxIDList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	if len(contractTxIDList) == 0 {
		ConvergeListResponse(ctx, []interface{}{}, txCount, nil)
		return
	}

	contractTxList, err := dbhandle.GetContractTransactionList(chainId, contractTxIDList)
	if err != nil {
		log.Errorf("GetContractTransactionList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	txsViewInterface := BuildContractTxView(chainId, contractTxList)
	ConvergeListResponse(ctx, txsViewInterface, txCount, nil)
}

// BuildContractTxView
//
//	@Description: 构造返回值
//	@param chainId
//	@param contractTxList
//	@return []interface{}
func BuildContractTxView(chainId string, contractTxList []*db.ContractTxListResult) []interface{} {
	//获取账户BNS地址
	var accountAddrs []string
	for _, txInfo := range contractTxList {
		accountAddrs = append(accountAddrs, txInfo.UserAddr)
	}
	accountMap, _ := dbhandle.QueryAccountExists(chainId, accountAddrs)

	txListViews := make(ContractTxListViewSlice, 0, len(contractTxList))
	for _, txInfo := range contractTxList {
		var senderBns string
		txStatus := 0
		if txInfo.TxStatusCode != SUCCESS {
			txStatus = 1
		}

		if account, ok := accountMap[txInfo.UserAddr]; ok {
			senderBns = account.BNS
		}

		showStatus := ShowStatusTrue
		if txInfo.ContractMessageBak != "" || txInfo.ReadSetBak != "" {
			showStatus = ShowStatusFalse
		}

		txView := &entity.ContractTxListView{
			TxId:           txInfo.TxId,
			BlockHeight:    txInfo.BlockHeight,
			TxStatus:       txStatus,
			ShowStatus:     showStatus,
			ContractName:   txInfo.ContractName,
			ContractAddr:   txInfo.ContractAddr,
			ContractMethod: txInfo.ContractMethod,
			Sender:         txInfo.Sender,
			SenderOrgId:    txInfo.SenderOrgId,
			Timestamp:      txInfo.Timestamp,
			UserAddr:       txInfo.UserAddr,
			UserAddrBns:    senderBns,
		}
		txListViews = append(txListViews, txView)
	}

	// 对交易列表进行排序
	sort.Sort(txListViews)

	// 将排序后的交易列表转换为[]interface{}类型
	txsViewInterface := make([]interface{}, len(txListViews))
	for i, tx := range txListViews {
		txsViewInterface[i] = tx
	}

	return txsViewInterface
}

// GetBlockTxListHandler get
type GetBlockTxListHandler struct {
}

// Handle deal
func (GetBlockTxListHandler *GetBlockTxListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetBlockTxListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetBlockTxList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取区块详情
	block, err := dbhandle.GetBlockByHash(params.BlockHash, params.ChainId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
	}

	//区块没有交易
	if block == nil || block.TxCount == 0 {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	//获取交易列表
	blockTxIDList, err := dbhandle.GetBlockTxIDList(params.ChainId, params.BlockHash, params.Offset, params.Limit)
	if err != nil {
		log.Errorf("GetTxList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	if len(blockTxIDList) == 0 {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	blockTxList, _ := dbhandle.GetBlockTransactionList(params.ChainId, blockTxIDList)

	//获取账户BNS地址
	var accountAddrs []string
	for _, txInfo := range blockTxList {
		accountAddrs = append(accountAddrs, txInfo.UserAddr)
	}
	accountMap, _ := dbhandle.QueryAccountExists(params.ChainId, accountAddrs)

	txListViews := arraylist.New()
	for _, txInfo := range blockTxList {
		var senderBns string
		txStatus := 0
		if txInfo.TxStatusCode != SUCCESS {
			txStatus = 1
		}

		if account, ok := accountMap[txInfo.UserAddr]; ok {
			senderBns = account.BNS
		}

		txView := &entity.BlockTxListView{
			TxId:         txInfo.TxId,
			BlockHeight:  block.BlockHeight,
			BlockHash:    block.BlockHash,
			TxStatus:     txStatus,
			ContractName: txInfo.ContractName,
			ContractAddr: txInfo.ContractAddr,
			Sender:       txInfo.Sender,
			SenderOrgId:  txInfo.SenderOrgId,
			Timestamp:    txInfo.Timestamp,
			UserAddr:     txInfo.UserAddr,
			UserAddrBns:  senderBns,
		}
		txListViews.Add(txView)
	}

	ConvergeListResponse(ctx, txListViews.Values(), int64(block.TxCount), nil)
}

// GetUserTxListHandler get
type GetUserTxListHandler struct {
}

// Handle deal
func (GetUserTxListHandler *GetUserTxListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetUserTxListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetUserTxList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var userAddrs []string
	if params.UserAddrs != "" {
		userAddrs = strings.Split(params.UserAddrs, ",")
	}

	//获取账户详情
	accountList, err := dbhandle.QueryAccountExists(params.ChainId, userAddrs)
	if err != nil {
		log.Errorf("GetUserTxList get account info err:%v", err)
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	var totalCount int64
	for _, account := range accountList {
		totalCount += account.TxNum
	}

	//账户没有交易
	if totalCount == 0 {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	//获取交易列表
	userTxIDList, err := dbhandle.GetUserTxIDList(params.ChainId, userAddrs, params.Offset, params.Limit)
	if err != nil {
		log.Errorf("GetUserTxIDList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	if len(userTxIDList) == 0 {
		ConvergeListResponse(ctx, []interface{}{}, totalCount, nil)
		return
	}

	userTxList, err := dbhandle.GetContractTransactionList(params.ChainId, userTxIDList)
	if err != nil {
		log.Errorf("GetUserTxList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	txsViewInterface := BuildContractTxView(params.ChainId, userTxList)
	ConvergeListResponse(ctx, txsViewInterface, totalCount, nil)
}

// GetTxListHandler get
type GetTxListHandler struct {
}

// Handle deal
func (getTxListHandler *GetTxListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetTxListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetTxList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var (
		senders   []string
		userAddrs []string
	)

	if params.Senders != "" {
		senders = strings.Split(params.Senders, ",")
	}
	if params.UserAddrs != "" {
		userAddrs = strings.Split(params.UserAddrs, ",")
	}

	// If BlockHash is not empty, use GetBlockTxListHandler
	// todo 区块中的交易列表，等前端换接口GetBlockTxListHandler
	if shouldUseBlockTxListHandler(params, senders) {
		blockTxListHandler := GetBlockTxListHandler{}
		blockTxListHandler.Handle(ctx)
		return
	}
	// todo 合约中的交易列表，等前端换接口GetContractTxListHandler
	if shouldUseContractTxListHandler(params, senders) {
		contractTxListHandler := GetContractTxListHandler{}
		contractTxListHandler.Handle(ctx)
		return
	}
	// todo 账户中的交易列表，等前端换接口GetUserTxListHandler
	if shouldUseUserTxListHandler(params, userAddrs) {
		userTxListHandler := GetUserTxListHandler{}
		userTxListHandler.Handle(ctx)
		return
	}

	//计算交易总数
	totalCount, err := getTxTotalCount(params, userAddrs, senders)
	if err != nil {
		log.Errorf("GetTxList totalCount err : %s", err)
		ConvergeHandleFailureResponse(ctx, err)
	}

	//账户没有交易
	if totalCount == 0 {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	//获取交易id列表
	txIdList, err := dbhandle.GetTransactionIDList(params.ChainId, params.ContractName, params.BlockHash, params.Offset,
		params.Limit, params.StartTime, params.EndTime, params.TxId, *params.TxStatus, senders, userAddrs)
	if err != nil {
		log.Errorf("GetTxList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//获取交易列表
	txList, _ := dbhandle.BatchQueryTxList(params.ChainId, txIdList)
	txsViewInterface := BuildTxView(params.ChainId, txList)
	ConvergeListResponse(ctx, txsViewInterface, totalCount, nil)
}

func shouldUseBlockTxListHandler(params *entity.GetTxListParams, senders []string) bool {
	return params.BlockHash != "" && params.StartTime == 0 && len(senders) == 0 && *params.TxStatus == -1
}

func shouldUseContractTxListHandler(params *entity.GetTxListParams, senders []string) bool {
	return params.ContractName != "" && params.StartTime == 0 && params.TxId == "" &&
		len(senders) == 0 && *params.TxStatus == -1
}

func shouldUseUserTxListHandler(params *entity.GetTxListParams, userAddrs []string) bool {
	return len(userAddrs) > 0 && params.StartTime == 0 && params.TxId == "" &&
		params.ContractName == "" && *params.TxStatus == -1
}

// getTxTotalCount
//
//	@Description: 获取交易总数
//	@param params
//	@param userAddrs
//	@param senders
//	@return int64
//	@return error
func getTxTotalCount(params *entity.GetTxListParams, userAddrs, senders []string) (int64, error) {
	var (
		err        error
		totalCount int64
	)

	if params.StartTime == 0 && params.TxId == "" && len(userAddrs) == 0 &&
		len(senders) == 0 && *params.TxStatus == -1 {
		totalCount, err = dbhandle.GetTotalTxNum(params.ChainId)
		if err != nil {
			log.Errorf("GetTxList GetTotalTxNum err:%v", err)
		}
	} else if params.StartTime == 0 && params.TxId == "" && len(userAddrs) > 0 &&
		*params.TxStatus == -1 {
		//查询某账户的交易列表
		//获取账户详情
		accountList, err1 := dbhandle.QueryAccountExists(params.ChainId, userAddrs)
		if err1 != nil {
			log.Errorf("GetUserTxList get account info err:%v", err1)
			return totalCount, err1
		}
		for _, account := range accountList {
			totalCount += account.TxNum
		}
	} else {
		//获取交易列表交易总数contractName, blockHash
		totalCount, err = dbhandle.GetTransactionListCount(params.ChainId, params.TxId, params.ContractName,
			params.BlockHash, params.StartTime, params.EndTime, *params.TxStatus, senders, userAddrs)
		if err != nil {
			log.Errorf("GetTxList totalCount err : %s", err)
			return totalCount, err
		}
	}

	return totalCount, nil
}

// BuildTxView
//
//	@Description: 构造交易列表的返回值
//	@param chainId
//	@param txList 交易数据集
//	@return []interface{}
func BuildTxView(chainId string, txList []*db.Transaction) []interface{} {
	//获取账户BNS地址
	var accountAddrs []string
	for _, txInfo := range txList {
		accountAddrs = append(accountAddrs, txInfo.UserAddr)
	}
	accountMap, _ := dbhandle.QueryAccountExists(chainId, accountAddrs)

	txsView := make(TxListViewSlice, 0, len(txList))
	for _, tx := range txList {
		var addrBns string
		txStatus := 0
		showStatus := ShowStatusTrue
		if tx.TxStatusCode != SUCCESS {
			txStatus = 1
		}

		if account, ok := accountMap[tx.UserAddr]; ok {
			addrBns = account.BNS
		}

		if bytes.Equal(tx.ContractResultBak, config.ContractResultMsg) ||
			tx.ContractMessageBak != "" ||
			tx.ContractParametersBak != "" ||
			tx.ReadSetBak != "" ||
			tx.WriteSetBak != "" {
			showStatus = ShowStatusFalse
		}

		txListView := &entity.TxListView{
			Id:                 tx.TxId,
			BlockHeight:        tx.BlockHeight,
			BlockHash:          tx.BlockHash,
			TxId:               tx.TxId,
			Sender:             tx.Sender,
			SenderOrgId:        tx.SenderOrgId,
			UserAddr:           tx.UserAddr,
			UserAddrBns:        addrBns,
			ContractName:       tx.ContractName,
			ContractAddr:       tx.ContractAddr,
			ContractMethod:     tx.ContractMethod,
			ContractParameters: tx.ContractParameters,
			TxStatus:           txStatus,
			ShowStatus:         showStatus,
			Timestamp:          tx.Timestamp,
			GasUsed:            tx.GasUsed,
			PayerAddr:          tx.PayerAddr,
		}
		txsView = append(txsView, txListView)
	}

	// 对交易列表进行排序
	sort.Sort(txsView)

	// 将排序后的交易列表转换为[]interface{}类型
	txsViewInterface := make([]interface{}, len(txsView))
	for i, tx := range txsView {
		txsViewInterface[i] = tx
	}

	return txsViewInterface
}

// GetTxNumByTimeHandler get
type GetTxNumByTimeHandler struct {
}

// Handle 按时间段查询交易量
func (getTxNumByTimeHandler *GetTxNumByTimeHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetTransactionNumByTimeHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetTxNumByTime param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//默认24小时
	var timeObj = time.Now()
	var err error
	interval := int64(3600)
	startTime := params.StartTime
	endTime := params.EndTime
	if startTime == 0 || endTime == 0 {
		startTime = timeObj.Add(-(24 * time.Hour)).Unix()
		endTime = timeObj.Unix()
	}
	if params.Interval != 0 {
		interval = params.Interval
	}

	//获取缓存
	// Convert Unix timestamp to date + hour format
	startKey := time.Unix(startTime, 0).Format("2006-01-02 15")
	endKey := time.Unix(endTime, 0).Format("2006-01-02 15")
	txMap := GetTxNumByTimeCache(ctx, params.ChainId, startKey, endKey, interval)
	if txMap == nil {
		//按时间段获取交易数据
		txMap, err = dbhandle.GetTxListNumByRange(params.ChainId, startTime, endTime, interval)
		if err != nil {
			log.Errorf("GetTxNumByTime err : %s", err.Error())
			ConvergeHandleFailureResponse(ctx, err)
			return
		}
		SetTxNumByTimeCache(ctx, params.ChainId, startKey, endKey, interval, txMap)
	}
	if len(txMap) == 0 {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
	}

	transactionViews := arraylist.New()
	for t := endTime; t > startTime; t -= interval {
		// 将时间戳转换为整数
		txKey := t / interval * interval
		if _, ok := txMap[txKey]; !ok {
			continue
		}

		decimalView := &entity.TransactionNumView{
			TxNum:     txMap[txKey],
			Timestamp: txKey,
		}

		//正序排序
		if params.SortType == SortTypeAsc {
			// 将视图添加到 transactionViews 列表的开头
			transactionViews.Insert(0, decimalView)
		} else {
			//倒序排序
			transactionViews.Add(decimalView)
		}
	}

	ConvergeListResponse(ctx, transactionViews.Values(), int64(transactionViews.Size()), nil)
}

// SetTxNumByTimeCache 缓存首页24小时交易缓存数据
func SetTxNumByTimeCache(ctx *gin.Context, chainId, startKey, endKey string, interval int64, txMap map[int64]int64) {
	if len(txMap) == 0 {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisOverviewTxNumTime, prefix, chainId, startKey, endKey, interval)
	retJson, err := json.Marshal(txMap)
	if err != nil {
		log.Errorf("【Redis】set cache failed, key:%v, result:%v", redisKey, retJson)
		return
	}
	// 设置键值对和过期时间
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), 20*time.Minute).Err()
}

// GetTxNumByTimeCache 获取首页24小时交易缓存数据
func GetTxNumByTimeCache(ctx *gin.Context, chainId, startKey, endKey string, interval int64) map[int64]int64 {
	txMap := make(map[int64]int64)
	//从缓存获取最新的block
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisOverviewTxNumTime, prefix, chainId, startKey, endKey, interval)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes == nil || redisRes.Val() == "" {
		return nil
	}

	err := json.Unmarshal([]byte(redisRes.Val()), &txMap)
	if err != nil {
		log.Errorf("【Redis】get cache failed, key:%v, result:%v", redisKey, redisRes)
		return nil
	}
	return txMap
}

// GetContractVersionListHandler get
type GetContractVersionListHandler struct {
}

// Handle 获取合约创建，更新，冻结，等合约类交易列表
func (getContractVersionListHandler *GetContractVersionListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetContractVersionListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetContractTxList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	txViews := arraylist.New()
	var senders []string
	if params.Senders != "" {
		senders = strings.Split(params.Senders, ",")
	}

	txList, count, err := dbhandle.GetUpgradeContractTxList(params.Offset, params.Limit, params.ChainId,
		params.ContractName, params.ContractAddr, senders, params.RuntimeType, *params.Status,
		params.StartTime, params.EndTime)
	if err != nil {
		log.Errorf("GetContractVersionList err : %s", err.Error())
		ConvergeListResponse(ctx, txViews.Values(), count, nil)
		return
	}

	var accountAddrs []string
	for _, txInfo := range txList {
		accountAddrs = append(accountAddrs, txInfo.UserAddr)
	}
	accountMap, _ := dbhandle.QueryAccountExists(params.ChainId, accountAddrs)

	thirdApplyUrl := config.GlobalConfig.WebConf.ThirdApplyUrl
	for _, tx := range txList {
		contractResultCode := 0
		if tx.ContractResultCode != 0 {
			contractResultCode = 1
		}
		txUrl := utils.SplicePath(thirdApplyUrl, params.ChainId, "/transaction/", tx.TxId)
		senderAddrBNS := GetAccountBNS(tx.UserAddr, accountMap)
		contractVersionView := &entity.ContractVersionView{
			TxId:               tx.TxId,
			ContractName:       tx.ContractName,
			ContractAddr:       tx.ContractAddr,
			Version:            tx.ContractVersion,
			SenderOrgId:        tx.SenderOrgId,
			Sender:             tx.Sender,
			SenderAddr:         tx.UserAddr,
			SenderAddrBNS:      senderAddrBNS,
			TxUrl:              txUrl,
			Timestamp:          tx.Timestamp,
			RuntimeType:        tx.ContractRuntimeType,
			ContractResultCode: contractResultCode,
		}
		txViews.Add(contractVersionView)
	}
	ConvergeListResponse(ctx, txViews.Values(), count, nil)
}
