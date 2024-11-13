
package cti_contract

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// CTIInfo 结构体表示情报信息
type CTIInfo struct {
    CTIId              string   `json:"cti_id"`
    CTIName            string   `json:"cti_name"`
    CTIType            int      `json:"cti_type"`
    CTITrafficType     int      `json:"cti_traffic_type"`
    OpenSource         int      `json:"open_source"`
    CreatorUserID      string   `json:"creator_user_id"`
    Tags               []string `json:"tags"`
    IOCs               []string `json:"iocs"`
    STIXData           string   `json:"stix_data"`
    Description        string   `json:"description"`
    DataSize           int      `json:"data_size"`
    DataHash           string   `json:"data_hash"`
    IPFSHash           string   `json:"ipfs_hash"`
    Need               int      `json:"need"`
    Value              int      `json:"value"`
    CompreValue        int      `json:"compre_value"`
    CreateTime         string   `json:"create_time"`
}

// CTIContract 是情报合约的结构体
type CTIContract struct {
    contractapi.Contract
}

// RegisterCTIInfo 注册情报信息
func (c *CTIContract) RegisterCTIInfo(ctx contractapi.TransactionContextInterface, ctiInfoJSON string, userID string, txSignature string, nonceSignature string) error {
    var ctiInfo CTIInfo
    err := json.Unmarshal([]byte(ctiInfoJSON), &ctiInfo)
    if err != nil {
        return fmt.Errorf("failed to unmarshal cti info: %v", err)
    }

    // TODO: 验证签名和随机数

    ctiInfoJSONBytes, _ := json.Marshal(ctiInfo)
    err = ctx.GetStub().PutState(ctiInfo.CTIId, ctiInfoJSONBytes)
    if err != nil {
        return fmt.Errorf("failed to put state: %v", err)
    }

    return nil
}

// QueryCTIInfo 根据ID查询情报信息
func (c *CTIContract) QueryCTIInfo(ctx contractapi.TransactionContextInterface, ctiID string) (*CTIInfo, error) {
    ctiInfoJSON, err := ctx.GetStub().GetState(ctiID)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if ctiInfoJSON == nil {
        return nil, fmt.Errorf("the cti %s does not exist", ctiID)
    }

    var ctiInfo CTIInfo
    err = json.Unmarshal(ctiInfoJSON, &ctiInfo)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal cti info: %v", err)
    }

    return &ctiInfo, nil
}

// 其他函数如分页查询等类似实现...