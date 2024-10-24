/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"fmt"
)

// GetUserListByAdder 根据adders查询交易
func GetUserListByAdder(chainId string, adders []string) (map[string]*db.User, error) {
	userExists := make(map[string]*db.User, 0)
	userList := make([]*db.User, 0)
	if len(adders) == 0 {
		return userExists, nil
	}
	tableName := db.GetTableName(chainId, db.TableUser)
	err := db.GormDB.Table(tableName).Where("userAddr IN ?", adders).Find(&userList).Error
	if err != nil {
		return userExists, err
	}

	for _, user := range userList {
		userExists[user.UserAddr] = user
	}
	return userExists, nil
}

// BatchInsertUser 批量保存
func BatchInsertUser(chainId string, userList []*db.User) error {
	if len(userList) == 0 {
		return nil
	}

	insertList := make([]*db.User, 0)
	//判断数据是否已经存在
	addrList := make([]string, 0)
	for _, user := range userList {
		addrList = append(addrList, user.UserAddr)
	}

	userMap, err := GetUserListByAdder(chainId, addrList)
	if err != nil {
		return err
	}

	for _, user := range userList {
		if _, ok := userMap[user.UserAddr]; !ok {
			insertList = append(insertList, user)
		}
	}

	if len(insertList) == 0 {
		return nil
	}

	//插入数据
	tableName := db.GetTableName(chainId, db.TableUser)
	return CreateInBatchesData(tableName, insertList)
}

// GetUserNum 获取用户数量
func GetUserNum(chainId, orgId string) (int64, error) {
	var count int64
	where := map[string]interface{}{}
	if orgId != "" {
		where["orgId"] = orgId
	}

	tableName := db.GetTableName(chainId, db.TableUser)
	err := db.GormDB.Table(tableName).Where(where).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count user err, cause : %s", err.Error())
	}
	return count, nil
}

// GetUserList 获取用户列表
func GetUserList(offset, limit int, chainId, orgId string, userIds, userAddrs []string) ([]*db.User, int64, error) {
	var count int64
	userList := make([]*db.User, 0)
	if chainId == "" {
		return userList, 0, db.ErrTableParams
	}

	where := map[string]interface{}{}
	whereIn := map[string]interface{}{}
	if orgId != "" {
		where["orgId"] = orgId
	}
	if len(userIds) > 0 {
		whereIn["userId"] = userIds
	}
	if len(userAddrs) > 0 {
		whereIn["userAddr"] = userAddrs
	}

	selectFile := &SelectFile{
		Where:   where,
		WhereIn: whereIn,
	}
	tableName := db.GetTableName(chainId, db.TableUser)
	query := BuildParamsQuery(tableName, selectFile)
	err := query.Count(&count).Error
	if err != nil {
		return userList, 0, err
	}

	query = query.Order("timestamp desc").Offset(offset * limit).Limit(limit)
	err = query.Find(&userList).Error
	if err != nil {
		return userList, 0, err
	}
	return userList, count, nil
}

// UpdateUserStatus 更新用户状态
func UpdateUserStatus(address, chainId string, status int) error {
	if address == "" || chainId == "" {
		return db.ErrTableParams
	}
	where := map[string]interface{}{
		"userAddr": address,
	}
	params := map[string]interface{}{
		"status": status,
	}
	tableName := db.GetTableName(chainId, db.TableUser)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}
	return nil
}

// GetUserCountByRange 获取指定时间内的交易数量
func GetUserCountByRange(chainId string, startTime, endTime int64) (int64, error) {
	var totalCount int64
	tableName := db.GetTableName(chainId, db.TableUser)
	query := db.GormDB.Table(tableName)
	// 添加时间范围条件
	if startTime > 0 && endTime > 0 {
		query = query.Where("timestamp BETWEEN ? AND ?", startTime, endTime)
	}
	err := query.Count(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}
