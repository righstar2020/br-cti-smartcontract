package user_contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	userPointContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/user-point-contract"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
)

// UserInfo 结构体表示用户信息

// UserContract 是用户信息合约的结构体
type UserContract struct {
	userPointContract.UserPointContract
}

// 注册用户(msgData:入参数据需要解析)
func (c *UserContract) RegisterUser(ctx contractapi.TransactionContextInterface, userRigisterData msgstruct.UserRegisterMsgData) (string, error) {
	
	// 验证必需字段
	if userRigisterData.PublicKey == "" {
		return "", fmt.Errorf("public key is required")
	}
	
	// 使用公钥的 SHA256 哈希值作为 UserID
	hash := sha256.New()
	hash.Write([]byte(userRigisterData.PublicKey))
	userID := hex.EncodeToString(hash.Sum(nil))

	// 检查生成的 UserID 是否已存在
	userAsBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if userAsBytes != nil {
		return "", fmt.Errorf("the user with ID %s already exists", userID)
	}

	// 获取交易的时间作为创建时间
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	createTime := time.Unix(txTimestamp.Seconds, int64(txTimestamp.Nanos)).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")

	// 创建新的用户对象，初始 value 设置为 0
	newUser := typestruct.UserInfo{
		UserID:         userID,
		UserName:       userRigisterData.UserName,
		PublicKey:      userRigisterData.PublicKey,
		PublicKeyType:  "ECC",
		CreateTime:     createTime,
	}
	// 将新用户对象序列化为 JSON 字节数组
	newUserAsBytes, err := json.Marshal(newUser)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user: %v", err)
	}
	// 将新用户数据存储到账本中
	err = ctx.GetStub().PutState(userID, newUserAsBytes)
	if err != nil {
		return "", fmt.Errorf("failed to put user into world state: %v", err)
	}
	// // 初始化 UserPointMap
	newUserPointInfo := typestruct.UserPointInfo{
		UserValue:   100, // 初始化用户的积分值为 100
		UserLevel: 1, // 初始化用户的等级为 1(初级用户)
		UserCTIMap: make(map[string]float64),    // 空的CTI映射
		CTIBuyMap: 	make(map[string]float64),      	// 空的CTI购买映射
		CTISaleMap: make(map[string]float64),      // 空的CTI销售映射
		UserModelMap: make(map[string]float64),      // 空的模型映射
		ModelBuyMap: make(map[string]float64),    // 空的模型购买映射
		ModelSaleMap: make(map[string]float64),   // 空的模型销售映射
	}
	
	

	if err := c.UserPointContract.RegisterUserPointInfo(ctx, userID, newUserPointInfo); err != nil {
		return "", fmt.Errorf("failed to register user point info: %v", err)
	}

	// 将新用户ID添加到用户列表中
	err = c.addUserToAccountList(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to add user to list: %v", err)
	}
	return userID, nil
}

// 新增：将用户ID添加到用户列表中
func (c *UserContract) addUserToAccountList(ctx contractapi.TransactionContextInterface, userID string) error {
	// 使用固定的键来存储用户列表
	const userListKey = "USER_ACCOUNT_LIST_KEY"
	
	// 获取现有的用户列表
	userListBytes, err := ctx.GetStub().GetState(userListKey)
	var userList []string
	
	if err != nil {
		return fmt.Errorf("failed to read user list: %v", err)
	}
	
	// 如果列表存在，则解析它
	if userListBytes != nil {
		err = json.Unmarshal(userListBytes, &userList)
		if err != nil {
			return fmt.Errorf("failed to unmarshal user list: %v", err)
		}
	}
	
	// 添加新用户ID到列表
	userList = append(userList, userID)
	
	// 将更新后的列表保存回账本
	updatedListBytes, err := json.Marshal(userList)
	if err != nil {
		return fmt.Errorf("failed to marshal user list: %v", err)
	}
	
	err = ctx.GetStub().PutState(userListKey, updatedListBytes)
	if err != nil {
		return fmt.Errorf("failed to save user list: %v", err)
	}
	
	return nil
}

// 新增：查询所有注册用户的列表
func (c *UserContract) QueryUserAccountList(ctx contractapi.TransactionContextInterface) ([]string, error) {
	const userListKey = "USER_ACCOUNT_LIST_KEY"
	
	// 获取用户列表
	userListBytes, err := ctx.GetStub().GetState(userListKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read user list: %v", err)
	}
	
	// 如果列表不存在，返回空列表
	if userListBytes == nil {
		return []string{}, nil
	}
	
	// 解析用户列表
	var userList []string
	err = json.Unmarshal(userListBytes, &userList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user list: %v", err)
	}
	
	return userList, nil
}

// QueryUserInfo 根据ID查询用户信息
func (c *UserContract) QueryUserInfo(ctx contractapi.TransactionContextInterface, userID string) (*typestruct.UserInfo, error) {
	userInfoJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userInfoJSON == nil {
		return nil, fmt.Errorf("the user %s does not exist", userID)
	}

	var userInfo typestruct.UserInfo
	err = json.Unmarshal(userInfoJSON, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %v", err)
	}

	return &userInfo, nil
}

