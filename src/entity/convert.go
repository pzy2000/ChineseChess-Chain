/*
Package entity comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package entity

import (
	loggers "chainmaker_web/src/logger"

	"github.com/gin-gonic/gin"
)

var log = loggers.GetLogger(loggers.MODULE_WEB)

// bindBody
func bindBody(ctx *gin.Context, body RequestBody) error {
	if err := ctx.ShouldBindJSON(body); err != nil {
		log.Error("resolve param error:", err)
		return err
	}
	return nil
}

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

// BindGetTransactionNumByTimeHandler bind param
func BindGetTransactionNumByTimeHandler(ctx *gin.Context) *GetTransactionNumByTimeParams {
	var body = &GetTransactionNumByTimeParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetLatestContractHandler bind param
func BindGetLatestContractHandler(ctx *gin.Context) *GetLatestContractParams {
	var body = &GetLatestContractParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetContractListHandler bind param
func BindGetContractListHandler(ctx *gin.Context) *GetContractListParams {
	var body = &GetContractListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetContractDetailHandler bind param
func BindGetContractDetailHandler(ctx *gin.Context) *GetContractDetailParams {
	var body = &GetContractDetailParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetFungibleContractListHandler bind param
func BindGetFungibleContractListHandler(ctx *gin.Context) *GetFTContractListParams {
	var body = &GetFTContractListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetFungibleContractHandler bind param
func BindGetFungibleContractHandler(ctx *gin.Context) *GetFungibleContractParams {
	var body = &GetFungibleContractParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetNonFungibleContractListHandler bind param
func BindGetNonFungibleContractListHandler(ctx *gin.Context) *GetNFTContractListParams {
	var body = &GetNFTContractListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetNonFungibleContractHandler bind param
func BindGetNonFungibleContractHandler(ctx *gin.Context) *GetNonFungibleContractParams {
	var body = &GetNonFungibleContractParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetEvidenceContractListHandler bind param
func BindGetEvidenceContractListHandler(ctx *gin.Context) *GetEvidenceContractListParams {
	var body = &GetEvidenceContractListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetEvidenceContractHandler bind param
func BindGetEvidenceContractHandler(ctx *gin.Context) *GetEvidenceContractParams {
	var body = &GetEvidenceContractParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetIdentityContractListHandler bind param
func BindGetIdentityContractListHandler(ctx *gin.Context) *GetIdentityContractListParams {
	var body = &GetIdentityContractListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetIdentityContractHandler bind param
func BindGetIdentityContractHandler(ctx *gin.Context) *GetIdentityContractParams {
	var body = &GetIdentityContractParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetEventListHandler bind param
func BindGetEventListHandler(ctx *gin.Context) *GetEventListParams {
	var body = &GetEventListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetContractCodeHandler bind param
func BindGetContractCodeHandler(ctx *gin.Context) *GetContractCodeParams {
	var body = &GetContractCodeParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetBlockDetailHandler bind param
func BindGetBlockDetailHandler(ctx *gin.Context) *GetBlockDetailParams {
	var body = &GetBlockDetailParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindChainOverviewDataHandler bind param
func BindChainOverviewDataHandler(ctx *gin.Context) *ChainOverviewDataParams {
	var body = &ChainOverviewDataParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetDetailHandler bind param
func BindGetDetailHandler(ctx *gin.Context) *SearchParams {
	var body = &SearchParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetTxDetailHandler bind param
func BindGetTxDetailHandler(ctx *gin.Context) *GetTxDetailParams {
	var body = &GetTxDetailParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetBlockListHandler bind param
func BindGetBlockListHandler(ctx *gin.Context) *GetBlockListParams {
	var body = &GetBlockListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetTxListHandler bind param
func BindGetTxListHandler(ctx *gin.Context) *GetTxListParams {
	var body = &GetTxListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetBlockTxListHandler bind param
func BindGetBlockTxListHandler(ctx *gin.Context) *GetBlockTxListParams {
	var body = &GetBlockTxListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetContractTxListHandler bind param
func BindGetContractTxListHandler(ctx *gin.Context) *GetContractTxListParams {
	var body = &GetContractTxListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetUserTxListHandler bind param
func BindGetUserTxListHandler(ctx *gin.Context) *GetUserTxListParams {
	var body = &GetUserTxListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetContractVersionListHandler bind param
func BindGetContractVersionListHandler(ctx *gin.Context) *GetContractVersionListParams {
	var body = &GetContractVersionListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetUserListHandler bind param
func BindGetUserListHandler(ctx *gin.Context) *GetUserListParams {
	var body = &GetUserListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetOrgListHandler bind param
func BindGetOrgListHandler(ctx *gin.Context) *GetOrgListParams {
	var body = &GetOrgListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetNodeListHandler bind param
func BindGetNodeListHandler(ctx *gin.Context) *ChainNodesParams {
	var body = &ChainNodesParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetAccountListHandler bind param
func BindGetAccountListHandler(ctx *gin.Context) *GetAccountListParams {
	var body = &GetAccountListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetAccountDetailHandler bind param
func BindGetAccountDetailHandler(ctx *gin.Context) *GetAccountDetailParams {
	var body = &GetAccountDetailParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetGasListHandler bind param
func BindGetGasListHandler(ctx *gin.Context) *GetGasListParams {
	var body = &GetGasListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetGasRecordListHandler bind param
func BindGetGasRecordListHandler(ctx *gin.Context) *GetGasRecordListParams {
	var body = &GetGasRecordListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindInnerGetChainInfoHandler bind param
func BindInnerGetChainInfoHandler(ctx *gin.Context) *InnerGetChainInfoParams {
	var body = &InnerGetChainInfoParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetGasInfoHandler bind param
func BindGetGasInfoHandler(ctx *gin.Context) *GetGasInfoParams {
	var body = &GetGasInfoParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindModifyUserStatusHandler bind param
func BindModifyUserStatusHandler(ctx *gin.Context) *ModifyUserStatusParams {
	var body = &ModifyUserStatusParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindSubscribeChainHandler 订阅链相关
func BindSubscribeChainHandler(ctx *gin.Context) *SubscribeChainParams {
	var body = &SubscribeChainParams{}
	if err := bindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetChainListHandler bind param
func BindGetChainListHandler(ctx *gin.Context) *GetChainListParams {
	var body = &GetChainListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetLatestChainHandler bind param
func BindGetLatestChainHandler(ctx *gin.Context) *GetLatestChainParams {
	var body = &GetLatestChainParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindDeleteSubscribeHandler bind param
func BindDeleteSubscribeHandler(ctx *gin.Context) *DeleteSubscribeParams {
	var body = &DeleteSubscribeParams{}
	if err := bindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindModifySubscribeHandler bind param
func BindModifySubscribeHandler(ctx *gin.Context) *ModifySubscribeParams {
	var body = &ModifySubscribeParams{}
	if err := bindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindCancelSubscribeHandler bind param
func BindCancelSubscribeHandler(ctx *gin.Context) *CancelSubscribeParams {
	var body = &CancelSubscribeParams{}
	if err := bindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetBlockLatestListHandler bind param
func BindGetBlockLatestListHandler(ctx *gin.Context) *GetBlockLatestListParams {
	var body = &GetBlockLatestListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetTxLatestListHandler bind param
func BindGetTxLatestListHandler(ctx *gin.Context) *GetTxLatestListParams {
	var body = &GetTxLatestListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetFungibleTransferListHandler bind param
func BindGetFungibleTransferListHandler(ctx *gin.Context) *GetFungibleTransferListParams {
	var body = &GetFungibleTransferListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetNonFungibleTransferListHandler bind param
func BindGetNonFungibleTransferListHandler(ctx *gin.Context) *GetNonFungibleTransferListParams {
	var body = &GetNonFungibleTransferListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetNFTListHandler bind param
func BindGetNFTListHandler(ctx *gin.Context) *GetNFTListParams {
	var body = &GetNFTListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetNFTDetailHandler bind param
func BindGetNFTDetailHandler(ctx *gin.Context) *GetNFTDetailParams {
	var body = &GetNFTDetailParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetFungiblePositionListHandler bind param
func BindGetFungiblePositionListHandler(ctx *gin.Context) *GetFungiblePositionListParams {
	var body = &GetFungiblePositionListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetUserFTPositionListHandler bind param
func BindGetUserFTPositionListHandler(ctx *gin.Context) *GetUserFTPositionListParams {
	var body = &GetUserFTPositionListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetNonFungiblePositionListHandler bind param
func BindGetNonFungiblePositionListHandler(ctx *gin.Context) *GetNonFungiblePositionListParams {
	var body = &GetNonFungiblePositionListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

//------更新操作--------

// BindModifyTxBlackListHandler bind param
func BindModifyTxBlackListHandler(ctx *gin.Context) *ModifyTxBlackListParams {
	var body = &ModifyTxBlackListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindDeleteTxBlackListHandler bind param
func BindDeleteTxBlackListHandler(ctx *gin.Context) *DeleteTxBlackListParams {
	var body = &DeleteTxBlackListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindUpdateEventSensitiveWordHandler bind param
func BindUpdateEventSensitiveWordHandler(ctx *gin.Context) *UpdateEventSensitiveWordParams {
	var body = &UpdateEventSensitiveWordParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindUpdateTxSensitiveWordHandler bind param
func BindUpdateTxSensitiveWordHandler(ctx *gin.Context) *UpdateTxSensitiveWordParams {
	var body = &UpdateTxSensitiveWordParams{}
	if err := bindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindEvidenceSensitiveWordHandler bind param
func BindEvidenceSensitiveWordHandler(ctx *gin.Context) *EvidenceSensitiveWordParams {
	var body = &EvidenceSensitiveWordParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindNFTSensitiveWordHandler bind param
func BindNFTSensitiveWordHandler(ctx *gin.Context) *NFTSensitiveWordParams {
	var body = &NFTSensitiveWordParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindUpdateContractNameSWHandler bind param
func BindUpdateContractNameSWHandler(ctx *gin.Context) *UpdateContractNameSWParams {
	var body = &UpdateContractNameSWParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindUpdateContractNameSWHandler bind param
func BindGetIDAContractListHandler(ctx *gin.Context) *GetIDAContractListParams {
	var body = &GetIDAContractListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetIDADataListHandler bind param
func BindGetIDADataListHandler(ctx *gin.Context) *GetIDADataListParams {
	var body = &GetIDADataListParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}

// BindGetIDADataDetailHandler bind param
func BindGetIDADataDetailHandler(ctx *gin.Context) *GetIDADataDetailParams {
	var body = &GetIDADataDetailParams{}
	if err := bindParams(ctx, body); err != nil {
		return nil
	}
	return body
}
