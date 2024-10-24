/*
Package db comment
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package db

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

const (
	CHIndexTypeRange = 1
	CHIndexTypeIndex = 2
)

var (
	// ErrRecordNotFoundErr 没有查到
	ErrRecordNotFoundErr = errors.New("record not found")
	// ErrTableParams 参数错误
	ErrTableParams = errors.New("table params err")
)

var (
	TableBlock                      = `block`
	TableTransaction                = `transaction`
	TableBlackTransaction           = `black_transaction`
	TableContract                   = `contract`
	TableContractEvent              = `contract_event`
	TableUser                       = "user"
	TableOrg                        = "org"
	TableNode                       = "node"
	TableContractUpgradeTransaction = "contract_upgrade_transaction"
	TableGas                        = "gas"
	TableGasRecord                  = "gas_record"
	TableChain                      = `chain`
	TableSubscribe                  = `subscribe`
	TableFungibleContract           = `fungible_contract`
	TableNonFungibleContract        = `non_fungible_contract`
	TableEvidenceContract           = `evidence_contract`
	TableEvidenceMetaData           = `evidence_contract_metadata`
	TableIdentityContract           = `identity_contract`
	TableFungibleTransfer           = `fungible_transfer`
	TableNonFungibleTransfer        = `non_fungible_transfer`
	TableFungiblePosition           = `fungible_position`
	TableNonFungiblePosition        = `non_fungible_position`
	// #nosec G101
	TableNonFungibleToken         = `non_fungible_Token`
	TableAccount                  = `account`
	TableCrossSubChainData        = `cross_sub_chain_data`
	TableCrossMainTransaction     = `cross_main_transaction`
	TableCrossTransactionTransfer = `cross_transaction_transfer`
	TableCrossCycleTransaction    = `cross_cycle_transaction`
	TableCrossBusinessTransaction = `cross_business_transaction`
	TableCrossSubChainCrossChain  = `cross_sub_chain_cross_chain`
	TableCrossChainContract       = `cross_chain_contract`

	TableIDAContract        = `ida_contract`
	TableIDAAssetDetail     = `ida_asset_detail`
	TableIDAAssetAttachment = `ida_asset_attachment`
	TableIDADataAsset       = `ida_data_asset`
	TableIDAApiAsset        = `ida_api_asset`
)

// TableInfo 表数据
type TableInfo struct {
	Name        string      // 表名称
	Structure   interface{} // 表结构
	Description string      // 表描述
}

// InitTableName 初始化表名
func InitTableName(prefix string) {
	TableBlock = prefix + TableBlock
	TableTransaction = prefix + TableTransaction
	TableBlackTransaction = prefix + TableBlackTransaction
	TableContract = prefix + TableContract
	TableContractEvent = prefix + TableContractEvent
	TableUser = prefix + TableUser
	TableOrg = prefix + TableOrg
	TableNode = prefix + TableNode
	TableContractUpgradeTransaction = prefix + TableContractUpgradeTransaction
	TableGas = prefix + TableGas
	TableGasRecord = prefix + TableGasRecord
	TableChain = prefix + TableChain
	TableSubscribe = prefix + TableSubscribe

	TableFungibleContract = prefix + TableFungibleContract
	TableNonFungibleContract = prefix + TableNonFungibleContract
	TableEvidenceContract = prefix + TableEvidenceContract
	TableEvidenceMetaData = prefix + TableEvidenceMetaData
	TableIdentityContract = prefix + TableIdentityContract
	TableFungibleTransfer = prefix + TableFungibleTransfer
	TableNonFungibleTransfer = prefix + TableNonFungibleTransfer
	TableFungiblePosition = prefix + TableFungiblePosition
	TableNonFungiblePosition = prefix + TableNonFungiblePosition
	TableNonFungibleToken = prefix + TableNonFungibleToken
	TableAccount = prefix + TableAccount

	TableCrossSubChainData = prefix + TableCrossSubChainData
	TableCrossMainTransaction = prefix + TableCrossMainTransaction
	TableCrossTransactionTransfer = prefix + TableCrossTransactionTransfer
	TableCrossCycleTransaction = prefix + TableCrossCycleTransaction
	TableCrossBusinessTransaction = prefix + TableCrossBusinessTransaction
	TableCrossSubChainCrossChain = prefix + TableCrossSubChainCrossChain
	TableCrossChainContract = prefix + TableCrossChainContract

	TableIDAContract = prefix + TableIDAContract
	TableIDAAssetDetail = prefix + TableIDAAssetDetail
	TableIDAAssetAttachment = prefix + TableIDAAssetAttachment
	TableIDADataAsset = prefix + TableIDADataAsset
	TableIDAApiAsset = prefix + TableIDAApiAsset
}

// GetTableName GetTableName
//
// GetTableName
//
//	@Description: 获取表名称
//	@param chainId 链ID
//	@param table 表名称
//	@return string 表名称
func GetTableName(chainId, table string) string {
	if table == TableSubscribe || table == TableChain || chainId == "" {
		return table
	}
	tableName := fmt.Sprintf("%s_%s", chainId, table)
	// 转换为小写
	return strings.ToLower(tableName)
}

// GetBlockTableNames 获取所有表名
func GetBlockTableNames() []TableInfo {
	return []TableInfo{
		{
			Name:        TableBlock,
			Structure:   &Block{},
			Description: "区块信息表",
		},
		{
			Name:        TableTransaction,
			Structure:   &Transaction{},
			Description: "交易信息表",
		},
		{
			Name:        TableBlackTransaction,
			Structure:   &BlackTransaction{},
			Description: "交易黑名单表",
		},
		{
			Name:        TableContract,
			Structure:   &Contract{},
			Description: "合约基础信息表",
		},
		{
			Name:        TableContractEvent,
			Structure:   &ContractEvent{},
			Description: "合约交易事件表",
		},
		{
			Name:        TableUser,
			Structure:   &User{},
			Description: "用户表",
		},
		{
			Name:        TableOrg,
			Structure:   &Org{},
			Description: "组织信息表",
		},
		{
			Name:        TableNode,
			Structure:   &Node{},
			Description: "节点信息表",
		},
		{
			Name:        TableContractUpgradeTransaction,
			Structure:   &UpgradeContractTransaction{},
			Description: "合约更新记录表",
		},
		{
			Name:        TableGas,
			Structure:   &Gas{},
			Description: "gas余额统计表",
		},
		{
			Name:        TableGasRecord,
			Structure:   &GasRecord{},
			Description: "gas消耗记录表",
		},
		{
			Name:        TableFungibleContract,
			Structure:   &FungibleContract{},
			Description: "同质化合约信息表",
		},
		{
			Name:        TableNonFungibleContract,
			Structure:   &NonFungibleContract{},
			Description: "非同质化合约信息表",
		},
		{
			Name:        TableEvidenceContract,
			Structure:   &EvidenceContract{},
			Description: "存证合约信息表",
		},
		{
			Name:        TableIdentityContract,
			Structure:   &IdentityContract{},
			Description: "身份认证合约信息表",
		},
		{
			Name:        TableFungibleTransfer,
			Structure:   &FungibleTransfer{},
			Description: "同质化合约交易流转信息表",
		},
		{
			Name:        TableNonFungibleTransfer,
			Structure:   &NonFungibleTransfer{},
			Description: "非同质化合约交易流转信息表",
		},
		{
			Name:        TableFungiblePosition,
			Structure:   &FungiblePosition{},
			Description: "同质化合约账户持仓统计表",
		},
		{
			Name:        TableNonFungiblePosition,
			Structure:   &NonFungiblePosition{},
			Description: "非同质化合约账户持仓统计表",
		},
		{
			Name:        TableNonFungibleToken,
			Structure:   &NonFungibleToken{},
			Description: "非同质化合约token数据表",
		},
		{
			Name:        TableAccount,
			Structure:   &Account{},
			Description: "账户信息表",
		},
		{
			Name:        TableTransaction,
			Structure:   &Transaction{},
			Description: "交易信息表",
		},
		{
			Name:        TableCrossSubChainData,
			Structure:   &CrossSubChainData{},
			Description: "跨链子链详情表",
		},
		{
			Name:        TableCrossMainTransaction,
			Structure:   &CrossMainTransaction{},
			Description: "跨链主链交易表",
		},
		{
			Name:        TableCrossTransactionTransfer,
			Structure:   &CrossTransactionTransfer{},
			Description: "跨链交易流转信息表",
		},

		{
			Name:        TableCrossCycleTransaction,
			Structure:   &CrossCycleTransaction{},
			Description: "跨链交易状态表",
		},
		{
			Name:        TableCrossBusinessTransaction,
			Structure:   &CrossBusinessTransaction{},
			Description: "跨链执行的业务交易表",
		},
		{
			Name:        TableCrossSubChainCrossChain,
			Structure:   &CrossSubChainCrossChain{},
			Description: "跨链子链跨链数据统计",
		},
		{
			Name:        TableCrossChainContract,
			Structure:   &CrossChainContract{},
			Description: "跨链合约表",
		},
		{
			Name:        TableIDAContract,
			Structure:   &IDAContract{},
			Description: "IDA合约信息表",
		},
		{
			Name:        TableIDAAssetDetail,
			Structure:   &IDAAssetDetail{},
			Description: "IDA资产信息表",
		},
		{
			Name:        TableTransaction,
			Structure:   &Transaction{},
			Description: "交易信息表",
		},
		{
			Name:        TableIDAAssetAttachment,
			Structure:   &IDAAssetAttachment{},
			Description: "IDA资产附件信息表",
		},
		{
			Name:        TableIDADataAsset,
			Structure:   &IDADataAsset{},
			Description: "IDA中data类型资产数据表",
		},
		{
			Name:        TableTransaction,
			Structure:   &Transaction{},
			Description: "交易信息表",
		},
		{
			Name:        TableIDAApiAsset,
			Structure:   &IDAApiAsset{},
			Description: "IDA中API类型资产数据表",
		},
	}
}

// GetClickhouseTableOptions
//
//	@Description: clickhouse获取表设计
//	@return map[string]string
//
//nolint:govet,lll
func GetClickhouseTableOptions() map[string]string {
	uniqueOptions := map[string]string{
		TableBlock:                      "ENGINE=ReplacingMergeTree() ORDER BY (blockHeight) PRIMARY KEY (blockHeight)",
		TableTransaction:                "ENGINE=ReplacingMergeTree() ORDER BY (timestamp, txId, contractAddr) PRIMARY KEY (timestamp, txId)",
		TableBlackTransaction:           "ENGINE=ReplacingMergeTree() ORDER BY txId  PRIMARY KEY (txId)",
		TableContract:                   "ENGINE=ReplacingMergeTree() ORDER BY (timestamp, addr)  PRIMARY KEY (timestamp, addr)",
		TableContractEvent:              "ENGINE=ReplacingMergeTree() ORDER BY (contractNameBak, timestamp, txId, eventIndex) PRIMARY KEY (contractNameBak, timestamp, txId, eventIndex)",
		TableUser:                       "ENGINE=ReplacingMergeTree() ORDER BY (userAddr) PRIMARY KEY (userAddr)",
		TableOrg:                        "ENGINE=ReplacingMergeTree() ORDER BY (orgId) PRIMARY KEY (orgId)",
		TableNode:                       "ENGINE=ReplacingMergeTree() ORDER BY (nodeId) PRIMARY KEY (nodeId)",
		TableContractUpgradeTransaction: "ENGINE=ReplacingMergeTree() ORDER BY (timestamp, txId) PRIMARY KEY (timestamp, txId)",
		TableGas:                        "ENGINE=ReplacingMergeTree() ORDER BY (address) PRIMARY KEY (address)",
		TableGasRecord:                  "ENGINE=ReplacingMergeTree() ORDER BY (timestamp, id) PRIMARY KEY (timestamp, id)",
		TableFungibleContract:           "ENGINE=ReplacingMergeTree() ORDER BY (timestamp, contractAddr) PRIMARY KEY (timestamp, contractAddr)",
		TableNonFungibleContract:        "ENGINE=ReplacingMergeTree() ORDER BY (timestamp, contractAddr) PRIMARY KEY (timestamp, contractAddr)",
		TableEvidenceContract:           "ENGINE=ReplacingMergeTree() ORDER BY (timestamp, id) PRIMARY KEY (timestamp, id)",
		TableIdentityContract:           "ENGINE=ReplacingMergeTree() ORDER BY (id) PRIMARY KEY (id)",
		TableFungibleTransfer:           "ENGINE=ReplacingMergeTree() ORDER BY (contractAddr, timestamp, id) PRIMARY KEY (contractAddr, timestamp, id)",
		TableNonFungibleTransfer:        "ENGINE=ReplacingMergeTree() ORDER BY (contractAddr, timestamp, id) PRIMARY KEY (contractAddr, timestamp, id)",
		TableFungiblePosition:           "ENGINE=ReplacingMergeTree() ORDER BY (contractAddr, id) PRIMARY KEY (contractAddr, id)",
		TableNonFungiblePosition:        "ENGINE=ReplacingMergeTree() ORDER BY (contractAddr, id) PRIMARY KEY (contractAddr, id)",
		TableNonFungibleToken:           "ENGINE=ReplacingMergeTree() ORDER BY (contractAddr, id) PRIMARY KEY (contractAddr, id)",
		TableAccount:                    "ENGINE=ReplacingMergeTree() ORDER BY (address) PRIMARY KEY (address)",
		TableCrossSubChainData:          "ENGINE=ReplacingMergeTree() ORDER BY (subChainId) PRIMARY KEY (subChainId)",
		TableCrossChainContract:         "ENGINE=ReplacingMergeTree() ORDER BY (subChainId, contractName)",
		TableCrossSubChainCrossChain:    "ENGINE=ReplacingMergeTree() ORDER BY (subChainId, id) PRIMARY KEY (subChainId, id)",
		TableCrossMainTransaction:       "ENGINE=ReplacingMergeTree() ORDER BY (timestamp, txId) PRIMARY KEY (timestamp, txId)",
		TableCrossTransactionTransfer:   "ENGINE=ReplacingMergeTree() ORDER BY (crossId, id) PRIMARY KEY (crossId, id)",
		TableCrossCycleTransaction:      "ENGINE=ReplacingMergeTree() ORDER BY (crossId, id) PRIMARY KEY (crossId, id)",
		TableCrossBusinessTransaction:   "ENGINE=ReplacingMergeTree() ORDER BY (subChainId, crossId, id) PRIMARY KEY (subChainId, crossId, id)",
	}

	return uniqueOptions
}

type CHIndexInfo struct {
	IndexType int
	Fields    []string
}

//nolint:govet,lll
func GetClickhouseDBIndexFields(table string) []CHIndexInfo {
	indexList := make([]CHIndexInfo, 0)
	switch table {
	case TableBlock:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"blockHeight"}},
		}
	case TableTransaction:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeRange, Fields: []string{"timestamp"}},
			{IndexType: CHIndexTypeIndex, Fields: []string{"sender", "blockHeight", "userAddr"}},
		}
	case TableContractEvent:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeRange, Fields: []string{"timestamp"}},
			{IndexType: CHIndexTypeIndex, Fields: []string{"contractNameBak"}},
		}
	case TableFungibleContract:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"contractNameBak", "contractAddr"}},
		}
	case TableNonFungibleContract:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"contractNameBak", "contractAddr"}},
		}
	case TableEvidenceContract:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeRange, Fields: []string{"timestamp"}},
			{IndexType: CHIndexTypeIndex, Fields: []string{"hash", "contractName"}},
		}
	case TableIdentityContract:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"contractName", "contractAddr"}},
		}
	case TableFungibleTransfer:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeRange, Fields: []string{"timestamp"}},
			{IndexType: CHIndexTypeIndex, Fields: []string{"fromAddr", "toAddr", "contractName"}},
		}
	case TableNonFungibleTransfer:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeRange, Fields: []string{"timestamp"}},
			{IndexType: CHIndexTypeIndex, Fields: []string{"fromAddr", "toAddr", "tokenId"}},
		}
	case TableFungiblePosition:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"ownerAddr"}},
		}
	case TableNonFungiblePosition:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"ownerAddr"}},
		}
	case TableNonFungibleToken:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"ownerAddr"}},
		}
	case TableAccount:
		indexList = []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"did", "bns"}},
		}
	}
	return indexList
}

func SqlCreateTableWithComment(db *gorm.DB, chainId string, tableInfo TableInfo) error {
	tableName := GetTableName(chainId, tableInfo.Name)
	if GormDB.Migrator().HasTable(tableName) {
		//表已经存在了，直接跳过
		return nil
	}

	//创建表
	err := db.Table(tableName).AutoMigrate(tableInfo.Structure)
	if err != nil {
		return err
	}

	// 添加表注释
	err = db.Exec(fmt.Sprintf("ALTER TABLE %s COMMENT '%s'", tableName, tableInfo.Description)).Error
	if err != nil {
		return nil
	}

	log.Infof("SqlCreateTableWithComment create table success, tableName:%v", tableName)
	return err
}

func ClickHouseCreateTableWithComment(chainId string, tableInfo TableInfo) error {
	var err error
	table := tableInfo.Name
	tableModel := tableInfo.Structure
	tableName := GetTableName(chainId, table)
	uniqueOptionsMap := GetClickhouseTableOptions()
	if GormDB.Migrator().HasTable(tableName) {
		//表已经存在了，直接跳过
		return nil
	}

	//创建表结构
	if uniqueOptionsMap[table] != "" {
		err = GormDB.Set("gorm:table_options", uniqueOptionsMap[table]).
			Table(tableName).AutoMigrate(tableModel)
	} else {
		err = GormDB.Table(tableName).AutoMigrate(tableModel)
	}
	if err != nil {
		return err
	}

	//创建索引
	indexFields := GetClickhouseDBIndexFields(table)
	for _, indexInfo := range indexFields {
		indexName := tableName + strings.Join(indexInfo.Fields, "_") + "_idx"
		switch indexInfo.IndexType {
		case CHIndexTypeIndex:
			// 创建数据跳过索引
			indexFieldsName := strings.Join(indexInfo.Fields, ", ")
			createIndexSQL := fmt.Sprintf("CREATE INDEX %s ON %s (%s) TYPE set(0) GRANULARITY 2;",
				indexName, tableName, indexFieldsName)
			err = GormDB.Exec(createIndexSQL).Error
			if err != nil {
				return err
			}
		case CHIndexTypeRange:
			// 创建范围查询的数据跳过索引
			indexFieldsName := strings.Join(indexInfo.Fields, ", ")
			createIndexSQL := fmt.Sprintf("CREATE INDEX %s ON %s (%s) TYPE minmax GRANULARITY 2;",
				indexName, tableName, indexFieldsName)
			err = GormDB.Exec(createIndexSQL).Error
			if err != nil {
				return err
			}
		}
	}

	//修改表字段类型
	switch table {
	case TableTransaction:
		_ = GormDB.Exec(fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN contractResult %s",
			tableName, "Array(UInt8)")).Error
		err = GormDB.Exec(fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN contractResultBak %s",
			tableName, "Array(UInt8)")).Error
	case TableCrossBusinessTransaction:
		err = GormDB.Exec(fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN contractResult %s",
			tableName, "Array(UInt8)")).Error
	}
	if err != nil {
		return err
	}
	return nil
}

// type SqlIndexInfo struct {
// 	isUnique bool
// 	fields   []string
// }

// // GetSqlDBIndexFields
// //
// //	@Description: 返回一个表对应的组合索引字段列表
// //	@param table 不包含chainId的表名称
// //	@return []*SqlIndexInfo 索引列表，只包括联合索引，单独索引去gorm表结构定义
// //
// //nolint:govet
// func GetSqlDBIndexFields(table string) []*SqlIndexInfo {
// 	indexList := make([]*SqlIndexInfo, 0)
// 	switch table {
// 	case TableBlock:
// 		//没有联合索引
// 	case TableContractUpgradeTransaction:
// 		//没有联合索引
// 	case TableTransaction:
// 		indexList = []*SqlIndexInfo{
// 			{isUnique: false, fields: []string{"sender", "timestamp"}},
// 			{isUnique: false, fields: []string{"blockHeight", "timestamp"}},
// 			{isUnique: false, fields: []string{"contractAddr", "timestamp"}},
// 			{isUnique: false, fields: []string{"userAddr", "timestamp"}},
// 		}
// 	case TableContractEvent:
// 		indexList = []*SqlIndexInfo{
// 			{isUnique: true, fields: []string{"txId", "eventIndex"}},
// 			{isUnique: false, fields: []string{"contractNameBak", "timestamp"}},
// 		}
// 	case TableGasRecord:
// 		indexList = []*SqlIndexInfo{
// 			{isUnique: true, fields: []string{"txId", "gasIndex"}},
// 			{isUnique: false, fields: []string{"businessType", "timestamp"}},
// 		}
// 	case TableIdentityContract:
// 		indexList = []*SqlIndexInfo{
// 			{isUnique: true, fields: []string{"txId", "eventIndex"}},
// 		}
// 	case TableFungibleTransfer:
// 		indexList = []*SqlIndexInfo{
// 			{isUnique: true, fields: []string{"txId", "eventIndex"}},
// 			{isUnique: false, fields: []string{"contractAddr", "timestamp"}},
// 			{isUnique: false, fields: []string{"contractAddr", "fromAddr"}},
// 			{isUnique: false, fields: []string{"contractAddr", "toAddr"}},
// 		}
// 	case TableNonFungibleTransfer:
// 		indexList = []*SqlIndexInfo{
// 			{isUnique: true, fields: []string{"txId", "eventIndex"}},
// 			{isUnique: false, fields: []string{"contractAddr", "timestamp"}},
// 			{isUnique: false, fields: []string{"contractAddr", "fromAddr"}},
// 			{isUnique: false, fields: []string{"contractAddr", "toAddr"}},
// 		}
// 	case TableFungiblePosition:
// 		indexList = []*SqlIndexInfo{
// 			{isUnique: true, fields: []string{"contractAddr", "ownerAddr"}},
// 			{isUnique: false, fields: []string{"contractAddr", "amount"}},
// 		}
// 	case TableNonFungiblePosition:
// 		indexList = []*SqlIndexInfo{
// 			{isUnique: true, fields: []string{"contractAddr", "ownerAddr"}},
// 			{isUnique: false, fields: []string{"contractAddr", "amount"}},
// 		}
// 	case TableNonFungibleToken:
// 		indexList = []*SqlIndexInfo{
// 			{isUnique: true, fields: []string{"tokenId", "contractAddr"}},
// 			{isUnique: false, fields: []string{"ownerAddr", "timestamp"}},
// 			{isUnique: false, fields: []string{"contractAddr", "timestamp"}},
// 			{isUnique: false, fields: []string{"ownerAddr", "contractAddr", "timestamp"}},
// 		}
// 	case TableCrossChainContract:
// 		indexList = []*SqlIndexInfo{
// 			{isUnique: true, fields: []string{"subChainId", "contractName"}},
// 		}
// 	case TableCrossCycleTransaction:
// 		indexList = []*SqlIndexInfo{
// 			{isUnique: false, fields: []string{"startTime", "endTime"}},
// 		}
// 	}

// 	return indexList
// }
