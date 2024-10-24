/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/sensitive"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/panjf2000/ants/v2"

	"net/http"
	"regexp"
	"sync"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	tms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tms/v20201229"
)

// nolint
const (
	Block          = "Block"
	Review         = "Review"
	FilterCapacity = 3000
)

var (
	contentBizType = "cwf"
)

var (
	globalChainMaxHeightMp sync.Map //chainid - blockheight(string-> int64)
)

func GetMaxHeight(chainId string) int64 {
	height, _ := globalChainMaxHeightMp.Load(chainId)
	return height.(int64)
}

func setMaxHeight(chainId string, height int64) {
	globalChainMaxHeightMp.Store(chainId, height)
	syncHeightGauge.WithLabelValues(chainId).Set(float64(height)) //
}

// Task task
type Task struct {
	f func() error
}

// NewTask task
func NewTask(f func() error) *Task {
	return &Task{
		f: f,
	}
}

// 执行进程
func (t *Task) execute() {
	err := t.f()
	if err != nil {
		log.Errorf("task exec fail:%v", err)
	}
}

// Pool  a pool
type Pool struct {
	workerNum  int
	EntryChan  chan *Task
	workerChan chan *Task
}

// NewPool new pool
func NewPool(num int) *Pool {
	return &Pool{
		workerNum:  num,
		EntryChan:  make(chan *Task),
		workerChan: make(chan *Task),
	}
}

// task exe
func (p *Pool) worker() {
	for task := range p.workerChan {
		task.execute()
	}
}

// Run a work
func (p *Pool) Run() {
	for i := 0; i < p.workerNum; i++ {
		go p.worker()
	}
	for task := range p.EntryChan {
		p.workerChan <- task
	}
}

// 敏感词过滤
//
//nolint:gocyclo
func filterTxAndEvent(transactions map[string]*db.Transaction, events []*db.ContractEvent) time.Duration {
	startTime := time.Now()
	if !sensitive.GetSensitiveEnable() {
		return time.Since(startTime)
	}
	var goRoutinePool *ants.Pool
	var err error
	if goRoutinePool, err = ants.NewPool(10, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return time.Since(startTime)
	}
	defer goRoutinePool.Release()

	var txIdEventMp sync.Map
	for i := 0; i < len(events); i++ {
		tempEvent := events[i]
		arr, arrOk := txIdEventMp.Load(tempEvent.TxId)
		if !arrOk {
			arr = []*db.ContractEvent{}
		}
		arry, ok := arr.([]*db.ContractEvent)
		if !ok {
			continue
		}
		arr = append(arry, tempEvent)
		txIdEventMp.Store(tempEvent.TxId, arr)
	}

	wg := sync.WaitGroup{}

	for i := range transactions {
		transaction := transactions[i]
		wg.Add(1)
		// contractName method
		_ = goRoutinePool.Submit(func() {
			defer wg.Done()
			var strBuiler strings.Builder
			eventArr, ok := txIdEventMp.Load(transaction.TxId)
			var txEvents []*db.ContractEvent
			if ok {
				txEvents, _ = eventArr.([]*db.ContractEvent)
			}
			// content add event data
			for _, eve := range txEvents {
				strBuiler.WriteString(eve.Topic + " " + eve.EventData + " ")
			}
			// add contract name and method name
			strBuiler.WriteString(transaction.ContractName + " " + transaction.ContractMethod + " ")
			// add contractResult,contractMessage
			strBuiler.WriteString(string(transaction.ContractResult) + " " + transaction.ContractMessage + " ")
			// add read write set
			var readParameters []*config.RwSet
			var writeParameters []*config.RwSet
			if len(transaction.ReadSet) > 0 {
				err = json.Unmarshal([]byte(transaction.ReadSet), &readParameters)
				if err != nil {
					log.Warn("Contract ReadSet Unmarshal Failed: " + err.Error())
				}
				for i, parameter := range readParameters {
					if i != 0 {
						strBuiler.WriteString(" ")
					}
					strBuiler.WriteString(parameter.Value)
				}
			}
			// add write set
			if len(transaction.WriteSet) > 0 {
				err = json.Unmarshal([]byte(transaction.WriteSet), &writeParameters)
				if err != nil {
					log.Warn("Contract WriteSet Unmarshal Failed: " + err.Error())
				}
				for i, parameter := range writeParameters {
					if i != 0 {
						strBuiler.WriteString(" ")
					}
					strBuiler.WriteString(parameter.Value)
				}
			}
			// add parameters
			var parameters []*pbCommon.KeyValuePair
			if len(transaction.ContractParameters) > 0 {

				errJson := json.Unmarshal([]byte(transaction.ContractParameters), &parameters)
				if errJson != nil {
					log.Warn("Contract parameters Unmarshal Failed: " + errJson.Error())
				}
				for i, parameter := range parameters {
					if i != 0 {
						strBuiler.WriteString(" ")
					}
					strBuiler.WriteString(string(parameter.Value))
				}
			}

			_, flag := FilteringSensitive(strBuiler.String())
			//log.Infof("transaction %s , got result %+v ", strBuiler.String(), flag)
			if !flag {
				return
			}
			if flag {
				transaction.ContractNameBak = transaction.ContractName
				transaction.ContractName = config.ContractWarnMsg
				transaction.ContractMethod = config.ContractWarnMsg
				// result
				transaction.ContractResultBak = make([]byte, len(transaction.ContractResult))
				copy(transaction.ContractResultBak, transaction.ContractResult)
				transaction.ContractResult = []byte(config.OtherWarnMsg)
				transaction.ContractMessageBak = transaction.ContractMessage
				transaction.ContractMessage = config.OtherWarnMsg
				// read set and write set
				transaction.WriteSetBak = transaction.WriteSet
				for _, parameter := range writeParameters {
					parameter.Key = config.OtherWarnMsg
					parameter.Value = config.OtherWarnMsg
				}
				writeParametersBytesBak, paramErr := json.Marshal(writeParameters)
				if paramErr != nil {
					log.Error("ParametersBak Marshal Failed: " + err.Error())
					return
				}
				transaction.WriteSet = string(writeParametersBytesBak)

				transaction.ReadSetBak = transaction.ReadSet
				for _, parameter := range readParameters {
					parameter.Key = config.OtherWarnMsg
					parameter.Value = config.OtherWarnMsg
				}
				readParametersBytesBak, paramErr := json.Marshal(readParameters)
				if paramErr != nil {
					log.Error("Contract ParametersBak Marshal Failed: " + err.Error())
					return
				}
				transaction.ReadSet = string(readParametersBytesBak)
				// parameters
				transaction.ContractParametersBak = transaction.ContractParameters
				for _, parameter := range parameters {
					parameter.Value = []byte(config.OtherWarnMsg)
				}
				parametersBytesBak, paramErr := json.Marshal(parameters)
				if paramErr != nil {
					log.Error("parameters Marshal Failed: " + err.Error())
					return
				}
				transaction.ContractParameters = string(parametersBytesBak)

				// events
				for _, e := range txEvents {
					e.Topic = config.OtherWarnMsg
					e.TopicBak = e.Topic
					e.EventData = config.OtherWarnMsg
					e.EventDataBak = e.EventData
				}

			}
		})
	}
	wg.Wait()
	return time.Since(startTime)
}

// FilteringSensitive 过滤敏感词
func FilteringSensitive(input string) (retKeyWords []string, flag bool) {
	flag = false
	// 敏感词设置
	if !sensitive.GetSensitiveEnable() {
		return
	}
	// 敏感词过滤，flag为true表示包含敏感词
	if len(input) == 0 {
		return
	}

	num := (len(input) + FilterCapacity - 1) / FilterCapacity
	var contentStr []string
	for i := 0; i < num; i++ {
		tmp := input[i*FilterCapacity : Min((i+1)*FilterCapacity-1, len(input))]
		contentStr = append(contentStr, tmp)
	}

	for _, content := range contentStr {
		credential := common.NewCredential(config.GlobalConfig.SensitiveConf.SecretId,
			config.GlobalConfig.SensitiveConf.SecretKey)
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "tms.tencentcloudapi.com"
		client, _ := tms.NewClient(credential, regions.Beijing, cpf)
		//忽略服务器证书校验
		// nolint
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.WithHttpTransport(tr)
		request := tms.NewTextModerationRequest()
		//过滤不可见字符
		content = FilterContent(content)
		if content == "" {
			continue
		}

		contentByte := []byte(content)
		// base64编码
		encodeString := base64.StdEncoding.EncodeToString(contentByte)
		request.Content = &encodeString
		request.BizType = &contentBizType
		response, err := client.TextModeration(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			sensitiveCallTotal.WithLabelValues(callSensitiveFailed, callSensitiveNotHit).Inc()
			log.Errorf("An API error has returned: %s \n content:%s", err, content)
			return
		}
		if err != nil {
			sensitiveCallTotal.WithLabelValues(callSensitiveFailed, callSensitiveNotHit).Inc()
			return
		}

		sug := *response.Response.Suggestion
		label := *response.Response.Label
		requestId := *response.Response.RequestId
		keywords := response.Response.Keywords
		var keyword []string
		for _, k := range keywords {
			keyword = append(keyword, *k)
			retKeyWords = append(retKeyWords, *k)
		}
		log.Infof("FilteringSensitive content: %s, label: %s: , sug: %s ,requestId: %s", content, label, sug, requestId)
		if sug == Block {
			flag = true
			sensitiveCallTotal.WithLabelValues(callSensitiveSuccess, callSensitiveHit).Inc()
			log.Infof("输入内容包含敏感词，已被系统过滤，过滤内容类型: requestId: %s block: %s  关键字:%v \n", requestId, label, keyword)
			return
		} else if sug == Review {
			sc := *response.Response.Score
			if sc > 70 {
				flag = true
				sensitiveCallTotal.WithLabelValues(callSensitiveSuccess, callSensitiveHit).Inc()
				log.Infof("输入内容包含敏感词，已被系统过滤，过滤内容类型: requestId: %s Review: %s  关键字:%v Score:%d\n", requestId, label, keywords, sc)
				return
			}
		}
		sensitiveCallTotal.WithLabelValues(callSensitiveSuccess, callSensitiveNotHit).Inc()
	}
	return
}

// FilterContent filter
func FilterContent(originStr string) string {
	reg, err := regexp.Compile("[^\u4e00-\u9fa5_a-zA-Z0-9\\s\\n]+")
	if err != nil {
		fmt.Println(err)
	}
	processedString := reg.ReplaceAllString(originStr, "")
	return processedString
}

// Min min
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// SingletonSync singleton
type SingletonSync struct {
	SyncStart bool
}

var singlesync *SingletonSync

// GetSyncStart get
// @desc
// @param ${param}
// @return bool
func (singletonSync *SingletonSync) GetSyncStart() bool {
	return singlesync.SyncStart
}

// SetSyncStart s
// @desc
// @param ${param}
func (singletonSync *SingletonSync) SetSyncStart(syncStart bool) {
	singlesync.SyncStart = syncStart
}

type APIResponse struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Err  string      `json:"err"`
}

// BlockHeightExtractNumber 处理BlockHeight
func BlockHeightExtractNumber(str string) int64 {
	trimmed := strings.TrimPrefix(str, "bh")
	if trimmed == "" {
		return 0
	}

	num, err := strconv.ParseInt(trimmed, 10, 64)
	if err == nil {
		return num
	}
	return 0
}

// GatewayIdExtractNumber 处理GatewayId
func GatewayIdExtractNumber(str string) int64 {
	trimmed := strings.TrimLeft(strings.TrimPrefix(str, "g"), "0")
	if trimmed == "" {
		return 0
	}

	gatewayId, err := strconv.ParseInt(trimmed, 10, 64)
	if err == nil {
		return gatewayId
	}
	return 0
}
