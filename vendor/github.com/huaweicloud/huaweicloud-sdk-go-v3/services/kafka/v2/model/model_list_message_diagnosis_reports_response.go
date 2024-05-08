package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListMessageDiagnosisReportsResponse Response Object
type ListMessageDiagnosisReportsResponse struct {

	// 诊断报告列表
	ReportList *[]KafkaMessageDiagnosisReportInfoEntity `json:"report_list,omitempty"`

	// 诊断报告总数
	TotalNum       *int32 `json:"total_num,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListMessageDiagnosisReportsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMessageDiagnosisReportsResponse struct{}"
	}

	return strings.Join([]string{"ListMessageDiagnosisReportsResponse", string(data)}, " ")
}
