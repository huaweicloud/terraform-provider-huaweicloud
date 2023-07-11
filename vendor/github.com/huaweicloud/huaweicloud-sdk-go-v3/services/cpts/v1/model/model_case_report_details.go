package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CaseReportDetails struct {

	// 用例下所有事务的基本性能数据信息
	CustomTransactions *[]CaseReportDetail `json:"customTransactions,omitempty"`

	// 用例下所有aw的基本性能数据信息
	DetailDatas *[]CaseReportDetail `json:"detailDatas,omitempty"`

	Performance *CaseReportDetail `json:"performance,omitempty"`
}

func (o CaseReportDetails) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CaseReportDetails struct{}"
	}

	return strings.Join([]string{"CaseReportDetails", string(data)}, " ")
}
