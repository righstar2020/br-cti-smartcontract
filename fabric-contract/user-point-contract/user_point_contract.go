package user_point_contract

import (
    "encoding/json"
    "fmt"
    
    "github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
     //UserContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/user-contract"
)

// UserPointInfo 结构体表示用户积分信息,映射到情报ID:ctiID
// type UserPointInfo struct {
//     UserValue           int                           `json:"user_value"` //用户积分    
//     UserCTIMap          map[string][]string           `json:"user_cti_map"` //用户拥有的情报map
//     CTIBuyMap           map[string]int                `json:"cti_buy_map"` //用户购买的情报map
//     CTISaleMap          map[string]int                `json:"cti_sale_map"` //用户销售的情报map
// }

// UserPointContract 是积分合约的结构体
type UserPointContract struct {
    contractapi.Contract
}
func (c *UserPointContract) InitUserPointInfoContract(ctx contractapi.TransactionContextInterface) error {
    //初始化积分合约
    userPointInfoMap := make(map[string]typestruct.UserPointInfo)
    userPointInfoMapJSONBytes, _ := json.Marshal(userPointInfoMap)
    err := ctx.GetStub().PutState("UserPointInfoMap", userPointInfoMapJSONBytes)
    if err != nil {
        return fmt.Errorf("failed to put state: %v", err)
    }
    return nil
    
}


// RegisterUserPointInfo 注册用户积分信息(初始化)
func (c *UserPointContract) RegisterUserPointInfo(ctx contractapi.TransactionContextInterface,userID string) error {
    var userInfo *typestruct.UserInfo
    userInfoBytes, err := ctx.GetStub().GetState(userID)
    err = json.Unmarshal(userInfoBytes, &userInfo)
    if err != nil {
        return fmt.Errorf("failed to unmarshal user info: %v", err)
    }
    var userPointInfo = typestruct.UserPointInfo{
        UserValue: userInfo.Value,
        UserCTIMap: make(map[string][]string),
        CTISaleMap: make(map[string]int),
    }
    var userPointInfoMap = make(map[string]typestruct.UserPointInfo)
    userPointInfoMapJSONBytes, _ := ctx.GetStub().GetState("UserPointInfoMap")
    err = json.Unmarshal(userPointInfoMapJSONBytes, &userPointInfoMap)
    if err != nil {
        return fmt.Errorf("failed to unmarshal user point info map: %v", err)
    }
    //更新用户积分信息
    userPointInfoMap[userID] = userPointInfo
    userPointInfoMapJSONBytes, _ = json.Marshal(userPointInfoMap)
    err = ctx.GetStub().PutState("UserPointInfoMap", userPointInfoMapJSONBytes)
    if err != nil {
        return fmt.Errorf("failed to put state: %v", err)
    }

    return nil
}

// QueryUserPointInfo 根据ID查询用户积分信息
func (c *UserPointContract) QueryUserPointInfo(ctx contractapi.TransactionContextInterface, userID string) (*typestruct.UserPointInfo, error) {
    // 从UserPointInfoMap中获取用户积分信息
    userPointInfoMapJSON, err := ctx.GetStub().GetState("UserPointInfoMap")
    if err != nil {
        return nil, fmt.Errorf("从世界状态读取失败: %v", err)
    }
    if userPointInfoMapJSON == nil {
        return nil, fmt.Errorf("用户积分信息映射不存在")
    }

    var userPointInfoMap map[string]typestruct.UserPointInfo
    err = json.Unmarshal(userPointInfoMapJSON, &userPointInfoMap)
    if err != nil {
        return nil, fmt.Errorf("解析用户积分信息映射失败: %v", err)
    }

    // 获取指定用户的积分信息
    userPointInfo, exists := userPointInfoMap[userID]
    if !exists {
        return nil, fmt.Errorf("用户 %s 不存在", userID)
    }

    return &userPointInfo, nil
}

// PurchaseCTI 用户使用积分购买情报
func (c *UserPointContract) PurchaseCTI(ctx contractapi.TransactionContextInterface, ctiID string, userID string, txSignature string, nonceSignature string) error {
    // 获取用户积分信息
    userPointInfo, err := c.QueryUserPointInfo(ctx, userID)
    if err != nil {
        return err
    }

    // 获取情报信息
    ctiInfo, err := c.QueryCTIInfo(ctx, ctiID)
    if err != nil {
        return err
    }

    // 检查用户是否有足够的积分
    if userPointInfo.UserValue < ctiInfo.Value {
        return fmt.Errorf("insufficient points for purchase")
    }

    // 更新用户积分信息
    userPointInfo.UserValue -= ctiInfo.Value
    userPointInfo.UserCTIMap[userID] = append(userPointInfo.UserCTIMap[userID], ctiID)
    userPointInfo.CTIBuyMap[ctiID] = ctiInfo.Value

    // 获取UserPointInfoMap
    var userPointInfoMap map[string]typestruct.UserPointInfo
    userPointInfoMapJSON, err := ctx.GetStub().GetState("UserPointInfoMap")
    if err != nil {
        return fmt.Errorf("failed to get UserPointInfoMap: %v", err)
    }
    err = json.Unmarshal(userPointInfoMapJSON, &userPointInfoMap)
    if err != nil {
        return fmt.Errorf("failed to unmarshal UserPointInfoMap: %v", err)
    }

    // 更新UserPointInfoMap
    userPointInfoMap[userID] = *userPointInfo
    userPointInfoMapJSONBytes, _ := json.Marshal(userPointInfoMap)
    err = ctx.GetStub().PutState("UserPointInfoMap", userPointInfoMapJSONBytes)
    if err != nil {
        return fmt.Errorf("failed to put state: %v", err)
    }

    return nil
}

// QueryCTIInfo 根据ID查询情报信息
func (c *UserPointContract) QueryCTIInfo(ctx contractapi.TransactionContextInterface, ctiID string) (*typestruct.CtiInfo, error) {
    ctiInfoJSON, err := ctx.GetStub().GetState(ctiID)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if ctiInfoJSON == nil {
        return nil, fmt.Errorf("the cti %s does not exist", ctiID)
    }

    var ctiInfo typestruct.CtiInfo
    err = json.Unmarshal(ctiInfoJSON, &ctiInfo)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal cti info: %v", err)
    }

    return &ctiInfo, nil
}