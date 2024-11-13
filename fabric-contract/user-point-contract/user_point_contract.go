package user_point_contract

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
     CTIContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/cti-contract"
)

// UserPointInfo 结构体表示用户积分信息
type UserPointInfo struct {
    UserPointMap        map[string]int                `json:"user_point_map"`
    UserCTIMap          map[string][]string           `json:"user_cti_map"`
    CTICompreValueMap   map[string]int                `json:"cti_compre_value_map"`
    CTISaleMap          map[string]int                `json:"cti_sale_map"`
}

// UserPointContract 是积分合约的结构体
type UserPointContract struct {
    contractapi.Contract
}

// RegisterUserPointInfo 注册用户积分信息
func (c *UserPointContract) RegisterUserPointInfo(ctx contractapi.TransactionContextInterface, userInfoJSON string, userID string) error {
    var userPointInfo UserPointInfo
    err := json.Unmarshal([]byte(userInfoJSON), &userPointInfo)
    if err != nil {
        return fmt.Errorf("failed to unmarshal user point info: %v", err)
    }

    userPointInfoJSONBytes, _ := json.Marshal(userPointInfo)
    err = ctx.GetStub().PutState(userID, userPointInfoJSONBytes)
    if err != nil {
        return fmt.Errorf("failed to put state: %v", err)
    }

    return nil
}

// QueryUserInfo 根据ID查询用户积分信息
func (c *UserPointContract) QueryUserInfo(ctx contractapi.TransactionContextInterface, userID string) (*UserPointInfo, error) {
    userPointInfoJSON, err := ctx.GetStub().GetState(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if userPointInfoJSON == nil {
        return nil, fmt.Errorf("the user %s does not exist", userID)
    }

    var userPointInfo UserPointInfo
    err = json.Unmarshal(userPointInfoJSON, &userPointInfo)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal user point info: %v", err)
    }

    return &userPointInfo, nil
}

// PurchaseCTI 用户使用积分购买情报
func (c *UserPointContract) PurchaseCTI(ctx contractapi.TransactionContextInterface, ctiID string, userID string, txSignature string, nonceSignature string) error {
    // 获取用户积分信息
    userPointInfo, err := c.QueryUserInfo(ctx, userID)
    if err != nil {
        return err
    }

    // 获取情报信息
    ctiInfo, err := c.QueryCTIInfo(ctx, ctiID)
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

// QueryCTIInfo 根据ID查询情报信息
func (c *UserPointContract) QueryCTIInfo(ctx contractapi.TransactionContextInterface, ctiID string) (*CTIContract.CTIInfo, error) {
    ctiInfoJSON, err := ctx.GetStub().GetState(ctiID)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if ctiInfoJSON == nil {
        return nil, fmt.Errorf("the cti %s does not exist", ctiID)
    }

    var ctiInfo CTIContract.CTIInfo
    err = json.Unmarshal(ctiInfoJSON, &ctiInfo)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal cti info: %v", err)
    }

    return &ctiInfo, nil
}