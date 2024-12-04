package data_contract

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
)

type DataContract struct {
	contractapi.Contract
}

//在这里写统计数据的函数(每次情报上链都会调用这些函数做统计)
//需要对外提供查询接口

// QueryLatestCTISummaryInfo 查询最新的num条情报精简信息
func (c *DataContract) QueryLatestCTISummaryInfo(ctx contractapi.TransactionContextInterface, num int) ([]typestruct.CtiSummaryInfo, error) {
	// 添加参数验证
	if num <= 0 {
		return nil, fmt.Errorf("查询数量必须大于0")
	}

	// 构建查询字符串
	queryString := `{"selector":{"cti_id":{"$exists":true}}}`

	// 执行查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("查询CTI数据失败: %v", err)
	}
	defer resultsIterator.Close()

	var ctiSummaryList []typestruct.CtiSummaryInfo
	count := 0

	// 遍历查询结果
	for resultsIterator.HasNext() && count < num {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条CTI数据失败: %v", err)
		}

		var ctiInfo typestruct.CtiInfo
		err = json.Unmarshal(queryResponse.Value, &ctiInfo)
		if err != nil {
			return nil, fmt.Errorf("解析CTI数据失败: %v", err)
		}

		// 构造精简信息
		ctiSummary := typestruct.CtiSummaryInfo{
			CTIId:         ctiInfo.CTIID,
			CTIHash:       ctiInfo.CTIHash,
			CTIType:       ctiInfo.CTIType,
			Tags:          ctiInfo.Tags,
			CreatorUserID: ctiInfo.CreatorUserID,
			CreateTime:    ctiInfo.CreateTime,
		}

		ctiSummaryList = append(ctiSummaryList, ctiSummary)
		count++
	}

	// 处理空结果
	if len(ctiSummaryList) == 0 {
		return []typestruct.CtiSummaryInfo{}, nil // 返回空切片而不是nil
	}

	return ctiSummaryList, nil
}

func (c *DataContract) GetDataStatistics(ctx contractapi.TransactionContextInterface) (string, error) {
	// 初始化统计信息
	stats := typestruct.DataSatisticsInfo{
		TotalCtiDataNum:    0,
		TotalCtiDataSize:   0,
		TotalModelDataNum:  0,
		TotalModelDataSize: 0,
		CTITypeDataNum:     make(map[string]int),
		IOCsDataNum:        make(map[string]int),
	}

	// 查询所有 CTI 信息
	ctiQuery := `{"selector":{"cti_id":{"$exists":true}}}`
	ctiIterator, err := ctx.GetStub().GetQueryResult(ctiQuery)
	if err != nil {
		return "", fmt.Errorf("failed to query CTI data: %v", err)
	}
	defer ctiIterator.Close()

	for ctiIterator.HasNext() {
		queryResponse, err := ctiIterator.Next()
		if err != nil {
			return "", fmt.Errorf("failed to get next CTI data: %v", err)
		}

		var ctiInfo typestruct.CtiInfo
		err = json.Unmarshal(queryResponse.Value, &ctiInfo)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal CTI data: %v", err)
		}

		// 更新统计信息
		stats.TotalCtiDataNum++
		stats.TotalCtiDataSize += ctiInfo.DataSize

		// 更新情报类型统计
		ctiType := fmt.Sprintf("Type_%d", ctiInfo.CTIType)
		stats.CTITypeDataNum[ctiType]++

		// 更新 IOCs 统计
		for _, ioc := range ctiInfo.IOCs {
			stats.IOCsDataNum[ioc]++
		}
	}

	// 查询所有 Model 信息
	modelQuery := `{"selector":{"model_id":{"$exists":true}}}`
	modelIterator, err := ctx.GetStub().GetQueryResult(modelQuery)
	if err != nil {
		return "", fmt.Errorf("failed to query Model data: %v", err)
	}
	defer modelIterator.Close()

	for modelIterator.HasNext() {
		queryResponse, err := modelIterator.Next()
		if err != nil {
			return "", fmt.Errorf("failed to get next Model data: %v", err)
		}

		var modelInfo typestruct.ModelInfo
		err = json.Unmarshal(queryResponse.Value, &modelInfo)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal Model data: %v", err)
		}

		// 更新统计信息
		stats.TotalModelDataNum++
		stats.TotalModelDataSize += modelInfo.ModelDataSize
	}

	// 序列化统计结果为 JSON
	statsBytes, err := json.Marshal(stats)
	if err != nil {
		return "", fmt.Errorf("failed to marshal statistics: %v", err)
	}

	return string(statsBytes), nil
}

// GetUpchainTrend 获取情报上链趋势数据
func (c *DataContract) GetUpchainTrend(ctx contractapi.TransactionContextInterface, timeRange string) (*typestruct.UpchainTrendInfo, error) {
	trendBytes, err := ctx.GetStub().GetState(UPCHAIN_TREND_KEY)
	if err != nil {
		return nil, fmt.Errorf("获取上链趋势数据失败: %v", err)
	}

	var trend typestruct.UpchainTrendInfo
	if trendBytes != nil {
		if err := json.Unmarshal(trendBytes, &trend); err != nil {
			return nil, fmt.Errorf("解析上链趋势数据失败: %v", err)
		}
	} else {
		trend = typestruct.UpchainTrendInfo{
			CTIUpchain:   make(map[string]int),
			ModelUpchain: make(map[string]int),
		}
		trendBytes, err := json.Marshal(trend)
		if err != nil {
			return nil, fmt.Errorf("序列化上链趋势数据失败: %v", err)
		}
		ctx.GetStub().PutState(UPCHAIN_TREND_KEY, trendBytes)
	}

	// 根据timeRange筛选数据
	// TODO: 实现时间范围过滤逻辑

	return &trend, nil
}

// GetAttackTypeRanking 获取攻击类型排行
func (c *DataContract) GetAttackTypeRanking(ctx contractapi.TransactionContextInterface) (*typestruct.AttackRankInfo, error) {
	rankBytes, err := ctx.GetStub().GetState(ATTACK_RANK_KEY)
	if err != nil {
		return nil, fmt.Errorf("获取攻击类型排行数据失败: %v", err)
	}

	var ranking typestruct.AttackRankInfo
	if rankBytes != nil {
		if err := json.Unmarshal(rankBytes, &ranking); err != nil {
			return nil, fmt.Errorf("解析攻击类型排行数据失败: %v", err)
		}
	} else {
		ranking = typestruct.AttackRankInfo{
			Rankings: []typestruct.RankItem{
				{Type: "流量攻击", Count: 0},
				{Type: "恶意软件", Count: 0},
				{Type: "钓鱼攻击", Count: 0},
				{Type: "Botnet", Count: 0},
				{Type: "应用层攻击", Count: 0},
			},
		}
		rankBytes, err := json.Marshal(ranking)
		if err != nil {
			return nil, fmt.Errorf("序列化攻击类型排行数据失败: %v", err)
		}
		ctx.GetStub().PutState(ATTACK_RANK_KEY, rankBytes)
	}

	return &ranking, nil
}

// GetIOCsDistribution 获取IOCs类型分布
func (c *DataContract) GetIOCsDistribution(ctx contractapi.TransactionContextInterface) (*typestruct.IOCsDistributionInfo, error) {
	distBytes, err := ctx.GetStub().GetState(IOCS_DIST_KEY)
	if err != nil {
		return nil, fmt.Errorf("获取IOCs分布数据失败: %v", err)
	}

	var distribution typestruct.IOCsDistributionInfo
	if distBytes != nil {
		if err := json.Unmarshal(distBytes, &distribution); err != nil {
			return nil, fmt.Errorf("解析IOCs分布数据失败: %v", err)
		}
	} else {
		distribution = typestruct.IOCsDistributionInfo{
			Distribution: map[string]float64{
				"IP":      0,
				"Hash":    0,
				"Port":    0,
				"Payload": 0,
				"URL":     0,
				"CVE":     0,
			},
		}
		distBytes, err := json.Marshal(distribution)
		if err != nil {
			return nil, fmt.Errorf("序列化IOCs分布数据失败: %v", err)
		}
		ctx.GetStub().PutState(IOCS_DIST_KEY, distBytes)
	}

	return &distribution, nil
}

// GetGlobalIOCsDistribution 获取全球IOCs地理分布
func (c *DataContract) GetGlobalIOCsDistribution(ctx contractapi.TransactionContextInterface) (*typestruct.GlobalIOCsInfo, error) {
	globalBytes, err := ctx.GetStub().GetState(GLOBAL_IOCS_KEY)
	if err != nil {
		return nil, fmt.Errorf("获取全球IOCs分布数据失败: %v", err)
	}

	var global typestruct.GlobalIOCsInfo
	if globalBytes != nil {
		if err := json.Unmarshal(globalBytes, &global); err != nil {
			return nil, fmt.Errorf("解析全球IOCs分布数据失败: %v", err)
		}
	} else {
		global = typestruct.GlobalIOCsInfo{
			Regions: make(map[string]int),
		}
		globalBytes, err := json.Marshal(global)
		if err != nil {
			return nil, fmt.Errorf("序列化全球IOCs分布数据失败: %v", err)
		}
		ctx.GetStub().PutState(GLOBAL_IOCS_KEY, globalBytes)
	}

	return &global, nil
}

// GetSystemOverview 获取系统概览数据
func (c *DataContract) GetSystemOverview(ctx contractapi.TransactionContextInterface) (*typestruct.SystemOverviewInfo, error) {
	overviewBytes, err := ctx.GetStub().GetState(SYSTEM_OVERVIEW_KEY)
	if err != nil {
		return nil, fmt.Errorf("获取系统概览数据失败: %v", err)
	}

	var overview typestruct.SystemOverviewInfo
	if overviewBytes != nil {
		if err := json.Unmarshal(overviewBytes, &overview); err != nil {
			return nil, fmt.Errorf("解析系统概览数据失败: %v", err)
		}
	} else {
		// 初始化系统概览数据
		overview = typestruct.SystemOverviewInfo{
			BlockHeight:       0,
			TotalTransactions: 0,
			CTIValue:          0,
			CTICount:          0,
			CTITransactions:   0,
			IOCsCount:         0,
			AccountCount:      0,
		}

		// 获取当前区块高度(需要从sdk server中获取)
		blockHeight := 0
		overview.BlockHeight = blockHeight
		overviewBytes, err := json.Marshal(overview)
		if err != nil {
			return nil, fmt.Errorf("序列化系统概览数据失败: %v", err)
		}
		ctx.GetStub().PutState(SYSTEM_OVERVIEW_KEY, overviewBytes)
	}

	return &overview, nil
}

// 定义统计数据的key前缀
const (
	STATS_KEY           = "STATS"
	UPCHAIN_TREND_KEY   = "UPCHAIN_TREND"
	ATTACK_RANK_KEY     = "ATTACK_RANK"
	IOCS_DIST_KEY       = "IOCS_DIST"
	GLOBAL_IOCS_KEY     = "GLOBAL_IOCS"
	SYSTEM_OVERVIEW_KEY = "SYS_OVERVIEW"
)

// UpdateCTIStatistics 更新CTI相关的所有统计数据
func (c *DataContract) UpdateCTIStatistics(ctx contractapi.TransactionContextInterface, ctiInfo *typestruct.CtiInfo) error {
	if err := c.updateBasicStats(ctx, ctiInfo); err != nil {
		return fmt.Errorf("failed to update basic stats: %v", err)
	}

	if err := c.updateUpchainTrend(ctx, "CTI"); err != nil {
		return fmt.Errorf("failed to update upchain trend: %v", err)
	}

	if err := c.updateAttackTypeRanking(ctx, ctiInfo.CTIType); err != nil {
		return fmt.Errorf("failed to update attack ranking: %v", err)
	}

	if err := c.updateIOCsDistribution(ctx, ctiInfo.IOCs); err != nil {
		return fmt.Errorf("failed to update IOCs distribution: %v", err)
	}

	if err := c.updateSystemOverview(ctx, ctiInfo); err != nil {
		return fmt.Errorf("failed to update system overview: %v", err)
	}

	return nil
}

// updateBasicStats 更新基本统计数据
func (c *DataContract) updateBasicStats(ctx contractapi.TransactionContextInterface, ctiInfo *typestruct.CtiInfo) error {
	statsBytes, err := ctx.GetStub().GetState(STATS_KEY)
	if err != nil {
		return err
	}

	var stats typestruct.DataSatisticsInfo
	if statsBytes != nil {
		if err := json.Unmarshal(statsBytes, &stats); err != nil {
			return err
		}
	} else {
		stats = typestruct.DataSatisticsInfo{
			TotalCtiDataNum:    0,
			TotalCtiDataSize:   0,
			TotalModelDataNum:  0,
			TotalModelDataSize: 0,
			CTITypeDataNum:     make(map[string]int),
			IOCsDataNum:        make(map[string]int),
		}
	}

	// 更新统计数据
	stats.TotalCtiDataNum++
	stats.TotalCtiDataSize += ctiInfo.DataSize

	ctiType := fmt.Sprintf("Type_%d", ctiInfo.CTIType)
	if _, ok := stats.CTITypeDataNum[ctiType]; !ok {
		stats.CTITypeDataNum[ctiType] = 0
	}
	stats.CTITypeDataNum[ctiType]++

	for _, ioc := range ctiInfo.IOCs {
		if _, ok := stats.IOCsDataNum[ioc]; !ok {
			stats.IOCsDataNum[ioc] = 0
		}
		stats.IOCsDataNum[ioc]++
	}

	// 保存更新后的统计数据
	statsBytes, err = json.Marshal(stats)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(STATS_KEY, statsBytes)
}

// updateUpchainTrend 更新上链趋势
func (c *DataContract) updateUpchainTrend(ctx contractapi.TransactionContextInterface, upchainType string) error {
	trendBytes, err := ctx.GetStub().GetState(UPCHAIN_TREND_KEY)
	if err != nil {
		return err
	}

	var trend typestruct.UpchainTrendInfo
	if trendBytes != nil {
		if err := json.Unmarshal(trendBytes, &trend); err != nil {
			return err
		}
	} else {
		trend = typestruct.UpchainTrendInfo{
			CTIUpchain:   make(map[string]int),
			ModelUpchain: make(map[string]int),
		}
	}

	today := time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02")
	day_and_hour := time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15")
	if upchainType == "CTI" {
		if trend.CTIUpchain == nil {
			trend.CTIUpchain = make(map[string]int)
		}
		if _, ok := trend.CTIUpchain[today]; !ok {
			trend.CTIUpchain[today] = 0
		}
		trend.CTIUpchain[today]++
		if _, ok := trend.CTIUpchain[day_and_hour]; !ok {
			trend.CTIUpchain[day_and_hour] = 0
		}
		trend.CTIUpchain[day_and_hour]++
	} else if upchainType == "Model" {
		if trend.ModelUpchain == nil {
			trend.ModelUpchain = make(map[string]int)
		}
		if _, ok := trend.ModelUpchain[today]; !ok {
			trend.ModelUpchain[today] = 0
		}
		trend.ModelUpchain[today]++
		if _, ok := trend.ModelUpchain[day_and_hour]; !ok {
			trend.ModelUpchain[day_and_hour] = 0
		}
		trend.ModelUpchain[day_and_hour]++
	}

	trendBytes, err = json.Marshal(trend)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(UPCHAIN_TREND_KEY, trendBytes)
}

// updateAttackTypeRanking 更新攻击类型排行
func (c *DataContract) updateAttackTypeRanking(ctx contractapi.TransactionContextInterface, ctiType int) error {
	rankBytes, err := ctx.GetStub().GetState(ATTACK_RANK_KEY)
	if err != nil {
		return err
	}

	var ranking typestruct.AttackRankInfo
	if rankBytes != nil {
		if err := json.Unmarshal(rankBytes, &ranking); err != nil {
			return err
		}
	} else {
		ranking = typestruct.AttackRankInfo{
			Rankings: []typestruct.RankItem{
				{Type: "TRAFFIC", Count: 0},
				{Type: "HONEYPOT", Count: 0},
				{Type: "BOTNET", Count: 0},
				{Type: "APP_LAYER", Count: 0},
				{Type: "OTHER", Count: 0},
			},
		}
	}

	// 根据CTIType更新对应的攻击类型计数
	attackType := getCTITypeString(ctiType)
	for i := range ranking.Rankings {
		if ranking.Rankings[i].Type == attackType {
			ranking.Rankings[i].Count++
			break
		}
	}

	rankBytes, err = json.Marshal(ranking)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ATTACK_RANK_KEY, rankBytes)
}

// updateIOCsDistribution 更新IOCs分布
func (c *DataContract) updateIOCsDistribution(ctx contractapi.TransactionContextInterface, iocs []string) error {
	distBytes, err := ctx.GetStub().GetState(IOCS_DIST_KEY)
	if err != nil {
		return err
	}

	var distribution typestruct.IOCsDistributionInfo
	if distBytes != nil {
		if err := json.Unmarshal(distBytes, &distribution); err != nil {
			return err
		}
	} else {
		distribution = typestruct.IOCsDistributionInfo{
			TotalCountMap: make(map[string]int),
			Distribution: make(map[string]float64),
		}
	}

	// 更新IOC类型分布
	typeCount := make(map[string]int)
	for _, ioc := range iocs {
		if _, ok := typeCount[ioc]; !ok {
			typeCount[ioc] = 0
		}
		typeCount[ioc]++
	}

	// 累加历史数据和新数据
	for iocType, count := range typeCount {
		if _, ok := distribution.TotalCountMap[iocType]; !ok {
			distribution.TotalCountMap[iocType] = 0
		}
		distribution.TotalCountMap[iocType] += count
	}

	// 计算总数
	var totalCount float64
	for _, count := range distribution.TotalCountMap {
		totalCount += float64(count)
	}

	// 重新计算百分比
	if totalCount > 0 {
		for iocType := range distribution.Distribution {
			distribution.Distribution[iocType] = (float64(distribution.TotalCountMap[iocType]) / totalCount) * 100
		}
	}

	distBytes, err = json.Marshal(distribution)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(IOCS_DIST_KEY, distBytes)
}

// updateSystemOverview 更新系统概览
func (c *DataContract) updateSystemOverview(ctx contractapi.TransactionContextInterface, ctiInfo *typestruct.CtiInfo) error {
	overviewBytes, err := ctx.GetStub().GetState(SYSTEM_OVERVIEW_KEY)
	if err != nil {
		return err
	}

	var overview typestruct.SystemOverviewInfo
	if overviewBytes != nil {
		if err := json.Unmarshal(overviewBytes, &overview); err != nil {
			return err
		}
	} else {
		overview = typestruct.SystemOverviewInfo{
			BlockHeight:       0,
			TotalTransactions: 0,
			CTIValue:          0,
			CTICount:          0,
			IOCsCount:         0,
			AccountCount:      0,
		}
	}

	// 更新系统概览数据
	overview.CTICount++
	overview.CTIValue += ctiInfo.Value
	overview.IOCsCount += len(ctiInfo.IOCs)

	blockHeight := 0
	overview.BlockHeight = blockHeight
	// 获取交易总数
	txID := ctx.GetStub().GetTxID()
	if txID != "" {
		overview.TotalTransactions++
	}

	overviewBytes, err = json.Marshal(overview)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(SYSTEM_OVERVIEW_KEY, overviewBytes)
}

// 辅助函数
func getCTITypeString(ctiType int) string {
	typeMap := map[int]string{
		1: "TRAFFIC",
		2: "HONEYPOT",
		3: "BOTNET",
		4: "APP_LAYER",
	}
	if typeName, ok := typeMap[ctiType]; ok {
		return typeName
	}
	return "OTHER"
}
