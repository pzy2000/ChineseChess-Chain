package utils

import (
	"bytes"
	"chainmaker_web/src/config"
	loggers "chainmaker_web/src/logger"
	"chainmaker_web/src/models/relayCrossChain"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

var log = loggers.GetLogger(loggers.MODULE_WEB)

func GetRelayCrossChainHttpResp(params interface{}, action string) ([]byte, error) {
	crossChainUrl := config.GlobalConfig.WebConf.RelayCrossChainUrl
	url := crossChainUrl + action
	jsonByte, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if resp.Body != nil {
			err = resp.Body.Close()
			if err != nil {
				return
			}
		}
	}()

	if resp.StatusCode != 200 {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	log.Infof("【http service】get relay cross log, params:%v, url:%v, respJson:%v, err:%v",
		string(jsonByte), url, string(body), err)
	return body, err
}

// GetCrossGatewayInfo
//
//	@Description: 根据gatewayId获取子链列表
//	@param gatewayId 网关id
//	@return *relayCrossChain.GatewayInfoData 子链列表
//	@return error
func GetCrossGatewayInfo(gatewayId int64) (*relayCrossChain.GatewayInfoData, error) {
	params := relayCrossChain.GetGatewayInfoReq{
		GatewayId: strconv.FormatInt(gatewayId, 10),
	}
	body, err := GetRelayCrossChainHttpResp(params, config.CrossGateWayIdUrl)
	if err != nil {
		log.Errorf("cross chain GetCrossGatewayInfo err:%v", err)
		return nil, err
	}
	var respJson relayCrossChain.GetGatewayInfoResponse
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Errorf("cross chain GetCrossGatewayInfo err:%v", err)
		return nil, err
	}

	log.Infof("cross chain GetCrossGatewayInfo GatewayId:%v, result:%v", params.GatewayId, string(body))
	return respJson.Data, nil
}

// GetCrossSubChainInfo
//
//	@Description: 根据子链id获取子链信息
//	@param subChainId
//	@return *relayCrossChain.SubChainInfoData
//	@return error
func GetCrossSubChainInfo(subChainId string) (*relayCrossChain.SubChainInfoData, error) {
	params := relayCrossChain.GetSubChainInfoReq{
		SubChainId: subChainId,
	}
	body, err := GetRelayCrossChainHttpResp(params, config.CrossSubChainInfoUrl)
	if err != nil {
		return nil, err
	}
	var respJson relayCrossChain.GetSubChainInfoResponse
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		return nil, err
	}

	return respJson.Data, nil
}
