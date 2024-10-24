package entity_cross

import (
	loggers "chainmaker_web/src/logger"

	"github.com/gin-gonic/gin"
)

var log = loggers.GetLogger(loggers.MODULE_WEB)

//// bindBody
//func bindBody(ctx *gin.Context, body RequestBody) error {
//	if err := ctx.ShouldBindJSON(body); err != nil {
//		log.Error("resolve param error:", err)
//		return err
//	}
//	return nil
//}

// bindParams
func bindParams(ctx *gin.Context, body RequestBody) error {
	requestMethod := ctx.Request.Method
	contentType := ctx.ContentType()

	var err error
	if requestMethod == "POST" && contentType == "application/json" {
		err = ctx.ShouldBindJSON(body)
	} else {
		err = ctx.ShouldBind(body)
	}

	if err != nil {
		log.Error("resolve param error:" + err.Error())
		return err
	}

	//log.Infof("Api Request method: %s, params: %v", requestMethod, body)
	return nil
}

//// bindParams
//func bindParams(ctx *gin.Context, body RequestBody) error {
//	if err := ctx.ShouldBind(body); err != nil {
//		log.Error("resolve param error:" + err.Error())
//		return err
//	}
//	return nil
//}

// BindGetChainIdHandler bind param
func BindGetChainIdHandler(ctx *gin.Context) *GetChainIdParams {
	var body = &GetChainIdParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetCrossSearchHandler bind param
func BindGetCrossSearchHandler(ctx *gin.Context) *CrossSearchParams {
	var body = &CrossSearchParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetCrossTxListHandler bind param
func BindGetCrossTxListHandler(ctx *gin.Context) *GetCrossTxListParams {
	var body = &GetCrossTxListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetCrossSubChainListHandler bind param
func BindGetCrossSubChainListHandler(ctx *gin.Context) *GetCrossSubChainListParams {
	var body = &GetCrossSubChainListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetCrossSubChainDetailHandler bind param
func BindGetCrossSubChainDetailHandler(ctx *gin.Context) *GetCrossSubChainDetailParams {
	var body = &GetCrossSubChainDetailParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetCrossTxDetailHandler bind param
func BindGetCrossTxDetailHandler(ctx *gin.Context) *GetCrossTxDetailParams {
	var body = &GetCrossTxDetailParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetSubChainCrossChainListHandler bind param
func BindGetSubChainCrossChainListHandler(ctx *gin.Context) *GetSubChainCrossChainListParams {
	var body = &GetSubChainCrossChainListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindCrossUpdateSubChainHandler bind param
func BindCrossUpdateSubChainHandler(ctx *gin.Context) *CrossUpdateSubChainParams {
	var body = &CrossUpdateSubChainParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}
