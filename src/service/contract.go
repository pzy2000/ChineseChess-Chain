/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"chainmaker.org/chainmaker/contract-utils/standard"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"chainmaker_web/src/config"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
)

// GetLatestContractListHandler get
type GetLatestContractListHandler struct{}

// Handle deal
func (getLatestContractListHandler *GetLatestContractListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetLatestContractHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetLatestContractList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//从缓存获取最新的Contract
	contractList, err := getContractListFromRedis(ctx, params.ChainId)
	if err != nil {
		log.Errorf("GetLatestContractList get redis fail err:%v", err)
	}
	count := int64(len(contractList))
	if count == 0 {
		// 获取ContractList
		contractList, err = dbhandle.GetLatestContractList(params.ChainId)
		if err != nil {
			log.Errorf("GetLatestContractList err : %s", err.Error())
			ConvergeHandleFailureResponse(ctx, err)
			return
		}
	}

	//获取创建合约账户对应的账户信息
	accountMap := GetContractAccountMap(params.ChainId, contractList)

	//数据渲染
	contractViews := arraylist.New()
	for i, contract := range contractList {
		//获取地址BNS
		senderAddrBns := GetAccountBNS(contract.CreatorAddr, accountMap)
		latestChainView := &entity.LatestContractView{
			Id:               i + 1,
			ContractName:     contract.Name,
			ContractAddr:     contract.Addr,
			ContractType:     contract.ContractType,
			Sender:           contract.CreateSender,
			SenderAddr:       contract.CreatorAddr,
			SenderAddrBNS:    senderAddrBns,
			Version:          contract.Version,
			TxNum:            contract.TxNum,
			CreateTimestamp:  contract.Timestamp,
			UpgradeTimestamp: contract.UpgradeTimestamp,
			UpgradeUser:      contract.UpgradeAddr,
			Timestamp:        contract.Timestamp,
		}
		contractViews.Add(latestChainView)
	}
	ConvergeListResponse(ctx, contractViews.Values(), count, nil)
}

// getContractListFromRedis 获取缓存数据
func getContractListFromRedis(ctx *gin.Context, chainId string) ([]*db.Contract, error) {
	contractList := make([]*db.Contract, 0)
	//从缓存获取最新的合约
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisLatestContractList, prefix, chainId)
	redisList := cache.GlobalRedisDb.ZRevRange(ctx, redisKey, 0, 9).Val()
	for _, resStr := range redisList {
		contractInfo := &db.Contract{}
		err := json.Unmarshal([]byte(resStr), contractInfo)
		if err != nil {
			log.Errorf("getContractListFromRedis json Unmarshal err : %s", err.Error())
			return contractList, err
		}
		contractList = append(contractList, contractInfo)
	}

	return contractList, nil
}

// GetContractListHandler handler
type GetContractListHandler struct{}

// Handle deal
func (getContractListHandler *GetContractListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetContractListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetContractList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var (
		senders      []string
		senderAddrs  []string
		upgraders    []string
		upgradeAddrs []string
	)

	chainId := params.ChainId
	offset := params.Offset
	limit := params.Limit
	if params.Creators != "" {
		senders = strings.Split(params.Creators, ",")
	}
	if params.CreatorAddrs != "" {
		senderAddrs = strings.Split(params.CreatorAddrs, ",")
	}
	if params.Upgraders != "" {
		upgraders = strings.Split(params.Upgraders, ",")
	}
	if params.UpgradeAddrs != "" {
		upgradeAddrs = strings.Split(params.UpgradeAddrs, ",")
	}

	contractList, count, err := dbhandle.GetContractList(chainId, offset, limit, params.Status, params.RuntimeType,
		params.ContractKey, senders, senderAddrs, upgraders, upgradeAddrs, params.StartTime, params.EndTime)
	if err != nil {
		log.Errorf("GetContractList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//获取创建合约账户对应的账户信息
	accountMap := GetContractAccountMap(chainId, contractList)
	contractListView := arraylist.New()
	for i, contract := range contractList {
		//获取地址BNS
		senderAddrBns := GetAccountBNS(contract.CreatorAddr, accountMap)

		listId := params.Offset*params.Limit + i + 1
		contractView := &entity.ContractListView{
			Id:               strconv.Itoa(listId),
			ContractName:     contract.Name,
			ContractSymbol:   contract.ContractSymbol,
			ContractAddr:     contract.Addr,
			ContractType:     contract.ContractType,
			Version:          contract.Version,
			Creator:          contract.CreateSender,
			CreatorAddr:      contract.CreatorAddr,
			CreatorAddrBns:   senderAddrBns,
			Upgrader:         contract.Upgrader,
			UpgradeAddr:      contract.UpgradeAddr,
			UpgradeOrgId:     contract.UpgradeOrgId,
			TxNum:            contract.TxNum,
			Status:           contract.ContractStatus,
			CreateTimestamp:  contract.Timestamp,
			UpgradeTimestamp: contract.UpgradeTimestamp,
			RuntimeType:      contract.RuntimeType,
		}
		contractListView.Add(contractView)
	}

	ConvergeListResponse(ctx, contractListView.Values(), count, nil)
}

// GetContractDetailHandler handler
type GetContractDetailHandler struct{}

// Handle deal
func (getContractDetailHandler *GetContractDetailHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetContractDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取合约
	contract, err := dbhandle.GetContractByNameOrAddr(params.ChainId, params.ContractKey)
	if err != nil {
		log.Errorf("GetContractDetail err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//获取创建合约账户对应的账户信息
	accountMap := GetContractAccountMap(params.ChainId, []*db.Contract{contract})
	//获取地址BNS
	senderAddrBns := GetAccountBNS(contract.CreatorAddr, accountMap)

	var dataAssetNum int64
	//判断是否是IDA合约
	if contract.ContractType == standard.ContractStandardNameCMIDA {
		idaContractInfo, err := dbhandle.GetIDAContractByAddr(params.ChainId, contract.Addr)
		if err != nil {
			if err != nil {
				log.Errorf("GetIDAContractByAddr err : %v", err)
				ConvergeHandleFailureResponse(ctx, err)
				return
			}
		}
		if idaContractInfo != nil {
			dataAssetNum = idaContractInfo.TotalNormalAssets
		}
	}
	contractDetailView := &entity.ContractDetailView{
		ContractName:    contract.Name,
		ContractNameBak: contract.NameBak,
		ContractAddr:    contract.Addr,
		ContractSymbol:  contract.ContractSymbol,
		ContractType:    contract.ContractType,
		Version:         contract.Version,
		ContractStatus:  contract.ContractStatus,
		TxId:            contract.CreateTxId,
		CreateSender:    contract.CreateSender,
		CreatorAddr:     contract.CreatorAddr,
		CreatorAddrBns:  senderAddrBns,
		Timestamp:       contract.Timestamp,
		DataAssetNum:    dataAssetNum,
		RuntimeType:     contract.RuntimeType,
	}

	ConvergeDataResponse(ctx, contractDetailView, nil)
}

// GetContractCodeHandler handler
type GetContractCodeHandler struct {
}

// Handle deal
func (getContractCodeHandler *GetContractCodeHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetContractCodeHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetContractCode param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// TODO：需要对应的服务进行修改，增加ChainId
	url := fmt.Sprintf(config.GlobalConfig.WebConf.TestnetUrl+"/contract/info?name=%s", params.ContractName)
	client := http.Client{Timeout: time.Second}
	// nolint
	resp, err := client.Get(url)
	if err != nil {
		log.Errorf("Get ContractCode from remote err, cause : %s", err.Error())
	}
	if resp != nil {
		err = resp.Body.Close()
		if err != nil {
			log.Errorf("close resp error, err:%s", err.Error())
		}
	}

	url = fmt.Sprintf(config.GlobalConfig.WebConf.OpennetUrl+"/contract/info?contractName=%s", params.ContractName)
	// nolint
	resp, err = client.Get(url)
	contractCodeView := &entity.ContractCodeView{}
	if err != nil {
		log.Errorf("Get ContractCode from remote err, cause : %s", err.Error())
		ConvergeDataResponse(ctx, contractCodeView, nil)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Errorf("Get ContractCode from remote err, StatusCode : %d", resp.StatusCode)
		ConvergeDataResponse(ctx, contractCodeView, nil)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Read resp body err, cause : %s", err.Error())
		ConvergeDataResponse(ctx, contractCodeView, nil)
		return
	}
	var respJson interface{}
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Errorf("Json Unmarshal from resp body err, cause : %s", err.Error())
		ConvergeDataResponse(ctx, contractCodeView, nil)
		return
	}
	respMap, ok := respJson.(map[string]interface{})

	if !ok || respMap["code"].(float64) != 0 {
		log.Error("contranct name don't exist")
		ConvergeDataResponse(ctx, contractCodeView, nil)
		return
	}

	// nolint
	dataMap := respMap["data"].(map[string]interface{})

	contractCodeView.ContractCode, ok = dataMap["source_code"].(string)
	if !ok {
		log.Warn("source_code don't exist")
	}
	contractCodeView.ContractAbi, ok = dataMap["abi_json"].(string)
	if !ok {
		log.Warn("abi_json don't exist")
	}
	contractCodeView.ContractByteCode, ok = dataMap["byte_code"].(string)
	if !ok {
		log.Warn("byte_code don't exist")
	}

	ConvergeDataResponse(ctx, contractCodeView, nil)
}

// GetEventListHandler handler
type GetEventListHandler struct {
}

// Handle deal
func (getEventListHandler *GetEventListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetEventListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetEventList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var contractInfo *db.Contract
	var err error
	//获取合约详情
	if params.ContractAddr != "" {
		contractInfo, err = dbhandle.GetContractByCacheOrAddr(params.ChainId, params.ContractAddr)
	} else {
		contractInfo, err = dbhandle.GetContractByCacheOrName(params.ChainId, params.ContractName)
	}
	if err != nil || contractInfo == nil {
		ConvergeHandleFailureResponse(ctx, err)
	}

	////获取交易列表交易总数
	//totalCount, err := dbhandle.GetEventListCount(params.ChainId, params.ContractName, params.ContractAddr, params.TxId)
	//if err != nil {
	//	log.Errorf("GetEventList totalCount err : %s", err)
	//	ConvergeHandleFailureResponse(ctx, err)
	//}

	if contractInfo.EventNum == 0 {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	eventIds, err := dbhandle.GetEventIdList(params.Offset, params.Limit, params.ChainId, params.ContractName,
		params.ContractAddr, params.TxId)
	if err != nil {
		log.Errorf("GetEventList err : %s", err.Error())
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	eventList, err := dbhandle.GetEventDataByIds(params.ChainId, eventIds)
	if err != nil {
		log.Errorf("GetTxList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	eventListView := make(ContractEventViewSlice, 0, len(eventList))
	for _, event := range eventList {
		contractEventView := &entity.ContractEventView{
			Topic:     event.Topic,
			EventInfo: event.EventData,
			Timestamp: event.Timestamp,
		}
		eventListView = append(eventListView, contractEventView)
	}

	// 对交易列表进行排序
	sort.Sort(eventListView)

	// 将排序后的交易列表转换为[]interface{}类型
	eventsViewInterface := make([]interface{}, len(eventListView))
	for i, tx := range eventListView {
		eventsViewInterface[i] = tx
	}
	ConvergeListResponse(ctx, eventsViewInterface, contractInfo.EventNum, nil)
}

type ContractEventViewSlice []*entity.ContractEventView

func (t ContractEventViewSlice) Len() int {
	return len(t)
}

func (t ContractEventViewSlice) Less(i, j int) bool {
	return t[i].Timestamp > t[j].Timestamp
}

func (t ContractEventViewSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
