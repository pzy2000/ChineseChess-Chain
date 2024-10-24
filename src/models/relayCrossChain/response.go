package relayCrossChain

type GetGatewayInfoResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *GatewayInfoData `json:"data"`
}

type GatewayInfoData struct {
	GatewayId    string          `json:"gatewayId"`
	GatewayName  string          `json:"gatewayName"`
	Address      string          `json:"address"`
	TxVerifyType int             `json:"txVerifyType"`
	CrossCa      string          `json:"crossCa"`
	SdkClientCrt string          `json:"sdkClientCrt"`
	SdkClientKey string          `json:"sdkClientKey"`
	Enable       bool            `json:"enable"`
	SubChainInfo []*SubChainInfo `json:"subChainInfo"`
}

type SubChainInfo struct {
	SubChainId        string `json:"subChainId"`
	ChainName         string `json:"chainName"`
	ChainId           string `json:"chainId"`
	ChainType         int32  `json:"chainType"`
	Introduction      string `json:"introduction"`
	Company           string `json:"company"`
	Name              string `json:"name"`
	Job               int    `json:"job"`
	PhoneNumber       string `json:"phoneNumber"`
	WechatNumber      string `json:"wechatNumber"`
	Email             string `json:"email"`
	SpvContractName   string `json:"spvContractName"`
	ExplorerAddress   string `json:"explorerAddress"`
	ExplorerTxAddress string `json:"txInfoExplorerAddress"`
}

type GetSubChainInfoResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    *SubChainInfoData `json:"data"`
}

type SubChainInfoData struct {
	GatewayId       string `json:"gatewayId"`
	GatewayName     string `json:"gatewayName"`
	Address         string `json:"address"`
	TxVerifyType    int    `json:"txVerifyType"`
	CrossCa         string `json:"crossCa"`
	SdkClientCrt    string `json:"sdkClientCrt"`
	SdkClientKey    string `json:"sdkClientKey"`
	Enable          bool   `json:"enable"`
	SubChainId      string `json:"subChainId"`
	ChainName       string `json:"chainName"`
	ChainId         string `json:"chainId"`
	ChainType       int    `json:"chainType"`
	Introduction    string `json:"introduction"`
	Company         string `json:"company"`
	Name            string `json:"name"`
	Job             int    `json:"job"`
	PhoneNumber     string `json:"phoneNumber"`
	WechatNumber    string `json:"wechatNumber"`
	Email           string `json:"email"`
	SpvContractName string `json:"spvContractName"`
	ExplorerAddress string `json:"explorerAddress"`
}
