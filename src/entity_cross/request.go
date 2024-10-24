/*
Package entity_cross comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package entity_cross

// nolint
const (
	OffsetDefault = 0
	OffsetMin     = 0
	LimitDefault  = 10
	LimitMax      = 200
	LimitMaxSpec  = 100000
	LimitMin      = 0

	//GetMainCrossConfig 主子链网配置
	GetMainCrossConfig = "GetMainCrossConfig"
	//CrossSearch 首页搜索
	CrossSearch = "CrossSearch"
	//CrossOverviewData 首页详情
	CrossOverviewData = "CrossOverviewData"
	//CrossLatestTxList 最新跨链交易列表
	CrossLatestTxList = "CrossLatestTxList"
	//CrossLatestSubChainList 最新子链列表
	CrossLatestSubChainList = "CrossLatestSubChainList"
	//GetCrossTxList 跨链交易列表
	GetCrossTxList         = "GetCrossTxList"
	GetCrossSubChainTxList = "GetCrossSubChainTxList"
	//CrossSubChainList 子链列表
	CrossSubChainList = "CrossSubChainList"
	//CrossSubChainDetail 子链详情
	CrossSubChainDetail = "CrossSubChainDetail"
	//GetCrossTxDetail 交易详情
	GetCrossTxDetail = "GetCrossTxDetail"
	//SubChainCrossChainList 子链跨链列表
	SubChainCrossChainList = "SubChainCrossChainList"
	//CrossUpdateSubChain 更新子链信息
	CrossUpdateSubChain    = "CrossUpdateSubChain"
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
	ChainId string
	Offset  int
	Limit   int
}

// IsLegal legal
func (rangeBody *RangeBody) IsLegal() bool {
	if rangeBody.Limit > LimitMax || rangeBody.Offset < OffsetMin {
		return false
	}
	if rangeBody.ChainId == "" {
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

// CrossSearchParams get
type CrossSearchParams struct {
	ChainId string
	Value   string
}

// IsLegal legal
func (params *CrossSearchParams) IsLegal() bool {
	return params.ChainId != "" && params.Value != ""
}

// GetCrossTxListParams get
type GetCrossTxListParams struct {
	RangeBody
	CrossId       string
	SubChainId    string
	FromChainName string
	ToChainName   string
	StartTime     int64
	EndTime       int64
}

// IsLegal legal
func (params *GetCrossTxListParams) IsLegal() bool {
	return true
}

// GetCrossSubChainListParams get
type GetCrossSubChainListParams struct {
	RangeBody
	SubChainId   string
	SubChainName string
}

// IsLegal legal
func (params *GetCrossSubChainListParams) IsLegal() bool {
	return true
}

// GetCrossSubChainDetailParams get
type GetCrossSubChainDetailParams struct {
	ChainId    string
	SubChainId string
}

// IsLegal legal
func (params *GetCrossSubChainDetailParams) IsLegal() bool {
	return params.ChainId != "" && params.SubChainId != ""
}

// GetCrossTxDetailParams get
type GetCrossTxDetailParams struct {
	ChainId string
	CrossId string
}

// IsLegal legal
func (params *GetCrossTxDetailParams) IsLegal() bool {
	return params.ChainId != "" && params.CrossId != ""
}

// GetSubChainCrossChainListParams get
type GetSubChainCrossChainListParams struct {
	ChainId    string
	SubChainId string
}

// IsLegal legal
func (params *GetSubChainCrossChainListParams) IsLegal() bool {
	return params.ChainId != "" && params.SubChainId != ""
}

// CrossUpdateSubChainParams get
type CrossUpdateSubChainParams struct {
	ChainId         string
	ChainRid        string
	SubChainId      string
	SubChainName    string
	GatewayId       string
	GatewayName     string
	GatewayAddr     string
	CrossCa         string
	SdkClientCrt    string
	SdkClientKey    string
	SpvContractName string
	Introduction    string
	ExplorerAddr    string
	ExplorerTxAddr  string
	TxNum           int64
	ChainType       int32
	BlockHeight     int64
}

// IsLegal legal
func (params *CrossUpdateSubChainParams) IsLegal() bool {
	return params.ChainId != "" && params.ChainRid != "" && params.GatewayId != "" && params.SpvContractName != ""
}
