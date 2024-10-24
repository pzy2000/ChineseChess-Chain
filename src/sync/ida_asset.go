/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"strconv"
	"strings"
	"time"

	"chainmaker.org/chainmaker/contract-utils/standard"
	"github.com/google/uuid"
)

// DealInsertIDAAssetsData 插入新增的IDA资产数据
func DealInsertIDAAssetsData(idaContractMap map[string]*db.IDAContract,
	idaEventResult *IDAEventData) *db.IDAAssetsDataDB {
	var (
		//insetAssetDetails insetAssetDetails
		insetAssetDetails []*db.IDAAssetDetail
		//insetIDAAttachments insetIDAAttachments
		insetIDAAttachments []*db.IDAAssetAttachment
		//insetAssetDatas insetAssetDatas
		insetAssetDatas []*db.IDADataAsset
		//insetAssetApis insetAssetApis
		insetAssetApis []*db.IDAApiAsset
	)

	if idaEventResult == nil {
		return nil
	}

	//IDACreatedMap新增的IDA资产数据
	for contractAddr, createdInfos := range idaEventResult.IDACreatedMap {
		contractInfo, exists := idaContractMap[contractAddr]
		if !exists {
			continue
		}

		for _, createdInfo := range createdInfos {
			//ida资产数据
			ida := createdInfo.IDAInfo
			//构造资产数据
			processCreatedAssetDetail(createdInfo, contractInfo, &insetAssetDetails)
			//构造Attachments数据
			processFileAttachments(ida, &insetIDAAttachments)
			//构造data和api数据
			processCategorySpecificData(ida, &insetAssetDatas, &insetAssetApis)
		}
	}

	//构造返回值
	return &db.IDAAssetsDataDB{
		IDAAssetDetail:     insetAssetDetails,
		IDAAssetAttachment: insetIDAAttachments,
		IDAAssetData:       insetAssetDatas,
		IDAAssetApi:        insetAssetApis,
	}
}

// processCreatedInfo 构造资产数据
func processCreatedAssetDetail(createdInfo *db.IDACreatedInfo, contractInfo *db.IDAContract,
	insetAssetDetails *[]*db.IDAAssetDetail) {
	//ida资产数据
	ida := createdInfo.IDAInfo
	//资产上链时间
	eventTime := createdInfo.EventTime

	//数据的规模描述1M，1G等
	dataScaleStr := GetIDAAssetDataScale(ida.Details.DataScale)
	// 使用对象类别，1: 政府用户, 2: 企业用户, 3: 个人用户, 4: 无限制用户
	userDescriptions := GetIDAUserCategories(ida.Details.UserCategories)
	//更新周期，1天，1周等
	updateTimeSpan := GetIDAAssetUpdateTimeSpan(ida.Source.UpdateCycle, ida.Details)

	newUUID := uuid.New().String()
	assetDetail := &db.IDAAssetDetail{
		ID:                newUUID,
		AssetCode:         ida.Basic.ID,
		ContractName:      contractInfo.ContractName,
		ContractAddr:      contractInfo.ContractAddr,
		AssetName:         ida.Basic.Name,
		AssetEnName:       ida.Basic.EnName,
		Category:          ida.Basic.Category,
		ImmediatelySupply: ida.Supply.ImmediatelySupply,
		DataScale:         dataScaleStr,
		IndustryTitle:     ida.Basic.Industry.Title,
		Summary:           ida.Basic.Summary,
		Creator:           ida.Basic.Creator,
		Holder:            ida.Ownership.Holder,
		TxID:              ida.Basic.TxID,
		UserCategories:    userDescriptions,
		UpdateCycleType:   ida.Source.UpdateCycle.UpdateCycleType,
		UpdateTimeSpan:    updateTimeSpan,
		CreatedTime:       eventTime,
		UpdatedTime:       eventTime,
	}

	//延迟供应时间
	if ida.Supply.DelayedSupplyTime != nil {
		assetDetail.SupplyTime = *ida.Supply.DelayedSupplyTime
	} else {
		assetDetail.SupplyTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	*insetAssetDetails = append(*insetAssetDetails, assetDetail)
}

// processFileAttachments 构造Attachments插入数据
func processFileAttachments(ida *standard.IDAInfo, insetIDAAttachments *[]*db.IDAAssetAttachment) {
	for _, attachment := range ida.Basic.FileAttachments {
		newUUID := uuid.New().String()
		assetAttachment := &db.IDAAssetAttachment{
			ID:          newUUID,
			AssetCode:   ida.Basic.ID,
			Url:         attachment.Url,
			ContextType: attachment.Type,
		}
		*insetIDAAttachments = append(*insetIDAAttachments, assetAttachment)
	}
}

// processCategorySpecificData 构造数据要素data和api数据
func processCategorySpecificData(ida *standard.IDAInfo, insetAssetDatas *[]*db.IDADataAsset,
	insetAssetApis *[]*db.IDAApiAsset) {
	switch ida.Basic.Category {
	case IDADataCategoryData:
		//构造data资产数据
		processDataCategory(ida, insetAssetDatas)
	case IDADataCategoryAPI:
		//构造api资产数据
		processAPICategory(ida, insetAssetApis)
	}
}

// processDataCategory 构造data资产数据
func processDataCategory(ida *standard.IDAInfo, insetAssetDatas *[]*db.IDADataAsset) {
	for _, columnInfo := range ida.Columns {
		newUUID := uuid.New().String()
		assetData := &db.IDADataAsset{
			ID:           newUUID,
			AssetCode:    ida.Basic.ID,
			FieldName:    columnInfo.Name,
			FieldType:    columnInfo.DataType,
			FieldLength:  columnInfo.DataLength,
			IsPrimaryKey: columnInfo.IsPrimaryKey,
			IsNotNull:    columnInfo.IsNotNull,
			PrivacyQuery: columnInfo.PrivacyQuery,
		}
		*insetAssetDatas = append(*insetAssetDatas, assetData)
	}
}

// processAPICategory 构造api资产数据
func processAPICategory(ida *standard.IDAInfo, insetAssetApis *[]*db.IDAApiAsset) {
	for _, apiInfo := range ida.APIs {
		newUUID := uuid.New().String()
		assetAPI := &db.IDAApiAsset{
			ID:           newUUID,
			AssetCode:    ida.Basic.ID,
			Header:       apiInfo.Header,
			Url:          apiInfo.Url,
			Params:       apiInfo.Params,
			Response:     apiInfo.Response,
			Method:       apiInfo.Method,
			ResponseType: apiInfo.ResponseType,
		}
		*insetAssetApis = append(*insetAssetApis, assetAPI)
	}
}

// GetIDAAssetUpdateTimeSpan 构造更新周期
func GetIDAAssetUpdateTimeSpan(updateCycle standard.UpdateCycle, details standard.Details) string {
	var updateTimeSpan string
	//静态使用查询时间跨度，其他使用更新周期
	if updateCycle.UpdateCycleType == UpdateCycleTypeStatic {
		updateTimeSpan = details.TimeSpan
	} else if updateCycle.UpdateCycleType == UpdateCycleTypePeriodic {
		switch updateCycle.UpdateCycleUnit {
		case IDAUpdateCycleMinute:
			updateTimeSpan = strconv.Itoa(updateCycle.Cycle) + "分钟"
		case IDAUpdateCycleHour:
			updateTimeSpan = strconv.Itoa(updateCycle.Cycle) + "小时"
		case IDAUpdateCycleday:
			updateTimeSpan = strconv.Itoa(updateCycle.Cycle) + "天"
		}
	}

	return updateTimeSpan
}

// GetIDAAssetDataScale 构造数据规模
func GetIDAAssetDataScale(dataScale standard.DataScale) string {
	var dataScaleStr string
	switch dataScale.Type {
	case IDADataScaleTypeNum:
		dataScaleStr = strconv.Itoa(dataScale.Scale) + "条"
	case IDADataScaleTypeM:
		dataScaleStr = strconv.Itoa(dataScale.Scale) + "M"
	case IDADataScaleTypeG:
		dataScaleStr = strconv.Itoa(dataScale.Scale) + "G"
	}
	return dataScaleStr
}

// GetIDAUserCategories 构造使用对象类别，1: 政府用户, 2: 企业用户, 3: 个人用户, 4: 无限制用户
func GetIDAUserCategories(userCategories []int) string {
	var userDescriptions []string
	for _, category := range userCategories {
		description, exists := UserCategoryDescriptions[category]
		if exists {
			userDescriptions = append(userDescriptions, description)
		}
	}

	return strings.Join(userDescriptions, ",")
}

// DealIDAContractUpdateData 更新IDA合约数据
func DealIDAContractUpdateData(idaContractMap map[string]*db.IDAContract,
	idaEventData *IDAEventData, minHeight int64) map[string]*db.IDAContract {
	//更新数据
	idaContractUpdate := make(map[string]*db.IDAContract, 0)
	if idaEventData == nil {
		return idaContractUpdate
	}

	//新增数据资产
	idaCreateds := idaEventData.IDACreatedMap
	//删除数据资产
	idaDeleteIds := idaEventData.IDADeletedCodeMap
	for contractAddr, idaInfos := range idaCreateds {
		idaContract, ok := idaContractMap[contractAddr]
		if !ok || idaContract.BlockHeight >= minHeight {
			//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
			continue
		}

		//更新IDA合约数据资产数量
		totalNormalAssets := idaContract.TotalNormalAssets + int64(len(idaInfos))
		totalAssets := idaContract.TotalAssets + int64(len(idaInfos))
		idaContractUpdate[contractAddr] = &db.IDAContract{
			ContractAddr:      contractAddr,
			TotalNormalAssets: totalNormalAssets,
			TotalAssets:       totalAssets,
			BlockHeight:       minHeight,
		}
	}

	//删除数据资产
	for contractAddr, deleteIds := range idaDeleteIds {
		idaContract, ok := idaContractMap[contractAddr]
		if !ok || idaContract.BlockHeight >= minHeight {
			//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
			continue
		}

		//更新IDA合约数据资产数量
		if _, ok := idaContractUpdate[contractAddr]; ok {
			idaContractUpdate[contractAddr].TotalNormalAssets -= int64(len(deleteIds))
		} else {
			assetNum := idaContract.TotalNormalAssets - int64(len(deleteIds))
			idaContractUpdate[contractAddr] = &db.IDAContract{
				ContractAddr:      contractAddr,
				TotalNormalAssets: assetNum,
				BlockHeight:       minHeight,
			}
		}
	}

	return idaContractUpdate
}

// DealUpdateIDAAssetsData 更新数据资产要素详情
func DealUpdateIDAAssetsData(idaEventResult *IDAEventData,
	assetDetailMap map[string]*db.IDAAssetDetail) *db.IDAAssetsUpdateDB {
	//UpdateAssetDetails 更新资产详情
	updateAssetDetails := make([]*db.IDAAssetDetail, 0)
	//insertIDAAttachments 插入Attachments
	insertIDAAttachments := make([]*db.IDAAssetAttachment, 0)
	//insertAssetDatas 插入data资产
	insertAssetDatas := make([]*db.IDADataAsset, 0)
	//insertAssetApis 插入api资产
	insertAssetApis := make([]*db.IDAApiAsset, 0)
	//deleteAttachmentCodes 删除code
	deleteAttachmentCodes := make([]string, 0)
	//deleteAssetDataCodes 删除code
	deleteAssetDataCodes := make([]string, 0)
	//deleteAssetApiCodes 删除code
	deleteAssetApiCodes := make([]string, 0)

	if idaEventResult == nil {
		return nil
	}

	for assetCode, updateList := range idaEventResult.IDAUpdatedMap {
		assetDetail, ok := assetDetailMap[assetCode]
		if !ok {
			continue
		}

		for _, updateField := range updateList {
			switch updateField.Field {
			case db.KeyIDABasic:
				handleIDABasic(updateField.Update, assetDetail, assetCode, &insertIDAAttachments, &deleteAttachmentCodes)
			case db.KeyIDASUppLy:
				handleIDASupply(updateField.Update, assetDetail)
			case db.KeyIDADetails:
				handleIDADetails(updateField.Update, assetDetail)
			case db.KeyIDAOwnership:
				handleIDAOwnership(updateField.Update, assetDetail)
			case db.KeyIDAColumns:
				handleIDAColumns(updateField.Update, assetCode, &insertAssetDatas, &deleteAssetDataCodes)
			case db.KeyIDAAPI:
				handleIDAAPI(updateField.Update, assetCode, &insertAssetApis, &deleteAssetApiCodes)
			}
			assetDetail.UpdatedTime = updateField.EventTime
		}
		updateAssetDetails = append(updateAssetDetails, assetDetail)
	}

	for _, deleteCodes := range idaEventResult.IDADeletedCodeMap {
		for _, assetCode := range deleteCodes {
			assetDetail, ok := assetDetailMap[assetCode]
			if !ok {
				continue
			}
			assetDetail.IsDeleted = true
			assetDetail.UpdatedTime = idaEventResult.EventTime
			updateAssetDetails = append(updateAssetDetails, assetDetail)
		}
	}

	return &db.IDAAssetsUpdateDB{
		UpdateAssetDetails:    updateAssetDetails,
		InsertAttachment:      insertIDAAttachments,
		InsertIDAAssetData:    insertAssetDatas,
		InsertIDAAssetApi:     insertAssetApis,
		DeleteAttachmentCodes: deleteAttachmentCodes,
		DeleteAssetDataCodes:  deleteAssetDataCodes,
		DeleteAssetApiCodes:   deleteAssetApiCodes,
	}
}

// handleIDABasic 处理数据资产basic更新数据
func handleIDABasic(updateData string, assetDetail *db.IDAAssetDetail, assetCode string,
	insertIDAAttachments *[]*db.IDAAssetAttachment, deleteAttachmentCodes *[]string) {
	basicInfo, err := UnmarshalIDAUpdatedBasic(updateData)
	if err != nil {
		log.Errorf("Error unmarshalling IDABasic:", err)
		return
	}

	assetDetail.AssetName = basicInfo.Name
	assetDetail.AssetEnName = basicInfo.EnName
	assetDetail.Category = basicInfo.Category
	assetDetail.IndustryTitle = basicInfo.Industry.Title
	assetDetail.Summary = basicInfo.Summary
	assetDetail.Creator = basicInfo.Creator
	assetDetail.TxID = basicInfo.TxID

	for _, attachment := range basicInfo.FileAttachments {
		newUUID := uuid.New().String()
		assetAttachment := &db.IDAAssetAttachment{
			ID:          newUUID,
			AssetCode:   assetCode,
			Url:         attachment.Url,
			ContextType: attachment.Type,
		}
		*insertIDAAttachments = append(*insertIDAAttachments, assetAttachment)
	}

	*deleteAttachmentCodes = append(*deleteAttachmentCodes, assetCode)
}

// handleIDASupply 处理数据资产Supply更新数据
func handleIDASupply(updateData string, assetDetail *db.IDAAssetDetail) {
	supplyInfo, err := UnmarshalIDAUpdatedSupply(updateData)
	if err != nil {
		log.Errorf("Error unmarshalling IDASupply:", err)
		return
	}

	assetDetail.ImmediatelySupply = supplyInfo.ImmediatelySupply
	assetDetail.SupplyTime = *supplyInfo.DelayedSupplyTime
}

// handleIDADetails 处理数据资产detail更新数据
func handleIDADetails(updateData string, assetDetail *db.IDAAssetDetail) {
	detailsInfo, err := UnmarshalIDAUpdatedDetails(updateData)
	if err != nil {
		log.Errorf("Error unmarshalling IDADetails:", err)
		return
	}

	assetDetail.DataScale = GetIDAAssetDataScale(detailsInfo.DataScale)
	assetDetail.UserCategories = GetIDAUserCategories(detailsInfo.UserCategories)
}

// handleIDAOwnership 处理数据资产Ownership更新数据
func handleIDAOwnership(updateData string, assetDetail *db.IDAAssetDetail) {
	ownershipInfo, err := UnmarshalIDAUpdatedOwnership(updateData)
	if err != nil {
		log.Errorf("Error unmarshalling IDAOwnership:", err)
		return
	}

	assetDetail.Holder = ownershipInfo.Holder
}

// handleIDAColumns 处理数据资产data数据
func handleIDAColumns(updateData, assetCode string, insertAssetDatas *[]*db.IDADataAsset,
	deleteAssetDataCodes *[]string) {
	columns, err := UnmarshalIDAUpdatedColumns(updateData)
	if err != nil {
		log.Errorf("Error unmarshalling IDAColumns:", err)
		return
	}

	for _, columnInfo := range columns {
		newUUID := uuid.New().String()
		assetData := &db.IDADataAsset{
			ID:           newUUID,
			AssetCode:    assetCode,
			FieldName:    columnInfo.Name,
			FieldType:    columnInfo.DataType,
			FieldLength:  columnInfo.DataLength,
			IsPrimaryKey: columnInfo.IsPrimaryKey,
			IsNotNull:    columnInfo.IsNotNull,
			PrivacyQuery: columnInfo.PrivacyQuery,
		}
		*insertAssetDatas = append(*insertAssetDatas, assetData)
	}
	*deleteAssetDataCodes = append(*deleteAssetDataCodes, assetCode)
}

// handleIDAAPI 处理数据资产api数据
func handleIDAAPI(updateData, assetCode string, insertAssetApis *[]*db.IDAApiAsset, deleteAssetApiCodes *[]string) {
	apiList, err := UnmarshalIDAUpdatedApis(updateData)
	if err != nil {
		log.Errorf("Error unmarshalling IDAAPI:", err)
		return
	}

	for _, apiInfo := range apiList {
		newUUID := uuid.New().String()
		assetAPI := &db.IDAApiAsset{
			ID:           newUUID,
			AssetCode:    assetCode,
			Header:       apiInfo.Header,
			Url:          apiInfo.Url,
			Params:       apiInfo.Params,
			Response:     apiInfo.Response,
			Method:       apiInfo.Method,
			ResponseType: apiInfo.ResponseType,
		}
		*insertAssetApis = append(*insertAssetApis, assetAPI)
	}
	*deleteAssetApiCodes = append(*deleteAssetApiCodes, assetCode)
}
