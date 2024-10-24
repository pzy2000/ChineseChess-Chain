/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"errors"

	"gorm.io/gorm"
)

// InsertIDAAttachments
func InsertIDAAttachments(chainId string, idaAttachments []*db.IDAAssetAttachment) error {
	if chainId == "" || len(idaAttachments) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableIDAAssetAttachment)
	return CreateInBatchesData(tableName, idaAttachments)
}

// DeleteIDAAttachments DeleteIDAAttachments
func DeleteIDAAttachments(chainId string, assetCodes []string) error {
	if len(assetCodes) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableIDAAssetAttachment)
	err := db.GormDB.Table(tableName).Where("assetCode in ?", assetCodes).Delete(nil).Error
	return err
}

// InsertIDAAssetApi
func InsertIDAAssetApi(chainId string, idaAssetApis []*db.IDAApiAsset) error {
	if chainId == "" || len(idaAssetApis) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableIDAApiAsset)
	return CreateInBatchesData(tableName, idaAssetApis)
}

// DeleteIDAApis DeleteIDAApis
func DeleteIDAApis(chainId string, assetCodes []string) error {
	if len(assetCodes) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableIDAApiAsset)
	err := db.GormDB.Table(tableName).Where("assetCode in ?", assetCodes).Delete(nil).Error
	return err
}

// InsertIDAAssetData
func InsertIDAAssetData(chainId string, idaDatas []*db.IDADataAsset) error {
	if chainId == "" || len(idaDatas) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableIDADataAsset)
	return CreateInBatchesData(tableName, idaDatas)
}

// DeleteIDADatas DeleteIDADatas
func DeleteIDADatas(chainId string, assetCodes []string) error {
	if len(assetCodes) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableIDADataAsset)
	err := db.GormDB.Table(tableName).Where("assetCode in ?", assetCodes).Delete(nil).Error
	return err
}

// GetIDAAssetDetail 获取资产详情
func GetIDAAssetAttachmentByCode(chainId, assetCode string) ([]*db.IDAAssetAttachment, error) {
	attachments := make([]*db.IDAAssetAttachment, 0)
	if chainId == "" || assetCode == "" {
		return attachments, db.ErrTableParams
	}

	where := map[string]interface{}{
		"assetCode": assetCode,
	}
	tableName := db.GetTableName(chainId, db.TableIDAAssetAttachment)
	err := db.GormDB.Table(tableName).Where(where).Find(&attachments).Error
	if err != nil {
		return attachments, err
	}

	return attachments, nil
}

// GetIDAAssetDetail 获取资产详情
func GetIDAAssetDataByCode(chainId, assetCode string) ([]*db.IDADataAsset, error) {
	dataAssets := make([]*db.IDADataAsset, 0)
	if chainId == "" || assetCode == "" {
		return dataAssets, db.ErrTableParams
	}

	where := map[string]interface{}{
		"assetCode": assetCode,
	}
	tableName := db.GetTableName(chainId, db.TableIDADataAsset)
	err := db.GormDB.Table(tableName).Where(where).Find(&dataAssets).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return dataAssets, nil
}

// GetIDAAssetDetail 获取资产详情
func GetIDAAssetApiByCode(chainId, assetCode string) ([]*db.IDAApiAsset, error) {
	apiAssets := make([]*db.IDAApiAsset, 0)
	if chainId == "" || assetCode == "" {
		return apiAssets, db.ErrTableParams
	}

	where := map[string]interface{}{
		"assetCode": assetCode,
	}
	tableName := db.GetTableName(chainId, db.TableIDAApiAsset)
	err := db.GormDB.Table(tableName).Where(where).Find(&apiAssets).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return apiAssets, nil
}
