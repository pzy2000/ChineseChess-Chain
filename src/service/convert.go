/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chainmaker_web/src/entity"
)

// ConvergeDataResponse 汇聚单一对象应答结果
func ConvergeDataResponse(ctx *gin.Context, data interface{}, err *entity.Error) {
	// 首先判断err是否为空
	if err == nil {
		successResponse := entity.NewSuccessDataResponse(data)
		//jsonStr, _ := json.Marshal(successResponse)
		//log.Infof("api return response: %v", string(jsonStr))
		ctx.JSON(http.StatusOK, successResponse)
	} else {
		ConvergeFailureResponse(ctx, err)
	}
}

// ConvergeListResponse 汇聚集合对象应答结果
func ConvergeListResponse(ctx *gin.Context, datas []interface{}, count int64, err *entity.Error) {
	// 首先判断err是否为空
	if err == nil {
		successResponse := entity.NewSuccessListResponse(datas, count)
		//jsonStr, _ := json.Marshal(successResponse)
		//log.Infof("Api Return response: %v", string(jsonStr))
		ctx.JSON(http.StatusOK, successResponse)
	} else {
		ConvergeFailureResponse(ctx, err)
	}
}

// ConvergeFailureResponse 汇聚失败应答
func ConvergeFailureResponse(ctx *gin.Context, err *entity.Error) {
	log.Errorf("Http request[%s]'s error = [%s]", ctx.Request.URL.String(), err.Error())
	failureResponse := entity.NewFailureResponse(err)
	ctx.JSON(http.StatusOK, failureResponse)
}

// ConvergeHandleFailureResponse 汇聚处理异常的应答
func ConvergeHandleFailureResponse(ctx *gin.Context, err error) {
	newError := entity.NewError(entity.ErrorHandleFailure, err.Error())
	ConvergeFailureResponse(ctx, newError)
}
