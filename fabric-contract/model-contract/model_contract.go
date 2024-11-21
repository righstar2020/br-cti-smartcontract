package model_contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
)

// ModelContract 是模型合约的结构体
type ModelContract struct {
	contractapi.Contract
}

// 注册 model信息
func (c *ModelContract) RegisterModelInfo(ctx contractapi.TransactionContextInterface, modelid string, modelname string, traffictype string, trafficfeatures []string, trafficprocesscode string, mlmethod string, mlinfo string, mltraincode string, ipfshash string, refctiid string, privateKey string) error {

	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	// 计算 privateKey 的哈希值
	hash2 := sha256.New()
	hash2.Write([]byte(privateKey))
	userID := hex.EncodeToString(hash2.Sum(nil))
	// 创建新的 CtiInfo 对象
	newModel := typestruct.ModelInfo{
		ModelID:            modelid,
		ModelName:          modelname,
		CreatorUserID:      userID,
		TrafficType:        traffictype,
		TrafficFeatures:    trafficfeatures,
		TrafficProcessCode: trafficprocesscode,
		MLMethod:           mlmethod,
		MLInfo:             mlinfo,
		MLTrainCode:        mltraincode,
		IPFSHashAddress:    ipfshash,
		RefCTIId:           refctiid,
		CreateTime:         time.Unix(txTimestamp.Seconds, int64(txTimestamp.Nanos)).UTC().Format(time.RFC3339),
	}

	// 将新 CTI 信息序列化为 JSON 字节数组
	modelAsBytes, err := json.Marshal(newModel)
	if err != nil {
		return fmt.Errorf("failed to marshal CTI info: %v", err)
	}

	// 使用 CTI ID 作为键将情报数据存储到账本中
	err = ctx.GetStub().PutState(modelid, modelAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put CTI info into world state: %v", err)
	}

	return nil
}

// QueryModelInfo 根据ID查询模型信息
func (c *ModelContract) QueryModelInfo(ctx contractapi.TransactionContextInterface, modelID string) (*typestruct.ModelInfo, error) {
	modelInfoJSON, err := ctx.GetStub().GetState(modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if modelInfoJSON == nil {
		return nil, fmt.Errorf("the model %s does not exist", modelID)
	}

	var modelInfo typestruct.ModelInfo
	err = json.Unmarshal(modelInfoJSON, &modelInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal model info: %v", err)
	}

	return &modelInfo, nil
}

// 其他函数如分页查询等类似实现...
