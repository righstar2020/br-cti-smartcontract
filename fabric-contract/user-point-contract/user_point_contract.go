package user_point_contract

import (
    "encoding/json"
    "fmt"
    
    "github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
)


// UserPointContract 是积分合约的结构体
type UserPointContract struct {
    contractapi.Contract
}



// QueryUserPointInfo 根据ID查询用户积分信息
func (c *UserPointContract) QueryUserPointInfo(ctx contractapi.TransactionContextInterface, userID string) (*typestruct.UserPointInfo, error) {
    // 从UserPointInfoMap中获取用户积分信息
    userPointInfoMapJSON, err := ctx.GetStub().GetState(userID+"_point_info")
    if err != nil {
        return nil, fmt.Errorf("从世界状态读取失败: %v", err)
    }
    if userPointInfoMapJSON == nil {
        return nil, fmt.Errorf("用户积分信息映射不存在")
    }

    var userPointInfo *typestruct.UserPointInfo
    err = json.Unmarshal(userPointInfoMapJSON, &userPointInfo)
    if err != nil {
        return nil, fmt.Errorf("解析用户积分信息映射失败: %v", err)
    }



    return userPointInfo, nil
}

// PurchaseCTI 用户使用积分购买情报
func (c *UserPointContract) PurchaseCTI(ctx contractapi.TransactionContextInterface, txData []byte) error {
    //解析msgData
    var purchaseCTITxData msgstruct.PurchaseCtiTxData
    err := json.Unmarshal(txData, &purchaseCTITxData)
    if err != nil {
        return fmt.Errorf("failed to unmarshal msg data: %v", err)
    }


    // 获取用户积分信息
    userPointInfo, err := c.QueryUserPointInfo(ctx, purchaseCTITxData.UserID)
    if err != nil {
        return err
    }

    // 获取情报信息
    ctiInfo, err := c.QueryCTIInfo(ctx, purchaseCTITxData.CTIID)
    if err != nil {
        return err
    }

    // 检查用户是否有足够的积分
    if userPointInfo.UserValue < ctiInfo.Value {
        return fmt.Errorf("insufficient points for purchase")
    }
    userID := purchaseCTITxData.UserID
    ctiID := purchaseCTITxData.CTIID
    // 更新用户积分信息
    userPointInfo.UserValue -= ctiInfo.Value
    userPointInfo.UserCTIMap[userID] = append(userPointInfo.UserCTIMap[userID], ctiID)
    userPointInfo.CTIBuyMap[ctiID] = ctiInfo.Value


    // 更新UserPointInfo
    userPointInfoMapJSONBytes, _ := json.Marshal(userPointInfo)
    err = ctx.GetStub().PutState(userID+"_point_info", userPointInfoMapJSONBytes)
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