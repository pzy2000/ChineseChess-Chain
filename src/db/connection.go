/*
Package db comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package db

import (
	"chainmaker_web/src/config"
	loggers "chainmaker_web/src/logger"
	"database/sql"
	"fmt"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// SqlDB DB client
	SqlDB *sql.DB
	//GormDB DB client
	GormDB *gorm.DB
	log    = loggers.GetLogger(loggers.MODULE_WEB)
)

// InitDbConn init database connection
func InitDbConn(dbConfig *config.DBConf) {
	var err error
	// 创建 MySQL 和 ClickHouse 数据库连接
	GormDB, err = ConnectDatabase(dbConfig, true)
	if err != nil {
		//创建数据库
		GormDB, err = ConnectDatabase(dbConfig, false)
		if err != nil {
			log.Errorf("failed to connect database: %v", err)
			panic(err)
		}
		CreateDatabase(GormDB, dbConfig.Database, dbConfig.DbProvider)
		//重新连接数据库
		GormDB, err = ConnectDatabase(dbConfig, true)
		if err != nil {
			panic(err)
		}
	}

	//初始化表结构
	InitTableName(dbConfig.Prefix)
	log.Infof("======3333333333333======:%v", config.SubscribeChains[0])
	//初始化数据库
	InitDBTable(dbConfig, config.SubscribeChains)
}

func InitDBTable(dbConfig *config.DBConf, chainList []*config.ChainInfo) {
	//初始化表
	switch dbConfig.DbProvider {
	case config.MySql:
		InitMysqlTable(chainList)
	case config.ClickHouse:
		InitClickHouseTable(chainList)
	case config.Pgsql:
		InitPgsqlTable(chainList)
	}
}

// ConnectDatabase 连接数据库
func ConnectDatabase(dbConfig *config.DBConf, useDataBase bool) (*gorm.DB, error) {
	var err error
	switch dbConfig.DbProvider {
	case config.MySql:
		dsn := dbConfig.ToMysqlUrl(useDataBase)
		GormDB, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,
			DontSupportRenameColumn:   true,  // rename column not supported before clickhouse 20.4
			SkipInitializeWithVersion: false, // smart configure based on used version
		}), &gorm.Config{
			//Logger: logger.Default.LogMode(logger.Info),
		})
	case config.ClickHouse:
		dsn := dbConfig.ToClickHouseUrl(useDataBase)
		GormDB, err = gorm.Open(clickhouse.New(clickhouse.Config{
			DSN:                          dsn,
			DontSupportRenameColumn:      true,  // rename column not supported before clickhouse 20.4
			DontSupportEmptyDefaultValue: false, // do not consider empty strings as valid default values
			SkipInitializeWithVersion:    false, // smart configure based on used version
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	case config.Pgsql:
		dsn := dbConfig.ToPgsqlUrl(useDataBase)
		GormDB, err = gorm.Open(postgres.Open(dsn),
			&gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
	}

	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB, _ := GormDB.DB()
	sqlDB.SetMaxIdleConns(config.DbMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.DbMaxOpenConns)
	return GormDB, nil
}

// CreateDatabase 创建数据库
func CreateDatabase(db *gorm.DB, database, dbProvider string) {
	var createDatabaseQuery string
	// 创建数据库
	switch dbProvider {
	case config.MySql:
		createDatabaseQuery = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s %s",
			database, config.MysqlDatabaseConf)
	case config.ClickHouse:
		createDatabaseQuery = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", database)
	case config.Pgsql:
		createDatabaseQuery = fmt.Sprintf("CREATE DATABASE %s", database)
	default:
		return
	}

	err := db.Exec(createDatabaseQuery).Error
	log.Infof("CREATE DATABASE %v", database)
	if err != nil {
		log.Errorf("CREATE DATABASE failed, err:%v", err)
	}
}

// InitMysqlTable 初始化数据库1
func InitMysqlTable(chainList []*config.ChainInfo) {
	err := GormDB.AutoMigrate(
		&Chain{},
		&Subscribe{},
	)
	if err != nil {
		panic(err)
	}

	//其他表按链ID分表
	blockTableNames := GetBlockTableNames()
	for _, chainInfo := range chainList {
		for _, tableInfo := range blockTableNames {
			err = SqlCreateTableWithComment(GormDB, chainInfo.ChainId, tableInfo)
			if err != nil {
				panic(err)
			}
		}
	}
}

func InitPgsqlTable(chainList []*config.ChainInfo) {
	err := GormDB.AutoMigrate(
		&Chain{},
		&Subscribe{},
	)
	if err != nil {
		panic(err)
	}

	//其他表按链ID分表
	blockTableNames := GetBlockTableNames()
	for _, chainInfo := range chainList {
		for _, tableInfo := range blockTableNames {
			err = SqlCreateTableWithComment(GormDB, chainInfo.ChainId, tableInfo)
			if err != nil {
				panic(err)
			}
			//if GormDB.Migrator().HasTable(tableName) {
			//	//表已经存在了，直接跳过
			//	continue
			//}
			//
			//err = GormDB.Table(GetTableName(chainInfo.ChainId, tableName)).AutoMigrate(tableModel)
			//if err != nil {
			//	panic(err)
			//}
		}
	}
}

// InitClickHouseTable 初始化数据库
func InitClickHouseTable(chainList []*config.ChainInfo) {
	err := GormDB.AutoMigrate(
		&Chain{},
		&Subscribe{},
	)
	if err != nil {
		panic(err)
	}

	//其他表按链ID分表
	blockTableNames := GetBlockTableNames()
	for _, chainInfo := range chainList {
		for _, tableInfo := range blockTableNames {
			err = ClickHouseCreateTableWithComment(chainInfo.ChainId, tableInfo)
			if err != nil {
				panic(err)
			}
		}
	}
}

// DeleteDbTable 删除表
// func DeleteDbTable(chainId string) {
// 	var err error
// 	tables := []string{
// 		TableBlock, TableTransaction, TableBlackTransaction, TableContract, TableContractEvent, TableUser, TableOrg,
// 		TableNode, TableContractUpgradeTransaction, TableGas, TableGasRecord,
// 		TableFungibleContract, TableNonFungibleContract, TableEvidenceContract, TableEvidenceMetaData,
// 		TableIdentityContract, TableFungibleTransfer, TableNonFungibleTransfer,
// 		TableFungiblePosition, TableNonFungiblePosition, TableNonFungibleToken, TableAccount,
// 	}

// 	for _, table := range tables {
// 		fullTableName := GetTableName(chainId, table)
// 		if SqlDB != nil {
// 			_, err = SqlDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", fullTableName))
// 		} else if GormDB != nil {
// 			err = GormDB.Migrator().DropTable(fullTableName)
// 		}
// 		if err != nil {
// 			log.Errorf("deleteTables err:%v, table:%v", err, fullTableName)
// 		}
// 	}
// }

// DeleteTablesByChainID 根据链ID删除相关表
func DeleteTablesByChainID(chainId string) error {
	blockTableNames := GetBlockTableNames()
	for _, tableInfo := range blockTableNames {
		tableName := GetTableName(chainId, tableInfo.Name)
		err := GormDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)).Error
		if err != nil {
			return err
		}
	}

	return nil
}
