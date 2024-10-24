package service

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/chain"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"chainmaker_web/src/entity_cross"
	"chainmaker_web/src/sync"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
)

const (
	//SearchCrossID searchType
	SearchCrossID = iota
	// SearchSubChainID searchType
	SearchSubChainID
	//SearchCrossUnKnow searchType
	SearchCrossUnKnow = -1
)

// GetMainCrossConfigHandler get
type GetMainCrossConfigHandler struct {
}

// Handle GetMainCrossConfigHandler 主子链网配置,是否是主链
func (handler *GetMainCrossConfigHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetChainIdHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	mainCrossConfigView := &entity_cross.MainCrossConfig{
		ShowTag: chain.GetIsMainChain(),
	}
	//返回response
	ConvergeDataResponse(ctx, mainCrossConfigView, nil)
}

// CrossSearchHandler sub
type CrossSearchHandler struct {
}

// Handle SubscribeChainHandler
func (handler *CrossSearchHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetCrossSearchHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var typeView int
	var valueView string
	_, err := strconv.ParseInt(params.Value, 10, 64)
	if err == nil {
		//数字，跨链ID
		crossTxInfo, err := dbhandle.GetCrossCycleById(params.ChainId, params.Value)
		if err == nil && crossTxInfo != nil {
			valueView = crossTxInfo.CrossId
			typeView = SearchCrossID
		}
	} else {
		//字符串，子链名称
		subChainInfo, err := dbhandle.GetCrossSubChainInfoByName(params.ChainId, params.Value)
		if err == nil && subChainInfo != nil {
			valueView = subChainInfo.SubChainId
			typeView = SearchSubChainID
		}
	}

	if valueView == "" {
		typeView = SearchCrossUnKnow
	}
	crossSearchView := &entity_cross.CrossSearchView{
		Type: typeView,
		Data: valueView,
	}
	//返回response
	ConvergeDataResponse(ctx, crossSearchView, nil)
}

// CrossOverviewDataHandler cancel
type CrossOverviewDataHandler struct {
}

// Handle deal
func (handler *CrossOverviewDataHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetChainIdHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	chainId := params.ChainId
	//获取缓存数据
	overviewData := GetCrossOverviewDataCache(ctx, chainId)
	if overviewData != nil {
		ConvergeDataResponse(ctx, overviewData, nil)
		return
	}

	//总区块高度
	totalBlockHeight, err := dbhandle.GetAllSubChainBlockHeight(chainId)
	if err != nil {
		log.Errorf("GetAllSubChainBlockHeight err : %v", err)
	}

	//本自然月的开始，结束时间
	startTime, endTime := GetCurrentMonthStartAndEndTime()
	//周期交易最短完成时间
	shortestTime, err := dbhandle.GetCycleShortestTime(chainId, startTime, endTime)
	if err != nil {
		log.Errorf("GetCycleShortestTime err : %v", err)
	}

	//周期交易最长完成时间
	longestTime, err := dbhandle.GetCycleLongestTime(chainId, startTime, endTime)
	if err != nil {
		log.Errorf("GetCycleShortestTime err : %v", err)
	}

	//周期交易平均完成时间
	averageTime, err := dbhandle.GetCycleAverageTime(chainId, startTime, endTime)
	if err != nil {
		log.Errorf("GetCycleShortestTime err : %v", err)
	}

	//子链总数
	subChainNum, err := dbhandle.GetCrossSubChainAllCount(chainId)
	if err != nil {
		log.Errorf("GetCrossSubChainAllCount err : %v", err)
	}

	//跨链交易总数
	cycleTxNum, err := dbhandle.GetCrossCycleTxAllCount(chainId)
	if err != nil {
		log.Errorf("GetCrossSubChainAllCount err : %v", err)
	}
	overviewData = &entity_cross.OverviewDataView{
		TotalBlockHeight: totalBlockHeight,
		ShortestTime:     shortestTime,
		LongestTime:      longestTime,
		AverageTime:      averageTime,
		SubChainNum:      subChainNum,
		TxNum:            cycleTxNum,
	}

	//设置缓存
	SetCrossOverviewDataCache(ctx, chainId, *overviewData)
	//返回response
	ConvergeDataResponse(ctx, overviewData, nil)
}

// GetCrossOverviewDataCache 获取首页缓存数据
func GetCrossOverviewDataCache(ctx *gin.Context, chainId string) *entity_cross.OverviewDataView {
	cacheResult := &entity_cross.OverviewDataView{}
	//从缓存获取最新的block
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossOverviewData, prefix, chainId)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes == nil || redisRes.Val() == "" {
		return nil
	}

	err := json.Unmarshal([]byte(redisRes.Val()), &cacheResult)
	if err != nil {
		log.Errorf("【Redis】get cache failed, key:%v, result:%v", redisKey, redisRes)
		return nil
	}
	return cacheResult
}

// SetCrossOverviewDataCache 缓存首页信息
func SetCrossOverviewDataCache(ctx *gin.Context, chainId string, overviewData entity_cross.OverviewDataView) {
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossOverviewData, prefix, chainId)
	retJson, err := json.Marshal(overviewData)
	if err != nil {
		log.Errorf("【Redis】set cache failed, key:%v, result:%v", redisKey, retJson)
		return
	}
	// 设置键值对和过期时间(40s 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), 40*time.Second).Err()
}

// CrossLatestTxListHandler modify
type CrossLatestTxListHandler struct {
}

// Handle CrossLatestTxListHandler 最新跨链交易
func (handler *CrossLatestTxListHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetChainIdHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// txList
	crossTxList, err := dbhandle.GetCrossLatestCycleTxList(params.ChainId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	txViews := arraylist.New()
	if len(crossTxList) == 0 {
		ConvergeListResponse(ctx, txViews.Values(), 0, nil)
	}

	for _, cycleTx := range crossTxList {
		var fromChainName string
		var toChainName string
		fromChainId := cycleTx.CrossTransactionTransfer.FromChainId
		fromIsMainChain := cycleTx.CrossTransactionTransfer.FromIsMainChain
		toChainId := cycleTx.CrossTransactionTransfer.ToChainId
		toIsMainChain := cycleTx.CrossTransactionTransfer.ToIsMainChain

		if fromIsMainChain {
			fromChainName = config.GlobalConfig.ChainConf.MainChainName
		} else {
			fromChainName, _ = dbhandle.GetCrossSubChainName(params.ChainId, fromChainId)
		}

		if toIsMainChain {
			toChainName = config.GlobalConfig.ChainConf.MainChainName
		} else {
			toChainName, _ = dbhandle.GetCrossSubChainName(params.ChainId, toChainId)
		}

		var status int32
		if cycleTx.CrossCycleTransaction.Status == int32(tcipCommon.CrossChainStateValue_CONFIRM_END) {
			status = 1
		} else if cycleTx.CrossCycleTransaction.Status == int32(tcipCommon.CrossChainStateValue_CANCEL_END) {
			status = 2
		}
		latestListView := &entity_cross.LatestTxListView{
			CrossId:         cycleTx.CrossId,
			Status:          status, //跨链状态（0:进行中，1:成功，2:失败）
			Timestamp:       cycleTx.CrossCycleTransaction.StartTime,
			FromChainName:   fromChainName,
			FromChainId:     fromChainId,
			FromIsMainChain: fromIsMainChain,
			ToChainName:     toChainName,
			ToChainId:       toChainId,
			ToIsMainChain:   toIsMainChain,
		}
		txViews.Add(latestListView)
	}
	ConvergeListResponse(ctx, txViews.Values(), int64(len(crossTxList)), nil)
}

// CrossLatestSubChainListHandler delete
type CrossLatestSubChainListHandler struct {
}

// Handle CrossLatestSubChainListHandler
func (handler *CrossLatestSubChainListHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetChainIdHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// crossSubChainList
	crossSubChainList, err := dbhandle.GetCrossLatestSubChainList(params.ChainId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	sunChainListViews := arraylist.New()
	if len(crossSubChainList) == 0 {
		ConvergeListResponse(ctx, sunChainListViews.Values(), 0, nil)
	}

	for _, subChain := range crossSubChainList {
		//跨链合约数
		crossContractNum, err := dbhandle.GetCrossContractCount(params.ChainId, subChain.SubChainId)
		if err != nil {
			log.Errorf("Get CrossContract Count err : %v", err)
		}
		latestListView := &entity_cross.LatestSubChainListView{
			SubChainId:       subChain.SubChainId,
			SubChainName:     subChain.ChainName,
			BlockHeight:      subChain.BlockHeight,
			Timestamp:        subChain.Timestamp,
			Status:           subChain.Status,
			CrossTxNum:       subChain.TxNum,
			CrossContractNum: crossContractNum,
		}
		sunChainListViews.Add(latestListView)
	}

	ConvergeListResponse(ctx, sunChainListViews.Values(), int64(len(crossSubChainList)), nil)
}

// GetCrossTxListHandler 跨链交易列表
type GetCrossTxListHandler struct {
}

// Handle CrossLatestTxListHandler 最新跨链交易
func (handler *GetCrossTxListHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetCrossTxListHandler(ctx)
	if params == nil || !params.IsLegal() || !params.RangeBody.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	chainId := params.ChainId
	fromChainName := params.FromChainName
	toChainName := params.ToChainName
	subChainId := params.SubChainId
	subChainInfo, err := dbhandle.GetCrossSubChainInfoById(chainId, subChainId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	isNull, fromChainId, toChainId, err := getFromToChainIdByName(chainId, subChainId, fromChainName, toChainName)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	if isNull {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	var totalCount int64
	if params.StartTime == 0 && params.EndTime == 0 && params.CrossId == "" && fromChainId == "" && toChainId == "" {
		if subChainId == "" {
			//跨链交易总数
			totalCount, err = dbhandle.GetCrossCycleTxAllCount(params.ChainId)
		} else {
			if subChainInfo != nil {
				totalCount = subChainInfo.TxNum
			}
		}
	} else {
		totalCount, err = dbhandle.GetCrossSubChainTxCount(params.StartTime, params.EndTime, chainId, params.CrossId,
			subChainId, fromChainId, toChainId)
	}
	if err != nil {
		log.Errorf("GetCrossCycleTxCount err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	if totalCount == 0 {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	//获取交易列表
	crossTxList, err := dbhandle.GetCrossSubChainTxList(params.Offset, params.Limit, params.StartTime,
		params.EndTime, params.ChainId, params.CrossId, params.SubChainId, fromChainId, toChainId)
	if err != nil {
		log.Errorf("GetTxList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//构造返回值
	txListViews := buildCrossTxListView(params.ChainId, crossTxList)
	ConvergeListResponse(ctx, txListViews.Values(), totalCount, nil)
}

func getFromToChainIdByName(chainId, subChainId, fromChainName, toChainName string) (bool, string, string, error) {
	var (
		fromChainId string
		toChainId   string
		err         error
		isNull      bool
	)
	if fromChainName != "" {
		fromChainId, err = GetSubChainIdByName(chainId, fromChainName)
		if err != nil {
			return isNull, fromChainId, toChainId, err
		}
		if fromChainId == "" {
			isNull = true
		}
	}

	if toChainName != "" {
		toChainId, err = GetSubChainIdByName(chainId, toChainName)
		if err != nil {
			return isNull, fromChainId, toChainId, err
		}
		if toChainId == "" {
			isNull = true
		}
	}

	if subChainId != "" {
		if fromChainId != "" && fromChainId != subChainId && toChainId == "" {
			toChainId = subChainId
		}
		if toChainId != "" && toChainId != subChainId && fromChainId == "" {
			fromChainId = subChainId
		}
	}

	return isNull, fromChainId, toChainId, nil
}

// 构造交易列表返回值
func buildCrossTxListView(chainId string, crossTxList []*db.CycleJoinTransferResult) *arraylist.List {
	txListViews := arraylist.New()
	if len(crossTxList) == 0 {
		return txListViews
	}

	for _, cycleTx := range crossTxList {
		var fromChainName string
		var toChainName string
		fromChainId := cycleTx.CrossTransactionTransfer.FromChainId
		fromIsMainChain := cycleTx.CrossTransactionTransfer.FromIsMainChain
		toChainId := cycleTx.CrossTransactionTransfer.ToChainId
		toIsMainChain := cycleTx.CrossTransactionTransfer.ToIsMainChain

		if fromIsMainChain {
			fromChainName = config.GlobalConfig.ChainConf.MainChainName
		} else {
			fromChainName, _ = dbhandle.GetCrossSubChainName(chainId, fromChainId)
		}

		if toIsMainChain {
			toChainName = config.GlobalConfig.ChainConf.MainChainName
		} else {
			toChainName, _ = dbhandle.GetCrossSubChainName(chainId, toChainId)
		}

		var status int32
		if cycleTx.CrossCycleTransaction.Status == int32(tcipCommon.CrossChainStateValue_CONFIRM_END) {
			status = 1
		} else if cycleTx.CrossCycleTransaction.Status == int32(tcipCommon.CrossChainStateValue_CANCEL_END) {
			status = 2
		}
		txView := &entity_cross.GetTxListView{
			CrossId:         cycleTx.CrossId,
			Status:          status, //跨链状态（0:进行中，1:成功，2:失败）
			Timestamp:       cycleTx.CrossCycleTransaction.StartTime,
			FromChainName:   fromChainName,
			FromChainId:     fromChainId,
			FromIsMainChain: fromIsMainChain,
			ToChainName:     toChainName,
			ToChainId:       toChainId,
			ToIsMainChain:   toIsMainChain,
		}
		txListViews.Add(txView)
	}
	return txListViews
}

// CrossSubChainListHandler get
type CrossSubChainListHandler struct {
}

// Handle CrossSubChainListHandler
func (handler *CrossSubChainListHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetCrossSubChainListHandler(ctx)
	if params == nil || !params.IsLegal() || !params.RangeBody.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	crossSubChainList, totalCount, err := dbhandle.GetCrossSubChainList(params.Offset, params.Limit, params.ChainId,
		params.SubChainId, params.SubChainName)
	if err != nil {
		log.Errorf("GetCrossSubChainList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	sunChainListViews := arraylist.New()
	if len(crossSubChainList) == 0 {
		ConvergeListResponse(ctx, sunChainListViews.Values(), 0, nil)
		return
	}

	for _, subChain := range crossSubChainList {
		//跨链合约数
		crossContractNum, err := dbhandle.GetCrossContractCount(params.ChainId, subChain.SubChainId)
		if err != nil {
			log.Errorf("Get CrossContract Count err : %v", err)
		}
		chainView := &entity_cross.GetSubChainListView{
			SubChainId:       subChain.SubChainId,
			SubChainName:     subChain.ChainName,
			BlockHeight:      subChain.BlockHeight,
			Timestamp:        subChain.Timestamp,
			Status:           subChain.Status,
			CrossTxNum:       subChain.TxNum,
			CrossContractNum: crossContractNum,
			ExplorerUrl:      subChain.ExplorerAddr,
		}
		sunChainListViews.Add(chainView)
	}

	ConvergeListResponse(ctx, sunChainListViews.Values(), totalCount, nil)
}

// CrossSubChainDetailHandler get
type CrossSubChainDetailHandler struct {
}

// Handle CrossSubChainDetailHandler 主子链网-子链详情
func (handler *CrossSubChainDetailHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetCrossSubChainDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//子链信息
	subChainInfo, err := dbhandle.GetCrossSubChainInfoById(params.ChainId, params.SubChainId)
	if err != nil {
		log.Errorf("GetCrossSubChainInfoById err : %v", err)
		ConvergeHandleFailureResponse(ctx, entity.ErrSelectFailed)
		return
	} else if subChainInfo == nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	//跨链合约数
	crossContractNum, err := dbhandle.GetCrossContractCount(params.ChainId, params.SubChainId)
	if err != nil {
		log.Errorf("Get CrossContract Count err : %v", err)
	}

	subChainView := &entity_cross.GetCrossSubChainDetailView{
		SubChainId:       subChainInfo.SubChainId,
		SubChainName:     subChainInfo.ChainName,
		BlockHeight:      subChainInfo.BlockHeight,
		ChainType:        subChainInfo.ChainType,
		CrossTxNum:       subChainInfo.TxNum,
		CrossContractNum: crossContractNum,
		Status:           subChainInfo.Status,
		GatewayId:        subChainInfo.GatewayId,
		GatewayName:      subChainInfo.GatewayName,
		GatewayAddr:      subChainInfo.GatewayAddr,
		Timestamp:        subChainInfo.Timestamp,
	}

	//返回response
	ConvergeDataResponse(ctx, subChainView, nil)
}

// GetCrossTxDetailHandler get
type GetCrossTxDetailHandler struct {
}

// Handle GetCrossTxDetailHandler 主子链网-跨链交易详情
func (handler *GetCrossTxDetailHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetCrossTxDetailHandler(ctx)
	newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
	if params == nil || !params.IsLegal() {
		ConvergeFailureResponse(ctx, newError)
		return
	}

	chainId := params.ChainId
	//跨链交易
	crossTxInfo, err := dbhandle.GetCrossCycleById(params.ChainId, params.CrossId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrSelectFailed)
		return
	} else if crossTxInfo == nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	//跨链交易流转
	crossTransfers, err := dbhandle.GetCrossCycleTransferById(params.ChainId, params.CrossId)
	if err != nil || len(crossTransfers) == 0 {
		ConvergeHandleFailureResponse(ctx, entity.ErrSelectFailed)
		return
	}

	transferInfo := crossTransfers[0]
	businessList, err := dbhandle.GetCrossBusinessTxByCross(params.ChainId, transferInfo.CrossId)
	var (
		fromGas  = "-"
		toGas    = "-"
		fromTxId string
		toTxId   string
	)
	fromTxInfo := &db.CrossBusinessTransaction{}
	toTxInfo := &db.CrossBusinessTransaction{}
	if err == nil && len(businessList) > 0 {
		for _, tx := range businessList {
			if transferInfo.FromChainId == tx.SubChainId {
				fromTxInfo = tx
				fromGas = strconv.FormatUint(tx.GasUsed, 10)
				fromTxId = fromTxInfo.TxId
			} else if transferInfo.ToChainId == tx.SubChainId {
				toTxInfo = tx
				toGas = strconv.FormatUint(tx.GasUsed, 10)
				toTxId = toTxInfo.TxId
			}
		}
	}

	fromIsMainChain := transferInfo.FromIsMainChain
	fromChainId := transferInfo.FromChainId
	fromChainName, fromTxUrl := getSubChainNameUrl(chainId, fromChainId, fromTxId, fromIsMainChain)
	crossDirection := &entity_cross.CrossDirection{
		FromChain: fromChainName,
	}
	fromChainTx := &entity_cross.TxChainInfo{
		ChainName:    fromChainName,
		ChainId:      fromChainId,
		ContractName: fromTxInfo.ContractName,
		IsMainChain:  fromIsMainChain,
		TxId:         fromTxId,
		TxStatus:     fromTxInfo.TxStatus,
		TxUrl:        fromTxUrl,
		Gas:          fromGas,
	}

	toIsMainChain := transferInfo.ToIsMainChain
	toChainId := transferInfo.ToChainId
	toChainName, toTxUrl := getSubChainNameUrl(chainId, toChainId, toTxId, toIsMainChain)
	crossDirection.ToChain = toChainName
	toChainTx := &entity_cross.TxChainInfo{
		ChainName:    toChainName,
		ChainId:      toChainId,
		ContractName: transferInfo.ContractName,
		IsMainChain:  toIsMainChain,
		TxId:         toTxId,
		TxStatus:     toTxInfo.TxStatus,
		TxUrl:        toTxUrl,
		Gas:          toGas,
	}

	txDetailView := &entity_cross.GetCrossTxDetailView{
		CrossId:        crossTxInfo.CrossId,
		Status:         crossTxInfo.Status,
		CrossDuration:  crossTxInfo.Duration,
		ContractName:   transferInfo.ContractName,
		ContractMethod: transferInfo.ContractMethod,
		Parameter:      transferInfo.Parameter,
		ContractResult: toTxInfo.CrossContractResult,
		CrossDirection: crossDirection,
		FromChainInfo:  fromChainTx,
		ToChainInfo:    toChainTx,
		Timestamp:      crossTxInfo.StartTime,
	}

	//返回response
	ConvergeDataResponse(ctx, txDetailView, nil)
}

func getSubChainNameUrl(chainId, subChainId, txId string, isMainChain bool) (string, string) {
	var (
		subChainName  string
		explorerTxUrl string
	)

	subChainName = config.GlobalConfig.ChainConf.MainChainName
	if !isMainChain {
		subChainInfo := dbhandle.GetCrossSubChainInfoCache(chainId, subChainId)
		if subChainInfo == nil {
			subChainInfo, _ = dbhandle.GetCrossSubChainInfoById(chainId, subChainId)
		}
		if subChainInfo != nil {
			subChainName = subChainInfo.ChainName
			if subChainInfo.ExplorerTxAddr != "" {
				explorerTxUrl = subChainInfo.ExplorerTxAddr + txId
			}
		}
	}

	return subChainName, explorerTxUrl
}

// SubChainCrossChainListHandler get
type SubChainCrossChainListHandler struct {
}

// Handle CrossSubChainListHandler
func (handler *SubChainCrossChainListHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetSubChainCrossChainListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	subChainCrossList, err := dbhandle.GetSubChainCrossChainList(params.ChainId, params.SubChainId)
	if err != nil {
		log.Errorf("GetSubChainCrossChainList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	sunChainCrossListViews := arraylist.New()
	if len(subChainCrossList) == 0 {
		ConvergeListResponse(ctx, sunChainCrossListViews.Values(), 0, nil)
		return
	}

	for _, crossChain := range subChainCrossList {
		crossChainView := &entity_cross.GetSubChainCrossView{
			ChainId:   crossChain.ChainId,
			ChainName: crossChain.ChainName,
			TxNum:     crossChain.TxNum,
		}
		sunChainCrossListViews.Add(crossChainView)
	}

	ConvergeListResponse(ctx, sunChainCrossListViews.Values(), int64(len(subChainCrossList)), nil)
}

// CrossUpdateSubChainHandler get
type CrossUpdateSubChainHandler struct {
}

// Handle CrossSubChainListHandler
func (handler *CrossUpdateSubChainHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindCrossUpdateSubChainHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	subChainInfo, err := dbhandle.GetCrossSubChainInfoById(params.ChainId, params.ChainRid)
	if err != nil {
		log.Errorf("GetSubChainCrossChainList err : %v", err)
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	if subChainInfo == nil {
		//insert
		err = insertSubChain(params)
	} else {
		//update
		if params.GatewayId != "" {
			subChainInfo.GatewayId = params.GatewayId
		}
		if params.GatewayName != "" {
			subChainInfo.GatewayName = params.GatewayName
		}
		if params.GatewayAddr != "" {
			subChainInfo.GatewayAddr = params.GatewayAddr
		}
		if params.SpvContractName != "" {
			subChainInfo.SpvContractName = params.SpvContractName
		}
		if params.CrossCa != "" {
			subChainInfo.CrossCa = params.CrossCa
		}
		if params.SdkClientCrt != "" {
			subChainInfo.SdkClientCrt = params.SdkClientCrt
		}
		if params.SdkClientKey != "" {
			subChainInfo.SdkClientKey = params.SdkClientKey
		}
		if params.TxNum != 0 {
			subChainInfo.TxNum = params.TxNum
		}
		if params.BlockHeight != 0 {
			subChainInfo.BlockHeight = params.BlockHeight
		}

		if params.CrossCa != "" && params.SdkClientCrt != "" && params.SdkClientKey != "" {
			//Grpc-获取子链健康状态
			chainOk, errGrpc := sync.CheckSubChainStatus(subChainInfo)
			if errGrpc != nil {
				subChainJson, _ := json.Marshal(subChainInfo)
				log.Errorf("[api update] CheckSubChainStatus failed, err:%v, subChainJson:%v",
					errGrpc, string(subChainJson))
			}
			log.Info("[api update] CheckSubChainStatus chainOk:%v", chainOk)
		}

		err = dbhandle.UpdateCrossSubChainById(params.ChainId, subChainInfo)
	}
	if err != nil {
		log.Errorf("CrossUpdateSubChainHandler insert err : %v, subChainInfo:%v", err, subChainInfo)
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	subChainInfo, err = dbhandle.GetCrossSubChainInfoById(params.ChainId, params.ChainRid)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
	}
	ConvergeDataResponse(ctx, subChainInfo, nil)
}

func insertSubChain(params *entity_cross.CrossUpdateSubChainParams) error {
	subChainList := make([]*db.CrossSubChainData, 0)
	subChainInfo := &db.CrossSubChainData{
		SubChainId:      params.ChainRid,
		ChainId:         params.SubChainId,
		ChainName:       params.SubChainName,
		GatewayId:       params.GatewayId,
		GatewayName:     params.GatewayName,
		GatewayAddr:     params.GatewayAddr,
		ChainType:       params.ChainType, //区块链架构（1 长安链，2 fabric，3 bcos， 4eth，5+ 扩展）
		SpvContractName: params.SpvContractName,
		CrossCa:         params.CrossCa,
		SdkClientCrt:    params.SdkClientCrt,
		SdkClientKey:    params.SdkClientKey,
		Introduction:    params.Introduction,
		ExplorerAddr:    params.ExplorerAddr,
		ExplorerTxAddr:  params.ExplorerTxAddr,
		TxVerifyType:    1,
		Enable:          true,
		TxNum:           params.TxNum,
		BlockHeight:     params.BlockHeight,
		Timestamp:       time.Now().Unix(),
	}
	//Grpc-获取子链健康状态
	chainOk, errGrpc := sync.CheckSubChainStatus(subChainInfo)
	if errGrpc != nil {
		subChainJson, _ := json.Marshal(subChainInfo)
		log.Errorf("[api update] CheckSubChainStatus failed, err:%v, subChainJson:%v",
			errGrpc, string(subChainJson))
	}

	status := dbhandle.SubChainStatusSuccess
	if !chainOk {
		status = dbhandle.SubChainStatusFail
	}
	subChainInfo.Status = status
	subChainList = append(subChainList, subChainInfo)
	err := dbhandle.InsertCrossSubChain(params.ChainId, subChainList)
	return err
}
