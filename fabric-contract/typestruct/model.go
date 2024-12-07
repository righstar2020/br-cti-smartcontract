package typestruct

// 用户结构
type UserInfo struct {
	UserID        string `json:"user_id"`         //用户ID(公钥sha256)
	UserName      string `json:"user_name"`       //用户名
	PublicKey     string `json:"public_key"`      //用户公钥
	PublicKeyType string `json:"public_key_type"` //用户公钥类型
	CreateTime    string `json:"create_time"`     //用户创建时间
}

type UserPointInfo struct {
	UserValue  int            `json:"user_value"`   //用户积分
	UserCTIMap map[string]int `json:"user_cti_map"` //用户拥有的情报map
	CTIBuyMap  map[string]int `json:"cti_buy_map"`  //用户购买的情报map
	CTISaleMap map[string]int `json:"cti_sale_map"` //用户销售的情报map
}
//用户详细信息
type UserDetailInfo struct {
	UserInfo
	UserPointInfo
}


// 情报数据结构定义
type CtiInfo struct {
	CTIID          string   `json:"cti_id"`           // 情报ID(链上生成)
	CTIHash        string   `json:"cti_hash"`         // 情报HASH(sha256链下生成)
	CTIName        string   `json:"cti_name"`         // 情报名称(可为空)
	CTIType        int      `json:"cti_type"`         // 情报类型（1:恶意流量、2:蜜罐情报、3:僵尸网络、4:应用层攻击、5:开源情报）
	CTITrafficType int      `json:"cti_traffic_type"` // 流量情报类型（0:非流量、1:5G、2:卫星网络、3:SDN）
	OpenSource     int      `json:"open_source"`      // 是否开源（0不开源，1开源）
	CreatorUserID  string   `json:"creator_user_id"`  // 创建者ID(公钥sha256)
	Tags           []string `json:"tags"`             // 情报标签数组
	IOCs           []string `json:"iocs"`             // 包含的沦陷指标（IP, Port, Payload,URL, Hash）
	StixData       string   `json:"stix_data"`        // STIX数据（JSON）可以有多条
	StatisticInfo  string   `json:"statistic_info"`   // 统计信息(JSON) 或者IPFS HASH
	Description    string   `json:"description"`      // 情报描述
	DataSize       int      `json:"data_size"`        // 数据大小（B）
	DataHash       string   `json:"data_hash"`        // 情报数据HASH（sha256）
	IPFSHash       string   `json:"ipfs_hash"`        // IPFS地址
	Need           int      `json:"need"`             // 情报需求量(销售数量)
	Value          int      `json:"value"`            // 情报价值（积分）
	CompreValue    int      `json:"compre_value"`     // 综合价值（积分激励算法定价）
	CreateTime     string   `json:"create_time"`      // 情报创建时间（由合约生成）
	Doctype        string    `json:"doctype"`          // 文档类型
}

// 情报查询结果
type CtiQueryResult struct {
	CTIInfos []CtiInfo `json:"cti_infos"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}
//用户拥有的情报(上传+购买的)
type UserOwnCTIInfos struct {
	UploadCTIInfos []CtiInfo `json:"upload_cti_infos"`
	PurchaseCTIInfos []CtiInfo `json:"purchase_cti_infos"`
	Total    int       `json:"total"`
}

type DataSatisticsInfo struct {
	TotalCtiDataNum    int            `json:"total_cti_data_num"`    // 情报数据总数
	TotalCtiDataSize   int            `json:"total_cti_data_size"`   // 情报数据总大小
	TotalModelDataNum  int            `json:"total_model_data_num"`  // 模型数据总数
	TotalModelDataSize int            `json:"total_model_data_size"` // 模型数据总大小
	CTITypeDataNum     map[string]int `json:"cti_type_data_num"`     // 情报分类型数据数量
	IOCsDataNum        map[string]int `json:"iocs_data_num"`         // IOCs分类型数据数量
	ModelTypeDataNum   map[string]int `json:"model_type_data_num"`   // 模型分类型数据数量
}

type CtiSummaryInfo struct {
	CTIId         string   `json:"cti_id"`          // 情报ID（链上生成）
	CTIHash       string   `json:"cti_hash"`        // 情报HASH(sha256链下生成)
	CTIType       int      `json:"cti_type"`        // 情报类型
	Tags          []string `json:"tags"`            // 情报标签数组
	CreatorUserID string   `json:"creator_user_id"` // 创建者ID
	CreateTime    string   `json:"create_time"`     // 创建时间
}

// 模型数据结构定义
type ModelInfo struct {
	ModelID            string   `json:"model_id"`              // 模型ID(链上生成)
	ModelHash          string   `json:"model_hash"`            // 模型hash
	ModelName          string   `json:"model_name"`            // 模型名称
	CreatorUserID      string   `json:"creator_user_id"`       // 模型创建者ID
	ModelDataType      int      `json:"model_data_type"`       // 模型数据类型(1:流量(数据集)、2:情报(文本))
	ModelType          int      `json:"model_type"`            // 模型类型(1:分类模型、2:回归模型、3:聚类模型、4:NLP模型)
	ModelAlgorithm     string   `json:"model_algorithm"`       // 模型算法
	ModelTrainFramework string  `json:"model_train_framework"` // 模型训练框架(Scikit-learn、Pytorch、TensorFlow)
	ModelOpenSource    int      `json:"model_open_source"`     // 是否开源
	ModelFeatures      []string `json:"model_features"`        // 模型特征
	ModelTags          []string `json:"model_tags"`            // 模型标签
	ModelDescription   string   `json:"model_description"`     // 模型描述
	ModelSize          int      `json:"model_size"`            // 模型大小
	ModelDataSize      int      `json:"model_data_size"`       // 模型训练数据大小
	ModelDataIPFSHash  string   `json:"model_data_ipfs_hash"`  // 模型训练数据IPFS地址
	ModelIPFSHash      string   `json:"model_ipfs_hash"`       // 模型IPFS地址
	Value        	   int      `json:"value"`           	   // 模型价值
	RefCTIId           string   `json:"ref_cti_id"`            // 关联情报ID(使用哪个情报训练的模型)
	CreateTime         string   `json:"create_time"`           // 模型创建时间（由合约生成）
}

type ModelQueryResult struct {
	ModelInfos []ModelInfo `json:"model_infos"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}

// TrafficTrendInfo 交易趋势信息
type TrafficTrendInfo struct {
	CTITraffic   map[string]int `json:"cti_traffic"`
	ModelTraffic map[string]int `json:"model_traffic"`
}

// RankItem 排名项
type RankItem struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// AttackRankInfo 攻击类型排名信息
type AttackRankInfo struct {
	Rankings []RankItem `json:"rankings"`
}

// IOCsDistributionInfo IOCs分布信息
type IOCsDistributionInfo struct {
	TotalCountMap map[string]int     `json:"total_count_map"`
	Distribution  map[string]float64 `json:"distribution"`
}

// GlobalIOCsInfo 全球IOCs分布信息
type GlobalIOCsInfo struct {
	Regions map[string]int `json:"regions"`
}

// SystemOverviewInfo 系统概览信息
type SystemOverviewInfo struct {
	BlockHeight       int `json:"block_height"`
	TotalTransactions int `json:"total_transactions"`
	CTIValue          int `json:"cti_value"`
	CTICount          int `json:"cti_count"`
	CTITransactions   int `json:"cti_transactions"`
	ModelValue        int `json:"model_value"`
	ModelCount        int `json:"model_count"`
	ModelTransactions int `json:"model_transactions"`
	IOCsCount         int `json:"iocs_count"`
	AccountCount      int `json:"account_count"`
}

// UpchainTrendInfo 上链趋势信息
type UpchainTrendInfo struct {
	CTIUpchain   map[string]int `json:"cti_upchain"`   // 情报上链趋势
	ModelUpchain map[string]int `json:"model_upchain"` // 模型上链趋势
}
