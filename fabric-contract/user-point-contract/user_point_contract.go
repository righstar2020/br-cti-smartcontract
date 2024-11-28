package user_point_contract

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
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

// TransferPoints 处理积分转移和相关状态更新
func (c *UserPointContract) TransferPoints(ctx contractapi.TransactionContextInterface, 
    fromID string, toID string, points int, ctiID string) error {
    
    // 获取买方积分信息
    fromPointInfo, err := c.QueryUserPointInfo(ctx, fromID)
    if err != nil {
        return fmt.Errorf("获取买方积分信息失败: %v", err)
    }

    // 获取卖方积分信息
    toPointInfo, err := c.QueryUserPointInfo(ctx, toID)
    if err != nil {
        return fmt.Errorf("获取卖方积分信息失败: %v", err)
    }

    // 更新买方积分信息
    fromPointInfo.UserValue -= points
    fromPointInfo.UserCTIMap[fromID] = append(fromPointInfo.UserCTIMap[fromID], ctiID)
    fromPointInfo.CTIBuyMap[ctiID] = points

    // 更新卖方积分信息
    toPointInfo.UserValue += points
    toPointInfo.CTISaleMap[ctiID] = points

    // 更新买方UserPointInfo
    fromPointInfoBytes, err := json.Marshal(fromPointInfo)
    if err != nil {
        return fmt.Errorf("序列化买方积分信息失败: %v", err)
    }
    err = ctx.GetStub().PutState(fromID+"_point_info", fromPointInfoBytes)
    if err != nil {
        return fmt.Errorf("更新买方积分信息失败: %v", err)
    }

    // 更新卖方UserPointInfo
    toPointInfoBytes, err := json.Marshal(toPointInfo)
    if err != nil {
        return fmt.Errorf("序列化卖方积分信息失败: %v", err)
    }
    err = ctx.GetStub().PutState(toID+"_point_info", toPointInfoBytes)
    if err != nil {
        return fmt.Errorf("更新卖方积分信息失败: %v", err)
    }

    return nil
}

// PurchaseCTI 修改后的购买CTI函数
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
    sellerID := ctiInfo.CreatorUserID

    // 处理积分转移
    err = c.TransferPoints(ctx, userID, sellerID, ctiInfo.Value, ctiID)
    if err != nil {
        return fmt.Errorf("积分转移失败: %v", err)
    }

    // 创建交易记录
    err = c.CreateBilateralTransactions(ctx, userID, sellerID, ctiInfo.Value, ctiID)
    if err != nil {
        return fmt.Errorf("创建交易记录失败: %v", err)
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

// UserStatistics 用户统计数据结构
type UserStatistics struct {
    TotalCTICount    int `json:"totalCTICount"`    // 链上总情报数
    UserCTICount     int `json:"userCTICount"`     // 我的情报数
    UserUploadCount  int `json:"userUploadCount"`  // 我的上链数
}

// QueryUserStatistics 查询用户统计数据
func (c *UserPointContract) QueryUserStatistics(ctx contractapi.TransactionContextInterface, userID string) (*UserStatistics, error) {
    // 获取用户积分信息
    userPointInfo, err := c.QueryUserPointInfo(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("获取用户积分信息失败: %v", err)
    }

    // 获取用户拥有的情报数量
    userCTICount := len(userPointInfo.UserCTIMap[userID])

    // 获取用户上传的情报数量
    userUploadCount := 0 //从上层接口获取

    // 获取链上总情报数量
    // 注：这里需要实现一个计数器或使用其他方式来追踪总情报数
    totalCTICountBytes, err := ctx.GetStub().GetState("total_cti_count")
    if err != nil {
        return nil, fmt.Errorf("获取总情报数失败: %v", err)
    }
    totalCTICount := 0
    if totalCTICountBytes != nil {
        err = json.Unmarshal(totalCTICountBytes, &totalCTICount)
        if err != nil {
            return nil, fmt.Errorf("解析总情报数失败: %v", err)
        }
    }

    statistics := &UserStatistics{
        TotalCTICount:    totalCTICount,
        UserCTICount:     userCTICount,
        UserUploadCount:  userUploadCount,
    }

    return statistics, nil
}

// PointTransaction 积分交易记录结构
type PointTransaction struct {
    TransactionType string    `json:"transactionType"`  // 交易类型：转入/转出
    Points         int       `json:"points"`           // 积分数量
    OtherParty    string    `json:"otherParty"`       // 对方账户
    InfoID        string    `json:"infoId"`           // 相关情报ID
    Timestamp     string    `json:"timestamp"`        // 交易时间
    Status        string    `json:"status"`           // 交易状态
}

// QueryPointTransactions 查询用户的积分交易记录
func (c *UserPointContract) QueryPointTransactions(ctx contractapi.TransactionContextInterface, userID string) ([]*PointTransaction, error) {
    // 从世界状态获取用户的交易记录
    transactionsJSON, err := ctx.GetStub().GetState(userID + "_transactions")
    if err != nil {
        return nil, fmt.Errorf("获取交易记录失败: %v", err)
    }

    var transactions []*PointTransaction
    if transactionsJSON != nil {
        err = json.Unmarshal(transactionsJSON, &transactions)
        if err != nil {
            return nil, fmt.Errorf("解析交易记录失败: %v", err)
        }
    }

    return transactions, nil
}

// AddPointTransaction 添加积分交易记录
func (c *UserPointContract) AddPointTransaction(ctx contractapi.TransactionContextInterface, userID string, transaction *PointTransaction) error {
    // 获取现有交易记录
    transactions, err := c.QueryPointTransactions(ctx, userID)
    if err != nil {
        return err
    }

    // 添加新交易记录
    transactions = append(transactions, transaction)

    // 保存更新后的交易记录
    transactionsJSON, err := json.Marshal(transactions)
    if err != nil {
        return fmt.Errorf("序列化交易记录失败: %v", err)
    }

    err = ctx.GetStub().PutState(userID+"_transactions", transactionsJSON)
    if err != nil {
        return fmt.Errorf("保存交易记录失败: %v", err)
    }

    return nil
}

// CreateBilateralTransactions 创建双方交易记录
func (c *UserPointContract) CreateBilateralTransactions(ctx contractapi.TransactionContextInterface, 
    fromID string, toID string, points int, infoID string) error {
    
    timestamp := time.Now().Format("2006-01-02 15:04")
    
    // 支出方交易记录
    outTransaction := &PointTransaction{
        TransactionType: "转出",
        Points:         -points,
        OtherParty:     toID,
        InfoID:         infoID,
        Timestamp:      timestamp,
        Status:         "已完成",
    }
    err := c.AddPointTransaction(ctx, fromID, outTransaction)
    if err != nil {
        return fmt.Errorf("添加支出方交易记录失败: %v", err)
    }

    // 收入方交易记录
    inTransaction := &PointTransaction{
        TransactionType: "转入",
        Points:         points,
        OtherParty:     fromID,
        InfoID:         infoID,
        Timestamp:      timestamp,
        Status:         "已完成",
    }
    err = c.AddPointTransaction(ctx, toID, inTransaction)
    if err != nil {
        return fmt.Errorf("添加收入方交易记录失败: %v", err)
    }

    return nil
}

