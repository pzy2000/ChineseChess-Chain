### API接口文档

#### 1. 整体请求格式
整体请求全部采用 **GET** 请求方式

URL格式：http://localhost/chainmaker?cmb={OpHandle}

contentType采用 **application/form-data**

请求的数据通过body体中以 **form-data** 方式发送

URL参数说明如下：

|字段| 描述 |
| :----- | :----- | 
|OpHandle|固定字符串，用于描述本次操作的类型| 

#### 2. 整体返回数据格式
默认的格式会包括Response一个字段，使用下面的格式：
```json
{
	"Response": {
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c",
        "Data": {}
    }
}
```

若返回的数据是一个数组的话，和上面的一样，可参考下面的格式：
```json
{
	"Response": {
        "TotalCount": 2,
        "GroupList": [
            {},
            {}
        ],
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }

}
```

若返回的数据是一个标识异常的结果，则可参考下面的格式：
```json
{
	"Response": {
        "Error": {
            "Code": "AuthFailure",
            "Message": "pamam is nil"
        },
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }

}

```

说明如下：

|字段| 描述 |取值范围 | 备注  |
| :----- | :----- | :----- | :----- |
|Response|应答关键字| ----- | 固定格式 |
|Data|返回对象| json格式  | 返回结果是对象（非集合）时使用 |
|TotalCount|总量，用于计算分页| 数字  | 返回结果是集合时使用 |
|GroupList|结果列表| 集合  | 返回结果是集合时使用 |


#### 2. 获取合约详情
```text
获取合约详情数据
```

#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetContractDetail&ContractKey=ef3a386d44c4d8b0cb59f06def82237456b41f76

#### 请求方式
> GET

#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetContractDetail | String | 是 | 方法名称
ContractKey | ef3a386d44c4d8b0cb59f06def82237456b41f76 | String | 是 | 合约名称/合约地址


#### 成功响应示例
```javascript
{
    "Response": {
        "Data": {
            "ContractName": "DID1234",
            "ContractNameBak": "DID1234",
            "ContractAddr": "9af6a4d25738a41224303a1dc0789fa49eb11ba0",
            "ContractSymbol": "",
            "ContractType": "CMDID",
            "ContractStatus": 0,
            "Version": "1.0",
            "TxId": "17beb3de53797af1ca2918527ca19ed3b9e99ed421434c27bf6bceb2cd64075f",
            "CreateSender": "e97e032d9fbc9f282cdb779fb365afea70bcb1d4",
            "CreatorAddr": "e97e032d9fbc9f282cdb779fb365afea70bcb1d4",
            "CreatorAddrBns": "",
            "Timestamp": 1711002675,
            "RuntimeType": "DOCKER_GO"
        },
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | - | object | -
Response.Data.ContractName | - | string | 合约名称
Response.Data.ContractNameBak | 情料相 | string | 合约名称
Response.Data.ContractSymbol | wewe | string | 合约简称
Response.Data.ContractType | 成位领资领声非成么看根放里土为高它明程别千马总将现山历 | string | 合约类型
Response.Data.ContractAddr | 低包类半阶长么分层和关料然 | string | 合约地址
Response.Data.Version | 部且反身到记果 | string | 合约版本
Response.Data.ContractStatus | -137716752190952 | integer | 合约状态（0:正常，1:冻结，2:注销）
Response.Data.TxId | 称与第但给包少给理六 | string | 交易ID
Response.Data.CreateSender | 安委由何单即如个器求育样 | string | 创建用户ID
Response.Data.CreatorAddr | 打土技细 | string | 创建用户地址
Response.Data.Timestamp | -2351064116902996 | integer | 上链时间
Response.Data.RuntimeType | EVM | string | 虚拟机类型
Response.Data.TotalSupply | 二那所该见来办争型务习资然直意流党里满离 | string | -
Response.RequestId | 解再次般识华较三听工经质条口 | string | -

### 3、获取链上用户列表
```text
获取链上用户列表
```

#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetUserList&Offset=0&Limit=10&UserAddrs=&UserIds=&OrgId=

#### 请求方式
> GET


#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetUserList | String | 是 | 方法名称
Offset | 0 | Integer | 是 | 页码
Limit | 10 | String | 是 | 条数
UserAddrs | - | String | 否 | user地址列表
UserIds | - | String | 否 | 用户id列表
OrgId | - | String | 否 | 组织id



#### 成功响应示例
```javascript
{
    "Response": {
        "GroupList": [
            {
                "Id": "1",
                "UserId": "09553323aea75de47a3298d717c237fb0d518888",
                "OrgId": "public",
                "Role": "admin",
                "Timestamp": 1711075554,
                "UserAddr": "09553323aea75de47a3298d717c237fb0d518888",
                "Status": 0
            },
            {
                "Id": "2",
                "UserId": "1dca44823a782008ef1d14dfd0e5a5a7b4228e12",
                "OrgId": "public",
                "Role": "admin",
                "Timestamp": 1711075554,
                "UserAddr": "1dca44823a782008ef1d14dfd0e5a5a7b4228e12",
                "Status": 0
            }
        ],
        "TotalCount": 2,
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.UserId | - | string | 用户ID
Response.GroupList.0.OrgId | 任阶局中团进单信光心由北管收已积之标权教统 | string | 组织ID
Response.GroupList.0.Role | 王立 | string | 用户身份
Response.GroupList.0.UserAddr | 张七白需发料 | string | 用户地址
Response.GroupList.0.Status | 1 | integer | 用户状态（0:正常，1:封禁）
Response.GroupList.0.Timestamp | -4954233422542663 | integer | -
Response.TotalCount | -6752290995179235 | integer | -
Response.RequestId | 直区 | string | -
### 4、首页详情数据
```text
获取首页详情数据
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chainmaker_pk&cmb=GetOverviewData

#### 请求方式
> GET

#### Content-Type
> form-data

#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chainmaker_pk | String | 是 | 链ID
cmb | GetOverviewData | String | 是 | 方法名称



#### 成功响应示例
```javascript
{
    "Response": {
        "Data": {
            "ChainId": "chainmaker_pk",
            "BlockHeight": 375,
            "UserCount": 11,
            "ContractCount": 37,
            "TxCount": 765,
            "OrgCount": 0,
            "RunningNode": 4,
            "CommonNode": 0,
            "ConsensusNode": 4
        },
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.Data | - | object | -
Response.Data.ContractName | goErc20 | string | 合约名称
Response.Data.ContractAddr | aba31ce4cd49f08073d2f115eb12610544242ff9 | string | 合约地址
Response.Data.TokenId | 1234 | string | TokenId
Response.Data.OwnerAddr | 第动那达指压部成好下来义往声都铁心个 | string | 持有地址
Response.Data.AddrType | 出切及强什己好来能满 | string | 持有地址类型
Response.Data.Timestamp | -1497110644697328 | integer | 上链时间
Response.Data.Metadata | - | object | -
Response.Data.Metadata.Name | 与 | string | 作品名称
Response.Data.Metadata.Author | 光品领我为义收二细根动得子学治教 | string | 作者
Response.Data.Metadata.ImageUrl | 节几革决证完后多打原把值求当华 | string | 图片地址
Response.Data.Metadata.Description | 以治适基为 | string | 作品描述
Response.Data.Metadata.SeriesHash | 火全热她 | string | 作品Hash
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | -

### 5、获取合约版本交易列表
```text
获取合约创建，升级相关交易列表
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetContractVersionList&Offset=0&Limit=10&ContractName=&ContractAddr=d4c6b94f36dd6a514bca8271518632370f522084&Senders=client1.sign.wx-org1.chainmaker.org,client1.sign.wx-org2.chainmaker.org&RuntimeType=&Status=

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetContractVersionList | String | 是 | 方法名称
Offset | 0 | Integer | 是 | 页码
Limit | 10 | String | 是 | 条数
ContractName | - | String | 否 | 合约名称/
ContractAddr | d4c6b94f36dd6a514bca8271518632370f522084 | String | 否 | 合约地址
Senders | client1.sign.wx-org1.chainmaker.org,client1.sign.wx-org2.chainmaker.org | String | 否 | userID集合，
RuntimeType | - | String | 否 | 虚拟机类型
Status | - | String | 否 | 合约运行状态（0:成功，1:失败）默认不传为查询全部数据



#### 成功响应示例
```javascript
{
    "Response": {
        "GroupList": [
            {
                "TxId": "17beb3de53797af1ca2918527ca19ed3b9e99ed421434c27bf6bceb2cd64075f",
                "Sender": "e97e032d9fbc9f282cdb779fb365afea70bcb1d4",
                "SenderAddr": "e97e032d9fbc9f282cdb779fb365afea70bcb1d4",
                "SenderAddrBNS": "",
                "SenderOrgId": "public",
                "Version": "1.0",
                "ContractName": "DID1234",
                "ContractAddr": "9af6a4d25738a41224303a1dc0789fa49eb11ba0",
                "TxUrl": "http://subchain-service:8888/chainmaker_pk/transaction/17beb3de53797af1ca2918527ca19ed3b9e99ed421434c27bf6bceb2cd64075f/",
                "ContractResultCode": 0,
                "Timestamp": 1711002675
            }
        ],
        "TotalCount": 1,
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.TxId | 角按和那活采应风二张公示 | string | 交易ID
Response.GroupList.0.Sender | 文前 | string | 创建账户ID
Response.GroupList.0.SenderAddr | 月保 | string | 创建账户地址
Response.GroupList.0.SenderAddrBNS | - | string | 创建账户地址BNS
Response.GroupList.0.SenderOrgId | 风应安高量量办加经回安法九证记重于展 | string | 创建账户组织
Response.GroupList.0.Version | 号压 | string | 合约版本
Response.GroupList.0.ContractName | 车青直商花行难回装 | string | 合约名称
Response.GroupList.0.ContractAddr | 采程科积布么种信定验界决精原选道群 | string | 合约地址
Response.GroupList.0.ContractResultCode | 1982922855930811 | integer | 合约结果码（0:成功，1:失败）
Response.GroupList.0.Timestamp | 5635886076384591 | integer | 上链时间
Response.GroupList.0.RuntimeType | 与 | string | 虚拟机类型
Response.GroupList.0.TxUrl | - | string | -
Response.TotalCount | 7857962009929503 | integer | -
Response.RequestId | 任次革动来体其林称 | string | -
### 6、首页搜索
```text
首页搜索接口，返回数据是否存在，并返回type对应的查询结果
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=Search&Type=7&Value=171262347a59fded92021a32421a5dad05424e03

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | Search | String | 是 | 方法名称
Type | 7 | Integer | 是 | 1:区块hash，2:区块高度，3：交易ID，4:合约名称，5:合约地址，6:账户地址，7:BNS
Value | 171262347a59fded92021a32421a5dad05424e03 | String | 是 | 搜索值



#### 成功响应示例
```javascript
{
    "Response": {
        "Data": {
            "Type": 0,
            "Data": "173ff09c5decd90f36aa435d6feab7203aa974aed78b62911fcde310df530ea6",
            "ChainId": "chainmaker_pk",
            "ContractType": ""
        },
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | Object | -
Response.Data | - | Object | -
Response.Data.Type | 2 | Integer | 返回值类型（-1:未找到。0:区块hash，2：交易ID，3：合约地址，4:账户地址）
Response.Data.Data | 6cd144e330edbe27f82bb44dbd06283836f000a3 | String | type对应的值
Response.Data.ChainId | chain1 | String | 链ID
Response.RequestId | 务理向展由作 | String | -
## 7、获取交易列表
```text
分页获取交易列表
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chainmaker_pk&cmb=GetTxList&Offset=0&Limit=10&TxId=&BlockHash=&ContractMethod=&ContractName=&ContractAddr=&UserAddrs=&Senders=&StartTime=&EndTime=

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chainmaker_pk | String | 是 | 链ID
cmb | GetTxList | String | 是 | 方法名称
Offset | 0 | String | 是 | 页码
Limit | 10 | String | 是 | 分页条
TxId | - | String | 否 | 交易ID
BlockHash | - | String | 否 | 区块hash
ContractMethod | - | String | 否 | 合约方法
ContractName | - | String | 否 | 合约名称/
ContractAddr | - | String | 否 | 合约地址
UserAddrs | - | String | 否 | user地址集合，多个以逗号分割123,345
Senders | - | String | 否 | userID集合，多个以逗号分割123,345
StartTime | - | Integer | 否 | 开始时间
EndTime | - | String | 否 | 结束时间
TxStatus | - | Integer | 否 | 交易状态。（0：成功，1:失败）不传默认全部数据



#### 成功响应示例
```javascript
{
    "Response": {
        "GroupList": [
            {
                "Id": "1",
                "TxId": "17b0ed03699b09e0ca176305434c3517bf68ad78221d4ce3b171218a5c88e443",
                "Sender": "e97e032d9fbc9f282cdb779fb365afea70bcb1d4",
                "SenderOrgId": "public",
                "BlockHeight": 12,
                "ContractName": "official_identity",
                "ContractAddr": "2ac650bf8cf7fb926ae2c908888690fef07176e6",
                "ContractMethod": "SetIdentity",
                "ContractParameters": "[{\"key\":\"address\",\"value\":\"ZTk3ZTAzMmQ5ZmJjOWYyODJjZGI3NzlmYjM2NWFmZWE3MGJjYjFkNA==\"},{\"key\":\"level\",\"value\":\"MQ==\"},{\"key\":\"pkPem\",\"value\":\"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb0VjejFVQmdpMERRZ0FFbnFTRnRPTDNWZTErMG9TNC9KUlVMVXBqQlQ5TgpMNjBzUmNqVDB1MXhUUTFNQVAza2h6SUVqVUg4NlhZZGxPYTFOR1V4S1dybnVaVnVkQ2t5RWM4RWJRPT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg==\"},{\"key\":\"metadata\",\"value\":\"eyJvcmdJZCI6IuW+ruiKryIsInRpbWVzdGFtcCI6IjE2ODI2NjE2NDYifQ==\"}]",
                "Status": "",
                "TxStatus": 0,
                "ShowStatus": 0,
                "BlockHash": "173ff09c5decd90f36aa435d6feab7203aa974aed78b62911fcde310df530ea6",
                "Timestamp": 1707124857,
                "UserAddr": "e97e032d9fbc9f282cdb779fb365afea70bcb1d4",
                "UserAddrBns": "",
                "GasUsed": 247,
                "PayerAddr": "",
                "Event": "[{\"index\":0,\"key\":\"setIdentity\",\"value\":\"e97e032d9fbc9f282cdb779fb365afea70bcb1d4,1,-----BEGIN PUBLIC KEY-----\\nMFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAEnqSFtOL3Ve1+0oS4/JRULUpjBT9N\\nL60sRcjT0u1xTQ1MAP3khzIEjUH86XYdlOa1NGUxKWrnuZVudCkyEc8EbQ==\\n-----END PUBLIC KEY-----\\n\"}]"
            },
            {
                "Id": "2",
                "TxId": "6dd14591604a4accbe99886f7dd236cfecc0c1c37af54251b2018b55414f913d",
                "Sender": "7e2d16337700a113c465ede743cdb34faa86e7ce",
                "SenderOrgId": "public",
                "BlockHeight": 12,
                "ContractName": "ACCOUNT_MANAGER",
                "ContractAddr": "a564b97f89f1c64b4ae12465a5870e879ffa0d3f",
                "ContractMethod": "CHARGE_GAS_FOR_MULTI_ACCOUNT",
                "ContractParameters": "[{\"key\":\"e97e032d9fbc9f282cdb779fb365afea70bcb1d4\",\"value\":\"MjQ3\"}]",
                "Status": "",
                "TxStatus": 0,
                "ShowStatus": 0,
                "BlockHash": "173ff09c5decd90f36aa435d6feab7203aa974aed78b62911fcde310df530ea6",
                "Timestamp": 1707124857,
                "UserAddr": "7e2d16337700a113c465ede743cdb34faa86e7ce",
                "UserAddrBns": "",
                "GasUsed": 0,
                "PayerAddr": "",
                "Event": ""
            }
        ],
        "TotalCount": 2,
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.TxId | e4b090ef87c24f3e9ab52b68b57c8e7dcc7ab73459974bcebafbc8a4b745f150 | string | 交易ID
Response.GroupList.0.ContractName | goErc20 | string | 合约名称
Response.GroupList.0.ContractAddr | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | 合约地址
Response.GroupList.0.ContractMethod | 品科保发接出产里展立外书见亲济难分及成改 | string | 合约方法
Response.GroupList.0.Sender | QmTqj47NaAq9u3QT8JBpMVaui3csgMomBw6aePh6ufa2cU | string | 发起用户ID
Response.GroupList.0.SenderOrgId | 组参 | string | 发起用户组织
Response.GroupList.0.UserAddr | 即信小生格识需 | string | 发起用户地址
Response.GroupList.0.UserAddrBns | - | string | 发起用户BNS
Response.GroupList.0.BlockHeight | 12 | integer | 区块高度
Response.GroupList.0.BlockHash | cd2572bef619f883c9bd12401943b0ae243b59bc2860c6af17a681a17b4c50b1 | string | 区块hash
Response.GroupList.0.TxStatus | -8688340935408671 | integer | 交易状态（0:成功，1失败）
Response.GroupList.0.GasUsed | 329470078576235 | integer | gas消耗
Response.GroupList.0.ContractParameters | - | string | -
Response.GroupList.0.PayerAddr | 复于任感阶提十界识度平圆育向定广料以斗海意思每地三作统 | string | 支付地址
Response.GroupList.0.Event | 出设山保矿 | string | 交易事件
Response.GroupList.0.Timestamp | 1698917816 | integer | 上链时间
Response.TotalCount | 10 | integer | 总数
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID

## 8、获取合约列表
```text
分页获取合约列表。（注意：ContractName改成了ContractKey）
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetContractList&Limit=10&Offset=0&ContractKey=d4c6b94f36dd6a514bca8271518632370f522084&Status=&RuntimeType=&Creators=123,456&CreatorAddrs=123，456&UpgradeAddrs=123,456&StartTime=&EndTime=

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetContractList | String | 是 | 方法名称
Limit | 10 | String | 是 | 分页条
Offset | 0 | String | 是 | 页码
ContractKey | d4c6b94f36dd6a514bca8271518632370f522084 | String | 否 | 合约名称/合约地址
Status | - | Integer | 否 | 合约状态(0:正常，1:冻结，2:注销)不传默认全部数据
RuntimeType | - | String | 否 | 虚拟机类型
Creators | 123,456 | String | 否 | 创建用户名称/BNS地址
CreatorAddrs | 123，456 | String | 是 | 创建用户地址
UpgradeAddrs | 123,456 | String | 否 | 更新用户地址集合
StartTime | - | Integer | 否 | 开始时间
EndTime | - | Integer | 否 | 结束时间



#### 成功响应示例
```javascript
{
    "Response": {
        "GroupList": [
            {
                "Id": "1",
                "ContractName": "DID1234",
                "ContractSymbol": "",
                "ContractAddr": "9af6a4d25738a41224303a1dc0789fa49eb11ba0",
                "Version": "1.0",
                "TxNum": 63,
                "Status": 0,
                "Creator": "e97e032d9fbc9f282cdb779fb365afea70bcb1d4",
                "CreatorUser": "",
                "CreatorAddr": "e97e032d9fbc9f282cdb779fb365afea70bcb1d4",
                "CreatorAddrBns": "",
                "UpgradeUser": "",
                "CreateTimestamp": 1711002675,
                "UpdateTimestamp": 0,
                "ContractType": "CMDID",
                "RuntimeType": "DOCKER_GO"
            }
        ],
        "TotalCount": 1,
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.ContractName | goErc20 | string | 合约名称
Response.GroupList.0.ContractSymbol | ERC20 | string | 合约简称
Response.GroupList.0.ContractAddr | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | 合约地址
Response.GroupList.0.ContractType | ERC20 | string | 合约类型（token类：CMDFA,  ERC20; NFT类：CMNFA，ERC721）
Response.GroupList.0.Version | 1.0 | string | 合约版本
Response.GroupList.0.RuntimeType | EVM | string | evm/docker-go
Response.GroupList.0.TxNum | 1000 | integer | 累计交易数
Response.GroupList.0.Status | - | integer | 合约状态（0:正常，1:冻结，2:注销）
Response.GroupList.0.Creator | 23232323 | string | 创建用户ID
Response.GroupList.0.CreatorAddr | ab108fc6c3850e01cee01e419d07f097186c3982 | string | 创建用户地址
Response.GroupList.0.CreatorAddrBns | 除气为得 | string | 创建用户BNS
Response.GroupList.0.UpgradeAddr | ab108fc6c3850e01cee01e419d07f097186c3982 | string | 更新用户地址
Response.GroupList.0.CreateTimestamp | 1699242950 | integer | 创建时间
Response.GroupList.0.UpdateTimestamp | 1699242950 | integer | 更新时间
Response.TotalCount | 10 | integer | 总数
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID



## 9、按时间段查询交易量
```text
按时间段查询交易量
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chainmaker_pk&cmb=GetTxNumByTime&SortType=1

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
Interval | 3600 | Integer | 否 | 间隔时长，默认3600s，单位秒
StartTime | 1704038800 | String | 否 | 开始时间，默认24小时
EndTime | 1704188000 | String | 否 | 结束时间
ChainId | chainmaker_pk | String | 是 | 链ID
cmb | GetTxNumByTime | String | 是 | 方法名称
SortType | 1 | Integer | 是 | 排序，0：倒序，1：正序



#### 成功响应示例
```javascript
{
    "Response": {
        "GroupList": [
            {
                "TxNum": 0,
                "Timestamp": 1711076400
            },
            {
                "TxNum": 28,
                "Timestamp": 1711072800
            },
            {
                "TxNum": 10,
                "Timestamp": 1711069200
            },
            {
                "TxNum": 0,
                "Timestamp": 1711065600
            },
            {
                "TxNum": 0,
                "Timestamp": 1711062000
            },
            {
                "TxNum": 0,
                "Timestamp": 1711058400
            },
            {
                "TxNum": 0,
                "Timestamp": 1711054800
            },
            {
                "TxNum": 0,
                "Timestamp": 1711051200
            },
            {
                "TxNum": 0,
                "Timestamp": 1711047600
            },
            {
                "TxNum": 0,
                "Timestamp": 1711044000
            },
            {
                "TxNum": 0,
                "Timestamp": 1711040400
            },
            {
                "TxNum": 0,
                "Timestamp": 1711036800
            },
            {
                "TxNum": 0,
                "Timestamp": 1711033200
            },
            {
                "TxNum": 0,
                "Timestamp": 1711029600
            },
            {
                "TxNum": 0,
                "Timestamp": 1711026000
            },
            {
                "TxNum": 0,
                "Timestamp": 1711022400
            },
            {
                "TxNum": 0,
                "Timestamp": 1711018800
            },
            {
                "TxNum": 32,
                "Timestamp": 1711015200
            },
            {
                "TxNum": 14,
                "Timestamp": 1711011600
            },
            {
                "TxNum": 48,
                "Timestamp": 1711008000
            },
            {
                "TxNum": 24,
                "Timestamp": 1711004400
            },
            {
                "TxNum": 12,
                "Timestamp": 1711000800
            },
            {
                "TxNum": 9,
                "Timestamp": 1710997200
            },
            {
                "TxNum": 0,
                "Timestamp": 1710993600
            }
        ],
        "TotalCount": 24,
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.TxNum | -1684401925513780 | integer | -
Response.GroupList.0.Timestamp | 5151256970057851 | integer | -
Response.TotalCount | -4687195283644155 | integer | -
Response.RequestId | 管全强影性以日地快组发后子常专层 | string | -
## 10、获取链列表
```text
获取订阅链列表
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chain1&Offset=0&Limit=10&cmb=GetChainList

#### 请求方式
> GET


#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 否 | 链ID
Offset | 0 | Integer | 是 | 页码
Limit | 10 | Integer | 是 | 条数
cmb | GetChainList | String | 是 | 方法名称
#### 请求Body参数
```javascript
{
	"ChainId": "chain1",
	"cmb": "GetChainList",
	"Offset": "0",
	"Limit": "10"
}
```



#### 成功响应示例
```javascript
{
    "Response": {
        "GroupList": [
            {
                "ChainId": "chainmaker_pk",
                "ChainVersion": "2030300",
                "Status": 0,
                "Consensus": "TBFT",
                "Timestamp": 1711075554,
                "AuthType": "public"
            }
        ],
        "TotalCount": 1,
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.Timestamp | -4796367723297775 | integer | 上链时间
Response.GroupList.0.Topic | 质 | string | topic
Response.GroupList.0.EventInfo | 多温 | string | 事件信息
Response.TotalCount | 8717360574102087 | integer | -
Response.RequestId | 克道支克明那群给式已感府法本教合问 | string | -
## 11、账户详情
```text
获取账户详情
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetAccountDetail&Address=123456&BNS=1234

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetAccountDetail | String | 是 | 方法名称
Address | 123456 | String | 否 | 账户地址和BNS必须传一个
BNS | 1234 | String | 否 | 账户地址和BNS必须传一个



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": {
			"Address": "aba31ce4cd49f08073d2f115eb12610544242ff9",
			"type": 1,
			"BNS": "BNS:12232",
			"DID": "DID:123123123"
		},
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.Data | - | object | -
Response.Data.Address | aba31ce4cd49f08073d2f115eb12610544242ff9 | string | 账户地址
Response.Data.type | 1 | integer | 账户类型，0:用户地址，1:合约地址
Response.Data.BNS | BNS:12232 | string | BNS
Response.Data.DID | DID:123123123 | string | DID
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | -
## 12、获取合约事件列表
```text
获取合约事件列表
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker?ChainId=chain1&cmb=GetLatestContractList

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetLatestContractList | String | 是 | 方法名称



#### 成功响应示例
```javascript
{
    "Response": {
        "GroupList": [
            {
                "EventInfo": "[\"did:cnbn:cnbn\",\"did:cnbn:20322ce7831247ffbb2b118fcf869ebf\",\"100000\",\"http://pre-api.cnbn.org.cn/api/v1/did/vc/50b85eff90124c63b1e1907e4fa49405\"]",
                "Topic": "VcIssueLog",
                "Timestamp": 1711078327
            },
            {
                "EventInfo": "[\"did:cnbn:cnbn\",\"did:cnbn:20322ce7831247ffbb2b118fcf869ebf\",\"100000\",\"32131321312313321312\"]",
                "Topic": "VcIssueLog",
                "Timestamp": 1711078327
            }
        ],
        "TotalCount": 64,
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.Timestamp | -4796367723297775 | integer | 上链时间
Response.GroupList.0.Topic | 质 | string | topic
Response.GroupList.0.EventInfo | 多温 | string | 事件信息
Response.TotalCount | 8717360574102087 | integer | -
Response.RequestId | 克道支克明那群给式已感府法本教合问 | string | -
## 13、获取NFT合约列表
```text
分页获取非同质化合约列表
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetNFTContractList&Limit=10&Offset=0&ContractKey=

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetNFTContractList | String | 是 | 方法名称
Limit | 10 | String | 是 | 分页条
Offset | 0 | String | 是 | 页码
ContractKey | - | String | 否 | 合约名称/合约地址



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"ContractName": "goErc20",
				"ContractAddr": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"ContractType": "CMDFA",
				"TxNum": 123,
				"TotalSupply": "12344",
				"HolderCount": 12,
				"Timestamp": 1699242950
			}
		],
		"TotalCount": 10,
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.ContractName | goErc20 | string | 合约名称
Response.GroupList.0.ContractAddr | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | 合约地址
Response.GroupList.0.ContractType | CMDFA | string | 合约类型
Response.GroupList.0.TxNum | 123 | integer | 累计交易数
Response.GroupList.0.TotalSupply | 12344 | string | 发行总量
Response.GroupList.0.HolderCount | 12 | integer | 持有人数
Response.GroupList.0.Timestamp | 1699242950 | integer | 创建时间
Response.TotalCount | 10 | integer | 总数
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID
## 14、获取NFT合约详情
```text
获取非同质化合约详情
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetNFTContractDetail&ContractAddr=d4c6b94f36dd6a514bca8271518632370f522084

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetNFTContractDetail | String | 是 | 方法名称
ContractAddr | d4c6b94f36dd6a514bca8271518632370f522084 | String | 是 | 合约地址



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": {
			"ContractName": "goErc20",
			"ContractNameBak": "满计约科原地技形识和争我明",
			"ContractAddr": "aba31ce4cd49f08073d2f115eb12610544242ff9",
			"ContractStatus": 0,
			"ContractType": "后拉且装",
			"RuntimeType": "成样解光近中采色角难越亲局",
			"Version": "事",
			"TxId": "17a05af0cfc1826fca85dc0d21083824de9ec658943c42738297fd372fc56d7f",
			"CreateUser": "171262347a59fded92021a32421a5dad05424e03",
			"TotalSupply": "123445",
			"HolderCount": 12,
			"CreateTimestamp": 1702460649,
			"UpdateTimestamp": 1702460649
		},
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.Data | - | object | -
Response.Data.ContractName | goErc20 | string | 合约名称
Response.Data.ContractNameBak | 满计约科原地技形识和争我明 | string | 合约名称
Response.Data.ContractAddr | aba31ce4cd49f08073d2f115eb12610544242ff9 | string | 合约地址
Response.Data.ContractStatus | - | integer | 合约状态
Response.Data.ContractType | 后拉且装 | string | 合约类型
Response.Data.RuntimeType | 成样解光近中采色角难越亲局 | string | 虚拟机类型
Response.Data.Version | 事 | string | 合约版本
Response.Data.TxId | 17a05af0cfc1826fca85dc0d21083824de9ec658943c42738297fd372fc56d7f | string | 创建交易
Response.Data.CreateUser | 171262347a59fded92021a32421a5dad05424e03 | string | 创建用户地址
Response.Data.TotalSupply | 123445 | string | 发行量
Response.Data.HolderCount | 12 | integer | 持有人数
Response.Data.CreateTimestamp | 1702460649 | integer | 创建时间
Response.Data.UpdateTimestamp | 1702460649 | integer | 更新时间
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | -
## 15、获取NFT合约持仓列表
```text
分页获取非同质化合约的持仓列表
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetNFTPositionList&Offset=0&Limit=10&ContractAddr=d4c6b94f36dd6a514bca8271518632370f522084&OwnerAddr=171262347a59fdewwwd92021a32421a5dad05424e03

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetNFTPositionList | String | 是 | 方法名称
Offset | 0 | String | 是 | 页码
Limit | 10 | String | 是 | 分页条
ContractAddr | d4c6b94f36dd6a514bca8271518632370f522084 | String | 否 | 合约地址
OwnerAddr | 171262347a59fdewwwd92021a32421a5dad05424e03 | String | 否 | 持仓地址



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"AddrType": 1,
				"ContractAddr": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"ContractName": "adasd",
				"OwnerAddr": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"Amount": "122.333",
				"HoldRatio": "12.3445%",
				"HoldRank ": 12
			}
		],
		"TotalCount": 10,
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.AddrType | 1 | integer | 持仓地址类型：0:用户，1:合约地址
Response.GroupList.0.ContractAddr | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | 合约地址
Response.GroupList.0.ContractName | adasd | string | 合约名称
Response.GroupList.0.OwnerAddr | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | 持仓地址
Response.GroupList.0.Amount | 122.333 | string | 持有量
Response.GroupList.0.HoldRatio | 12.3445% | string | 持有比例
Response.GroupList.0.HoldRank | 12 | integer | 持有排名
Response.TotalCount | 10 | integer | 总数
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID
## 16、获取NFT流转列表
```text
分页获取同质化交易流转列表
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetNFTTransferList&Offset=0&Limit=10&TxId=&TokenId=10010&ContractAddr=&UserAddr=

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetNFTTransferList | String | 是 | 方法名称
Offset | 0 | String | 是 | 页码
Limit | 10 | String | 是 | 分页条
TxId | - | String | 否 | 交易ID
TokenId | 10010 | String | 否 | tokenID
ContractAddr | - | String | 否 | 合约地址
UserAddr | - | String | 否 | 流转地址



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"TxId": "e4b090ef87c24f3e9ab52b68b57c8e7dcc7ab73459974bcebafbc8a4b745f150",
				"ContractName": "goErc20",
				"ContractAddr": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"ContractMethod": "transfer",
				"From": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"To": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"TokenId": "1234",
				"ImageUrl": "http://1234.jpg",
				"Timestamp ": 1698917816
			}
		],
		"TotalCount": 10,
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.TxId | e4b090ef87c24f3e9ab52b68b57c8e7dcc7ab73459974bcebafbc8a4b745f150 | string | 交易ID
Response.GroupList.0.ContractName | goErc20 | string | 合约名称
Response.GroupList.0.ContractAddr | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | 合约地址
Response.GroupList.0.ContractMethod | transfer | string | 合约方法
Response.GroupList.0.From | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | from地址
Response.GroupList.0.To | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | to地址
Response.GroupList.0.TokenId | 1234 | string | tokenID
Response.GroupList.0.ImageUrl | http://1234.jpg | string | 图片地址
Response.GroupList.0.Timestamp | 1698917816 | integer | 上链时间
Response.TotalCount | 10 | integer | 总数
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID
## 17、/获取NFT列表
```text
分页获取NFT列表
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetNFTList&Offset=0&Limit=10&TokenId=&ContractKey=&OwnerAddrs=&StartTime=&EndTime=

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetNFTList | String | 是 | 方法名称
Offset | 0 | String | 是 | 页码
Limit | 10 | String | 是 | 分页条
TokenId | - | String | 否 | tokenID
ContractKey | - | String | 否 | 合约地址/合约名称
OwnerAddrs | - | String | 否 | user地址集合，多个以逗号分割123,345
StartTime | - | String | 否 | 开始时间
EndTime | - | String | 否 | 结束时间



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"AddrType": 1,
				"ContractAddr": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"ContractName": "adasd",
				"OwnerAddr": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"OwnerAddrBNS": "四在我提百活可知增所参产",
				"TokenId": "122",
				"Timestamp": 12121213232,
				"CategoryName": "wewer",
				"ImageUrl": "http://1234.jpg"
			}
		],
		"TotalCount": 10,
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.AddrType | 1 | integer | 持仓地址类型：0:用户地址，1:合约地址
Response.GroupList.0.ContractAddr | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | 合约地址
Response.GroupList.0.ContractName | adasd | string | 合约名称
Response.GroupList.0.OwnerAddr | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | 持仓地址
Response.GroupList.0.OwnerAddrBNS | 四在我提百活可知增所参产 | string | 持仓地址BNS
Response.GroupList.0.TokenId | 122 | string | tokenID
Response.GroupList.0.Timestamp | 12121213232 | integer | 生成时间
Response.GroupList.0.CategoryName | wewer | string | 分组名称
Response.GroupList.0.ImageUrl | http://1234.jpg | string | 图片地址
Response.TotalCount | 10 | integer | 总数
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID
## 18、获取NFT详情
```text
获取NFT详情
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetNFTDetail&TokenId=10010&ContractAddr=1213343454545

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetNFTDetail | String | 是 | 方法名称
TokenId | 10010 | String | 是 | TokenID
ContractAddr | 1213343454545 | String | 是 | 合约地址



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": {
			"ContractName": "goErc20",
			"ContractAddr": "aba31ce4cd49f08073d2f115eb12610544242ff9",
			"TokenId": "1234",
			"OwnerAddr": "论车见省更真南须铁素状了或理强片象之",
			"AddrType": "同响斗由关社易越想北改",
			"Timestamp": 1253848085710939,
			"IsViolation": true,
			"Metadata": {
				"Name": "图机列种号么点铁最",
				"Author": "",
				"ImageUrl": "流它例族有积",
				"Description": "技场应毛",
				"SeriesHash": "主段地阶点听系部支空白天响压格取史难三"
			}
		},
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.Data | - | object | -
Response.Data.ContractName | goErc20 | string | 合约名称
Response.Data.ContractAddr | aba31ce4cd49f08073d2f115eb12610544242ff9 | string | 合约地址
Response.Data.TokenId | 1234 | string | TokenId
Response.Data.OwnerAddr | 论车见省更真南须铁素状了或理强片象之 | string | 持有地址
Response.Data.AddrType | 同响斗由关社易越想北改 | string | 持有地址类型
Response.Data.Timestamp | 1253848085710939 | integer | 上链时间
Response.Data.IsViolation | true | boolean | 是否违规（true：违规）包括图片或文本违规
Response.Data.Metadata | - | object | -
Response.Data.Metadata.Name | 图机列种号么点铁最 | string | 作品名称
Response.Data.Metadata.Author | - | string | 作者
Response.Data.Metadata.ImageUrl | 流它例族有积 | string | 图片地址
Response.Data.Metadata.Description | 技场应毛 | string | 作品描述
Response.Data.Metadata.SeriesHash | 主段地阶点听系部支空白天响压格取史难三 | string | 作品Hash
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | -
## 19、获取FT合约列表
```text
分页获取同质化合约列表
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetFTContractList&Limit=10&Offset=0&ContractKey=EVM_ERC2021212

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetFTContractList | String | 是 | 方法名称
Limit | 10 | String | 是 | 分页条
Offset | 0 | String | 是 | 页码
ContractKey | EVM_ERC2021212 | String | 否 | 合约名称/合约地址



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"ContractName": "goErc20",
				"ContractSymbol": "Erc20",
				"ContractAddr": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"ContractType": "ERC20",
				"TxNum": 123,
				"TotalSupply": "12344",
				"HolderCount": 123,
				"Timestamp": 1699242950
			}
		],
		"TotalCount": 10,
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.ContractName | goErc20 | string | 合约名称
Response.GroupList.0.ContractSymbol | Erc20 | string | 合约简称
Response.GroupList.0.ContractAddr | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | 合约地址
Response.GroupList.0.ContractType | ERC20 | string | 合约类型
Response.GroupList.0.TxNum | 123 | integer | 累计交易数
Response.GroupList.0.TotalSupply | 12344 | string | 发行总量
Response.GroupList.0.HolderCount | 123 | integer | 持有人数
Response.GroupList.0.Timestamp | 1699242950 | integer | 创建时间
Response.TotalCount | 10 | integer | 总数
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID
## 20、获取FT合约详情
```text
获取同质化合约详情
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetFTContractDetail&ContractAddr=aba31ce4cd49f08073d2f115eb12610544242ff9

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetFTContractDetail | String | 是 | 方法名称
ContractAddr | aba31ce4cd49f08073d2f115eb12610544242ff9 | String | 是 | 合约地址



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": {
			"ContractName": "goErc20",
			"ContractSymbol": "erc",
			"ContractAddr": "aba31ce4cd49f08073d2f115eb12610544242ff9",
			"ContractStatus": 0,
			"ContractType": "ERC20",
			"RuntimeType": "ERC",
			"Version": "1.0",
			"TxId": "17a05af0cfc1826fca85dc0d21083824de9ec658943c42738297fd372fc56d7f",
			"CreateSender": "171262347a59fded92021a32421a5dad05424e03",
			"CreatorAddr": "171262347a59fded92021a32421a5dad05424e03",
			"CreatorAddrBNS": "最也需安议何心处能查说教",
			"TotalSupply": "123445",
			"HolderCount": 12,
			"CreateTimestamp": 1702460649,
			"UpdateTimestamp": 1702460649,
			"TxNum": "非她切物受保号做会响主决些正到候求集同度声非"
		},
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.Data | - | object | -
Response.Data.ContractName | goErc20 | string | 合约名称
Response.Data.ContractSymbol | erc | string | token简称
Response.Data.ContractAddr | aba31ce4cd49f08073d2f115eb12610544242ff9 | string | 合约地址
Response.Data.ContractStatus | - | integer | 合约状态(0:正常，1:冻结，2:注销)
Response.Data.ContractType | ERC20 | string | 合约类型
Response.Data.RuntimeType | ERC | string | 虚拟机类型
Response.Data.Version | 1.0 | string | 合约版本
Response.Data.TxId | 17a05af0cfc1826fca85dc0d21083824de9ec658943c42738297fd372fc56d7f | string | 创建交易
Response.Data.CreateSender | 171262347a59fded92021a32421a5dad05424e03 | string | 创建用户ID
Response.Data.CreatorAddr | 171262347a59fded92021a32421a5dad05424e03 | string | 创建用户地址
Response.Data.CreatorAddrBNS | 最也需安议何心处能查说教 | string | 创建用户地址BNS
Response.Data.TotalSupply | 123445 | string | 发行量
Response.Data.HolderCount | 12 | integer | 持有人数
Response.Data.CreateTimestamp | 1702460649 | integer | 创建时间
Response.Data.UpdateTimestamp | 1702460649 | integer | 更新时间
Response.Data.TxNum | 非她切物受保号做会响主决些正到候求集同度声非 | string | 交易数
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | -
## 21、获取FT合约持仓列表
```text
分页获取同质化合约持仓列表
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chain1&cmb=GetFTPositionList&Offset=0&Limit=10

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ContractAddr | aba31ce4cd49f08073d2f115eb12610544242ff9 | String | 否 | 合约地址
OwnerAddr | b1dbe7f50d1d82a0318c38c5835dd2a5038dfce6 | String | 否 | 持仓地址
ChainId | chain1 | String | 是 | 链ID
cmb | GetFTPositionList | String | 是 | 方法名称
Offset | 0 | String | 是 | 页码
Limit | 10 | String | 是 | 分页条



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"AddrType": 1,
				"OwnerAddr": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"OwnerAddrBNS": "中行非党具场",
				"ContractName": "界其共义线持技结办我",
				"ContractAddr": "空你",
				"ContractType": "分斗变看米断周值管火理",
				"ContractSymbol": "过清",
				"Amount": "122.333",
				"HoldRatio": "12.3445%",
				"HoldRank": 2
			}
		],
		"TotalCount": 10,
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.AddrType | 1 | integer | 持仓地址类型：0:合约地址，1:用户地址
Response.GroupList.0.OwnerAddr | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | 持仓地址
Response.GroupList.0.OwnerAddrBNS | 中行非党具场 | string | 持仓地址BNS
Response.GroupList.0.ContractName | 界其共义线持技结办我 | string | 合约名称
Response.GroupList.0.ContractAddr | 空你 | string | 合约地址
Response.GroupList.0.ContractType | 分斗变看米断周值管火理 | string | 合约类型
Response.GroupList.0.ContractSymbol | 过清 | string | token符号
Response.GroupList.0.Amount | 122.333 | string | 持有量
Response.GroupList.0.HoldRatio | 12.3445% | string | 持有比例
Response.GroupList.0.HoldRank | 2 | integer | 持有排名
Response.TotalCount | 10 | integer | 总数
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID
## 22、获取FT流转列表
```text
分页获取同质化交易流转列表
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetFTTransferList&Offset=0&Limit=10&TxId=17a09dcaf9a2484dcad8e7a648c1923b7dbe731c43574a6e9f0d6d7372048619&ContractAddr=aba31ce4cd49f08073d2f115eb12610544242ff9&UserAddr=18fc4e7429af8419d5bb307e34db398b9a2331c6

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetFTTransferList | String | 是 | 方法名称
Offset | 0 | String | 是 | 页码
Limit | 10 | String | 是 | 分页条
TxId | 17a09dcaf9a2484dcad8e7a648c1923b7dbe731c43574a6e9f0d6d7372048619 | String | 否 | 交易ID
ContractAddr | aba31ce4cd49f08073d2f115eb12610544242ff9 | String | 否 | 合约地址
UserAddr | 18fc4e7429af8419d5bb307e34db398b9a2331c6 | String | 否 | 流转地址



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"TxId": "e4b090ef87c24f3e9ab52b68b57c8e7dcc7ab73459974bcebafbc8a4b745f150",
				"ContractName": "goErc20",
				"ContractAddr": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"ContractMethod": "transfer",
				"From": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"To": "6cd144e330edbe27f82bb44dbd06283836f000a3",
				"Amount": "123.3434",
				"Timestamp ": 1698917816
			}
		],
		"TotalCount": 10,
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.TxId | e4b090ef87c24f3e9ab52b68b57c8e7dcc7ab73459974bcebafbc8a4b745f150 | string | 交易ID
Response.GroupList.0.ContractName | goErc20 | string | 合约名称
Response.GroupList.0.ContractAddr | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | 合约地址
Response.GroupList.0.ContractMethod | transfer | string | 合约方法
Response.GroupList.0.From | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | from地址
Response.GroupList.0.To | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | to地址
Response.GroupList.0.Amount | 123.3434 | string | 转账ETH
Response.GroupList.0.Timestamp | 1698917816 | integer | 上链时间
Response.TotalCount | 10 | integer | 总数
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID
## 23、其他接口
```text
暂无描述
```
#### Header参数
参数名 | 示例值 | 参数描述
--- | --- | ---
暂无参数
#### Query参数
参数名 | 示例值 | 参数描述
--- | --- | ---
暂无参数
#### Body参数
参数名 | 示例值 | 参数描述
--- | --- | ---
暂无参数



## 24、交易详情
```text
获取账户详情
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chainmaker_pk&cmb=GetTxDetail&TxId=123456

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chainmaker_pk | String | 是 | 链ID
cmb | GetTxDetail | String | 是 | 方法名称
TxId | 123456 | String | 是 | 交易id



#### 成功响应示例
```javascript
{
    "Response": {
        "Data": {
            "TxId": "17beb3de53797af1ca2918527ca19ed3b9e99ed421434c27bf6bceb2cd64075f",
            "BlockHash": "c5b1eef811a933498d3337cf7504209f568e6eecbb492adab8f656946bbaf0de",
            "BlockHeight": 296,
            "Sender": "e97e032d9fbc9f282cdb779fb365afea70bcb1d4",
            "SenderOrgId": "public",
            "ContractName": "DID1234",
            "ContractNameBak": "DID1234",
            "ContractAddr": "9af6a4d25738a41224303a1dc0789fa49eb11ba0",
            "ContractMessage": "Success",
            "ContractVersion": "1.0",
            "TxStatusCode": "SUCCESS",
            "ContractResultCode": 0,
            "ContractResult": "",
            "RwSetHash": "69844a54d3f72eb36c2322332ff5f287b8f2dd7a5c01cbb6feaf4d0b6378a2b2",
            "ContractMethod": "INIT_CONTRACT",
            "ContractParameters": "",
            "Endorsement": "",
            "TxType": "INVOKE_CONTRACT",
            "Timestamp": 1711002675,
            "UserAddr": "e97e032d9fbc9f282cdb779fb365afea70bcb1d4",
            "UserAddrBns": "",
            "ContractRead": "[{\"index\":0,\"key\":\"Contract:DID1234\",\"value\":\"\"}]",
            "ContractWrite": "",
            "GasUsed": 19002,
            "Payer": "",
            "PayerBns": "",
            "Event": "[{\"index\":0,\"key\":\"SetDidDocument\",\"value\":\"did:cnbn:cnbn,{\\\"@context\\\":[\\\"https://www.w3.org/ns/did/v1\\\"],\\\"authentication\\\":[\\\"did:cnbn:cnbn#key-1\\\"],\\\"controller\\\":[\\\"did:cnbn:cnbn\\\"],\\\"created\\\":\\\"2024-02-28T15:37:18+08:00\\\",\\\"id\\\":\\\"did:cnbn:cnbn\\\",\\\"service\\\":[{\\\"id\\\":\\\"did:cnbn:cnbn#service-1\\\",\\\"serviceEndpoint\\\":\\\"http://192.168.1.181:30002/api/v1\\\",\\\"type\\\":\\\"IssuerService\\\"}],\\\"updated\\\":\\\"2024-02-28T15:37:18+08:00\\\",\\\"verificationMethod\\\":[{\\\"address\\\":\\\"8c52c5b8c4ea7d837acc891c3d7eb5fc4f6d2282\\\",\\\"controller\\\":\\\"did:cnbn:cnbn\\\",\\\"id\\\":\\\"did:cnbn:cnbn#key-1\\\",\\\"publicKeyPem\\\":\\\"-----BEGIN PUBLIC KEY-----\\\\nMFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAEYt5ptt0R2pLyBb58gIDjggjNlg6W\\\\nmHhmA5QGutQcfcv3G5M1AaEbeRG3QJeQkECIfiJ7sX4+CYdgIAxpQrZmzA==\\\\n-----END PUBLIC KEY-----\\\\n\\\",\\\"type\\\":\\\"SM2VerificationKey2022\\\"}]}\"}]",
            "RuntimeType": "DOCKER_GO",
            "ShowStatus": 0
        },
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | - | object | -
Response.Data.TxId | 统点每少 | string | -
Response.Data.BlockHash | 名需月至就子历 | string | -
Response.Data.BlockHeight | -3829859871894292 | integer | -
Response.Data.Sender | 走调有展从主北导应目 | string | -
Response.Data.SenderOrgId | 较级消林单流量 | string | -
Response.Data.ContractName | - | string | -
Response.Data.ContractNameBak | 效空体改连办干只于干派象着是战育交海车感 | string | -
Response.Data.ContractAddr | 们很指按在月众原 | string | -
Response.Data.ContractMessage | 交位按根 | string | -
Response.Data.ContractVersion | 备题层区比委细文现离道定定为应教 | string | -
Response.Data.TxStatusCode | 律 | string | -
Response.Data.ContractResultCode | 3775190835749747 | integer | -
Response.Data.ContractResult | 思需六子层 | string | -
Response.Data.RwSetHash | 选自养问连南界 | string | -
Response.Data.ContractMethod | 华大务解老以业报直强都我质 | string | -
Response.Data.ContractParameters | 铁走半织 | string | -
Response.Data.Endorsement | 边酸器装加 | string | -
Response.Data.TxType | 等六老 | string | -
Response.Data.Timestamp | 5359380928501503 | integer | -
Response.Data.UserAddr | 速器老 | string | 发起用户地址
Response.Data.UserAddrBns | 大际 | string | 发起用户BNS
Response.Data.ContractRead | 又 | string | -
Response.Data.ContractWrite | 老究积流术与 | string | -
Response.Data.GasUsed | 6111495808077567 | integer | -
Response.Data.Payer | 率示 | string | -
Response.Data.PayerBns | 三然龙活 | string | gas代付地址BNS
Response.Data.Event | 去目想 | string | -
Response.Data.RuntimeType | 必公响圆细置 | string | -
Response.Data.ShowStatus | 2399591123875923 | integer | -
Response.RequestId | 以花表则值量共 | string | -
## 25、最新交易列表
```text
获取最新的10条交易记录
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetLatestTxList

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetLatestTxList | String | 是 | 方法名称



#### 成功响应示例
```javascript
{
    "Response": {
        "GroupList": [
            {
                "Id": 1,
                "TxId": "2797b503f7ec4ea3a2e50e4337dcd554e05adf5b2e3c40a8a01c4833b943c534",
                "BlockHash": "6797c0731bab118db789627555094d2b0a56b7f614a59f7a7ce8da9fb0740509",
                "BlockHeight": 379,
                "Status": "成功",
                "Timestamp": 1711078328,
                "ContractName": "ACCOUNT_MANAGER",
                "ContractNameBak": "ACCOUNT_MANAGER",
                "ContractAddr": "a564b97f89f1c64b4ae12465a5870e879ffa0d3f",
                "Sender": "e8d472f7313891a1a34976a0b2b60ff6283d0a26",
                "UserAddr": "e8d472f7313891a1a34976a0b2b60ff6283d0a26",
                "UserAddrBns": "",
                "GasUsed": 0
            },
            {
                "Id": 2,
                "TxId": "17bef8ac96ff7eedca2ce3a131a44d90e9c9035a9da14c3f8f18f667109df98c",
                "BlockHash": "6797c0731bab118db789627555094d2b0a56b7f614a59f7a7ce8da9fb0740509",
                "BlockHeight": 379,
                "Status": "成功",
                "Timestamp": 1711078328,
                "ContractName": "DID1234",
                "ContractNameBak": "DID1234",
                "ContractAddr": "9af6a4d25738a41224303a1dc0789fa49eb11ba0",
                "Sender": "8c52c5b8c4ea7d837acc891c3d7eb5fc4f6d2282",
                "UserAddr": "8c52c5b8c4ea7d837acc891c3d7eb5fc4f6d2282",
                "UserAddrBns": "",
                "GasUsed": 712
            },
            {
                "Id": 3,
                "TxId": "52f1a429e4054171abfc8dc726e5f7d768c7aea3c9694f7998292809dc622d5e",
                "BlockHash": "7a69eea88c4b1eedeb8ff8904da2a407dd431cf146a753c0d200f978012452d3",
                "BlockHeight": 378,
                "Status": "成功",
                "Timestamp": 1711078328,
                "ContractName": "ACCOUNT_MANAGER",
                "ContractNameBak": "ACCOUNT_MANAGER",
                "ContractAddr": "a564b97f89f1c64b4ae12465a5870e879ffa0d3f",
                "Sender": "7e2d16337700a113c465ede743cdb34faa86e7ce",
                "UserAddr": "7e2d16337700a113c465ede743cdb34faa86e7ce",
                "UserAddrBns": "",
                "GasUsed": 0
            },
            {
                "Id": 4,
                "TxId": "17bef8ac6b146428ca4be62ed3dac6c202908a2f1a704932988490f1da2b4d0b",
                "BlockHash": "7a69eea88c4b1eedeb8ff8904da2a407dd431cf146a753c0d200f978012452d3",
                "BlockHeight": 378,
                "Status": "成功",
                "Timestamp": 1711078327,
                "ContractName": "DID1234",
                "ContractNameBak": "DID1234",
                "ContractAddr": "9af6a4d25738a41224303a1dc0789fa49eb11ba0",
                "Sender": "8c52c5b8c4ea7d837acc891c3d7eb5fc4f6d2282",
                "UserAddr": "8c52c5b8c4ea7d837acc891c3d7eb5fc4f6d2282",
                "UserAddrBns": "",
                "GasUsed": 421
            },
            {
                "Id": 5,
                "TxId": "5b48906012074b44a22f33fab226b60b8969be8de3554229981fdfff76424774",
                "BlockHash": "ac44e3c30f86add2b678f4e46d9bd1ad7905218c540f5253de9d536d8ad9c5ea",
                "BlockHeight": 377,
                "Status": "成功",
                "Timestamp": 1711078327,
                "ContractName": "ACCOUNT_MANAGER",
                "ContractNameBak": "ACCOUNT_MANAGER",
                "ContractAddr": "a564b97f89f1c64b4ae12465a5870e879ffa0d3f",
                "Sender": "f074cb594541dc020cf3ad07825d2dce7db8b467",
                "UserAddr": "f074cb594541dc020cf3ad07825d2dce7db8b467",
                "UserAddrBns": "",
                "GasUsed": 0
            },
            {
                "Id": 6,
                "TxId": "17bef8ac494316f0ca30c2129befd08bbef528ad556246b1a740e09585a4a6e9",
                "BlockHash": "ac44e3c30f86add2b678f4e46d9bd1ad7905218c540f5253de9d536d8ad9c5ea",
                "BlockHeight": 377,
                "Status": "成功",
                "Timestamp": 1711078327,
                "ContractName": "DID1234",
                "ContractNameBak": "DID1234",
                "ContractAddr": "9af6a4d25738a41224303a1dc0789fa49eb11ba0",
                "Sender": "8c52c5b8c4ea7d837acc891c3d7eb5fc4f6d2282",
                "UserAddr": "8c52c5b8c4ea7d837acc891c3d7eb5fc4f6d2282",
                "UserAddrBns": "",
                "GasUsed": 445
            },
            {
                "Id": 7,
                "TxId": "0adce4e9c8ab4c7ebfe25ad58a64fbeebf3a95f681d74ec7bf77e5993bc9fd5f",
                "BlockHash": "3b3946e2585ea4c70597eeaf8d071fa4c9fb29e3d10208e2752a8ebf92d932c6",
                "BlockHeight": 376,
                "Status": "成功",
                "Timestamp": 1711078326,
                "ContractName": "ACCOUNT_MANAGER",
                "ContractNameBak": "ACCOUNT_MANAGER",
                "ContractAddr": "a564b97f89f1c64b4ae12465a5870e879ffa0d3f",
                "Sender": "6094bfc00fc238b8ba11e303cdefb71de4a2b2df",
                "UserAddr": "6094bfc00fc238b8ba11e303cdefb71de4a2b2df",
                "UserAddrBns": "",
                "GasUsed": 0
            },
            {
                "Id": 8,
                "TxId": "17bef8ac2935b73bca88e6790f18d3982ed5f09503854ca08f63481f6dc954d7",
                "BlockHash": "3b3946e2585ea4c70597eeaf8d071fa4c9fb29e3d10208e2752a8ebf92d932c6",
                "BlockHeight": 376,
                "Status": "成功",
                "Timestamp": 1711078326,
                "ContractName": "DID1234",
                "ContractNameBak": "DID1234",
                "ContractAddr": "9af6a4d25738a41224303a1dc0789fa49eb11ba0",
                "Sender": "8c52c5b8c4ea7d837acc891c3d7eb5fc4f6d2282",
                "UserAddr": "8c52c5b8c4ea7d837acc891c3d7eb5fc4f6d2282",
                "UserAddrBns": "",
                "GasUsed": 423
            },
            {
                "Id": 9,
                "TxId": "2676a4c767ea472ebff6225de10c16fd95137b3fcb714214a62da32094ddc5bf",
                "BlockHash": "08b95ad7d463e10d2e695c2a44f31257db4e03e3df3c5e9da1844052f6a735e1",
                "BlockHeight": 375,
                "Status": "成功",
                "Timestamp": 1711075301,
                "ContractName": "ACCOUNT_MANAGER",
                "ContractNameBak": "ACCOUNT_MANAGER",
                "ContractAddr": "a564b97f89f1c64b4ae12465a5870e879ffa0d3f",
                "Sender": "7e2d16337700a113c465ede743cdb34faa86e7ce",
                "UserAddr": "7e2d16337700a113c465ede743cdb34faa86e7ce",
                "UserAddrBns": "",
                "GasUsed": 0
            },
            {
                "Id": 10,
                "TxId": "17bef5ebaa4ba8cecaced7afbfdbda9ea6bacf7165034b8985f5eef7135312f4",
                "BlockHash": "08b95ad7d463e10d2e695c2a44f31257db4e03e3df3c5e9da1844052f6a735e1",
                "BlockHeight": 375,
                "Status": "成功",
                "Timestamp": 1711075300,
                "ContractName": "DID1234",
                "ContractNameBak": "DID1234",
                "ContractAddr": "9af6a4d25738a41224303a1dc0789fa49eb11ba0",
                "Sender": "8c52c5b8c4ea7d837acc891c3d7eb5fc4f6d2282",
                "UserAddr": "8c52c5b8c4ea7d837acc891c3d7eb5fc4f6d2282",
                "UserAddrBns": "",
                "GasUsed": 631
            }
        ],
        "TotalCount": 10,
        "RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
    }
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.Id | 7969300059635883 | integer | -
Response.GroupList.0.TxId | 回 | string | -
Response.GroupList.0.BlockHash | 据 | string | -
Response.GroupList.0.BlockHeight | -5142500384308315 | integer | -
Response.GroupList.0.Status | 政十代列外海做少验原派又想因身门专效克前基安周 | string | -
Response.GroupList.0.Timestamp | -7129110104921755 | integer | -
Response.GroupList.0.ContractName | 八问算装转 | string | -
Response.GroupList.0.ContractAddr | 务不四族办新易个把子可平电候近反记白组好 | string | -
Response.GroupList.0.Sender | 两没 | string | 交易发起用户id
Response.GroupList.0.UserAddr | 例克值又军用义越青其元半方把该二信利于采到合三写必 | string | 交易发起用户地址
Response.GroupList.0.UserAddrBns | 器公力间们相队众必 | string | 交易发起用户BNS
Response.GroupList.0.GasUsed | 6319830960183175 | integer | -
Response.TotalCount | -337861162145264 | integer | -
Response.RequestId | 同就来经传国改且算题身八斗电很代根至意 | string | -
## 26、首页-最新合约列表
```text
获取最新的10条交易记录
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker?ChainId=chain1&cmb=GetLatestContractList

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | -
cmb | GetLatestContractList | String | 是 | -



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"Id": -7189614234881695,
				"TxId": "打转",
				"BlockHash": "把龙保特建铁满变先着管确使内志",
				"BlockHeight": -643004647503852,
				"Status": "府议于示验律万",
				"Timestamp": 2172909803875619,
				"ContractName": "",
				"ContractAddr": "类治式次员",
				"ContractType": "许直共回八关战些想清科准较布土传",
				"Sender": "",
				"SenderAddr": "花",
				"SenderAddrBNS": "种劳高能代采江通青办类力战型严老条间能定大",
				"GasUsed": 5953137758179351
			}
		],
		"TotalCount": 5470508674081671,
		"RequestId": "就府行生实计"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.Id | -7189614234881695 | integer | -
Response.GroupList.0.TxId | 打转 | string | -
Response.GroupList.0.BlockHash | 把龙保特建铁满变先着管确使内志 | string | -
Response.GroupList.0.BlockHeight | -643004647503852 | integer | -
Response.GroupList.0.Status | 府议于示验律万 | string | -
Response.GroupList.0.Timestamp | 2172909803875619 | integer | -
Response.GroupList.0.ContractName | - | string | -
Response.GroupList.0.ContractAddr | 类治式次员 | string | -
Response.GroupList.0.ContractType | 许直共回八关战些想清科准较布土传 | string | 合约类型
Response.GroupList.0.Sender | - | string | -
Response.GroupList.0.SenderAddr | 花 | string | 创建用户地址
Response.GroupList.0.SenderAddrBNS | 种劳高能代采江通青办类力战型严老条间能定大 | string | 创建用户BNS
Response.GroupList.0.GasUsed | 5953137758179351 | integer | -
Response.TotalCount | 5470508674081671 | integer | -
Response.RequestId | 就府行生实计 | string | -



## 27、存证合约详情
```text
获取存在合约hash列表
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetEvidenceContract&Offset=0&Limit=10&SenderAddrs=171262347a59fded92021a32421a5dad05424e03,123&Hashs=a412,233&TxId=179ca990d4da1134ca9fac9b94fddf44fed1d6f02318420e9877f04beabdb018

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetEvidenceContract | String | 是 | 方法名称
Offset | 0 | Integer | 是 | 页码
Limit | 10 | String | 是 | 条数
SenderAddrs | 171262347a59fded92021a32421a5dad05424e03,123 | String | 否 | 创建账户地址列表
Hashs | a412,233 | String | 是 | hash列表
TxId | 179ca990d4da1134ca9fac9b94fddf44fed1d6f02318420e9877f04beabdb018 | String | 是 | 交易ID



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"Id": "热列些办六率集就持上号包正水史",
				"UserAddr": "通容养龙提新人前人眼包北界根员是",
				"PayerAddr": "二里",
				"SenderAddr": "接间回它却车前县务规便连但",
				"MetaData": "家专米",
				"Hash": "拉其",
				"Timestamp": -4311593605584640,
				"BlockHeight": 12
			}
		],
		"TotalCount": 863666302166503,
		"RequestId": "已"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.Id | 热列些办六率集就持上号包正水史 | string | EvidenceId
Response.GroupList.0.UserAddr | 通容养龙提新人前人眼包北界根员是 | string | 链ID
Response.GroupList.0.PayerAddr | 二里 | string | 交易ID
Response.GroupList.0.SenderAddr | 接间回它却车前县务规便连但 | string | 创建账户地址
Response.GroupList.0.MetaData | 家专米 | string | MetaData
Response.GroupList.0.Hash | 拉其 | string | hash
Response.GroupList.0.Timestamp | -4311593605584640 | integer | 上链时间
Response.GroupList.0.BlockHeight | 12 | integer | 区块高度
Response.TotalCount | 863666302166503 | integer | -
Response.RequestId | 已 | string | -
## 28、获取gas列表
```text
获取Gas列表
```


#### 接口URL
> http://{{explore-k8s}}/chainmaker/?ChainId=chainmaker_testnet_pk&cmb=GetGasList&Offset=0&Limit=10&UserAddrs=

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chainmaker_testnet_pk | String | 是 | 链ID
cmb | GetGasList | String | 是 | 方法名称
Offset | 0 | Integer | 是 | 页码
Limit | 10 | String | 是 | 条数
UserAddrs | - | String | 否 | user地址列表



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"GasBalance": -8305665248833567,
				"GasTotal": 1391681634375795,
				"GasUsed": -5590931443457971,
				"Address": "色委重个华争起他争装其社",
				"ChainId": "象那电消报们门布或标",
				"Timestamp": 5633662367982359
			}
		],
		"TotalCount": -2180985831731708,
		"RequestId": ""
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.GasBalance | -8305665248833567 | integer | gas余额
Response.GroupList.0.GasTotal | 1391681634375795 | integer | gas总量
Response.GroupList.0.GasUsed | -5590931443457971 | integer | gas消耗
Response.GroupList.0.Address | 色委重个华争起他争装其社 | string | 用户地址
Response.GroupList.0.ChainId | 象那电消报们门布或标 | string | 链ID
Response.GroupList.0.Timestamp | 5633662367982359 | integer | -
Response.TotalCount | -2180985831731708 | integer | -
Response.RequestId | - | string | -
## 29、获取gas详情
```text
获取gas详情数据（浏览器和管理后台在用）
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chainmaker_pk&cmb=GetGasInfo&UserAddrs=eba4047337993ac5ed7c599e1c8c9087065e413b

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chainmaker_pk | String | 是 | 链ID
cmb | GetGasInfo | String | 是 | 方法名称
UserAddrs | eba4047337993ac5ed7c599e1c8c9087065e413b | String | 是 | 用户地址列表,，多个使用逗号分隔



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": {
			"GasBalance": 1223
		},
		"RequestId": "是技本想义该眼近飞切建阶儿"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | - | object | -
Response.Data.GasBalance | 1223 | integer | gas余额
Response.RequestId | 是技本想义该眼近飞切建阶儿 | string | -
## 30、获取gas消耗列表
```text
获取Gas列表
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=GetGasRecordList&Offset=0&Limit=10&UserAddrs=123,171262347a59fded92021a32421a5dad05424e03&BusinessType=2

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetGasRecordList | String | 是 | 方法名称
Offset | 0 | Integer | 是 | 页码
Limit | 10 | String | 是 | 条数
UserAddrs | 123,171262347a59fded92021a32421a5dad05424e03 | String | 否 | user地址列表
BusinessType | 2 | String | 否 | 1:gas充值，2:gas消耗(不传默认获取所有的gas充值，消耗数据)



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"GasAmount": 4021813523401847,
				"Address": "由科型及日又应社目质还本于还之",
				"PayerAddress": "切专保商长农记",
				"BusinessType": 3381503912129959,
				"ChainId": "断改你们最研题称按",
				"TxId": "这海计争此水林证解十子转给受把合进米没引教",
				"Timestamp": -855660177113004
			}
		],
		"TotalCount": 8311880862060947,
		"RequestId": "容力低人算面加即现专亲"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.GasAmount | 4021813523401847 | integer | gas消耗
Response.GroupList.0.Address | 由科型及日又应社目质还本于还之 | string | 账户地址
Response.GroupList.0.PayerAddress | 切专保商长农记 | string | 支付地址
Response.GroupList.0.BusinessType | 3381503912129959 | integer | 1：充值，2:消耗
Response.GroupList.0.ChainId | 断改你们最研题称按 | string | 链ID
Response.GroupList.0.TxId | 这海计争此水林证解十子转给受把合进米没引教 | string | 交易ID
Response.GroupList.0.Timestamp | -855660177113004 | integer | 上链时间
Response.TotalCount | 8311880862060947 | integer | -
Response.RequestId | 容力低人算面加即现专亲 | string | -




## 31、修改交易黑名单
```text
添加，删除黑名单，控制交易是否可见
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=ModifyTxBlackList&TxId=179ca9803a2d2f87ca3fdd90d2dc82a32771b721498140b1a900fd03a8a9e8f7&Status=0

#### 请求方式
> GET



#### 请求Header参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
x-api-key | AdSGavFKmbrzesqkZPQaVD2UGnqbbC | String | 是 | -
#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | ModifyTxBlackList | String | 是 | 方法名称
TxId | 179ca9803a2d2f87ca3fdd90d2dc82a32771b721498140b1a900fd03a8a9e8f7 | Integer | 是 | 交易ID
Status | 0 | String | 是 | 0:添加黑名单，1:移除黑名单



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": "ok",
		"RequestId": "能"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | ok | string | -
Response.RequestId | 能 | string | -
## 32、修改用户状态
```text
修改用户的状态
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=ModifyUserStatus&Address=171262347a59fded92021a32421a5dad05424e03&Status=0

#### 请求方式
> GET



#### 请求Header参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
x-api-key | AdSGavFKmbrzesqkZPQaVD2UGnqbbC | String | 是 | -
#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | ModifyUserStatus | String | 是 | 方法名称
Address | 171262347a59fded92021a32421a5dad05424e03 | Integer | 是 | 用户地址
Status | 0 | String | 是 | 0:正常，1:删除， 2:封禁



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": "ok",
		"RequestId": "能"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | ok | string | -
Response.RequestId | 能 | string | -
## 33、更新交易敏感词
```text
更新交易数据敏感词
涉及字段名称：
"ContractResult"
 "ContractMessage"
 "ContractParameters"
 "ReadSet"
 "WriteSet"
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=UpdateTxSensitiveWord&TxId=17a05682faf2fa99cac0f8780532c0418f6d2517d4ac4159b6a3edf269f4b4d0&Status=1&Column=&WarnMsg=

#### 请求方式
> GET



#### 请求Header参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
x-api-key | AdSGavFKmbrzesqkZPQaVD2UGnqbbC | String | 是 | -
#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | UpdateTxSensitiveWord | String | 是 | 方法名称
TxId | 17a05682faf2fa99cac0f8780532c0418f6d2517d4ac4159b6a3edf269f4b4d0 | Integer | 是 | 交易ID
Status | 1 | String | 是 | 0:敏感词隐藏，1:恢复数据
Column | - | String | 否 | 默认修改全部字段。
WarnMsg | - | String | 否 | 更新后数据，默认：「上链内容违反相关法律规定，内容已屏蔽」



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": "ok",
		"RequestId": "能"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | ok | string | -
Response.RequestId | 能 | string | -





## 34、更新存证合约敏感词
```text
更新存证数据敏感词
涉及字段名称：EvidenceMetaData
```
#### 接口状态
> 需修改

#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=UpdateEvidenceSensitiveWord&Hash=12213312312312&Status=0&WarnMsg=&Column=

#### 请求方式
> POST



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | UpdateEvidenceSensitiveWord | String | 是 | 方法名称
Hash | 12213312312312 | Integer | 是 | 存证hash
Status | 0 | String | 是 | 0:敏感词隐藏，1:恢复数据
WarnMsg | - | String | 是 | 更新后数据，默认：「上链内容违反相关法律规定，内容已屏蔽」
Column | - | String | 是 | 字段名称：EvidenceMetaData



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": "ok",
		"RequestId": "能"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | ok | string | -
Response.RequestId | 能 | string | -
## 35、更新合约名称敏感词
```text
更新存证数据敏感词
```
#### 接口状态
> 需修改

#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=UpdateContractNameSensitiveWord&ContractName=123&Status=0&WarnMsg=

#### 请求方式
> POST



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | UpdateContractNameSensitiveWord | String | 是 | 方法名称
ContractName | 123 | Integer | 是 | 合约名称
Status | 0 | String | 是 | 0:敏感词隐藏，1:恢复数据
WarnMsg | - | String | 是 | 更新后数据，默认：「上链内容违反相关法律规定，内容已屏蔽」



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": "ok",
		"RequestId": "能"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | ok | string | -
Response.RequestId | 能 | string | -
## 36、更新Token数据敏感词
```text
更新存证数据敏感词
涉及字段名称：EvidenceMetaData
```
#### 接口状态
> 需修改

#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=UpdateNFTSensitiveWord&TokenId=1234&Status=0&WarnMsg=&Column=

#### 请求方式
> POST



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | UpdateNFTSensitiveWord | String | 是 | 方法名称
TokenId | 1234 | Integer | 是 | TokenID
Status | 0 | String | 是 | 0:敏感词隐藏，1:恢复数据
WarnMsg | - | String | 是 | 更新后数据，默认：「上链内容违反相关法律规定，内容已屏蔽」
Column | - | String | 是 | 字段名称：TokenMetaData



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": "ok",
		"RequestId": "能"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | ok | string | -
Response.RequestId | 能 | string | -




## 37、主子链网配置
```text
主子链网配置。根据链配置确定是否展示主子链tag标签
```


#### 接口URL
> http://{{explore-k8s}}/chainmaker/?ChainId=chain1&cmb=GetMainCrossConfig

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 主链ID
cmb | GetMainCrossConfig | String | 是 | 方法名称



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": {
			"ShowTag": true
		},
		"RequestId": "dasadasdasdasdasd"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | - | object | -
Response.Data.ShowTag | true | boolean | 是否显示主子链网标签
（true：显示，false：不显示）
Response.RequestId | dasadasdasdasdasd | string | -
## 38、主子链/首页搜索
```text
首页搜索接口，返回数据是否存在，并返回type对应的查询结果
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chainmaker_pk&cmb=CrossSearch&Value=1

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chainmaker_pk | String | 是 | 主链ID
cmb | CrossSearch | String | 是 | 方法名称
Value | 1 | String | 是 | 搜索值（跨链id是数字，子链名称是字符串）



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": {
			"Type": 0,
			"Data": "6cd144e330edbe27f82bb44dbd06283836f000a3"
		},
		"RequestId": "没书织其九类金称后主院北很际真"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | - | object | -
Response.Data.Type | - | integer | 返回值类型（-1:未找到。0:跨链ID，1：子链id
Response.Data.Data | 6cd144e330edbe27f82bb44dbd06283836f000a3 | string | type对应的值
Response.RequestId | 没书织其九类金称后主院北很际真 | string | -
## 39、/主子链/首页详情数据
```text
获取首页详情数据，页面1min刷一次
（数据缓存40s）
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chain1&cmb=CrossOverviewData

#### 请求方式
> GET

#### Content-Type
> form-data

#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 主链ID
cmb | CrossOverviewData | String | 是 | 方法名称



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": {
			"TotalBlockHeight": 6816665082881503,
			"ShortestTime": -6048560609355775,
			" LongestTime": -2455753652176300,
			"AverageTime": 6932314693659315,
			"SubChainNum": -1583045839907028,
			"TxNum": 1926911444947099
		},
		"RequestId": "已代员见样料较"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | - | object | -
Response.Data.TotalBlockHeight | 6816665082881503 | integer | 总区块高度
Response.Data.ShortestTime | -6048560609355775 | integer | 最短时间（单位s）
Response.Data. LongestTime | -2455753652176300 | integer | 最长时间（单位s）
Response.Data.AverageTime | 6932314693659315 | integer | 平均时间（单位s）
Response.Data.SubChainNum | -1583045839907028 | integer | 子链数量
Response.Data.TxNum | 1926911444947099 | integer | 跨链交易数
Response.RequestId | 已代员见样料较 | string | -
## 40、/主子链/获取最新跨链交易列表
```text
分页获取跨链交易列表
（缓存40s）。根据FromIsMainChain，ToIsMainChain判断链是否是主链。如果是主链跳转到基础的浏览器中，如果是子链跳转到子链详情。
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chain1&cmb=CrossLatestTxList

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 主链ID
cmb | CrossLatestTxList | String | 是 | 方法名称



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"CrossId": "e4b090ef87c24f3e9ab52b68b57c8e7dcc7ab73459974bcebafbc8a4b745f150",
				"FromChainName": "来源链名称",
				"FromChainId": "来源链id",
				"FromIsMainChain": true,
				"ToChainName": "目标链名称",
				"ToChainId": "目标链id",
				"ToIsMainChain": false,
				"Status": "0",
				"Timestamp ": 12345667
			}
		],
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c",
		"TotalCount": 3422222705903463
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.CrossId | e4b090ef87c24f3e9ab52b68b57c8e7dcc7ab73459974bcebafbc8a4b745f150 | string | 跨链ID
Response.GroupList.0.FromChainName | 来源链名称 | string | 发起链名称
Response.GroupList.0.FromChainId | 来源链id | string | 发起链id
Response.GroupList.0.FromIsMainChain | true | boolean | 是否是主链
Response.GroupList.0.ToChainName | 目标链名称 | string | 目标链名称
Response.GroupList.0.ToChainId | 目标链id | string | 目标链id
Response.GroupList.0.ToIsMainChain | false | boolean | 是否是主链
Response.GroupList.0.Status | 0 | string | 跨链状态（0:成功，1:失败）
Response.GroupList.0.Timestamp | 12345667 | integer | 跨链发起时间
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID
Response.TotalCount | 3422222705903463 | integer | 总数
## 41、主子链/获取最新子链列表
```text
分页获取子链列表
缓存数据，新增子链更新缓存。
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chain1&cmb=CrossLatestSubChainList

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 主链ID
cmb | CrossLatestSubChainList | String | 是 | 方法名称



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"SubChainId": "1234",
				"SubChainName": "联盟链",
				"BlockHeight": 12,
				"CrossTxNum": 232,
				"CrossContractNum": 12,
				"Status": "0",
				"Timestamp": 1234566
			}
		],
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c",
		"TotalCount": 7361452864452887
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.SubChainId | 1234 | string | 子链ID
Response.GroupList.0.SubChainName | 联盟链 | string | 子链名称
Response.GroupList.0.BlockHeight | 12 | integer | 子链高度
Response.GroupList.0.CrossTxNum | 232 | integer | 跨链交易数
Response.GroupList.0.CrossContractNum | 12 | integer | 跨链合约数
Response.GroupList.0.Status | 0 | string | 子链状态（0:正常，1:异常）
Response.GroupList.0.Timestamp | 1234566 | integer | 同步时间
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID
Response.TotalCount | 7361452864452887 | integer | 总数
## 42、/主子链/获取跨链交易列表
```text
分页获取跨链交易列表。
根据FromIsMainChain，ToIsMainChain判断链是否是主链。如果是主链跳转到基础的浏览器中，如果是子链跳转到子链详情。
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chain1&cmb=GetCrossTxList&Offset=0&Limit=10&CrossId=&SubChainId=&FromChainName=&ToChainName=物流链&StartTime=&EndTime=

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 主链ID
cmb | GetCrossTxList | String | 是 | 方法名称
Offset | 0 | String | 是 | 页码
Limit | 10 | String | 是 | 分页条
CrossId | - | String | 否 | 跨链ID
SubChainId | - | String | 否 | 子链id(子链详情必传)
FromChainName | - | String | 否 | 发起链名称
ToChainName | 物流链 | String | 否 | 接收链名称
StartTime | - | Integer | 否 | 开始时间
EndTime | - | String | 否 | 结束时间



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"CrossId": "计",
				"FromChainName": "下",
				"FromChainId": "火器元真示专是料除更条设型人划",
				"FromIsMainChain": true,
				"ToChainName": "三会",
				"ToChainId": "属验确于争反济示起来算都对组党",
				"ToIsMainChain": false,
				"Status": "团得维件且广山团等查用拉约地",
				"Timestamp ": "张方但金之"
			}
		],
		"RequestId": "",
		"TotalCount": -6634530323870067
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.CrossId | 计 | string | 跨链ID
Response.GroupList.0.FromChainName | 下 | string | 来源链名称
Response.GroupList.0.FromChainId | 火器元真示专是料除更条设型人划 | string | 来源链ID
Response.GroupList.0.FromIsMainChain | true | boolean | 是否是主链
Response.GroupList.0.ToChainName | 三会 | string | 目标链名称
Response.GroupList.0.ToChainId | 属验确于争反济示起来算都对组党 | string | 目标链ID
Response.GroupList.0.ToIsMainChain | false | boolean | 是否是主链
Response.GroupList.0.Status | 团得维件且广山团等查用拉约地 | string | //跨链状态（0:进行中，1:成功，2:失败）
Response.GroupList.0.Timestamp | 张方但金之 | string | 跨链发起时间
Response.RequestId | - | string | -
Response.TotalCount | -6634530323870067 | integer | -
## 43、主子链/获取子链列表
```text
分页获取子链列表
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chain1&cmb=CrossSubChainList&Offset=0&Limit=10&SubChainId=chainmaker001222&SubChainName=

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | CrossSubChainList | String | 是 | 方法名称
Offset | 0 | String | 是 | 页码
Limit | 10 | String | 是 | 条数
SubChainId | chainmaker001222 | String | 否 | 子链id
SubChainName | - | String | 否 | 子链名称



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"CrossId": "e4b090ef87c24f3e9ab52b68b57c8e7dcc7ab73459974bcebafbc8a4b745f150",
				"SourceChainName": "来源链名称",
				"SourceChainId": "来源链id",
				"TargetChainName": "目标链名称",
				"TargetChainId": "目标链id",
				"Status": "0",
				"Timestamp ": "1234567"
			}
		],
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c",
		"TotalCount": 5015337229919087
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.GroupList | - | array | 列表数据
Response.GroupList.0 | - | object | -
Response.GroupList.0.CrossId | e4b090ef87c24f3e9ab52b68b57c8e7dcc7ab73459974bcebafbc8a4b745f150 | string | 跨链ID
Response.GroupList.0.SourceChainName | 来源链名称 | string | 来源链名称
Response.GroupList.0.SourceChainId | 来源链id | string | 来源链id
Response.GroupList.0.TargetChainName | 目标链名称 | string | 目标链名称
Response.GroupList.0.TargetChainId | 目标链id | string | 目标链id
Response.GroupList.0.Status | 0 | string | -
Response.GroupList.0.Timestamp | 1234567 | string | 跨链发起时间
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | 请求ID
Response.TotalCount | 5015337229919087 | integer | 总数
## 44、/主子链/获取子链详情
```text
获取子链详情
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chain1&cmb=CrossSubChainDetail&SubChainId=chainmaker001

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | CrossSubChainDetail | String | 是 | 方法名称
SubChainId | chainmaker001 | String | 是 | 子链Id



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": {
			"SubChainID": "1234",
			"SubChainName": "联盟链",
			"BlockHeight": -4687333168437631,
			"ChainType": 1,
			"CrossContractNum": 123,
			"CrossTxNum": 123,
			"GatewayId": "商备总军七所书说",
			"GatewayName": "边者许变统被那团效议式省亲美许今程定满接象造发二头收称易生接酸",
			"GatewayAddr": "特给资积人两权包即认打除便速南改"
		},
		"RequestId": "fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | 响应
Response.Data | - | object | -
Response.Data.SubChainID | 1234 | string | 子链ID
Response.Data.SubChainName | 联盟链 | string | 子链名称
Response.Data.BlockHeight | -4687333168437631 | integer | 区块高度
Response.Data.ChainType | 1 | integer | 区块链架构（1 长安链，2 fabric，3 bcos， 4eth，5+ 扩展）
Response.Data.CrossContractNum | 123 | integer | 跨链合约数
Response.Data.CrossTxNum | 123 | integer | 跨链交易数
Response.Data.GatewayId | 商备总军七所书说 | string | 跨链网关ID
Response.Data.GatewayName | 边者许变统被那团效议式省亲美许今程定满接象造发二头收称易生接酸 | string | 跨链网关名称
Response.Data.GatewayAddr | 特给资积人两权包即认打除便速南改 | string | 跨链网关地址
Response.RequestId | fa11a5f1-f45c-42bb-bbd2-a5aa684f7c1c | string | -
## 45、主子链/获取跨链交易详情
```text
获取跨链交易详情。根据IsMainChain判断是否是主链。如果是主链TxUrl会返回空，点击交易后根据ChainId跳转到基础浏览器的交易详情页面。如果是子链TxUrl如果不为空就跳转到相应的http地址，如果为空就不跳转。
```


#### 接口URL
> http://127.0.0.1:7999/chainmaker/?ChainId=chain1&cmb=GetCrossTxDetail&CrossId=0

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 链ID
cmb | GetCrossTxDetail | String | 是 | 方法名称
CrossId | 0 | Integer | 是 | 跨链ID



#### 成功响应示例
```javascript
{
	"Response": {
		"Data": {
			"CrossId": "龙正节给他证所器定广电门识存题总",
			"Status": 438515600621279,
			"CrossDuration": 1234,
			"Timestamp": -2387579318033864,
			"ContractName": "",
			"ContractMethod": "属布内术在际年素色话便",
			"Parameter": "样经造拉自常",
			"ContractResult": "人计八无何听便精可龙",
			"CrossDirection": {
				"FromChain": "前",
				"ToChain": "量机"
			},
			"FromChainInfo": {
				"ChainName": "六反",
				"ChainId": "任",
				"TxId": "照后他式型队准活出决完方区族之每书列布写点查次间",
				"TxStatus": "广小角素",
				"TxUrl": "研热百车即门目拉片话红究阶里明家组",
				"IsMainChain": "准该低当目问便所",
				"gas": "\"-\""
			},
			"ToChainInfo": {
				"ChainName": "利了文多值动确她持个什员件酸由意细米",
				"ChainId": "被属美作容样以酸联米他造代学斗达听圆开决起记林头装取和图组",
				"TxId": "",
				"TxStatus": "米见新群省活青技",
				"TxUrl": "实话了看温",
				"IsMainChain": "",
				"gas": "收号声天"
			}
		},
		"RequestId": "科王话土研员广只特常养南口得种全"
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.Data | - | object | -
Response.Data.CrossId | 龙正节给他证所器定广电门识存题总 | string | 跨链ID
Response.Data.Status | 438515600621279 | integer | 跨链状态（0:新建，1：待执行，2:待提交，3:确认结束，4:回滚结束）
Response.Data.CrossDuration | 1234 | integer | 跨链完成时间（单位s）
Response.Data.Timestamp | -2387579318033864 | integer | 跨链发起时间
Response.Data.ContractName | - | string | 合约名称
Response.Data.ContractMethod | 属布内术在际年素色话便 | string | 合约方法
Response.Data.Parameter | 样经造拉自常 | string | 合约入参
Response.Data.ContractResult | 人计八无何听便精可龙 | string | 合约执行结果json
Response.Data.CrossDirection | - | object | 跨链方向
Response.Data.CrossDirection.FromChain | 前 | string | 来源链名称
Response.Data.CrossDirection.ToChain | 量机 | string | 目标链名称
Response.Data.FromChainInfo | - | object | 来源链数据
Response.Data.FromChainInfo.ChainName | 六反 | string | 链名称
Response.Data.FromChainInfo.ChainId | 任 | string | 链ID
Response.Data.FromChainInfo.TxId | 照后他式型队准活出决完方区族之每书列布写点查次间 | string | 交易ID
Response.Data.FromChainInfo.TxStatus | 广小角素 | string | 交易状态
Response.Data.FromChainInfo.TxUrl | 研热百车即门目拉片话红究阶里明家组 | string | 交易跳转地址
Response.Data.FromChainInfo.IsMainChain | 准该低当目问便所 | string | 是否是主链（0:子链，1:主链）
Response.Data.FromChainInfo.gas | "-" | string | gas消耗
Response.Data.ToChainInfo | - | object | 目标链数据
Response.Data.ToChainInfo.ChainName | 利了文多值动确她持个什员件酸由意细米 | string | 链名称
Response.Data.ToChainInfo.ChainId | 被属美作容样以酸联米他造代学斗达听圆开决起记林头装取和图组 | string | 链ID
Response.Data.ToChainInfo.TxId | - | string | 交易ID
Response.Data.ToChainInfo.TxStatus | 米见新群省活青技 | string | 交易状态
Response.Data.ToChainInfo.TxUrl | 实话了看温 | string | 交易跳转地址
Response.Data.ToChainInfo.IsMainChain | - | string | 是否是主链（0:子链，1:主链）
Response.Data.ToChainInfo.gas | 收号声天 | string | gas消耗
Response.RequestId | 科王话土研员广只特常养南口得种全 | string | -
## 46、/主子链/获取子链历史跨链列表
```text
分页获取跨链交易列表
进行中的跨链交易不显示
```


#### 接口URL
> http://{{explore-host}}/chainmaker/?ChainId=chain1&cmb=SubChainCrossChainList&SubChainId=chainmaker001

#### 请求方式
> GET



#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
ChainId | chain1 | String | 是 | 主链ID
cmb | SubChainCrossChainList | String | 是 | 方法名称
SubChainId | chainmaker001 | String | 是 | 子链id(子链详情必传)



#### 成功响应示例
```javascript
{
	"Response": {
		"GroupList": [
			{
				"ChainName": "维约政权一大叫保收对展东些全界空",
				"ChainId": "素线给",
				"TxNum ": 904284532414191
			}
		],
		"RequestId": "此传部两支科切",
		"TotalCount": -2591062961529776
	}
}
```
参数名 | 示例值 | 参数类型 | 参数描述
--- | --- | --- | ---
Response | - | object | -
Response.GroupList | - | array | -
Response.GroupList.0 | - | object | -
Response.GroupList.0.ChainName | 维约政权一大叫保收对展东些全界空 | string | 链名称
Response.GroupList.0.ChainId | 素线给 | string | 链ID
Response.GroupList.0.TxNum | 904284532414191 | integer | 跨链交易数
Response.RequestId | 此传部两支科切 | string | -
Response.TotalCount | -2591062961529776 | integer | -



## 47、/订阅/新增订阅
```text
暂无描述
```
#### 接口状态
> 开发中

#### 接口URL
> http://{{explore-host}}/chainmaker?cmb=SubscribeChain

#### 请求方式
> POST

#### Content-Type
> json

#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
cmb | SubscribeChain | String | 是 | -
#### 请求Body参数
```javascript
{
    "ChainId": "chain2",
    "AuthType": "public",
    "NodeList": [
        {
            "Addr": "9.135.180.61:12301",
            "OrgCA": "",
            "TLSHostName": "chainmaker.org",
            "Tls": false
        }
    ],
    "UserKey": "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgelqvKuVT3RMtKqC3\nZCJI4Tac8tYkS50erVjLPe+rhBegCgYIKoEcz1UBgi2hRANCAARe9ivpYoa0yirj\nVPmCAm8XYbqhuj2RWdPhTr81/B15t/0Zp/oL0g48l4vWp/89X3n9S94g2VIYECbN\nVpYRWCRs\n-----END PRIVATE KEY-----\n",
    "HashType": 0
}
```



## 48、/订阅/修改订阅
```text
暂无描述
```
#### 接口URL
> http://{{explore-host}}/chainmaker?cmb=ModifySubscribe

#### 请求方式
> POST

#### Content-Type
> json

#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
cmb | ModifySubscribe | String | 是 | -
#### 请求Body参数
```javascript
{
    "ChainId": "chain2",
    "AuthType": "public",
    "NodeList": [
        {
            "Addr": "9.135.180.61:12301",
            "OrgCA": "",
            "TLSHostName": "chainmaker.org",
            "Tls": false
        },
        {
            "Addr": "9.135.180.61:12302",
            "OrgCA": "",
            "TLSHostName": "chainmaker.org",
            "Tls": false
        }
    ],
    "UserKey": "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgelqvKuVT3RMtKqC3\nZCJI4Tac8tYkS50erVjLPe+rhBegCgYIKoEcz1UBgi2hRANCAARe9ivpYoa0yirj\nVPmCAm8XYbqhuj2RWdPhTr81/B15t/0Zp/oL0g48l4vWp/89X3n9S94g2VIYECbN\nVpYRWCRs\n-----END PRIVATE KEY-----\n",
    "HashType": 0
}
```



## 49、订阅/暂停订阅
```text
暂无描述
```

#### 接口URL
> http://{{explore-host}}/chainmaker?cmb=CancelSubscribe

#### 请求方式
> POST

#### Content-Type
> json

#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
cmb | CancelSubscribe | String | 是 | -
#### 请求Body参数
```javascript
{
    "ChainId": "chain1",
    "Status": 0
}
```



## 50、/订阅/删除订阅
```text
暂无描述
```
#### 接口URL
> http://{{explore-host}}/chainmaker?cmb=DeleteSubscribe

#### 请求方式
> POST

#### Content-Type
> json

#### 请求Query参数
参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述
--- | --- | --- | --- | ---
cmb | DeleteSubscribe | String | 是 | -
#### 请求Body参数
```javascript
{
    "ChainId": "chain2"
}
```


