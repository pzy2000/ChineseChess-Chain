package relayCrossChain

type GetGatewayInfoReq struct {
	GatewayId string `json:"gatewayId" form:"gatewayId"`
}

type GetSubChainInfoReq struct {
	SubChainId string `json:"subChainId" form:"subChainId"`
}
