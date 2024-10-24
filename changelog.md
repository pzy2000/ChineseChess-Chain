## 230版本更新

### 增加公钥模式


### 增加多链

```yml
chains:
  -
    chain_id: chainmaker_testnet_chain
    auth_type: permissionedwithcert
    nodes:
      -
        enable: true
        chain_id: chainmaker_testnet_chain
        org_id: org1.cmtestnet
        tls: true
        tls_host: common1.tls.org1.cmtestnet
        ca_paths: ../configs/crypto-config/org1.cmtestnet/certs/ca/org1.cmtestnet
        remotes: 172.21.0.6:12302
        user:
          priv_key_file: ../configs/crypto-config/org1.cmtestnet/certs/user/admin1/admin1.sign.key
          cert_file: ../configs/crypto-config/org1.cmtestnet/certs/user/admin1/admin1.sign.crt

  -
    chain_id: chainmaker_testnet_pk
    auth_type: public
    nodes:
      - enable: true
        chain_id: chainmaker_testnet_pk
        auth_type: public
        hash_type: SM3
        remotes: 152.136.217.46:17301
        user:
          priv_key_file: ../configs/crypto-config/node1/user/client1/client1.key
```

### 监控适配


### 多签模式显示正常合约

### 敏感词过滤

##### 根据合约查询，改为根据合约地址查询详情信息