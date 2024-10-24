package service

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
)

// PositionMaxRank 持仓排名最大数量
const PositionMaxRank = 100000

// ContractPositionRank 连表查询持仓和排名
type ContractPositionRank struct {
	OwnerAddr string
	Amount    decimal.Decimal
	BNS       string
	AddrType  int
	HoldRank  int64
}

// GetFTPositionListHandler get
type GetFTPositionListHandler struct{}

func (GetFTPositionListHandler *GetFTPositionListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetFungiblePositionListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//根据账户查询持仓列表
	if params.OwnerAddr != "" && params.ContractAddr == "" {
		userFTPositionListHandler := GetUserFTPositionListHandler{}
		userFTPositionListHandler.Handle(ctx)
		return
	}

	if params.ContractAddr == "" {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	fungibleContract, err := dbhandle.GetFungibleContractByAddr(params.ChainId, params.ContractAddr)
	if err != nil || fungibleContract == nil {
		log.Errorf("GetFungibleContract err : %s", err)
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	positionList, totalCount, err := getPositionListAndTotalCount(params.Offset, params.Limit, params.ChainId,
		params.ContractAddr, params.OwnerAddr)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	positionListView := arraylist.New()
	for _, position := range positionList {
		holdRatio := position.Amount.Div(fungibleContract.TotalSupply).Mul(decimal.NewFromInt(100))
		holdRatioStr := holdRatio.Round(4).String() + "%"

		positionView := &entity.PositionListView{
			ContractName:   fungibleContract.ContractName,
			ContractSymbol: fungibleContract.Symbol,
			ContractAddr:   fungibleContract.ContractAddr,
			ContractType:   fungibleContract.ContractType,
			AddrType:       position.AddrType,
			OwnerAddr:      position.OwnerAddr,
			OwnerAddrBNS:   position.BNS,
			Amount:         position.Amount.String(),
			HoldRatio:      holdRatioStr,
			HoldRank:       position.HoldRank,
		}
		positionListView.Add(positionView)
	}

	ConvergeListResponse(ctx, positionListView.Values(), totalCount, nil)
}

// getPositionListAndTotalCount
//
//	@Description: 获取排行榜列表和数据总量
//	@param offset
//	@param limit
//	@param chainId
//	@param contractAddr
//	@param ownerAddr
//	@return []*ContractPositionRank
//	@return int64
//	@return error
func getPositionListAndTotalCount(offset, limit int, chainId, contractAddr, ownerAddr string) ([]*ContractPositionRank,
	int64, error) {
	//持仓总数
	var (
		totalCount        int64
		isUpdateRankCache bool
		positionRankList  = make([]*ContractPositionRank, 0)
	)

	if ownerAddr != "" {
		accountInfo, err := dbhandle.GetAccountByAddr(chainId, ownerAddr)
		if err != nil || accountInfo == nil {
			return positionRankList, totalCount, err
		}

		positionList, err := dbhandle.GetFungiblePositionList(offset, limit, chainId, contractAddr, ownerAddr)
		if err != nil {
			return positionRankList, totalCount, err
		}

		totalCount = int64(len(positionList))
		holdRank := GetPositionRankCacheByAddr(chainId, contractAddr, ownerAddr)
		if holdRank == 0 {
			isUpdateRankCache = true
		}
		for _, position := range positionList {
			positionRankList = append(positionRankList, &ContractPositionRank{
				OwnerAddr: position.OwnerAddr,
				Amount:    position.Amount,
				BNS:       accountInfo.BNS,
				AddrType:  accountInfo.AddrType,
				HoldRank:  holdRank,
			})
		}
	} else {
		//缓存获取总数据量
		totalCount = GetContractPositionOwnerCount(chainId, contractAddr)
		if totalCount == 0 {
			isUpdateRankCache = true
			//缓存失效，从数据库获取
			positionList, err := dbhandle.GetFTPositionJoinAccount(offset, limit, chainId, contractAddr, ownerAddr)
			if err != nil {
				return positionRankList, totalCount, err
			}

			if len(positionList) > 0 {
				totalCount = int64(offset*limit + len(positionList))
			}
			for i, position := range positionList {
				holdRank := int64(offset*limit + i + 1)
				positionRankList = append(positionRankList, &ContractPositionRank{
					OwnerAddr: position.OwnerAddr,
					Amount:    position.Amount,
					BNS:       position.BNS,
					AddrType:  position.AddrType,
					HoldRank:  holdRank,
				})
			}
		} else {
			//缓存未失效，从缓存获取分页排名数据
			rankOwnerMap, ownerList, err := GetPositionRankMapCache(offset, limit, chainId, contractAddr)
			if err != nil {
				return positionRankList, totalCount, err
			}

			//根据地址列表获取详细数据
			positionList, err := dbhandle.GetFTPositionByAddrJoinAccount(chainId, contractAddr, ownerList)
			if err != nil {
				return positionRankList, totalCount, err
			}

			positionMap := make(map[string]*db.ContractPositionAccount, len(positionList))
			for _, position := range positionList {
				positionMap[position.OwnerAddr] = position
			}

			for _, address := range ownerList {
				position, ok := positionMap[address]
				if !ok {
					continue
				}

				positionRankList = append(positionRankList, &ContractPositionRank{
					OwnerAddr: position.OwnerAddr,
					Amount:    position.Amount,
					BNS:       position.BNS,
					AddrType:  position.AddrType,
					HoldRank:  rankOwnerMap[position.OwnerAddr],
				})
			}
		}
	}

	if isUpdateRankCache {
		//异步更新排行榜缓存
		go UpdatePositionRankListCache(chainId, contractAddr)
	}

	return positionRankList, totalCount, nil
}

// UpdatePositionRankListCache
//
//	@Description: 异步更新-持仓排行榜缓存
//	@param chainId
//	@param contractAddr
func UpdatePositionRankListCache(chainId, contractAddr string) {
	var (
		offset        int
		limit         = 1000
		goRoutinePool *ants.Pool
		err           error
		wg            sync.WaitGroup
	)

	startTime := time.Now()
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error("UpdatePositionRankListCache new ants pool error:" + err.Error())
		return
	}

	defer goRoutinePool.Release()
	// 原子计数器，用于存储已完成任务的数量
	var completedTasks int32
	// 原子布尔变量，用于存储是否已翻到最后一页
	var isLastPage int32
	for {
		wg.Add(1)
		errSub := goRoutinePool.Submit(func(offset int) func() {
			return func() {
				defer wg.Done()
				//插入数据
				positionList, _ := dbhandle.GetFungiblePositionList(offset, limit, chainId, contractAddr, "")
				// 写入缓存
				dbhandle.SetFTPositionListCache(chainId, contractAddr, positionList)
				// 更新已完成任务的数量
				atomic.AddInt32(&completedTasks, int32(len(positionList)))
				// 检查是否已翻到最后一页
				if len(positionList) < limit || positionList == nil {
					atomic.StoreInt32(&isLastPage, 1)
				}
			}
		}(offset))
		if errSub != nil {
			log.Error("UpdatePositionRankList submit Failed : " + errSub.Error())
		}

		// 如果已翻到最后一页，或已完成的任务数量达到最大排名，退出循环
		if atomic.LoadInt32(&isLastPage) == 1 ||
			int(atomic.LoadInt32(&completedTasks)) >= PositionMaxRank {
			break
		}
		offset++
	}
	wg.Wait()
	log.Infof("【redis】set redis success, key:contract_position_owner_list total:%v, duration_time:%vms",
		completedTasks, time.Since(startTime).Milliseconds())
}

// UpdateNFTPositionRankListCache
//
//	@Description: 异步更新-持仓排行榜缓存
//	@param chainId
//	@param contractAddr
func UpdateNFTPositionRankListCache(chainId, contractAddr string) {
	var (
		offset        int
		limit         = 1000
		goRoutinePool *ants.Pool
		err           error
		wg            sync.WaitGroup
	)

	startTime := time.Now()
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error("UpdateNFTPositionRankListCache new ants pool error:" + err.Error())
		return
	}

	defer goRoutinePool.Release()
	// 原子计数器，用于存储已完成任务的数量
	var completedTasks int32
	// 原子布尔变量，用于存储是否已翻到最后一页
	var isLastPage int32
	for {
		wg.Add(1)
		errSub := goRoutinePool.Submit(func(offset int) func() {
			return func() {
				defer wg.Done()
				//插入数据
				positionList, _ := dbhandle.GetNFTPositionList(offset, limit, chainId, contractAddr, "")
				// 写入缓存
				dbhandle.SetNFTPositionListCache(chainId, contractAddr, positionList)
				// 更新已完成任务的数量
				atomic.AddInt32(&completedTasks, int32(len(positionList)))
				// 检查是否已翻到最后一页
				if len(positionList) < limit || positionList == nil {
					atomic.StoreInt32(&isLastPage, 1)
				}
			}
		}(offset))
		if errSub != nil {
			log.Error("InsertFungibleTransfer submit Failed : " + errSub.Error())
		}

		// 如果已翻到最后一页，或已完成的任务数量达到最大排名，退出循环
		if atomic.LoadInt32(&isLastPage) == 1 ||
			int(atomic.LoadInt32(&completedTasks)) >= PositionMaxRank {
			break
		}
		offset++
	}

	wg.Wait()
	log.Infof("【redis】set redis success, key:contract_position_owner_list total:%v, duration_time:%vms",
		completedTasks, time.Since(startTime).Milliseconds())
}

// GetUserFTPositionListHandler get
type GetUserFTPositionListHandler struct{}

func (GetUserFTPositionListHandler *GetUserFTPositionListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetUserFTPositionListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取地址BNS
	accountInfo, err := dbhandle.GetAccountByAddr(params.ChainId, params.OwnerAddr)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}
	if accountInfo == nil {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
	}

	//持仓数据
	totalCount, err := dbhandle.GetFTPositionCountByAddr(params.ChainId, params.OwnerAddr)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}
	if totalCount == 0 {
		ConvergeListResponse(ctx, []interface{}{}, totalCount, nil)
		return
	}

	positionList, err := dbhandle.GetFTPositionListByAddr(params.Offset, params.Limit, params.ChainId, params.OwnerAddr)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	positionListView := arraylist.New()
	for _, position := range positionList {
		positionView := &entity.PositionListView{
			ContractName:   position.ContractName,
			ContractSymbol: position.Symbol,
			ContractAddr:   position.ContractAddr,
			ContractType:   position.ContractType,
			AddrType:       accountInfo.AddrType,
			OwnerAddr:      position.OwnerAddr,
			OwnerAddrBNS:   accountInfo.BNS,
			Amount:         position.Amount.String(),
		}
		positionListView.Add(positionView)
	}
	ConvergeListResponse(ctx, positionListView.Values(), totalCount, nil)
}

// GetNonFungiblePositionListHandler get
type GetNonFungiblePositionListHandler struct{}

// Handle deal
func (getNonFungiblePositionListHandler *GetNonFungiblePositionListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetNonFungiblePositionListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	nftContract, err := dbhandle.GetNonFungibleContractByAddr(params.ChainId, params.ContractAddr)
	if err != nil || nftContract == nil {
		log.Errorf("getNonFungiblePositionListHandler GetNonFungibleContract err : %s", err)
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	positionList, totalCount, err := getNFTPositionListAndTotalCount(params.Offset, params.Limit, params.ChainId,
		params.ContractAddr, params.OwnerAddr)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	positionListView := arraylist.New()
	for _, position := range positionList {
		holdRatio := position.Amount.Div(nftContract.TotalSupply).Mul(decimal.NewFromInt(100))
		holdRatioStr := holdRatio.Round(4).String() + "%"

		positionView := &entity.NonPositionListView{
			ContractName: nftContract.ContractName,
			ContractAddr: nftContract.ContractAddr,
			AddrType:     position.AddrType,
			OwnerAddr:    position.OwnerAddr,
			OwnerAddrBNS: position.BNS,
			Amount:       position.Amount.String(),
			HoldRatio:    holdRatioStr,
			HoldRank:     position.HoldRank,
		}
		positionListView.Add(positionView)
	}
	ConvergeListResponse(ctx, positionListView.Values(), totalCount, nil)
}

// getNFTPositionListAndTotalCount
//
//	@Description: 获取排行榜列表和数据总量
//	@param offset
//	@param limit
//	@param chainId
//	@param contractAddr
//	@param ownerAddr
//	@return []*ContractPositionRank
//	@return int64
//	@return error
func getNFTPositionListAndTotalCount(offset, limit int, chainId, contractAddr, ownerAddr string) (
	[]*ContractPositionRank, int64, error) {
	//持仓总数
	var (
		totalCount        int64
		isUpdateRankCache bool
		positionRankList  = make([]*ContractPositionRank, 0)
	)

	if ownerAddr != "" {
		accountInfo, err := dbhandle.GetAccountByAddr(chainId, ownerAddr)
		if err != nil || accountInfo == nil {
			return positionRankList, totalCount, err
		}

		positionList, err := dbhandle.GetNFTPositionList(offset, limit, chainId, contractAddr, ownerAddr)
		if err != nil {
			return positionRankList, totalCount, err
		}

		totalCount = int64(len(positionList))
		holdRank := GetPositionRankCacheByAddr(chainId, contractAddr, ownerAddr)
		if holdRank == 0 {
			isUpdateRankCache = true
		}
		for _, position := range positionList {
			positionRankList = append(positionRankList, &ContractPositionRank{
				OwnerAddr: position.OwnerAddr,
				Amount:    position.Amount,
				BNS:       accountInfo.BNS,
				AddrType:  accountInfo.AddrType,
				HoldRank:  holdRank,
			})
		}
	} else {
		//缓存获取总数据量
		totalCount = GetContractPositionOwnerCount(chainId, contractAddr)
		if totalCount == 0 {
			isUpdateRankCache = true
			//缓存失效，从数据库获取
			positionList, err := dbhandle.GetNFTPositionJoinAccount(offset, limit, chainId, contractAddr, ownerAddr)
			if err != nil {
				return positionRankList, totalCount, err
			}

			if len(positionList) > 0 {
				totalCount = int64(offset*limit + len(positionList))
			}
			for i, position := range positionList {
				holdRank := int64(offset*limit + i + 1)
				positionRankList = append(positionRankList, &ContractPositionRank{
					OwnerAddr: position.OwnerAddr,
					Amount:    position.Amount,
					BNS:       position.BNS,
					AddrType:  position.AddrType,
					HoldRank:  holdRank,
				})
			}
		} else {
			//缓存未失效，从缓存获取分页排名数据
			rankOwnerMap, ownerList, err := GetPositionRankMapCache(offset, limit, chainId, contractAddr)
			if err != nil {
				return positionRankList, totalCount, err
			}

			//根据地址列表获取详细数据
			positionList, err := dbhandle.GetNFTPositionByAddrJoinAccount(chainId, contractAddr, ownerList)
			if err != nil {
				return positionRankList, totalCount, err
			}
			positionMap := make(map[string]*db.ContractPositionAccount, len(positionList))
			for _, position := range positionList {
				positionMap[position.OwnerAddr] = position
			}

			for _, address := range ownerList {
				position, ok := positionMap[address]
				if !ok {
					continue
				}

				positionRankList = append(positionRankList, &ContractPositionRank{
					OwnerAddr: position.OwnerAddr,
					Amount:    position.Amount,
					BNS:       position.BNS,
					AddrType:  position.AddrType,
					HoldRank:  rankOwnerMap[position.OwnerAddr],
				})
			}
		}
	}

	if isUpdateRankCache {
		//异步更新排行榜缓存
		go UpdateNFTPositionRankListCache(chainId, contractAddr)
	}

	return positionRankList, totalCount, nil
}

// GetPositionAccountMap
//
//	@Description: 根据合约地址获取账户列表
//	@param chainId
//	@param contractList
//	@return map[string]*db.Account
func GetPositionAccountMap(chainId string, positionList []*db.PositionWithRank) map[string]*db.Account {
	var accountAddrs []string
	for _, position := range positionList {
		accountAddrs = append(accountAddrs, position.OwnerAddr)
	}
	accountMap, _ := dbhandle.QueryAccountExists(chainId, accountAddrs)
	return accountMap
}

// GetContractPositionOwnerCount
//
//	@Description:  获取同质化持仓列表总数
//	@param chainId
//	@param contractAddr
//	@return int64
func GetContractPositionOwnerCount(chainId, contractAddr string) int64 {
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisContractPositionOwnerList, prefix, chainId, contractAddr)
	// 查询合约持仓列表总数
	count, err := cache.GlobalRedisDb.ZCard(ctx, redisKey).Result()
	if err != nil {
		panic("failed to get the number of members")
	}

	return count
}

// GetPositionRankMapCache
//
//	@Description: 获取同质化持仓排名列表缓存
//	@param offset
//	@param limit
//	@param chainId
//	@param contractAddr
//	@return map[string]int64
//	@return []string
//	@return error
func GetPositionRankMapCache(offset, limit int, chainId, contractAddr string) (map[string]int64, []string, error) {
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisContractPositionOwnerList, prefix, chainId, contractAddr)
	// 分页查询，每页10条数据，这里查询第1页
	start := offset * limit
	end := start + limit - 1

	rankMap := make(map[string]int64, 0)
	ownerList := make([]string, 0)
	// 使用ZREVRANGE命令按排名由高到低获取数据
	results, err := cache.GlobalRedisDb.ZRevRangeWithScores(ctx, redisKey, int64(start), int64(end)).Result()
	if err != nil {
		return rankMap, ownerList, err
	}

	// 输出查询结果
	for i, result := range results {
		ownerAddr, _ := result.Member.(string)
		rankMap[ownerAddr] = int64(start + i + 1)
		ownerList = append(ownerList, ownerAddr)
	}

	return rankMap, ownerList, nil
}

// GetPositionRankCacheByAddr
//
//	@Description: 获取同质化持仓排名缓存
//	@param chainId
//	@param contractAddr
//	@param ownerAddr
//	@return int64
func GetPositionRankCacheByAddr(chainId, contractAddr, ownerAddr string) int64 {
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisContractPositionOwnerList, prefix, chainId, contractAddr)
	// 获取排名
	rank, err := cache.GlobalRedisDb.ZRevRank(ctx, redisKey, ownerAddr).Result()
	if err != nil {
		log.Errorf("GetPositionRankCacheByAddr ZRevRank err:%v, redisKey:%v, ownerAddr:%v",
			err, redisKey, ownerAddr)
		return 0
	}
	return rank + 1
}
