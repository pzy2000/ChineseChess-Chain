/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// GetDBSubscribeChains get
// @desc
// @param ${param}
// @return []*dao.Subscribe
// @return error
func GetDBSubscribeChains() ([]*db.Subscribe, error) {
	statusList := []int{
		db.SubscribeOK,
		db.SubscribeFailed,
	}

	subscribe := make([]*db.Subscribe, 0)
	err := db.GormDB.Table(db.TableSubscribe).Where("status in ?", statusList).Find(&subscribe).Error
	if err != nil {
		return subscribe, err
	}
	return subscribe, nil
}

// GetSubscribeByChainIds
//
//	@Description: 获取订阅列表
//	@param chainIds 链ID数组
//	@return []*db.Subscribe
//	@return error
func GetSubscribeByChainIds(chainIds []string) ([]*db.Subscribe, error) {
	subscribe := make([]*db.Subscribe, 0)
	err := db.GormDB.Table(db.TableSubscribe).Where("chainId in ?", chainIds).Find(&subscribe).Error
	if err != nil {
		return subscribe, fmt.Errorf("GetSubscribeByChainIds By chainIds err, cause : %s", err.Error())
	}

	return subscribe, err
}

// GetSubscribeByChainId
//
//	@Description:
//	@param chainId
//	@return *db.Subscribe
//	@return error
func GetSubscribeByChainId(chainId string) (*db.Subscribe, error) {
	subscribe := &db.Subscribe{}
	where := map[string]interface{}{
		"chainId": chainId,
	}
	err := db.GormDB.Table(db.TableSubscribe).Where(where).First(subscribe).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return subscribe, nil
}

// SetSubscribeStatus
//
//	@Description: 修改订阅状态
//	@param chainId
//	@param status （订阅状态 0:成功， 1:失败，2:暂停，3删除）
//	@return error
func SetSubscribeStatus(chainId string, status int) error {
	if chainId == "" {
		return db.ErrTableParams
	}

	where := map[string]interface{}{
		"chainId": chainId,
	}
	params := map[string]interface{}{
		"status": status,
	}
	err := db.GormDB.Table(db.TableSubscribe).Where(where).Updates(params).Error
	return err
}

// InsertSubscribe
//
//	@Description: 插入订阅数据
//	@param subscribe
//	@return error
func InsertSubscribe(subscribe *db.Subscribe) error {
	if subscribe == nil || subscribe.ChainId == "" {
		return db.ErrTableParams
	}

	return InsertData(db.TableSubscribe, subscribe)
	//return db.GormDB.Table(db.TableSubscribe).Create(subscribe).Error
}

// UpdateSubscribe
//
//	@Description: 更新订阅数据
//	@param subscribe
//	@return error
func UpdateSubscribe(subscribe *db.Subscribe) error {
	where := map[string]interface{}{
		"chainId": subscribe.ChainId,
	}
	params := map[string]interface{}{
		"orgId":    subscribe.OrgId,
		"userCert": subscribe.UserCert,
		"userKey":  subscribe.UserKey,
		"authType": subscribe.AuthType,
		"hashType": subscribe.HashType,
		"nodeList": subscribe.NodeList,
		"status":   subscribe.Status,
	}
	return db.GormDB.Table(db.TableSubscribe).Model(&db.Subscribe{}).Where(where).Updates(params).Error
}

// BuildSubscribeInfo BuildSubscribeInfo
func BuildSubscribeInfo(chainInfo *config.ChainInfo, status int) *db.Subscribe {
	if chainInfo == nil {
		return nil
	}
	var (
		userKey  string
		userCert string
	)

	nodes, _ := json.Marshal(chainInfo.NodesList)
	if chainInfo.UserInfo != nil {
		userKey = chainInfo.UserInfo.UserKey
		userCert = chainInfo.UserInfo.UserCert
	}
	subscribeInfo := &db.Subscribe{
		ChainId:  chainInfo.ChainId,
		OrgId:    chainInfo.OrgId,
		UserKey:  userKey,
		UserCert: userCert,
		NodeList: string(nodes),
		Status:   status,
		AuthType: chainInfo.AuthType,
		HashType: chainInfo.HashType,
	}
	return subscribeInfo
}

// InsertOrUpdateSubscribe
//
//	@Description:  存在就更新，不存在就插入
//	@param chainInfo
//	@param Status
//	@return error
func InsertOrUpdateSubscribe(chainInfo *config.ChainInfo, status int) error {
	if chainInfo == nil {
		return nil
	}
	var (
		userKey  string
		userCert string
	)

	nodes, _ := json.Marshal(chainInfo.NodesList)
	subscribe, _ := GetSubscribeByChainId(chainInfo.ChainId)
	if chainInfo.UserInfo != nil {
		userKey = chainInfo.UserInfo.UserKey
		userCert = chainInfo.UserInfo.UserCert
	}
	subscribeInfo := &db.Subscribe{
		ChainId:  chainInfo.ChainId,
		OrgId:    chainInfo.OrgId,
		UserKey:  userKey,
		UserCert: userCert,
		NodeList: string(nodes),
		Status:   status,
		AuthType: chainInfo.AuthType,
		HashType: chainInfo.HashType,
	}
	if subscribe == nil {
		err := InsertSubscribe(subscribeInfo)
		return err
	}

	//订阅状态一样，且不是成功就不用更新了，成功订阅有能是更新节点，所以也要更新
	if status == subscribe.Status && status != db.SubscribeOK {
		return nil
	}

	//更新订阅状态
	err := UpdateSubscribe(subscribeInfo)
	return err
}

// DeleteSubscribe
//
//	@Description: 删除订阅
//	@param chainId
//	@return error
func DeleteSubscribe(chainId string) error {
	subscribe := &db.Subscribe{}
	where := map[string]interface{}{
		"chainId": chainId,
	}
	err := db.GormDB.Table(db.TableSubscribe).Where(where).Delete(&subscribe).Error
	if err != nil {
		log.Error("[DB] Delete DeleteSubscribe Failed: " + err.Error())
	}
	return err
}
