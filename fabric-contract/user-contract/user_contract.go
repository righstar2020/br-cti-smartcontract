package user_contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
)

// UserInfo 结构体表示用户信息

// UserContract 是用户信息合约的结构体
type UserContract struct {
	contractapi.Contract
}

// 注册用户
func (c *UserContract) RegisterUser(ctx contractapi.TransactionContextInterface, userName, PrivatecKey string)  error {
	// 使用公钥的 SHA256 哈希值作为 UserID
	hash := sha256.New()
	hash.Write([]byte(PrivatecKey))
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
		UserName:       userName,
		PrivateKey:    PrivatecKey,
		PrivateKeyType: "RSA",
		Value:          0,
		CreateTime:     createTime,
	}

	// // 初始化 UserPointMap
	// newPointInfo := PointInfo{
	// 	UserPointMap: map[string]int{
	// 		userID: 0, // 初始化用户的积分值为 0
	// 	},
	// 	UserCtiMap:        make(map[string][]string), // 空的CTI映射
	// 	CtiCompreValueMap: make(map[string]int),      // 空的CTI综合价值映射
	// 	CtiSaleMap:        make(map[string]int),      // 空的CTI购买数量映射
	// }
	// 将新用户对象序列化为 JSON 字节数组
	newUserAsBytes, err := json.Marshal(newUser)
	if err != nil {
		return  fmt.Errorf("failed to marshal user: %v", err)
	}

	// // 将用户的积分信息序列化为 JSON 字节数组
	// newPointInfoAsBytes, err := json.Marshal(newPointInfo)
	// if err != nil {
	// 	return " ", fmt.Errorf("failed to marshal point info: %v", err)
	// }

	// 将新用户数据存储到账本中
	err = ctx.GetStub().PutState(userID, newUserAsBytes)
	if err != nil {
		return  fmt.Errorf("failed to put user into world state: %v", err)
	}

	// // 将用户的积分映射存储到账本中
	// err = ctx.GetStub().PutState(userID+"_points", newPointInfoAsBytes)
	// if err != nil {
	// 	return " ", fmt.Errorf("failed to put user points into world state: %v", err)
	// }
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

// 更新信息
func (c *UserContract) UpdateUserInfo(ctx contractapi.TransactionContextInterface, PrivatecKey string, newUserName string, userid string) error  {
	// 从账本中获取当前用户信息
	hash := sha256.New()
	hash.Write([]byte(PrivatecKey))
	userID := hex.EncodeToString(hash.Sum(nil))
	if userid == userID {
		userAsBytes, err := ctx.GetStub().GetState(userID)
		if err != nil {
			return fmt.Errorf("failed to read from world state: %v", err)
		}
		if userAsBytes == nil {
			return  fmt.Errorf("the user with ID %s does not exist", userID)
		}

		// 将用户数据反序列化为 UserInfo 对象
		var user typestruct.UserInfo
		err = json.Unmarshal(userAsBytes, &user)
		if err != nil {
			return fmt.Errorf("failed to unmarshal user data: %v", err)
		}

		// 更新用户信息
		if newUserName != "" {
			user.UserName = newUserName
		}

		// 将更新后的用户对象序列化为 JSON 字节数组
		updatedUserAsBytes, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("failed to marshal updated user: %v", err)
		}

		// 将更新后的用户数据存储回账本
		err = ctx.GetStub().PutState(userID, updatedUserAsBytes)
		if err != nil {
			return fmt.Errorf("failed to update user in world state: %v", err)
		}
	} else {
		fmt.Print("Key incorrect")
	}
	return nil
}
