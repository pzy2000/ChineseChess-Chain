/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
)

// GetAccountType 获取账户类型
func GetAccountType(chainId, address string) int {
	//判断是否是合约地址
	contractInfo, err := dbhandle.GetContractByCacheOrAddr(chainId, address)
	addrType := AddrTypeUser
	if contractInfo != nil && err == nil {
		addrType = AddrTypeContract
	}
	return addrType
}

// BuildAccountInsertOrUpdate
//
//	@Description: 计算需要新增，更新的账户信息
//	@param chainId
//	@param topicEventResult 处理好的event事件数据
//	@return []*db.Account 新增账户列表
//	@return []*db.Account 更新账户DID，BNS
//	@return map[string]*db.Account 所有的涉及到的账户数据
//	@return error
func BuildAccountInsertOrUpdate(chainId string, minHeight int64, delayGetDBResult *GetDBResult,
	topicEventResult *TopicEventResult, accountTx, accountNFT map[string]int64) ([]*db.Account, []*db.Account,
	map[string]*db.Account, error) {
	var (
		accountInsertMap = make(map[string]*db.Account)
		accountUpdateMap = make(map[string]*db.Account)
		accountInsert    []*db.Account
		accountUpdate    []*db.Account
		accountMap       = delayGetDBResult.AccountDBMap
	)

	//数据库不存在的账户即为新增数据
	for _, addr := range topicEventResult.OwnerAdders {
		if addr == "" {
			continue
		}
		//数据库不存在的账户地址
		if _, ok := accountMap[addr]; !ok {
			//获取账户类型
			addrType := GetAccountType(chainId, addr)
			accountInfo := &db.Account{
				AddrType:    addrType,
				Address:     addr,
				BlockHeight: minHeight,
			}
			accountInsertMap[addr] = accountInfo
		}
	}

	//处理BNS账户
	bnsBind := topicEventResult.BNSBindEventData
	bnsUnBindAccount := delayGetDBResult.AccountBNSList
	processBNSAccounts(bnsBind, bnsUnBindAccount, accountInsertMap, accountUpdateMap, accountMap)
	//处理DID账户
	didBind := topicEventResult.DIDAccount
	didUnBindAccount := delayGetDBResult.AccountDIDList
	processDIDAccounts(didBind, didUnBindAccount, accountInsertMap, accountUpdateMap, accountMap)

	//统计账户交易数量
	dealAccountTxNum(chainId, minHeight, accountTx, accountInsertMap, accountUpdateMap, accountMap)
	//统计账户NFT数量
	dealAccountNFTNum(chainId, minHeight, accountNFT, accountInsertMap, accountUpdateMap, accountMap)

	if len(accountInsertMap) > 0 {
		for _, account := range accountInsertMap {
			accountInsert = append(accountInsert, account)
			accountMap[account.Address] = account
		}
	}
	if len(accountUpdateMap) > 0 {
		for _, account := range accountUpdateMap {
			//已经更新过了
			if account.BlockHeight >= minHeight {
				continue
			}
			account.BlockHeight = minHeight
			accountUpdate = append(accountUpdate, account)
			accountMap[account.Address] = account
		}
	}
	return accountInsert, accountUpdate, accountMap, nil
}

/**
 * @description: 统计账户NFT数量
 * @param {string} chainId 链id
 * @param {int64} minHeight  本次订阅区块高度
 * @param {*} accountInsertMap 新增账户
 * @param {*} accountUpdateMap 更新账户
 * @param {map[string]*db.Account} accountMap 本次订阅涉及账户
 * @param {map[string]int64} accountNFT 账户NFT数量
 */
func dealAccountNFTNum(chainId string, minHeight int64, accountNFT map[string]int64,
	accountInsertMap, accountUpdateMap, accountMap map[string]*db.Account) {
	for address, num := range accountNFT {
		if address == "" {
			continue
		}

		if account, ok := accountInsertMap[address]; ok {
			nftNum := account.NFTNum + num
			if nftNum < 0 {
				nftNum = 0
			}
			account.NFTNum = nftNum
		} else if accountUp, okUp := accountUpdateMap[address]; okUp {
			nftNum := accountUp.NFTNum + num
			if nftNum < 0 {
				nftNum = 0
			}
			accountUp.NFTNum = nftNum
		} else if accountDB, okDB := accountMap[address]; okDB {
			nftNum := accountDB.NFTNum + num
			if nftNum < 0 {
				nftNum = 0
			}
			accountDB.NFTNum = nftNum
			accountUpdateMap[address] = accountDB
		} else {
			var nftNum int64
			if num > 0 {
				nftNum = num
			}
			//获取账户类型
			addrType := GetAccountType(chainId, address)
			accountInfo := &db.Account{
				AddrType:    addrType,
				Address:     address,
				NFTNum:      nftNum,
				BlockHeight: minHeight,
			}
			accountInsertMap[address] = accountInfo
		}
	}
}

/**
 * @description: 统计账户交易数量
 * @param {string} chainId
 * @param {int64} minHeight
 * @param {*} accountInsertMap  新增账户
 * @param {*} accountUpdateMap  更新账户
 * @param {map[string]*db.Account} accountMap 本次订阅涉及账户
 * @param {map[string]int64} accountTx 本次订阅账户交易量
 */
func dealAccountTxNum(chainId string, minHeight int64, accountTx map[string]int64,
	accountInsertMap, accountUpdateMap, accountMap map[string]*db.Account) {
	for address, num := range accountTx {
		if address == "" || num == 0 {
			continue
		}

		if account, ok := accountInsertMap[address]; ok {
			txNum := account.TxNum + num
			if txNum < 0 {
				txNum = 0
			}
			account.TxNum = txNum
		} else if accountUp, okUp := accountUpdateMap[address]; okUp {
			txNum := accountUp.TxNum + num
			if txNum < 0 {
				txNum = 0
			}
			accountUp.TxNum = txNum
		} else if accountDB, okDB := accountMap[address]; okDB {
			txNum := accountDB.TxNum + num
			if txNum < 0 {
				txNum = 0
			}
			accountDB.TxNum = txNum
			accountUpdateMap[address] = accountDB
		} else {
			var txNum int64
			if num > 0 {
				txNum = num
			}
			//获取账户类型
			addrType := GetAccountType(chainId, address)
			accountInfo := &db.Account{
				AddrType:    addrType,
				Address:     address,
				TxNum:       txNum,
				BlockHeight: minHeight,
			}
			accountInsertMap[address] = accountInfo
		}
	}
}

// processBNSAccounts
//
//	@Description: 更新账户BNS
//	@param chainId
//	@param bnsBindEventData 绑定bns
//	@param accountInsertMap
//	@param accountUpdateMap
//	@param accountMap
func processBNSAccounts(bnsBindEventData []*db.BNSTopicEventData, unBindBNSs []*db.Account, accountInsertMap,
	accountUpdateMap, accountMap map[string]*db.Account) {
	//绑定BNS
	for _, event := range bnsBindEventData {
		accountAddr := event.Value
		accountBNS := event.Domain
		if accountDB, ok := accountMap[accountAddr]; ok {
			// 数据库存在，更新BNS
			if updateAcc, okUp := accountUpdateMap[accountAddr]; okUp {
				updateAcc.BNS = accountBNS
				accountUpdateMap[accountAddr] = updateAcc
			} else {
				accountDB.BNS = accountBNS
				accountUpdateMap[accountAddr] = accountDB
			}
		} else {
			// 数据库不存在，判断是否已经在insert了
			if account, okIn := accountInsertMap[accountAddr]; okIn {
				account.BNS = accountBNS
			}
		}
	}

	//解绑BNS
	for _, account := range unBindBNSs {
		accountAddr := account.Address
		if accountDB, ok := accountMap[accountAddr]; ok {
			// 数据库存在，更新BNS
			if updateAcc, okUp := accountUpdateMap[accountAddr]; okUp {
				updateAcc.BNS = ""
				accountUpdateMap[accountAddr] = updateAcc
			} else {
				accountDB.BNS = ""
				accountUpdateMap[accountAddr] = accountDB
			}
		}
	}

}

// 处理DID账户
func processDIDAccounts(didAccount map[string][]string, unBindDIDs []*db.Account, accountInsertMap, accountUpdateMap,
	accountMap map[string]*db.Account) {
	// 将 unBindDIDs 转换为一个包含地址的集合，用于解绑DID
	unbindDIDSet := make(map[string]bool)
	for _, account := range unBindDIDs {
		unbindDIDSet[account.Address] = true
	}

	for did, didAccounts := range didAccount {
		for _, accountAddr := range didAccounts {
			if accountDB, ok := accountMap[accountAddr]; ok {
				// 数据库存在，更新DID
				if updateAcc, okUp := accountUpdateMap[accountAddr]; okUp {
					updateAcc.DID = did
					accountUpdateMap[accountAddr] = updateAcc
				} else {
					accountDB.DID = did
					accountUpdateMap[accountAddr] = accountDB
				}

				// 从 unbindDIDSet 中移除已处理的地址
				delete(unbindDIDSet, accountAddr)
			} else {
				// 数据库不存在的账户地址
				if account, okIn := accountInsertMap[accountAddr]; okIn {
					account.DID = did
				}
			}
		}
	}

	// 处理 unbindDIDSet 中剩余的地址，这些数据将要解绑DID
	for accountAddr := range unbindDIDSet {
		if accountDB, ok := accountMap[accountAddr]; ok {
			// 数据库存在，更新DID
			if updateAcc, okUp := accountUpdateMap[accountAddr]; okUp {
				updateAcc.DID = ""
				accountUpdateMap[accountAddr] = updateAcc
			} else {
				accountDB.DID = ""
				accountUpdateMap[accountAddr] = accountDB
			}
		}
	}
}

func DealAccountTxNFTNum(txList map[string]*db.Transaction, contractEvents []*db.ContractEventData) (
	map[string]int64, map[string]int64) {
	accountTxNum := make(map[string]int64)
	accountNFTNum := make(map[string]int64)
	// 更新交易数量
	for _, tx := range txList {
		accountTxNum[tx.UserAddr]++
	}

	for _, event := range contractEvents {
		// 只解析交易流转的topic
		if _, ok := TopicEventDataKey[event.Topic]; !ok {
			continue
		}

		if event.EventData == nil || event.EventData.TokenId == "" {
			continue
		}

		fromAddr := event.EventData.FromAddress
		toAddr := event.EventData.ToAddress
		if fromAddr != "" {
			accountNFTNum[fromAddr]--
		}
		if toAddr != "" {
			accountNFTNum[toAddr]++
		}
	}

	return accountTxNum, accountNFTNum
}
