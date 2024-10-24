package dbhandle

import (
	"chainmaker_web/src/db"
)

// InsertCrossSubTransaction 插入跨链交易
func InsertCrossSubTransaction(chainId string, crossTxs []*db.CrossMainTransaction) error {
	if len(crossTxs) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossMainTransaction)
	return CreateInBatchesData(tableName, crossTxs)
}

// InsertCrossBusinessTransaction 插入跨链交易
func InsertCrossBusinessTransaction(chainId string, crossTxs []*db.CrossBusinessTransaction) error {
	if len(crossTxs) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossBusinessTransaction)
	return CreateInBatchesData(tableName, crossTxs)
}

// GetCrossBusinessTxById 根据Id获取业务交易流转
func GetCrossBusinessTxById(chainId string, txList []string) ([]*db.CrossBusinessTransaction, error) {
	businessTxs := make([]*db.CrossBusinessTransaction, 0)
	if len(txList) == 0 {
		return businessTxs, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossBusinessTransaction)
	err := db.GormDB.Table(tableName).Where("txId IN ?", txList).Find(&businessTxs).Error
	if err != nil {
		return businessTxs, err
	}

	return businessTxs, nil
}

// GetCrossBusinessTxByCross 根据Id获取业务交易流转
func GetCrossBusinessTxByCross(chainId, crossId string) ([]*db.CrossBusinessTransaction, error) {
	businessTxs := make([]*db.CrossBusinessTransaction, 0)
	if chainId == "" || crossId == "" {
		return businessTxs, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossBusinessTransaction)
	where := map[string]interface{}{
		"crossId": crossId,
	}
	err := db.GormDB.Table(tableName).Where(where).Find(&businessTxs).Error
	if err != nil {
		return businessTxs, err
	}

	return businessTxs, nil
}
