/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"fmt"

	"gorm.io/gorm"
)

// InsertIDAContract 插入合约
func InsertIDAContract(chainId string, idaContracts []*db.IDAContract) error {
	if chainId == "" || len(idaContracts) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableIDAContract)
	return CreateInBatchesData(tableName, idaContracts)
}

// GetIDAContractMapByAddrs 根据addr查询
func GetIDAContractMapByAddrs(chainId string, addrList []string) (map[string]*db.IDAContract, error) {
	idaContractMap := make(map[string]*db.IDAContract)
	if len(addrList) == 0 {
		return idaContractMap, nil
	}

	contracts := make([]*db.IDAContract, 0)
	tableName := db.GetTableName(chainId, db.TableIDAContract)
	err := db.GormDB.Table(tableName).Where("contractAddr in ?", addrList).Find(&contracts).Error
	if err != nil {
		return idaContractMap, err
	}

	for _, v := range contracts {
		idaContractMap[v.ContractAddr] = v
	}
	return idaContractMap, nil
}

// UpdateIDAContractByAddr 更新更新状态
func UpdateIDAContractByAddr(chainId, contractAddr string, idaContract *db.IDAContract) error {
	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}

	params := map[string]interface{}{}
	if idaContract.TotalNormalAssets > 0 {
		params["totalNormalAssets"] = idaContract.TotalNormalAssets
	}
	if idaContract.TotalAssets > 0 {
		params["totalAssets"] = idaContract.TotalAssets
	}
	if idaContract.BlockHeight > 0 {
		params["blockHeight"] = idaContract.BlockHeight
	}
	if len(params) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableIDAContract)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	return err
}

// GetIDAContractList 获取合约列表
func GetIDAContractList(offset, limit int, chainId, contractKey string) ([]*db.IDAContract, int64, error) {
	var count int64
	contracts := make([]*db.IDAContract, 0)
	tableName := db.GetTableName(chainId, db.TableIDAContract)
	query := db.GormDB.Table(tableName)
	if contractKey != "" {
		query = query.Where("contractNameBak = ? or contractAddr = ?", contractKey, contractKey)
	}

	err := query.Count(&count).Error
	if err != nil {
		return contracts, 0, err
	}

	query = query.Order("timestamp desc")
	query = query.Offset(offset * limit).Limit(limit)
	err = query.Find(&contracts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("GetIDAContractList err, cause : %s", err.Error())
	}
	return contracts, count, nil
}

// GetIDAContractByAddr 根据合约名称获取合约
func GetIDAContractByAddr(chainId, contractAddr string) (*db.IDAContract, error) {
	if chainId == "" || contractAddr == "" {
		return nil, db.ErrTableParams
	}

	contractInfo := &db.IDAContract{}
	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	tableName := db.GetTableName(chainId, db.TableIDAContract)
	err := db.GormDB.Table(tableName).Where(where).First(&contractInfo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return contractInfo, nil
}
