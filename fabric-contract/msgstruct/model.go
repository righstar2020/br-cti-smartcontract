package msgstruct
//数据传输结构
//不需要签名的消息
type UserRegisterMsgData struct {
	UserName string `json:"user_name" default:""` //用户名称
	PublicKey string `json:"public_key" default:""` //用户公钥(pem string)
}
type TxMsgRawData struct {
	UserID string `json:"user_id"` //用户ID
	TxData string `json:"tx_data"` //交易数据 base64
	Nonce string `json:"nonce"` //随机数(base64)
	TxSignature string`json:"tx_signature"` //交易签名(Base64 ASN.1 DER)
	NonceSignature string `json:"nonce_signature"` //随机数签名(Base64 ASN.1 DER)
}


//交易数据结构(需要签名的数据)
type TxMsgData struct {
	UserID string `json:"user_id"` //用户ID
	TxData []byte `json:"tx_data"` //交易数据
	Nonce string `json:"nonce"` //随机数(base64)
	TxSignature string `json:"tx_signature"` //交易签名(Base64 ASN.1 DER)
	NonceSignature string `json:"nonce_signature"` //随机数签名(Base64 ASN.1 DER)
}
//----------------------------------情报----------------------------------
//情报交易数据结构
type CtiTxData struct {
	CTIID          string   `json:"cti_id"`           // 情报ID(链上生成)
	CTIHash        string   `json:"cti_hash"`         // 情报HASH(sha256链下生成)
	CTIName        string   `json:"cti_name"`         // 情报名称(可为空)
	CTIType        int      `json:"cti_type"`         // 情报类型（1:恶意流量、2:蜜罐情报、3:僵尸网络、4:应用层攻击、5:开源情报）
	CTITrafficType int      `json:"cti_traffic_type"` // 流量情报类型（0:非流量、1:5G、2:卫星网络、3:SDN）
	OpenSource     int      `json:"open_source"`      // 是否开源（0不开源，1开源）
	CreatorUserID  string   `json:"creator_user_id"`  // 创建者ID(公钥sha256)
	Tags           []string `json:"tags"`             // 情报标签数组
	IOCs           []string `json:"iocs"`             // 包含的沦陷指标（IP, Port, Payload,URL, Hash）
	StatisticInfo  string  `json:"statistic_info"`   // 统计信息(JSON []byte)
	StixData       string   `json:"stix_data"`        // STIX数据（JSON []byte）可以有多条
	StixIPFSHash   string   `json:"stix_ipfs_hash"`   // STIX数据,IPFS地址
	Description    string   `json:"description"`      // 情报描述
	DataSize       int      `json:"data_size"`        // 数据大小（B）
	DataSourceHash string   `json:"data_source_hash"` // 数据源HASH（sha256）
	DataSourceIPFSHash string   `json:"data_source_ipfs_hash"` // 数据源IPFS地址
	Need           int      `json:"need"`             // 情报需求量(销售数量)
	Value          float64      `json:"value"`            // 情报价值（积分）
	CompreValue    float64      `json:"compre_value"`     // 综合价值（积分激励算法定价）
	IncentiveMechanism int `json:"incentive_mechanism"` // 激励机制(1:积分激励、2:三方博弈、3:演化博弈)
}

type PurchaseCtiTxData struct {
	CTIID string `json:"cti_id"` // 情报ID
	UserID string `json:"user_id"` // 用户ID
}

//----------------------------------模型----------------------------------
//模型交易数据结构
type ModelTxData struct {
	ModelID          string   `json:"model_id"`           // 模型ID
	ModelHash        string   `json:"model_hash"`         // 模型hash
	ModelName        string   `json:"model_name"`         // 模型名称
	CreatorUserID    string   `json:"creator_user_id"`    // 模型创建者ID
	ModelDataType    int      `json:"model_data_type"`    // 模型数据类型(1:流量(数据集)、2:情报(文本))
	ModelType        int      `json:"model_type"`         // 模型类型(1:分类模型、2:回归模型、3:聚类模型、4:NLP模型)
	ModelAlgorithm   string   `json:"model_algorithm"`    // 模型算法
	ModelTrainFramework string `json:"model_train_framework"` // 模型训练框架(Scikit-learn、Pytorch、TensorFlow)
	ModelOpenSource  int      `json:"model_open_source"`  // 是否开源
	ModelFeatures    []string `json:"model_features"`     // 模型特征
	ModelTags        []string `json:"model_tags"`         // 模型标签
	ModelDescription string   `json:"model_description"`  // 模型描述
	ModelSize        int      `json:"model_size"`         // 模型大小
	ModelDataSize    int      `json:"model_data_size"`    // 模型训练数据大小
	ModelDataIPFSHash string   `json:"model_data_ipfs_hash"` // 模型训练数据IPFS地址
	IncentiveMechanism int `json:"incentive_mechanism"`   // 激励机制(1:积分激励、2:三方博弈、3:演化博弈)
	Value            float64      `json:"value"`                // 模型价值
	ModelIPFSHash    string   `json:"model_ipfs_hash"`      // 模型IPFS地址
	RefCTIId         string   `json:"ref_cti_id"`           // 关联情报ID(使用哪个情报训练的模型)
}
type PurchaseModelTxData struct {
	ModelID string `json:"model_id"` // 模型ID
	UserID string `json:"user_id"` // 用户ID
}
// 模型算法映射表
var ModelAlgorithms = map[string][]string{
	//分类模型
    "classification": {
        "SVM",                  // 支持向量机
        "RandomForest",         // 随机森林
        "DecisionTree",         // 决策树
        "GaussianNB",          // 高斯朴素贝叶斯
        "KNN",                  // K近邻分类器
        "LogisticRegression",   // 逻辑回归
    },
	//回归模型
    "regression": {
        "LinearRegression",     // 线性回归
        "Ridge",               // 岭回归
        "Lasso",               // Lasso回归
        "SVR",                 // 支持向量回归
        "RandomForestRegressor", // 随机森林回归
    },
	//聚类模型
    "clustering": {
        "KMeans",              // K均值聚类
        "DBSCAN",              // 密度聚类
        "AgglomerativeClustering", // 层次聚类
        "GaussianMixture",     // 高斯混合模型
        "MeanShift",           // 均值漂移
    },
	//NLP模型
    "nlp": {
        "TfidfVectorizer",     // TF-IDF特征提取
        "CountVectorizer",     // 词频特征提取
        "MultinomialNB",       // 多项式朴素贝叶斯
        "BernoulliNB",         // 伯努利朴素贝叶斯
    },
}
//----------------------------------评论----------------------------------
//评论交易数据结构
type CommentTxData struct {
	CommentID string `json:"comment_id"` // 评论ID
	UserID string `json:"user_id"` // 用户ID
	UserLevel float64 `json:"user_level"` // 用户等级(只记录评论发生时用户等级)
	CommentDocType string `json:"comment_doc_type"` // 评论文档类型(cti:情报、model:模型)
	CommentRefID string `json:"comment_ref_id"` // 评论关联ID(情报ID、模型ID)
	CommentScore float64 `json:"comment_score"` // 评论分数
	CommentStatus int `json:"comment_status"` // 评论状态(1:待审核、2:已审核、3:已拒绝)
	CommentContent string `json:"comment_content"` // 评论内容
	CreateTime string `json:"create_time"` // 创建时间
	DocType string `json:"doctype"` // 文档类型
}

//评论审核交易数据结构
type ApproveCommentTxData struct {
	UserID string `json:"user_id"` // 审核用户ID
	CommentID string `json:"comment_id"` // 评论ID
	Status int `json:"status"` // 审核状态(1:通过、2:拒绝)
}
//----------------------------------激励机制----------------------------------
//激励交易数据结构
type IncentiveTxData struct {
	RefID string `json:"ref_id"` // 关联ID(情报ID、模型ID)
	Doctype string `json:"doctype"` // 文档类型(cti、model)
	Mechanism int `json:"mechanism"` // 激励机制(1:积分激励、2:三方博弈、3:演化博弈)
}

