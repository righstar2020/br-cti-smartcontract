//----------------------------------信誉合约(用户信誉，评论可信度)----------------------------------
package reputation_contract
import (
	commentContract "github.com/righstar2020/br-cti-smartcontract/fabric-contract/comment-contract"
)

type ReputationContract struct {
	commentContract.CommentContract
}

