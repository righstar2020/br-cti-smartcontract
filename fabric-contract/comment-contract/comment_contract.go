package comment_contract

import (
	"encoding/json"
	"fmt"

	"encoding/base64"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/msgstruct"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
	ctiContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/cti-contract"
	modelContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/model-contract"
	userContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/user-contract"
	userPointContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/user-point-contract"
)

// 评论类型
const (
	COMMENT_TYPE_CTI = 1 // 情报评论
	COMMENT_TYPE_MODEL = 2 // 模型评论
)

// 评论状态
const (
	COMMENT_STATUS_PENDING = 1 // 待审核
	COMMENT_STATUS_APPROVED = 2 // 已审核
	COMMENT_STATUS_REJECTED = 3 // 已拒绝
)

// CommentContract 是评论合约的结构体
type CommentContract struct {
	ctiContract.CTIContract
	modelContract.ModelContract
	userContract.UserContract
	userPointContract.UserPointContract
}

// RegisterComment 注册评论信息
func (c *CommentContract) RegisterComment(ctx contractapi.TransactionContextInterface, userID string, nonce string, commentTxData msgstruct.CommentTxData) (*typestruct.CommentInfo, error) {
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
	// 生成评论 ID: 评论类型(1位) + 时间戳(12位,年月日时分) + 随机数(6位)
	commentID := fmt.Sprintf("%s%s%s", commentTxData.CommentDocType, timestamp, randomNum)

	// 创建新的评论信息对象
	newComment := typestruct.CommentInfo{
		CommentID:      commentID,
		UserID:         userID,
		CommentDocType: commentTxData.CommentDocType,
		CommentRefID:   commentTxData.CommentRefID,
		CommentScore:   commentTxData.CommentScore,
		CommentStatus:  COMMENT_STATUS_PENDING, // 默认为待审核状态
		CommentContent: commentTxData.CommentContent,
		CreateTime:     time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"),
		Doctype:        "comment",
	}

	// 序列化为JSON
	commentAsBytes, err := json.Marshal(newComment)
	if err != nil {
		return nil, fmt.Errorf("序列化评论信息失败: %v", err)
	}

	// 存储到账本
	err = ctx.GetStub().PutState(commentID, commentAsBytes)
	if err != nil {
		return nil, fmt.Errorf("存储评论信息失败: %v", err)
	}

	return &newComment, nil
}

// ApproveComment 审核评论
func (c *CommentContract) ApproveComment(ctx contractapi.TransactionContextInterface,userID string, commentID string, status int) error {
	// 获取用户积分信息
	userPointInfo, err := c.QueryUserPointInfo(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户积分信息失败: %v", err)
	}

	//检查用户积分是否大于1000
	if userPointInfo.UserValue < 1000 {
		return fmt.Errorf("初级用户没有审核权限")
	}

	// 获取评论信息
	commentAsBytes, err := ctx.GetStub().GetState(commentID)
	if err != nil {
		return fmt.Errorf("获取评论信息失败: %v", err)
	}
	if commentAsBytes == nil {
		return fmt.Errorf("评论不存在: %s", commentID)
	}

	// 反序列化评论信息
	var comment typestruct.CommentInfo
	err = json.Unmarshal(commentAsBytes, &comment)
	if err != nil {
		return fmt.Errorf("反序列化评论信息失败: %v", err)
	}

	// 更新评论状态
	comment.CommentStatus = status

	// 重新序列化
	commentAsBytes, err = json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("序列化更新后的评论信息失败: %v", err)
	}

	// 更新账本
	err = ctx.GetStub().PutState(commentID, commentAsBytes)
	if err != nil {
		return fmt.Errorf("更新评论状态失败: %v", err)
	}

	return nil
}

// QueryComment 查询特定ID的评论
func (c *CommentContract) QueryComment(ctx contractapi.TransactionContextInterface, commentID string) (*typestruct.CommentInfo, error) {
	// 获取评论信息
	commentAsBytes, err := ctx.GetStub().GetState(commentID)
	if err != nil {
		return nil, fmt.Errorf("获取评论信息失败: %v", err)
	}
	if commentAsBytes == nil {
		return nil, fmt.Errorf("评论不存在: %s", commentID)
	}

	// 反序列化评论信息
	var comment typestruct.CommentInfo
	err = json.Unmarshal(commentAsBytes, &comment)
	if err != nil {
		return nil, fmt.Errorf("反序列化评论信息失败: %v", err)
	}

	return &comment, nil
}
//查询指定RefID的所有评论
func (c *CommentContract) QueryAllCommentsByRefID(ctx contractapi.TransactionContextInterface, refID string) ([]typestruct.CommentInfo, error) {
	queryString := fmt.Sprintf(`{"selector":{"comment_ref_id":"%s","doctype":"comment"}}`, refID)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("执行查询失败: %v", err)
	}
	defer resultsIterator.Close()

	comments := []typestruct.CommentInfo{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Printf("获取下一个查询结果失败: %v", err)
			continue
		}
		var comment typestruct.CommentInfo
		err = json.Unmarshal(queryResponse.Value, &comment)
		if err != nil {
			fmt.Printf("解析评论信息失败: %v", err)
			continue
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
// QueryCommentsByRefIDWithPagination 分页查询特定RefID的评论
func (c *CommentContract) QueryCommentsByRefIDWithPagination(ctx contractapi.TransactionContextInterface, refID string, page int, pageSize int) (*typestruct.CommentQueryResult, error) {
	// 构建查询字符串
	queryString := fmt.Sprintf(`{"selector":{"comment_ref_id":"%s","doctype":"comment"}}`, refID)

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

	comments := []typestruct.CommentInfo{}

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

		var comment typestruct.CommentInfo
		err = json.Unmarshal(queryResponse.Value, &comment)
		if err != nil {
			fmt.Printf("解析评论信息失败: %v", err)
			continue
		}

		comments = append(comments, comment)
		count++

		// 如果达到页面大小，停止
		if len(comments) >= pageSize {
			break
		}
	}

	// 构造返回结构
	queryResult := &typestruct.CommentQueryResult{
		CommentInfos: comments,
		Total:        totalCount,
		Page:         page,
		PageSize:     pageSize,
	}

	return queryResult, nil
}

