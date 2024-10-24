package saveTasks

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"

	"github.com/google/uuid"
)

// GetSubChainSaveList 获取子链插入更新数据
func GetSubChainSaveList(chainId string, saveSubChainList []*db.CrossSubChainData) (
	[]*db.CrossSubChainData, []*db.CrossSubChainData, error) {
	insertSubChainList := make([]*db.CrossSubChainData, 0)
	updateSubChainList := make([]*db.CrossSubChainData, 0)
	subChainIds := make([]string, 0)
	for _, subChainData := range saveSubChainList {
		subChainIds = append(subChainIds, subChainData.SubChainId)
	}

	crossSubChainDBMap, err := dbhandle.GetCrossSubChainById(chainId, subChainIds)
	if err != nil {
		return insertSubChainList, updateSubChainList, err
	}

	for _, subChainData := range saveSubChainList {
		if _, exists := crossSubChainDBMap[subChainData.SubChainId]; exists {
			updateSubChainList = append(updateSubChainList, subChainData)
		} else {
			insertSubChainList = append(insertSubChainList, subChainData)
		}
	}

	return insertSubChainList, updateSubChainList, nil
}

// SaveRelayCrossChainToDB 存储主子链数据
func SaveRelayCrossChainToDB(chainId string, crossChainResult *db.CrossChainResult) error {
	if crossChainResult == nil {
		return nil
	}

	var err error
	//插入子链数据
	//获取数据库子链数据
	insertSubChainList, updateSubChainList, err := GetSubChainSaveList(chainId, crossChainResult.SaveSubChainList)
	if err != nil {
		log.Errorf("SaveRelayCrossChainToDB GetSubChainSaveList failed, SaveSubChainList:%v",
			crossChainResult.SaveSubChainList)
		return err
	}

	err = dbhandle.InsertCrossSubChain(chainId, insertSubChainList)
	if err != nil {
		log.Errorf("SaveRelayCrossChainToDB insert sub chain failed, SubChain:%v",
			insertSubChainList)
		return err
	}

	//更新子链网关数据
	for _, subChain := range updateSubChainList {
		err = dbhandle.UpdateCrossSubChainById(chainId, subChain)
		if err != nil {
			log.Errorf("SaveRelayCrossChainToDB Update sub chain failed subChain:%v ", subChain)
			return err
		}
	}

	//存储跨链交易
	err = dbhandle.InsertCrossSubTransaction(chainId, crossChainResult.CrossMainTransaction)
	if err != nil {
		log.Errorf("SaveRelayCrossChainToDB insert tx failed, tx:%v",
			crossChainResult.CrossMainTransaction)
		return err
	}

	//跨链合约
	err = SaveCrossChainContract(chainId, crossChainResult.CrossChainContractMap)
	if err != nil {
		return err
	}

	//存储跨链交易流转数据
	err = SaveCrossTransfer(chainId, crossChainResult.CrossTransfer)
	if err != nil {
		return err
	}

	//保存跨链周期数据
	insertCycleTx := crossChainResult.InsertCrossCycleTx
	updateCycleTx := crossChainResult.UpdateCrossCycleTx
	err = SaveCrossCycleTx(chainId, insertCycleTx, updateCycleTx)
	if err != nil {
		return err
	}

	//业务交易数据
	err = SaveBusinessTransaction(chainId, crossChainResult.BusinessTxMap)
	if err != nil {
		return err
	}

	//更新子链区块高度
	for spvContractName, blockHeight := range crossChainResult.SubChainBlockHeight {
		subChain := &db.CrossSubChainData{
			SpvContractName: spvContractName,
			BlockHeight:     blockHeight,
		}
		err = dbhandle.UpdateCrossChainHeightBySpv(chainId, subChain)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveCrossChainContract 保存跨链合约数据
func SaveCrossChainContract(chainId string, chainContractMap map[string]map[string]string) error {
	if len(chainContractMap) == 0 {
		return nil
	}

	for subChainId, contractMap := range chainContractMap {
		var contractNames []string
		crossContractDBMap := make(map[string]*db.CrossChainContract, 0)
		insertCrossContracts := make([]*db.CrossChainContract, 0)
		for _, name := range contractMap {
			contractNames = append(contractNames, name)
		}
		// 查询数据库中是否存在记录
		crossContractDB, err := dbhandle.GetCrossContractByName(chainId, subChainId, contractNames)
		if err != nil {
			log.Errorf("Failed to get cross contract by name: %v", err)
			return err
		}
		for _, contract := range crossContractDB {
			crossContractDBMap[contract.ContractName] = contract
		}

		for _, name := range contractMap {
			if _, ok := crossContractDBMap[name]; !ok {
				newUUID := uuid.New().String()
				insertCrossContracts = append(insertCrossContracts, &db.CrossChainContract{
					ID:           newUUID,
					SubChainId:   subChainId,
					ContractName: name,
				})
			}
		}
		// 如果记录不存在，则插入新记录
		err = dbhandle.InsertCrossContract(chainId, subChainId, insertCrossContracts)
		if err != nil {
			log.Errorf("Failed to insert cross contract: %v", err)
			return err
		}
	}

	return nil
}

// SaveCrossTransfer 报错跨链周期流转数据
func SaveCrossTransfer(chainId string, crossTransfer map[string]*db.CrossTransactionTransfer) error {
	if len(crossTransfer) == 0 {
		return nil
	}
	crossIds := make([]string, 0)
	for _, transfer := range crossTransfer {
		crossIds = append(crossIds, transfer.CrossId)
	}
	crossIdsExistenceMap, err := dbhandle.CheckCrossIdsExistenceTransfer(chainId, crossIds)
	if err != nil {
		return err
	}
	insertTransfer := make([]*db.CrossTransactionTransfer, 0)
	for _, transfer := range crossTransfer {
		if !crossIdsExistenceMap[transfer.CrossId] {
			insertTransfer = append(insertTransfer, transfer)
		}
	}

	err = dbhandle.InsertCrossTxTransfers(chainId, insertTransfer)
	if err != nil {
		log.Errorf("SaveRelayCrossChainToDB insert transfer failed, transfer:%v",
			insertTransfer)
		return err
	}

	return nil
}

// SaveCrossCycleTx 保存跨链周期数据
func SaveCrossCycleTx(chainId string, insertCycleTxs []*db.CrossCycleTransaction,
	updateCycleTxs map[string]*db.CrossCycleTransaction) error {
	if len(insertCycleTxs) == 0 && len(updateCycleTxs) == 0 {
		return nil
	}

	for _, cycleTx := range updateCycleTxs {
		//更新
		err := dbhandle.UpdateCrossCycleTx(chainId, cycleTx)
		if err != nil {
			log.Errorf("SaveCrossCycleTx failed UpdateCrossCycleTx:%v ", cycleTx)
			return err
		}
	}

	if len(insertCycleTxs) > 0 {
		err := dbhandle.InsertCrossCycleTx(chainId, insertCycleTxs)
		if err != nil {
			log.Errorf("SaveCrossCycleTx failed, insertCycleTxs:%v", insertCycleTxs)
			return err
		}
	}

	return nil
}

// SaveBusinessTransaction 保存主子链业务交易数据
func SaveBusinessTransaction(chainId string, businessTxMap map[string]*db.CrossBusinessTransaction) error {
	if len(businessTxMap) == 0 {
		return nil
	}
	insertTxList := make([]*db.CrossBusinessTransaction, 0)
	for _, txInfo := range businessTxMap {
		insertTxList = append(insertTxList, txInfo)
	}
	err := dbhandle.InsertCrossBusinessTransaction(chainId, insertTxList)
	if err != nil {
		log.Errorf("SaveBusinessTransaction failed, txlist:%v", insertTxList)
		return err
	}
	return nil
}

// SaveCrossSubChainCrossToDB 保存子链跨链数据
func SaveCrossSubChainCrossToDB(chainId string, inserts []*db.CrossSubChainCrossChain,
	updates []*db.CrossSubChainCrossChain) error {
	if len(inserts) < 0 && len(updates) == 0 {
		return nil
	}

	err := dbhandle.InsertCrossSubChainCross(chainId, inserts)
	if err != nil {
		log.Errorf("SaveCrossSubChainCrossToDB InsertCrossSubChainCross failed,data:%v ", inserts)
		return err
	}
	for _, subChainCross := range updates {
		err = dbhandle.UpdateCrossSubChainCross(chainId, subChainCross)
		if err != nil {
			log.Errorf("SaveCrossSubChainCrossToDB UpdateCrossSubChainCross failed,data:%v ", subChainCross)
			return err
		}
	}

	return nil
}

// UpdateCrossSubChainData 保存子链跨链数据
func UpdateCrossSubChainData(chainId string, updates []*db.CrossSubChainData) error {
	if len(updates) == 0 {
		return nil
	}

	for _, subChainData := range updates {
		err := dbhandle.UpdateCrossSubChainById(chainId, subChainData)
		if err != nil {
			log.Errorf("SaveCrossSubChainCrossToDB UpdateCrossSubChainCross failed,data:%v ", subChainData)
			return err
		}
	}

	return nil
}
