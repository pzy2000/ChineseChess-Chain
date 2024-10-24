/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"fmt"
)

/*=======同质化表sql=========*/

// InsertFungibleTransfer 批量保存
func InsertFungibleTransfer(chainId string, transfers []*db.FungibleTransfer) error {
	if len(transfers) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableFungibleTransfer)
	return CreateInBatchesData(tableName, transfers)
}

// GetFTTransferList 同质化流转列表
func GetFTTransferList(offset, limit int, chainId, contractAddr, userAddr string) (
	[]*db.FungibleTransfer, error) {
	transfers := make([]*db.FungibleTransfer, 0)
	if chainId == "" || contractAddr == "" {
		return transfers, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableFungibleTransfer)
	if userAddr == "" {
		where := map[string]interface{}{
			"contractAddr": contractAddr,
		}
		// 当userAddr为空时，执行简单查询
		query := db.GormDB.Table(tableName).Where(where)
		err := query.Order("timestamp desc").Offset(offset * limit).Limit(limit).Find(&transfers).Error
		if err != nil {
			return transfers, fmt.Errorf("GetFTTransferList err, cause : %s", err.Error())
		}
	} else {
		subQuerySql := GetTransferSubQuerySql(tableName, tableName)
		// 当userAddr不为空时，执行UNION查询
		sqlQuery := fmt.Sprintf(`
			SELECT * FROM (%s) AS t
			ORDER BY timestamp DESC
			LIMIT ? OFFSET ?;
		`, subQuerySql)

		// 执行原始SQL查询并将结果映射到transfers切片
		query := db.GormDB.Raw(sqlQuery, contractAddr, userAddr, contractAddr, userAddr, limit, offset*limit)
		err := query.Scan(&transfers).Error
		if err != nil {
			return transfers, fmt.Errorf("GetFTTransferList err, cause : %s", err.Error())
		}
	}

	return transfers, nil
}

func GetFTTransferCount(chainId, contractAddr, userAddr string) (int64, error) {
	var count int64
	tableName := db.GetTableName(chainId, db.TableFungibleTransfer)
	if userAddr == "" {
		where := map[string]interface{}{
			"contractAddr": contractAddr,
		}
		// 当userAddr为空时，执行简单查询
		query := db.GormDB.Table(tableName).Where(where)
		err := query.Count(&count).Error
		if err != nil {
			return 0, err
		}
	} else {
		subQuerySql := GetTransferSubQuerySql(tableName, tableName)
		// 当userAddr不为空时，执行UNION查询
		sqlQuery := fmt.Sprintf(`
			SELECT COUNT(*) FROM (%s) AS t;
		`, subQuerySql)

		// 执行原始SQL查询并获取结果数量
		query := db.GormDB.Raw(sqlQuery, contractAddr, userAddr, contractAddr, userAddr)
		err := query.Row().Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("GetFTTransferCount err, cause : %s", err.Error())
		}
	}

	return count, nil
}

// UpdateTransferContractName 更新合约名称
func UpdateTransferContractName(chainId, contractName, contractAddr string) error {
	if chainId == "" || contractName == "" || contractAddr == "" {
		return nil
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	params := map[string]interface{}{
		"contractName": contractName,
	}
	tableName := db.GetTableName(chainId, db.TableFungibleTransfer)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

/*=======非同质化表sql=========*/

//func GetNFTTransferList(offset, limit int, chainId, contractAddr, userAddr string) (
//	[]*db.NonFungibleTransfer, error) {
//	transfers := make([]*db.NonFungibleTransfer, 0)
//	if chainId == "" || contractAddr == "" {
//		return transfers, db.ErrTableParams
//	}
//
//	whereOr := map[string]interface{}{}
//	where := map[string]interface{}{
//		"contractAddr": contractAddr,
//	}
//	if userAddr != "" {
//		whereOr["fromAddr"] = userAddr
//		whereOr["toAddr"] = userAddr
//	}
//	selectFile := &SelectFile{
//		Where:   where,
//		WhereOr: whereOr,
//	}
//
//	tableName := db.GetTableName(chainId, db.TableNonFungibleTransfer)
//	query := BuildParamsQuery(tableName, selectFile)
//	query = query.Order("timestamp desc").Offset(offset * limit).Limit(limit)
//	err := query.Find(&transfers).Error
//	if err != nil {
//		return transfers, fmt.Errorf("GetNFTTransferList err, cause : %s", err.Error())
//	}
//
//	return transfers, nil
//}

// UpdateNonTransferContractName 更新合约名称
func UpdateNonTransferContractName(chainId, contractName, contractAddr string) error {
	if chainId == "" || contractName == "" || contractAddr == "" {
		return nil
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	params := map[string]interface{}{
		"contractName": contractName,
	}
	tableName := db.GetTableName(chainId, db.TableNonFungibleTransfer)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// InsertNonFungibleTransfer 批量保存
func InsertNonFungibleTransfer(chainId string, transfers []*db.NonFungibleTransfer) error {
	if len(transfers) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableNonFungibleTransfer)
	return CreateInBatchesData(tableName, transfers)
}

// GetNFTTransferCount
//
//	@Description: 统计非同质化流转数据
//	@param chainId
//	@param tokenId
//	@param contractAddr
//	@param userAddr
//	@return int64
//	@return error
func GetNFTTransferCount(chainId, contractAddr, userAddr, tokenId string) (int64, error) {
	var count int64
	tableName := db.GetTableName(chainId, db.TableNonFungibleTransfer)
	if userAddr == "" {
		where := map[string]interface{}{
			"contractAddr": contractAddr,
		}
		if tokenId != "" {
			where["tokenId"] = tokenId
		}
		// 当userAddr为空时，执行简单查询
		query := db.GormDB.Table(tableName).Where(where)
		err := query.Count(&count).Error
		if err != nil {
			return 0, err
		}
	} else {
		subQuerySql := GetTransferSubQuerySql(tableName, tableName)
		// 当userAddr不为空时，执行UNION查询
		sqlQuery := fmt.Sprintf(`
			SELECT COUNT(*) FROM (%s) AS t;
		`, subQuerySql)

		// 执行原始SQL查询并获取结果数量
		query := db.GormDB.Raw(sqlQuery, contractAddr, userAddr, contractAddr, userAddr)
		err := query.Row().Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("GetNFTTransferCount err, cause : %s", err.Error())
		}
	}

	return count, nil
}

func GetNFTTransferList(offset, limit int, chainId, contractAddr, userAddr, tokenId string) (
	[]*db.NonFungibleTransfer, error) {
	transfers := make([]*db.NonFungibleTransfer, 0)
	if chainId == "" || contractAddr == "" {
		return transfers, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableNonFungibleTransfer)
	if userAddr == "" {
		where := map[string]interface{}{
			"contractAddr": contractAddr,
		}
		if tokenId != "" {
			where["tokenId"] = tokenId
		}
		// 当userAddr为空时，执行简单查询
		query := db.GormDB.Table(tableName).Where(where)
		err := query.Order("timestamp desc").Offset(offset * limit).Limit(limit).Find(&transfers).Error
		if err != nil {
			return transfers, fmt.Errorf("GetNFTTransferList err, cause : %s", err.Error())
		}
	} else {
		subQuerySql := GetTransferSubQuerySql(tableName, tableName)
		// 当userAddr不为空时，执行UNION查询
		sqlQuery := fmt.Sprintf(`
			SELECT * FROM (%s) AS t
			ORDER BY timestamp DESC
			LIMIT ? OFFSET ?;
		`, subQuerySql)

		// 执行原始SQL查询并将结果映射到transfers切片
		query := db.GormDB.Raw(sqlQuery, contractAddr, userAddr, contractAddr, userAddr, limit, offset*limit)
		err := query.Scan(&transfers).Error
		if err != nil {
			return transfers, fmt.Errorf("GetNFTTransferList err, cause : %s", err.Error())
		}
	}

	return transfers, nil
}

func GetTransferSubQuerySql(tableName1, tableName2 string) string {
	sqlQuery := fmt.Sprintf(`
			(
  				SELECT * FROM %s WHERE contractAddr = ? AND fromAddr = ?
			)
			UNION
			(
 				 SELECT * FROM %s WHERE contractAddr = ? AND toAddr = ?
			)
		`, tableName1, tableName2)

	return sqlQuery
}
