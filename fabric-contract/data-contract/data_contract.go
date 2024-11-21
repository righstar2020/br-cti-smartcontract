
package data_contract
import (
	"fmt"
	"encoding/json"
	"github.com/righstar2020/br-cti-smartcontract/fabric-contract/typestruct"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)
type DataSatisticsInfo struct {
	TotalCtiDataNum int `json:"total_cti_data_num"` // 情报数据总数
	TotalCtiDataSize int `json:"total_cti_data_size"` // 情报数据总大小
	TotalModelDataNum int `json:"total_model_data_num"` // 模型数据总数
	TotalModelDataSize int `json:"total_model_data_size"` // 模型数据总大小
	CTITypeDataNum map[string]int `json:"cti_type_data_num"` // 情报分类型数据数量
	IOCsDataNum map[string]int `json:"iocs_data_num"` // IOCs分类型数据数量
}

type CtiSummaryInfo struct {
	CTIId	 string `json:"cti_id"` // 情报ID
	CTIType int `json:"cti_type"` // 情报类型
	CTITrafficType int `json:"cti_traffic_type"` // 情报流量类型
	IOCsDataNum map[string]int `json:"iocs_data_num"` // IOCs分类型数据数量
	DataCreateTime string `json:"data_create_time"` // 数据创建时间
}

type DataContract struct {
    contractapi.Contract
}

//在这里写统计数据的函数(每次情报上链都会调用这些函数做统计)
//需要对外提供查询接口

// QueryCTIInfo 根据ID查询情报信息
func (c *DataContract) QueryCTIInfo(ctx contractapi.TransactionContextInterface, ctiID string) (*typestruct.CtiInfo, error) {
	ctiInfoJSON, err := ctx.GetStub().GetState(ctiID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if ctiInfoJSON == nil {
		return nil, fmt.Errorf("the cti %s does not exist", ctiID)
	}

	var ctiInfo typestruct.CtiInfo
	err = json.Unmarshal(ctiInfoJSON, &ctiInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cti info: %v", err)
	}

	return &ctiInfo, nil
}