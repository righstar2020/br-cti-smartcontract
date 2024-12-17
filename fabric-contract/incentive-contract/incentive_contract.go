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
func (c *IncentiveContract) RegisterDocIncentiveInfo(ctx contractapi.TransactionContextInterface, refID string, doctype string, nonce string) (*typestruct.DocIncentiveInfo, error) {
	var ctiInfo *typestruct.CtiInfo
	var modelInfo *typestruct.ModelInfo
	var err error
	historyValue := 0
	need := 0
	incentiveMechanism := 1
	if doctype == "cti" {
		ctiInfo, err = c.CTIContract.QueryCTIInfo(ctx, refID)
		historyValue = ctiInfo.Value
		need = ctiInfo.Need
		incentiveMechanism = ctiInfo.IncentiveMechanism
	} else if doctype == "model" {
		modelInfo, err = c.ModelContract.QueryModelInfo(ctx, refID)
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
	commentScore := 0.0
	for _, commentInfo := range commentInfos {
		commentScore += commentInfo.CommentScore
	}
	//计算评价分数
	commentScore = commentScore / float64(len(commentInfos))
	//生成激励ID
	incentiveID, err := c.GenerateIncentiveID(ctx, refID, doctype, nonce)
	if err != nil {
		return nil, fmt.Errorf("生成激励ID失败: %v", err)
	}
	docIncentiveInfo := typestruct.DocIncentiveInfo{
		IncentiveID: incentiveID,
		RefID: refID,
		Doctype: doctype,
		HistoryValue: historyValue,
		IncentiveMechanism: incentiveMechanism,
		CommentScore: int(commentScore),//取整
		Need: need,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	//根据激励机制计算激励值
	if incentiveMechanism == INCENTIVE_TYPE_POINT {
		incentiveValue, err := c.CalculateCommonPointIncentive(ctx, &docIncentiveInfo)
		if err != nil {
			return nil, fmt.Errorf("计算激励值失败: %v", err)
		}
		docIncentiveInfo.IncentiveValue = incentiveValue
	} else if incentiveMechanism == INCENTIVE_TYPE_GAME {
		incentiveValue, err := c.CalculateThreePartyGameIncentive(ctx, &docIncentiveInfo)
		if err != nil {
			return nil, fmt.Errorf("计算激励值失败: %v", err)
		}
		docIncentiveInfo.IncentiveValue = incentiveValue
	} else if incentiveMechanism == INCENTIVE_TYPE_EVOLUTION {
		incentiveValue, err := c.CalculateEvolutionGameIncentive(ctx, &docIncentiveInfo)
		if err != nil {
			return nil, fmt.Errorf("计算激励值失败: %v", err)
		}
		docIncentiveInfo.IncentiveValue = incentiveValue
	} else {
		incentiveValue, err := c.CalculateCommonPointIncentive(ctx, &docIncentiveInfo)
		if err != nil {
			return nil, fmt.Errorf("计算激励值失败: %v", err)
		}
		docIncentiveInfo.IncentiveValue = incentiveValue
	}
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

//查询文档激励信息
func (c *IncentiveContract) QueryAllDocIncentiveInfo(ctx contractapi.TransactionContextInterface, refID string, doctype string) ([]*typestruct.DocIncentiveInfo, error) {
	query := fmt.Sprintf(`{"selector":{"ref_id":"%s","doctype":"%s"}}`, refID, doctype)
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
			return nil, fmt.Errorf("获取下一个查询结果失败: %v", err)
		}

		// 将查询结果反序列化为 DocIncentiveInfo 结构体
		var docIncentiveInfo typestruct.DocIncentiveInfo
		err = json.Unmarshal(queryResponse.Value, &docIncentiveInfo)
		if err != nil {
			return nil, fmt.Errorf("反序列化文档激励信息失败: %v", err)
		}

		docIncentiveInfos = append(docIncentiveInfos, &docIncentiveInfo)
	}

	return docIncentiveInfos, nil
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
	incentiveID := fmt.Sprintf("%s_%s_%s_%s", timestamp, randomNum, refID, doctype)
	return incentiveID, nil
}



//----------------------------------不同激励机制计算积分----------------------------------
//--------------------------------------积分激励--------------------------------------
func (c *IncentiveContract) CalculateCommonPointIncentive(ctx contractapi.TransactionContextInterface, docIncentiveInfo *typestruct.DocIncentiveInfo) (int, error) {
	alpha := 0.2
	beta := 0.3
	gamma := 0.5
	//综合历史value、评论分数、需求量
	historyValue := float64(docIncentiveInfo.HistoryValue)
	//取log
	logCommentScore := math.Log(float64(docIncentiveInfo.CommentScore))*10
	logNeed := math.Log(float64(docIncentiveInfo.Need))*10
	incentiveValue := alpha * historyValue + beta * logCommentScore + gamma * logNeed
	return int(incentiveValue), nil
}

//--------------------------------------三方博弈--------------------------------------
func (c *IncentiveContract) CalculateThreePartyGameIncentive(ctx contractapi.TransactionContextInterface, docIncentiveInfo *typestruct.DocIncentiveInfo) (int, error) {
	return docIncentiveInfo.IncentiveValue, nil
}

//--------------------------------------演化博弈--------------------------------------	
func (c *IncentiveContract) CalculateEvolutionGameIncentive(ctx contractapi.TransactionContextInterface, docIncentiveInfo *typestruct.DocIncentiveInfo) (int, error) {
	return docIncentiveInfo.IncentiveValue, nil
}

