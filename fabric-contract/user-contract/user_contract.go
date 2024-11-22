package user_contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
)

// UserInfo 结构体表示用户信息

// UserContract 是用户信息合约的结构体
type UserContract struct {
	contractapi.Contract
}

// 注册用户(msgData:入参数据需要解析)
func (c *UserContract) RegisterUser(ctx contractapi.TransactionContextInterface, msgData []byte)  error {

	//解析消息结构体
	var userMsgData msgstruct.UserRegisterMsgData
	err := json.Unmarshal(msgData, &userMsgData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal msg data: %v", err)
	}
	
	// 验证必需字段
	if userMsgData.PublicKey == "" {
		return fmt.Errorf("public key is required")
	}
	
	// 使用公钥的 SHA256 哈希值作为 UserID
	hash := sha256.New()
	hash.Write([]byte(userMsgData.PublicKey))
	userID := hex.EncodeToString(hash.Sum(nil))

	// 检查生成的 UserID 是否已存在
	userAsBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return  fmt.Errorf("failed to read from world state: %v", err)
	}
	if userAsBytes != nil {
		return  fmt.Errorf("the user with ID %s already exists", userID)
	}

	// 获取交易的时间作为创建时间
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return  fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	createTime := time.Unix(txTimestamp.Seconds, int64(txTimestamp.Nanos)).UTC().Format(time.RFC3339)

	// 创建新的用户对象，初始 value 设置为 0
	newUser := typestruct.UserInfo{
		UserID:         userID,
		PublicKey:      userMsgData.PublicKey,
		PublicKeyType:  "ECC",
		Value:          0,
		CreateTime:     createTime,
	}
	// 将新用户对象序列化为 JSON 字节数组
	newUserAsBytes, err := json.Marshal(newUser)
	if err != nil {
		return  fmt.Errorf("failed to marshal user: %v", err)
	}
	// 将新用户数据存储到账本中
	err = ctx.GetStub().PutState(userID, newUserAsBytes)
	if err != nil {
		return  fmt.Errorf("failed to put user into world state: %v", err)
	}
	// // 初始化 UserPointMap
	newUserPointInfo := typestruct.UserPointInfo{
		UserValue:   0, // 初始化用户的积分值为 0
		UserCTIMap: make(map[string][]string),    // 空的CTI映射
		CTIBuyMap: 	make(map[string]int),      	// 空的CTI购买映射
		CTISaleMap: make(map[string]int),      // 空的CTI销售映射
	}
	
	// 将用户的积分信息序列化为 JSON 字节数组
	newPointInfoAsBytes, err := json.Marshal(newUserPointInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal point info: %v", err)
	}

	// 将用户的积分映射存储到账本中
	err = ctx.GetStub().PutState(userID+"_point_info", newPointInfoAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put user points into world state: %v", err)
	}
	return  nil
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


