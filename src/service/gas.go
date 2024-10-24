package service

import (
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
)

// GetGasListHandler handler
type GetGasListHandler struct {
}

// Handle deal
func (getGasListHandler *GetGasListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetGasListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var userAddrs []string
	if params.UserAddrs != "" {
		userAddrs = strings.Split(params.UserAddrs, ",")
	}
	//gas列表
	gasList, totalCount, err := dbhandle.GetGasList(params.Offset, params.Limit, params.ChainId, userAddrs)
	if err != nil {
		log.Errorf("GetGasList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	gasesView := arraylist.New()
	for i, gas := range gasList {
		listId := params.Offset*params.Limit + i + 1
		userListView := &entity.GasListView{
			Id:         strconv.Itoa(listId),
			ChainId:    params.ChainId,
			Timestamp:  gas.CreatedAt.Unix(),
			GasBalance: gas.GasBalance,
			GasTotal:   gas.GasTotal,
			GasUsed:    gas.GasUsed,
			Address:    gas.Address,
		}
		gasesView.Add(userListView)
	}

	ConvergeListResponse(ctx, gasesView.Values(), totalCount, nil)
}

// GetGasRecordListHandler handler
type GetGasRecordListHandler struct {
}

// Handle deal
func (getGasRecordListHandler *GetGasRecordListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetGasRecordListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetUserList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var userAddrs []string
	if params.UserAddrs != "" {
		userAddrs = strings.Split(params.UserAddrs, ",")
	}
	//GasRecord列表
	gases, totalCount, err := dbhandle.GetGasRecordList(params.Offset, params.Limit, params.ChainId, userAddrs,
		params.StartTime, params.EndTime, params.BusinessType)
	if err != nil {
		log.Errorf("GetGasRecordList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	usersView := arraylist.New()
	for i, gas := range gases {
		listId := params.Offset*params.Limit + i + 1
		userListView := &entity.GasRecordListView{
			Id:        strconv.Itoa(listId),
			ChainId:   params.ChainId,
			GasAmount: gas.GasAmount,
			Address:   gas.Address,
			//PayerAddress: gas.PayerAddress,
			BusinessType: gas.BusinessType,
			TxId:         gas.TxId,
			Timestamp:    gas.Timestamp,
		}
		usersView.Add(userListView)
	}
	ConvergeListResponse(ctx, usersView.Values(), totalCount, nil)
}

// GetGasInfoHandler handler
type GetGasInfoHandler struct {
}

// Handle deal
func (getGasInfoHandler *GetGasInfoHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetGasInfoHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetUserList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}
	var userAddrs []string
	if params.UserAddrs != "" {
		userAddrs = strings.Split(params.UserAddrs, ",")
	}
	//根据地址获取gas数据
	gasBalance, err := dbhandle.GetGasByAddrInfo(params.ChainId, userAddrs)
	if err != nil {
		log.Errorf("GetGasByAddrInfo err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	ConvergeDataResponse(ctx, &entity.GasInfoView{GasBalance: gasBalance}, nil)
}
