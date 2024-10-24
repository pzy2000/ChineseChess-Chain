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

// InsertEvidenceContract 插入合约
func InsertEvidenceContract(chainId string, contracts []*db.EvidenceContract) error {
	if len(contracts) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableEvidenceContract)
	return CreateInBatchesData(tableName, contracts)
}

// InsertIdentityContract 插入合约
func InsertIdentityContract(chainId string, contracts []*db.IdentityContract) error {
	if len(contracts) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableIdentityContract)
	return CreateInBatchesData(tableName, contracts)
}

// GetEvidenceContract 根据合约地址获取合约信息
func GetEvidenceContract(offset, limit, resCode int, chainId, contractName, txId, searchValue string,
	hashList, senderAddrs []string) ([]*db.EvidenceContract, int64, error) {
	var count int64
	evidenceList := make([]*db.EvidenceContract, 0)
	if chainId == "" {
		return evidenceList, 0, db.ErrTableParams
	}
	whereIn := map[string]interface{}{}
	where := map[string]interface{}{}
	whereOr := map[string]interface{}{}
	if contractName != "" {
		where["contractName"] = contractName
	}
	if txId != "" {
		where["txId"] = txId
	}
	if resCode > 0 {
		where["contractResultCode"] = resCode - 1
	}
	if len(senderAddrs) > 0 {
		whereIn["senderAddr"] = senderAddrs
	}
	if len(hashList) > 0 {
		whereIn["hash"] = hashList
	}

	if searchValue != "" {
		whereOr["txId"] = searchValue
		whereOr["senderAddr"] = searchValue
		whereOr["hash"] = searchValue
	}

	selectFile := &SelectFile{
		Where:   where,
		WhereIn: whereIn,
		WhereOr: whereOr,
	}
	tableName := db.GetTableName(chainId, db.TableEvidenceContract)
	query := BuildParamsQuery(tableName, selectFile)
	err := query.Count(&count).Error
	if err != nil {
		return evidenceList, 0, err
	}

	query = query.Order("timestamp desc").Offset(offset * limit).Limit(limit)
	err = query.Find(&evidenceList).Error
	if err != nil {
		return evidenceList, 0, err
	}
	return evidenceList, count, nil
}

// UpdateEvidenceContractName 更新合约名称
func UpdateEvidenceContractName(chainId, contractName, contractNameNew string) error {
	if chainId == "" || contractName == "" || contractNameNew == "" {
		return nil
	}

	where := map[string]interface{}{
		"contractName": contractName,
	}
	params := map[string]interface{}{
		"contractName": contractNameNew,
	}
	tableName := db.GetTableName(chainId, db.TableEvidenceContract)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// GetEvidenceContractByHash 根据合约地址获取合约信息
func GetEvidenceContractByHash(chainId, hash string) (*db.EvidenceContract, error) {
	evidenceContract := &db.EvidenceContract{}
	where := map[string]interface{}{
		"hash": hash,
	}
	tableName := db.GetTableName(chainId, db.TableEvidenceContract)
	err := db.GormDB.Table(tableName).Where(where).First(&evidenceContract).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return evidenceContract, nil
}

// GetEvidenceContractByHashLit 根据hash查询存证合约
func GetEvidenceContractByHashLit(chainId string, hashList []string) (map[string]*db.EvidenceContract, error) {
	evidenceExists := make(map[string]*db.EvidenceContract, 0)
	evidenceList := make([]*db.EvidenceContract, 0)
	if len(hashList) == 0 {
		return evidenceExists, nil
	}
	tableName := db.GetTableName(chainId, db.TableEvidenceContract)
	err := db.GormDB.Table(tableName).Where("hash IN ?", hashList).Find(&evidenceList).Error
	if err != nil {
		return evidenceExists, err
	}

	for _, evidence := range evidenceList {
		evidenceExists[evidence.Hash] = evidence
	}
	return evidenceExists, nil
}

// UpdateEvidenceBak 更新敏感词
func UpdateEvidenceBak(chainId string, evidence *db.EvidenceContract) error {
	if chainId == "" || evidence == nil {
		return nil
	}

	where := map[string]interface{}{
		"hash": evidence.Hash,
	}
	params := map[string]interface{}{
		"metaData":    evidence.MetaData,
		"metaDataBak": evidence.MetaDataBak,
	}
	tableName := db.GetTableName(chainId, db.TableEvidenceContract)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

//// GetIdentityContractList 身份认证合约列表
//func GetIdentityContractList(offset, limit int, chainId, contractKey string) (
//	[]*db.EvidenceList, int64, error) {
//	evidenceList := make([]*db.EvidenceList, 0)
//	whereOr := map[string]interface{}{}
//	if contractKey != "" {
//		whereOr = map[string]interface{}{
//			"contractName": contractKey,
//			"contractAddr": contractKey,
//		}
//	}
//	selectFile := &db.SelectFile{
//		WhereOr:    whereOr,
//		ResultType: &evidenceList,
//		SortField:  "createdAt",
//		SortAsc:    false,
//		Offset:     offset,
//		Size:       limit,
//	}
//	tableName := db.GetTableName(chainId, db.TableIdentityContract)
//	count, err := db.StorageClient.SelectGroupContractList(tableName, selectFile)
//	if err != nil {
//		return nil, 0, fmt.Errorf("GetEvidenceContractList err, cause : %s", err.Error())
//	}
//	return evidenceList, count, nil
//}

// GetIdentityContract 根据合约地址获取合约信息
func GetIdentityContract(offset, limit int, chainId, contractAddr string, userAddrs []string) (
	[]*db.IdentityContract, int64, error) {
	var count int64
	contractList := make([]*db.IdentityContract, 0)
	if chainId == "" || contractAddr == "" {
		return contractList, 0, db.ErrTableParams
	}
	whereIn := map[string]interface{}{}
	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	if len(userAddrs) > 0 {
		whereIn["userAddr"] = userAddrs
	}
	selectFile := &SelectFile{
		Where:   where,
		WhereIn: whereIn,
	}
	tableName := db.GetTableName(chainId, db.TableIdentityContract)
	query := BuildParamsQuery(tableName, selectFile)
	err := query.Count(&count).Error
	if err != nil {
		return contractList, 0, err
	}
	query = query.Offset(offset * limit).Limit(limit)
	err = query.Find(&contractList).Error
	if err != nil {
		return contractList, 0, fmt.Errorf("GetIdentityContract err, cause : %s", err.Error())
	}

	return contractList, count, nil
}

// UpdateIdentityContractName 更新合约名称
func UpdateIdentityContractName(chainId, contractName, contractAddr string) error {
	if chainId == "" || contractName == "" || contractAddr == "" {
		return nil
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	params := map[string]interface{}{
		"contractName": contractName,
	}
	tableName := db.GetTableName(chainId, db.TableIdentityContract)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}
