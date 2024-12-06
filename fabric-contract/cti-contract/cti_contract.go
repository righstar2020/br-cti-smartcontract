package cti_contract

import (
	"encoding/json"
	"fmt"

	"encoding/base64"
	"time"

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
}

// 注册 CTI 信息
func (c *CTIContract) RegisterCTIInfo(ctx contractapi.TransactionContextInterface, userID string, nonce string,ctiTxData msgstruct.CtiTxData) (*typestruct.CtiInfo, error) {
	

	// 从base64编码的nonce中提取随机数
	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	nonceNum := 100000

	if err == nil && len(nonceBytes) >= 3 {
		// 使用前3个字节生成6位随机数
		nonceNum = int(nonceBytes[0])*10000 + int(nonceBytes[1])*100 + int(nonceBytes[2])
		nonceNum = nonceNum % 1000000 // 确保是6位数
	}
	ctiType := 0
	if ctiTxData.CTIType != 0 {
		ctiType = ctiTxData.CTIType
	}
	timestamp := time.Now().Format("0601021504")
	randomNum := fmt.Sprintf("%06d", nonceNum)
	// 生成CTI ID: 类型(2位) + 时间戳(12位,年月日时分) + 随机数(6位)
	ctiID := fmt.Sprintf("%02d%s%s", ctiType, timestamp, randomNum)
	doctype := "cti"
	// 创建新的 CtiInfo 对象
	newCTI := typestruct.CtiInfo{
		CTIID:          ctiID,                                                                      // 生成唯一的 CTI ID
		CTIHash:        ctiTxData.CTIHash,                                                          // 情报HASH(链下生成)
		CTIName:        ctiTxData.CTIName,                                                          // 情报名称
		CTIType:        ctiType,                                                                    // 情报类型
		CTITrafficType: ctiTxData.CTITrafficType,                                                   // 流量情报类型
		CreatorUserID:  userID,                                                                     // 创建者ID
		OpenSource:     ctiTxData.OpenSource,                                                       // 是否开源
		Tags:           ctiTxData.Tags,                                                             // 情报标签
		IOCs:           ctiTxData.IOCs,                                                             // 情报IOCs
		StixData:       ctiTxData.StixData,                                                         // STIX数据
		StatisticInfo:  ctiTxData.StatisticInfo,                                                    // 统计信息
		Description:    ctiTxData.Description,                                                      // 情报描述
		DataSize:       ctiTxData.DataSize,                                                         // 数据大小（B）
		IPFSHash:       ctiTxData.IPFSHash,                                                         // IPFS 地址
		Need:           0,                                                                          // 情报需求量
		Value:          ctiTxData.Value,                                                            // 情报价值（积分）
		CompreValue:    0,                                                                          // 综合价值（积分激励算法定价）
		CreateTime:     time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"), // 情报创建时间
		Doctype:        doctype,
	}

	// 将新 CTI 信息序列化为 JSON 字节数组
	ctiAsBytes, err := json.Marshal(newCTI)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CTI info: %v", err)
	}

	// 使用 CTI ID 作为键将情报数据存储到账本中
	err = ctx.GetStub().PutState(ctiID, ctiAsBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to put CTI info into world state: %v", err)
	}

	return &newCTI, nil
}

func (c *CTIContract) QueryCTIInfo(ctx contractapi.TransactionContextInterface, ctiID string) (*typestruct.CtiInfo, error) {
	// 根据 CTIID 查询数据
	ctiAsBytes, err := ctx.GetStub().GetState(ctiID)
	if err != nil {
		return nil, fmt.Errorf("failed to get state for CTI with ID %s: %v", ctiID, err)
	}
	if ctiAsBytes == nil {
		return nil, fmt.Errorf("the CTI with ID %s does not exist", ctiID)
	}

	// 将获取到的字节数据反序列化为 CtiInfo 结构体
	var ctiInfo typestruct.CtiInfo
	err = json.Unmarshal(ctiAsBytes, &ctiInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal CTI info: %v", err)
	}

	// 返回查询到的 CTI 信息
	return &ctiInfo, nil
}



// 根据CTIHash查询情报信息
func (c *CTIContract) QueryCTIInfoByCTIHash(ctx contractapi.TransactionContextInterface, ctiHash string) (*typestruct.CtiInfo, error) {
	// 构建查询字符串，根据CTIHash进行查询
	queryString := fmt.Sprintf(`{"selector":{"cti_hash":"%s"}}`, ctiHash)

	// 执行查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer resultsIterator.Close()

	// 遍历查询结果
	if resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get query result: %v", err)
		}

		// 将查询结果反序列化为 CtiInfo 结构体
		var ctiInfo typestruct.CtiInfo
		err = json.Unmarshal(queryResponse.Value, &ctiInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal CTI info: %v", err)
		}

		return &ctiInfo, nil
	}

	return nil, fmt.Errorf("the cti with hash %s does not exist", ctiHash)
}

// QueryCTIInfoByCreatorUserID 根据创建者ID查询所有相关情报信息
func (c *CTIContract) QueryCTIInfoByCreatorUserID(ctx contractapi.TransactionContextInterface, userID string) ([]typestruct.CtiInfo, error) {
	// 构建查询字符串，根据创建者ID进行查询
	queryString := fmt.Sprintf(`{"selector":{"creator_user_id":"%s"}}`, userID)

	// 执行查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer resultsIterator.Close()

	// 定义一个切片存储查询结果
	var ctiInfos []typestruct.CtiInfo

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Printf("failed to get query result: %v", err)
			continue
		}

		// 将查询结果反序列化为 CtiInfo 结构体
		var ctiInfo typestruct.CtiInfo
		err = json.Unmarshal(queryResponse.Value, &ctiInfo)
		if err != nil {
			fmt.Printf("failed to unmarshal CTI info: %v", err)
			continue
		}

		// 将结果追加到切片
		ctiInfos = append(ctiInfos, ctiInfo)
	}

	return ctiInfos, nil
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
			fmt.Printf("failed to get next query result: %v", err)
			continue
		}

		// 将查询结果反序列化为 CtiInfo 结构体
		var ctiInfo typestruct.CtiInfo
		err = json.Unmarshal(queryResponse.Value, &ctiInfo)
		if err != nil {
			fmt.Printf("failed to unmarshal CTI info: %v", err)
			continue
		}

		// 将结果追加到切片
		ctiInfos = append(ctiInfos, ctiInfo)
	}

	// 返回查询结果
	return ctiInfos, nil
}

func (c *CTIContract) QueryCTIInfoByTypeWithPagination(ctx contractapi.TransactionContextInterface, ctiType int, pageSize int, bookmark string) (string, error) {
	// 构建查询字符串，根据情报类型查询
	queryString := fmt.Sprintf(`{"selector":{"cti_type":%d}}`, ctiType)

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
			fmt.Printf("failed to get next query result: %v", err)
			continue
		}

		var ctiInfo typestruct.CtiInfo
		err = json.Unmarshal(queryResponse.Value, &ctiInfo)
		if err != nil {
			fmt.Printf("failed to unmarshal query result: %v", err)
			continue
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

// QueryAllCTIInfoWithPagination 分页查询所有情报信息
func (c *CTIContract) QueryAllCTIInfoWithPagination(ctx contractapi.TransactionContextInterface, pageSize int, bookmark string) (string, error) {
	// 构建查询字符串，查询 Doctype 为 "cti" 的所有情报
	queryString := `{"selector":{"doctype":"cti"}}`

	// 执行带分页的查询
	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return "", fmt.Errorf("执行分页查询失败: %v", err)
	}
	defer resultsIterator.Close()

	var ctiInfos []typestruct.CtiInfo

	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Printf("获取下一个查询结果失败: %v", err)
			continue
		}

		var ctiInfo typestruct.CtiInfo
		err = json.Unmarshal(queryResponse.Value, &ctiInfo)
		if err != nil {
			//解析失败，跳过
			fmt.Printf("failed to unmarshal query result: %v", err)
			continue
		}

		ctiInfos = append(ctiInfos, ctiInfo)
	}

	// 构造返回结构
	response := struct {
		CtiInfos []typestruct.CtiInfo `json:"cti_infos"`
		Bookmark string               `json:"bookmark"`
	}{
		CtiInfos: ctiInfos,
		Bookmark: metadata.Bookmark,
	}

	// 序列化为 JSON 字符串
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("序列化响应数据失败: %v", err)
	}

	return string(responseBytes), nil
}

