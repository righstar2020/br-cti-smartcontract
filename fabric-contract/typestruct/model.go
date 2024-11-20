package typestruct

// 用户结构
type UserInfo struct {
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	PrivateKey    string `json:"public_key"`
	PrivateKeyType string `json:"public_key_type"`
	Value          int    `json:"value"`
	CreateTime     string `json:"create_time"`
}

type UserPointInfo struct {
    UserValue           int                           `json:"user_value"` //用户积分    
    UserCTIMap          map[string][]string           `json:"user_cti_map"` //用户拥有的情报map
    CTIBuyMap           map[string]int                `json:"cti_buy_map"` //用户购买的情报map
    CTISaleMap          map[string]int                `json:"cti_sale_map"` //用户销售的情报map
}

// 情报数据结构定义
type CtiInfo struct {
	CTIID          string   `json:"cti_id"`           // 情报ID
	CTIName        string   `json:"cti_name"`         // 情报名称
	CTIType        int      `json:"cti_type"`         // 情报类型（1-10）10是流量类型的情报
	CTITrafficType int      `json:"cti_traffic_type"` // 流量情报（0：非流量、1：卫星网络、2：5G、3：SDN）
	OpenSource     int      `json:"open_source"`      // 是否开源（0不开源，1开源）
	CreatorUserID  string   `json:"creator_user_id"`  // 创建者ID
	Tags           []string `json:"tags"`             // 情报标签数组
	IOCs           []string   `json:"iocs"`             // 包含的沦陷指标（IP, Port, URL, Hash）
	StixData       string    `json:"stix_data"`        // STIX数据（JSON）可以有多条
	Description string `json:"description"`  // 情报描述
	DataSize    int    `json:"data_size"`    // 数据大小（B）
	DataHash    string `json:"data_hash"`    // 情报数据HASH（sha256）
	IPFSHash    string `json:"ipfs_hash"`    // IPFS地址
	Need        int    `json:"need"`         // 情报需求量
	Value       int    `json:"value"`        // 情报价值（积分）
	CompreValue int    `json:"compre_value"` // 综合价值（积分激励算法定价）
	SaleCount   int    `json:"sale_count"`   //销售数量
	CreateTime  string `json:"create_time"`  // 情报创建时间（由合约生成）
}
