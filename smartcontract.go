package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/meshplus/gosdk/abi"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
	"github.com/meshplus/gosdk/hvm"
	"github.com/meshplus/gosdk/rpc"
	"github.com/meshplus/gosdk/utils/java"
)

// DeployHVM 部署合约
func DeployHVM(api *rpc.RPC, key *account.SM2Key, jarPath string) (contractAddress string, err error) {
	// 获取账户的地址
	from := key.GetAddress()
	// 构造交易
	transaction := rpc.NewTransaction(from.Hex())
	// 设置交易：智能合约虚拟机类型为HVM虚拟机
	transaction.VMType(rpc.HVM)
	// 将bin设置到payload中
	payload, _ := rpc.DecompressFromJar(jarPath)

	transaction = transaction.Deploy(common.ToHex([]byte(payload)))
	// 交易签名
	transaction.Sign(key)
	// 部署合约
	res, err := api.DeployContract(transaction)
	if err != nil {
		return
	}
	// 获取已经部署的合约地址
	contractAddress = res.ContractAddress
	return
}

// UpgradeHVM 升级java合约
func UpgradeHVM(api *rpc.RPC, key *account.SM2Key, jarPath string, contractAddr string) (contractAddress string, err error) {
	// 获取账户的地址
	from := key.GetAddress()
	// 构造交易
	transaction := rpc.NewTransaction(from.Hex())
	// 设置交易：智能合约虚拟机类型为HVM虚拟机
	transaction.VMType(rpc.HVM)
	// 将bin设置到payload中
	payload, _ := rpc.DecompressFromJar(jarPath)

	transaction = transaction.Maintain(1, contractAddr, common.ToHex([]byte(payload)))
	// 交易签名
	transaction.Sign(key)
	// 升级合约
	res, err := api.MaintainContract(transaction)
	if err != nil {
		return
	}
	// 获取已经部署的合约地址
	contractAddress = res.ContractAddress
	return
}

// InvokeHVM 调用HVM合约
func InvokeHVM(api *rpc.RPC, key *account.SM2Key, contractAddress string, ABI hvm.Abi, functionName string, args string) (result string, err error) {
	// 获取账户的地址
	from := key.GetAddress()
	// 构造交易
	transaction := rpc.NewTransaction(from.Hex())
	// 设置交易：智能合约虚拟机类型为HVM虚拟机
	transaction.VMType(rpc.HVM)
	// 获取解析后的方法ABI
	abiParsed, err := ABI.GetBeanAbi(functionName)
	if err != nil {
		return
	}
	// 获取解析后的入参
	argsParsed, err := parseHVMArgs(args, abiParsed)
	if err != nil {
		return
	}
	// 生成payload
	payload, err := hvm.GenPayload(abiParsed, argsParsed...)
	// 设置调用交易
	transaction.Invoke(contractAddress, payload)
	// 交易签名
	transaction.Sign(key)
	res, err := api.InvokeContract(transaction)
	if err != nil {
		return
	}
	// 获取已经部署的合约地址
	txReceipt := res.Ret
	result = java.DecodeJavaResult(txReceipt)
	return
}

// 将HVM合约调用的入参从json格式转化成[]interface{}
func parseHVMArgs(args string, bean *hvm.BeanAbi) ([]interface{}, error) {
	var argsMap map[string]interface{}
	err := json.Unmarshal([]byte(args), &argsMap)
	if err != nil {
		return nil, fmt.Errorf("parse args from json to map error. args = %s", args)
	}
	var argsParsed []interface{}
	for _, b := range bean.Inputs {
		if s, ok := argsMap[b.Name]; ok {
			/*sString, sStringOk := s.(string)
			if !sStringOk {
				return nil, fmt.Errorf("arg is not string type, argName=%s, argVal=%v", b.Name, s)
			}*/
			argsParsed = append(argsParsed, s)
		}
	}
	return argsParsed, nil
}

// Deploy 部署合约
func Deploy(api *rpc.RPC, key *account.SM2Key, bin string, abi string, args string) (contractAddress string, err error) {
	// 获取账户的地址
	from := key.GetAddress()
	// 构造交易
	transaction := rpc.NewTransaction(from.Hex())
	// 设置交易：智能合约虚拟机类型为EVM虚拟机
	transaction.VMType(rpc.EVM)
	// 设置交易：部署合约
	transaction = transaction.Deploy(bin)
	// 解析ABI
	ABI, err := parseAbi(abi)
	if err != nil {
		return
	}
	// 解析args，函数名为""表示解析构造函数的入参
	argsParsed, err := parseArgs(args, ABI, "")
	if err != nil {
		return
	}
	// 如果有构造函数追加payload
	if argsParsed != nil {
		ok := setConstructorArgs(transaction, abi, argsParsed...)
		if !ok {
			err = fmt.Errorf("invalid args, please check your constructor args")
			return
		}
	}
	// 交易签名
	transaction.Sign(key)
	// 部署合约
	res, err := api.DeployContract(transaction)
	if err != nil {
		return
	}
	// 获取已经部署的合约地址
	contractAddress = res.ContractAddress
	return
}

// Upgrade 升级evm合约
func Upgrade(api *rpc.RPC, key *account.SM2Key, bin string, abi string, args string, contractAddr string) (contractAddress string, err error) {
	// 获取账户的地址
	from := key.GetAddress()
	// 构造交易
	transaction := rpc.NewTransaction(from.Hex())
	// 设置交易：智能合约虚拟机类型为EVM虚拟机
	transaction.VMType(rpc.EVM)
	// 设置交易：升级合约
	transaction = transaction.Maintain(1, contractAddr, bin)
	// 解析ABI
	ABI, err := parseAbi(abi)
	if err != nil {
		return
	}
	// 解析args，函数名为""表示解析构造函数的入参
	argsParsed, err := parseArgs(args, ABI, "")
	if err != nil {
		return
	}
	// 如果有构造函数追加payload
	if argsParsed != nil {
		ok := setConstructorArgs(transaction, abi, argsParsed...)
		if !ok {
			err = fmt.Errorf("invalid args, please check your constructor args")
			return
		}
	}
	// 交易签名
	transaction.Sign(key)
	// 部署合约
	res, err := api.MaintainContract(transaction)
	if err != nil {
		return
	}
	// 获取已经部署的合约地址
	contractAddress = res.ContractAddress
	return
}

// Invoke 调用合约
func Invoke(api *rpc.RPC, key *account.SM2Key, contractAddress string, abi string, functionName string, args string) (result string, err error) {
	// 获取账户的地址
	from := key.GetAddress()
	// 构造交易
	transaction := rpc.NewTransaction(from.Hex())
	// 设置交易：智能合约虚拟机类型为EVM虚拟机
	transaction.VMType(rpc.EVM)
	// 解析ABI
	ABI, err := parseAbi(abi)
	if err != nil {
		return
	}
	// 解析args
	argsParsed, err := parseArgs(args, ABI, functionName)
	if err != nil {
		return
	}
	// 包装payload
	payload, err, ok := setFunctionArgs(ABI, functionName, argsParsed...)
	if err != nil {
		return
	}
	if !ok {
		err = fmt.Errorf("invalid args, please check your invoke function args")
		return
	}
	// 设置调用函数
	transaction.Invoke(contractAddress, payload)
	// 交易签名
	transaction.Sign(key)
	res, err := api.InvokeContract(transaction)
	if err != nil {
		return
	}
	// 获取已经部署的合约地址
	txReceipt := res.Ret
	result, err = parseResult(txReceipt, ABI, functionName)
	return
}

// parseAbi 解析ABI
func parseAbi(rawAbi string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(rawAbi))
}

// parseArgs 将json形式的rawArgs 转化成 []interface{} 类型的合约入参
func parseArgs(rawArgs string, ABI abi.ABI, functionName string) ([]interface{}, error) {
	// argsRaw(json) 转化成 []interface{}
	var args []interface{}
	// 如果参数为空，返回nil
	if rawArgs == "" || rawArgs == "{}" {
		return nil, nil
	}
	// 将参数取出
	kv := make(map[string]interface{})
	err := json.Unmarshal([]byte(rawArgs), &kv)
	if nil != err {
		return nil, err
	}
	// 按照ABI对参数进行重新排序
	if functionName == "" {
		for _, v := range ABI.Constructor.Inputs {
			args = append(args, kv[v.Name])
		}
	} else {
		for _, v := range ABI.Methods[functionName].Inputs {
			args = append(args, kv[v.Name])
		}
	}
	return args, nil
}

// setConstructorArgs 部署合约时构造函数参数的设置
func setConstructorArgs(t *rpc.Transaction, abiString string, args ...interface{}) (ok bool) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(fmt.Sprintf("fail to setConstructorArgs. err=%v", err))
			return
		}
	}()

	t = t.DeployStringArgs(abiString, args...)
	if t == nil {
		return false
	}
	ok = true
	return
}

// setFunctionArgs 设置合约函数的入参
func setFunctionArgs(ABI abi.ABI, functionName string, args ...interface{}) (payload []byte, err error, ok bool) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(fmt.Sprintf("fail to setFunctionArgs. err=%v", err))
			return
		}
	}()

	payload, err = ABI.Encode(functionName, args...)
	ok = true
	return
}

// parseResult 解析合约调用的返回结果
func parseResult(txReceipt string, ABI abi.ABI, functionName string) (string, error) {
	if txReceipt == "0x0" {
		return "", nil
	}
	result, err := ABI.Decode(functionName, common.FromHex(txReceipt))
	if nil != err {
		return "", err
	}
	jsonResult, err := json.Marshal(result)
	if nil != err {
		return "", err
	}
	return string(jsonResult), nil
}

