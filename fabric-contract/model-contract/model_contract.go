package model_contract

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ModelInfo 结构体表示模型信息
type ModelInfo struct {
    ModelID            string   `json:"model_id"`
    ModelName          string   `json:"model_name"`
    CreatorUserID      string   `json:"creator_user_id"`
    TrafficType        string   `json:"traffic_type"`
    TrafficFeatures    []string `json:"traffic_features"`
    TrafficProcessCode string   `json:"traffic_process_code"`
    MLMethod           string   `json:"ml_method"`
    MLInfo             string   `json:"ml_info"`
    MLTrainCode        string   `json:"ml_train_code"`
    IPFSHashAddress    string   `json:"ipfs_hash_address"`
    RefCTIId           string   `json:"ref_cti_id"`
    CreateTime         string   `json:"create_time"`
}

// ModelContract 是模型合约的结构体
type ModelContract struct {
    contractapi.Contract
}

// RegisterModelInfo 注册模型信息
func (c *ModelContract) RegisterModelInfo(ctx contractapi.TransactionContextInterface, modelInfoJSON string, userID string, txSignature string, nonceSignature string) error {
    var modelInfo ModelInfo
    err := json.Unmarshal([]byte(modelInfoJSON), &modelInfo)
    if err != nil {
        return fmt.Errorf("failed to unmarshal model info: %v", err)
    }

    // TODO: 验证签名和随机数

    modelInfoJSONBytes, _ := json.Marshal(modelInfo)
    err = ctx.GetStub().PutState(modelInfo.ModelID, modelInfoJSONBytes)
    if err != nil {
        return fmt.Errorf("failed to put state: %v", err)
    }

    return nil
}

// QueryModelInfo 根据ID查询模型信息
func (c *ModelContract) QueryModelInfo(ctx contractapi.TransactionContextInterface, modelID string) (*ModelInfo, error) {
    modelInfoJSON, err := ctx.GetStub().GetState(modelID)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if modelInfoJSON == nil {
        return nil, fmt.Errorf("the model %s does not exist", modelID)
    }

    var modelInfo ModelInfo
    err = json.Unmarshal(modelInfoJSON, &modelInfo)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal model info: %v", err)
    }

    return &modelInfo, nil
}

// 其他函数如分页查询等类似实现...