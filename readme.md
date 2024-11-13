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
部署合约
```shell
./network.sh deployCC -ccn br-cti-contract -ccp ./fabric-contract -ccl go
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
peer chaincode invoke $TARGET_TLS_OPTIONS -C mychannel -n br-cti-contract -c '{"function":"InitLedger","Args":[]}'
```