/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"encoding/json"
	"math/big"

	"chainmaker.org/chainmaker/contract-utils/standard"

	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
)

func GetTransferEventDta(contractType, topic, senderUser string, eventData []string) *db.TransferTopicEventData {
	var topicEventData *db.TransferTopicEventData
	switch contractType {
	case ContractStandardNameCMDFA:
		topicEventData = DealDockerDFAEventData(topic, eventData)
	case ContractStandardNameCMNFA:
		topicEventData = DealDockerNFAEventData(topic, eventData, senderUser)
	case ContractStandardNameEVMDFA:
		topicEventData = DealEVMDFAEventData(topic, eventData)
	case ContractStandardNameEVMNFA:
		topicEventData = DealEVMNFAEventData(topic, eventData)
	}
	if topicEventData != nil {
		//有部分初始化会给默认FromAddress
		if topicEventData.FromAddress == topicEventData.ToAddress {
			topicEventData.FromAddress = ""
		}
	}

	return topicEventData
}

// DealBackListEventData 解析黑名单交易
func DealBackListEventData(contractName, topic string, eventData []string) ([]string, []string) {
	addBlack := make([]string, 0)
	deleteBlack := make([]string, 0)
	//第一条记录是链ID
	if contractName != syscontract.SystemContract_TRANSACTION_MANAGER.String() || len(eventData) <= 1 {
		return addBlack, deleteBlack
	}

	//加入黑名单
	if topic == TopicTxAddBlack {
		for i := 1; i < len(eventData); i++ {
			addBlack = append(addBlack, eventData[i])
		}
	} else if topic == TopicTxDeleteBlack {
		//解封黑名单
		for i := 1; i < len(eventData); i++ {
			deleteBlack = append(deleteBlack, eventData[i])
		}
	}
	return addBlack, deleteBlack
}

// DealIdentityEventData 解析身份认证数据
func DealIdentityEventData(topic string, eventData []string) *db.IdentityEventData {
	var identityEvent *db.IdentityEventData
	if topic == "setIdentity" {
		if len(eventData) < 3 {
			return identityEvent
		}
		identityEvent = &db.IdentityEventData{
			UserAddr: eventData[0],
			Level:    eventData[1],
			PkPem:    eventData[2],
		}
	}
	return identityEvent
}

// DealDockerDFAEventData 同质化docker解析eventdata
func DealDockerDFAEventData(topic string, eventData []string) *db.TransferTopicEventData {
	var topicEventData *db.TransferTopicEventData
	switch topic {
	case "mint":
		if len(eventData) < 2 {
			return topicEventData
		}
		topicEventData = &db.TransferTopicEventData{
			ToAddress: eventData[0],
			Amount:    eventData[1],
		}
	case "transfer":
		if len(eventData) < 3 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			ToAddress:   eventData[1],
			Amount:      eventData[2],
		}
	case "burn":
		if len(eventData) < 2 {
			return topicEventData
		}
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			Amount:      eventData[1],
		}
	case "approve":
		if len(eventData) < 3 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			ToAddress:   eventData[1],
			Amount:      eventData[2],
		}
	}
	return topicEventData
}

// DealDockerNFAEventData 非同质化docker解析eventdata
func DealDockerNFAEventData(topic string, eventData []string, senderUser string) *db.TransferTopicEventData {
	var topicEventData *db.TransferTopicEventData
	switch topic {
	case "Mint":
		//发行
		if len(eventData) < 5 {
			return topicEventData
		}
		if isZeroAddress(eventData[0]) {
			eventData[0] = ""
		}
		topicEventData = &db.TransferTopicEventData{
			FromAddress:  eventData[0],
			ToAddress:    eventData[1],
			TokenId:      eventData[2],
			CategoryName: eventData[3],
			Metadata:     eventData[4],
		}
	case "TransferFrom":
		//转账
		if len(eventData) < 3 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			ToAddress:   eventData[1],
			TokenId:     eventData[2],
		}
	case "Burn":
		//销毁
		if len(eventData) == 0 {
			return topicEventData
		}
		topicEventData = &db.TransferTopicEventData{
			FromAddress: senderUser,
			TokenId:     eventData[0],
		}
	case "SetApproval":
		//授权
		if len(eventData) < 4 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			ToAddress:   eventData[1],
			TokenId:     eventData[2],
			Approval:    eventData[3],
		}
	case "SetApprovalForAll":
		//全部授权
		if len(eventData) < 3 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			ToAddress:   eventData[1],
			Approval:    eventData[2],
		}
	case "SetApprovalByCategory":
		//分类授权
		if len(eventData) < 4 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress:  eventData[0],
			ToAddress:    eventData[1],
			CategoryName: eventData[2],
			Approval:     eventData[3],
		}
	}

	return topicEventData
}

// DealEVMDFAEventData 解析eventdata
func DealEVMDFAEventData(topic string, eventData []string) *db.TransferTopicEventData {
	var topicEventData *db.TransferTopicEventData
	switch topic {
	case EVMEventTopicTransfer:
		if len(eventData) < 3 {
			return topicEventData
		}
		if len(eventData[0]) < 24 || len(eventData[1]) < 24 {
			return topicEventData
		}
		// 解析数据
		fromAddress := eventData[0][24:]
		toAddress := eventData[1][24:]
		// 将十六进制字符串转换为big.Int类型
		bigInt := new(big.Int)
		bigInt.SetString(eventData[2], 16)
		// 将big.Int类型转换为十进制字符串
		amountStr := bigInt.String()
		if isZeroAddress(fromAddress) {
			fromAddress = ""
		}
		if isZeroAddress(toAddress) {
			toAddress = ""
		}
		topicEventData = &db.TransferTopicEventData{
			FromAddress: fromAddress,
			ToAddress:   toAddress,
			Amount:      amountStr,
		}
	}
	return topicEventData
}

// DealEVMNFAEventData 解析eventdata
func DealEVMNFAEventData(topic string, eventData []string) *db.TransferTopicEventData {
	var topicEventData *db.TransferTopicEventData
	switch topic {
	case EVMEventTopicTransfer:
		if len(eventData) < 3 {
			return topicEventData
		}
		if len(eventData[0]) < 24 || len(eventData[1]) < 24 {
			return topicEventData
		}
		// 解析数据
		fromAddress := eventData[0][24:]
		toAddress := eventData[1][24:]
		// 将十六进制字符串转换为big.Int类型
		bigInt := new(big.Int)
		bigInt.SetString(eventData[2], 16)
		// 将big.Int类型转换为十进制字符串
		tokenId := bigInt.String()
		if isZeroAddress(fromAddress) {
			fromAddress = ""
		}
		if isZeroAddress(toAddress) {
			toAddress = ""
		}
		topicEventData = &db.TransferTopicEventData{
			FromAddress: fromAddress,
			ToAddress:   toAddress,
			TokenId:     tokenId,
		}
	}
	return topicEventData
}

// DealUserBNSEventData 解析BNS eventdata
func DealUserBNSEventData(contractName, topic string, eventData []string) (*db.BNSTopicEventData, string) {
	var topicEventData *db.BNSTopicEventData
	var unBindDomain string
	if contractName != PayloadContractNameBNS {
		return topicEventData, unBindDomain
	}

	switch topic {
	case standard.Topic_Bind:
		//绑定BNS
		if len(eventData) < 3 {
			return topicEventData, unBindDomain
		}

		//ResourceType, _ := strconv.ParseInt(eventData[2], 10, 64)
		////BNS解析资源类型,string "1“-链地址，”2"-DID,"3"-去中心化网站，"4“-合约，"5"-子链
		//if ResourceType > 1 {
		//	//暂时只解析链地址
		//	return topicEventData
		//}

		// 解析数据
		topicEventData = &db.BNSTopicEventData{
			Domain:       eventData[0],
			Value:        eventData[1],
			ResourceType: eventData[2],
		}
	case standard.Topic_UnBind:
		//解绑BNS
		if len(eventData) == 0 {
			return topicEventData, unBindDomain
		}
		// 解析数据
		unBindDomain = eventData[0]
	}

	return topicEventData, unBindDomain
}

// DealUserDIDEventData 解析DID eventdata
func DealUserDIDEventData(contractName, topic string, eventData []string) *db.DIDTopicEventData {
	var topicEventData *db.DIDTopicEventData
	switch topic {
	case DIDSetDidDocument:
		//绑定，解绑DID
		if len(eventData) < 2 {
			return topicEventData
		}
		didDocument := &db.DidDocument{}
		err := json.Unmarshal([]byte(eventData[1]), &didDocument)
		if err != nil {
			log.Errorf("DealUserDIDEventData json Unmarshal err, err:%v", err)
		}
		// 解析数据
		topicEventData = &db.DIDTopicEventData{
			Did:                eventData[0],
			VerificationMethod: didDocument.VerificationMethod,
		}
	}
	return topicEventData
}

// DealEventIDACreated 解析IDA event
func DealEventIDACreated(eventData []string) []*standard.IDAInfo {
	//standard.EventIDACreated
	idaInfoList := make([]*standard.IDAInfo, 0)
	if len(eventData) == 0 {
		return idaInfoList
	}
	err := json.Unmarshal([]byte(eventData[0]), &idaInfoList)
	if err != nil {
		log.Errorf("DealEventIDACreated json Unmarshal err, err:%v, eventData:%v", err, eventData)
	}
	return idaInfoList
}

// DealEventIDAUpdated 解析IDA event
func DealEventIDAUpdated(eventData []string) *db.EventIDAUpdatedData {
	updateData := &db.EventIDAUpdatedData{}
	if len(eventData) < 3 {
		return updateData
	}
	updateData.IDACode = eventData[0]
	updateData.Field = eventData[1]
	updateData.Update = eventData[2]
	return updateData
}

// DealEventIDADeleted 解析IDA event
func DealEventIDADeleted(eventData []string) []string {
	idaCodes := make([]string, 0)
	if len(eventData) == 0 {
		return idaCodes
	}
	idaCodes = eventData
	return idaCodes
}

// DealEventIDACreated 解析IDA event
func UnmarshalIDAUpdatedBasic(updateJson string) (*standard.Basic, error) {
	basicInfo := &standard.Basic{}
	if updateJson == "" {
		return basicInfo, nil
	}
	err := json.Unmarshal([]byte(updateJson), &basicInfo)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedBasic json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return basicInfo, err
}

func UnmarshalIDAUpdatedSupply(updateJson string) (*standard.Supply, error) {
	supplyInfo := &standard.Supply{}
	if updateJson == "" {
		return supplyInfo, nil
	}
	err := json.Unmarshal([]byte(updateJson), &supplyInfo)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedSupply json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return supplyInfo, err
}

func UnmarshalIDAUpdatedDetails(updateJson string) (*standard.Details, error) {
	detailInfo := &standard.Details{}
	if updateJson == "" {
		return detailInfo, nil
	}
	err := json.Unmarshal([]byte(updateJson), &detailInfo)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedDetails json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return detailInfo, err
}

func UnmarshalIDAUpdatedOwnership(updateJson string) (*standard.Ownership, error) {
	ownershipInfo := &standard.Ownership{}
	if updateJson == "" {
		return ownershipInfo, nil
	}
	err := json.Unmarshal([]byte(updateJson), &ownershipInfo)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedOwnership json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return ownershipInfo, err
}

func UnmarshalIDAUpdatedColumns(updateJson string) ([]*standard.ColumnInfo, error) {
	columns := make([]*standard.ColumnInfo, 0)
	if updateJson == "" {
		return columns, nil
	}
	err := json.Unmarshal([]byte(updateJson), &columns)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedColumns json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return columns, err
}

func UnmarshalIDAUpdatedApis(updateJson string) ([]*standard.APIInfo, error) {
	apis := make([]*standard.APIInfo, 0)
	if updateJson == "" {
		return apis, nil
	}
	err := json.Unmarshal([]byte(updateJson), &apis)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedApis json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return apis, err
}
