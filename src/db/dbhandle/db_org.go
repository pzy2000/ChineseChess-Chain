/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"fmt"

	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"

	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
)

// SaveOrgByConfig load ogr info
func SaveOrgByConfig(chainConfig *pbconfig.ChainConfig) error {
	if chainConfig.ChainId == "" {
		return db.ErrTableParams
	}
	if len(chainConfig.TrustRoots) == 0 {
		return nil
	}

	orgCount := len(chainConfig.TrustRoots)
	orgList := make([]*db.Org, 0, orgCount)
	var orgIds []string
	for _, trustRoot := range chainConfig.TrustRoots {
		org := &db.Org{
			OrgId:  trustRoot.OrgId,
			Status: config.StatusNormal,
		}
		orgList = append(orgList, org)
		orgIds = append(orgIds, trustRoot.OrgId)
	}
	//获取交易表名称
	tableName := db.GetTableName(chainConfig.ChainId, db.TableOrg)
	// 查询已经存在的 orgId
	existingOrgs := make([]*db.Org, 0)
	err := db.GormDB.Table(tableName).Where("orgId in ?", orgIds).Find(&existingOrgs).Error
	if err != nil {
		log.Error("Query existing orgs failed:", err)
		return err
	}

	// 创建一个 map 存储已经存在的 orgId
	existingOrgMap := make(map[string]bool)
	for _, org := range existingOrgs {
		existingOrgMap[org.OrgId] = true
	}

	// 过滤出不存在的 orgId 并插入数据库
	newOrgList := make([]*db.Org, 0, orgCount)
	for _, org := range orgList {
		if !existingOrgMap[org.OrgId] {
			newOrgList = append(newOrgList, org)
		}
	}

	if len(newOrgList) == 0 {
		return nil
	}

	err = CreateInBatchesData(tableName, newOrgList)
	if err != nil {
		log.Error("SaveOrgByConfig failed err:", err)
		return err
	}
	return nil
}

// GetOrgList 获取es的org列表
func GetOrgList(chainId, orgId string, offset, limit int) ([]*db.Org, int64, error) {
	orgList := make([]*db.Org, 0)
	var count int64
	where := map[string]interface{}{
		"status": config.StatusNormal,
	}
	if orgId != "" {
		where["orgId"] = orgId
	}
	tableName := db.GetTableName(chainId, db.TableOrg)
	query := db.GormDB.Table(tableName).Where(where)
	err := query.Count(&count).Error
	if err != nil {
		return orgList, 0, fmt.Errorf("count orgs err, cause : %s", err.Error())
	}

	query = query.Order("createdAt desc").Offset(offset * limit).Limit(limit)
	err = query.Find(&orgList).Error
	return orgList, count, nil
}

// GetOrgNum 获取组织数量
func GetOrgNum(chainId string) (int64, error) {
	var count int64
	where := map[string]interface{}{
		"status": config.StatusNormal,
	}
	tableName := db.GetTableName(chainId, db.TableOrg)
	err := db.GormDB.Table(tableName).Where(where).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count orgs err, cause : %s", err.Error())
	}
	return count, nil
}
