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

	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

// InsertNonFungibleContract 插入合约
func InsertNonFungibleContract(chainId string, contracts []*db.NonFungibleContract) error {
	if chainId == "" || len(contracts) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableNonFungibleContract)
	return CreateInBatchesData(tableName, contracts)
}

// UpdateNonFungibleContract 更新合约
func UpdateNonFungibleContract(chainId string, contract *db.NonFungibleContract) error {
	if chainId == "" || contract == nil || contract.ContractAddr == "" {
		return nil
	}

	//获取缓存数据
	cachedRes, err := GetNonFungibleContractByAddr(chainId, contract.ContractAddr)
	if err != nil {
		return err
	}

	where := map[string]interface{}{
		"contractAddr": contract.ContractAddr,
	}
	params := map[string]interface{}{}
	if contract.HolderCount > 0 {
		params["holderCount"] = contract.HolderCount
		cachedRes.HolderCount = contract.HolderCount
	}
	if contract.TransferNum > 0 {
		params["transferNum"] = contract.TransferNum
		cachedRes.TransferNum = contract.TransferNum
	}
	//是否大于0
	if contract.TotalSupply.GreaterThan(decimal.Zero) {
		params["totalSupply"] = contract.TotalSupply
		cachedRes.TotalSupply = contract.TotalSupply
	}
	if contract.BlockHeight > 0 {
		params["blockHeight"] = contract.BlockHeight
		cachedRes.BlockHeight = contract.BlockHeight
	}

	if len(params) == 0 {
		return nil
	}

	// 获取表名
	tableName := db.GetTableName(chainId, db.TableNonFungibleContract)
	err = db.GormDB.Table(tableName).Model(&db.NonFungibleContract{}).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	//更新缓存
	SetNFTContractDataCache(chainId, cachedRes)
	return nil
}

// GetNonFungibleContractList 非同质化合约列表
func GetNonFungibleContractList(offset, limit int, chainId, contractKey string) (
	[]*db.NonFungibleContractWithTxNum, int64, error) {
	var count int64
	contracts := make([]*db.NonFungibleContractWithTxNum, 0)
	whereOr := map[string]interface{}{}
	if contractKey != "" {
		whereOr = map[string]interface{}{
			"contractNameBak": contractKey,
			"contractAddr":    contractKey,
		}
	}

	selectFile := &SelectFile{
		WhereOr: whereOr,
	}
	tableName := db.GetTableName(chainId, db.TableNonFungibleContract)
	contractTableName := db.GetTableName(chainId, db.TableContract)
	query := BuildParamsQuery(tableName, selectFile)
	query = query.Joins(fmt.Sprintf("LEFT JOIN %s ON %s.contractAddr = %s.addr", contractTableName, tableName,
		contractTableName))
	query = query.Select(fmt.Sprintf("%s.*, %s.txNum", tableName, contractTableName))
	err := query.Count(&count).Error
	if err != nil {
		return contracts, 0, err
	}
	query = query.Order("timestamp desc")
	query = query.Offset(offset * limit).Limit(limit)
	err = query.Find(&contracts).Error
	if err != nil {
		return contracts, 0, err
	}

	return contracts, count, nil
}

// GetNonFungibleContractByAddr 根据合约地址获取合约信息
func GetNonFungibleContractByAddr(chainId, contractAddr string) (*db.NonFungibleContract, error) {
	var contractInfo *db.NonFungibleContract
	if chainId == "" || contractAddr == "" {
		return nil, db.ErrTableParams
	}

	contractInfo = GetNFTContractDataCache(chainId, contractAddr)
	if contractInfo != nil {
		return contractInfo, nil
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	tableName := db.GetTableName(chainId, db.TableNonFungibleContract)
	err := db.GormDB.Table(tableName).Where(where).First(&contractInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return contractInfo, nil
}

// UpdateNonFungibleContractName 更新合约名称
func UpdateNonFungibleContractName(chainId string, contract *db.Contract) error {
	if chainId == "" || contract == nil || contract.Addr == "" {
		return nil
	}

	where := map[string]interface{}{
		"contractAddr": contract.Addr,
	}
	params := map[string]interface{}{
		"contractName":    contract.Name,
		"contractNameBak": contract.NameBak,
	}
	tableName := db.GetTableName(chainId, db.TableNonFungibleContract)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// QueryNonFungibleContractExists 根据addr查询
func QueryNonFungibleContractExists(chainId string, addrList []string) (map[string]*db.NonFungibleContract, error) {
	contractExists := make(map[string]*db.NonFungibleContract)
	if len(addrList) == 0 {
		return contractExists, nil
	}

	// 尝试从缓存中获取gas信息
	missingAddrs := make([]string, 0)
	for _, addr := range addrList {
		contractInfo := GetNFTContractDataCache(chainId, addr)
		if contractInfo != nil {
			contractExists[contractInfo.ContractAddr] = contractInfo
		} else {
			missingAddrs = append(missingAddrs, addr)
		}
	}

	if len(missingAddrs) == 0 {
		return contractExists, nil
	}

	contracts := make([]*db.NonFungibleContract, 0)
	tableName := db.GetTableName(chainId, db.TableNonFungibleContract)
	err := db.GormDB.Table(tableName).Where("contractAddr in ?", addrList).Find(&contracts).Error
	if err != nil {
		return contractExists, err
	}

	for _, v := range contracts {
		contractExists[v.ContractAddr] = v
		SetNFTContractDataCache(chainId, v)
	}
	return contractExists, nil
}

// GetNFTContractByNameOrAddr 根据合约地址获取合约信息
func GetNFTContractByNameOrAddr(chainId, contractKey string) (*db.NonFungibleContract, error) {
	if chainId == "" || contractKey == "" {
		return nil, db.ErrTableParams
	}

	contractInfo := GetNFTContractDataCache(chainId, contractKey)
	if contractInfo != nil {
		return contractInfo, nil
	}

	tableName := db.GetTableName(chainId, db.TableNonFungibleContract)
	err := db.GormDB.Table(tableName).Where("contractNameBak = ?", contractKey).
		Or("contractAddr = ?", contractKey).First(&contractInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	SetNFTContractDataCache(chainId, contractInfo)
	return contractInfo, nil
}
