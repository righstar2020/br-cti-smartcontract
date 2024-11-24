package cti_contract

import (
	"encoding/json"
	"fmt"

	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
)

const (
	// 攻击类型
	AttackType_Traffic    = 1 // 恶意流量
	AttackType_Malware    = 2 // 恶意软件
	AttackType_Phishing   = 3 // 钓鱼地址
	AttackType_Botnet     = 4 // 僵尸网络
	AttackType_AppLayer   = 5 // 应用层攻击
	AttackType_OpenSource = 6 // 开源情报
)

// 流量情报类型
const (
	CTITrafficType_5G        = 1 // 5G
	CTITrafficType_Satellite = 2 // 卫星网络
	CTITrafficType_SDN       = 3 // SDN
)

// ----------示例数据-----------------/
// 情报标签
var Tags_List = []string{"卫星网络", "SDN网络", "5G网络", "恶意软件", "DDoS", "钓鱼", "僵尸网络", "APT", "IOT"}

// 情报IOCs
var IOCs_List = []string{"IP", "端口", "流特征", "HASH", "URL", "CVE"}

// 情报统性信息
var SatisticInfo = map[string]interface{}{
	"location": map[string]int{
		"中国":  1,
		"美国":  2,
		"俄罗斯": 3,
		"欧洲":  4,
		"亚洲":  5,
		"非洲":  6,
		"南美洲": 7,
		"北美洲": 8,
		"大洋洲": 9,
	},
}

// CTIContract 是情报合约的结构体
type CTIContract struct {
	contractapi.Contract
}

// 注册 CTI 信息
func (c *CTIContract) RegisterCTIInfo(ctx contractapi.TransactionContextInterface,txData []byte)  error {
	//解析交易数据	
	var ctiTxData msgstruct.CtiTxData
	err := json.Unmarshal(txData, &ctiTxData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal msg data: %v", err)
	}
	
	// 生成随机的 CTI ID
	ctiID := uuid.New().String()	

	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	// 创建新的 CtiInfo 对象
	newCTI := typestruct.CtiInfo{
		CTIID:          ctiID,                                                                               // 生成唯一的 CTI ID
		CTIName:        ctiTxData.CTIName,                                                                             // 情报名称
		CTITrafficType: ctiTxData.CTITrafficType,                                                                      // 流量类型
		OpenSource:     ctiTxData.OpenSource,                                                                          // 是否开源
		Tags:           ctiTxData.Tags,                                                                                // 情报标签
		IOCs:           ctiTxData.IOCs,                                                                                         // 情报IOCs
		StixData:       ctiTxData.StixData,                                                                            // STIX数据
		StatisticInfo:  ctiTxData.StatisticInfo,                                                                         // 统计信息
		Description:    ctiTxData.Description,                                                                         // 情报描述
		DataSize:       ctiTxData.DataSize,                                                                            // 数据大小（B）
		IPFSHash:       ctiTxData.IPFSHash,                                                                                 // IPFS 地址
		Need:           ctiTxData.Need,                                                                                                 // 情报需求量
		Value:          ctiTxData.Value,                                                                                                 // 情报价值（积分）
		CompreValue:    ctiTxData.CompreValue,                                                                                                 // 综合价值（积分激励算法定价）
		CreateTime:     time.Unix(txTimestamp.Seconds, int64(txTimestamp.Nanos)).UTC().Format(time.RFC3339),              // 情报创建时间
	}

	// 将新 CTI 信息序列化为 JSON 字节数组
	ctiAsBytes, err := json.Marshal(newCTI)
	if err != nil {
		return fmt.Errorf("failed to marshal CTI info: %v", err)
	}

	// 使用 CTI ID 作为键将情报数据存储到账本中
	err = ctx.GetStub().PutState(ctiID, ctiAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put CTI info into world state: %v", err)
	}

	return nil
}

// QueryCTIInfo 根据ID查询情报信息
func (c *CTIContract) QueryCTIInfo(ctx contractapi.TransactionContextInterface, ctiID string) (*typestruct.CtiInfo, error) {
	ctiInfoJSON, err := ctx.GetStub().GetState(ctiID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if ctiInfoJSON == nil {
		return nil, fmt.Errorf("the cti %s does not exist", ctiID)
	}

	var ctiInfo typestruct.CtiInfo
	err = json.Unmarshal(ctiInfoJSON, &ctiInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cti info: %v", err)
	}

	return &ctiInfo, nil
}

func (c *CTIContract) QueryCTIInfoByCTIIDWithPagination(ctx contractapi.TransactionContextInterface, ctiIDPrefix string, pageSize int, bookmark string) (string, error) {
	// 构建查询字符串，匹配以 ctiIDPrefix 开头的 CTIID
	queryString := fmt.Sprintf(`{"selector":{"cti_id":{"$regex":"^%s"}}}`, ctiIDPrefix)

	// 执行带分页的查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return "", fmt.Errorf("failed to execute paginated query: %v", err)
	}
	defer resultsIterator.Close()

	var ctiInfos []typestruct.CtiInfo

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf("failed to get next query result: %v", err)
		}

		var ctiInfo typestruct.CtiInfo
		err = json.Unmarshal(queryResponse.Value, &ctiInfo)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal query result: %v", err)
		}

		ctiInfos = append(ctiInfos, ctiInfo)
	}

	// 构造返回结构
	response := struct {
		CtiInfos []typestruct.CtiInfo `json:"cti_infos"`
		Bookmark string               `json:"bookmark"` // 分页标记
	}{
		CtiInfos: ctiInfos,
		Bookmark: metadata.Bookmark,
	}

	// 序列化为 JSON 字符串
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %v", err)
	}

	return string(responseBytes), nil // 返回 JSON 字符串
}



// 根据 CTIType 查询所有相关的 CTIInfo
func (c *CTIContract) QueryCTIInfoByType(ctx contractapi.TransactionContextInterface, ctiType int) ([]typestruct.CtiInfo, error) {
	// 构建查询字符串，根据 CTIType 进行查询
	queryString := fmt.Sprintf(`{"selector":{"cti_type":%d}}`, ctiType)

	// 执行查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rich query: %v", err)
	}
	defer resultsIterator.Close()

	// 定义一个切片存储查询结果
	var ctiInfos []typestruct.CtiInfo

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		// 将查询结果反序列化为 CtiInfo 结构体
		var ctiInfo typestruct.CtiInfo
		err = json.Unmarshal(queryResponse.Value, &ctiInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}

		// 将结果追加到切片
		ctiInfos = append(ctiInfos, ctiInfo)
	}

	// 返回查询结果
	return ctiInfos, nil
}
