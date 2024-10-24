/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
)

// GetIDAContractListHandler get
type GetIDAContractListHandler struct {
}

// Handle GetIDAContractListHandler
func (handler *GetIDAContractListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetIDAContractListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	contractList, count, err := dbhandle.GetIDAContractList(params.Offset, params.Limit, params.ChainId,
		params.ContractKey)
	if err != nil {
		log.Errorf("GetIDAContractList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	contractListView := arraylist.New()
	for _, contract := range contractList {
		contractView := &entity.IDAContractListView{
			ContractName: contract.ContractName,
			ContractAddr: contract.ContractAddr,
			ContractType: contract.ContractType,
			DataAssetNum: contract.TotalNormalAssets,
			Timestamp:    contract.Timestamp,
		}
		contractListView.Add(contractView)
	}
	ConvergeListResponse(ctx, contractListView.Values(), count, nil)
}

// GetIDADataListHandler get
type GetIDADataListHandler struct {
}

// Handle GetIDADataListHandler
func (handler *GetIDADataListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetIDADataListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var (
		err             error
		totalCount      int64
		idaContractInfo *db.IDAContract
	)
	if params.ContractAddr != "" && params.AssetCode == "" {
		//合约数据
		idaContractInfo, err = dbhandle.GetIDAContractByAddr(params.ChainId, params.ContractAddr)
		if idaContractInfo != nil {
			totalCount = idaContractInfo.TotalAssets
		}
	} else {
		totalCount, err = dbhandle.GetIDAAssetCount(params.ChainId, params.ContractAddr, params.AssetCode)
	}
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	assetListView := arraylist.New()
	if totalCount == 0 {
		ConvergeListResponse(ctx, assetListView.Values(), 0, nil)
		return
	}

	assetList, err := dbhandle.GetIDAAssetList(params.Offset, params.Limit, params.ChainId,
		params.ContractAddr, params.AssetCode)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	for _, asset := range assetList {
		assetView := &entity.IDAAssetListView{
			AssetCode:   asset.AssetCode,
			Creator:     asset.Creator,
			IsDeleted:   asset.IsDeleted,
			CreatedTime: asset.CreatedAt.Unix(),
			UpdatedTime: asset.UpdatedAt.Unix(),
		}
		assetListView.Add(assetView)
	}
	ConvergeListResponse(ctx, assetListView.Values(), totalCount, nil)
}

// GetIDADataDetailHandler get
type GetIDADataDetailHandler struct{}

// Handle deal
func (handler *GetIDADataDetailHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetIDADataDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取资产详情
	assetDetail, err := dbhandle.GetIDAAssetDetailByCode(params.ChainId, params.ContractAddr, params.AssetCode)
	if err != nil || assetDetail == nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	//获取资产附件详情
	attachments, err := dbhandle.GetIDAAssetAttachmentByCode(params.ChainId, params.AssetCode)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	annexUrls := make([]string, 0)
	for _, attachment := range attachments {
		annexUrls = append(annexUrls, attachment.Url)
	}

	//获取资产附件详情
	dataAssetsDB, err := dbhandle.GetIDAAssetDataByCode(params.ChainId, params.AssetCode)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}
	apiAssetsDB, err := dbhandle.GetIDAAssetApiByCode(params.ChainId, params.AssetCode)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	dataAssets := make([]entity.DataAsset, 0)
	apiAssets := make([]entity.ApiAsset, 0)
	for _, data := range dataAssetsDB {
		dataAssets = append(dataAssets, entity.DataAsset{
			Name:         data.FieldName,
			Type:         data.FieldType,
			Length:       data.FieldLength,
			IsPrimaryKey: data.IsPrimaryKey,
			IsNotNull:    data.IsNotNull,
			PrivacyQuery: data.PrivacyQuery,
		})
	}

	for _, apiInfo := range apiAssetsDB {
		apiAssets = append(apiAssets, entity.ApiAsset{
			Header:       apiInfo.Header,
			Params:       apiInfo.Params,
			Response:     apiInfo.Response,
			Method:       apiInfo.Method,
			ResponseType: apiInfo.ResponseType,
			Url:          apiInfo.Url,
		})
	}

	assetDetailView := &entity.IDAAssetDetailView{
		AssetCode:         assetDetail.AssetCode,
		AssetName:         assetDetail.AssetName,
		AssetEnName:       assetDetail.AssetEnName,
		Category:          assetDetail.Category,
		ImmediatelySupply: assetDetail.ImmediatelySupply,
		SupplyTime:        assetDetail.SupplyTime,
		DataScale:         assetDetail.DataScale,
		IndustryTitle:     assetDetail.IndustryTitle,
		Summary:           assetDetail.Summary,
		AnnexUrls:         annexUrls,
		UserCategories:    assetDetail.UserCategories,
		UpdateCycleType:   assetDetail.UpdateCycleType,
		TimeSpan:          assetDetail.UpdateTimeSpan,
		DataAsset:         dataAssets,
		ApiAsset:          apiAssets,
		IsDeleted:         assetDetail.IsDeleted,
		CreatedTime:       assetDetail.CreatedAt.Unix(),
		UpdatedTime:       assetDetail.UpdatedAt.Unix(),
	}
	ConvergeDataResponse(ctx, assetDetailView, nil)
}
