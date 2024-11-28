package data_contract
import (
	"fmt"
	"encoding/json"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
	"strings"
	"net"
)


type DataContract struct {
    contractapi.Contract
}

//在这里写统计数据的函数(每次情报上链都会调用这些函数做统计)
//需要对外提供查询接口

// QueryCTIInfo 根据ID查询情报信息
func (c *DataContract) QueryCTISummaryInfoByCTIID(ctx contractapi.TransactionContextInterface, ctiID string) (*typestruct.CtiSummaryInfo, error) {
	// 从账本中查询指定 CTIID 的 CtiInfo
	ctiAsBytes, err := ctx.GetStub().GetState(ctiID)
	if err != nil {
		return nil, fmt.Errorf("failed to get CTI info: %v", err)
	}
	if ctiAsBytes == nil {
		return nil, fmt.Errorf("CTI with ID %s does not exist", ctiID)
	}

	// 反序列化为 CtiInfo 结构体
	var ctiInfo typestruct.CtiInfo
	err = json.Unmarshal(ctiAsBytes, &ctiInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal CTI info: %v", err)
	}

	// 初始化 IOCsDataNum
	iocsDataNum := make(map[string]int)
	for _, ioc := range ctiInfo.IOCs {
		iocsDataNum[ioc]++ // 按类型统计
	}

	// 构造 CtiSummaryInfo
	ctiSummary := &typestruct.CtiSummaryInfo{
		CTIId:         ctiInfo.CTIID,
		CTIType:       ctiInfo.CTIType,
		CTITrafficType: ctiInfo.CTITrafficType,
		IOCsDataNum:   iocsDataNum,
		DataCreateTime: ctiInfo.CreateTime,
	}

	return ctiSummary, nil
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

// GetCTITrafficTrend 获取情报交易趋势数据
func (c *DataContract) GetCTITrafficTrend(ctx contractapi.TransactionContextInterface, timeRange string) (*typestruct.TrafficTrendInfo, error) {
    trend := &typestruct.TrafficTrendInfo{
        CTITraffic: make(map[string]int),    // 情报交易量
        ModelTraffic: make(map[string]int),   // 模型交易量
    }
    
    // 根据timeRange查询指定时间范围内的交易数据
    // 返回按时间分组的交易量统计
    return trend, nil
}

// GetAttackTypeRanking 获取攻击类型排行
func (c *DataContract) GetAttackTypeRanking(ctx contractapi.TransactionContextInterface) (*typestruct.AttackRankInfo, error) {
    ranking := &typestruct.AttackRankInfo{
        Rankings: []typestruct.RankItem{
            {Type: "流量攻击", Count: 0},
            {Type: "恶意软件", Count: 0},
            {Type: "钓鱼攻击", Count: 0},
            {Type: "Botnet", Count: 0},
            {Type: "应用层攻击", Count: 0},
        },
    }
    
    // 统计各类型攻击数量并排序
    return ranking, nil
}

// GetIOCsDistribution 获取IOCs类型分布
func (c *DataContract) GetIOCsDistribution(ctx contractapi.TransactionContextInterface) (*typestruct.IOCsDistributionInfo, error) {
    distribution := &typestruct.IOCsDistributionInfo{
        Distribution: map[string]float64{
            "IP": 0,
            "Hash": 0,
            "Domain": 0,
            "URL": 0,
            "CVE": 0,
        },
    }
    
    // 统计各类型IOC的占比
    return distribution, nil
}

// GetGlobalIOCsDistribution 获取全球IOCs地理分布
func (c *DataContract) GetGlobalIOCsDistribution(ctx contractapi.TransactionContextInterface) (*typestruct.GlobalIOCsInfo, error) {
    global := &typestruct.GlobalIOCsInfo{
        Regions: make(map[string]int),  // 按地区统计IOC数量
    }
    
    // 统计各地区的IOC数量
    return global, nil
}

// GetSystemOverview 获取系统概览数据
func (c *DataContract) GetSystemOverview(ctx contractapi.TransactionContextInterface) (*typestruct.SystemOverviewInfo, error) {
    overview := &typestruct.SystemOverviewInfo{
        BlockHeight: 0,        // 区块高度
        TotalTransactions: 0,  // 区块交易总数
        CTIValue: 0,          // 情报价值总分
        CTICount: 0,          // 情报数量
        CTITransactions: 0,    // 情报交易数
        IOCsCount: 0,         // IOCs数量
        AccountCount: 0,       // 账户数量
    }
    
    // 获取系统整体统计数据
    return overview, nil
}

// 定义统计数据的key前缀
const (
    STATS_KEY = "STATS"
    TRAFFIC_KEY = "TRAFFIC"
    ATTACK_RANK_KEY = "ATTACK_RANK" 
    IOCS_DIST_KEY = "IOCS_DIST"
    GLOBAL_IOCS_KEY = "GLOBAL_IOCS"
    SYSTEM_OVERVIEW_KEY = "SYS_OVERVIEW"
)

// UpdateCTIStatistics 更新CTI相关的所有统计数据
func (c *DataContract) UpdateCTIStatistics(ctx contractapi.TransactionContextInterface, ctiInfo *typestruct.CtiInfo) error {
    if err := c.updateBasicStats(ctx, ctiInfo); err != nil {
        return fmt.Errorf("failed to update basic stats: %v", err)
    }
    
    if err := c.updateTrafficTrend(ctx, "CTI"); err != nil {
        return fmt.Errorf("failed to update traffic trend: %v", err)
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
    }

    // 更新统计数据
    stats.TotalCtiDataNum++
    stats.TotalCtiDataSize += ctiInfo.DataSize
    
    ctiType := fmt.Sprintf("Type_%d", ctiInfo.CTIType)
    stats.CTITypeDataNum[ctiType]++
    
    for _, ioc := range ctiInfo.IOCs {
        stats.IOCsDataNum[ioc]++
    }

    // 保存更新后的统计数据
    statsBytes, err = json.Marshal(stats)
    if err != nil {
        return err
    }
    
    return ctx.GetStub().PutState(STATS_KEY, statsBytes)
}

// updateTrafficTrend 更新流量趋势
func (c *DataContract) updateTrafficTrend(ctx contractapi.TransactionContextInterface, trafficType string) error {
    trendBytes, err := ctx.GetStub().GetState(TRAFFIC_KEY)
    if err != nil {
        return err
    }

    var trend typestruct.TrafficTrendInfo
    if trendBytes != nil {
        if err := json.Unmarshal(trendBytes, &trend); err != nil {
            return err
        }
    }

    today := time.Now().Format("2006-01-02")
    if trafficType == "CTI" {
        if trend.CTITraffic == nil {
            trend.CTITraffic = make(map[string]int)
        }
        trend.CTITraffic[today]++
    } else if trafficType == "Model" {
        if trend.ModelTraffic == nil {
            trend.ModelTraffic = make(map[string]int)
        }
        trend.ModelTraffic[today]++
    }

    trendBytes, err = json.Marshal(trend)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(TRAFFIC_KEY, trendBytes)
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
    }

    // 根据CTIType更新对应的攻击类型计数
    attackType := getAttackTypeString(ctiType)
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
    }

    // 更新IOC类型分布
    total := float64(len(iocs))
    typeCount := make(map[string]int)
    for _, ioc := range iocs {
        iocType := getIOCType(ioc)
        typeCount[iocType]++
    }

    // 计算百分比
    for iocType, count := range typeCount {
        distribution.Distribution[iocType] = float64(count) / total * 100
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
    }

    // 更新系统概览数据
    overview.CTICount++
    overview.CTIValue += ctiInfo.Value
    overview.IOCsCount += len(ctiInfo.IOCs)
    // 获取当前区块高度
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
func getAttackTypeString(ctiType int) string {
    typeMap := map[int]string{
        1: "流量攻击",
        2: "恶意软件",
        3: "钓鱼攻击",
        4: "Botnet",
        5: "应用层攻击",
    }
    if typeName, ok := typeMap[ctiType]; ok {
        return typeName
    }
    return "其他"
}

func getIOCType(ioc string) string {
    // 这里需要根据实际的IOC格式来判断类型
    // 这只是一个示例实现
    if strings.Contains(ioc, ".") {
        if net.ParseIP(ioc) != nil {
            return "IP"
        }
        return "Domain"
    }
    if strings.Contains(ioc, "://") {
        return "URL"
    }
    if strings.HasPrefix(ioc, "CVE-") {
        return "CVE"
    }
    if len(ioc) == 32 || len(ioc) == 40 || len(ioc) == 64 {
        return "Hash"
    }
    return "Other"
}

