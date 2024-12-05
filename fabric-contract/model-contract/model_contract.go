package model_contract

import (
	"encoding/json"
	"fmt"
    "time"
	"encoding/base64"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
)

// ModelContract 是模型合约的结构体
type ModelContract struct {
	contractapi.Contract
}

// RegisterModelInfo 注册模型信息
func (c *ModelContract) RegisterModelInfo(ctx contractapi.TransactionContextInterface,userID string,txData []byte,nonce string) error {
	//解析交易数据
	var modelTxData msgstruct.ModelTxData
	err := json.Unmarshal(txData, &modelTxData)
	if err != nil {
        return fmt.Errorf("failed to unmarshal model info: %v", err)
    }
    // 从base64编码的nonce中提取随机数
    nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
    nonceNum := 100000
    
    if err == nil && len(nonceBytes) >= 3 {
        // 使用前3个字节生成6位随机数
        nonceNum = int(nonceBytes[0])*10000 + int(nonceBytes[1])*100 + int(nonceBytes[2])
        nonceNum = nonceNum % 1000000 // 确保是6位数
    }
    modelType := 0
    if modelTxData.ModelType != 0 {
        modelType = modelTxData.ModelType
    }
    timestamp := time.Now().Format("0601021504")
    randomNum := fmt.Sprintf("%06d", nonceNum)
    // 生成Model ID: 类型(2位) + 时间戳(12位,年月日时分) + 随机数(6位)
    modelID := fmt.Sprintf("%02d%s%s", modelType, timestamp, randomNum)
    //创建模型信息
    modelInfo := typestruct.ModelInfo{
        ModelID: modelID,
        ModelName: modelTxData.ModelName,
        ModelTrafficType: modelTxData.ModelTrafficType,
        ModelType: modelTxData.ModelType,
        ModelHash: modelTxData.ModelHash,
        ModelOpenSource: modelTxData.ModelOpenSource,
        ModelFeatures: modelTxData.ModelFeatures,
        ModelTags: modelTxData.ModelTags,
        ModelDescription: modelTxData.ModelDescription,
        ModelDataSize: modelTxData.ModelDataSize,
        ModelIPFSHash: modelTxData.ModelIPFSHash,
        ModelCreateTime: time.Now().Format("2024-01-01 00:00:00"),
    }

    modelInfoJSONBytes, _ := json.Marshal(modelInfo)
    err = ctx.GetStub().PutState(modelID, modelInfoJSONBytes)
    if err != nil {
        return fmt.Errorf("failed to put state: %v", err)
    }

	// 使用 Model ID 作为键将情报数据存储到账本中
	err = ctx.GetStub().PutState(modelID, modelInfoJSONBytes)
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
		if model.ModelFeatures == nil {
			model.ModelFeatures = []string{}
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
//QueryModelInfoByCreatorUserID 根据创建者ID查询
func (c *ModelContract) QueryModelInfoByCreatorUserID(ctx contractapi.TransactionContextInterface, userId string) ([]typestruct.ModelInfo, error) {
	// 构建查询字符串
	queryString := fmt.Sprintf(`{"selector":{"model_creator_user_id":"%s"}}`, userId)

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
