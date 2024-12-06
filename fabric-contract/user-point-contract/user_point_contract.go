package user_point_contract

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"encoding/base64"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	ctiContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/cti-contract"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
	modelContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/model-contract"
)

// UserPointContract 是积分合约的结构体
type UserPointContract struct {
	ctiContract.CTIContract
	modelContract.ModelContract
}

//注册用户积分
func (c *UserPointContract) RegisterUserPointInfo(ctx contractapi.TransactionContextInterface, userID string, userPointInfo typestruct.UserPointInfo) error {
	userPointInfoBytes, err := json.Marshal(userPointInfo)
	if err != nil {
		return fmt.Errorf("序列化用户积分信息失败: %v", err)
	}
	return ctx.GetStub().PutState(userID+"_point_info", userPointInfoBytes)
}

// QueryUserPointInfo 根据ID查询用户积分信息
func (c *UserPointContract) QueryUserPointInfo(ctx contractapi.TransactionContextInterface, userID string) (*typestruct.UserPointInfo, error) {
	// 从UserPointInfoMap中获取用户积分信息
	userPointInfoMapJSON, err := ctx.GetStub().GetState(userID + "_point_info")
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
	fromPointInfo.UserCTIMap[ctiID] = points
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
func (c *UserPointContract) PurchaseCTI(ctx contractapi.TransactionContextInterface, purchaseCTITxData msgstruct.PurchaseCtiTxData, nonce string) (string,error) {
	

	// 获取用户积分信息
	userPointInfo, err := c.QueryUserPointInfo(ctx, purchaseCTITxData.UserID)
	if err != nil {
		return "",err
	}

	// 获取情报信息
	ctiInfo, err := c.QueryCTIInfo(ctx, purchaseCTITxData.CTIID)
	if err != nil {
		return "",err
	}

	// 检查用户是否有足够的积分
	if userPointInfo.UserValue < ctiInfo.Value {
		return "",fmt.Errorf("insufficient points for purchase")
	}
	userID := purchaseCTITxData.UserID
	ctiID := purchaseCTITxData.CTIID
	sellerID := ctiInfo.CreatorUserID

	// 处理积分转移
	err = c.TransferPoints(ctx, userID, sellerID, ctiInfo.Value, ctiID)
	if err != nil {
		return "",fmt.Errorf("积分转移失败: %v", err)
	}

	// 创建交易记录
	transaction_id,err := c.CreateBilateralTransactions(ctx, userID, sellerID, ctiInfo.Value, ctiID, nonce)
	if err != nil {
		return "",fmt.Errorf("创建交易记录失败: %v", err)
	}

	// 更新CTI交易总数
	err = c.UpdateCTITransactionCount(ctx)
	if err != nil {
		return transaction_id,fmt.Errorf("更新交易计数失败: %v", err)
	}

	return transaction_id,nil
}
// 购买模型
func (c *UserPointContract) PurchaseModel(ctx contractapi.TransactionContextInterface, purchaseModelTxData msgstruct.PurchaseModelTxData, nonce string) (string,error) {
	// 获取用户积分信息
	userPointInfo, err := c.QueryUserPointInfo(ctx, purchaseModelTxData.UserID)
	if err != nil {
		return "",err
	}

	// 获取模型信息
	modelInfo, err := c.QueryModelInfo(ctx, purchaseModelTxData.ModelID)
	if err != nil {
		return "",err
	}

	// 检查用户是否有足够的积分
	if userPointInfo.UserValue < modelInfo.Value {
		return "",fmt.Errorf("insufficient points for purchase")
	}

	userID := purchaseModelTxData.UserID
	modelID := purchaseModelTxData.ModelID
	sellerID := modelInfo.CreatorUserID

	// 处理积分转移
	err = c.TransferPoints(ctx, userID, sellerID, modelInfo.Value, modelID)
	if err != nil {
		return "",fmt.Errorf("积分转移失败: %v", err)
	}

	// 创建交易记录
	transaction_id,err := c.CreateBilateralTransactions(ctx, userID, sellerID, modelInfo.Value, modelID, nonce)
	if err != nil {
		return "",fmt.Errorf("创建交易记录失败: %v", err)
	}

	// 更新模型交易总数
	err = c.UpdateModelTransactionCount(ctx)
	if err != nil {
		return transaction_id,fmt.Errorf("更新交易计数失败: %v", err)
	}

	return transaction_id,nil
}


// UserStatistics 用户统计数据结构
type UserStatistics struct {
	TotalCTICount   int `json:"totalCTICount"`   // 链上总情报数
	UserCTICount    int `json:"userCTICount"`    // 我的情报数
	UserUploadCount int `json:"userUploadCount"` // 我的上链数
}

// 更新用户CTI的统计信息
func (c *UserPointContract) UpdateUserCTIStatistics(ctx contractapi.TransactionContextInterface, userID string, ctiCount int) error {
	// 获取现有统计数据
	totalCTICount := 0
	userCTICount := 0
	userUploadCount := 0

	totalCTICountBytes, err := ctx.GetStub().GetState("total_cti_count")
	if err != nil {
		return fmt.Errorf("获取总情报数失败: %v", err)
	}
	if totalCTICountBytes != nil {
		err = json.Unmarshal(totalCTICountBytes, &totalCTICount)
		if err != nil {
			return fmt.Errorf("解析总情报数失败: %v", err)
		}
	}
	var userStatistics UserStatistics
	// 使用用户ID作为key存储统计数据
	key := fmt.Sprintf("USER_CTI_STATS_%s", userID)
	statisticsBytes, err := ctx.GetStub().GetState(key)
	if err != nil || statisticsBytes == nil {
		userStatistics = UserStatistics{
			TotalCTICount:   totalCTICount,
			UserCTICount:    0,
			UserUploadCount: 0,
		}
	} else {
		err = json.Unmarshal(statisticsBytes, &userStatistics)
		if err != nil {
			return fmt.Errorf("解析用户统计数据失败: %v", err)
		}
	}

	// 更新统计数据
	userStatistics.TotalCTICount = totalCTICount + ctiCount
	userStatistics.UserCTICount = userCTICount + ctiCount
	userStatistics.UserUploadCount = userUploadCount + ctiCount
	// 将统计数据序列化为JSON
	statisticsBytes, err = json.Marshal(userStatistics)
	if err != nil {
		return fmt.Errorf("序列化用户统计数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(key, statisticsBytes)
	if err != nil {
		return fmt.Errorf("保存用户统计数据失败: %v", err)
	}
	err = ctx.GetStub().PutState("total_cti_count", []byte(strconv.Itoa(userStatistics.TotalCTICount)))
	if err != nil {
		return fmt.Errorf("保存总情报数失败: %v", err)
	}

	return nil
}

// QueryUserStatistics 查询用户统计数据
func (c *UserPointContract) QueryUserStatistics(ctx contractapi.TransactionContextInterface, userID string) (*UserStatistics, error) {
	// 获取用户积分信息
	userPointInfo, err := c.QueryUserPointInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户积分信息失败: %v", err)
	}

	// 获取用户拥有的情报数量
	userCTICount := len(userPointInfo.UserCTIMap)

	// 获取用户上传的情报数量
	key := fmt.Sprintf("USER_CTI_STATS_%s", userID)
	statisticsBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("获取用户统计数据失败: %v", err)
	}

	userUploadCount := 0
	if statisticsBytes != nil {
		var stats UserStatistics
		err = json.Unmarshal(statisticsBytes, &stats)
		if err != nil {
			return nil, fmt.Errorf("解析用户统计数据失败: %v", err)
		}
		userCTICount += stats.UserCTICount
		userUploadCount = stats.UserUploadCount
	}

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
		TotalCTICount:   totalCTICount,
		UserCTICount:    userCTICount,
		UserUploadCount: userUploadCount,
	}

	return statistics, nil
}

// PointTransaction 积分交易记录结构
type PointTransaction struct {
	TransactionID   string `json:"transaction_id"`   // 交易ID
	TransactionType string `json:"transaction_type"` // 交易类型：in/out
	Points          int    `json:"points"`           // 积分数量
	OtherParty      string `json:"other_party"`      // 对方账户
	InfoID          string `json:"info_id"`          // 相关情报ID
	Timestamp       string `json:"timestamp"`        // 交易时间
	Status          string `json:"status"`           // 交易状态(success/fail)
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
	fromID string, toID string, points int, infoID string,nonce string) (string,error) {

	timestamp := time.Now().Format("2006-01-02 15:04")
	// 从base64编码的nonce中提取随机数
	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	nonceNum := 100000

	if err == nil && len(nonceBytes) >= 3 {
		// 使用前3个字节生成6位随机数
		nonceNum = int(nonceBytes[0])*10000 + int(nonceBytes[1])*100 + int(nonceBytes[2])
		nonceNum = nonceNum % 1000000 // 确保是6位数
	}
	timesID := time.Now().Format("0601021504")
	randomNum := fmt.Sprintf("%06d", nonceNum)
	// 生成交易ID: 时间戳(12位,年月日时分) + 随机数(6位)
	transaction_id := fmt.Sprintf("%s%s", timesID, randomNum)
	// 支出方交易记录
	outTransaction := &PointTransaction{
		TransactionID:   transaction_id,
		TransactionType: "out",
		Points:          -points,
		OtherParty:      toID,
		InfoID:          infoID,
		Timestamp:       timestamp,
		Status:          "success",
	}
	if err := c.AddPointTransaction(ctx, fromID, outTransaction); err != nil {
		return "",fmt.Errorf("添加支出方交易记录失败: %v", err)
	}
	
	// 收入方交易记录
	inTransaction := &PointTransaction{
		TransactionID:   transaction_id,
		TransactionType: "in",
		Points:          points,
		OtherParty:      fromID,
		InfoID:          infoID,
		Timestamp:       timestamp,
		Status:          "success",
	}
	if err := c.AddPointTransaction(ctx, toID, inTransaction); err != nil {
		return "",fmt.Errorf("添加收入方交易记录失败: %v", err)
	}

	return transaction_id,nil
}

// UpdateCTITransactionCount 更新CTI交易总数
func (c *UserPointContract) UpdateCTITransactionCount(ctx contractapi.TransactionContextInterface) error {
	// 获取当前交易计数
	key := "total_cti_transactions"
	countBytes, err := ctx.GetStub().GetState(key)

	currentCount := 0
	if err != nil {
		return fmt.Errorf("获取交易总数失败: %v", err)
	}
	if countBytes != nil {
		currentCount, err = strconv.Atoi(string(countBytes))
		if err != nil {
			return fmt.Errorf("解析交易总数失败: %v", err)
		}
	}

	// 增加计数并保存
	newCount := currentCount + 1
	err = ctx.GetStub().PutState(key, []byte(strconv.Itoa(newCount)))
	if err != nil {
		return fmt.Errorf("保存更新后的交易总数失败: %v", err)
	}

	return nil
}
// UpdateModelTransactionCount 更新模型交易总数
func (c *UserPointContract) UpdateModelTransactionCount(ctx contractapi.TransactionContextInterface) error {
	// 获取当前交易计数
	key := "total_model_transactions"
	countBytes, err := ctx.GetStub().GetState(key)

	currentCount := 0
	if err != nil {
		return fmt.Errorf("获取交易总数失败: %v", err)
	}
	if countBytes != nil {
		currentCount, err = strconv.Atoi(string(countBytes))
		if err != nil {
			return fmt.Errorf("解析交易总数失败: %v", err)
		}
	}

	// 增加计数并保存
	newCount := currentCount + 1
	err = ctx.GetStub().PutState(key, []byte(strconv.Itoa(newCount)))
	if err != nil {
		return fmt.Errorf("保存更新后的交易总数失败: %v", err)
	}

	return nil
}

// GetCTITransactionCount 获取当前CTI交易总数
func (c *UserPointContract) GetCTITransactionCount(ctx contractapi.TransactionContextInterface) (int, error) {
	key := "total_cti_transactions"
	countBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return 0, fmt.Errorf("获取交易总数失败: %v", err)
	}

	currentCount := 0
	if countBytes != nil {
		currentCount, err = strconv.Atoi(string(countBytes))
		if err != nil {
			return 0, fmt.Errorf("解析交易总数失败: %v", err)
		}
	}

	return currentCount, nil
}

// QueryUserPurchasedCTIs 查询用户已购买的CTI信息列表
func (c *UserPointContract) QueryUserPurchasedCTIs(ctx contractapi.TransactionContextInterface, userID string) ([]typestruct.CtiInfo, error) {
	// 获取用户积分信息
	userPointInfo, err := c.QueryUserPointInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户积分信息失败: %v", err)
	}

	// 获取用户购买的CTI ID列表
	ctiIDs := []string{}
	for ctiID := range userPointInfo.CTIBuyMap {
		ctiIDs = append(ctiIDs, ctiID)
	}
	if len(ctiIDs) == 0 {
		return []typestruct.CtiInfo{}, nil
	}

	// 查询每个CTI的详细信息
	var ctiInfoList []typestruct.CtiInfo
	for _, ctiID := range ctiIDs {
		ctiInfo, err := c.QueryCTIInfo(ctx, ctiID)
		if err == nil {
			ctiInfoList = append(ctiInfoList, *ctiInfo)
		}
	}

	return ctiInfoList, nil
}

// UserModelStatistics 用户模型统计数据结构
type UserModelStatistics struct {
	TotalModelCount   int `json:"totalModelCount"`   // 链上总模型数
	UserModelCount    int `json:"userModelCount"`    // 我的模型数
	UserUploadCount   int `json:"userUploadCount"`   // 我的上链数
}

// UpdateUserModelStatistics 更新用户模型的统计信息
func (c *UserPointContract) UpdateUserModelStatistics(ctx contractapi.TransactionContextInterface, userID string, modelCount int) error {
	// 获取现有统计数据
	totalModelCount := 0
	
	totalModelCountBytes, err := ctx.GetStub().GetState("total_model_count")
	if err != nil {
		return fmt.Errorf("获取总模型数失败: %v", err)
	}
	if totalModelCountBytes != nil {
		err = json.Unmarshal(totalModelCountBytes, &totalModelCount)
		if err != nil {
			return fmt.Errorf("解析总模型数失败: %v", err)
		}
	}

	var userStatistics UserModelStatistics
	// 使用用户ID作为key存储统计数据
	key := fmt.Sprintf("USER_MODEL_STATS_%s", userID)
	statisticsBytes, err := ctx.GetStub().GetState(key)
	if err != nil || statisticsBytes == nil {
		userStatistics = UserModelStatistics{
			TotalModelCount: totalModelCount,
			UserModelCount:  0,
			UserUploadCount: 0,
		}
	} else {
		err = json.Unmarshal(statisticsBytes, &userStatistics)
		if err != nil {
			return fmt.Errorf("解析用户统计数据失败: %v", err)
		}
	}

	// 更新统计数据
	userStatistics.TotalModelCount = totalModelCount + modelCount
	userStatistics.UserModelCount += modelCount
	userStatistics.UserUploadCount += modelCount

	// 将统计数据序列化为JSON
	statisticsBytes, err = json.Marshal(userStatistics)
	if err != nil {
		return fmt.Errorf("序列化用户统计数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(key, statisticsBytes)
	if err != nil {
		return fmt.Errorf("保存用户统计数据失败: %v", err)
	}
	err = ctx.GetStub().PutState("total_model_count", []byte(strconv.Itoa(userStatistics.TotalModelCount)))
	if err != nil {
		return fmt.Errorf("保存总模型数失败: %v", err)
	}

	return nil
}
