
package data_contract
import (
	"fmt"
	"encoding/json"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
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
