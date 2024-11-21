package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	ctiContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/cti-contract"
	dataContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/data-contract"
	modelContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/model-contract"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
	userContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/user-contract"
	userPointContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/user-point-contract"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/utils"
)

// 主合约结构体
type MainContract struct {
	dataContract.DataContract
	contractapi.Contract
	modelContract.ModelContract
	ctiContract.CTIContract
	userContract.UserContract
	userPointContract.UserPointContract
}

// 初始化主合约
func (c *MainContract) InitLedger(ctx contractapi.TransactionContextInterface) (string, error) {
	// 可以在这里初始化一些初始数据
	// 创建一个默认用户
	defaultUser := typestruct.UserInfo{
		UserID:         "01",
		UserName:       "aision",
		PrivateKey:     "123456",
		PrivateKeyType: "RSA",
		Value:          100,                                   // 默认设置用户的 value
		CreateTime:     time.Now().UTC().Format(time.RFC3339), // 获取当前时间并格式化

	}

	// 将用户对象序列化为 JSON 字节数组
	userAsBytes, err := json.Marshal(defaultUser)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user: %v", err)
	}

	// 将用户数据存储到账本中
	err = ctx.GetStub().PutState(defaultUser.UserID, userAsBytes)
	if err != nil {
		return "", fmt.Errorf("failed to put user into world state: %v", err)
	}
	return "success", nil
}

// 注册模型信息
func (c *MainContract) RegisterModelInfo(ctx contractapi.TransactionContextInterface, modelname string, traffictype string, trafficfeatures []string, trafficprocesscode string, mlmethod string, mlinfo string, mltraincode string, ipfshash string, refctiid string, privateKey string) error {
	return c.ModelContract.RegisterModelInfo(ctx, modelname, traffictype, trafficfeatures, trafficprocesscode, mlmethod, mlinfo, mltraincode, ipfshash, refctiid, privateKey)
}

// 注册情报信息
func (c *MainContract) RegisterCTIInfo(ctx contractapi.TransactionContextInterface, ctiName string, ctitype int, ctiTrafficType int, openSource int, tags []string, iocs []string, stixdata string, statisticInfo string, description string, dataSize int, ipfsHash string, need int, value int, compreValue int, privateKey string) error {
	return c.CTIContract.RegisterCTIInfo(ctx, ctiName, ctitype, ctiTrafficType, openSource, tags, iocs, stixdata, statisticInfo, description, dataSize, ipfsHash, need, value, compreValue, privateKey)
}

// 注册用户信息
func (c *MainContract) RegisterUserInfo(ctx contractapi.TransactionContextInterface, userName, Privatekey string) error {
	return c.UserContract.RegisterUser(ctx, userName, Privatekey)
}

// 注册用户积分信息
func (c *MainContract) RegisterUserPointInfo(ctx contractapi.TransactionContextInterface, userID string) error {
	return c.UserPointContract.RegisterUserPointInfo(ctx, userID)
}

// 用户购买情报
func (c *MainContract) PurchaseCTI(ctx contractapi.TransactionContextInterface, ctiID string, userID string, txSignature string, nonceSignature string) error {
	return c.UserPointContract.PurchaseCTI(ctx, ctiID, userID, txSignature, nonceSignature)
}

// 查询模型信息
func (c *MainContract) QueryModelInfo(ctx contractapi.TransactionContextInterface, modelID string) (*typestruct.ModelInfo, error) {
	return c.ModelContract.QueryModelInfo(ctx, modelID)
}

// 查询情报信息
func (c *MainContract) QueryCTIInfo(ctx contractapi.TransactionContextInterface, ctiID string) (*typestruct.CtiInfo, error) {
	return c.CTIContract.QueryCTIInfo(ctx, ctiID)
}

// 查询用户上传的情报
func (c *MainContract) QueryCTIInfoByCreatorUserID(ctx contractapi.TransactionContextInterface, privateKey string) (*typestruct.CtiInfo, error) {
	return c.CTIContract.QueryCTIInfo(ctx, privateKey)
}

// 根据cti类别查询
func (c *MainContract) QueryCTIInfoByType(ctx contractapi.TransactionContextInterface, ctiType int) ([]typestruct.CtiInfo, error) {
	return c.CTIContract.QueryCTIInfoByType(ctx, ctiType)
}

// 查询用户信息
func (c *MainContract) QueryUserInfo(ctx contractapi.TransactionContextInterface, userID string) (*typestruct.UserInfo, error) {
	return c.UserContract.QueryUserInfo(ctx, userID)
}

// 查询用户积分信息
func (c *MainContract) QueryUserPointInfo(ctx contractapi.TransactionContextInterface, userID string) (*typestruct.UserPointInfo, error) {
	return c.UserPointContract.QueryUserPointInfo(ctx, userID)
}

// 模型信息分页查询
func (c *MainContract) QueryModelInfoByModelIDWithPagination(ctx contractapi.TransactionContextInterface, modelIDPrefix string, pageSize int, bookmark string) ([]byte, error) {
	return c.ModelContract.QueryModelInfoByModelIDWithPagination(ctx, modelIDPrefix, pageSize, bookmark)
}

// 根据流量类型查询模型信息
func (c *MainContract) QueryModelsByTrafficType(ctx contractapi.TransactionContextInterface, trafficType string) ([]typestruct.ModelInfo, error) {
	return c.ModelContract.QueryModelsByTrafficType(ctx, trafficType)
}

// 查询用户所上传的模型信息
func (c *MainContract) QueryModelsByPrivateKey(ctx contractapi.TransactionContextInterface, privateKey string) ([]typestruct.ModelInfo, error) {
	return c.ModelContract.QueryModelsByPrivateKey(ctx, privateKey)
}

// 根据相关CTI查询
func (c *MainContract) QueryModelsByRefCTIId(ctx contractapi.TransactionContextInterface, refCTIId string) ([]typestruct.ModelInfo, error) {
	return c.ModelContract.QueryModelsByRefCTIId(ctx, refCTIId)
}

// 更新用户信息
func (c *MainContract) UpdateUserInfo(ctx contractapi.TransactionContextInterface, PrivatecKey string, newUserName string, userid string) error {
	return c.UserContract.UpdateUserInfo(ctx, PrivatecKey, newUserName, userid)
}

// 分页查询
func (c *MainContract) QueryCTIInfoByCTIIDWithPagination(ctx contractapi.TransactionContextInterface, ctiIDPrefix string, pageSize int, bookmark string) ([]byte, error) {
	return c.CTIContract.QueryCTIInfoByCTIIDWithPagination(ctx, ctiIDPrefix, pageSize, bookmark)
}

// 生成 RSA 公私钥对
func (c *MainContract) GenerateRSAKeyPair(ctx contractapi.TransactionContextInterface) (string, error) {
	// 调用 utils 包中的 GenerateRSAKeyPair 方法
	keyPair, err := utils.GenerateRSAKeyPair()
	if err != nil {
		return "", fmt.Errorf("failed to generate RSA key pair: %v", err)
	}
	return keyPair, nil
}

// 从私钥获取公钥
func (c *MainContract) GetPublicKeyFromPrivateKey(ctx contractapi.TransactionContextInterface, privateKeyPEM string) (string, error) {
	// 调用 utils 包中的 GetPublicKeyFromPrivateKey 方法
	publicKey, err := utils.GetPublicKeyFromPrivateKey(privateKeyPEM)
	if err != nil {
		return "", fmt.Errorf("failed to get public key from private key: %v", err)
	}
	return publicKey, nil
}

// 生成交易随机数
func (c *MainContract) GetTransactionNonce(ctx contractapi.TransactionContextInterface, userID string, signatureContent string) (string, error) {
	// 生成一个随机的 32 字节数组
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	// 将随机字节转换为 base64 编码的字符串
	nonce := base64.StdEncoding.EncodeToString(randomBytes)
	//保存随机数和签名
	err = ctx.GetStub().PutState(nonce, []byte(signatureContent))
	if err != nil {
		return "", fmt.Errorf("failed to put state: %v", err)
	}

	return nonce, nil
}

// 主函数
func main() {
	chaincode, err := contractapi.NewChaincode(&MainContract{})
	if err != nil {
		fmt.Printf("Error creating main chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting main chaincode: %s", err.Error())
	}
}
