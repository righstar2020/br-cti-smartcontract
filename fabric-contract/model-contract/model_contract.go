package model_contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/utils"
)

// ModelContract 是模型合约的结构体
type ModelContract struct {
	contractapi.Contract
}

// 注册 model信息
func (c *ModelContract) RegisterModelInfo(ctx contractapi.TransactionContextInterface, modelname string, traffictype string, trafficfeatures []string, trafficprocesscode string, mlmethod string, mlinfo string, mltraincode string, ipfshash string, refctiid string, privateKey string) error {

	modelID, err := utils.GenerateNextModelID(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate random CTI ID: %v", err)
	}
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	// 计算 privateKey 的哈希值
	hash2 := sha256.New()
	hash2.Write([]byte(privateKey))
	userID := hex.EncodeToString(hash2.Sum(nil))
	// 创建新的 ModelInfo 对象
	newModel := typestruct.ModelInfo{
		ModelID:            modelID,
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

	// 将新 Model 信息序列化为 JSON 字节数组
	modelAsBytes, err := json.Marshal(newModel)
	if err != nil {
		return fmt.Errorf("failed to marshal CTI info: %v", err)
	}

	// 使用 Model ID 作为键将情报数据存储到账本中
	err = ctx.GetStub().PutState(modelID, modelAsBytes)
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

// 分页查询模型信息
func (c *ModelContract) QueryModelInfoByModelIDWithPagination(ctx contractapi.TransactionContextInterface, modelIDPrefix string, pageSize int, bookmark string) (string, error) {
	// 构建查询字符串，匹配以 modelIDPrefix 开头的 ModelID
	queryString := fmt.Sprintf(`{"selector":{"model_id":{"$regex":"^%s"}}}`, modelIDPrefix)

	// 执行带分页的查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return "", fmt.Errorf("failed to execute paginated query: %v", err)
	}
	defer resultsIterator.Close()

	var models []typestruct.ModelInfo

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf("failed to get next query result: %v", err)
		}

		var model typestruct.ModelInfo
		err = json.Unmarshal(queryResponse.Value, &model)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal query result: %v", err)
		}

		// 确保 TrafficFeatures 等字段不会为 nil
		if model.TrafficFeatures == nil {
			model.TrafficFeatures = []string{}
		}

		models = append(models, model)
	}

	// 构造返回结构
	response := struct {
		Models   []typestruct.ModelInfo `json:"models"`
		Bookmark string                 `json:"bookmark"`
	}{
		Models:   models,
		Bookmark: metadata.Bookmark,
	}

	// 序列化为 JSON 字符串
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %v", err)
	}

	return string(responseBytes), nil // 返回 JSON 字符串
}


// 根据流量类型查询
func (c *ModelContract) QueryModelsByTrafficType(ctx contractapi.TransactionContextInterface, trafficType string) ([]typestruct.ModelInfo, error) {
	// 构建查询字符串
	queryString := fmt.Sprintf(`{"selector":{"traffic_type":"%s"}}`, trafficType)

	// 执行查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer resultsIterator.Close()

	var models []typestruct.ModelInfo

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var model typestruct.ModelInfo
		err = json.Unmarshal(queryResponse.Value, &model)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}

		models = append(models, model)
	}

	// 返回结果
	return models, nil
}

func (c *ModelContract) QueryModelsByPrivateKey(ctx contractapi.TransactionContextInterface, privateKey string) ([]typestruct.ModelInfo, error) {
	hash := sha256.New()
	hash.Write([]byte(privateKey))
	creatorUserID := hex.EncodeToString(hash.Sum(nil))// 一致的SHA256生成逻辑

    queryString := fmt.Sprintf(`{"selector":{"creator_user_id":"%s"}}`, creatorUserID)

    resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %v", err)
    }
    defer resultsIterator.Close()

    var models []typestruct.ModelInfo
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, fmt.Errorf("failed to get next query result: %v", err)
        }

        var model typestruct.ModelInfo
        err = json.Unmarshal(queryResponse.Value, &model)
        if err != nil {
            return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
        }

        models = append(models, model)
    }

    return models, nil
}


// 根据CTIid查询
func (c *ModelContract) QueryModelsByRefCTIId(ctx contractapi.TransactionContextInterface, refCTIId string) ([]typestruct.ModelInfo, error) {
	// 构建查询字符串
	queryString := fmt.Sprintf(`{"selector":{"ref_cti_id":"%s"}}`, refCTIId)

	// 执行查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer resultsIterator.Close()

	var models []typestruct.ModelInfo

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var model typestruct.ModelInfo
		err = json.Unmarshal(queryResponse.Value, &model)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}

		models = append(models, model)
	}

	// 返回结果
	return models, nil
}

