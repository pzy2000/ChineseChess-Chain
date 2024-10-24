/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/utils"
	"context"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"chainmaker.org/chainmaker/common/v2/ca"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	"chainmaker.org/chainmaker/common/v2/evmutils"
	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	tcipApi "chainmaker.org/chainmaker/tcip-go/v2/api"
	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	commonutils "chainmaker.org/chainmaker/utils/v2"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gogo/protobuf/proto"
)

// MemberAddrIdCert MemberAddrIdCert
type MemberAddrIdCert struct {
	UserAddr             string
	UserId               string
	HasCert              bool
	CertCommonName       string
	CertOrganization     []string
	CertOrganizationUnit []string
}

func getConnection(destGatewayInfo *tcipCommon.GatewayInfo, sdkClientKey string) (tcipApi.RpcCrossChainClient,
	*grpc.ClientConn, error) {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}
	var tlsClient ca.CAClient
	var (
		err error
	)

	tlsClient = ca.CAClient{
		ServerName: destGatewayInfo.ServerName,
		CaCerts:    []string{destGatewayInfo.Tlsca},
		CertBytes:  []byte(destGatewayInfo.ClientCert),
		KeyBytes:   []byte(sdkClientKey),
	}

	c, err := tlsClient.GetCredentialsByCA()
	if err != nil {
		return nil, nil, err
	}
	conn, err := grpc.Dial(
		destGatewayInfo.Address,
		grpc.WithTransportCredentials(*c),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024),
			grpc.MaxCallSendMsgSize(1024*1024),
		),
		grpc.WithKeepaliveParams(kacp),
	)
	if err != nil {
		return nil, nil, err
	}
	return tcipApi.NewRpcCrossChainClient(conn), conn, nil
}

// CheckSubChainStatus 检查子链健康状态
func CheckSubChainStatus(subChainInfo *db.CrossSubChainData) (bool, error) {
	ctx := context.Background()
	destGatewayInfo := &tcipCommon.GatewayInfo{
		ServerName: "chainmaker.org",
		Tlsca:      subChainInfo.CrossCa,
		ClientCert: subChainInfo.SdkClientCrt,
		Address:    subChainInfo.GatewayAddr,
	}
	client, conn, err := getConnection(destGatewayInfo, subChainInfo.SdkClientKey)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	res, err := client.PingPong(ctx, &emptypb.Empty{})
	if res == nil {
		return false, err
	}
	return res.ChainOk, err
}

// GetTransfer GetTransfer
func GetTransfer(chainId string, contractName, tokenId string) []byte {
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return []byte{}
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return []byte{}
	}
	client := sdkClient.queryClient
	txResponse, err := client.QueryContract(contractName, "TokenMetadata",
		[]*common.KeyValuePair{
			{
				Key:   "tokenId",
				Value: []byte(tokenId),
			},
		}, -1)
	if err != nil || txResponse == nil {
		return []byte{}
	}
	if txResponse.Code == common.TxStatusCode_SUCCESS {
		return txResponse.ContractResult.Result
	}

	return []byte{}
}

// GetOwner GetOwner
func GetOwner(chainId string, contractName, tokenId string) string {
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return ""
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return ""
	}
	client := sdkClient.queryClient

	txResponse, err := client.QueryContract(contractName, "OwnerOf",
		[]*common.KeyValuePair{
			{
				Key:   "tokenId",
				Value: []byte(tokenId),
			},
		}, -1)
	if err != nil || txResponse == nil {
		return ""
	}

	if txResponse.Code == common.TxStatusCode_SUCCESS {
		return string(txResponse.ContractResult.Result)
	}
	return ""
}

// DockerGetTotalSupply totalSupply
func DockerGetTotalSupply(chainId string, name string) (string, error) {
	totalSupply := "0"
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return totalSupply, fmt.Errorf("sdkClients failed")
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return totalSupply, fmt.Errorf("sdkClients failed")
	}
	client := sdkClient.queryClient
	txResponse, err := client.QueryContract(name, "TotalSupply",
		[]*common.KeyValuePair{}, -1)
	if err != nil || txResponse == nil {
		return totalSupply, fmt.Errorf("totalSupply QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	if txResponse.Code == common.TxStatusCode_SUCCESS {
		totalSupply = string(txResponse.ContractResult.Result)
	}

	return totalSupply, err
}

// EvmGetTotalSupply totalSupply
func EvmGetTotalSupply(evmType, chainId, address string) (string, error) {
	totalSupply := "0"
	var ercAbi *abi.ABI
	//ercAbi, err := parseAbiJson(evmType)
	if evmType == ContractStandardNameEVMDFA {
		ercAbi = config.GlobalAbiERC20
	} else if evmType == ContractStandardNameEVMNFA {
		ercAbi = config.GlobalAbiERC721
	}

	if ercAbi == nil {
		return totalSupply, nil
	}

	dataByte, err := ercAbi.Pack("totalSupply")
	if err != nil {
		return totalSupply, fmt.Errorf("ercAbi Pack err: %s", err)
	}
	dataString := hex.EncodeToString(dataByte)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return totalSupply, fmt.Errorf("sdkClients failed: %s", evmType)
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return totalSupply, fmt.Errorf("sdkClients failed: %s", evmType)
	}
	client := sdkClient.queryClient
	txResponse, err := client.QueryContract(address, "totalSupply", kvs, -1)
	if err != nil || txResponse == nil {
		return totalSupply, fmt.Errorf("totalSupply QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	if txResponse.Code != common.TxStatusCode_SUCCESS {
		return totalSupply, fmt.Errorf("totalSupply QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	total := big.NewInt(0)
	output := txResponse.ContractResult.Result
	err = ercAbi.UnpackIntoInterface(&total, "totalSupply", output)
	return total.String(), err
}

// GetBalanceOf 代币持有量
func GetBalanceOf(evmType, chainId, ownerAddr, contractAddr string) (string, error) {
	var balanceStr string
	var err error
	if evmType == ContractStandardNameCMDFA ||
		evmType == ContractStandardNameCMNFA {
		balanceStr, err = DockerGetBalanceOf(evmType, chainId, ownerAddr, contractAddr)
		if err != nil {
			return "", err
		}
	} else if evmType == ContractStandardNameEVMDFA ||
		evmType == ContractStandardNameEVMNFA {
		balanceStr, err = EvmGetBalanceOf(evmType, chainId, ownerAddr, contractAddr)
		if err != nil {
			return "", err
		}
	}

	return balanceStr, nil
}

// DockerGetBalanceOf 代币持有量
func DockerGetBalanceOf(evmType, chainId, ownerAddr, contractAddr string) (string, error) {
	var balanceStr string
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return balanceStr, fmt.Errorf("sdkClients failed: %s", evmType)
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return balanceStr, fmt.Errorf("sdkClients failed: %s", evmType)
	}
	client := sdkClient.queryClient
	txResponse, err := client.QueryContract(contractAddr, "balanceOf", []*common.KeyValuePair{}, -1)
	if err != nil || txResponse == nil {
		return balanceStr, fmt.Errorf("balanceOf QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	if txResponse.Code != common.TxStatusCode_SUCCESS {
		return balanceStr, fmt.Errorf("balanceOf QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	balanceStr = string(txResponse.ContractResult.Result)
	return balanceStr, err
}

// DockerGetDecimals 发行小数位数
func DockerGetDecimals(chainId, contractAddr string) (int, error) {
	var decimals int
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return decimals, fmt.Errorf("sdkClients failed")
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return decimals, fmt.Errorf("sdkClients failed")
	}
	client := sdkClient.queryClient
	txResponse, err := client.QueryContract(contractAddr, "Decimals", []*common.KeyValuePair{}, -1)
	if err != nil || txResponse == nil {
		return decimals, fmt.Errorf("Decimals QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	if txResponse.Code != common.TxStatusCode_SUCCESS {
		return decimals, fmt.Errorf("Decimals QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	strDecimal := string(txResponse.ContractResult.Result)
	decimals, err = strconv.Atoi(strDecimal) // 使用Atoi函数将字符串转换为int类型
	return decimals, err
}

// EvmGetDecimals 发行小数位数
func EvmGetDecimals(evmType, chainId, contractAddr string) (int, error) {
	var decimals int
	var ercAbi *abi.ABI
	//ercAbi, err := parseAbiJson(evmType)
	if evmType == ContractStandardNameEVMDFA {
		ercAbi = config.GlobalAbiERC20
	} else if evmType == ContractStandardNameEVMNFA {
		ercAbi = config.GlobalAbiERC721
	}

	if ercAbi == nil {
		return decimals, nil
	}

	dataByte, err := ercAbi.Pack("decimals")
	if err != nil {
		return decimals, fmt.Errorf("ercAbi Pack err: %s", err)
	}
	dataString := hex.EncodeToString(dataByte)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return decimals, fmt.Errorf("sdkClients failed: %s", evmType)
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return decimals, fmt.Errorf("sdkClients failed: %s", evmType)
	}
	client := sdkClient.queryClient
	txResponse, err := client.QueryContract(contractAddr, "decimals", kvs, -1)
	if err != nil || txResponse == nil {
		return decimals, fmt.Errorf("Decimals QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	if txResponse.Code != common.TxStatusCode_SUCCESS {
		return decimals, fmt.Errorf("Decimals QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	// 将字节数组转换为 *big.Int 值
	value := new(big.Int).SetBytes(txResponse.ContractResult.Result)
	decimals = int(value.Int64())
	return decimals, err
}

// EvmGetBalanceOf 代币持有量
func EvmGetBalanceOf(evmType, chainId, ownerAddr, contractAddr string) (string, error) {
	var balanceStr string
	balanceOf := big.NewInt(0)
	var ercAbi *abi.ABI
	//ercAbi, err := parseAbiJson(evmType)
	if evmType == ContractStandardNameEVMDFA {
		ercAbi = config.GlobalAbiERC20
	} else if evmType == ContractStandardNameEVMNFA {
		ercAbi = config.GlobalAbiERC721
	}

	if ercAbi == nil {
		return balanceStr, nil
	}

	ownerBytes, err := hex.DecodeString(ownerAddr)
	if err != nil {
		return balanceStr, fmt.Errorf("invalid owner address: %s", ownerAddr)
	}
	addr := evmutils.BytesToAddress(ownerBytes)

	dataByte, err := ercAbi.Pack("balanceOf", addr)
	if err != nil {
		return balanceStr, fmt.Errorf("ercAbi Pack err: %s", err)
	}
	dataString := hex.EncodeToString(dataByte)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return balanceStr, fmt.Errorf("sdkClients failed: %s", evmType)
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return balanceStr, fmt.Errorf("sdkClients failed: %s", evmType)
	}
	client := sdkClient.queryClient
	txResponse, err := client.QueryContract(contractAddr, "balanceOf", kvs, -1)
	if err != nil || txResponse == nil {
		return balanceStr, fmt.Errorf("balanceOf QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	if txResponse.Code != common.TxStatusCode_SUCCESS {
		return balanceStr, fmt.Errorf("balanceOf QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	output := txResponse.ContractResult.Result
	err = ercAbi.UnpackIntoInterface(&balanceOf, "balanceOf", output)
	return balanceOf.String(), err
}

// DockerGetContractType docker-go获取合约类型
func DockerGetContractType(chainId, contractName string) (string, error) {
	contractType := ContractStandardNameOTHER
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return contractType, fmt.Errorf("[Pool] pool has not chain, chainId:%v", chainId)
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return contractType, fmt.Errorf("[Pool] don't analyze pool, chainId:%v", chainId)
	}
	client := sdkClient.queryClient
	txResponse, err := client.QueryContract(contractName, "Standards", []*common.KeyValuePair{}, -1)
	if err != nil || txResponse == nil {
		log.Errorf("【sdk】DockerGetContractType method Standards err :%v, contractName:%v", err, contractName)
		return contractType, err
	}

	if txResponse.Code != common.TxStatusCode_SUCCESS {
		log.Errorf("【sdk】DockerGetContractType method Standards err, txResponse.Code:%v, contractName:%v",
			txResponse.Code, contractName)
		return contractType, err
	}

	var types []string
	err = json.Unmarshal(txResponse.ContractResult.Result, &types)
	log.Infof("sdk QueryContract DockerGetContractType err[%v] chainId[%v], contractName[%v], types[%v]",
		err, chainId, contractName, types)
	if err != nil {
		log.Errorf("【sdk】DockerGetContractType method Standards err, json unmarsh err:%v, contractName:%v",
			err, contractName)
		return contractType, err
	}

	if len(types) > 0 {
		if types[0] == ContractStandardNameCMBC {
			types = types[1:]
		}
		if len(types) > 0 {
			contractType = types[0]
		}
	}

	return contractType, nil
}

// DockerGetContractSymbol docker-go获取合约简称
func DockerGetContractSymbol(chainId, contractAddr string) (string, error) {
	var contractSymbol string
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return contractSymbol, fmt.Errorf("[Pool] pool has not chain, chainId:%v", chainId)
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return contractSymbol, fmt.Errorf("[Pool] don't analyze pool, chainId:%v", chainId)
	}
	client := sdkClient.queryClient
	txResponse, err := client.QueryContract(contractAddr, "Symbol",
		[]*common.KeyValuePair{}, -1)
	if err == nil && txResponse != nil {
		if txResponse.Code == common.TxStatusCode_SUCCESS {
			contractSymbol = string(txResponse.ContractResult.Result)
		}
	}

	return contractSymbol, err
}

// EVMGetContractSymbol 获取合约简称
func EVMGetContractSymbol(chainId, contractAddr, evmType string) (string, error) {
	var contractSymbol string
	var ercAbi *abi.ABI
	//ercAbi, err := parseAbiJson(evmType)
	if evmType == ContractStandardNameEVMDFA {
		ercAbi = config.GlobalAbiERC20
	} else if evmType == ContractStandardNameEVMNFA {
		ercAbi = config.GlobalAbiERC721
	}

	if ercAbi == nil {
		return contractSymbol, nil
	}

	dataByte, err := ercAbi.Pack("symbol")
	if err != nil {
		return contractSymbol, fmt.Errorf("ercAbi Pack err: %s", err)
	}
	dataString := hex.EncodeToString(dataByte)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return contractSymbol, fmt.Errorf("sdkClients failed: %s", evmType)
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return contractSymbol, fmt.Errorf("sdkClients failed: %s", evmType)
	}
	client := sdkClient.queryClient
	txResponse, err := client.QueryContract(contractAddr, "symbol", kvs, -1)
	if err != nil || txResponse == nil {
		return contractSymbol, fmt.Errorf("symbol QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	if txResponse.Code != common.TxStatusCode_SUCCESS {
		return contractSymbol, fmt.Errorf("symbol QueryContract failed err: %v,txResponse: %v", err, txResponse)
	}

	output := txResponse.ContractResult.Result
	err = ercAbi.UnpackIntoInterface(&contractSymbol, "symbol", output)
	return contractSymbol, err
}

// GetContractMultiSign 获取多签状态
func GetContractMultiSign(chainId, txId string) (*syscontract.MultiSignInfo, error) {
	multiSignInfo := &syscontract.MultiSignInfo{}
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return nil, fmt.Errorf("[Pool] pool has not chain, chainId:%v", chainId)
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return nil, fmt.Errorf("[Pool] don't analyze pool, chainId:%v", chainId)
	}
	pairs := []*common.KeyValuePair{
		{
			Key:   syscontract.MultiVote_TX_ID.String(),
			Value: []byte(txId),
		},
		{
			Key:   "truncateValueLen",
			Value: []byte("10000"),
		},
		{
			Key:   "truncateModel",
			Value: []byte("hash"),
		},
	}

	client := sdkClient.queryClient
	txResp, err := client.QueryContract(syscontract.SystemContract_MULTI_SIGN.String(),
		syscontract.MultiSignFunction_QUERY.String(), pairs, -1)
	if txResp == nil || err != nil {
		return nil, fmt.Errorf("GetTxByTxId MultiSignContractQuery failed:%v", err)
	}
	if err = proto.Unmarshal(txResp.ContractResult.Result, multiSignInfo); err != nil {
		return nil, fmt.Errorf("GetTxByTxId unmarshal failed:%v", err)
	}

	return multiSignInfo, err
}

// ExtractFunctionSignatures 获取合约4字节列表
func ExtractFunctionSignatures(bytecode []byte) [][]byte {
	var signatures [][]byte
	for i := 0; i < len(bytecode)-4; i++ {
		// 查找 PUSH4 (0x63) 指令
		if bytecode[i] == 0x63 {
			sig := bytecode[i+1 : i+5]
			signatures = append(signatures, sig)
		}
	}

	return signatures
}

// EVMGetMethodName 获取EVM合约方法
func EVMGetMethodName(ercAbi *abi.ABI, methodId []byte) (string, error) {
	// 检查给定的四个字节是否对应于 EVM 方法
	method, err := ercAbi.MethodById(methodId)
	if method != nil && err == nil {
		return method.Name, nil
	}
	return "", err
}

//// parseAbiJson 解析EVM abi
//func parseAbiJson(evmType string) (*abi.ABI, error) {
//	var jsonPath string
//	if evmType == ContractStandardNameEVMDFA {
//		jsonPath = ERC20AbiJson
//	} else if evmType == ContractStandardNameEVMNFA {
//		jsonPath = ERC721AbiJson
//	} else {
//		return nil, fmt.Errorf("parseAbiJson unsupported EVM type: %s", evmType)
//	}
//
//	// 获取当前执行的文件名
//	_, filename, _, _ := runtime.Caller(0)
//	// 获取configs目录的路径
//	dir := path.Join(path.Dir(filename), "..", "..", "configs")
//	abiFilePath := path.Join(dir, jsonPath) // 获取ABI文件的路径
//	//wd, _ := os.Getwd()
//	//erc20AbiPath := filepath.Join(wd, jsonPath)
//	ercAbiJson, err := os.Open(abiFilePath)
//	if err != nil {
//		return nil, fmt.Errorf("readFile erc20_abi failed: %v", err)
//	}
//	ercAbi, err := abi.JSON(ercAbiJson)
//	if err != nil {
//		return nil, fmt.Errorf("unmarshal erc20_abi failed: %v", err)
//	}
//
//	return &ercAbi, nil
//}

// GetContractAddr 获取合约地址
func GetContractAddr(chainId, contractName string) (string, error) {
	contractAddr := ""
	val, ok := sdkClientPool.sdkClients.Load(chainId)
	if !ok {
		return contractAddr, fmt.Errorf("[Pool] pool has not chain, chainId:%v", chainId)
	}
	sdkClient, ok := val.(*SingleSdkClientPool)
	if !ok {
		return contractAddr, fmt.Errorf("[Pool] pool has not chain, chainId:%v", chainId)
	}
	client := sdkClient.queryClient
	if contractName == "" {
		return contractAddr, fmt.Errorf("contract Name is null, chainId:%v", chainId)
	}

	res, err := client.GetContractInfo(contractName)
	if err != nil {
		log.Errorf("[Pool] get contract addr err, chainId:%v, name:%v,  err,%v", chainId, contractName, err)
		return contractAddr, fmt.Errorf("get contract addr err, chainId:%v, name:%v,  err,%v",
			chainId, contractName, err)
	}
	return res.Address, nil
}

// getMemberIdAddrAndCert 计算Addr,Id,cert
// nolint:gocyclo
func getMemberIdAddrAndCert(chainId, hashType string, member *accesscontrol.Member) (*MemberAddrIdCert, error) {
	ctx := context.Background()
	ret := MemberAddrIdCert{}
	if member == nil {
		return &ret, nil
	}
	//缓存key
	memberKey, memberKeyErr := getMemberInfoKey(chainId, hashType, int32(member.MemberType), member.MemberInfo)
	if memberKeyErr != nil {
		return &ret, memberKeyErr
	}
	//获取缓存
	redisRes := cache.GlobalRedisDb.Get(ctx, memberKey)
	if redisRes != nil && redisRes.Val() != "" {
		err := json.Unmarshal([]byte(redisRes.Val()), &ret)
		if err == nil {
			return &ret, nil
		}
		log.Errorf("getMemberIdAddrAndCert json Unmarshal err : %v", err)
	}

	var x509Cert *x509.Certificate
	var err error
	var userAddr string
	var userId string

	if member.MemberType.String() == accesscontrol.MemberType_CERT.String() {
		x509Cert, err = utils.ParseCertificate(member.MemberInfo)
		if err != nil {
			err = fmt.Errorf("[Parse SerializedMember] parse cert failed: %s", err.Error())
			return &ret, err
		}
		userId = x509Cert.Subject.CommonName
		cert, certErr := utils.X509CertToChainMakerCert(x509Cert)
		if certErr != nil {
			return &ret, certErr
		}
		userAddr, err = commonutils.CertToAddrStr(cert, pbconfig.AddrType_ETHEREUM)

		if err != nil {
			return &ret, err
		}
		ret.HasCert = true
	} else if member.MemberType.String() == accesscontrol.MemberType_CERT_HASH.String() {
		chainClient := GetChainClient(chainId)
		if chainClient == nil {
			err = fmt.Errorf("ParseCertificate get chainClient failed:  chainId%s", chainId)
			return &ret, err
		}
		certInfos, infoErr := chainClient.QueryCert([]string{hex.EncodeToString(member.MemberInfo)})
		if infoErr != nil {
			infoErr = fmt.Errorf("QuertCert failed: %v", err)
			return &ret, infoErr
		}
		if len(certInfos.CertInfos) > 0 {
			x509Cert, err = utils.ParseCertificate(certInfos.CertInfos[0].Cert)
			if err != nil {
				err = fmt.Errorf("ParseCertificate failed: %s", err.Error())
				return &ret, err
			}
			userId = x509Cert.Subject.CommonName
			cert, certErr := utils.X509CertToChainMakerCert(x509Cert)
			if certErr != nil {
				return &ret, certErr
			}
			userAddr, err = commonutils.CertToAddrStr(cert, pbconfig.AddrType_ETHEREUM)

			if err != nil {
				return &ret, err
			}
		}
		ret.HasCert = true
	} else if member.MemberType.String() == accesscontrol.MemberType_ALIAS.String() {
		chainClient := GetChainClient(chainId)
		if chainClient == nil {
			err = fmt.Errorf("ParseCertificate get chainClient failed:  chainId%s", chainId)
			return &ret, err
		}

		aliasInfos, errAli := chainClient.QueryCertsAlias([]string{string(member.MemberInfo)})
		if errAli != nil {
			return &ret,
				fmt.Errorf("Quert AliasInfos failed: %s ", err.Error())
		}
		//nolint:staticcheck
		for _, aliasInfo := range aliasInfos.AliasInfos {
			cert := aliasInfo.NowCert.Cert
			if cert == nil {
				for k := len(aliasInfo.HisCerts) - 1; k >= 0; k-- {
					if aliasInfo.HisCerts[k].Cert != nil {
						cert = aliasInfo.HisCerts[k].Cert
						break
					}
				}
			}
			x509Cert, err = utils.ParseCertificate(cert)
			if err != nil {
				err = fmt.Errorf("ParseCertificate failed: %s", err.Error())
				return &ret, err
			}
			userId = x509Cert.Subject.CommonName
			certP, certErr := utils.X509CertToChainMakerCert(x509Cert)
			if certErr != nil {
				return &ret, certErr
			}
			userAddr, err = commonutils.CertToAddrStr(certP, pbconfig.AddrType_ETHEREUM)

			if err != nil {
				return &ret, err
			}
			break
		}
		ret.HasCert = true
	} else if member.MemberType.String() == accesscontrol.MemberType_PUBLIC_KEY.String() {
		publicKeyStr := member.MemberInfo
		chainHashType := hashType
		publicKey, pkErr := asym.PublicKeyFromPEM(publicKeyStr)
		if pkErr != nil {
			return &ret, pkErr
		}
		userAddr, err = commonutils.PkToAddrStr(publicKey, pbconfig.AddrType_ETHEREUM,
			crypto.HashAlgoMap[chainHashType])
		if err != nil {
			return &ret, err
		}
		//public模式没有userId，使用userAddr作为userId入库
		userId = userAddr

		//userId, err = helper.CreateLibp2pPeerIdWithPublicKey(publicKey)
		//if err != nil {
		//	return &ret, err
		//}
	}
	ret.UserId = userId
	ret.UserAddr = userAddr
	if ret.HasCert {
		ret.CertCommonName = x509Cert.Subject.CommonName
		ret.CertOrganization = x509Cert.Subject.Organization
		ret.CertOrganizationUnit = x509Cert.Subject.OrganizationalUnit
	}

	//缓存数据
	retJson, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("Error getMemberIdAddrAndCert json marshal err: %v，memberKey：%v", err, memberKey)
	} else {
		// 设置键值对和过期时间
		_ = cache.GlobalRedisDb.Set(ctx, memberKey, string(retJson), time.Hour).Err()
	}
	return &ret, nil
}
