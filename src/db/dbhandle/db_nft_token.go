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

// SelectTokenByID 查询token信息
//func SelectTokenByID(chainId string, tokenIds []string) (map[string]*db.NonFungibleToken, error) {
//	tokenMap := make(map[string]*db.NonFungibleToken, 0)
//	tokenList := make([]*db.NonFungibleToken, 0)
//	if len(tokenIds) == 0 {
//		return tokenMap, nil
//	}
//	tableName := db.GetTableName(chainId, db.TableNonFungibleToken)
//	err := db.GormDB.Table(tableName).Where("tokenId in ?", tokenIds).Find(&tokenList).Error
//	if err != nil {
//		return tokenMap, err
//	}
//	for _, token := range tokenList {
//		tokenMap[token.TokenId] = token
//	}
//	return tokenMap, nil
//}

// SelectTokenByID
//
//	@Description: 根据token和合约地址查询token数据是否已经存在
//	@param chainId
//	@param tokenIds tokenid
//	@param contractAddrs 合约地址
//	@return map[string]*db.NonFungibleToken
//	@return error
func SelectTokenByID(chainId string, tokenIds, contractAddrs []string) (map[string]*db.NonFungibleToken, error) {
	tokenMap := make(map[string]*db.NonFungibleToken, 0)
	tokenList := make([]*db.NonFungibleToken, 0)
	if len(tokenIds) == 0 {
		return tokenMap, nil
	}
	tableName := db.GetTableName(chainId, db.TableNonFungibleToken)
	err := db.GormDB.Table(tableName).Where("tokenId IN ? AND contractAddr IN ?",
		tokenIds, contractAddrs).Find(&tokenList).Error
	if err != nil {
		return tokenMap, err
	}
	for _, token := range tokenList {
		tokenKey := token.TokenId + "_" + token.ContractAddr
		tokenMap[tokenKey] = token
	}
	return tokenMap, nil
}

// InsertNonFungibleToken 批量保存
func InsertNonFungibleToken(chainId string, tokenList []*db.NonFungibleToken) error {
	if len(tokenList) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableNonFungibleToken)
	return CreateInBatchesData(tableName, tokenList)
}

// UpdateNonFungibleToken 批量更新token
func UpdateNonFungibleToken(chainId string, tokenInfo *db.NonFungibleToken) error {
	if chainId == "" || tokenInfo == nil || tokenInfo.TokenId == "" {
		return nil
	}

	where := map[string]interface{}{
		"tokenId":      tokenInfo.TokenId,
		"contractAddr": tokenInfo.ContractAddr,
	}
	params := map[string]interface{}{
		//"addrType":  tokenInfo.AddrType,
		"ownerAddr": tokenInfo.OwnerAddr,
	}
	tableName := db.GetTableName(chainId, db.TableNonFungibleToken)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteNonFungibleToken 删除token
func DeleteNonFungibleToken(chainId string, tokenList []*db.NonFungibleToken) error {
	if len(tokenList) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableNonFungibleToken)
	for _, token := range tokenList {
		where := map[string]interface{}{
			"tokenId":      token.TokenId,
			"contractAddr": token.ContractAddr,
		}

		err := db.GormDB.Table(tableName).Where(where).Delete(&db.NonFungibleToken{}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// GetNonFungibleTokenList GetNonFungibleTokenList
func GetNonFungibleTokenList(offset int, limit int, chainId, tokenId, contractAddr string, ownerAddrs []string) (
	[]*db.NonFungibleToken, error) {
	tokenList := make([]*db.NonFungibleToken, 0)
	if chainId == "" {
		return tokenList, db.ErrTableParams
	}
	where := map[string]interface{}{}
	whereIn := map[string]interface{}{}
	if tokenId != "" {
		where["tokenId"] = tokenId
	}
	if contractAddr != "" {
		where["contractAddr"] = contractAddr
	}
	if len(ownerAddrs) > 0 {
		whereIn["ownerAddr"] = ownerAddrs
	}
	selectFile := &SelectFile{
		Where:   where,
		WhereIn: whereIn,
	}
	tableName := db.GetTableName(chainId, db.TableNonFungibleToken)
	query := BuildParamsQuery(tableName, selectFile)
	query = query.Order("timestamp desc").Offset(offset * limit).Limit(limit)
	err := query.Find(&tokenList).Error
	return tokenList, err
}

// GetNFTTokenCount
//
//	@Description: 获取token数量
//	@param chainId
//	@param tokenId
//	@param contractAddr
//	@param ownerAddrs
//	@return int64
//	@return error
func GetNFTTokenCount(chainId, tokenId, contractAddr string, ownerAddrs []string) (int64, error) {
	var count int64
	where := map[string]interface{}{}
	if tokenId != "" {
		where["tokenId"] = tokenId
	}
	if contractAddr != "" {
		where["contractAddr"] = contractAddr
	}

	tableName := db.GetTableName(chainId, db.TableNonFungibleToken)
	query := db.GormDB.Table(tableName)
	if len(where) > 0 {
		query = query.Where(where)
	}
	if len(ownerAddrs) > 0 {
		query = query.Where("ownerAddr in ?", ownerAddrs)
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetNonFungibleTokenDetail GetNonFungibleTokenDetail
func GetNonFungibleTokenDetail(chainId, tokenId, contractAddr string) (*db.NonFungibleToken, error) {
	tokenInfo := &db.NonFungibleToken{}
	if chainId == "" || tokenId == "" || contractAddr == "" {
		return tokenInfo, db.ErrTableParams
	}

	where := map[string]interface{}{
		"tokenId":      tokenId,
		"contractAddr": contractAddr,
	}
	tableName := db.GetTableName(chainId, db.TableNonFungibleToken)
	err := db.GormDB.Table(tableName).Where(where).First(&tokenInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return tokenInfo, nil
}

// UpdateNonFungibleTokenBak 更新敏感词
func UpdateNonFungibleTokenBak(chainId string, tokenData *db.NonFungibleToken) error {
	if chainId == "" || tokenData == nil {
		return nil
	}

	where := map[string]interface{}{
		"tokenId": tokenData.TokenId,
	}
	params := map[string]interface{}{
		"metaData":    tokenData.MetaData,
		"metaDataBak": tokenData.MetaDataBak,
	}
	tableName := db.GetTableName(chainId, db.TableNonFungibleToken)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateTokenContractName 更新合约名称
func UpdateTokenContractName(chainId, contractName, contractAddr string) error {
	if chainId == "" || contractName == "" || contractAddr == "" {
		return nil
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	params := map[string]interface{}{
		"contractName": contractName,
	}
	tableName := db.GetTableName(chainId, db.TableNonFungibleToken)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}
