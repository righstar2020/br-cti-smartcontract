package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"bytes"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	ctiContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/cti-contract"
	dataContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/data-contract"
	modelContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/model-contract"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
	userContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/user-contract"
	userPointContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/user-point-contract"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
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

// NonceRecord 结构体用于存储 nonce 相关信息
type NonceRecord struct {
	UserID    string    `json:"userId"`
	Timestamp time.Time `json:"timestamp"`
	Signature []byte    `json:"signature"`
}

// 初始化主合约
func (c *MainContract) InitLedger(ctx contractapi.TransactionContextInterface) (string, error) {
	return "success", nil
}

// 注册用户信息
func (c *MainContract) RegisterUserInfo(ctx contractapi.TransactionContextInterface, msgData []byte) (string, error) {
	return c.UserContract.RegisterUser(ctx, msgData)
}

// 查询模型信息
func (c *MainContract) QueryModelInfo(ctx contractapi.TransactionContextInterface, modelID string) (*typestruct.ModelInfo, error) {
	return c.ModelContract.QueryModelInfo(ctx, modelID)
}

// 查询情报信息
func (c *MainContract) QueryCTIInfo(ctx contractapi.TransactionContextInterface, ctiID string) (*typestruct.CtiInfo, error) {
	return c.CTIContract.QueryCTIInfo(ctx, ctiID)
}

// 查询情报信息(hash)
func (c *MainContract) QueryCTIInfoByCTIHash(ctx contractapi.TransactionContextInterface, ctiHash string) (*typestruct.CtiInfo, error) {
	return c.CTIContract.QueryCTIInfoByCTIHash(ctx, ctiHash)
}




// 查询用户上传的情报
func (c *MainContract) QueryCTIInfoByCreatorUserID(ctx contractapi.TransactionContextInterface, userID string) ([]typestruct.CtiInfo, error) {
	return c.CTIContract.QueryCTIInfoByCreatorUserID(ctx, userID)
}

// 根据cti类别查询
func (c *MainContract) QueryCTIInfoByType(ctx contractapi.TransactionContextInterface, ctiType int) ([]typestruct.CtiInfo, error) {
	return c.CTIContract.QueryCTIInfoByType(ctx, ctiType)
}

// 查询用户信息
func (c *MainContract) QueryUserInfo(ctx contractapi.TransactionContextInterface, userID string) (*typestruct.UserInfo, error) {
	return c.UserContract.QueryUserInfo(ctx, userID)
}

// 查询用户积分信息
func (c *MainContract) QueryUserPointInfo(ctx contractapi.TransactionContextInterface, userID string) (*typestruct.UserPointInfo, error) {
	return c.UserPointContract.QueryUserPointInfo(ctx, userID)
}

// 模型信息分页查询
func (c *MainContract) QueryModelInfoByModelIDWithPagination(ctx contractapi.TransactionContextInterface, modelIDPrefix string, pageSize int, bookmark string) (string, error) {
	return c.ModelContract.QueryModelInfoByModelIDWithPagination(ctx, modelIDPrefix, pageSize, bookmark)
}

// 根据流量类型查询模型信息
func (c *MainContract) QueryModelsByTrafficType(ctx contractapi.TransactionContextInterface, trafficType string) ([]typestruct.ModelInfo, error) {
	return c.ModelContract.QueryModelsByTrafficType(ctx, trafficType)
}

// 查询用户所上传的模型信息
func (c *MainContract) QueryModelsByUserID(ctx contractapi.TransactionContextInterface, userId string) ([]typestruct.ModelInfo, error) {
	return c.ModelContract.QueryModelInfoByCreatorUserID(ctx, userId)
}

// 根据相关CTI查询
func (c *MainContract) QueryModelsByRefCTIId(ctx contractapi.TransactionContextInterface, refCTIId string) ([]typestruct.ModelInfo, error) {
	return c.ModelContract.QueryModelsByRefCTIId(ctx, refCTIId)
}
//分页查询所有情报信息
func (c *MainContract) QueryAllCTIInfoWithPagination(ctx contractapi.TransactionContextInterface, pageSize int, bookmark string) (string, error) {
	return c.CTIContract.QueryAllCTIInfoWithPagination(ctx, pageSize, bookmark)
}

// 根据类型分页查询
func (c *MainContract) QueryCTIInfoByTypeWithPagination(ctx contractapi.TransactionContextInterface, ctiType int, pageSize int, bookmark string) (string, error) {
	return c.CTIContract.QueryCTIInfoByTypeWithPagination(ctx, ctiType, pageSize, bookmark)
}

// CTI精简信息
func (c *MainContract) QueryCTISummaryInfoByCTIID(ctx contractapi.TransactionContextInterface, ctiID string) (*typestruct.CtiSummaryInfo, error) {
	return c.DataContract.QueryCTISummaryInfoByCTIID(ctx, ctiID)
}

// 统计信息，没有入参
func (c *MainContract) GetDataStatistics(ctx contractapi.TransactionContextInterface) (string, error) {
	return c.DataContract.GetDataStatistics(ctx)
}


//--------------------------------------------------------------------以下需要签名验证--------------------------------------------------------------------

// 注册模型信息
func (c *MainContract) RegisterModelInfo(ctx contractapi.TransactionContextInterface, txMsgData []byte) error {
	//验证交易签名(返回交易数据和验证结果)
	txData, err := c.VerifyTxSignature(ctx, txMsgData)
	if err != nil {
		return err
	}
	//验证通过后，注册模型信息
	return c.ModelContract.RegisterModelInfo(ctx, txData)
}

//注册情报信息
func (c *MainContract) RegisterCTIInfo(ctx contractapi.TransactionContextInterface, txMsgData []byte) error {
	//验证交易签名(返回交易数据和验证结果)
	txData, err := c.VerifyTxSignature(ctx, txMsgData)
	if err != nil {
		return err
	}

	//验证通过后，注册情报信息	
	return c.CTIContract.RegisterCTIInfo(ctx, txData)
}




// 用户购买情报
func (c *MainContract) PurchaseCTI(ctx contractapi.TransactionContextInterface, txMsgData []byte) error {
	//验证交易签名(返回交易数据和验证结果)
	txData, err := c.VerifyTxSignature(ctx, txMsgData)
	if err != nil {
		return fmt.Errorf("transaction signature verification failed")
	}
	return c.UserPointContract.PurchaseCTI(ctx, txData)
}

//验证交易随机数和签名
func (c *MainContract) VerifyTxSignature(ctx contractapi.TransactionContextInterface, msgData []byte) ([]byte, error) {
	//解析msgData	
	var txMsgData msgstruct.TxMsgData
	err := json.Unmarshal(msgData, &txMsgData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal msg data: %v", err)
	}

	return txMsgData.TxData, nil

	//暂时取消交易签名验证
	// 在处理实际交易时验证
	// err = c.VerifyTransactionReplay(ctx, txMsgData.Nonce, txMsgData.UserID, txMsgData.NonceSignature)
	// if err != nil {
	// 	// 交易被重放或 nonce 无效
	// 	return nil, fmt.Errorf("transaction replay or nonce invalid: %v", err)
	// }
	// //获取用户公钥pem
	// userInfo, err := ctx.GetStub().GetState(txMsgData.UserID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get user info: %v", err)
	// }
	// //解析用户信息
	// var user typestruct.UserInfo
	// err = json.Unmarshal(userInfo, &user)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to unmarshal user info: %v", err)
	// }
	// //验证交易签名
	// verify, err := utils.Verify(ctx, txMsgData.TxData, []byte(user.PublicKey), txMsgData.TxSignature)	
	// if err != nil {
	// 	return nil, fmt.Errorf("transaction signature verification failed: %v", err)
	// }
	// if !verify {
	// 	return nil, fmt.Errorf("transaction signature verification failed")
	// }
	// return txMsgData.TxData, nil
}

// 生成交易随机数(用于避免重放攻击)
func (c *MainContract) GetTransactionNonce(ctx contractapi.TransactionContextInterface, userID string, txSignature []byte) (string, error) {
	// 生成一个随机的 32 字节数组
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	
	// 将随机字节转换为 base64 编码的字符串
	nonce := base64.StdEncoding.EncodeToString(randomBytes)
	
	// 创建包含时间戳的 nonce 记录
	nonceRecord := struct {
		UserID    string    `json:"userId"`
		Timestamp time.Time `json:"timestamp"`
		Signature []byte    `json:"signature"`
	}{
		UserID:    userID,
		Timestamp: time.Now().UTC(),
		Signature: txSignature,
	}
	
	// 序列化 nonce 记录
	nonceBytes, err := json.Marshal(nonceRecord)
	if err != nil {
		return "", fmt.Errorf("failed to marshal nonce record: %v", err)
	}
	
	// 检查 nonce 是否已存在
	existing, err := ctx.GetStub().GetState(nonce)
	if err != nil {
		return "", fmt.Errorf("failed to check existing nonce: %v", err)
	}
	if existing != nil {
		return "", fmt.Errorf("nonce already exists")
	}
	
	// 保存 nonce 记录
	err = ctx.GetStub().PutState(nonce, nonceBytes)
	if err != nil {
		return "", fmt.Errorf("failed to put nonce state: %v", err)
	}
	//返回随机数(base64 string)
	return nonce, nil
}

// VerifyTransactionReplay 验证交易是否被重放
func (c *MainContract) VerifyTransactionReplay(ctx contractapi.TransactionContextInterface, nonce string, userID string, txSignature []byte) error {
	// 获取 nonce 记录
	nonceBytes, err := ctx.GetStub().GetState(nonce)
	if err != nil {
		return fmt.Errorf("failed to get nonce record: %v", err)
	}
	
	// 检查 nonce 是否存在
	if nonceBytes == nil {
		return fmt.Errorf("nonce does not exist")
	}
	
	// 解析 nonce 记录
	var nonceRecord NonceRecord
	err = json.Unmarshal(nonceBytes, &nonceRecord)
	if err != nil {
		return fmt.Errorf("failed to unmarshal nonce record: %v", err)
	}
	
	// 验证用户ID是否匹配
	if nonceRecord.UserID != userID {
		return fmt.Errorf("user ID mismatch")
	}
	
	// 验证时间戳是否在有效期内（例如：30分钟）
	if time.Since(nonceRecord.Timestamp) > 30*time.Minute {
		return fmt.Errorf("nonce has expired")
	}
	
	// 验证签名是否匹配
	if !bytes.Equal(nonceRecord.Signature, txSignature) {
		return fmt.Errorf("signature mismatch")
	}
	
	// 验证通过后，删除已使用的 nonce
	err = ctx.GetStub().DelState(nonce)
	if err != nil {
		return fmt.Errorf("failed to delete used nonce: %v", err)
	}
	
	return nil
}

// 清理过期的 nonce（可以定期调用）
func (c *MainContract) CleanExpiredNonces(ctx contractapi.TransactionContextInterface) error {
	// 获取所有 nonce
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return fmt.Errorf("failed to get nonce iterator: %v", err)
	}
	defer iterator.Close()
	
	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return fmt.Errorf("failed to iterate nonces: %v", err)
		}
		
		var nonceRecord NonceRecord
		err = json.Unmarshal(kv.Value, &nonceRecord)
		if err != nil {
			continue // 跳过无法解析的记录
		}
		
		// 如果 nonce 已过期（例如：30分钟），则删除
		if time.Since(nonceRecord.Timestamp) > 30*time.Minute {
			err = ctx.GetStub().DelState(kv.Key)
			if err != nil {
				return fmt.Errorf("failed to delete expired nonce: %v", err)
			}
		}
	}
	
	return nil
}

//--------------------------------------------------------------------以下为数据查询(数据统计)--------------------------------------------------------------------

// 获取情报交易趋势数据
func (c *MainContract) GetCTITrafficTrend(ctx contractapi.TransactionContextInterface, timeRange string) (*typestruct.TrafficTrendInfo, error) {
	return c.DataContract.GetCTITrafficTrend(ctx, timeRange)
}

// 获取攻击类型排行
func (c *MainContract) GetAttackTypeRanking(ctx contractapi.TransactionContextInterface) (*typestruct.AttackRankInfo, error) {
	return c.DataContract.GetAttackTypeRanking(ctx)
}

// 获取IOCs类型分布
func (c *MainContract) GetIOCsDistribution(ctx contractapi.TransactionContextInterface) (*typestruct.IOCsDistributionInfo, error) {
	return c.DataContract.GetIOCsDistribution(ctx)
}

// 获取全球IOCs地理分布
func (c *MainContract) GetGlobalIOCsDistribution(ctx contractapi.TransactionContextInterface) (*typestruct.GlobalIOCsInfo, error) {
	return c.DataContract.GetGlobalIOCsDistribution(ctx)
}

// 获取系统概览数据
func (c *MainContract) GetSystemOverview(ctx contractapi.TransactionContextInterface) (*typestruct.SystemOverviewInfo, error) {
	return c.DataContract.GetSystemOverview(ctx)
}


// 获取用户统计数据
func (c *MainContract) GetUserStatistics(ctx contractapi.TransactionContextInterface, userID string) (*userPointContract.UserStatistics, error) {
	userOwnCtiInfo, err := c.UserPointContract.QueryUserStatistics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistics: %v", err)
	}
	UserUploadCtiList,err := c.CTIContract.QueryCTIInfoByCreatorUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistics: %v", err)
	}
	userOwnCtiInfo.UserUploadCount = len(UserUploadCtiList)
	return userOwnCtiInfo, nil
}

// 查询用户积分交易记录
func (c *MainContract) QueryPointTransactions(ctx contractapi.TransactionContextInterface, userID string) ([]*userPointContract.PointTransaction, error) {
	return c.UserPointContract.QueryPointTransactions(ctx, userID)
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
