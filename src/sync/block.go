/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

// DealBlockInfo
//
//	@Description:
//	@param blockInfo
//	@param hashType
//	@return *db.Block
//	@return error
func DealBlockInfo(blockInfo *common.BlockInfo, hashType string) (*db.Block, error) {
	if blockInfo == nil {
		return nil, fmt.Errorf("blockInfo is nil")
	}
	chainId := blockInfo.Block.Header.ChainId
	blockHeight := blockInfo.Block.Header.BlockHeight

	newUUID := uuid.New().String()
	modBlock := &db.Block{
		ID:            newUUID,
		BlockHeight:   int64(blockHeight),
		BlockHash:     hex.EncodeToString(blockInfo.Block.Header.BlockHash),
		BlockVersion:  int32(blockInfo.Block.Header.BlockVersion),
		PreBlockHash:  hex.EncodeToString(blockInfo.Block.Header.PreBlockHash),
		ConsensusArgs: utils.Base64Encode(blockInfo.Block.Header.ConsensusArgs),
		DagHash:       hex.EncodeToString(blockInfo.Block.Header.DagHash),
		Timestamp:     blockInfo.Block.Header.BlockTimestamp,
		TxCount:       int(blockInfo.Block.Header.TxCount),
	}
	modBlock.RwSetHash = hex.EncodeToString(blockInfo.Block.Header.RwSetRoot)
	modBlock.Signature = utils.Base64Encode(blockInfo.Block.Header.Signature)
	modBlock.TxRootHash = hex.EncodeToString(blockInfo.Block.Header.TxRoot)
	member := blockInfo.Block.Header.Proposer
	if member != nil {
		modBlock.OrgId = member.OrgId
		//根据proposer信息填充block中的地址,id信息
		getInfos, err := getMemberIdAddrAndCert(chainId, hashType, member)
		if err != nil {
			log.Error("getMemberIdAddrAndCert Failed: " + err.Error())
			return modBlock, err
		}
		modBlock.ProposerAddr = getInfos.UserAddr
		modBlock.ProposerId = getInfos.UserId
	}

	//解析 dag
	dagBytes, _ := json.Marshal(blockInfo.Block.Dag)
	modBlock.BlockDag = string(dagBytes)
	return modBlock, nil
}

// BuildLatestBlockListCache
//
//	@Description:设置最新区块缓存列表
//	@param chainId
//	@param modBlock 区块数据
func BuildLatestBlockListCache(chainId string, modBlock *db.Block) {
	var blockList []*db.Block
	//获取最新区块列表
	blockListCache, _ := dbhandle.GetLatestBlockListCache(chainId)
	if len(blockListCache) > 0 {
		//缓存存在
		blockList = append(blockList, modBlock)
	} else {
		//缓存可能丢失
		blockList, _ = dbhandle.GetLatestBlockListCache(chainId)
	}
	if len(blockList) == 0 {
		return
	}
	// 缓存交易信息
	dbhandle.SetLatestBlockListCache(chainId, blockList)
}

// BuildOverviewMaxBlockHeightCache
//
//	@Description: 缓存最高区块高度
//	@param chainId
//	@param blockInfo
func BuildOverviewMaxBlockHeightCache(chainId string, blockInfo *db.Block) {
	maxBlockHeight := blockInfo.BlockHeight
	dbhandle.SetMaxBlockHeightCache(chainId, maxBlockHeight)
}
