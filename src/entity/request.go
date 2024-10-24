/*
Package entity comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package entity

import "chainmaker_web/src/db"

// nolint
const (
	OffsetDefault = 0
	OffsetMin     = 0
	LimitDefault  = 10
	LimitMax      = 200
	LimitMaxSpec  = 100000
	LimitMin      = 0

	Project = "chainmaker"
	CMB     = "cmb"

	GetOverviewData = "GetOverviewData"
	Decimal         = "Decimal"
	GetChainConfig  = "GetChainConfig"
	GetChainList    = "GetChainList"
	Search          = "Search"

	GetNodeList      = "GetNodeList"
	GetOrgList       = "GetOrgList"
	GetAccountList   = "GetAccountList"
	GetAccountDetail = "GetAccountDetail"

	GetTxDetail       = "GetTxDetail"
	GetLatestTxList   = "GetLatestTxList"
	GetTxList         = "GetTxList"
	GetContractTxList = "GetContractTxList"
	GetBlockTxList    = "GetBlockTxList"
	GetUserTxList     = "GetUserTxList"
	GetTxNumByTime    = "GetTxNumByTime"

	GetContractVersionList = "GetContractVersionList"
	GetUserList            = "GetUserList"

	GetBlockDetail       = "GetBlockDetail"
	GetLatestBlockList   = "GetLatestBlockList"
	GetBlockList         = "GetBlockList"
	GetContractEventList = "GetContractEventList"
	GetContractCode      = "GetContractCode"

	GetLatestContractList = "GetLatestContractList"
	GetContractDetail     = "GetContractDetail"
	GetContractList       = "GetContractList"
	GetEventList          = "GetEventList"

	GetFTContractList       = "GetFTContractList"
	GetFTContractDetail     = "GetFTContractDetail"
	GetNFTContractList      = "GetNFTContractList"
	GetNFTContractDetail    = "GetNFTContractDetail"
	GetEvidenceContractList = "GetEvidenceContractList"
	GetEvidenceContract     = "GetEvidenceContract"
	GetIdentityContractList = "GetIdentityContractList"
	GetIdentityContract     = "GetIdentityContract"

	GetFTTransferList  = "GetFTTransferList"
	GetNFTTransferList = "GetNFTTransferList"
	GetNFTList         = "GetNFTList"
	GetNFTDetail       = "GetNFTDetail"

	GetFTPositionList     = "GetFTPositionList"
	GetUserFTPositionList = "GetUserFTPositionList"
	GetNFTPositionList    = "GetNFTPositionList"

	GetGasRecordList = "GetGasRecordList"
	GetGasList       = "GetGasList"
	GetGasInfo       = "GetGasInfo"

	// 订阅接口访问
	SubscribeChain  = "SubscribeChain"
	ModifySubscribe = "ModifySubscribe"
	DeleteSubscribe = "DeleteSubscribe"
	CancelSubscribe = "CancelSubscribe"

	//ModifyTxBlackList 更新操作
	ModifyTxBlackList = "ModifyTxBlackList"
	ModifyUserStatus  = "ModifyUserStatus"

	UpdateTxSensitiveWord           = "UpdateTxSensitiveWord"
	UpdateEventSensitiveWord        = "UpdateEventSensitiveWord"
	UpdateEvidenceSensitiveWord     = "UpdateEvidenceSensitiveWord"
	RecoverEvidenceSensitiveWord    = "RecoverEvidenceSensitiveWord"
	UpdateNFTSensitiveWord          = "UpdateNFTSensitiveWord"
	UpdateContractNameSensitiveWord = "UpdateContractNameSensitiveWord"

	//数据要素相关
	//GetIDAContractList 获取IDA合约列表
	GetIDAContractList = "GetIDAContractList"
	//GetIDADataList ida资产列表
	GetIDADataList = "GetIDADataList"
	//GetIDADataDetail IDA资产详情
	GetIDADataDetail = "GetIDADataDetail"
)

// RequestBody body
type RequestBody interface {
	// IsLegal 是否合法
	IsLegal() bool
}

// ChainBody chain
type ChainBody struct {
	ChainId string
}

// IsLegal legal
func (chainBody *ChainBody) IsLegal() bool {
	// 不为空即合法
	return chainBody.ChainId != ""
}

// RangeBody range
type RangeBody struct {
	Offset int
	Limit  int
}

// IsLegal legal
func (rangeBody *RangeBody) IsLegal() bool {
	if rangeBody.Limit > LimitMax || rangeBody.Offset < OffsetMin {
		return false
	}
	return true
}

// NewRangeBody new
func NewRangeBody() *RangeBody {
	return &RangeBody{
		Offset: OffsetDefault,
		Limit:  LimitDefault,
	}
}

// GetChainIdParams get
type GetChainIdParams struct {
	ChainId string
}

// IsLegal legal
func (params *GetChainIdParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetUserListParams get
type GetUserListParams struct {
	ChainId   string
	UserIds   string
	UserAddrs string
	OrgId     string
	RangeBody
}

// IsLegal legal
func (params *GetUserListParams) IsLegal() bool {
	if params.Limit > LimitMaxSpec || params.Offset < OffsetMin || params.Limit < LimitMin {
		return false
	}
	return len(params.ChainId) > 0
}

// GetAccountListParams get
type GetAccountListParams struct {
	ChainId  string
	AddrType *int
	AddrList []string
	RangeBody
}

// IsLegal legal
func (params *GetAccountListParams) IsLegal() bool {
	if params.Limit > LimitMaxSpec || params.Offset < OffsetMin || params.Limit < LimitMin {
		return false
	}
	return len(params.ChainId) > 0
}

// GetAccountDetailParams get
type GetAccountDetailParams struct {
	ChainId string
	Address string
	BNS     string
}

// IsLegal legal
func (params *GetAccountDetailParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetTransactionNumByTimeParams get
type GetTransactionNumByTimeParams struct {
	ChainId   string
	SortType  int
	StartTime int64
	EndTime   int64
	Interval  int64
}

// IsLegal legal
func (params *GetTransactionNumByTimeParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetEventListParams get
type GetEventListParams struct {
	RangeBody
	ChainId      string
	ContractName string
	ContractAddr string
	TxId         string
}

// IsLegal legal
func (params *GetEventListParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetLatestContractParams get
type GetLatestContractParams struct {
	ChainId string
}

// IsLegal legal
func (params *GetLatestContractParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetLatestChainParams get
type GetLatestChainParams struct {
	Number int
}

// IsLegal legal
func (params *GetLatestChainParams) IsLegal() bool {
	return params.Number > 0
}

// GetLatestBlockParams get
type GetLatestBlockParams struct {
	ChainId string
}

// IsLegal legal
func (params *GetLatestBlockParams) IsLegal() bool {
	return true
}

// GetBlockDetailParams get
type GetBlockDetailParams struct {
	ChainId     string
	BlockHash   string
	BlockHeight *int64
}

// IsLegal legal
func (params *GetBlockDetailParams) IsLegal() bool {
	if params.BlockHash == "" && params.BlockHeight == nil {
		return false
	}
	return params.ChainId != ""

}

// GetContractDetailParams get
type GetContractDetailParams struct {
	ChainId     string
	ContractKey string
}

// IsLegal legal
func (params *GetContractDetailParams) IsLegal() bool {
	return len(params.ChainId) > 0 && len(params.ContractKey) > 0
}

// GetContractCodeParams get
type GetContractCodeParams struct {
	ChainId      string
	ContractName string
}

// IsLegal legal
func (params *GetContractCodeParams) IsLegal() bool {
	if params.ContractName != "" && params.ChainId != "" {
		return true
	}
	return false

}

// GetContractListParams get
type GetContractListParams struct {
	RangeBody
	ChainId      string
	Creators     string
	CreatorAddrs string
	Upgraders    string
	UpgradeAddrs string
	RuntimeType  string
	ContractKey  string
	Status       *int32 //合约状态 -1:全部合约 0：正常 1：已冻结 2：已注销
	StartTime    int64
	EndTime      int64
}

// IsLegal legal
func (params *GetContractListParams) IsLegal() bool {
	return params.ChainId != ""

}

// GetFTContractListParams get
type GetFTContractListParams struct {
	RangeBody
	ChainId     string
	ContractKey string
}

// IsLegal legal
func (params *GetFTContractListParams) IsLegal() bool {
	return params.ChainId != ""

}

// GetFungibleContractParams get
type GetFungibleContractParams struct {
	ChainId      string
	ContractAddr string
}

// IsLegal legal
func (params *GetFungibleContractParams) IsLegal() bool {
	return len(params.ChainId) > 0 && len(params.ContractAddr) > 0
}

// GetNFTContractListParams get
type GetNFTContractListParams struct {
	RangeBody
	ChainId     string
	ContractKey string
}

// IsLegal legal
func (params *GetNFTContractListParams) IsLegal() bool {
	return params.ChainId != ""

}

// GetNonFungibleContractParams get
type GetNonFungibleContractParams struct {
	ChainId      string
	ContractAddr string
}

// IsLegal legal
func (params *GetNonFungibleContractParams) IsLegal() bool {
	return len(params.ChainId) > 0 && len(params.ContractAddr) > 0
}

// GetEvidenceContractListParams get
type GetEvidenceContractListParams struct {
	RangeBody
	ChainId     string
	ContractKey string
}

// IsLegal legal
func (params *GetEvidenceContractListParams) IsLegal() bool {
	return params.ChainId != ""

}

// GetEvidenceContractParams get
type GetEvidenceContractParams struct {
	RangeBody
	ChainId      string
	ContractName string
	TxId         string
	SenderAddrs  string
	Hashs        string
	Code         int
	Search       string
}

// IsLegal legal
func (params *GetEvidenceContractParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetIdentityContractListParams get
type GetIdentityContractListParams struct {
	RangeBody
	ChainId     string
	ContractKey string
}

// IsLegal legal
func (params *GetIdentityContractListParams) IsLegal() bool {
	return params.ChainId != ""

}

// GetIdentityContractParams get
type GetIdentityContractParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	UserAddrs    []string
}

// IsLegal legal
func (params *GetIdentityContractParams) IsLegal() bool {
	return len(params.ChainId) > 0 && len(params.ContractAddr) > 0
}

// ChainOverviewDataParams param
type ChainOverviewDataParams struct {
	ChainId string
}

// IsLegal legal
func (params *ChainOverviewDataParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetChainListParams get
type GetChainListParams struct {
	ChainId   string
	StartTime int64
	EndTime   int64
	RangeBody
}

// IsLegal legal
func (params *GetChainListParams) IsLegal() bool {
	return true
}

// GetOrgListParams get
type GetOrgListParams struct {
	ChainId string
	OrgId   string
	RangeBody
}

// IsLegal legal
func (params *GetOrgListParams) IsLegal() bool {
	if params.Limit > LimitMaxSpec || params.Offset < OffsetMin {
		return false
	}
	return params.ChainId != ""
}

// GetTxDetailParams get
type GetTxDetailParams struct {
	ChainId string
	TxId    string
}

// IsLegal legal
func (params *GetTxDetailParams) IsLegal() bool {
	return len(params.TxId) > 0
}

// GetBlockLatestListParams get
type GetBlockLatestListParams struct {
	ChainId string
}

// IsLegal legal
func (params *GetBlockLatestListParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetTxLatestListParams get
type GetTxLatestListParams struct {
	ChainId string
}

// IsLegal legal
func (params *GetTxLatestListParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetContractVersionListParams get
type GetContractVersionListParams struct {
	RangeBody
	ChainId      string
	ContractName string
	ContractAddr string
	Senders      string
	Status       *int
	RuntimeType  string
	StartTime    int64
	EndTime      int64
}

// IsLegal legal
func (params *GetContractVersionListParams) IsLegal() bool {
	if params.Status == nil {
		defaultValue := -1
		params.Status = &defaultValue
	}

	return params.ChainId != ""
}

// GetContractTxListParams get
type GetContractTxListParams struct {
	RangeBody
	ChainId        string
	ContractName   string
	ContractAddr   string
	UserAddrs      string
	ContractMethod string
	TxStatus       *int
}

// IsLegal legal
func (params *GetContractTxListParams) IsLegal() bool {
	if params.TxStatus == nil {
		defaultValue := -1
		params.TxStatus = &defaultValue
	}

	return params.ChainId != ""
}

// InnerGetContractVersionListParams get
type InnerGetContractVersionListParams struct {
	RangeBody
	ChainId      string
	ContractName string
	Creator      string
	RuntimeType  string
	Status       int64
	StartTime    int64
	EndTime      int64
}

// IsLegal legal
func (params *InnerGetContractVersionListParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// InnerGetTxListParams get
type InnerGetTxListParams struct {
	RangeBody
	ChainId      string
	ContractName string
	Creator      string
	TxId         string
	TxStatus     int
	StartTime    int64
	EndTime      int64
	UserAddr     string
	UserAddrs    []string
}

// IsLegal legal
func (params *InnerGetTxListParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// InnerModifyTxShowStatusParams inner
type InnerModifyTxShowStatusParams struct {
	TxId    string
	ChainId string
}

// IsLegal legal
func (params *InnerModifyTxShowStatusParams) IsLegal() bool {
	return len(params.TxId) > 0
}

// GetBlockListParams param
type GetBlockListParams struct {
	RangeBody
	ChainId   string
	BlockKey  string // Height or block Hash
	StartTime int64
	EndTime   int64
}

// IsLegal legal
func (params *GetBlockListParams) IsLegal() bool {
	if params.Limit > LimitMax || params.Offset < OffsetMin {
		return false
	}
	return params.ChainId != ""
}

// GetTxListParams param
type GetTxListParams struct {
	RangeBody
	TxId           string
	ChainId        string
	ContractName   string
	ContractAddr   string
	ContractMethod string
	BlockHash      string
	UserAddrs      string
	Senders        string
	TxStatus       *int
	StartTime      int64
	EndTime        int64
}

// IsLegal legal
func (params *GetTxListParams) IsLegal() bool {
	if params.TxStatus == nil {
		//默认全部交易
		defaultValue := -1
		params.TxStatus = &defaultValue
	}

	if params.Limit > LimitMaxSpec || params.Offset < OffsetMin || params.Limit < LimitMin {
		return false
	}
	return params.ChainId != ""
}

// GetBlockTxListParams param
type GetBlockTxListParams struct {
	RangeBody
	ChainId   string
	BlockHash string
}

// IsLegal legal
func (params *GetBlockTxListParams) IsLegal() bool {
	if params.Limit > LimitMax || params.Offset < OffsetMin || params.Limit < LimitMin {
		return false
	}
	return params.ChainId != "" && params.BlockHash != ""
}

// GetUserTxListParams param
type GetUserTxListParams struct {
	RangeBody
	ChainId   string
	UserAddrs string
}

// IsLegal legal
func (params *GetUserTxListParams) IsLegal() bool {
	if params.Limit > LimitMax || params.Offset < OffsetMin || params.Limit < LimitMin {
		return false
	}
	return params.ChainId != "" && params.UserAddrs != ""
}

// SearchParams search
type SearchParams struct {
	Type    string
	Value   string
	ChainId string
}

// IsLegal legal
func (params *SearchParams) IsLegal() bool {
	if params.Value == "" || params.ChainId == "" {
		return false
	}
	return true
}

// GetDetailParams get
type GetDetailParams struct {
	Id      string
	ChainId string
}

// IsLegal legal
func (params *GetDetailParams) IsLegal() bool {
	return true
}

// ChainNodesParams param
type ChainNodesParams struct {
	RangeBody
	ChainId  string
	NodeName string
	OrgId    string
	NodeId   string
}

// IsLegal legal
func (params *ChainNodesParams) IsLegal() bool {
	if params.Limit > LimitMaxSpec || params.Offset < OffsetMin || params.Limit < LimitMin {
		return false
	}
	return params.ChainId != ""
}

// SubscribeChainParams 订阅链相关
type SubscribeChainParams struct {
	ChainId     string
	OrgId       string
	Addr        string
	OrgCA       string
	Tls         bool
	UserCert    string
	UserKey     string
	TLSHostName string
	AuthType    string
	HashType    int
	NodeList    []db.SubscribeNode
}

// IsLegal legal
func (params *SubscribeChainParams) IsLegal() bool {
	if params.ChainId == "" || params.HashType < 0 {
		return false
	}
	return true
}

// CancelSubscribeParams cancel
type CancelSubscribeParams struct {
	ChainId string
	Status  int
}

// IsLegal legal
func (params *CancelSubscribeParams) IsLegal() bool {
	return params.ChainId != ""
}

// ModifySubscribeParams modify
type ModifySubscribeParams struct {
	ChainId     string
	OrgId       string
	Addr        string
	Tls         bool
	OrgCA       string
	UserCert    string
	UserKey     string
	TLSHostName string
	AuthType    string
	HashType    int
	NodeList    []db.SubscribeNode
}

// IsLegal legal
func (params *ModifySubscribeParams) IsLegal() bool {
	if params.ChainId == "" || len(params.NodeList) <= 0 {
		return false
	}
	return true
}

// DeleteSubscribeParams delete
type DeleteSubscribeParams struct {
	ChainId string
}

// IsLegal legal
func (params *DeleteSubscribeParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetGasListParams get
type GetGasListParams struct {
	ChainId   string
	UserAddrs string
	RangeBody
}

// IsLegal legal
func (params *GetGasListParams) IsLegal() bool {
	if params.Limit > LimitMaxSpec || params.Offset < OffsetMin || params.Limit < LimitMin {
		return false
	}
	return len(params.ChainId) > 0
}

// GetGasRecordListParams get
type GetGasRecordListParams struct {
	ChainId      string
	BusinessType int
	UserAddrs    string
	StartTime    int64
	EndTime      int64
	RangeBody
}

// IsLegal legal
func (params *GetGasRecordListParams) IsLegal() bool {
	if params.Limit > LimitMaxSpec || params.Offset < OffsetMin || params.Limit < LimitMin {
		return false
	}
	return len(params.ChainId) > 0
}

// GetGasInfoParams get
type GetGasInfoParams struct {
	ChainId   string
	UserAddrs string
}

// IsLegal legal
func (params *GetGasInfoParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// InnerGetChainInfoParams get
type InnerGetChainInfoParams struct {
	ChainId   string
	Address   string
	Addresses []string
}

// IsLegal legal
func (params *InnerGetChainInfoParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// ModifyUserStatusParams get
type ModifyUserStatusParams struct {
	ChainId string
	Address string
	Status  int
}

// IsLegal legal
func (params *ModifyUserStatusParams) IsLegal() bool {
	return len(params.ChainId) > 0 && len(params.Address) > 0
}

// InnerEvidenceListParams get
type InnerEvidenceListParams struct {
	ChainId     string
	SenderAddr  string
	TxId        string
	Code        int
	Hash        string
	SenderAddrs []string
	RangeBody
}

// IsLegal legal
func (params *InnerEvidenceListParams) IsLegal() bool {
	if params.Limit > LimitMaxSpec || params.Offset < OffsetMin || params.Limit < LimitMin {
		return false
	}
	return len(params.ChainId) > 0
}

// GetNFTListParams get
type GetNFTListParams struct {
	RangeBody
	ChainId     string
	TokenId     string
	ContractKey string
	OwnerAddrs  string
}

// IsLegal legal
func (params *GetNFTListParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetNFTDetailParams get
type GetNFTDetailParams struct {
	ChainId      string
	TokenId      string
	ContractAddr string
}

// IsLegal legal
func (params *GetNFTDetailParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetTransferListParams get
type GetTransferListParams struct {
	RangeBody
	ChainId      string
	TokenId      string
	ContractName string
}

// IsLegal legal
func (params *GetTransferListParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetFungibleTransferListParams get
type GetFungibleTransferListParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	UserAddr     string
}

// IsLegal legal
func (params *GetFungibleTransferListParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetNonFungibleTransferListParams get
type GetNonFungibleTransferListParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	UserAddr     string
	TokenId      string
}

// IsLegal legal
func (params *GetNonFungibleTransferListParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetFungiblePositionListParams get
type GetFungiblePositionListParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	OwnerAddr    string
}

// IsLegal legal
func (params *GetFungiblePositionListParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetUserFTPositionListParams get
type GetUserFTPositionListParams struct {
	RangeBody
	ChainId   string
	OwnerAddr string
}

// IsLegal legal
func (params *GetUserFTPositionListParams) IsLegal() bool {
	return params.ChainId != "" && params.OwnerAddr != ""
}

// GetNonFungiblePositionListParams get
type GetNonFungiblePositionListParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	OwnerAddr    string
}

// IsLegal legal
func (params *GetNonFungiblePositionListParams) IsLegal() bool {
	return params.ChainId != "" && params.ContractAddr != ""
}

// ModifyTxBlackListParams update
type ModifyTxBlackListParams struct {
	ChainId string
	TxId    string
	Status  *int
}

// IsLegal legal
func (params *ModifyTxBlackListParams) IsLegal() bool {
	if params.TxId == "" {
		return false
	}
	return params.ChainId != ""
}

// DeleteTxBlackListParams update
type DeleteTxBlackListParams struct {
	ChainId string
	TxId    string
}

// IsLegal legal
func (params *DeleteTxBlackListParams) IsLegal() bool {
	if params.TxId == "" {
		return false
	}
	return params.ChainId != ""
}

// UpdateTxSensitiveWordParams update
type UpdateTxSensitiveWordParams struct {
	ChainId string
	TxId    string
	Column  string
	Status  int
	WarnMsg string
}

// IsLegal legal
func (params *UpdateTxSensitiveWordParams) IsLegal() bool {
	if params.ChainId == "" || params.TxId == "" {
		return false
	}
	return true
}

// UpdateEventSensitiveWordParams update
type UpdateEventSensitiveWordParams struct {
	ChainId string
	TxId    string
	Index   int
	Column  string
	Status  *int
	WarnMsg string
}

// IsLegal legal
func (params *UpdateEventSensitiveWordParams) IsLegal() bool {
	if params.Column == "" || params.ChainId == "" || params.TxId == "" || params.Status == nil {
		return false
	}
	return true
}

// EvidenceSensitiveWordParams update
type EvidenceSensitiveWordParams struct {
	ChainId string
	Hash    string
	Column  string
	Status  *int
	WarnMsg string
}

// IsLegal legal
func (params *EvidenceSensitiveWordParams) IsLegal() bool {
	if params.Column == "" || params.ChainId == "" || params.Hash == "" || params.Status == nil {
		return false
	}
	return true
}

// NFTSensitiveWordParams update
type NFTSensitiveWordParams struct {
	ChainId      string
	TokenId      string
	ContractAddr string
	Column       string
	Status       *int
	WarnMsg      string
}

// IsLegal legal
func (params *NFTSensitiveWordParams) IsLegal() bool {
	if params.Column == "" || params.ChainId == "" || params.TokenId == "" || params.Status == nil {
		return false
	}
	return true
}

// UpdateContractNameSWParams update
type UpdateContractNameSWParams struct {
	ChainId      string
	ContractName string
	Status       *int
	WarnMsg      string
}

// IsLegal legal
func (params *UpdateContractNameSWParams) IsLegal() bool {
	if params.ChainId == "" || params.ContractName == "" || params.Status == nil {
		return false
	}
	return true
}

// GetIDAContractListParams ida合约列表
type GetIDAContractListParams struct {
	ChainId     string
	ContractKey string
	RangeBody
}

// IsLegal legal
func (params *GetIDAContractListParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetIDADataListParams ida资产列表
type GetIDADataListParams struct {
	ChainId      string
	AssetCode    string
	ContractAddr string
	RangeBody
}

// IsLegal legal
func (params *GetIDADataListParams) IsLegal() bool {
	return params.ChainId != "" && params.ContractAddr != ""
}

// GetIDADataDetailParams ida资产详情
type GetIDADataDetailParams struct {
	ChainId      string
	AssetCode    string
	ContractAddr string
	RangeBody
}

// IsLegal legal
func (params *GetIDADataDetailParams) IsLegal() bool {
	return params.ChainId != "" && params.AssetCode != "" && params.ContractAddr != ""
}
