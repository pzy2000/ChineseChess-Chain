/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"context"
	"fmt"
	"strconv"
	"time"
)

// BatchInsertNode 批量插入node数据，数据库存在就忽略
func BatchInsertNode(chainId string, nodeList []*db.Node) error {
	if len(nodeList) == 0 {
		return nil
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	delKey := make([]string, 0)
	nodeIds := make([]string, 0, len(nodeList))
	for _, node := range nodeList {
		nodeIds = append(nodeIds, node.NodeId)
	}
	nodeMap, err := GetNodeInById(chainId, nodeIds)
	if err != nil || len(nodeMap) == len(nodeList) {
		return err
	}

	insertNodes := make([]*db.Node, 0)
	//updateNodeIds := make([]string, 0)
	for _, node := range nodeList {
		if _, ok := nodeMap[node.NodeId]; !ok {
			insertNodes = append(insertNodes, node)
			redisKey := fmt.Sprintf(cache.RedisOverviewNodeCount, prefix, chainId, node.Role)
			delKey = append(delKey, redisKey)
		}
	}

	// //更新节点
	// deErr := UpdateNode(chainId, updateNodeIds, config.StatusNormal)
	// if deErr != nil {
	// 	log.Error("[DB] update Node Info Failed : " + deErr.Error())
	// 	return err
	// }

	err = InsertNodes(chainId, insertNodes)
	if err != nil {
		return err
	}

	//删除缓存
	ctx := context.Background()
	redisConsensusKey := fmt.Sprintf(cache.RedisOverviewNodeCount, prefix, chainId, "consensus")
	redisCommonKey := fmt.Sprintf(cache.RedisOverviewNodeCount, prefix, chainId, "common")
	delKey = append(delKey, redisConsensusKey)
	delKey = append(delKey, redisCommonKey)
	_ = cache.GlobalRedisDb.Del(ctx, delKey...).Err()
	return err
}

// InsertNodes 插入节点
func InsertNodes(chainId string, insertNodes []*db.Node) error {
	if len(insertNodes) == 0 {
		return nil
	}
	tableName := db.GetTableName(chainId, db.TableNode)
	err := db.GormDB.Table(tableName).Create(&insertNodes).Error
	return err
}

// GetNodeInById 根据id获取节点数据
func GetNodeInById(chainId string, nodeIds []string) (map[string]*db.Node, error) {
	nodes := make([]*db.Node, 0)
	nodeMap := make(map[string]*db.Node, 0)
	if len(nodeIds) == 0 {
		return nodeMap, nil
	}
	tableName := db.GetTableName(chainId, db.TableNode)
	err := db.GormDB.Table(tableName).Where("nodeId in ?", nodeIds).Find(&nodes).Error
	if err != nil {
		return nodeMap, err
	}
	for _, node := range nodes {
		nodeMap[node.NodeId] = node
	}

	return nodeMap, err
}

// UpdateNode delete
// @desc
// @param ${param}
// @return error
func UpdateNode(chainId string, nodeIds []string, status int) error {
	if len(nodeIds) == 0 {
		return nil
	}
	tableName := db.GetTableName(chainId, db.TableNode)
	for _, nodeId := range nodeIds {
		where := map[string]interface{}{
			"nodeId": nodeId,
		}
		params := map[string]interface{}{
			"status": status,
		}
		err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
		if err != nil {
			log.Errorf("update node status fail, chainId : %s, nodeId : %s , cause: %s",
				chainId, nodeId, err.Error())
			return err
		}
	}
	return nil
}

// UpdateNode delete
// @desc
// @param ${param}
// @return error
func DeleteNodeById(chainId string, nodeIds []string) error {
	tableName := db.GetTableName(chainId, db.TableNode)
	err := db.GormDB.Table(tableName).Where("nodeId in ?", nodeIds).Delete(nil).Error
	if err != nil {
		log.Errorf("DeleteNodeById node status fail, chainId : %s, nodeIds : %s , cause: %s",
			chainId, nodeIds, err.Error())
		return err
	}
	return nil
}

// GetNodesRef get
// @desc
// @param ${param}
// @return []*NodeIds
// @return error
func GetNodesRef(chainId string) ([]string, error) {
	nodes := make([]*db.Node, 0)
	where := map[string]interface{}{
		"status": config.StatusNormal,
	}
	tableName := db.GetTableName(chainId, db.TableNode)
	err := db.GormDB.Table(tableName).Where(where).Find(&nodes).Error
	if err != nil {
		return nil, fmt.Errorf("GetNodesRef By chainId err, cause : %s", err.Error())
	}

	nodeIds := make([]string, 0)
	for _, node := range nodes {
		nodeIds = append(nodeIds, node.NodeId)
	}
	return nodeIds, nil
}

// GetNodeList 获取es的node列表
func GetNodeList(chainId, nodeName, orgId, nodeId string, offset, limit int) ([]*db.Node, int64, error) {
	nodes := make([]*db.Node, 0)
	var count int64
	where := map[string]interface{}{
		"status": config.StatusNormal,
	}
	if nodeName != "" {
		where["nodeName"] = nodeName
	}
	if orgId != "" {
		where["orgId"] = orgId
	}
	if nodeId != "" {
		where["nodeId"] = nodeId
	}
	tableName := db.GetTableName(chainId, db.TableNode)
	query := db.GormDB.Table(tableName).Where(where)
	err := query.Count(&count).Error
	if err != nil {
		return nodes, 0, err
	}
	query = query.Offset(offset * limit).Limit(limit)
	err = query.Find(&nodes).Error
	if err != nil {
		return nodes, 0, err
	}
	return nodes, count, nil
}

// GetNodeNumByOrg 获取节点数量
func GetNodeNumByOrg(chainId string, orgId string) (int64, error) {
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisOverviewNodeCount, prefix, chainId, orgId)
	countRes, err := GetNodeNumCache(redisKey)
	if err == nil && countRes != 0 {
		return countRes, nil
	}

	var count int64
	where := map[string]interface{}{
		"status": config.StatusNormal,
	}
	if orgId != "" {
		where["orgId"] = orgId
	}
	tableName := db.GetTableName(chainId, db.TableNode)
	err = db.GormDB.Table(tableName).Where(where).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count nodes err, cause : %s", err.Error())
	}

	// 设置键值对和过期时间
	ctx := context.Background()
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, count, time.Hour).Err()
	return count, nil
}

// GetNodeNum 获取节点数量
func GetNodeNum(chainId, role string) (int64, error) {
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisOverviewNodeCount, prefix, chainId, role)
	countRes, err := GetNodeNumCache(redisKey)
	if err == nil && countRes != 0 {
		return countRes, nil
	}

	var count int64
	where := map[string]interface{}{
		"status": config.StatusNormal,
	}
	if role != "" {
		where["role"] = role
	}
	tableName := db.GetTableName(chainId, db.TableNode)
	err = db.GormDB.Table(tableName).Where(where).Count(&count).Error
	if err != nil {
		return 0, err
	}

	// 设置键值对和过期时间
	ctx := context.Background()
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, count, time.Hour).Err()
	return count, nil
}

// GetNodeNumCache 获取node数量缓存
func GetNodeNumCache(redisKey string) (int64, error) {
	//获取缓存
	var count int64
	ctx := context.Background()
	result, err := cache.GlobalRedisDb.Get(ctx, redisKey).Result()
	if err != nil {
		return 0, err
	}
	count, err = strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}
