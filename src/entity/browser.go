/*
Package entity comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package entity

import "time"

// DecimalView decimal
type DecimalView struct {
	BlockHeight       int64 // 区块高度
	RecentBlockHeight int64
	TxNum             int64 // 交易数量
	RecentTxNum       int64
	ContractNum       int64
	RecentContractNum int64
	UserNum           int64
	RecentUserNum     int64
	OrgNum            int64 // 链数量
	NodeNum           int64
	AuthType          string
}

// ChainIdView view
type ChainIdView struct {
	ChainId  string
	AuthType string
}

// OverviewDataView OverviewData
type OverviewDataView struct {
	ChainId       string
	BlockHeight   int64
	UserCount     int64
	ContractCount int64
	TxCount       int64
	OrgCount      int64
	RunningNode   int64
	CommonNode    int64
	ConsensusNode int64
}

// TransactionNumView tran
type TransactionNumView struct {
	TxNum     int64 // 交易数量
	Timestamp int64 // 时间戳
}

// LatestContractView ContractListView contract
type LatestContractView struct {
	Id               int
	ContractName     string //合约名
	ContractNameBak  string //合约名备份
	ContractAddr     string
	ContractType     string
	Version          string
	Sender           string
	SenderAddr       string
	SenderAddrBNS    string
	UpgradeUser      string
	TxNum            int64
	CreateTimestamp  int64 // 时间戳
	UpgradeTimestamp int64 // 时间戳
	Timestamp        int64 // 时间戳
}

// LatestChainView latest
type LatestChainView struct {
	ChainName string //链名
	Timestamp int64  // 时间戳
}

// LatestBlockView latest
type LatestBlockView struct {
	BlockHash        string
	PreBlockHash     string
	Timestamp        int64
	BlockHeight      int64
	TxCount          int
	ProposalNodeId   string
	ProposalNodeAddr string
}

// BlockDetailView block
// param view
type BlockDetailView struct {
	BlockHash        string
	PreBlockHash     string
	Timestamp        int64
	BlockHeight      int64
	TxCount          int
	ProposalNodeId   string
	RwSetHash        string
	TxRootHash       string
	Dag              string
	OrgId            string
	ProposalNodeAddr string
}

// ContractDetailView view
type ContractDetailView struct {
	ContractName    string
	ContractNameBak string
	ContractAddr    string
	ContractSymbol  string
	ContractType    string
	ContractStatus  int32
	Version         string
	TxId            string
	CreateSender    string
	CreatorAddr     string
	CreatorAddrBns  string
	Timestamp       int64
	DataAssetNum    int64
	RuntimeType     string
}

// NFTContractDetailView view
type NFTContractDetailView struct {
	ContractName    string
	ContractAddr    string
	Version         string
	ContractStatus  int32
	TxId            string
	CreateSender    string
	CreatorAddr     string
	CreatorAddrBNS  string
	CreateTimestamp int64
	UpdateTimestamp int64
	RuntimeType     string
	Status          int
	ContractType    string
	TotalSupply     string
	HolderCount     int64
	TxNum           int64
}

// FTContractDetailView view
type FTContractDetailView struct {
	ContractName    string
	ContractSymbol  string
	ContractAddr    string
	Version         string
	ContractStatus  int32
	TxId            string
	CreateSender    string
	CreatorAddr     string
	CreatorAddrBNS  string
	CreateTimestamp int64
	UpdateTimestamp int64
	RuntimeType     string
	Status          int
	ContractType    string
	TotalSupply     string
	HolderCount     int64
	TxNum           int64
}

// ContractEventView view
type ContractEventView struct {
	EventInfo string
	Topic     string
	Timestamp int64
}

// ContractCodeView view
type ContractCodeView struct {
	ContractCode     string
	ContractByteCode string
	ContractAbi      string
}

// ContractListView contract
type ContractListView struct {
	Id               string
	ContractName     string
	ContractSymbol   string
	ContractAddr     string
	Version          string
	TxNum            int64
	Status           int32
	Creator          string
	CreatorAddr      string
	CreatorAddrBns   string
	Upgrader         string
	UpgradeAddr      string
	UpgradeOrgId     string
	CreateTimestamp  int64
	UpgradeTimestamp int64
	ContractType     string
	RuntimeType      string
}

// InnerContractListView inner
type InnerContractListView struct {
	Id              string
	ContractName    string
	Addr            string
	Version         string
	RuntimeType     string
	Creator         string
	CreatorAddr     string
	Upgrader        string
	UpgradeOrgId    string
	Status          int
	UpdateTimestamp int64
	CreateTimestamp int64
	ContractType    string
}

// ChainDecimalView chain
type ChainDecimalView struct {
	BlockHeight     int64
	OrgNum          int
	ContractNum     int
	TransactionNum  int64
	ContractExecNum int64
	TxNum           float32
}

// ChainListView chain
type ChainListView struct {
	ChainId      string
	ChainVersion string
	Status       int
	Consensus    string
	Timestamp    int64
	AuthType     string
}

// TxDetailView tx
// param view
type TxDetailView struct {
	TxId               string
	BlockHash          string
	BlockHeight        int64
	Sender             string
	SenderOrgId        string
	ContractName       string
	ContractNameBak    string
	ContractAddr       string
	ContractMessage    string
	ContractVersion    string
	TxStatusCode       string
	TxStatus           int
	ContractResultCode uint32
	ContractResult     []byte
	RwSetHash          string
	ContractMethod     string
	ContractParameters string
	Endorsement        string
	TxType             string
	Timestamp          int64
	UserAddr           string
	UserAddrBns        string
	ContractRead       string
	ContractWrite      string
	GasUsed            uint64
	Payer              string
	PayerBns           string
	Event              string
	RuntimeType        string
	ShowStatus         int
}

// ContractVersionView contract
type ContractVersionView struct {
	TxId               string
	Sender             string
	SenderAddr         string
	SenderAddrBNS      string
	SenderOrgId        string
	Version            string
	ContractName       string
	ContractAddr       string
	TxUrl              string
	RuntimeType        string
	ContractResultCode int
	Timestamp          int64
}

// InnerContractVersionView inner
type InnerContractVersionView struct {
	ContractName string
	Sender       string
	SenderAddr   string
	SenderOrgId  string
	Version      string
	RuntimeType  string
	TxUrl        string
	Status       int
	Timestamp    int64
}

// Endorsement endorsement
type Endorsement struct {
}

// LatestBlockListView latest
type LatestBlockListView struct {
	Id               int64
	BlockHeight      int64
	BlockHash        string
	TxNum            int
	ProposalNodeId   string
	ProposalNodeAddr string
	Timestamp        int64
}

// LatestTxListView latest
type LatestTxListView struct {
	Id              int
	TxId            string
	BlockHash       string
	BlockHeight     int64
	Status          string
	Timestamp       int64
	ContractName    string
	ContractNameBak string
	ContractAddr    string
	Sender          string
	UserAddr        string
	UserAddrBns     string
	GasUsed         uint64
}

// NodesView view
type NodesView struct {
	NodeId      string
	NodeName    string
	NodeAddress string
	Role        string
	OrgId       string
	Status      int
	Timestamp   int64
}

// OrgView view
type OrgView struct {
	ChainId   string //链标识
	OrgId     string
	Status    int // 0:正常 1: 已删掉
	UserCount int64
	NodeCount int64
}

// BlockListView block
type BlockListView struct {
	BlockHeight      int64
	BlockHash        string
	TxNum            int
	ProposalNodeId   string
	ProposalNodeAddr string
	Timestamp        int64
}

// TxListView tx
// param view
type TxListView struct {
	Id                 string
	TxId               string
	Sender             string
	SenderOrgId        string
	BlockHeight        int64
	ContractName       string
	ContractAddr       string
	ContractMethod     string
	ContractParameters string
	Status             string
	TxStatus           int
	ShowStatus         int
	BlockHash          string
	Timestamp          int64
	UserAddr           string
	UserAddrBns        string
	GasUsed            uint64
	PayerAddr          string
	//Event              string
}

// BlockTxListView tx
// param view
type BlockTxListView struct {
	TxId         string
	Sender       string
	SenderOrgId  string
	BlockHeight  int64
	ContractName string
	ContractAddr string
	TxStatus     int
	BlockHash    string
	Timestamp    int64
	UserAddr     string
	UserAddrBns  string
}

// ContractTxListView tx
// param view
type ContractTxListView struct {
	TxId           string
	Sender         string
	SenderOrgId    string
	BlockHeight    int64
	ContractName   string
	ContractAddr   string
	ContractMethod string
	TxStatus       int
	ShowStatus     int
	BlockHash      string
	Timestamp      int64
	UserAddr       string
	UserAddrBns    string
}

// DetailView detail
type DetailView struct {
	Type      int64
	Id        int64
	BlockHash string
}

// InnerTxListView inner
type InnerTxListView struct {
	Id           string
	TxId         string
	Sender       string
	SenderOrgId  string
	ContractName string
	TxStatus     int
	Timestamp    int64
	UserAddr     string
	GasUsed      uint64
}

// UserListView view
type UserListView struct {
	Id        string
	UserId    string
	OrgId     string
	Role      string
	Timestamp int64
	UserAddr  string
	Status    int
}

// AccountListView view
type AccountListView struct {
	AddrType  int
	Address   string
	BNS       string
	DID       string
	Timestamp int64
}

// SearchView search
type SearchView struct {
	Type         int
	Data         string
	ChainId      string
	ContractType string
}

// TextContent text
type TextContent struct {
	Content string `json:"content"`
}

// NewSearchView new
//func NewSearchView(searchType int, data string, chainId string) *SearchView {
//	view := SearchView{
//		Type:    searchType,
//		Data:    data,
//		ChainId: chainId,
//		ContractType: chainId,
//	}
//	return &view
//}

// GasListView view
type GasListView struct {
	Id         string
	GasBalance int64
	GasTotal   int64
	GasUsed    int64
	Address    string
	ChainId    string
	Timestamp  int64
}

// GasRecordListView view
type GasRecordListView struct {
	Id        string
	GasAmount int64
	Address   string
	//PayerAddress string
	BusinessType int
	TxId         string
	ChainId      string
	Timestamp    int64
}

// GasInfoView view
type GasInfoView struct {
	GasBalance int64
}

// EvidenceListView view
type EvidenceListView struct {
	Id          string
	ChainId     string
	TxId        string
	SenderAddr  string
	Hash        string
	Timestamp   int64
	Code        int
	MetaData    string
	BlockHeight int64
	ResultCode  uint32
}

// FungibleContractListView contract
type FungibleContractListView struct {
	ContractName   string
	ContractAddr   string
	ContractSymbol string
	ContractType   string
	TotalSupply    string
	HolderCount    int64
	Timestamp      int64
	TxNum          int64
}

// NonFungibleContractListView contract
type NonFungibleContractListView struct {
	ContractName string
	ContractAddr string
	ContractType string
	TotalSupply  string
	HolderCount  int64
	Timestamp    int64
	TxNum        int64
}

// EvidenceContractListView contract
type EvidenceContractListView struct {
	ContractName string
	ContractAddr string
	Timestamp    int64
}

// EvidenceContractView contract
type EvidenceContractView struct {
	Id           string
	ChainId      string
	ContractName string
	SenderAddr   string
	TxId         string
	BlockHeight  int64
	Timestamp    int64
	EvidenceId   string
	Hash         string
	MetaData     string
	Code         int
	ResultCode   uint32
}

// EvidenceMetaData 存证合约MetaData
type EvidenceMetaData struct {
	Key   string
	Value string
}

// IdentityContractListView contract
type IdentityContractListView struct {
	ContractName string
	ContractAddr string
	Timestamp    int64
}

// IdentityContractView contract
type IdentityContractView struct {
	ContractName string
	ContractAddr string
	UserAddr     string
	Level        string
	PkPem        string
}

// FungibleTransferListView tx
// param view
type FungibleTransferListView struct {
	TxId           string
	ContractName   string
	ContractAddr   string
	ContractMethod string
	ContractSymbol string
	From           string
	FromBNS        string
	To             string
	ToBNS          string
	Amount         string
	Timestamp      int64
}

// NonFungibleTransferListView tx
// param view
type NonFungibleTransferListView struct {
	TxId           string
	ContractName   string
	ContractAddr   string
	ContractMethod string
	From           string
	FromBNS        string
	To             string
	ToBNS          string
	TokenId        string
	Timestamp      int64
}

// NFTListView detail
type NFTListView struct {
	Timestamp    int64
	ContractName string
	ContractAddr string
	TokenId      string
	AddrType     int
	OwnerAddr    string
	OwnerAddrBNS string
	CategoryName string
	ImageUrl     string
	UrlType      string
}

// MetadataJson Metadata
type MetadataJson struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	OrgName     string `json:"orgName"`
	ImageUrl    string `json:"imageUrl"`
	Image       string `json:"image"`
	Description string `json:"description"`
	SeriesHash  string `json:"seriesHash"`
}

// NFTDetailView detail
type NFTDetailView struct {
	Timestamp    int64
	ContractName string
	ContractAddr string
	TokenId      string
	AddrType     int
	OwnerAddr    string
	OwnerAddrBNS string
	CategoryName string
	IsViolation  bool
	Metadata     Metadata
}

// Metadata Metadata
type Metadata struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	OrgName     string `json:"orgName"`
	ImageUrl    string `json:"imageUrl"`
	UrlType     string `json:"urlType"`
	Description string `json:"description"`
	SeriesHash  string `json:"seriesHash"`
}

// PositionListView detail
type PositionListView struct {
	AddrType       int
	OwnerAddr      string
	OwnerAddrBNS   string
	ContractAddr   string
	ContractName   string
	ContractSymbol string
	ContractType   string
	Amount         string
	HoldRatio      string
	HoldRank       int64
}

// NonPositionListView detail
type NonPositionListView struct {
	AddrType     int
	OwnerAddr    string
	OwnerAddrBNS string
	ContractAddr string
	ContractName string
	Amount       string
	HoldRatio    string
	HoldRank     int64
}

// AccountDetailView account
type AccountDetailView struct {
	Address string
	Type    int
	BNS     string
	DID     string
}

// NonFungibleContractListView contract
type IDAContractListView struct {
	ContractName string
	ContractAddr string
	ContractType string
	DataAssetNum int64
	Timestamp    int64
}

// NonFungibleContractListView contract
type IDAAssetListView struct {
	AssetCode   string
	Creator     string
	IsDeleted   bool
	CreatedTime int64
	UpdatedTime int64
}

// IDAAssetDetailView detail
type IDAAssetDetailView struct {
	AssetCode         string
	AssetName         string
	AssetEnName       string
	Category          int
	ImmediatelySupply bool //是否及时供应
	SupplyTime        time.Time
	DataScale         string
	IndustryTitle     string
	Summary           string
	AnnexUrls         []string
	UserCategories    string
	UpdateCycleType   int
	TimeSpan          string
	CreatedTime       int64
	UpdatedTime       int64
	IsDeleted         bool
	DataAsset         []DataAsset
	ApiAsset          []ApiAsset
}

// DataAsset DataAsset
type DataAsset struct {
	Name         string
	Type         string
	Length       int
	IsPrimaryKey int
	IsNotNull    int
	PrivacyQuery int
}

// ApiAsset ApiAsset
type ApiAsset struct {
	Header       string
	Params       string
	Response     string
	Method       string
	ResponseType string
	Url          string
}
