package service

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"encoding/json"
	"strings"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
)

const (
	//UrlTypeImage token类型为图片
	UrlTypeImage = "image"
	//UrlTypeVideo token类型为视频
	UrlTypeVideo = "video"
)

// GetFungibleTransferListHandler get
type GetFungibleTransferListHandler struct {
}

// Handle deal
func (getFungibleTransferListHandler *GetFungibleTransferListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetFungibleTransferListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetFungibleTransferList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	fungibleContract, err := dbhandle.GetFungibleContractByAddr(params.ChainId, params.ContractAddr)
	if err != nil || fungibleContract == nil {
		log.Error("getFungibleTransferListHandler GetFungibleContract err : %s", err)
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	//统计同质化合约流转记录数
	totalCount, err := getFTTransferTotalCount(params.ChainId, params.UserAddr, fungibleContract)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}
	if totalCount == 0 {
		ConvergeListResponse(ctx, []interface{}{}, totalCount, nil)
		return
	}

	//流转列表
	transferList, err := dbhandle.GetFTTransferList(params.Offset, params.Limit, params.ChainId,
		params.ContractAddr, params.UserAddr)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//获取创建合约账户对应的账户信息
	accountMap := GetFungibleTransferAccountMap(params.ChainId, transferList)
	views := arraylist.New()
	for _, transfer := range transferList {
		//获取地址BNS
		fromAddrBNS := GetAccountBNS(transfer.FromAddr, accountMap)
		toAddrBNS := GetAccountBNS(transfer.ToAddr, accountMap)

		transferView := &entity.FungibleTransferListView{
			TxId:           transfer.TxId,
			ContractName:   transfer.ContractName,
			ContractAddr:   transfer.ContractAddr,
			ContractMethod: transfer.ContractMethod,
			ContractSymbol: fungibleContract.Symbol,
			From:           transfer.FromAddr,
			FromBNS:        fromAddrBNS,
			To:             transfer.ToAddr,
			ToBNS:          toAddrBNS,
			Amount:         transfer.Amount.String(),
			Timestamp:      transfer.Timestamp,
		}
		views.Add(transferView)
	}

	ConvergeListResponse(ctx, views.Values(), totalCount, nil)
}

// getNFTTransferTotalCount
//
//	@Description: 统计非同质化合约流转记录数
//	@param params
//	@return int64
//	@return error
func getFTTransferTotalCount(chainId, userAddr string, ftContract *db.FungibleContract) (int64, error) {
	var contractAddr string
	if ftContract != nil && userAddr == "" {
		return ftContract.TransferNum, nil
	}

	//需要查询count
	if ftContract != nil {
		contractAddr = ftContract.ContractAddr
	}
	totalCount, err := dbhandle.GetFTTransferCount(chainId, contractAddr, userAddr)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

// GetNonFungibleTransferListHandler get
type GetNonFungibleTransferListHandler struct{}

// Handle deal
func (getNonFungibleTransferListHandler *GetNonFungibleTransferListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetNonFungibleTransferListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetFungibleTransferList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//统计非同质化合约流转记录数
	totalCount, err := getNFTTransferTotalCount(params)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}
	if totalCount == 0 {
		ConvergeListResponse(ctx, []interface{}{}, totalCount, nil)
		return
	}

	//流转列表列表
	transferList, err := dbhandle.GetNFTTransferList(params.Offset, params.Limit, params.ChainId, params.ContractAddr,
		params.UserAddr, params.TokenId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	if len(transferList) == 0 {
		ConvergeListResponse(ctx, []interface{}{}, totalCount, nil)
		return
	}

	//获取创建合约账户对应的账户信息
	accountMap := GetNFTTransferAccountMap(params.ChainId, transferList)
	views := arraylist.New()
	for _, transfer := range transferList {
		//获取地址BNS
		fromAddrBNS := GetAccountBNS(transfer.FromAddr, accountMap)
		toAddrBNS := GetAccountBNS(transfer.ToAddr, accountMap)

		transferView := &entity.NonFungibleTransferListView{
			TxId:           transfer.TxId,
			ContractName:   transfer.ContractName,
			ContractAddr:   transfer.ContractAddr,
			ContractMethod: transfer.ContractMethod,
			From:           transfer.FromAddr,
			FromBNS:        fromAddrBNS,
			To:             transfer.ToAddr,
			ToBNS:          toAddrBNS,
			TokenId:        transfer.TokenId,
			Timestamp:      transfer.Timestamp,
		}
		views.Add(transferView)
	}

	ConvergeListResponse(ctx, views.Values(), totalCount, nil)
}

// getNFTTransferTotalCount
//
//	@Description: 统计非同质化合约流转记录数
//	@param params
//	@return int64
//	@return error
func getNFTTransferTotalCount(params *entity.GetNonFungibleTransferListParams) (int64, error) {
	contractAddr := params.ContractAddr
	userAddr := params.UserAddr
	tokenId := params.TokenId
	if contractAddr != "" && userAddr == "" && tokenId == "" {
		//没有搜索条件，使用合约表transferNum
		contractInfo, err := dbhandle.GetNonFungibleContractByAddr(params.ChainId, contractAddr)
		if err != nil || contractInfo == nil {
			return 0, err
		}
		return contractInfo.TransferNum, nil
	}

	//需要查询count
	totalCount, err := dbhandle.GetNFTTransferCount(params.ChainId, contractAddr, userAddr, tokenId)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

// GetNFTListHandler get
type GetNFTListHandler struct{}

// Handle deal
func (getNFTListHandler *GetNFTListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetNFTListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetNFTList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var (
		contractAddr string
		ownerAddrs   []string
		contractInfo *db.NonFungibleContract
		err          error
	)

	if params.OwnerAddrs != "" {
		ownerAddrs = strings.Split(params.OwnerAddrs, ",")
	}

	if params.ContractKey != "" {
		contractInfo, err = dbhandle.GetNFTContractByNameOrAddr(params.ChainId, params.ContractKey)
		if contractInfo == nil || err != nil {
			log.Infof("getNFTListHandler GetContractByCacheOrNameAddr err:%v", err)
			ConvergeListResponse(ctx, []interface{}{}, 0, nil)
			return
		}
		contractAddr = contractInfo.ContractAddr
	}

	totalCount, err := GetNFTTokenTotalCount(params.ChainId, params.TokenId, contractInfo, ownerAddrs)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	if totalCount == 0 {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	//流转详情
	tokenList, err := dbhandle.GetNonFungibleTokenList(params.Offset, params.Limit, params.ChainId,
		params.TokenId, contractAddr, ownerAddrs)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	var accountAddrs []string
	for _, tokenInfo := range tokenList {
		accountAddrs = append(accountAddrs, tokenInfo.OwnerAddr)
	}
	accountMap, _ := dbhandle.QueryAccountExists(params.ChainId, accountAddrs)
	tokenListView := arraylist.New()
	for _, token := range tokenList {
		var addrType int
		var addrBNS string
		if account, ok := accountMap[token.OwnerAddr]; ok {
			addrType = account.AddrType
			addrBNS = account.BNS
		}
		tokenView := &entity.NFTListView{
			AddrType:     addrType,
			TokenId:      token.TokenId,
			Timestamp:    token.Timestamp,
			ContractName: token.ContractName,
			ContractAddr: token.ContractAddr,
			OwnerAddr:    token.OwnerAddr,
			OwnerAddrBNS: addrBNS,
			CategoryName: token.CategoryName,
		}
		//解析metadata
		if token.MetaData != "" {
			metadata := BuildMetadata(token.MetaData)
			if metadata.ImageUrl != "" {
				tokenView.ImageUrl = metadata.ImageUrl
				tokenView.UrlType = isImageOrVideo(metadata.ImageUrl)
			}
		}
		tokenListView.Add(tokenView)
	}
	ConvergeListResponse(ctx, tokenListView.Values(), totalCount, nil)
}

func GetNFTTokenTotalCount(chainId, tokenId string, contractInfo *db.NonFungibleContract, ownerAddrs []string) (int64,
	error) {
	var (
		err          error
		totalCount   int64
		contractAddr string
	)

	//只有合约地址，使用合约发行总量
	if contractInfo != nil && tokenId == "" && len(ownerAddrs) == 0 {
		totalCount = contractInfo.TotalSupply.IntPart()
		return totalCount, err
	} else if contractInfo == nil && tokenId == "" && len(ownerAddrs) > 0 {
		//账户的token列表
		accountMap, errMap := dbhandle.QueryAccountExists(chainId, ownerAddrs)
		if errMap != nil {
			return 0, errMap
		}

		for _, account := range accountMap {
			totalCount += account.NFTNum
		}
		return totalCount, nil
	}

	//否则就需要使用sql进行count计算
	if contractInfo != nil {
		contractAddr = contractInfo.ContractAddr
	}
	totalCount, err = dbhandle.GetNFTTokenCount(chainId, tokenId, contractAddr, ownerAddrs)
	return totalCount, err
}

// GetNFTDetailHandler get
type GetNFTDetailHandler struct{}

// Handle deal
func (getNFTDetailHandler *GetNFTDetailHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetNFTDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//token详情
	tokenInfo, err := dbhandle.GetNonFungibleTokenDetail(params.ChainId, params.TokenId, params.ContractAddr)
	if err != nil || tokenInfo == nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	//获取地址BNS
	accountInfo, _ := dbhandle.GetAccountByAddr(params.ChainId, tokenInfo.OwnerAddr)
	isViolation := true
	tokenView := &entity.NFTDetailView{
		TokenId:      tokenInfo.TokenId,
		Timestamp:    tokenInfo.Timestamp,
		ContractName: tokenInfo.ContractName,
		ContractAddr: tokenInfo.ContractAddr,
		AddrType:     accountInfo.AddrType,
		OwnerAddr:    tokenInfo.OwnerAddr,
		OwnerAddrBNS: accountInfo.BNS,
		CategoryName: tokenInfo.CategoryName,
	}
	//解析metadata
	if tokenInfo.MetaData != "" {
		metadataJson := BuildMetadata(tokenInfo.MetaData)
		if metadataJson.ImageUrl != "" {
			isViolation = false
			urlType := isImageOrVideo(metadataJson.ImageUrl)
			tokenView.Metadata = entity.Metadata{
				Name:        metadataJson.Name,
				Author:      metadataJson.Author,
				OrgName:     metadataJson.OrgName,
				ImageUrl:    metadataJson.ImageUrl,
				Description: metadataJson.Description,
				SeriesHash:  metadataJson.SeriesHash,
				UrlType:     urlType,
			}
		}
	}
	tokenView.IsViolation = isViolation
	ConvergeDataResponse(ctx, tokenView, nil)
}

// GetFungibleTransferAccountMap
//
//	@Description: 根据合约地址获取账户列表
//	@param chainId
//	@param contractList
//	@return map[string]*db.Account
func GetFungibleTransferAccountMap(chainId string, fungibleTransfer []*db.FungibleTransfer) map[string]*db.Account {
	var accountAddrs []string
	for _, transfer := range fungibleTransfer {
		if transfer.FromAddr != "" {
			accountAddrs = append(accountAddrs, transfer.FromAddr)
		}
		if transfer.ToAddr != "" {
			accountAddrs = append(accountAddrs, transfer.ToAddr)
		}
	}
	accountMap, _ := dbhandle.QueryAccountExists(chainId, accountAddrs)
	return accountMap
}

// GetNFTTransferAccountMap
//
//	@Description: 根据合约地址获取账户列表
//	@param chainId
//	@param contractList
//	@return map[string]*db.Account
func GetNFTTransferAccountMap(chainId string, fungibleTransfer []*db.NonFungibleTransfer) map[string]*db.Account {
	var accountAddrs []string
	for _, transfer := range fungibleTransfer {
		if transfer.FromAddr != "" {
			accountAddrs = append(accountAddrs, transfer.FromAddr)
		}
		if transfer.ToAddr != "" {
			accountAddrs = append(accountAddrs, transfer.ToAddr)
		}
	}
	accountMap, _ := dbhandle.QueryAccountExists(chainId, accountAddrs)
	return accountMap
}

// BuildMetadata
//
//	@Description: 解析图片数据
//	@param metaData
//	@return entity.MetadataJson
func BuildMetadata(metaData string) entity.MetadataJson {
	metadataJson := entity.MetadataJson{}
	err := json.Unmarshal([]byte(metaData), &metadataJson)
	if err != nil {
		log.Errorf("合约参数未遵循长安链合约标准协议，解析失败")
	}

	//兼容Image和ImageUrl格式
	if metadataJson.Image != "" {
		metadataJson.ImageUrl = metadataJson.Image
	}
	return metadataJson
}
