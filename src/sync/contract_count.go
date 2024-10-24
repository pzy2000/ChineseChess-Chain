/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"

	"github.com/shopspring/decimal"
)

// DealFungibleContractUpdateData 计算同质化合约持有人数
func DealFungibleContractUpdateData(holdCountMap map[string]int64, totalSupplyMap map[string]decimal.Decimal,
	ftContractMap map[string]*db.FungibleContract, minHeight int64) map[string]*db.FungibleContract {
	updateFTContractMap := make(map[string]*db.FungibleContract, 0)
	if len(ftContractMap) == 0 {
		return updateFTContractMap
	}

	if len(holdCountMap) == 0 && len(totalSupplyMap) == 0 {
		return updateFTContractMap
	}

	for _, contract := range ftContractMap {
		if contract.BlockHeight >= minHeight {
			//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
			continue
		}

		count, okCount := holdCountMap[contract.ContractAddr]
		total, okTotal := totalSupplyMap[contract.ContractAddr]
		//hold和total都没有改，就不需要更新
		if !okCount && !okTotal {
			continue
		}

		if okCount {
			contract.HolderCount = contract.HolderCount + count
		}
		if okTotal {
			// 将字符串转换为 decimal.Decimal 值
			//totalSupplyDB, _ := decimal.NewFromString(contract.TotalSupply)
			// 使用 Add 方法将两个 Decimal 值相加
			//totalRes := contract.TotalSupply.Add(total)
			contract.TotalSupply = contract.TotalSupply.Add(total)
		}
		//contract.BlockHeight = minHeight
		updateFTContractMap[contract.ContractAddr] = contract
	}

	return updateFTContractMap
	//更新update数据
	//delayedUpdateData.ContractResult.UpdateFungibleContract = updateFungible
}

// DealNonFungibleContractUpdateData 计算同质化合约持有人数
func DealNonFungibleContractUpdateData(holdCountMap map[string]int64, totalSupplyMap map[string]decimal.Decimal,
	nftContractMap map[string]*db.NonFungibleContract, minHeight int64) map[string]*db.NonFungibleContract {
	updateFungibleMap := make(map[string]*db.NonFungibleContract, 0)
	if len(nftContractMap) == 0 {
		return updateFungibleMap
	}

	if len(holdCountMap) == 0 && len(totalSupplyMap) == 0 {
		return updateFungibleMap
	}

	for _, contract := range nftContractMap {
		if contract.BlockHeight >= minHeight {
			//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
			continue
		}

		count, okCount := holdCountMap[contract.ContractAddr]
		total, okTotal := totalSupplyMap[contract.ContractAddr]
		//hold和total都没有改，就不需要更新
		if !okCount && !okTotal {
			continue
		}

		if okCount {
			contract.HolderCount = contract.HolderCount + count
		}

		if okTotal {
			// 将字符串转换为 decimal.Decimal 值
			//totalSupplyDB, _ := decimal.NewFromString(contract.TotalSupply)
			// 使用 Add 方法将两个 Decimal 值相加
			//totalRes := contract.TotalSupply.Add(total)
			contract.TotalSupply = contract.TotalSupply.Add(total)
		}
		//contract.BlockHeight = minHeight
		updateFungibleMap[contract.ContractAddr] = contract
	}

	return updateFungibleMap
	//更新update数据
	//delayedUpdateData.ContractResult.UpdateNonFungible = updateFungible
}

// DealContractTotalSupply 计算发行总量
func DealContractTotalSupply(contractEvents []*db.ContractEventData,
	contractMap map[string]*db.Contract) map[string]decimal.Decimal {
	//存储合约发行量
	contractAdders := make(map[string]decimal.Decimal)
	for _, event := range contractEvents {
		//只解析交易流转的topic
		if _, ok := TopicEventDataKey[event.Topic]; !ok {
			continue
		}
		contract, isOk := contractMap[event.ContractName]
		if !isOk || event.EventData == nil {
			continue
		}

		if event.EventData.Amount == "" && event.EventData.TokenId == "" {
			continue
		}

		fromAddress := event.EventData.FromAddress
		toAddress := event.EventData.ToAddress
		amountDecimal := decimal.NewFromInt(1)
		if contract.ContractType == ContractStandardNameCMDFA ||
			contract.ContractType == ContractStandardNameEVMDFA {
			//同质化合约
			amountDecimal = StringAmountDecimal(event.EventData.Amount, contract.Decimals)
		}

		if fromAddress == "" {
			//没有from地址就是增发,累加
			if decimalVal, ok := contractAdders[contract.Addr]; ok {
				// 使用 Add 方法将两个 Decimal 值相加
				amountRes := decimalVal.Add(amountDecimal)
				contractAdders[contract.Addr] = amountRes
			} else {
				contractAdders[contract.Addr] = amountDecimal
			}
		} else if toAddress == "" {
			//没有to地址就只销毁，减法
			if decimalVal, ok := contractAdders[contract.Addr]; ok {
				// 使用 Sub 方法将两个 Decimal 值相减
				amountRes := decimalVal.Sub(amountDecimal)
				contractAdders[contract.Addr] = amountRes
			} else {
				dec := decimal.NewFromInt(0)
				amountRes := dec.Sub(amountDecimal)
				contractAdders[contract.Addr] = amountRes
			}
		} else {
			//from-》to转帐，TotalSupply不变
			continue
		}
	}

	return contractAdders
}

// DealContractHoldCount 计算同质化合约持有人数
func DealContractHoldCount(positionOperates *db.BlockPosition) map[string]int64 {
	holdCountList := make(map[string]int64, 0)
	//同质化
	insertPosition := positionOperates.InsertFungiblePosition
	deletePosition := positionOperates.DeleteFungiblePosition
	//新增持有人
	if len(insertPosition) > 0 {
		for _, position := range insertPosition {
			holdCountList[position.ContractAddr]++
		}
	}

	if len(deletePosition) > 0 {
		for _, position := range deletePosition {
			holdCountList[position.ContractAddr]--
		}
	}

	//非同质化
	insertNonPosition := positionOperates.InsertNonFungible
	deleteNonPosition := positionOperates.DeleteNonFungible
	if len(insertNonPosition) > 0 {
		for _, position := range insertNonPosition {
			holdCountList[position.ContractAddr]++
		}
	}

	if len(deleteNonPosition) > 0 {
		for _, position := range deleteNonPosition {
			holdCountList[position.ContractAddr]--
		}
	}
	return holdCountList
}

// NFTContractTransferNum
//
//	@Description: 统计NFT合约交易流转数量
//	@param transferList
//	@param nftContractMap
//	@param minHeight
//	@return map[string]*db.NonFungibleContract
func NFTContractTransferNum(transferList []*db.NonFungibleTransfer, nftContractMap map[string]*db.NonFungibleContract,
	minHeight int64) map[string]*db.NonFungibleContract {
	updateNFTContractMap := make(map[string]*db.NonFungibleContract, 0)
	if len(transferList) == 0 || len(nftContractMap) == 0 {
		return updateNFTContractMap
	}

	// Create a map to store the count of each ContractAddr in the transferList
	contractAddrCount := make(map[string]int)

	// Loop through the transferList and count the occurrences of each ContractAddr
	for _, transfer := range transferList {
		contractAddrCount[transfer.ContractAddr]++
	}

	// Loop through the nftContractMap and update the TransferNum
	for contractAddr, contract := range nftContractMap {
		if contract.BlockHeight >= minHeight {
			//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
			continue
		}

		if count, exists := contractAddrCount[contractAddr]; exists {
			// Update the TransferNum
			contract.TransferNum += int64(count)

			// Add the updated contract to the updateNFTContractMap
			updateNFTContractMap[contractAddr] = contract
		}
	}

	return updateNFTContractMap
}

func MergeNFTContractMaps(minHeight int64, updateNFTContractMap,
	nftContractTransferMap map[string]*db.NonFungibleContract) []*db.NonFungibleContract {
	// Loop through the nftContractTransfetMap
	for contractAddr, transferContract := range nftContractTransferMap {
		// Get the corresponding contract from updateNFTContractMap
		if updateContract, exists := updateNFTContractMap[contractAddr]; exists {
			// Update the TransferNum
			updateContract.TransferNum = transferContract.TransferNum
		} else {
			// If the contract doesn't exist in updateNFTContractMap, add it
			updateNFTContractMap[contractAddr] = transferContract
		}
	}

	// Convert the merged map to a slice
	mergedNFTContractSlice := make([]*db.NonFungibleContract, 0, len(updateNFTContractMap))
	for _, contract := range updateNFTContractMap {
		contract.BlockHeight = minHeight
		mergedNFTContractSlice = append(mergedNFTContractSlice, contract)
	}

	return mergedNFTContractSlice
}

func FTContractTransferNum(transferList []*db.FungibleTransfer, ftContractMap map[string]*db.FungibleContract,
	minHeight int64) map[string]*db.FungibleContract {
	updateFTContractMap := make(map[string]*db.FungibleContract, 0)
	if len(transferList) == 0 || len(ftContractMap) == 0 {
		return updateFTContractMap
	}

	// Create a map to store the count of each ContractAddr in the transferList
	contractAddrCount := make(map[string]int)

	// Loop through the transferList and count the occurrences of each ContractAddr
	for _, transfer := range transferList {
		contractAddrCount[transfer.ContractAddr]++
	}

	// Loop through the nftContractMap and update the TransferNum
	for contractAddr, contract := range ftContractMap {
		if contract.BlockHeight >= minHeight {
			//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
			continue
		}

		if count, exists := contractAddrCount[contractAddr]; exists {
			// Update the TransferNum
			contract.TransferNum += int64(count)
			// Add the updated contract to the updateNFTContractMap
			updateFTContractMap[contractAddr] = contract
		}
	}

	return updateFTContractMap
}

func MergeFTContractMaps(minHeight int64, updateFTContractMap,
	ftContractTransferMap map[string]*db.FungibleContract) []*db.FungibleContract {
	// Loop through the nftContractTransfetMap
	for contractAddr, transferContract := range ftContractTransferMap {
		// Get the corresponding contract from updateNFTContractMap
		if updateContract, exists := updateFTContractMap[contractAddr]; exists {
			// Update the TransferNum
			updateContract.TransferNum = transferContract.TransferNum
		} else {
			// If the contract doesn't exist in updateNFTContractMap, add it
			updateFTContractMap[contractAddr] = transferContract
		}
	}

	// Convert the merged map to a slice
	mergedFTContractSlice := make([]*db.FungibleContract, 0, len(updateFTContractMap))
	for _, contract := range updateFTContractMap {
		contract.BlockHeight = minHeight
		mergedFTContractSlice = append(mergedFTContractSlice, contract)
	}

	return mergedFTContractSlice
}
