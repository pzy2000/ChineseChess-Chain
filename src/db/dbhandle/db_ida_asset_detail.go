/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// InsertIDADetail
func InsertIDADetail(chainId string, idaDetails []*db.IDAAssetDetail) error {
	if chainId == "" || len(idaDetails) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableIDAAssetDetail)
	return CreateInBatchesData(tableName, idaDetails)
}

// GetIDAAssetList 获取资产列表
func GetIDAAssetList(offset, limit int, chainId, contractAddr, assetCode string) ([]*db.IDAAssetDetail, error) {
	assetList := make([]*db.IDAAssetDetail, 0)
	tableName := db.GetTableName(chainId, db.TableIDAAssetDetail)
	where := map[string]interface{}{}
	if contractAddr != "" {
		where["contractAddr"] = contractAddr
	}
	if assetCode != "" {
		where["assetCode"] = assetCode
	}
	query := db.GormDB.Table(tableName).Where(where).Order("updatedAt desc")
	query = query.Offset(offset * limit).Limit(limit)
	err := query.Find(&assetList).Error
	if err != nil {
		return assetList, fmt.Errorf("GetIDAAssetList err, cause : %s", err.Error())
	}
	return assetList, nil
}

// GetIDAAssetCount 获取资产数量
func GetIDAAssetCount(chainId, contractAddr, assetCode string) (int64, error) {
	var count int64
	tableName := db.GetTableName(chainId, db.TableIDAAssetDetail)
	where := map[string]interface{}{}
	if contractAddr != "" {
		where["contractAddr"] = contractAddr
	}
	if assetCode != "" {
		where["assetCode"] = assetCode
	}

	query := db.GormDB.Table(tableName).Where(where)
	err := query.Count(&count).Error
	return count, err
}

// GetIDAAssetDetail 获取资产详情
func GetIDAAssetDetailByCode(chainId, contractAddr, assetCode string) (*db.IDAAssetDetail, error) {
	assetDetail := &db.IDAAssetDetail{}
	if chainId == "" || assetCode == "" || contractAddr == "" {
		return assetDetail, db.ErrTableParams
	}

	where := map[string]interface{}{
		"assetCode":    assetCode,
		"contractAddr": contractAddr,
	}
	tableName := db.GetTableName(chainId, db.TableIDAAssetDetail)
	err := db.GormDB.Table(tableName).Where(where).First(&assetDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return assetDetail, nil
}

// UpdateNonFungibleTokenBak 更新敏感词
func UpdateIDADetailByCode(chainId string, idaDetail *db.IDAAssetDetail) error {
	if chainId == "" || idaDetail == nil {
		return nil
	}

	where := map[string]interface{}{
		"assetCode": idaDetail.AssetCode,
	}

	params := map[string]interface{}{
		"assetName":         idaDetail.AssetName,
		"assetEnName":       idaDetail.AssetEnName,
		"category":          idaDetail.Category,
		"immediatelySupply": idaDetail.ImmediatelySupply,
		"supplyTime":        idaDetail.SupplyTime,
		"dataScale":         idaDetail.DataScale,
		"industryTitle":     idaDetail.IndustryTitle,
		"summary":           idaDetail.Summary,
		"creator":           idaDetail.Creator,
		"holder":            idaDetail.Holder,
		"txId":              idaDetail.TxID,
		"userCategories":    idaDetail.UserCategories,
		"UpdatedTime":       idaDetail.UpdatedTime,
		"isDeleted":         idaDetail.IsDeleted,
	}
	tableName := db.GetTableName(chainId, db.TableIDAAssetDetail)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// GetIDAAssetDetail 获取资产详情
func GetIDAAssetDetailMapByCodes(chainId string, assetCodes []string) (map[string]*db.IDAAssetDetail, error) {
	assetList := make([]*db.IDAAssetDetail, 0)
	assetMap := make(map[string]*db.IDAAssetDetail, 0)
	if chainId == "" || len(assetCodes) == 0 {
		return assetMap, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableIDAAssetDetail)
	err := db.GormDB.Table(tableName).Where("assetCode in ?", assetCodes).Find(&assetList).Error
	if err != nil {
		return nil, err
	}

	for _, asset := range assetList {
		assetMap[asset.AssetCode] = asset
	}
	return assetMap, nil
}
