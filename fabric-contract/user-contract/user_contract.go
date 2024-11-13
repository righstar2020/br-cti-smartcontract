package user_contract

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// UserInfo 结构体表示用户信息
type UserInfo struct {
    UserID           string `json:"user_id"`
    UserName         string `json:"user_name"`
    PublicKey        string `json:"public_key"`
    PublicKeyType    string `json:"public_key_type"`
    Value            int    `json:"value"`
    CreateTime       string `json:"create_time"`
}

// UserContract 是用户信息合约的结构体
type UserContract struct {
    contractapi.Contract
}

// RegisterUserInfo 注册用户信息
func (c *UserContract) RegisterUserInfo(ctx contractapi.TransactionContextInterface, userInfoJSON string, userID string) error {
    var userInfo UserInfo
    err := json.Unmarshal([]byte(userInfoJSON), &userInfo)
    if err != nil {
        return fmt.Errorf("failed to unmarshal user info: %v", err)
    }

    userInfoJSONBytes, _ := json.Marshal(userInfo)
    err = ctx.GetStub().PutState(userInfo.UserID, userInfoJSONBytes)
    if err != nil {
        return fmt.Errorf("failed to put state: %v", err)
    }

    return nil
}

// QueryUserInfo 根据ID查询用户信息
func (c *UserContract) QueryUserInfo(ctx contractapi.TransactionContextInterface, userID string) (*UserInfo, error) {
    userInfoJSON, err := ctx.GetStub().GetState(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if userInfoJSON == nil {
        return nil, fmt.Errorf("the user %s does not exist", userID)
    }

    var userInfo UserInfo
    err = json.Unmarshal(userInfoJSON, &userInfo)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal user info: %v", err)
    }

    return &userInfo, nil
}

// UpdateUserInfo 更新用户信息
func (c *UserContract) UpdateUserInfo(ctx contractapi.TransactionContextInterface, userInfoJSON string, userID string, txSignature string, nonceSignature string) error {
    var userInfo UserInfo
    err := json.Unmarshal([]byte(userInfoJSON), &userInfo)
    if err != nil {
        return fmt.Errorf("failed to unmarshal user info: %v", err)
    }

    // TODO: 验证签名和随机数

    userInfoJSONBytes, _ := json.Marshal(userInfo)
    err = ctx.GetStub().PutState(userInfo.UserID, userInfoJSONBytes)
    if err != nil {
        return fmt.Errorf("failed to put state: %v", err)
    }

    return nil
}