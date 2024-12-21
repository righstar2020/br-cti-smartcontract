package incentive_contract

import (
	"fmt"
	"time"
	"math"
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
	"encoding/base64"
	ctiContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/cti-contract"
	modelContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/model-contract"
	commentContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/comment-contract"
)

// 激励机制类型
const (
	INCENTIVE_TYPE_POINT = 1 // 积分激励
	INCENTIVE_TYPE_GAME = 2  // 三方博弈
	INCENTIVE_TYPE_EVOLUTION = 3 // 演化博弈
)

// IncentiveContract 是激励合约的结构体
type IncentiveContract struct {
	ctiContract.CTIContract
	modelContract.ModelContract
	commentContract.CommentContract
}

//注册文档激励信息
func (c *IncentiveContract) RegisterDocIncentiveInfo(ctx contractapi.TransactionContextInterface, refID string, doctype string, nonce string, totalUserNum int) (*typestruct.DocIncentiveInfo, error) {
	var ctiInfo *typestruct.CtiInfo
	var modelInfo *typestruct.ModelInfo
	var err error
	historyValue := 10.0
	need := 1
	incentiveMechanism := 1
	if doctype == "cti" {
		ctiInfo, err = c.CTIContract.QueryCTIInfo(ctx, refID)
		if err != nil {
			return nil, fmt.Errorf("获取CTI信息失败: %v", err)
		}
		historyValue = ctiInfo.Value
		need = ctiInfo.Need
		incentiveMechanism = ctiInfo.IncentiveMechanism
	} else if doctype == "model" {
		modelInfo, err = c.ModelContract.QueryModelInfo(ctx, refID)
		if err != nil {
			return nil, fmt.Errorf("获取模型信息失败: %v", err)
		}
		historyValue = modelInfo.Value
		need = modelInfo.Need
		incentiveMechanism = modelInfo.IncentiveMechanism
	} else {
		return nil, fmt.Errorf("文档类型错误: %v", doctype)
	}
	if err != nil {
		return nil, fmt.Errorf("获取文档信息失败: %v", err)
	}
	// 获取评论信息
	commentInfos, err := c.QueryAllCommentsByRefID(ctx, refID)
	if err != nil {
		return nil, fmt.Errorf("获取评论信息失败: %v", err)
	}
	commentScore := 60.0 //初始评论分数
	for _, commentInfo := range commentInfos {
		commentScore += commentInfo.CommentScore
	}
	if len(commentInfos) > 0 {
		commentScore = commentScore / float64(len(commentInfos)+1) //加上初始评论分数
	}
	//生成激励ID
	incentiveID, err := c.GenerateIncentiveID(ctx, refID, doctype, nonce)
	if err != nil {
		return nil, fmt.Errorf("生成激励ID失败: %v", err)
	} 
	docIncentiveInfo := typestruct.DocIncentiveInfo{
		IncentiveID: incentiveID,
		RefID: refID,
		IncentiveDoctype: doctype,
		HistoryValue: historyValue,
		IncentiveMechanism: incentiveMechanism,
		CommentScore: commentScore,
		Need: need,
		TotalUserNum: totalUserNum,
		IncentiveValue: 10,//初始激励值
		Doctype: "incentive",
		CreateTime: time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"),
	}
	//异常参数处理
	if docIncentiveInfo.HistoryValue < 1 {
		docIncentiveInfo.HistoryValue = 10
	}
	if docIncentiveInfo.CommentScore < 1 {
		docIncentiveInfo.CommentScore = 1
	}
	if docIncentiveInfo.Need < 1 {
		docIncentiveInfo.Need = 1
	}
	if docIncentiveInfo.TotalUserNum < 1 {
		docIncentiveInfo.TotalUserNum = 1
	}

	//根据激励机制计算激励值
	if incentiveMechanism == INCENTIVE_TYPE_POINT {
		incentiveValue, err := c.CalculateCommonPointIncentive(ctx, &docIncentiveInfo)
		if err != nil {
			return nil, fmt.Errorf("计算激励值失败: %v", err)
		}
		fmt.Println("计算激励值成功: ", incentiveValue)
		docIncentiveInfo.IncentiveValue = incentiveValue
	} else if incentiveMechanism == INCENTIVE_TYPE_GAME {
		incentiveValue, err := c.CalculateThreePartyGameIncentive(ctx, &docIncentiveInfo)
		if err != nil {
			return nil, fmt.Errorf("计算激励值失败: %v", err)
		}
		fmt.Println("计算激励值成功: ", incentiveValue)
		docIncentiveInfo.IncentiveValue = incentiveValue
	} else if incentiveMechanism == INCENTIVE_TYPE_EVOLUTION {
		incentiveValue, err := c.CalculateEvolutionGameIncentive(ctx, &docIncentiveInfo)
		if err != nil {
			return nil, fmt.Errorf("计算激励值失败: %v", err)
		}
		fmt.Println("计算激励值成功: ", incentiveValue)
		docIncentiveInfo.IncentiveValue = incentiveValue
	} else {
		incentiveValue, err := c.CalculateCommonPointIncentive(ctx, &docIncentiveInfo)
		if err != nil {
			return nil, fmt.Errorf("计算激励值失败: %v", err)
		}
		fmt.Println("计算激励值成功: ", incentiveValue)
		docIncentiveInfo.IncentiveValue = incentiveValue
	}
	//异常参数处理
	if docIncentiveInfo.IncentiveValue < 1 {
		docIncentiveInfo.IncentiveValue = 1
	}
	fmt.Println("docIncentiveInfo: ", docIncentiveInfo)
	//将激励信息写入区块链
	docIncentiveInfoBytes, err := json.Marshal(docIncentiveInfo)
	if err != nil {
		return nil, fmt.Errorf("序列化文档激励信息失败: %v", err)
	}
	err = ctx.GetStub().PutState(incentiveID, docIncentiveInfoBytes)
	if err != nil {
		return nil, fmt.Errorf("写入文档激励信息失败: %v", err)
	}
	//更新文档信息
	if doctype == "cti" {
		c.UpdateCTIValue(ctx, refID, docIncentiveInfo.IncentiveValue)
	} else if doctype == "model" {
		c.UpdateModelValue(ctx, refID, docIncentiveInfo.IncentiveValue)
	}
	return &docIncentiveInfo, nil
}
//生成ID
func (c *IncentiveContract) GenerateIncentiveID(ctx contractapi.TransactionContextInterface, refID string, doctype string,nonce string) (string, error) {
	// 从base64编码的nonce中提取随机数
	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	nonceNum := 100000

	if err == nil && len(nonceBytes) >= 3 {
		// 使用前3个字节生成6位随机数
		nonceNum = int(nonceBytes[0])*10000 + int(nonceBytes[1])*100 + int(nonceBytes[2])
		nonceNum = nonceNum % 1000000 // 确保是6位数
	}
	timestamp := time.Now().Format("0601021504")
	randomNum := fmt.Sprintf("%06d", nonceNum)
	incentiveID := fmt.Sprintf("%s_incentive_%s%s",doctype,timestamp, randomNum)
	return incentiveID, nil
}
//查询文档激励信息(全部)
func (c *IncentiveContract) QueryAllDocIncentiveInfo(ctx contractapi.TransactionContextInterface, refID string, doctype string) ([]*typestruct.DocIncentiveInfo, error) {
	query := fmt.Sprintf(`{"selector":{"ref_id":"%s","doctype":"%s"}}`, refID, "incentive")
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("查询文档激励信息失败: %v", err)
	}
	defer resultsIterator.Close()
	docIncentiveInfos := []*typestruct.DocIncentiveInfo{}
	// 遍历查询结果
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Printf("获取下一个查询结果失败: %v", err)
		}

		// 将查询结果反序列化为 DocIncentiveInfo 结构体
		var docIncentiveInfo typestruct.DocIncentiveInfo
		err = json.Unmarshal(queryResponse.Value, &docIncentiveInfo)
		if err != nil {
			fmt.Printf("反序列化文档激励信息失败: %v", err)
		}

		docIncentiveInfos = append(docIncentiveInfos, &docIncentiveInfo)
	}

	return docIncentiveInfos, nil
}
//查询文档激励信息(分页)
func (c *IncentiveContract) QueryDocIncentiveInfoByPage(ctx contractapi.TransactionContextInterface, refID string, doctype string, page int, pageSize int) (*typestruct.IncentiveQueryResult, error) {
	// 构建查询字符串
	queryString := fmt.Sprintf(`{"selector":{"ref_id":"%s","doctype":"%s"}}`, refID, "incentive")

	// 获取总数
	_, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(999999999), "") 
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

	incentiveInfos := []typestruct.DocIncentiveInfo{}

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

		var incentiveInfo typestruct.DocIncentiveInfo
		err = json.Unmarshal(queryResponse.Value, &incentiveInfo)
		if err != nil {
			fmt.Printf("解析查询结果失败: %v", err)
			continue
		}

		incentiveInfos = append(incentiveInfos, incentiveInfo)
		count++

		// 如果达到页面大小，停止
		if len(incentiveInfos) >= pageSize {
			break
		}
	}

	// 构造返回结构
	queryResult := &typestruct.IncentiveQueryResult{
		IncentiveInfos: incentiveInfos,
		Total:         totalCount,
		Page:          page,
		PageSize:      pageSize,
	}

	return queryResult, nil
}


//----------------------------------不同激励机制计算积分----------------------------------
//--------------------------------------积分激励--------------------------------------
func (c *IncentiveContract) CalculateCommonPointIncentive(ctx contractapi.TransactionContextInterface, docIncentiveInfo *typestruct.DocIncentiveInfo) (float64, error) {
	alpha := 0.5
	beta := 0.2
	gamma := 0.3
    //异常入参处理
	if docIncentiveInfo.CommentScore < 1 {
		docIncentiveInfo.CommentScore = 1
	}
	if docIncentiveInfo.Need < 1 {
		docIncentiveInfo.Need = 1
	}
	//综合历史value、评论分数、需求量
	historyValue := float64(docIncentiveInfo.HistoryValue)
	//取log
	logCommentScore := math.Log(float64(docIncentiveInfo.CommentScore))*10
	logNeed := math.Log(float64(docIncentiveInfo.Need))*10
	incentiveValue := alpha * historyValue + beta * logCommentScore + gamma * logNeed
	return math.Round(incentiveValue*100)/100, nil
}

//--------------------------------------三方博弈--------------------------------------

// GameParameters 定义三方博弈所需的参数
type GameParameters struct {
    K1 float64 // 请求者的CTI需求系数
    K2 float64 // 平台的成本系数
    K3 float64 // 提供商的成本系数
    Beta float64 // 提供商服务质量系数
    Theta float64 // 平台服务费比例
    Lambda float64 // 价值权重系数
}

func (c *IncentiveContract) CalculateThreePartyGameIncentive(ctx contractapi.TransactionContextInterface, docIncentiveInfo *typestruct.DocIncentiveInfo) (float64, error) {
	// 初始化固定参数
	params := GameParameters{
		K1: 0.5,
		K2: 0.3,
		K3: 0.4,
		Beta: 0.6,
		Theta: 0.2,
		Lambda: 0.5,
	}
	//异常入参处理
	if docIncentiveInfo.CommentScore < 1 {
		docIncentiveInfo.CommentScore = 1
	}
	if docIncentiveInfo.Need < 1 {
		docIncentiveInfo.Need = 1
	}
	if docIncentiveInfo.TotalUserNum < 1 {
		docIncentiveInfo.TotalUserNum = 1
	}
	// 从docIncentiveInfo获取变量参数
	Y := float64(docIncentiveInfo.TotalUserNum) 
	need := float64(docIncentiveInfo.Need)
	commentScore := docIncentiveInfo.CommentScore
	
	// 计算最优价格
	optimalPrice := (params.Beta*Y*params.K2*params.K3 + params.Beta*need*params.K2*(params.Theta-1) - commentScore*params.K3*params.Theta) / (2*Y*params.K1*params.K2*params.K3)
	
	// 计算服务质量
	Q := (commentScore*params.Beta*optimalPrice) / (params.K3*Y)
	
	// 计算综合价值
	baseValue := docIncentiveInfo.HistoryValue
	rawValue := params.Lambda*baseValue + (1-params.Lambda)*optimalPrice*Q
	
	if baseValue < 1 {
		baseValue = 1
	}
	// 使用 tanh 函数进行规制，让结果在 baseValue 周围波动
	// 将波动范围限制在 baseValue 的 ±30%
	fluctuation := 0.3 * baseValue * math.Tanh((rawValue-baseValue)/baseValue)
	incentiveValue := baseValue + fluctuation
	
	// 确保结果为正数
	if incentiveValue <= 0 {
		incentiveValue = baseValue
	}
	// 返回四舍五入到2位小数的结果
	return math.Round(incentiveValue*100)/100, nil
}

//--------------------------------------演化博弈--------------------------------------	
func (c *IncentiveContract) CalculateEvolutionGameIncentive(ctx contractapi.TransactionContextInterface, docIncentiveInfo *typestruct.DocIncentiveInfo) (float64, error) {
	return docIncentiveInfo.IncentiveValue, nil
}

