
package cti_contract

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)
const (
    // 攻击类型
    AttackType_Traffic         = 1 // 恶意流量
    AttackType_Malware         = 2 // 恶意软件  
    AttackType_Phishing        = 3 // 钓鱼地址
    AttackType_Botnet          = 4 // 僵尸网络
    AttackType_AppLayer        = 5 // 应用层攻击
    AttackType_OpenSource      = 6 // 开源情报
)
//流量情报类型
const (
	CTITrafficType_5G = 1 // 5G
	CTITrafficType_Satellite = 2 // 卫星网络
	CTITrafficType_SDN = 3 // SDN
)
//----------示例数据-----------------/
//情报标签
var Tags_List = []string{"卫星网络", "SDN网络", "5G网络", "恶意软件", "DDoS", "钓鱼", "僵尸网络", "APT", "IOT"}
//情报IOCs
var IOCs_List = []string{"IP", "端口", "流特征", "HASH", "URL", "CVE"}
//情报统性信息
var SatisticInfo = map[string]string{"location": "中国"}

// CTIInfo 结构体表示情报信息
type CTIInfo struct {
    CTIId              string   `json:"cti_id"` // 情报ID
    CTIName            string   `json:"cti_name"` // 情报名称
    CTIType            int      `json:"cti_type"` // 情报类型 
    CTITrafficType     int      `json:"cti_traffic_type"` // 情报流量类型(5G、卫星网络、SDN)
    OpenSource         int      `json:"open_source"` // 情报来源
    CreatorUserID      string   `json:"creator_user_id"` // 创建者用户ID
    Tags               []string `json:"tags"` // 情报标签
    IOCs               []string `json:"iocs"` // 情报IOCs
    SatisticInfo       map[string]string `json:"satistic_info"` // 情报统性信息
    STIXData           string   `json:"stix_data"` // STIX数据
    Description        string   `json:"description"` // 情报描述
    DataSize           int      `json:"data_size"` // 情报数据大小
    DataHash           string   `json:"data_hash"` // 情报数据哈希
    IPFSHash           string   `json:"ipfs_hash"` // IPFS哈希
    Need               int      `json:"need"` // 情报需求
    Value              int      `json:"value"` // 情报价值(用户指定)
    CompreValue        int      `json:"compre_value"` // 情报综合价值(平台评估)
    CreateTime         string   `json:"create_time"` // 创建时间
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