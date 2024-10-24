/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package db

import "github.com/shopspring/decimal"

// SubscribeNode json
type SubscribeNode struct {
	Addr        string
	OrgCA       string
	TLSHostName string
	Tls         bool
}

// FungibleContractWithTxNum 连表查询同质化合约和交易量
type FungibleContractWithTxNum struct {
	FungibleContract
	TxNum int64 `json:"txNum" gorm:"column:txNum"` // 交易数据量
}

// NonFungibleContractWithTxNum 连表查询非同质化合约和交易量
type NonFungibleContractWithTxNum struct {
	NonFungibleContract
	TxNum int64 `json:"txNum" gorm:"column:txNum"` // 交易数据量
}

// PositionWithRank 连表查询持仓和排名
type PositionWithRank struct {
	FungiblePosition
	HoldRank int64 `json:"holdRank" gorm:"column:holdRank"` // 持仓排名
}

// FTPositionJoinContract 连表查询持仓和合约
type FTPositionJoinContract struct {
	FungiblePosition
	ContractType string `json:"contractType" gorm:"column:contractType"` // 合约类型
}

// ContractPositionAccount 连表查询持仓和排名
type ContractPositionAccount struct {
	OwnerAddr string          `json:"ownerAddr" gorm:"column:ownerAddr"`
	Amount    decimal.Decimal `json:"amount" gorm:"column:amount"`
	BNS       string          `json:"bns" gorm:"column:bns"`
	AddrType  int             `json:"addrType" gorm:"column:addrType"`
}

// CycleJoinTransferResult 跨链交易连表查询
//
//nolint:govet
type CycleJoinTransferResult struct {
	CrossId string `gorm:"column:crossId"`
	CrossCycleTransaction
	CrossTransactionTransfer
}

type BlockTxListResult struct {
	TxId         string `json:"txId" gorm:"column:txId"`
	Sender       string `json:"sender" gorm:"column:sender"`
	SenderOrgId  string `json:"senderOrgId" gorm:"column:senderOrgId"`
	UserAddr     string `json:"userAddr" gorm:"column:userAddr"`
	TxStatusCode string `json:"txStatusCode" gorm:"column:txStatusCode"`
	ContractName string `json:"contractName" gorm:"column:contractName"`
	ContractAddr string `json:"contractAddr" gorm:"column:contractAddr"`
	Timestamp    int64  `json:"timestamp" gorm:"column:timestamp"`
}

type ContractTxListResult struct {
	TxId               string `json:"txId" gorm:"column:txId"`
	Sender             string `json:"sender" gorm:"column:sender"`
	SenderOrgId        string `json:"senderOrgId" gorm:"column:senderOrgId"`
	UserAddr           string `json:"userAddr" gorm:"column:userAddr"`
	TxStatusCode       string `json:"txStatusCode" gorm:"column:txStatusCode"`
	ContractName       string `json:"contractName" gorm:"column:contractName"`
	ContractAddr       string `json:"contractAddr" gorm:"column:contractAddr"`
	ContractMethod     string `json:"contractMethod" gorm:"column:contractMethod"`
	ContractMessageBak string `json:"contractMessageBak" gorm:"column:contractMessageBak"`
	ReadSetBak         string `json:"readSetBak" gorm:"column:readSetBak"`
	BlockHeight        int64  `json:"blockHeight" gorm:"column:blockHeight"`
	Timestamp          int64  `json:"timestamp" gorm:"column:timestamp"`
}
