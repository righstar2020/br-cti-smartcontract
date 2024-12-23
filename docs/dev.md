# 智能合约开发说明
## 1.说明
使用fabric 2.2.0 链码SDK
没有使用最新的版本
编译请在fabric-contract文件夹下进行

## 2.依赖说明
go 依赖以太坊加密库1.9.25
```shell
go get github.com/ethereum/go-ethereum/crypto@v1.9.25
```

## 3.链代码
### 3.1 环境启动
可参考fabric-sample内教程
启动环境并选择couchdb数据库
1.进入fabric-sample/test-network文件夹
```shell
cd ~/go/src/github.com/hyperledger/fabric-samples/test-network
```
2.停止网络(如网络已启动)
```shell
./network.sh down
```
3.启动环境并选择couchdb数据库
```shell
./network.sh up createChannel -s couchdb
```
### 3.2 合约部署
4.部署合约(每次更改都需要重启环境)
```shell
./network.sh deployCC -ccn main_contract -ccp ../br-cti-smartcontract/fabric-contract -ccl go
```
5.升级合约
```shell
./network.sh deployCC -ccn main_contract -ccp ../br-cti-smartcontract/fabric-contract -ccl go -ccv 1.0
```
6.查看合约版本
```shell
peer lifecycle chaincode querycommitted -C mychannel -n main_contract
```
7.环境变量(证书)配置
```shell
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export TARGET_TLS_OPTIONS="-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
```
8.初始化账本
```shell
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n main_contract --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'
```

### 3.3 链码执行
6.环境变量(证书)配置
```shell
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export TARGET_TLS_OPTIONS="-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
```


7.CLI执行链码函数

```shell
//调用命令-c后可替换
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n main_contract --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'
```

8.测试
```shell
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n main_contract --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"RegisterUserInfo","Args":["{\"userName\":\"测试用户\",\"publicKey\":\"用户公钥\"}"]}'
```


9.查询接口
下面 $TARGET_TLS_OPTIONS 有bug暂时无法使用，需要修改
```shell  
//初始化账本
peer chaincode invoke $TARGET_TLS_OPTIONS -C mychannel -n main_contract -c '{"function":"InitLedger","Args":[]}'

//查询用户信息
peer chaincode query $TARGET_TLS_OPTIONS -C mychannel -n main_contract -c '{"function":"QueryUserInfo","Args":["用户ID"]}'

//查询情报信息
-c '{"function":"QueryCTIInfo","Args":["情报ID"]}'

//根据情报Hash查询情报信息
-c '{"function":"QueryCTIInfoByCTIHash","Args":["情报Hash"]}'

//查询用户上传的情报
-c '{"function":"QueryCTIInfoByCreatorUserID","Args":["用户ID"]}'

//根据情报类型查询
-c '{"function":"QueryCTIInfoByType","Args":["情报类型"]}'

//查询用户积分信息
-c '{"function":"QueryUserPointInfo","Args":["用户ID"]}'

//模型信息分页查询
-c '{"function":"QueryModelInfoByModelIDWithPagination","Args":["模型ID前缀", "每页数量", "书签"]}'

//根据流量类型查询模型
-c '{"function":"QueryModelsByTrafficType","Args":["流量类型"]}'

//查询用户上传的模型
-c '{"function":"QueryModelsByUserID","Args":["用户ID"]}'

//根据关联情报查询模型
-c '{"function":"QueryModelsByRefCTIId","Args":["关联情报ID"]}'

//分页查询所有情报
-c '{"function":"QueryAllCTIInfoWithPagination","Args":["每页数量", "书签"]}'

//根据类型分页查询情报
-c '{"function":"QueryCTIInfoByTypeWithPagination","Args":["情报类型", "每页数量", "书签"]}'

//查询情报精简信息
-c '{"function":"QueryCTISummaryInfoByCTIID","Args":["情报ID"]}'

//获取数据统计信息
-c '{"function":"GetDataStatistics","Args":[]}'

//获取情报交易趋势
-c '{"function":"GetCTITrafficTrend","Args":["时间范围"]}'

//获取攻击类型排行
-c '{"function":"GetAttackTypeRanking","Args":[]}'

//获取IOCs类型分布
-c '{"function":"GetIOCsDistribution","Args":[]}'

//获取全球IOCs地理分布
-c '{"function":"GetGlobalIOCsDistribution","Args":[]}'

//获取系统概览数据
-c '{"function":"GetSystemOverview","Args":[]}'

//获取用户统计数据
-c '{"function":"GetUserStatistics","Args":["用户ID"]}'

//查询用户积分交易记录
-c '{"function":"QueryPointTransactions","Args":["用户ID"]}'
```
9.注册接口
```shell
//注册CTI
peer chaincode invoke $TARGET_TLS_OPTIONS -C mychannel -n main_contract -c '{"function":"RegisterCTIInfo","Args":["情报信息JSON"]}'
//注册用户信息
peer chaincode invoke $TARGET_TLS_OPTIONS -C mychannel -n main_contract -c '{"function":"RegisterUserInfo","Args":["用户信息JSON"]}'
```