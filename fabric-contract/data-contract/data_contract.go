
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
