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

// InsertFungibleContract 插入合约
func InsertFungibleContract(chainId string, contracts []*db.FungibleContract) error {
	if chainId == "" || len(contracts) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableFungibleContract)
	return CreateInBatchesData(tableName, contracts)

}

// UpdateFungibleContractName 更新合约名称
func UpdateFungibleContractName(chainId string, contract *db.Contract) error {
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
	tableName := db.GetTableName(chainId, db.TableFungibleContract)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateFungibleContract 更新合约
func UpdateFungibleContract(chainId string, contract *db.FungibleContract) error {
	if chainId == "" || contract == nil || contract.ContractAddr == "" {
		return nil
	}

	//获取缓存数据
	cachedRes, err := GetFungibleContractByAddr(chainId, contract.ContractAddr)
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

	//是否大于0
	if contract.TotalSupply.GreaterThan(decimal.Zero) {
		params["totalSupply"] = contract.TotalSupply
		cachedRes.TotalSupply = contract.TotalSupply
	}

	if contract.BlockHeight != 0 {
		params["blockHeight"] = contract.BlockHeight
		cachedRes.BlockHeight = contract.BlockHeight
	}
	if contract.TransferNum != 0 {
		params["transferNum"] = contract.TransferNum
		cachedRes.TransferNum = contract.TransferNum
	}

	if len(params) == 0 {
		return nil
	}
	// 获取表名
	tableName := db.GetTableName(chainId, db.TableFungibleContract)
	err = db.GormDB.Table(tableName).Model(&db.FungibleContract{}).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	// 更新缓存
	SetFTContractDataCache(chainId, cachedRes)
	return nil
}

// GetFungibleContractList 同质化合约列表
func GetFungibleContractList(offset, limit int, chainId, contractKey string) (
	[]*db.FungibleContractWithTxNum, int64, error) {
	var count int64
	contracts := make([]*db.FungibleContractWithTxNum, 0)
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
	tableName := db.GetTableName(chainId, db.TableFungibleContract)
	contractTableName := db.GetTableName(chainId, db.TableContract)
	query := BuildParamsQuery(tableName, selectFile)
	query = query.Joins(fmt.Sprintf("LEFT JOIN %s ON %s.contractAddr = %s.addr", contractTableName, tableName,
		contractTableName))
	query = query.Select(fmt.Sprintf("%s.*, %s.txNum", tableName, contractTableName))
	err := query.Count(&count).Error
	if err != nil {
		return contracts, 0, err
	}
	query = query.Order("timestamp desc").Offset(offset * limit).Limit(limit)
	err = query.Find(&contracts).Error
	if err != nil {
		return contracts, 0, err
	}
	return contracts, count, nil
}

// GetFungibleContractByAddr 根据合约地址获取合约信息
func GetFungibleContractByAddr(chainId, contractAddr string) (*db.FungibleContract, error) {
	if chainId == "" || contractAddr == "" {
		return nil, db.ErrTableParams
	}

	contractInfo := GetFTContractDataCache(chainId, contractAddr)
	if contractInfo != nil {
		return contractInfo, nil
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	tableName := db.GetTableName(chainId, db.TableFungibleContract)
	err := db.GormDB.Table(tableName).Where(where).First(&contractInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return contractInfo, nil
}

// QueryFungibleContractExists 根据addr查询
func QueryFungibleContractExists(chainId string, addrList []string) (map[string]*db.FungibleContract, error) {
	contractExists := make(map[string]*db.FungibleContract)
	if len(addrList) == 0 {
		return contractExists, nil
	}

	// 尝试从缓存中获取gas信息
	missingAddrs := make([]string, 0)
	for _, addr := range addrList {
		contractInfo := GetFTContractDataCache(chainId, addr)
		if contractInfo != nil {
			contractExists[contractInfo.ContractAddr] = contractInfo
		} else {
			missingAddrs = append(missingAddrs, addr)
		}
	}

	if len(missingAddrs) == 0 {
		return contractExists, nil
	}

	contracts := make([]*db.FungibleContract, 0)
	tableName := db.GetTableName(chainId, db.TableFungibleContract)
	err := db.GormDB.Table(tableName).Where("contractAddr in ?", missingAddrs).Find(&contracts).Error
	if err != nil {
		return contractExists, err
	}

	for _, v := range contracts {
		contractExists[v.ContractAddr] = v
		SetFTContractDataCache(chainId, v)
	}
	return contractExists, nil
}

// GetFTContractByNameOrAddr 根据合约地址获取合约信息
func GetFTContractByNameOrAddr(chainId, contractKey string) (*db.FungibleContract, error) {
	if chainId == "" || contractKey == "" {
		return nil, db.ErrTableParams
	}

	contractInfo := GetFTContractDataCache(chainId, contractKey)
	if contractInfo != nil {
		return contractInfo, nil
	}

	tableName := db.GetTableName(chainId, db.TableFungibleContract)
	err := db.GormDB.Table(tableName).Where("contractNameBak = ?", contractKey).
		Or("contractAddr = ?", contractKey).First(&contractInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	SetFTContractDataCache(chainId, contractInfo)
	return contractInfo, nil
}
