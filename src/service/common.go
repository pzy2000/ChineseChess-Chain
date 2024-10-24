package service

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"path"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// StringAmountDecimal  string转decimal
func StringAmountDecimal(amount string, decimals int) decimal.Decimal {
	// 将字符串转换为 decimal.Decimal 值
	amountDecimal, _ := decimal.NewFromString(amount)
	// 创建一个新的 decimal.Decimal 值，表示 10^decimals
	divisor := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))
	// 使用 Div 方法将 amountDecimal 除以 divisor
	resultDecimal := amountDecimal.Div(divisor)

	return resultDecimal
}

// GetCurrentMonthStartAndEndTime 自然月的开始结束时间
func GetCurrentMonthStartAndEndTime() (int64, int64) {
	now := time.Now()

	// 获取当前月份的第一天
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 获取下个月的第一天
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

	// 获取当前月份的最后一天（下个月的第一天减去一秒）
	endOfMonth := startOfNextMonth.Add(-time.Second)

	// 返回 Unix 时间戳（以秒为单位）
	return startOfMonth.Unix(), endOfMonth.Unix()
}

// GetAccountBNS
//
//	@Description: 根据合约地址获取账户信息
//	@param address
//	@param accountMap
//	@return string
func GetAccountBNS(address string, accountMap map[string]*db.Account) string {
	var addrBns string
	if account, ok := accountMap[address]; ok {
		addrBns = account.BNS
	}
	return addrBns
}

// GetAccountBNSByAddr
//
//	@Description: 根据账户地址获取BNS
//	@param chainId
//	@param address 账户地址
//	@return string BNS
func GetAccountBNSByAddr(chainId, address string) string {
	accountInfo, _ := dbhandle.GetAccountByAddr(chainId, address)
	var addrBns string
	if accountInfo != nil {
		addrBns = accountInfo.BNS
	}
	return addrBns
}

// GetContractAccountMap
//
//	@Description: 根据合约地址获取账户列表
//	@param chainId
//	@param contractList
//	@return map[string]*db.Account
func GetContractAccountMap(chainId string, contractList []*db.Contract) map[string]*db.Account {
	var accountAddrs []string
	for _, contract := range contractList {
		accountAddrs = append(accountAddrs, contract.CreatorAddr)
	}
	accountMap, _ := dbhandle.QueryAccountExists(chainId, accountAddrs)
	return accountMap
}

// isImageOrVideo
//
//	@Description: 判断token地址是图片还是url
//	@param url
//	@return string
func isImageOrVideo(url string) string {
	fileExt := strings.ToLower(path.Ext(url))

	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	videoExtensions := []string{".mp4", ".mkv", ".flv", ".avi", ".mov", ".wmv"}

	for _, ext := range imageExtensions {
		if fileExt == ext {
			return UrlTypeImage
		}
	}

	for _, ext := range videoExtensions {
		if fileExt == ext {
			return UrlTypeVideo
		}
	}

	return UrlTypeImage
}

type TxListViewSlice []*entity.TxListView

func (t TxListViewSlice) Len() int {
	return len(t)
}

func (t TxListViewSlice) Less(i, j int) bool {
	if t[i].BlockHeight != t[j].BlockHeight {
		return t[i].BlockHeight > t[j].BlockHeight
	}
	return t[i].Timestamp > t[j].Timestamp
}

func (t TxListViewSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type ContractTxListViewSlice []*entity.ContractTxListView

func (t ContractTxListViewSlice) Len() int {
	return len(t)
}

func (t ContractTxListViewSlice) Less(i, j int) bool {
	if t[i].BlockHeight != t[j].BlockHeight {
		return t[i].BlockHeight > t[j].BlockHeight
	}
	return t[i].Timestamp > t[j].Timestamp
}

func (t ContractTxListViewSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func GetSubChainIdByName(chainId, subChainName string) (string, error) {
	if subChainName == "" {
		return "", nil
	}

	mainChainName := config.GlobalConfig.ChainConf.MainChainName
	mainChainId := config.GlobalConfig.ChainConf.MainChainId
	if subChainName == mainChainName {
		return mainChainId, nil
	}

	subChainInfo, err := dbhandle.GetCrossSubChainInfoByName(chainId, subChainName)
	if err != nil {
		log.Errorf("GetSubChainIdByName err : %v", err)
		return "", err
	}

	if subChainInfo == nil {
		return "", nil
	}

	return subChainInfo.SubChainId, nil
}
