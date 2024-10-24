/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"encoding/json"

	"github.com/google/uuid"

	"chainmaker.org/chainmaker/contract-utils/standard"
	"chainmaker.org/chainmaker/pb-go/v2/common"
)

type TopicEventResult struct {
	AddBlack          []string
	DeleteBlack       []string
	IdentityContract  []*db.IdentityContract
	ContractEventData []*db.ContractEventData
	OwnerAdders       []string
	DIDAccount        map[string][]string
	DIDUnBindList     []string
	BNSBindEventData  []*db.BNSTopicEventData
	BNSUnBindDomain   []string
	IDAEventData      *IDAEventData
}

func newTopicEventResult() *TopicEventResult {
	return &TopicEventResult{
		AddBlack:          make([]string, 0),
		DeleteBlack:       make([]string, 0),
		IdentityContract:  make([]*db.IdentityContract, 0),
		ContractEventData: make([]*db.ContractEventData, 0),
		OwnerAdders:       make([]string, 0),
		DIDAccount:        make(map[string][]string, 0),
		DIDUnBindList:     make([]string, 0),
		BNSBindEventData:  make([]*db.BNSTopicEventData, 0),
		BNSUnBindDomain:   make([]string, 0),
	}
}

type IDAEventData struct {
	IDACreatedMap     map[string][]*db.IDACreatedInfo
	IDAUpdatedMap     map[string][]*db.EventIDAUpdatedData
	IDADeletedCodeMap map[string][]string
	EventTime         int64
}

func newIDAEventResult() *IDAEventData {
	return &IDAEventData{
		IDACreatedMap:     make(map[string][]*db.IDACreatedInfo),
		IDAUpdatedMap:     make(map[string][]*db.EventIDAUpdatedData),
		IDADeletedCodeMap: make(map[string][]string),
	}
}

// DealContractEvents 解析所有合约事件
func DealContractEvents(txInfo *common.Transaction) []*db.ContractEvent {
	contractEvents := make([]*db.ContractEvent, 0)
	//失败的操作不处理
	if txInfo.Result.ContractResult == nil ||
		txInfo.Result.ContractResult.Code != 0 {
		return contractEvents
	}

	resEvent := txInfo.Result.ContractResult.ContractEvent
	// 处理合约交易事件
	for i, event := range resEvent {
		eventDataJson, _ := json.Marshal(event.EventData)
		newUUID := uuid.New().String()
		contractEvent := &db.ContractEvent{
			ID:              newUUID,
			Topic:           event.Topic,
			EventIndex:      i + 1,
			TxId:            event.TxId,
			ContractName:    event.ContractName,
			ContractNameBak: event.ContractName,
			ContractVersion: event.ContractVersion,
			EventData:       string(eventDataJson),
			Timestamp:       txInfo.Payload.Timestamp,
		}
		contractEvents = append(contractEvents, contractEvent)
	}
	return contractEvents
}

func parseEventData(event *db.ContractEvent) []string {
	var eventData []string
	if event.EventDataBak != "" {
		_ = json.Unmarshal([]byte(event.EventDataBak), &eventData)
	} else if event.EventData != "" {
		_ = json.Unmarshal([]byte(event.EventData), &eventData)
	}
	return eventData
}

func processIDAEvent(event *db.ContractEvent, eventData []string, idaEventResult *IDAEventData) {
	idaCreatedMap := idaEventResult.IDACreatedMap
	idaUpdatedMap := idaEventResult.IDAUpdatedMap
	idaDeletedCodeMap := idaEventResult.IDADeletedCodeMap
	idaInfos, idaUpdateData, idaDeleteIds := BuildIDAEventData(event.ContractType, event.Topic, eventData)
	if len(idaInfos) > 0 {
		createdInfos := make([]*db.IDACreatedInfo, 0)
		for _, idaInfo := range idaInfos {
			createdInfo := &db.IDACreatedInfo{
				IDAInfo:      idaInfo,
				ContractAddr: event.ContractAddr,
				EventTime:    event.Timestamp,
			}
			createdInfos = append(createdInfos, createdInfo)
		}

		// 检查并初始化切片
		if _, ok := idaCreatedMap[event.ContractAddr]; !ok {
			idaCreatedMap[event.ContractAddr] = make([]*db.IDACreatedInfo, 0)
		}
		idaCreatedMap[event.ContractAddr] = append(idaCreatedMap[event.ContractAddr], createdInfos...)
	}

	if idaUpdateData != nil {
		idaUpdateData.EventTime = event.Timestamp
		if _, ok := idaUpdatedMap[idaUpdateData.IDACode]; !ok {
			idaUpdatedMap[idaUpdateData.IDACode] = make([]*db.EventIDAUpdatedData, 0)
		}
		idaUpdatedMap[idaUpdateData.IDACode] = append(idaUpdatedMap[idaUpdateData.IDACode], idaUpdateData)
	}

	if len(idaDeleteIds) > 0 {
		// 检查并初始化切片
		if _, ok := idaDeletedCodeMap[event.ContractAddr]; !ok {
			idaDeletedCodeMap[event.ContractAddr] = make([]string, 0)
		}
		idaDeletedCodeMap[event.ContractAddr] = append(idaDeletedCodeMap[event.ContractAddr], idaDeleteIds...)
		idaEventResult.EventTime = event.Timestamp
	}
}

// DealTopicEventData 解析eventDate
func DealTopicEventData(contractEvent []*db.ContractEvent, contractInfoMap map[string]*db.Contract,
	txInfoMap map[string]*db.Transaction) *TopicEventResult {
	ownerAddrMap := make(map[string]string, 0)
	didAccountMap := make(map[string][]string, 0)
	bnsAccountMap := make(map[string]*db.BNSTopicEventData, 0)
	idaEventResult := newIDAEventResult()
	topicEventResult := newTopicEventResult()

	if len(contractEvent) == 0 {
		return topicEventResult
	}

	for _, event := range contractEvent {
		eventData := parseEventData(event)
		if len(eventData) == 0 {
			continue
		}

		//解析IDA数据
		processIDAEvent(event, eventData, idaEventResult)
		//解析BNS事件，BNS绑定，解绑
		processBNSEvent(event, eventData, bnsAccountMap, topicEventResult)
		//解析DID事件，设置DID
		processDIDEvent(event, eventData, didAccountMap, topicEventResult)
		//解析黑名单交易，添加，删除黑名单
		processBlackListEvent(topicEventResult, event.ContractName, event.Topic, eventData)

		//合约信息
		contract := contractInfoMap[event.ContractName]
		if contract == nil {
			continue
		}

		//解析身份认证合约
		BuildIdentityEventData(topicEventResult, contractInfoMap, event, eventData)

		//流转数据解析
		var senderUser string
		if txInfo, ok := txInfoMap[event.TxId]; ok {
			senderUser = txInfo.UserAddr
		}
		//根据eventData解析transfer流转记录
		ownerAddrMap = BuildTransferEventData(topicEventResult, ownerAddrMap, contractInfoMap, event,
			senderUser, eventData)
	}

	//持仓地址列表
	ownerAdders := make([]string, 0)
	for _, addr := range ownerAddrMap {
		ownerAdders = append(ownerAdders, addr)
	}
	topicEventResult.OwnerAdders = ownerAdders

	//DID列表
	if len(didAccountMap) > 0 {
		topicEventResult.DIDAccount = didAccountMap
	}

	for _, value := range bnsAccountMap {
		topicEventResult.BNSBindEventData = append(topicEventResult.BNSBindEventData, value)
	}

	//IDA数据
	topicEventResult.IDAEventData = idaEventResult
	return topicEventResult
}

// 解析BNS事件，BNS绑定，解绑
func processBNSEvent(event *db.ContractEvent, eventData []string, bnsAccountMap map[string]*db.BNSTopicEventData,
	topicEventResult *TopicEventResult) {
	bnsBindEventData, bnsUnBindDomain := DealUserBNSEventData(event.ContractName, event.Topic, eventData)
	if bnsBindEventData != nil {
		bnsAccountMap[bnsBindEventData.Domain] = bnsBindEventData
	}
	//bns解绑
	if bnsUnBindDomain != "" {
		//如果前面有绑定，需要删除
		delete(bnsAccountMap, bnsUnBindDomain)
		topicEventResult.BNSUnBindDomain = append(topicEventResult.BNSUnBindDomain, bnsUnBindDomain)
	}
}

// processBlackListEvent 构造交易黑名单数据
func processBlackListEvent(topicEventResult *TopicEventResult, contractName, topic string, eventData []string) {
	//解析黑名单交易，添加，删除黑名单
	addBlack, deleteBlack := DealBackListEventData(contractName, topic, eventData)
	if len(addBlack) > 0 {
		topicEventResult.AddBlack = append(topicEventResult.AddBlack, addBlack...)
	}
	if len(deleteBlack) > 0 {
		topicEventResult.DeleteBlack = append(topicEventResult.DeleteBlack, deleteBlack...)
	}
}

// 解析DID事件，设置DID
func processDIDEvent(event *db.ContractEvent, eventData []string, didAccountMap map[string][]string,
	topicEventResult *TopicEventResult) {
	did, didAddrs := BuildDIDEventData(event.ContractName, event.Topic, eventData)
	if did != "" {
		//如果重复绑定，会覆盖didAccountMap的数据，以最后一次绑定为准
		didAccountMap[did] = didAddrs
		//DID解绑列表，
		//所有涉及到的DID都先解绑，在重新绑定。
		topicEventResult.DIDUnBindList = append(topicEventResult.DIDUnBindList, did)
	}
}

// BuildTransferEventData 构造流转数据
func BuildTransferEventData(topicEventResult *TopicEventResult, ownerAddrMap map[string]string,
	contractInfoMap map[string]*db.Contract, event *db.ContractEvent, senderUser string,
	eventData []string) map[string]string {
	//合约信息
	contract := contractInfoMap[event.ContractName]
	if contract == nil {
		return ownerAddrMap
	}

	//根据eventData解析transfer流转记录
	topicEventData := GetTransferEventDta(contract.ContractType, event.Topic, senderUser, eventData)
	if topicEventData == nil {
		return ownerAddrMap
	}

	if topicEventData.TokenId != "" || topicEventData.Amount != "" {
		//统计持仓地址
		if topicEventData.FromAddress != "" {
			if _, ok := ownerAddrMap[topicEventData.FromAddress]; !ok {
				ownerAddrMap[topicEventData.FromAddress] = topicEventData.FromAddress
			}
		}
		if topicEventData.ToAddress != "" {
			if _, ok := ownerAddrMap[topicEventData.ToAddress]; !ok {
				ownerAddrMap[topicEventData.ToAddress] = topicEventData.ToAddress
			}
		}
	}

	transferData := &db.ContractEventData{
		Topic:        event.Topic,
		Index:        event.EventIndex,
		TxId:         event.TxId,
		ContractName: event.ContractName,
		EventData:    topicEventData,
		Timestamp:    event.Timestamp,
	}
	topicEventResult.ContractEventData = append(topicEventResult.ContractEventData, transferData)
	return ownerAddrMap
}

// BuildIdentityEventData 解析身份认证合约
func BuildIdentityEventData(topicEventResult *TopicEventResult, contractInfoMap map[string]*db.Contract,
	event *db.ContractEvent, eventData []string) {
	//合约信息
	contract := contractInfoMap[event.ContractName]
	if contract == nil || contract.ContractType != ContractStandardNameCMID {
		return
	}

	//解析身份认证合约
	identityEventData := DealIdentityEventData(event.Topic, eventData)
	if identityEventData == nil {
		return
	}

	newUUID := uuid.New().String()
	tempIdentity := &db.IdentityContract{
		ID:           newUUID,
		TxId:         event.TxId,
		EventIndex:   event.EventIndex,
		ContractName: contract.Name,
		ContractAddr: contract.Addr,
		UserAddr:     identityEventData.UserAddr,
		Level:        identityEventData.Level,
		PkPem:        identityEventData.PkPem,
	}
	topicEventResult.IdentityContract = append(topicEventResult.IdentityContract, tempIdentity)
}

// BuildBNSEventData 构造BNS数据
func BuildBNSEventData(topicEventResult *TopicEventResult, contractName, topic string, eventData []string) {
	//解析BNS事件，BNS绑定，解绑
	bnsBindEventData, bnsUnBindDomain := DealUserBNSEventData(contractName, topic, eventData)
	if bnsBindEventData != nil {
		//BNS绑定事件
		topicEventResult.BNSBindEventData = append(topicEventResult.BNSBindEventData, bnsBindEventData)
	}
	if bnsUnBindDomain != "" {
		//BNS解绑事件
		topicEventResult.BNSUnBindDomain = append(topicEventResult.BNSUnBindDomain, bnsUnBindDomain)
	}
}

// BuildDIDEventData 构造DID数据
func BuildDIDEventData(contractName, topic string, eventData []string) (string, []string) {
	var accountAddrs []string
	didDocument := DealUserDIDEventData(contractName, topic, eventData)
	if didDocument == nil {
		return "", accountAddrs
	}

	//一个DID绑定的所有账户
	for _, value := range didDocument.VerificationMethod {
		accountAddrs = append(accountAddrs, value.Address)
	}

	return didDocument.Did, accountAddrs
}

// UpdateContractTxAndEventNum 更新合约交易数和事件数
func UpdateContractTxAndEventNum(minHeight int64, contractMap map[string]*db.Contract,
	txList map[string]*db.Transaction, contractEvent []*db.ContractEvent) []*db.Contract {
	contractTxNumMap := make(map[string]int64, 0)
	contractEventNumMap := make(map[string]int64, 0)
	updateContractMap := make(map[string]*db.Contract, 0)
	updateContractNum := make([]*db.Contract, 0)
	if len(contractMap) == 0 || len(txList) == 0 {
		return updateContractNum
	}

	//统计本次交易数据量
	for _, txInfo := range txList {
		if contract, ok := contractMap[txInfo.ContractNameBak]; ok {
			//说明已经更新过了
			if contract.BlockHeight >= minHeight {
				continue
			}
			contractTxNumMap[contract.Addr]++
		}
	}

	//统计本次合约事件数据量
	for _, event := range contractEvent {
		if contract, ok := contractMap[event.ContractNameBak]; ok {
			//说明已经更新过了
			if contract.BlockHeight >= minHeight {
				continue
			}
			contractEventNumMap[contract.Addr]++
		}
	}

	for addr, txNum := range contractTxNumMap {
		var contractInfo *db.Contract
		if contract, ok := updateContractMap[addr]; ok {
			contractInfo = contract
		} else if contract, ok = contractMap[addr]; ok {
			contractInfo = contract
		} else {
			continue
		}
		contractInfo.TxNum = contractInfo.TxNum + txNum
		contractInfo.BlockHeight = minHeight
		updateContractMap[addr] = contractInfo
	}

	for addr, eventNum := range contractEventNumMap {
		var contractInfo *db.Contract
		if contract, ok := updateContractMap[addr]; ok {
			contractInfo = contract
		} else if contract, ok = contractMap[addr]; ok {
			contractInfo = contract
		} else {
			continue
		}
		contractInfo.EventNum = contractInfo.EventNum + eventNum
		contractInfo.BlockHeight = minHeight
		updateContractMap[addr] = contractInfo
	}

	for _, contract := range updateContractMap {
		updateContractNum = append(updateContractNum, contract)
	}

	return updateContractNum
}

// DealEvidence 处理存证合约
func DealEvidence(blockHeight int64, txInfo *common.Transaction, userInfo *MemberAddrIdCert) (
	evidences []*db.EvidenceContract, err error) {
	evidences = make([]*db.EvidenceContract, 0)
	if txInfo.Payload.Method != PayloadMethodEvidence &&
		txInfo.Payload.Method != PayloadMethodEvidenceBatch {
		return evidences, nil
	}

	contractName := txInfo.Payload.ContractName
	tempEvidence := &db.EvidenceContract{
		ContractName:       contractName,
		TxId:               txInfo.Payload.TxId,
		SenderAddr:         userInfo.UserAddr,
		Timestamp:          txInfo.Payload.Timestamp,
		BlockHeight:        blockHeight,
		ContractResult:     txInfo.Result.ContractResult.Result,
		ContractResultCode: txInfo.Result.ContractResult.Code,
	}
	if txInfo.Payload.Method == PayloadMethodEvidence {
		for _, parameter := range txInfo.Payload.Parameters {
			if parameter.Key == "hash" {
				tempEvidence.Hash = string(parameter.Value)
			}
			if parameter.Key == "metadata" {
				tempEvidence.MetaData = string(parameter.Value)
			}
			if parameter.Key == "id" {
				tempEvidence.EvidenceId = string(parameter.Value)
			}
		}
		newUUID := uuid.New().String()
		tempEvidence.ID = newUUID
		evidences = append(evidences, tempEvidence)
	} else if txInfo.Payload.Method == PayloadMethodEvidenceBatch {
		for _, parameter := range txInfo.Payload.Parameters {
			if parameter.Key == "evidences" {
				standardEvidences := make([]standard.Evidence, 0)
				err = json.Unmarshal(parameter.Value, &standardEvidences)
				if err != nil {
					return evidences, err
				}
				for _, e := range standardEvidences {
					newUUID := uuid.New().String()
					tempEvidence.EvidenceId = e.Id
					tempEvidence.Hash = e.Hash
					tempEvidence.MetaData = e.Metadata
					tempEvidence.ID = newUUID
					evidences = append(evidences, tempEvidence)
				}
			}
		}
	}
	return evidences, err
}

// BuildIDAEventData 解析IDA数据
func BuildIDAEventData(contractType, topic string, eventData []string) (
	[]*standard.IDAInfo, *db.EventIDAUpdatedData, []string) {
	idaIds := make([]string, 0)
	idaInfoList := make([]*standard.IDAInfo, 0)
	var idaUpdateData *db.EventIDAUpdatedData
	//判断是否是IDA合约
	if contractType != standard.ContractStandardNameCMIDA {
		return idaInfoList, idaUpdateData, idaIds
	}

	switch topic {
	case standard.EventIDACreated:
		idaInfoList = DealEventIDACreated(eventData)
	case standard.EventIDAUpdated:
		idaUpdateData = DealEventIDAUpdated(eventData)
	case standard.EventIDADeleted:
		idaIds = DealEventIDADeleted(eventData)
	}

	return idaInfoList, idaUpdateData, idaIds
}
