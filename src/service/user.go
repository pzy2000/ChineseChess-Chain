package service

import (
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
)

// GetUserListHandler handler
type GetUserListHandler struct{}

// Handle deal
func (getUserListHandler *GetUserListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetUserListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetUserList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var (
		userIds   []string
		userAddrs []string
	)

	if params.UserIds != "" {
		userIds = strings.Split(params.UserIds, ",")
	}
	if params.UserAddrs != "" {
		userAddrs = strings.Split(params.UserAddrs, ",")
	}

	//用户列表
	users, totalCount, err := dbhandle.GetUserList(params.Offset, params.Limit, params.ChainId, params.OrgId, userIds,
		userAddrs)
	if err != nil {
		log.Errorf("GetUserList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	startNum := params.Offset * params.Limit
	usersView := arraylist.New()
	for i, user := range users {
		listId := startNum + i + 1
		userListView := &entity.UserListView{
			Id:        strconv.Itoa(listId),
			UserId:    user.UserId,
			OrgId:     user.OrgId,
			Role:      user.Role,
			Timestamp: user.Timestamp,
			UserAddr:  user.UserAddr,
			Status:    user.Status,
		}
		usersView.Add(userListView)
	}

	ConvergeListResponse(ctx, usersView.Values(), totalCount, nil)
}

// GetAccountListHandler handler
type GetAccountListHandler struct{}

// Handle deal
func (getAccountListHandler *GetAccountListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetAccountListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetUserList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取账户列表
	accountList, totalCount, err := dbhandle.GetAccountList(params.Offset, params.Limit,
		params.ChainId, params.AddrType)
	if err != nil {
		log.Errorf("GetAccountList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	accountListView := arraylist.New()
	for _, account := range accountList {
		userListView := &entity.AccountListView{
			AddrType:  account.AddrType,
			Address:   account.Address,
			BNS:       account.BNS,
			DID:       account.DID,
			Timestamp: account.CreatedAt.Unix(),
		}
		accountListView.Add(userListView)
	}

	ConvergeListResponse(ctx, accountListView.Values(), totalCount, nil)
}

// GetAccountDetailHandler handler
type GetAccountDetailHandler struct{}

// Handle deal
func (getAccountDetailHandler *GetAccountDetailHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetAccountDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetUserList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	accountView := &entity.AccountDetailView{}
	//获取账户信息
	accountInfo, err := dbhandle.GetAccountDetail(params.ChainId, params.Address, params.BNS)
	if err != nil || accountInfo == nil {
		log.Errorf("GetAccountDetail err : %s", err)
		//ConvergeHandleFailureResponse(ctx, entity.ErrSelectFailed)
		ConvergeDataResponse(ctx, accountView, nil)
		return
	}

	accountView = &entity.AccountDetailView{
		Address: accountInfo.Address,
		Type:    accountInfo.AddrType,
		BNS:     accountInfo.BNS,
		DID:     accountInfo.DID,
	}

	//返回response
	ConvergeDataResponse(ctx, accountView, nil)
}

// ModifyUserStatusHandler handler
type ModifyUserStatusHandler struct {
}

// Handle deal
func (modifyUserStatusHandler *ModifyUserStatusHandler) Handle(ctx *gin.Context) {
	params := entity.BindModifyUserStatusHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetUserList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}
	// 验证API密钥
	apiKey := ctx.GetHeader("x-api-key")
	if !ValidateAPIKey(apiKey) {
		newError := entity.NewError(entity.ErrorParamWrong, "api key is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	err := dbhandle.UpdateUserStatus(params.Address, params.ChainId, params.Status)
	if err != nil {
		log.Errorf("UpdateUserStatus err : %s", err)
		ConvergeHandleFailureResponse(ctx, entity.ErrUpdateFail)
		return
	}

	ConvergeDataResponse(ctx, "OK", nil)
}
