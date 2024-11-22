package model_contract

import (
	"encoding/json"
	"fmt"
    "time"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
)

// ModelContract 是模型合约的结构体
type ModelContract struct {
    contractapi.Contract
}

// RegisterModelInfo 注册模型信息
func (c *ModelContract) RegisterModelInfo(ctx contractapi.TransactionContextInterface,txData []byte) error {
	//解析交易数据
	var modelTxData msgstruct.ModelTxData
	err := json.Unmarshal(txData, &modelTxData)
	if err != nil {
        return fmt.Errorf("failed to unmarshal model info: %v", err)
    }
    //创建模型ID(使用uuid生成)
    modelID := uuid.New().String()
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