package main

import (
    "github.com/meshplus/gosdk/rpc"
    "github.com/meshplus/gosdk/account"
    "github.com/meshplus/gosdk/common"
    "github.com/meshplus/gosdk/abi"
    "fmt"
    "strings"
)

func main() {
    logger := common.GetLogger("main")

    // 新建账户
    accountJson, err := account.NewAccountSm2("123")
    if err != nil {
        logger.Error(err)
        return
    }
    fmt.Println("accountJson:", accountJson)

    key, err := account.NewAccountSm2FromAccountJSON(accountJson, "123")
    if err != nil {
        logger.Error(err)
        return
    }
    fmt.Println("account address:", key.GetAddress())

    // 创建RPC客户端
    rpcAPI := rpc.NewRPCWithPath("./conf")

    // 编译合约
    code, _ := common.ReadFileAsString("./smartcontract/solidity/CTISmartContract.sol")
    cr, stdErr := rpcAPI.CompileContract(code)
    if stdErr != nil {
        logger.Error(stdErr)
        return
    }
    fmt.Println("contract abi:", cr.Abi[0])

    // 部署合约
    tranDeploy := rpc.NewTransaction(key.GetAddress().String()).Deploy(cr.Bin[0]).DeployArgs(cr.Abi[0], int64(1), []byte("demo"))
    tranDeploy.Sign(key)
    txDeploy, stdErr := rpcAPI.DeployContract(tranDeploy)
    if stdErr != nil {
        logger.Error(stdErr)
        return
    }
    fmt.Println("contract address:", txDeploy.ContractAddress)

    // 调用合约方法
    ABI, _ := abi.JSON(strings.NewReader(cr.Abi[0]))

    // 准备参数
    ctiId := "test"
    ctiName := "test"
    publisher := "test"
    _type := uint64(1)
    data := "test"
    hash := "0x6653"
    dataSize := uint64(len([]byte(data)))
    value := uint64(100)
    chainId := "0x6653"

    // 打包参数
    packed, err := ABI.Pack("registerCTI", ctiId, ctiName, publisher, _type, data, hash, dataSize, value, chainId)
    if err != nil {
        logger.Error(err)
        return
    }

    tranInvoke := rpc.NewTransaction(key.GetAddress().String()).Invoke(txDeploy.ContractAddress, packed)
    tranInvoke.Sign(key)
    txInvoke, stdErr := rpcAPI.InvokeContract(tranInvoke)
    if stdErr != nil {
        logger.Error(stdErr)
        return
    }
    fmt.Println("invoke transaction hash:", txInvoke.TxHash)

    // 查询单个CTI
    packedQuery, err := ABI.Pack("queryCTI", ctiId)
    if err != nil {
        logger.Error(err)
        return
    }
    tranQuery := rpc.NewTransaction(key.GetAddress().String()).Invoke(txDeploy.ContractAddress, packedQuery)
    tranQuery.Sign(key)
    txQuery, stdErr := rpcAPI.InvokeContract(tranQuery)
    if stdErr != nil {
        logger.Error(stdErr)
        return
    }
    fmt.Println("query transaction hash:", txQuery.TxHash)

    // 解码查询结果
    var id string
    var name string
    var pub string
    var t uint64
    var d string
    var h string
    var ds uint64
    var v uint64
    var cid string
    result := []interface{}{&id, &name, &pub, &t, &d, &h, &ds, &v, &cid}
    if err = ABI.UnpackResult(&result, "queryCTI", txQuery.Ret); err != nil {
        logger.Error(err)
        return
    }
    fmt.Println("id, name, pub, t, d, h, ds, v, cid:", id, name, pub, t, d, h, ds, v, cid)

    // 查询所有CTI
    packedQueryAll, err := ABI.Pack("queryAllCTIs")
    if err != nil {
        logger.Error(err)
        return
    }
    tranQueryAll := rpc.NewTransaction(key.GetAddress().String()).Invoke(txDeploy.ContractAddress, packedQueryAll)
    tranQueryAll.Sign(key)
    txQueryAll, stdErr := rpcAPI.InvokeContract(tranQueryAll)
    if stdErr != nil {
        logger.Error(stdErr)
        return
    }
    fmt.Println("queryAll transaction hash:", txQueryAll.TxHash)

    // 解码查询所有结果
    var ids []string
    var names []string
    var types []uint64
    var datas []string
    var hashes []string
    var dataSizes []uint64
    var values []uint64
    var chainIds []string
    resultAll := []interface{}{&ids, &names, &types, &datas, &hashes, &dataSizes, &values, &chainIds}
    if err = ABI.UnpackResult(&resultAll, "queryAllCTIs", txQueryAll.Ret); err != nil {
        logger.Error(err)
        return
    }
    fmt.Println("ids, names, types, datas, hashes, dataSizes, values, chainIds:", ids, names, types, datas, hashes, dataSizes, values, chainIds)
}