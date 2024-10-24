package dbhandle

import "chainmaker_web/src/db"

// InsertCrossTxTransfers 插入跨链交易流转
func InsertCrossTxTransfers(chainId string, crossTxTransfers []*db.CrossTransactionTransfer) error {
	if len(crossTxTransfers) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	err := CreateInBatchesData(tableName, crossTxTransfers)
	if err != nil {
		return err
	}
	return nil
}

// CheckCrossIdsExistenceTransfer 查询数据是否已经保存
func CheckCrossIdsExistenceTransfer(chainId string, crossIds []string) (map[string]bool, error) {
	txTransferMap := make(map[string]bool, 0)
	if chainId == "" || len(crossIds) == 0 {
		return txTransferMap, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	// 查询与 crossIds 匹配的唯一 CrossId
	var uniqueCrossIds []string
	err := db.GormDB.Table(tableName).
		Select("DISTINCT crossId").
		Where("crossId IN ?", crossIds).
		Pluck("crossId", &uniqueCrossIds).Error

	if err != nil {
		return txTransferMap, err
	}

	// 将查询结果保存到 map 中
	for _, crossId := range uniqueCrossIds {
		txTransferMap[crossId] = true
	}

	return txTransferMap, nil
}

// GetCrossCycleTransferById 根据Id获取交易流转
func GetCrossCycleTransferById(chainId, crossId string) ([]*db.CrossTransactionTransfer, error) {
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	cycleTransfers := make([]*db.CrossTransactionTransfer, 0)
	where := map[string]interface{}{
		"crossId": crossId,
	}
	err := db.GormDB.Table(tableName).Where(where).Find(&cycleTransfers).Error
	if err != nil {
		return cycleTransfers, err
	}

	return cycleTransfers, nil
}

func GetCrossCycleTransferByCrossIds(chainId string, crossIds []string) ([]*db.CrossTransactionTransfer, error) {
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	cycleTransfers := make([]*db.CrossTransactionTransfer, 0)
	err := db.GormDB.Table(tableName).Where("crossId in ?", crossIds).Find(&cycleTransfers).Error
	if err != nil {
		return cycleTransfers, err
	}

	return cycleTransfers, nil
}

// GetCrossCycleTransferByHeight 根据height获取交易流转
func GetCrossCycleTransferByHeight(chainId string, blockHeight []int64) ([]*db.CrossTransactionTransfer, error) {
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	cycleTransfers := make([]*db.CrossTransactionTransfer, 0)
	err := db.GormDB.Table(tableName).Where("blockHeight IN ?", blockHeight).Find(&cycleTransfers).Error
	if err != nil {
		return cycleTransfers, err
	}

	return cycleTransfers, nil
}
