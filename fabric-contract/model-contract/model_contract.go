package model_contract

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
)

// ModelContract 是模型合约的结构体
type ModelContract struct {
}

// RegisterModelInfo 注册模型信息
func (c *ModelContract) RegisterModelInfo(ctx contractapi.TransactionContextInterface, userID string, nonce string,modelTxData msgstruct.ModelTxData) (*typestruct.ModelInfo, error) {

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
	doctype := "model"
	//创建模型信息
	modelInfo := typestruct.ModelInfo{
		ModelID:            modelID,
		ModelName:          modelTxData.ModelName,
		ModelTrafficType:   modelTxData.ModelTrafficType,
		ModelType:          modelTxData.ModelType,
		ModelHash:          modelTxData.ModelHash,
		ModelOpenSource:    modelTxData.ModelOpenSource,
		ModelCreatorUserID: userID,
		ModelFeatures:      modelTxData.ModelFeatures,
		ModelTags:          modelTxData.ModelTags,
		ModelDescription:   modelTxData.ModelDescription,
		ModelDataSize:      modelTxData.ModelDataSize,
		ModelIPFSHash:      modelTxData.ModelIPFSHash,
		ModelCreateTime:    time.Now().Format("2024-01-01 00:00:00"),
		Doctype:            doctype,
	}

	modelInfoJSONBytes, _ := json.Marshal(modelInfo)
	err = ctx.GetStub().PutState(doctype, modelInfoJSONBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to put state: %v", err)
	}

	// 使用 Model ID 作为键将情报数据存储到账本中
	err = ctx.GetStub().PutState(modelID, modelInfoJSONBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to put CTI info into world state: %v", err)
	}

	return &modelInfo, nil
}

func (c *ModelContract) QueryModelInfo(ctx contractapi.TransactionContextInterface, modelID string) (*typestruct.ModelInfo, error) {
	// 根据 CTIID 查询数据
	ctiAsBytes, err := ctx.GetStub().GetState(modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get state for CTI with ID %s: %v", modelID, err)
	}
	if ctiAsBytes == nil {
		return nil, fmt.Errorf("the CTI with ID %s does not exist", modelID)
	}

	// 将获取到的字节数据反序列化为 CtiInfo 结构体
	var modelInfo typestruct.ModelInfo
	err = json.Unmarshal(ctiAsBytes, &modelInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal CTI info: %v", err)
	}

	// 返回查询到的 CTI 信息
	return &modelInfo, nil
}

// / QueryAllCTIInfoWithPagination 分页查询所有模型信息
func (c *ModelContract) QueryModelInfoByModelIDWithPagination(ctx contractapi.TransactionContextInterface, pageSize int, bookmark string) (string, error) {
	// 构建查询字符串，查询 Doctype 为 "model" 的所有信息
	queryString := `{"selector":{"doctype":"model"}}`

	// 执行带分页的查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return "", fmt.Errorf("执行分页查询失败: %v", err)
	}
	defer resultsIterator.Close()

	var modelInfos []typestruct.ModelInfo

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Printf("获取下一个查询结果失败: %v", err)
			continue
		}

		var modelInfo typestruct.ModelInfo
		err = json.Unmarshal(queryResponse.Value, &modelInfo)
		if err != nil {
			// 解析失败，跳过
			fmt.Printf("failed to unmarshal query result: %v", err)
			continue
		}

		modelInfos = append(modelInfos, modelInfo)
	}

	// 构造返回结构
	response := struct {
		ModelInfos []typestruct.ModelInfo `json:"model_infos"`
		Bookmark   string                 `json:"bookmark"`
	}{
		ModelInfos: modelInfos,
		Bookmark:   metadata.Bookmark,
	}

	// 序列化为 JSON 字符串
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("序列化响应数据失败: %v", err)
	}

	return string(responseBytes), nil
}

// 根据流量类型查询
// 根据 CTIType 查询所有相关的 CTIInfo
func (c *ModelContract) QueryModelsByTrafficType(ctx contractapi.TransactionContextInterface, modelType int) ([]typestruct.ModelInfo, error) {
	// 构建查询字符串，根据 CTIType 进行查询
	queryString := fmt.Sprintf(`{"selector":{"model_traffic_type":%d}}`, modelType)

	// 执行查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rich query: %v", err)
	}
	defer resultsIterator.Close()

	// 定义一个切片存储查询结果
	var modelInfos []typestruct.ModelInfo

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Printf("failed to get next query result: %v", err)
			continue
		}

		// 将查询结果反序列化为 CtiInfo 结构体
		var modelInfo typestruct.ModelInfo
		err = json.Unmarshal(queryResponse.Value, &modelInfo)
		if err != nil {
			fmt.Printf("failed to unmarshal CTI info: %v", err)
			continue
		}

		// 将结果追加到切片
		modelInfos = append(modelInfos, modelInfo)
	}

	// 返回查询结果
	return modelInfos, nil
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

// QueryModelInfoByCreatorUserID 根据创建者ID查询
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
