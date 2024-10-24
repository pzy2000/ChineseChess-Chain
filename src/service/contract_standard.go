package service

import (
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"strings"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
)

// GetFTContractListHandler handler
type GetFTContractListHandler struct{}

// Handle 同质化合约列表
func (getFTContractListHandler *GetFTContractListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetFungibleContractListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetFungibleContractList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	contractList, count, err := dbhandle.GetFungibleContractList(params.Offset, params.Limit, params.ChainId,
		params.ContractKey)
	if err != nil {
		log.Errorf("GetContractList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	contractListView := arraylist.New()
	for _, contract := range contractList {
		contractView := &entity.FungibleContractListView{
			TxNum:          contract.TxNum,
			ContractName:   contract.ContractName,
			ContractSymbol: contract.Symbol,
			ContractAddr:   contract.ContractAddr,
			ContractType:   contract.ContractType,
			TotalSupply:    contract.TotalSupply.String(),
			HolderCount:    contract.HolderCount,
			Timestamp:      contract.Timestamp,
		}
		contractListView.Add(contractView)
	}
	ConvergeListResponse(ctx, contractListView.Values(), count, nil)
}

// GetFTContractDetailHandler handler
type GetFTContractDetailHandler struct{}

// Handle 同质化合约详情
func (GetFTContractDetailHandler *GetFTContractDetailHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetFungibleContractHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetContractDetail param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取合约
	contract, err := dbhandle.GetContractByAddr(params.ChainId, params.ContractAddr)
	if err != nil || contract == nil {
		log.Errorf("GetFungibleContract err : %s", err)
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	fungibleContract, err := dbhandle.GetFungibleContractByAddr(params.ChainId, params.ContractAddr)
	if err != nil || fungibleContract == nil {
		log.Errorf("GetFungibleContract err : %s", err)
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	fungibleContractView := &entity.FTContractDetailView{
		Version:         contract.Version,
		ContractName:    contract.Name,
		ContractSymbol:  contract.ContractSymbol,
		ContractAddr:    contract.Addr,
		ContractStatus:  contract.ContractStatus,
		TxId:            contract.CreateTxId,
		CreateSender:    contract.CreatorAddr,
		CreatorAddr:     contract.CreatorAddr,
		CreatorAddrBNS:  contract.CreatorAddr,
		RuntimeType:     contract.RuntimeType,
		ContractType:    contract.ContractType,
		CreateTimestamp: contract.Timestamp,
		UpdateTimestamp: contract.UpgradeTimestamp,
		TotalSupply:     fungibleContract.TotalSupply.String(),
		HolderCount:     fungibleContract.HolderCount,
		TxNum:           contract.TxNum,
	}
	ConvergeDataResponse(ctx, fungibleContractView, nil)
}

// GetNFTContractListHandler handler
type GetNFTContractListHandler struct{}

// Handle 非同质化合约列表
func (getNFTContractListHandler *GetNFTContractListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetNonFungibleContractListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetNonFungibleContractList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	contractList, count, err := dbhandle.GetNonFungibleContractList(params.Offset, params.Limit, params.ChainId,
		params.ContractKey)
	if err != nil {
		log.Errorf("GetContractList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	contractListView := arraylist.New()
	for _, contract := range contractList {
		contractView := &entity.NonFungibleContractListView{
			TxNum:        contract.TxNum,
			ContractName: contract.ContractName,
			ContractAddr: contract.ContractAddr,
			ContractType: contract.ContractType,
			TotalSupply:  contract.TotalSupply.String(),
			HolderCount:  contract.HolderCount,
			Timestamp:    contract.Timestamp,
		}
		contractListView.Add(contractView)
	}
	ConvergeListResponse(ctx, contractListView.Values(), count, nil)
}

// GetNFTContractDetailHandler handler
type GetNFTContractDetailHandler struct{}

// Handle 非同质化合约详情
func (getNFTContractDetailHandler *GetNFTContractDetailHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetNonFungibleContractHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetNonFungibleContract param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取合约
	contract, err := dbhandle.GetContractByAddr(params.ChainId, params.ContractAddr)
	if err != nil || contract == nil {
		log.Errorf("GetNonFungibleContract err : %s", err)
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	nonFungibleContract, err := dbhandle.GetNonFungibleContractByAddr(params.ChainId, params.ContractAddr)
	if err != nil || nonFungibleContract == nil {
		log.Errorf("GetNonFungibleContract err : %s", err)
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	//获取地址BNS
	creatorAddrBNS := GetAccountBNSByAddr(params.ChainId, contract.CreatorAddr)

	fungibleContractView := &entity.NFTContractDetailView{
		Version:         contract.Version,
		ContractName:    contract.Name,
		ContractAddr:    contract.Addr,
		ContractStatus:  contract.ContractStatus,
		TxId:            contract.CreateTxId,
		CreateSender:    contract.CreateSender,
		CreatorAddr:     contract.CreatorAddr,
		CreatorAddrBNS:  creatorAddrBNS,
		RuntimeType:     contract.RuntimeType,
		ContractType:    contract.ContractType,
		TotalSupply:     nonFungibleContract.TotalSupply.String(),
		HolderCount:     nonFungibleContract.HolderCount,
		CreateTimestamp: contract.Timestamp,
		UpdateTimestamp: contract.UpgradeTimestamp,
		TxNum:           contract.TxNum,
	}
	ConvergeDataResponse(ctx, fungibleContractView, nil)
}

// GetEvidenceContractHandler handler
type GetEvidenceContractHandler struct{}

// Handle 存证合约详情
func (getEvidenceContractHandler *GetEvidenceContractHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetEvidenceContractHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetContractDetail param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}
	var hashList []string
	var senderAddrs []string
	if params.Hashs != "" {
		hashList = strings.Split(params.Hashs, ",")
	}
	if params.SenderAddrs != "" {
		senderAddrs = strings.Split(params.SenderAddrs, ",")
	}
	//获取合约
	contractList, count, err := dbhandle.GetEvidenceContract(params.Offset, params.Limit, params.Code, params.ChainId,
		params.ContractName, params.TxId, params.Search, hashList, senderAddrs)
	if err != nil {
		log.Errorf("GetEvidenceContract err : %s", err)
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	contractListView := arraylist.New()
	for _, contract := range contractList {
		contractView := &entity.EvidenceContractView{
			Id:           contract.ID,
			ChainId:      params.ChainId,
			ContractName: contract.ContractName,
			TxId:         contract.TxId,
			BlockHeight:  contract.BlockHeight,
			EvidenceId:   contract.EvidenceId,
			SenderAddr:   contract.SenderAddr,
			Hash:         contract.Hash,
			MetaData:     contract.MetaData,
			Code:         int(contract.ContractResultCode + 1),
			ResultCode:   contract.ContractResultCode,
			Timestamp:    contract.Timestamp,
		}
		contractListView.Add(contractView)
	}
	ConvergeListResponse(ctx, contractListView.Values(), count, nil)
}

// SaveEvidenceMetaData 处理EvidenceMetaData
//func dealEvidenceMetaData(metaData string) []entity.EvidenceMetaData {
//	evidenceDataList := make([]entity.EvidenceMetaData, 0)
//	if metaData == "" {
//		return evidenceDataList
//	}
//
//	// 解码 JSON 字符串到一个 map[string]interface{} 变量中
//	var jsonData map[string]interface{}
//	err := json.Unmarshal([]byte(metaData), &jsonData)
//	if err != nil {
//		log.Errorf("dealEvidenceMetaData error decoding JSON :%v", err.Error())
//		return evidenceDataList
//	}
//
//	for key, value := range jsonData {
//		// 将值转换为 JSON 字符串（如果需要）
//		var valueStr string
//		switch v := value.(type) {
//		case string, bool:
//			valueStr = fmt.Sprintf("%v", v)
//		case int, int32, int64, float64, float32:
//			valueStr = fmt.Sprintf("%v", v)
//		default:
//			jsonBytes, err := json.Marshal(v)
//			if err != nil {
//				log.Errorf("dealEvidenceMetaData Error encoding value to JSON::%v", err)
//				valueStr = fmt.Sprintf("%v", v)
//			} else {
//				valueStr = string(jsonBytes)
//			}
//		}
//
//		evidenceData := entity.EvidenceMetaData{
//			Key:   key,
//			Value: valueStr,
//		}
//		evidenceDataList = append(evidenceDataList, evidenceData)
//	}
//
//	return evidenceDataList
//}

//// GetIdentityContractListHandler handler
//type GetIdentityContractListHandler struct{}

//// Handle 身份认证合约列表
//func (getIdentityContractListHandler *GetIdentityContractListHandler) Handle(ctx *gin.Context) {
//	params := entity.BindGetIdentityContractListHandler(ctx)
//	if params == nil || !params.IsLegal() {
//		newError := entity.NewError(entity.ErrorParamWrong, " GetEvidenceContractList param is wrong")
//		ConvergeFailureResponse(ctx, newError)
//		return
//	}
//
//	contractList, count, err := dbhandle.GetIdentityContractList(params.Offset, params.Limit, params.ChainId,
//	params.ContractKey)
//	if err != nil {
//		log.Errorf("GetIdentityContractList err : %s", err.Error())
//		ConvergeHandleFailureResponse(ctx, err)
//		return
//	}
//
//	contractListView := arraylist.New()
//	for _, contract := range contractList {
//		contractView := &entity.IdentityContractListView{
//			ContractName: contract.ContractName,
//			ContractAddr: contract.ContractAddr,
//			Timestamp:    contract.CreatedAt,
//		}
//		contractListView.Add(contractView)
//	}
//	ConvergeListResponse(ctx, contractListView.Values(), count, nil)
//}

// GetIdentityContractHandler handler
type GetIdentityContractHandler struct{}

// Handle 身份认证合约详情
func (getIdentityContractHandler *GetIdentityContractHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetIdentityContractHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetIdentityContract param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取合约
	contractList, count, err := dbhandle.GetIdentityContract(params.Offset, params.Limit, params.ChainId,
		params.ContractAddr, params.UserAddrs)
	if err != nil {
		log.Errorf("GetEvidenceContract err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	contractListView := arraylist.New()
	for _, contract := range contractList {
		contractView := &entity.IdentityContractView{
			ContractName: contract.ContractName,
			ContractAddr: contract.ContractAddr,
			UserAddr:     contract.UserAddr,
			Level:        contract.Level,
			PkPem:        contract.PkPem,
		}
		contractListView.Add(contractView)
	}
	ConvergeListResponse(ctx, contractListView.Values(), count, nil)
}
