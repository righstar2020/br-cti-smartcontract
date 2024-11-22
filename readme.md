## 说明
使用fabric 2.2.0 链码SDK
没有使用最新的版本
编译请在fabric-contract文件夹下进行

## 依赖说明
go 依赖以太坊加密库1.9.25
```shell
go get github.com/ethereum/go-ethereum/crypto@v1.9.25
```

## 链代码
可参考fabric-sample内教程
启动环境并选择couchdb数据库
```shell
./network.sh up createChannel -s couchdb
```

部署合约
```shell
./network.sh deployCC -ccn main_contract -ccp ../br-cti-smartcontract/fabric-contract -ccl go
```
证书配置
```shell
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/minter@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051
export TARGET_TLS_OPTIONS="-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
```
CLI执行链码函数
```shell
//调用命令-c后可替换
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n main_contract --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'


//初始化
-c '{"function":"InitLedger","Args":[]}'

//注册用户
 -c '{"function":"RegisterUserInfo","Args":["lxp","123456"]}'

//查询用户
 -c '{"function":"QueryUserInfo","Args":["8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"]}'

//修改用户
 -c '{"function":"UpdateUserInfo","Args":["123456","test","8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"]}'

//注册CTI
  -c '{"function":"RegisterCTIInfo","Args":["Example CTI Info", "1", "2", "1", "[\"Malware\", \"DDoS\", \"APT\"]", "[\"IP\", \"URL\", \"HASH\"]", "{\"type\":\"malware\",\"name\":\"Example Malware\",\"description\":\"This is a test stix data.\"}", "{\"uniqueIP\":100,\"maliciousURL\":50}", "Example description for CTI", "1024", "Qm12345abcde", "10", "100", "50", "privateKeyExample"]}'

//查询CTI
-c '{"function":"QueryCTIInfo","Args":["CTI_1"]}'

//分页查询
-c '{"function":"QueryModelInfoByModelIDWithPagination","Args":["MODEL_", "5", ""]}'
-c '{"function":"QueryCTIInfoByCTIIDWithPagination","Args":["CTI_", "5", ""]}'

//CTI类型查询
-c '{"function":"QueryCTIInfoByType","Args":["1"]}'

//根据私钥查询用户上传cti信息
-c '{"function":"QueryCTIInfoByPrivateKey","Args":["123456"]}'

//注册模型信息
-c '{"function":"RegisterModelInfo","Args":["test Model", "5G", "[\"Feature1\", \"Feature2\"]", "example_traffic_process_code", "Supervised Learning", "This is an ML model info.", "example_ml_train_code", "example_ipfs_hash", "CTI_1", "123456"]}'

//Modelid查询
-c '{"function":"QueryModelInfo","Args":["MODEL_1"]}'

//modeltype查询
-c '{"function":"QueryModelsByTrafficType","Args":["5G"]}'

// 查询用户所上传的模型信息
-c '{"function":"QueryModelsByPrivateKey","Args":["123456"]}'

//根据refctiid查询模型信息
-c '{"function":"QueryModelsByRefCTIId","Args":["CTI_1"]}'

//数据统计
 -c '{"function":"GetDataStatistics","Args":[]}'

//精简CTI信息
-c '{"function":"QueryCTISummaryInfoByCTIID","Args":["CTI_1"]}'



```
