/*
Package db comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package db

import (
	"chainmaker.org/chainmaker/contract-utils/standard"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/shopspring/decimal"
)

const (
	KeyUserInfo          = "UserInfo"
	KeyIDABasic          = "IDABasic"
	KeyIDAOwnership      = "IDAOwnership"
	KeyIDASource         = "IDASource"
	KeyIDAScenarios      = "IDAScenarios"
	KeyIDASUppLy         = "IDASupply"
	KeyIDADetails        = "IDADetails"
	KeyIDAPrivacy        = "IDAPrivacy"
	KeyIDAStatus         = "IDAStatus"
	KeyIDAColumns        = "IDAColumns"
	KeyIDAAPI            = "IDAApi"
	KeyIDAcertifications = "IDAcertifications"
	KeyIDADataSet        = "IDADataSet"
	KeyIDAEnName         = "IDAEnName"
	KeyRegistercount     = "RegisterCount"
	KeyPlatformInfo      = "PlatformInfo"
	KeyPlatformcount     = "PlatformCount"
)

// BlockPosition 计算同质化，非同质化持仓数据
type BlockPosition struct {
	//InsertFungiblePosition新增同质化持仓
	InsertFungiblePosition []*FungiblePosition
	//UpdateFungiblePosition 更新同质化持仓
	UpdateFungiblePosition []*FungiblePosition
	//DeleteFungiblePosition 删除同质化持仓
	DeleteFungiblePosition []*FungiblePosition
	//InsertNonFungible 新增非同质化持仓
	InsertNonFungible []*NonFungiblePosition
	//UpdateNonFungible 更新非同质化持仓
	UpdateNonFungible []*NonFungiblePosition
	//DeleteNonFungible 删除飞同质化持仓
	DeleteNonFungible []*NonFungiblePosition
}

// PositionData 持仓信息
type PositionData struct {
	//AddrType 持仓地址类型
	AddrType int
	//OwnerAddr 持仓地址
	OwnerAddr string
	//ContractAddr 合约地址
	ContractAddr string
	//ContractName 合约名称
	ContractName string
	//Symbol 合约简称
	Symbol string
	//ContractType 合约类型
	ContractType string
	//Amount 持有量
	Amount decimal.Decimal
	//Decimals 小数位数
	Decimals int
}

// TokenResult 解析token数据
type TokenResult struct {
	//InsertUpdateToken 新增，更新token
	InsertUpdateToken []*NonFungibleToken
	//DeleteToken 删除token
	DeleteToken []*NonFungibleToken
}

// TransferTopicEventData 解析event数据
type TransferTopicEventData struct {
	//FromAddress from地址
	FromAddress string
	//ToAddress to地址
	ToAddress string
	//Amount 流转数量
	Amount string
	//TokenId 流转token
	TokenId      string
	CategoryName string
	Metadata     string
	Approval     string
}

// BNSTopicEventData 合约event解析BNS数据
type BNSTopicEventData struct {
	//Domain BNS地址
	Domain string
	//Value 账户地址
	Value string
	//BNS解析资源类型,string "1“-链地址，”2"-DID,"3"-去中心化网站，"4“-合约，"5"-子链
	ResourceType string
}

// DIDTopicEventData 合约event解析DID数据
type DIDTopicEventData struct {
	//Did DID值
	Did string
	//VerificationMethod did document
	VerificationMethod []VerificationMethod
}

// DidDocument did document
type DidDocument struct {
	//Id document id 不是did
	Id string `json:"id"`
	//VerificationMethod did document
	VerificationMethod []VerificationMethod `json:"verificationMethod"`
}

// VerificationMethod Method
type VerificationMethod struct {
	//Id document id 不是did
	Id         string `json:"id"`
	Type       string `json:"type"`
	Controller string `json:"controller"`
	//PublicKeyPem 公钥
	PublicKeyPem string `json:"publicKeyPem"`
	//Address did绑定账户地址
	Address string `json:"address"`
}

// ContractEventData 解析event数据
type ContractEventData struct {
	Index        int
	Topic        string
	TxId         string
	ContractName string
	EventData    *TransferTopicEventData
	Timestamp    int64
}

// IdentityEventData IdentityEventData
type IdentityEventData struct {
	UserAddr string
	Level    string
	PkPem    string
}

// SenderPayerUser 交易解析出的账户信息
type SenderPayerUser struct {
	//SenderUserId 发送交易的userid
	SenderUserId string
	// SenderUserAddr 发送交易的地址
	SenderUserAddr string
	//SenderOrgId 发起交易的组织
	SenderOrgId string
	//SenderRole 发起交易的地址身份
	SenderRole string
	//PayerUserId 支付gas的账户id
	PayerUserId string
	//PayerUserAddr 支付gas的地址
	PayerUserAddr string
}

// UpdateTxBlack 交易黑名单
type UpdateTxBlack struct {
	//AddTxBlack 添加交易黑名单
	AddTxBlack []*BlackTransaction
	// DeleteTxBlack删除交易黑名单
	DeleteTxBlack []*Transaction
}

// GetContractResult 解析合约数据
type GetContractResult struct {
	//UpdateContractTxEventNum 更新合约交易数量
	UpdateContractTxEventNum []*Contract
	//IdentityContract 身份合约
	IdentityContract []*IdentityContract
	//UpdateFungibleContract 更新同质化合约
	UpdateFungibleContract []*FungibleContract
	//UpdateNonFungible 更新非同质化合约
	UpdateNonFungible []*NonFungibleContract
	//UpdateIdaContract 更新ida合约
	UpdateIdaContract map[string]*IDAContract
}

// GetContractWriteSet 合约读写集
type GetContractWriteSet struct {
	//ContractResult 合约信息
	ContractResult *common.Contract
	//ByteCode 合约ByteCode，获取合约类型使用
	ByteCode []byte
	//Decimal 合约小数
	Decimal string
	//Symbol 合约简称
	Symbol string
}

// UpdateAccountResult 更新账户数据
type UpdateAccountResult struct {
	//InsertAccount 新增账户
	InsertAccount []*Account
	//UpdateAccount 更新账户
	UpdateAccount []*Account
}

// CrossChainResult 跨链主子链数据
type CrossChainResult struct {
	//SaveSubChainList 保存子链数据
	SaveSubChainList []*CrossSubChainData
	//CrossMainTransaction 跨链交易主链交易
	CrossMainTransaction []*CrossMainTransaction
	//CrossTransfer 跨链交易流转数据
	CrossTransfer      map[string]*CrossTransactionTransfer
	InsertCrossCycleTx []*CrossCycleTransaction
	SaveCrossCycleTx   map[string]*CrossCycleTransaction
	UpdateCrossCycleTx map[string]*CrossCycleTransaction
	//BusinessTxMap 跨链交易-具体业务交易
	BusinessTxMap map[string]*CrossBusinessTransaction
	//SubChainBlockHeight 跨链子链高度
	SubChainBlockHeight map[string]int64
	//CrossChainContractMap 跨链合约数据
	CrossChainContractMap map[string]map[string]string
	//GateWayIds 跨链网关
	GateWayIds []int64
}

// EventIDAUpdatedData 合约event解析DID数据
type EventIDAUpdatedData struct {
	//Did id
	IDACode string
	//Field 字段
	Field     string
	Update    string
	EventTime int64
}

type IDACreatedInfo struct {
	IDAInfo      *standard.IDAInfo
	ContractAddr string
	EventTime    int64
}

type IDAAssetsDataDB struct {
	IDAAssetDetail     []*IDAAssetDetail
	IDAAssetAttachment []*IDAAssetAttachment
	IDAAssetData       []*IDADataAsset
	IDAAssetApi        []*IDAApiAsset
}

type IDAAssetsUpdateDB struct {
	UpdateAssetDetails    []*IDAAssetDetail
	InsertAttachment      []*IDAAssetAttachment
	InsertIDAAssetData    []*IDADataAsset
	InsertIDAAssetApi     []*IDAApiAsset
	DeleteAttachmentCodes []string
	DeleteAssetDataCodes  []string
	DeleteAssetApiCodes   []string
}
