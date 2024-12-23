package model_contract

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"math"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
)

// ModelContract 是模型合约的结构体
type ModelContract struct {
}

// RegisterModelInfo 注册模型信息
func (c *ModelContract) RegisterModelInfo(ctx contractapi.TransactionContextInterface, userID string, nonce string, modelTxData msgstruct.ModelTxData) (*typestruct.ModelInfo, error) {

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
	block_time, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return nil, fmt.Errorf("获取交易时间戳失败: %v", err)
	}
	currentTime := time.Unix(int64(block_time.GetSeconds()), 0)
	//只需要精确到分钟
	timestamp := currentTime.Format("0601021504")
	randomNum := fmt.Sprintf("%06d", nonceNum)
	// 生成Model ID: 类型(2位) + 时间戳(12位,年月日时分) + 随机数(6位)
	modelID := fmt.Sprintf("%02d%s%s", modelType, timestamp, randomNum)
	doctype := "model"
	existingModel, err := ctx.GetStub().GetState(modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if existingModel != nil {
		return nil, fmt.Errorf("model with ID %s already exists", modelID)
	}
	//创建模型信息
	modelInfo := typestruct.ModelInfo{
		ModelID:             modelID,
		ModelHash:           modelTxData.ModelHash,
		ModelName:           modelTxData.ModelName,
		CreatorUserID:       userID,
		ModelDataType:       modelTxData.ModelDataType,
		ModelType:           modelTxData.ModelType,
		ModelAlgorithm:      modelTxData.ModelAlgorithm,
		ModelTrainFramework: modelTxData.ModelTrainFramework,
		ModelOpenSource:     modelTxData.ModelOpenSource,
		ModelFeatures:       modelTxData.ModelFeatures,
		ModelTags:           modelTxData.ModelTags,
		ModelDescription:    modelTxData.ModelDescription,
		ModelSize:           modelTxData.ModelSize,
		ModelDataSize:       modelTxData.ModelDataSize,
		ModelDataIPFSHash:   modelTxData.ModelDataIPFSHash,
		ModelIPFSHash:       modelTxData.ModelIPFSHash,
		Value:               modelTxData.Value,
		IncentiveMechanism: modelTxData.IncentiveMechanism,
		RefCTIId:            modelTxData.RefCTIId,
		CreateTime:          currentTime.In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"),
		Doctype:             doctype,
	}
	//保留两位小数
	modelInfo.Value = math.Round(modelInfo.Value*100) / 100
	modelInfoJSONBytes, err := json.Marshal(modelInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal model info: %v", err)
	}

	// **步骤 2：使用 modelID 作为键将模型信息存储到账本中**
	err = ctx.GetStub().PutState(modelID, modelInfoJSONBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to put model info into world state: %v", err)
	}
	return &modelInfo, nil
}

func (c *ModelContract) QueryModelInfo(ctx contractapi.TransactionContextInterface, modelID string) (*typestruct.ModelInfo, error) {
	// 根据 ModelID 查询数据
	modelAsBytes, err := ctx.GetStub().GetState(modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get state for model with ID %s: %v", modelID, err)
	}
	if modelAsBytes == nil {
		return nil, fmt.Errorf("the model with ID %s does not exist", modelID)
	}

	// 将获取到的字节数据反序列化为 ModelInfo 结构体
	var modelInfo typestruct.ModelInfo
	err = json.Unmarshal(modelAsBytes, &modelInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal model info: %v", err)
	}

	// 返回查询到的模型信息
	return &modelInfo, nil
}

// QueryAllModelInfoWithPagination 分页查询所有模型信息
func (c *ModelContract) QueryAllModelInfoWithPagination(ctx contractapi.TransactionContextInterface, page int, pageSize int) (*typestruct.ModelQueryResult, error) {
	// 构建查询字符串，查询 Doctype 为 "model" 的所有信息
	queryString := `{"selector":{"doctype":"model"}}`
	_, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(999999999), "") // 极限可获取总数
	if err != nil {
		return nil, fmt.Errorf("获取总数失败: %v", err)
	}
	totalCount := int(metadata.FetchedRecordsCount)

	// 执行查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("执行查询失败: %v", err)
	}
	defer resultsIterator.Close()

	modelInfos := []typestruct.ModelInfo{}

	// 计算偏移量
	offset := pageSize * (page - 1)
	count := 0

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Printf("获取下一个查询结果失败: %v", err)
			continue
		}

		// 跳过偏移量之前的结果
		if count < offset {
			count++
			continue
		}

		var modelInfo typestruct.ModelInfo
		err = json.Unmarshal(queryResponse.Value, &modelInfo)
		if err != nil {
			fmt.Printf("解析查询结果失败: %v", err)
			continue
		}

		modelInfos = append(modelInfos, modelInfo)
		count++

		// 如果达到页面大小，停止
		if len(modelInfos) >= pageSize {
			break
		}
	}

	// 构造返回结构
	modelQueryResult := typestruct.ModelQueryResult{
		ModelInfos: modelInfos,
		Total:      totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	return &modelQueryResult, nil
}
// QueryModelsByModelTypeWithPagination 根据模型类型分页查询
func (c *ModelContract) QueryModelsByTypeWithPagination(ctx contractapi.TransactionContextInterface, modelType int, page int, pageSize int) (*typestruct.ModelQueryResult, error) {
	// 构建查询字符串，根据 ModelType 进行查询
	queryString := fmt.Sprintf(`{"selector":{"doctype":"model", "model_type":%d}}`, modelType)
	_, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(999999999), "") // 极限可获取总数
	if err != nil {
		return nil, fmt.Errorf("获取总数失败: %v", err)
	}
	totalCount := int(metadata.FetchedRecordsCount)
	// 执行查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rich query: %v", err)
	}
	defer resultsIterator.Close()

	// 定义一个切片存储查询结果
	var modelInfos []typestruct.ModelInfo

	// 计算偏移量
	offset := pageSize * (page - 1)
	count := 0

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Printf("failed to get next query result: %v", err)
			continue
		}

		// 跳过偏移量之前的结果
		if count < offset {
			count++
			continue
		}

		// 将查询结果反序列化为 ModelInfo 结构体
		var modelInfo typestruct.ModelInfo
		err = json.Unmarshal(queryResponse.Value, &modelInfo)
		if err != nil {
			fmt.Printf("failed to unmarshal model info: %v", err)
			continue
		}

		// 将结果追加到切片
		modelInfos = append(modelInfos, modelInfo)
		count++

		// 如果达到页面大小，停止
		if len(modelInfos) >= pageSize {
			break
		}
	}

	// 返回查询结果
	return &typestruct.ModelQueryResult{
		ModelInfos: modelInfos,
		Total:      totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
//按激励机制分页查询
func (c *ModelContract) QueryModelsByIncentiveMechanismWithPagination(ctx contractapi.TransactionContextInterface,page int, pageSize int,incentiveMechanism int) (*typestruct.ModelQueryResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 构造查询字符串
	queryString := fmt.Sprintf(`{"selector":{"doctype":"model","incentive_mechanism":%d}}`, incentiveMechanism)

	_, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(999999999), "") // 极限可获取总数
	if err != nil {
		return nil, fmt.Errorf("获取总数失败: %v", err)
	}
	totalCount := int(metadata.FetchedRecordsCount)

	// 执行查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("执行查询失败: %v", err)
	}
	defer resultsIterator.Close()
	
	modelInfos := []typestruct.ModelInfo{}

	// 计算偏移量
	offset := pageSize * (page - 1)
	count := 0

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Printf("获取下一个查询结果失败: %v", err)
			continue
		}

		// 跳过偏移量之前的结果
		if count < offset {
			count++
			continue
		}

		var modelInfo typestruct.ModelInfo
		err = json.Unmarshal(queryResponse.Value, &modelInfo)
		if err != nil {
			fmt.Printf("解析查询结果失败: %v", err)
			continue
		}

		modelInfos = append(modelInfos, modelInfo)
		count++

		// 如果达到页面大小，停止
		if len(modelInfos) >= pageSize {
			break
		}
	}

	// 构造返回结构
	queryResult := &typestruct.ModelQueryResult{
		ModelInfos: modelInfos,
		Total:    totalCount,
		Page:     page,
		PageSize: pageSize,
	}

	return queryResult, nil
}



// QueryModelsByRefCTIId 根据CTIid查询
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
//根据创建者ID查询相关情报总数量
func (c *ModelContract) QueryModelTotalCountByCreatorUserID(ctx contractapi.TransactionContextInterface, userID string) (int, error) {
	// 构建查询字符串，根据创建者ID进行查询
	queryString := fmt.Sprintf(`{"selector":{"creator_user_id":"%s","doctype":"model"}}`, userID)

	// 执行查询
	_,metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString,9999999,"")
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}
	totalCount:=metadata.FetchedRecordsCount
	if totalCount >0 {
		return int(totalCount),nil
	}

	return 0, nil
}
// QueryModelInfoByCreatorUserID 根据创建者ID查询
func (c *ModelContract) QueryModelInfoByCreatorUserID(ctx contractapi.TransactionContextInterface, userId string) ([]typestruct.ModelInfo, error) {
	// 构建查询字符串
	queryString := fmt.Sprintf(`{"selector":{"creator_user_id":"%s","doctype":"model"}}`, userId)



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
//----------------------------------更新模型信息函数----------------------------------
//更新模型信息函数(Value)
func (c *ModelContract) UpdateModelValue(ctx contractapi.TransactionContextInterface, modelID string, value float64) error {
	modelInfo, err := c.QueryModelInfo(ctx, modelID)
	if err != nil {
		return fmt.Errorf("failed to query model info: %v", err)
	}
	modelInfo.Value = value
	modelAsBytes, err := json.Marshal(modelInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal model info: %v", err)
	}
	return ctx.GetStub().PutState(modelID, modelAsBytes)
}

//更新模型信息函数(Need)
func (c *ModelContract) UpdateModelNeedAdd(ctx contractapi.TransactionContextInterface, modelID string, need int) error {
	modelInfo, err := c.QueryModelInfo(ctx, modelID)
	if err != nil {
		return fmt.Errorf("failed to query model info: %v", err)
	}
	if modelInfo.Need == 0 {
		modelInfo.Need = need
	} else {
		modelInfo.Need += need
	}
	modelAsBytes, err := json.Marshal(modelInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal model info: %v", err)
	}
	return ctx.GetStub().PutState(modelID, modelAsBytes)
}

