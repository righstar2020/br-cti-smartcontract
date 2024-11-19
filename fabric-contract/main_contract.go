package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	modelContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/model-contract"
	ctiContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/cti-contract"
	userContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/user-contract"
	userPointContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/user-point-contract"
	dataContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/data-contract"
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
func (c *MainContract) InitLedger(ctx contractapi.TransactionContextInterface) string {
    // 可以在这里初始化一些初始数据
    return "success"
}

// 注册模型信息
func (c *MainContract) RegisterModelInfo(ctx contractapi.TransactionContextInterface, modelInfoJSON string, userID string, txSignature string, nonceSignature string) error {
    return c.ModelContract.RegisterModelInfo(ctx, modelInfoJSON, userID, txSignature, nonceSignature)
}

// 注册情报信息
func (c *MainContract) RegisterCTIInfo(ctx contractapi.TransactionContextInterface, ctiInfoJSON string, userID string, txSignature string, nonceSignature string) error {
    return c.CTIContract.RegisterCTIInfo(ctx, ctiInfoJSON, userID, txSignature, nonceSignature)
}

// 注册用户信息
func (c *MainContract) RegisterUserInfo(ctx contractapi.TransactionContextInterface, userInfoJSON string, userID string) error {
    return c.UserContract.RegisterUserInfo(ctx, userInfoJSON, userID)
}

// 注册用户积分信息
func (c *MainContract) RegisterUserPointInfo(ctx contractapi.TransactionContextInterface, userInfoJSON string, userID string) error {
    return c.UserPointContract.RegisterUserPointInfo(ctx, userInfoJSON, userID)
}

// 用户购买情报
func (c *MainContract) PurchaseCTI(ctx contractapi.TransactionContextInterface, ctiID string, userID string, txSignature string, nonceSignature string) error {
    // 获取用户积分信息
    userPointInfo, err := c.UserPointContract.QueryUserInfo(ctx, userID)
    if err != nil {
        return err
    }

    // 获取情报信息
    ctiInfo, err := c.CTIContract.QueryCTIInfo(ctx, ctiID)
    if err != nil {
        return err
    }

    // 检查用户是否有足够的积分
    if userPointInfo.UserPointMap[userID] < ctiInfo.Value {
        return fmt.Errorf("insufficient points for purchase")
    }

    // 更新用户积分信息
    userPointInfo.UserPointMap[userID] -= ctiInfo.Value
    userPointInfo.UserCTIMap[userID] = append(userPointInfo.UserCTIMap[userID], ctiID)
    userPointInfoJSONBytes, _ := json.Marshal(userPointInfo)
    err = ctx.GetStub().PutState(userID, userPointInfoJSONBytes)
    if err != nil {
        return fmt.Errorf("failed to put state: %v", err)
    }

    return nil
}

// 查询模型信息
func (c *MainContract) QueryModelInfo(ctx contractapi.TransactionContextInterface, modelID string) (*modelContract.ModelInfo, error) {
    return c.ModelContract.QueryModelInfo(ctx, modelID)
}

// 查询情报信息
func (c *MainContract) QueryCTIInfo(ctx contractapi.TransactionContextInterface, ctiID string) (*ctiContract.CTIInfo, error) {
    return c.CTIContract.QueryCTIInfo(ctx, ctiID)
}

// 查询用户信息
func (c *MainContract) QueryUserInfo(ctx contractapi.TransactionContextInterface, userID string) (*userContract.UserInfo, error) {
    return c.UserContract.QueryUserInfo(ctx, userID)
}

// 查询用户积分信息
func (c *MainContract) QueryUserPointInfo(ctx contractapi.TransactionContextInterface, userID string) (*userPointContract.UserPointInfo, error) {
    return c.UserPointContract.QueryUserInfo(ctx, userID)
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
        return "",fmt.Errorf("failed to put state: %v", err)
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


