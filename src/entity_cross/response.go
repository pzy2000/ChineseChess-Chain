package entity_cross

// MainCrossConfig 主子链网配置
type MainCrossConfig struct {
	ShowTag bool
}

// CrossSearchView 主子链网搜索
type CrossSearchView struct {
	Type int
	Data string
}

// OverviewDataView OverviewData
type OverviewDataView struct {
	TotalBlockHeight int64
	ShortestTime     int64
	LongestTime      int64
	AverageTime      int64
	SubChainNum      int64
	TxNum            int64
}

// LatestTxListView latest
type LatestTxListView struct {
	CrossId         string
	FromChainName   string
	FromChainId     string
	FromIsMainChain bool
	ToChainName     string
	ToChainId       string
	ToIsMainChain   bool
	Timestamp       int64
	Status          int32
}

// LatestSubChainListView latest
type LatestSubChainListView struct {
	SubChainId       string
	SubChainName     string
	BlockHeight      int64
	CrossTxNum       int64
	CrossContractNum int64
	Timestamp        int64
	Status           int32
}

// GetTxListView latest
type GetTxListView struct {
	CrossId         string
	FromChainName   string
	FromChainId     string
	FromIsMainChain bool
	ToChainName     string
	ToChainId       string
	ToIsMainChain   bool
	Timestamp       int64
	Status          int32
}

// GetSubChainListView get
type GetSubChainListView struct {
	SubChainId       string
	SubChainName     string
	BlockHeight      int64
	CrossTxNum       int64
	CrossContractNum int64
	Timestamp        int64
	Status           int32
	ExplorerUrl      string
}

// GetCrossSubChainDetailView get
type GetCrossSubChainDetailView struct {
	SubChainId       string
	SubChainName     string
	BlockHeight      int64
	ChainType        int32 //区块链架构（1 长安链，2 fabric，3, bcos， 4,eth，5+ 扩展）
	CrossTxNum       int64
	CrossContractNum int64
	Timestamp        int64
	Status           int32
	GatewayId        string
	GatewayName      string
	GatewayAddr      string
}

// GetCrossTxDetailView get
type GetCrossTxDetailView struct {
	CrossId        string
	Status         int32
	CrossDuration  int64
	ContractName   string
	ContractMethod string
	Parameter      string
	ContractResult string
	CrossDirection *CrossDirection
	FromChainInfo  *TxChainInfo
	ToChainInfo    *TxChainInfo
	Timestamp      int64
}

type CrossDirection struct {
	FromChain string
	ToChain   string
}

type TxChainInfo struct {
	ChainName    string
	ChainId      string
	ContractName string
	IsMainChain  bool
	TxId         string
	TxStatus     int32
	TxUrl        string
	Gas          string
}

// GetSubChainCrossView get
type GetSubChainCrossView struct {
	ChainId   string
	ChainName string
	TxNum     int64
}
