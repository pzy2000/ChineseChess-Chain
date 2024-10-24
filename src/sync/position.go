/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"time"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"
)

type EventTokenInfo struct {
	Event *db.ContractEventData
	//token最终归属
	FinalOwner string
	FirstOwner string
	AddrType   int
}

// BuildPositionList 构造持仓数据
func BuildPositionList(contractEvents []*db.ContractEventData, contractInfoMap map[string]*db.Contract,
	accountMap map[string]*db.Account) map[string]*db.PositionData {
	var positionList = make(map[string]*db.PositionData)
	for _, event := range contractEvents {
		//只解析交易流转的topic
		if _, ok := TopicEventDataKey[event.Topic]; !ok {
			continue
		}

		if event.EventData == nil {
			continue
		}
		amount := event.EventData.Amount
		tokenId := event.EventData.TokenId
		if amount == "" && tokenId == "" {
			continue
		}

		//合约信息
		contract, ok := contractInfoMap[event.ContractName]
		if !ok {
			continue
		}
		//非标准合约不解析
		if contract.ContractType != ContractStandardNameCMDFA &&
			contract.ContractType != ContractStandardNameEVMDFA &&
			contract.ContractType != ContractStandardNameCMNFA &&
			contract.ContractType != ContractStandardNameEVMNFA {
			continue
		}

		amountDecimal := decimal.NewFromInt(1)
		fromAddress := event.EventData.FromAddress
		toAddress := event.EventData.ToAddress
		//同质化
		if contract.ContractType == ContractStandardNameCMDFA ||
			contract.ContractType == ContractStandardNameEVMDFA {
			amountDecimal = StringAmountDecimal(amount, contract.Decimals)
		}

		//构造持仓数据
		PositionFromAddressData(positionList, fromAddress, amountDecimal, accountMap, contract)
		PositionToAddressData(positionList, toAddress, amountDecimal, accountMap, contract)
	}

	return positionList
}

// PositionFromAddressData 根据from地址计算position数据
func PositionFromAddressData(positionList map[string]*db.PositionData, address string, amount decimal.Decimal,
	accountMap map[string]*db.Account, contract *db.Contract) {
	if address == "" {
		return
	}
	positionKey := address + "_" + contract.Addr
	position, ok := positionList[positionKey]
	if ok {
		// 使用 Sub 方法将两个 Decimal 值相减
		amountRes := position.Amount.Sub(amount)
		position.Amount = amountRes
	} else {
		dec := decimal.NewFromInt(0)
		amountRes := dec.Sub(amount)
		//获取地址类型
		addrType := AddrTypeUser
		if value, isOk := accountMap[address]; isOk {
			addrType = value.AddrType
		}
		positionList[positionKey] = &db.PositionData{
			AddrType:     addrType,
			OwnerAddr:    address,
			ContractAddr: contract.Addr,
			ContractName: contract.Name,
			Symbol:       contract.ContractSymbol,
			Amount:       amountRes,
			Decimals:     contract.Decimals,
			ContractType: contract.ContractType,
		}
	}
}

// PositionToAddressData 根据to地址计算position数据
func PositionToAddressData(positionList map[string]*db.PositionData, address string, amount decimal.Decimal,
	accountMap map[string]*db.Account, contract *db.Contract) {
	if address == "" {
		return
	}

	positionKey := address + "_" + contract.Addr
	position, ok := positionList[positionKey]
	if ok {
		// 使用 Add 方法将两个 Decimal 值相加
		amountRes := position.Amount.Add(amount)
		position.Amount = amountRes
	} else {
		//获取地址类型
		addrType := AddrTypeUser
		if value, isOk := accountMap[address]; isOk {
			addrType = value.AddrType
		}
		positionList[positionKey] = &db.PositionData{
			AddrType:     addrType,
			OwnerAddr:    address,
			ContractAddr: contract.Addr,
			ContractName: contract.Name,
			Symbol:       contract.ContractSymbol,
			Amount:       amount,
			Decimals:     contract.Decimals,
			ContractType: contract.ContractType,
		}
	}
}

// BuildUpdatePositionData 解析所有持仓数据
func BuildUpdatePositionData(minHeight int64, positionList map[string]*db.PositionData,
	positionDBMap map[string][]*db.FungiblePosition,
	nonPositionDBMap map[string][]*db.NonFungiblePosition) *db.BlockPosition {
	positionOp := &db.BlockPosition{}
	for _, position := range positionList {
		if position.ContractType == ContractStandardNameCMDFA ||
			position.ContractType == ContractStandardNameEVMDFA {
			//同质化
			dealFungiblePosition(minHeight, position, positionDBMap, positionOp)
		} else if position.ContractType == ContractStandardNameCMNFA ||
			position.ContractType == ContractStandardNameEVMNFA {
			//非同质化
			dealNonFungiblePosition(minHeight, position, nonPositionDBMap, positionOp)
		}
	}

	return positionOp
}

// dealFungiblePosition 处理同质化合约数据
func dealFungiblePosition(minHeight int64, position *db.PositionData, positionDBMap map[string][]*db.FungiblePosition,
	positionOp *db.BlockPosition) {
	var (
		isHave         bool
		insertPosition = make([]*db.FungiblePosition, 0)
		updatePosition = make([]*db.FungiblePosition, 0)
		deletePosition = make([]*db.FungiblePosition, 0)
	)
	//判断数据库是否存在
	positionDbValue := &db.FungiblePosition{}
	if dbList, ok := positionDBMap[position.OwnerAddr]; ok {
		for _, positionDb := range dbList {
			if position.ContractAddr == positionDb.ContractAddr {
				isHave = true
				positionDbValue = positionDb
				break
			}
		}
	}
	//存在就更新，删除
	if isHave {
		//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
		if positionDbValue.BlockHeight >= minHeight {
			return
		}

		// 将字符串转换为 Decimal 值
		//amountDecimalDB, err := decimal.NewFromString(positionDbValue.Amount)
		//if err != nil {
		//	log.Error("dealFungiblePosition", "amountDecimalDB", err)
		//	return
		//}
		// 使用 Add 方法将两个 Decimal 值相加
		amountRes := positionDbValue.Amount.Add(position.Amount)
		//判断是否大于0
		if amountRes.GreaterThan(decimal.Zero) {
			positionDbValue.Amount = amountRes
			positionDbValue.BlockHeight = minHeight
			updatePosition = append(updatePosition, positionDbValue)
		} else {
			deletePosition = append(deletePosition, positionDbValue)
		}
	} else {
		//不存在就添加
		newUUID := uuid.New().String()
		tempPosition := &db.FungiblePosition{
			ID: newUUID,
			//AddrType:     position.AddrType,
			OwnerAddr:    position.OwnerAddr,
			ContractAddr: position.ContractAddr,
			ContractName: position.ContractName,
			Symbol:       position.Symbol,
			Amount:       position.Amount,
			BlockHeight:  minHeight,
		}
		insertPosition = append(insertPosition, tempPosition)
	}
	if len(insertPosition) > 0 {
		positionOp.InsertFungiblePosition = append(positionOp.InsertFungiblePosition, insertPosition...)
	}
	if len(updatePosition) > 0 {
		positionOp.UpdateFungiblePosition = append(positionOp.UpdateFungiblePosition, updatePosition...)
	}
	if len(deletePosition) > 0 {
		positionOp.DeleteFungiblePosition = append(positionOp.DeleteFungiblePosition, deletePosition...)
	}
}

// dealNonFungiblePosition 处理同质化合约数据
func dealNonFungiblePosition(minHeight int64, position *db.PositionData,
	positionDBMap map[string][]*db.NonFungiblePosition, positionOp *db.BlockPosition) {
	var (
		isHave         bool
		insertPosition = make([]*db.NonFungiblePosition, 0)
		updatePosition = make([]*db.NonFungiblePosition, 0)
		deletePosition = make([]*db.NonFungiblePosition, 0)
	)
	//判断数据库是否存在
	positionDbValue := &db.NonFungiblePosition{}
	if dbList, ok := positionDBMap[position.OwnerAddr]; ok {
		for _, positionDb := range dbList {
			if position.ContractAddr == positionDb.ContractAddr {
				isHave = true
				positionDbValue = positionDb
				break
			}
		}
	}
	//存在就更新，删除
	if isHave {
		//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
		if positionDbValue.BlockHeight >= minHeight {
			return
		}

		// 将字符串转换为 Decimal 值
		//amountDecimalDB, err := decimal.NewFromString(positionDbValue.Amount)
		//if err != nil {
		//	log.Error("dealFungiblePosition", "amountDecimalDB", err)
		//	return
		//}
		// 使用 Add 方法将两个 Decimal 值相加
		amountRes := positionDbValue.Amount.Add(position.Amount)
		//判断是否大于0
		if amountRes.GreaterThan(decimal.Zero) {
			positionDbValue.Amount = amountRes
			positionDbValue.BlockHeight = minHeight
			updatePosition = append(updatePosition, positionDbValue)
		} else {
			deletePosition = append(deletePosition, positionDbValue)
		}
	} else {
		//不存在就添加
		newUUID := uuid.New().String()
		tempPosition := &db.NonFungiblePosition{
			ID: newUUID,
			//	AddrType:     position.AddrType,
			OwnerAddr:    position.OwnerAddr,
			ContractAddr: position.ContractAddr,
			ContractName: position.ContractName,
			Amount:       position.Amount,
			BlockHeight:  minHeight,
		}
		insertPosition = append(insertPosition, tempPosition)
	}
	if len(insertPosition) > 0 {
		positionOp.InsertNonFungible = append(positionOp.InsertNonFungible, insertPosition...)
	}
	if len(updatePosition) > 0 {
		positionOp.UpdateNonFungible = append(positionOp.UpdateNonFungible, updatePosition...)
	}
	if len(deletePosition) > 0 {
		positionOp.DeleteNonFungible = append(positionOp.DeleteNonFungible, deletePosition...)
	}
}

// DealNonFungibleToken 处理Token
func DealNonFungibleToken(chainId string, contractEvents []*db.ContractEventData, contractMap map[string]*db.Contract,
	accountMap map[string]*db.Account) *db.TokenResult {
	tokenResult := &db.TokenResult{
		InsertUpdateToken: make([]*db.NonFungibleToken, 0),
		DeleteToken:       make([]*db.NonFungibleToken, 0),
	}

	// 用于存储每个tokenId的所有事件，tokenId+ContractAddr确定唯一token
	tokenEvents := make(map[string][]*db.ContractEventData)
	for _, event := range contractEvents {
		// 只解析交易流转的topic
		if _, ok := TopicEventDataKey[event.Topic]; !ok {
			continue
		}

		if event.EventData == nil || event.EventData.TokenId == "" {
			continue
		}

		//合约不存在
		contract, ok := contractMap[event.ContractName]
		if !ok {
			continue
		}

		tokenKey := event.EventData.TokenId + "_" + contract.Addr
		// 将事件添加到tokenEvents
		tokenEvents[tokenKey] = append(tokenEvents[tokenKey], event)
	}

	if len(tokenEvents) == 0 {
		return tokenResult
	}

	// 用于存储每个tokenId的最后持有人信息
	tokenOwners := GetTokenFinalOwner(chainId, tokenEvents, accountMap)
	for tokenKey := range tokenEvents {
		ownerEvent, ok := tokenOwners[tokenKey]
		if !ok || ownerEvent.Event == nil {
			// token持有人没变
			continue
		}
		//合约不存在
		contract, ok := contractMap[ownerEvent.Event.ContractName]
		if !ok {
			continue
		}
		tokenParam, shouldDelete := processTokenEvents(ownerEvent, contract)
		if shouldDelete {
			tokenResult.DeleteToken = append(tokenResult.DeleteToken, tokenParam)
		} else {
			tokenResult.InsertUpdateToken = append(tokenResult.InsertUpdateToken, tokenParam)
		}
	}

	return tokenResult
}

// processTokenEvents
//
//	@Description: 构造token信息，判断是否被删除
//	@param eventTokenInfo
//	@param contractInfoMap
//	@param accountMap
//	@return *db.NonFungibleToken
//	@return bool
func processTokenEvents(eventTokenInfo *EventTokenInfo, contract *db.Contract) (
	*db.NonFungibleToken, bool) {
	parsEventData := eventTokenInfo.Event.EventData
	newUUID := uuid.New().String()
	tokenParam := &db.NonFungibleToken{
		ID: newUUID,
		//AddrType:     eventTokenInfo.AddrType,
		TokenId:      parsEventData.TokenId,
		ContractAddr: contract.Addr,
		ContractName: contract.Name,
		MetaData:     parsEventData.Metadata,
		CategoryName: parsEventData.CategoryName,
		OwnerAddr:    eventTokenInfo.FinalOwner,
		Timestamp:    time.Now().Unix(),
	}

	// 持有人为空就是销毁token了
	if eventTokenInfo.FinalOwner == "" {
		return tokenParam, true
	}

	return tokenParam, false
}

// GetTokenFinalOwner
//
//	@Description: 获取token的最终持有人
//	@param tokenEvents
//	@return map[string]string
func GetTokenFinalOwner(chainId string, tokenEvents map[string][]*db.ContractEventData,
	accountMap map[string]*db.Account) map[string]*EventTokenInfo {
	tokenInfos := make(map[string]*EventTokenInfo)
	kongAddr := "kong"
	for tokenKey, tokens := range tokenEvents {
		var finalOwner string
		var firstOwner string
		//持有人是否发生变化
		var isEffective bool
		//地址持有某一个token的数量，当创建或者销毁时fromAddress，toAddress会为kong
		balances := make(map[string]int)
		eventMintTokenInfo := &EventTokenInfo{}
		for _, event := range tokens {
			if eventMintTokenInfo.Event == nil ||
				event.EventData.Metadata != "" ||
				event.EventData.CategoryName != "" {
				eventMintTokenInfo.Event = event
			}

			fromAddress := event.EventData.FromAddress
			toAddress := event.EventData.ToAddress
			//转出或销毁，持有量减1
			if fromAddress != "" {
				balances[fromAddress]--
			} else {
				balances[kongAddr]--
			}

			//转入或增发，持有量加1
			if toAddress != "" {
				balances[toAddress]++
			} else {
				balances[kongAddr]++
			}
		}

		for address, balance := range balances {
			//持有地址，持有量大于0的最终持有这个token
			if balance > 0 {
				isEffective = true
				//如果最终的持有地址是kong，说明token销毁了， finalOwner=""
				if address != kongAddr {
					finalOwner = address
				}
			} else if balance < 0 {
				//token从这里转出
				//如果最终的持有地址是kong，说明token初始化
				if address != kongAddr {
					firstOwner = address
				}
			}
		}

		//a-b，b-a，代币在同一个区块又流转回去了，不需要更新，isEffective = false
		if !isEffective {
			continue
		}

		// finalOwner持有账户类型
		addrType := AddrTypeUser
		if account, accOk := accountMap[finalOwner]; accOk {
			addrType = account.AddrType
		}

		//最终finalOwner持有tokenId
		eventMintTokenInfo.FinalOwner = finalOwner
		eventMintTokenInfo.FirstOwner = firstOwner
		eventMintTokenInfo.AddrType = addrType
		tokenInfos[tokenKey] = eventMintTokenInfo
	}

	return tokenInfos
}

//// DealNonFungibleTokenOld 处理Token
//func DealNonFungibleTokenOld(contractEvents []*db.ContractEventData, contractInfoMap map[string]*db.Contract,
//	accountMap map[string]*db.Account) *db.TokenResult {
//	tokenResult := &db.TokenResult{
//		InsertUpdateToken: make([]*db.NonFungibleToken, 0),
//		DeleteToken:       make([]*db.NonFungibleToken, 0),
//	}
//	tokenList := make(map[string][]*db.ContractEventData)
//	contractNameList := make([]string, 0)
//	for _, event := range contractEvents {
//		//只解析交易流转的topic
//		if _, ok := TopicEventDataKey[event.Topic]; !ok {
//			continue
//		}
//
//		if event.EventData == nil || event.EventData.TokenId == "" {
//			continue
//		}
//
//		tokenId := event.EventData.TokenId
//		contractNameList = append(contractNameList, event.ContractName)
//		tokenList[tokenId] = append(tokenList[tokenId], event)
//	}
//
//	if len(tokenList) == 0 {
//		return tokenResult
//	}
//
//	//token最终归属
//	tokenFinalOwner := make(map[string]string, 0)
//	for tokenId, tokens := range tokenList {
//		var finalOwner string
//		//持有人是否发生变化
//		var isEffective bool
//		//地址持有某一个token的数量，当创建或者销毁时fromAddress，toAddress会为kong
//		balances := make(map[string]int)
//		for _, event := range tokens {
//			fromAddress := event.EventData.FromAddress
//			toAddress := event.EventData.ToAddress
//			//转出或销毁，持有量减1
//			if fromAddress != "" {
//				balances[fromAddress]--
//			} else {
//				balances["kong"]--
//			}
//
//			//转入或增发，持有量加1
//			if toAddress != "" {
//				balances[toAddress]++
//			} else {
//				balances["kong"]++
//			}
//		}
//
//		for address, balance := range balances {
//			//持有地址，持有量大于0的最终持有这个token
//			if balance > 0 {
//				isEffective = true
//				//如果最终的持有地址是kong，说明token销毁了， finalOwner=""
//				if address != "kong" {
//					finalOwner = address
//				}
//				break
//			}
//		}
//
//		//a-b，b-a，代币在同一个区块又流转回去了，不需要更新，isEffective = false
//		if !isEffective {
//			continue
//		}
//
//		//最终finalOwner持有tokenId
//		tokenFinalOwner[tokenId] = finalOwner
//	}
//
//	for tokenId, tokens := range tokenList {
//		parsEventData := tokens[0].EventData
//		contractName := tokens[0].ContractName
//		//合约不存在
//		contract, ok := contractInfoMap[contractName]
//		if !ok {
//			continue
//		}
//
//		//token持有人没变
//		finalOwner, ok := tokenFinalOwner[tokenId]
//		if !ok {
//			continue
//		}
//
//		newUUID := uuid.New().String()
//		tokenParam := &db.NonFungibleToken{
//			ID:           newUUID,
//			TokenId:      parsEventData.TokenId,
//			ContractAddr: contract.Addr,
//			ContractName: contract.Name,
//			MetaData:     parsEventData.Metadata,
//			CategoryName: parsEventData.CategoryName,
//			OwnerAddr:    finalOwner,
//			Timestamp:    time.Now().Unix(),
//		}
//
//		//持有人为空就是销毁token了
//		if finalOwner == "" {
//			tokenResult.DeleteToken = append(tokenResult.DeleteToken, tokenParam)
//		} else {
//			//finalOwner持有账户类型
//			addrType := AddrTypeUser
//			if account, accOk := accountMap[finalOwner]; accOk {
//				addrType = account.AddrType
//			}
//
//			tokenParam.AddrType = addrType
//			tokenResult.InsertUpdateToken = append(tokenResult.InsertUpdateToken, tokenParam)
//		}
//	}
//
//	return tokenResult
//}
