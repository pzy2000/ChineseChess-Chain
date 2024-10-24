package service

import (
	"bytes"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"

	"github.com/gin-gonic/gin"
)

// SensitiveWordColumn 敏感词字段
const (
	SensitiveContractName       = "ContractName"
	SensitiveContractResult     = "ContractResult"
	SensitiveContractMessage    = "ContractMessage"
	SensitiveContractParameters = "ContractParameters"
	SensitiveReadSet            = "ReadSet"
	SensitiveWriteSet           = "WriteSet"
	SensitiveTopic              = "Topic"
	SensitiveEventData          = "EventData"
	SensitiveEvidenceMetaData   = "EvidenceMetaData"
	SensitiveTokenMetaData      = "TokenMetaData"
)

const (
	SensitiveStatusAdd    = 1
	SensitiveStatusDelete = 0
)

// ModifyTxBlackListHandler update
type ModifyTxBlackListHandler struct{}

// Handle 修改交易黑名单
func (ModifyTxBlackListHandler *ModifyTxBlackListHandler) Handle(ctx *gin.Context) {
	params := entity.BindModifyTxBlackListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// 验证API密钥
	apiKey := ctx.GetHeader("x-api-key")
	if !ValidateAPIKey(apiKey) {
		newError := entity.NewError(entity.ErrorParamWrong, "api key is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	if params.Status != nil && *params.Status == SensitiveStatusAdd {
		//加入黑名单
		//查询交易是否还在
		txInfo, err := dbhandle.GetTransactionByTxId(params.TxId, params.ChainId)
		if err == nil && txInfo == nil {
			//交易已经不在了
			ConvergeDataResponse(ctx, "OK", nil)
			return
		} else if err != nil || txInfo == nil {
			log.Errorf("GetTransactionByTxId err : %s", err)
			ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
			return
		}

		//加入黑名单
		transactions := make([]*db.BlackTransaction, 0)
		transactions = append(transactions, (*db.BlackTransaction)(txInfo))
		err = dbhandle.InsertBlackTransactions(params.ChainId, transactions)
		if err != nil {
			log.Errorf("InsertBlackTransactions err : %s", err.Error())
			ConvergeHandleFailureResponse(ctx, err)
			return
		}
	} else if params.Status != nil && *params.Status == SensitiveStatusDelete {
		//移除黑名单
		//查询交易是否在黑名单
		blackTxInfo, err := dbhandle.GetBlackTxInfoByTxId(params.ChainId, params.TxId)
		if err == nil && blackTxInfo == nil {
			ConvergeDataResponse(ctx, "OK", nil)
			return
		} else if err != nil || blackTxInfo == nil {
			log.Errorf("GetBlackTxInfoByTxId err : %s", err)
			ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
			return
		}

		//删除黑名单
		transactions := make([]*db.Transaction, 0)
		transactions = append(transactions, (*db.Transaction)(blackTxInfo))
		err = dbhandle.DeleteBlackTransaction(params.ChainId, transactions)
		if err != nil {
			log.Errorf("DeleteBlackTransaction err : %s", err.Error())
			ConvergeHandleFailureResponse(ctx, err)
			return
		}
	}
	ConvergeDataResponse(ctx, "OK", nil)
}

// DeleteTxBlackListHandler update
type DeleteTxBlackListHandler struct{}

// Handle 删除交易黑名单
func (deleteTxBlackListHandler *DeleteTxBlackListHandler) Handle(ctx *gin.Context) {
	params := entity.BindDeleteTxBlackListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//查询交易是否在黑名单
	blackTxInfo, err := dbhandle.GetBlackTxInfoByTxId(params.ChainId, params.TxId)
	if err == db.ErrRecordNotFoundErr {
		ConvergeDataResponse(ctx, "OK", nil)
		return
	} else if err != nil {
		log.Errorf("GetBlackTxInfoByTxId err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//删除黑名单
	transactions := make([]*db.Transaction, 0)
	transactions = append(transactions, (*db.Transaction)(blackTxInfo))
	err = dbhandle.DeleteBlackTransaction(params.ChainId, transactions)
	if err != nil {
		log.Errorf("DeleteBlackTransaction err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	ConvergeDataResponse(ctx, "OK", nil)
}

// UpdateTxSensitiveWordHandler update
type UpdateTxSensitiveWordHandler struct{}

// Handle 更新交易敏感词
func (updateTxSensitiveWordHandler *UpdateTxSensitiveWordHandler) Handle(ctx *gin.Context) {
	params := entity.BindUpdateTxSensitiveWordHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}
	// 验证API密钥
	apiKey := ctx.GetHeader("x-api-key")
	if !ValidateAPIKey(apiKey) {
		newError := entity.NewError(entity.ErrorParamWrong, "api key is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	txInfo, err := dbhandle.GetTransactionByTxId(params.TxId, params.ChainId)
	if err != nil || txInfo == nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	var isUpdate bool
	//合约ContractResult
	if params.Status == 1 {
		//隐藏敏感词
		if params.Column == "" {
			txInfo, isUpdate = TxSensitiveAddAll(params.WarnMsg, txInfo)
		} else {
			txInfo, isUpdate = TxSensitiveAdd(params.Column, params.WarnMsg, txInfo)
		}
	} else if params.Status == 0 {
		//敏感词可见
		if params.Column == "" {
			txInfo, isUpdate = TxSensitiveDeleteAll(txInfo)
		} else {
			txInfo, isUpdate = TxSensitiveDelete(params.Column, txInfo)
		}
	} else {
		newError := entity.NewError(entity.ErrorParamWrong, "Status is wrong")
		ConvergeFailureResponse(ctx, newError)
	}

	//不需要更新
	if !isUpdate {
		ConvergeDataResponse(ctx, "OK", nil)
		return
	}

	//更新交易数据
	err = dbhandle.UpdateTransactionBak(params.ChainId, txInfo)
	if err != nil {
		log.Errorf("UpdateTxSensitiveWord update err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	ConvergeDataResponse(ctx, "OK", nil)
}

func updateField(current *string, bak *string, warnMsg string) bool {
	var isUpdate bool
	if *bak == "" && *current != "" {
		*bak = *current
		isUpdate = true
	}
	if *current != warnMsg {
		*current = warnMsg
		isUpdate = true
	}
	return isUpdate
}

func TxSensitiveAdd(column, warnMsg string, txInfo *db.Transaction) (*db.Transaction, bool) {
	if warnMsg == "" {
		warnMsg = config.OtherWarnMsg
	}

	isUpdate := false
	switch column {
	case SensitiveContractResult:
		if !bytes.Equal(txInfo.ContractResult, config.ContractResultMsg) {
			// 创建一个新的字节切片，长度和容量与 txInfo.ContractResult 相同
			newSlice := make([]byte, len(txInfo.ContractResult))
			// 将 txInfo.ContractResult 的内容复制到新的切片中
			copy(newSlice, txInfo.ContractResult)
			// 将新的切片分配给 txInfo.ContractResultBak
			txInfo.ContractResultBak = newSlice
			isUpdate = true
		}
		if string(txInfo.ContractResult) != warnMsg {
			txInfo.ContractResult = []byte(warnMsg)
			isUpdate = true
		}
	case SensitiveContractMessage:
		isUpdate = updateField(&txInfo.ContractMessage, &txInfo.ContractMessageBak, warnMsg)
	case SensitiveContractParameters:
		isUpdate = updateField(&txInfo.ContractParameters, &txInfo.ContractParametersBak, warnMsg)
	case SensitiveReadSet:
		isUpdate = updateField(&txInfo.ReadSet, &txInfo.ReadSetBak, warnMsg)
	case SensitiveWriteSet:
		isUpdate = updateField(&txInfo.WriteSet, &txInfo.WriteSetBak, warnMsg)
	}

	return txInfo, isUpdate
}

// TxSensitiveAddAll 添加所有交易敏感词
func TxSensitiveAddAll(warnMsg string, txInfo *db.Transaction) (*db.Transaction, bool) {
	var isUpdate bool
	if warnMsg == "" {
		warnMsg = config.OtherWarnMsg
	}

	columns := []string{
		SensitiveContractResult,
		SensitiveContractMessage,
		SensitiveContractParameters,
		SensitiveReadSet,
		SensitiveWriteSet,
	}

	for _, column := range columns {
		_, columnUpdated := TxSensitiveAdd(column, warnMsg, txInfo)
		if !isUpdate && columnUpdated {
			isUpdate = true
		}
	}

	return txInfo, isUpdate
}

// TxSensitiveDeleteAll 删除所有交易敏感词
func TxSensitiveDeleteAll(txInfo *db.Transaction) (*db.Transaction, bool) {
	var isUpdate bool
	columns := []string{
		SensitiveContractResult,
		SensitiveContractMessage,
		SensitiveContractParameters,
		SensitiveReadSet,
		SensitiveWriteSet,
	}

	for _, column := range columns {
		_, columnUpdated := TxSensitiveDelete(column, txInfo)
		if !isUpdate && columnUpdated {
			isUpdate = true
		}
	}

	return txInfo, isUpdate
}

// TxSensitiveDelete 删除交易敏感词
func TxSensitiveDelete(column string, txInfo *db.Transaction) (*db.Transaction, bool) {
	var isUpdate bool
	switch column {
	case SensitiveContractResult:
		if len(txInfo.ContractResultBak) != 0 {
			// 创建一个新的字节切片，长度和容量与 txInfo.ContractResultBak 相同
			newSlice := make([]byte, len(txInfo.ContractResultBak))
			// 将 txInfo.ContractResultBak 的内容复制到新的切片中
			copy(newSlice, txInfo.ContractResultBak)
			// 将新的切片分配给 txInfo.ContractResult
			txInfo.ContractResult = newSlice
			// 将 txInfo.ContractResultBak 设置为空字节切片
			txInfo.ContractResultBak = []byte("")
			isUpdate = true
		}
	case SensitiveContractMessage:
		if txInfo.ContractMessageBak != "" {
			txInfo.ContractMessage = txInfo.ContractMessageBak
			txInfo.ContractMessageBak = ""
			isUpdate = true
		}
	case SensitiveContractParameters:
		if txInfo.ContractParametersBak != "" {
			txInfo.ContractParameters = txInfo.ContractParametersBak
			txInfo.ContractParametersBak = ""
			isUpdate = true
		}
	case SensitiveReadSet:
		if txInfo.ReadSetBak != "" {
			txInfo.ReadSet = txInfo.ReadSetBak
			txInfo.ReadSetBak = ""
			isUpdate = true
		}
	case SensitiveWriteSet:
		if txInfo.WriteSetBak != "" {
			txInfo.WriteSet = txInfo.WriteSetBak
			txInfo.WriteSetBak = ""
			isUpdate = true
		}
	}

	return txInfo, isUpdate
}

// UpdateEvidenceSensitiveWordHandler update
type UpdateEvidenceSensitiveWordHandler struct{}

// Handle 更新事件敏感词
func (updateEvidenceSensitiveWordHandler *UpdateEvidenceSensitiveWordHandler) Handle(ctx *gin.Context) {
	params := entity.BindEvidenceSensitiveWordHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	if params.Column != SensitiveEvidenceMetaData {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	evidenceData, err := dbhandle.GetEvidenceContractByHash(params.ChainId, params.Hash)
	if err != nil || evidenceData == nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	var isUpdate bool
	warnMsg := params.WarnMsg
	if warnMsg == "" {
		warnMsg = config.OtherWarnMsg
	}

	if params.Status != nil && *params.Status == SensitiveStatusAdd {
		if evidenceData.MetaDataBak == "" && evidenceData.MetaData != "" {
			evidenceData.MetaDataBak = evidenceData.MetaData
			isUpdate = true
		}
		if evidenceData.MetaData != warnMsg {
			evidenceData.MetaData = warnMsg
			isUpdate = true
		}
	} else if params.Status != nil && *params.Status == SensitiveStatusDelete {
		if evidenceData.MetaDataBak != "" {
			evidenceData.MetaData = evidenceData.MetaDataBak
			evidenceData.MetaDataBak = ""
			isUpdate = true
		}
	}

	//不需要更新
	if !isUpdate {
		ConvergeDataResponse(ctx, "OK", nil)
	}

	//更新交易数据
	err = dbhandle.UpdateEvidenceBak(params.ChainId, evidenceData)
	if err != nil {
		log.Errorf("AddEvidenceSensitiveWord err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	ConvergeDataResponse(ctx, "OK", nil)
}

// UpdateEventSensitiveWordHandler update
type UpdateEventSensitiveWordHandler struct{}

// Handle 更新事件敏感词
func (updateEventSensitiveWordHandler *UpdateEventSensitiveWordHandler) Handle(ctx *gin.Context) {
	params := entity.BindUpdateEventSensitiveWordHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	eventList, err := dbhandle.GetEventDataByTxIds(params.ChainId, []string{params.TxId})
	if err != nil || len(eventList) == 0 {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}
	eventData := &db.ContractEvent{}
	for _, event := range eventList {
		if event.TxId == params.TxId && event.EventIndex == params.Index {
			eventData = event
			break
		}
	}

	if eventData.TxId == "" {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	isAdd := params.Status != nil && *params.Status == SensitiveStatusAdd
	isDelete := params.Status != nil && *params.Status == SensitiveStatusDelete
	warnMsg := params.WarnMsg
	if warnMsg == "" {
		warnMsg = config.OtherWarnMsg
	}
	if isAdd || isDelete {
		isUpdate := updateEventData(params.Column, warnMsg, eventData, isAdd)
		if isUpdate {
			// 更新交易数据
			err = dbhandle.UpdateContractEventBak(params.ChainId, eventData)
			if err != nil {
				log.Errorf("UpdateEventSensitiveWord err : %s", err.Error())
				ConvergeHandleFailureResponse(ctx, err)
				return
			}
		}
	}

	ConvergeDataResponse(ctx, "OK", nil)
}

func updateEventData(column, warnMsg string, eventData *db.ContractEvent, isAdd bool) bool {
	var isUpdate bool

	switch column {
	case SensitiveTopic:
		if isAdd {
			if eventData.TopicBak == "" && eventData.Topic != "" {
				eventData.TopicBak = eventData.Topic
				isUpdate = true
			}
			if eventData.Topic != warnMsg {
				eventData.Topic = warnMsg
				isUpdate = true
			}
		} else {
			if eventData.TopicBak != "" {
				eventData.Topic = eventData.TopicBak
				eventData.TopicBak = ""
				isUpdate = true
			}
		}
	case SensitiveEventData:
		if isAdd {
			if eventData.EventDataBak == "" && eventData.EventData != "" {
				eventData.EventDataBak = eventData.EventData
				isUpdate = true
			}
			if eventData.EventData != warnMsg {
				eventData.EventData = warnMsg
				isUpdate = true
			}
		} else {
			if eventData.EventDataBak != "" {
				eventData.EventData = eventData.EventDataBak
				eventData.EventDataBak = ""
				isUpdate = true
			}
		}
	}

	return isUpdate
}

// UpdateNFTSensitiveWordHandler update
type UpdateNFTSensitiveWordHandler struct{}

// Handle 更新事件敏感词
func (updateNFTSensitiveWordHandler *UpdateNFTSensitiveWordHandler) Handle(ctx *gin.Context) {
	params := entity.BindNFTSensitiveWordHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	if params.Column != SensitiveTokenMetaData {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	tokenData, err := dbhandle.GetNonFungibleTokenDetail(params.ChainId, params.TokenId, params.ContractAddr)
	if err != nil || tokenData == nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	var isUpdate bool
	warnMsg := params.WarnMsg
	if warnMsg == "" {
		warnMsg = config.OtherWarnMsg
	}
	if params.Status != nil && *params.Status == SensitiveStatusAdd {
		if tokenData.MetaDataBak == "" && tokenData.MetaData != "" {
			tokenData.MetaDataBak = tokenData.MetaData
			isUpdate = true
		}
		if tokenData.MetaData != warnMsg {
			tokenData.MetaData = warnMsg
			isUpdate = true
		}
	} else if params.Status != nil && *params.Status == SensitiveStatusDelete {
		if tokenData.MetaDataBak != "" {
			tokenData.MetaData = tokenData.MetaDataBak
			tokenData.MetaDataBak = ""
			isUpdate = true
		}
	}

	//不需要更新
	if !isUpdate {
		ConvergeDataResponse(ctx, "OK", nil)
	}

	//更新交易数据
	err = dbhandle.UpdateNonFungibleTokenBak(params.ChainId, tokenData)
	if err != nil {
		log.Errorf("AddEvidenceSensitiveWord err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	ConvergeDataResponse(ctx, "OK", nil)
}

// UpdateContractNameSensitiveWordHandler update
type UpdateContractNameSensitiveWordHandler struct{}

// Handle 更新事件敏感词
func (updateContractNameSW *UpdateContractNameSensitiveWordHandler) Handle(ctx *gin.Context) {
	params := entity.BindUpdateContractNameSWHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	contractData, err := dbhandle.GetContractByName(params.ChainId, params.ContractName)
	if err != nil || contractData == nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	var isUpdate bool
	warnMsg := params.WarnMsg
	if warnMsg == "" {
		warnMsg = config.OtherWarnMsg
	}
	if params.Status != nil && *params.Status == SensitiveStatusAdd {
		if contractData.NameBak == "" && contractData.Name != "" {
			contractData.NameBak = contractData.Name
			isUpdate = true
		}
		if contractData.Name != warnMsg {
			contractData.Name = warnMsg
			isUpdate = true
		}
	} else if params.Status != nil && *params.Status == SensitiveStatusDelete {
		if contractData.NameBak != "" {
			contractData.Name = contractData.NameBak
			contractData.NameBak = ""
			isUpdate = true
		}
	}

	//不需要更新
	if !isUpdate {
		ConvergeDataResponse(ctx, "OK", nil)
	}

	//更新合约数据
	err = UpdateContractSensitiveWord(params.ChainId, contractData)
	if err != nil {
		log.Errorf("UpdateContractSensitiveWord err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	ConvergeDataResponse(ctx, "OK", nil)
}

// UpdateContractSensitiveWord 更新合约敏感词
// @desc
// @param ${param}
// @return error
func UpdateContractSensitiveWord(chainId string, contract *db.Contract) error {
	var err error
	if chainId == "" || contract == nil {
		return nil
	}
	err = dbhandle.UpdateContractNameBak(chainId, contract)
	if err != nil {
		return err
	}
	err = dbhandle.UpdateUpgradeContractName(chainId, contract)
	if err != nil {
		return err
	}
	err = dbhandle.UpdateTransactionContractName(chainId, contract)
	if err != nil {
		return err
	}
	err = dbhandle.UpdateContractEventSensitiveWord(chainId, contract)
	if err != nil {
		return err
	}
	err = dbhandle.UpdateTransferContractName(chainId, contract.Name, contract.Addr)
	if err != nil {
		return err
	}
	err = dbhandle.UpdateNonTransferContractName(chainId, contract.Name, contract.Addr)
	if err != nil {
		return err
	}
	err = dbhandle.UpdateFungibleContractName(chainId, contract)
	if err != nil {
		return err
	}
	err = dbhandle.UpdateNonFungibleContractName(chainId, contract)
	if err != nil {
		return err
	}
	//err = dbhandle.UpdateEvidenceContractName(chainId, contract.Name, contract.Addr)
	//if err != nil {
	//	return err
	//}
	err = dbhandle.UpdateIdentityContractName(chainId, contract.Name, contract.Addr)
	if err != nil {
		return err
	}

	_ = dbhandle.UpdatePositionContractName(chainId, contract.Name, contract.Addr)
	_ = dbhandle.UpdateNonPositionContractName(chainId, contract.Name, contract.Addr)
	_ = dbhandle.UpdateTokenContractName(chainId, contract.Name, contract.Addr)

	return err
}
