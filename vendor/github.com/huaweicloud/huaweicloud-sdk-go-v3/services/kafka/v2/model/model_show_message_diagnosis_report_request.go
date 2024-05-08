package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowMessageDiagnosisReportRequest Request Object
type ShowMessageDiagnosisReportRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 消息积压诊断报告ID
	ReportId string `json:"report_id"`
}

func (o ShowMessageDiagnosisReportRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMessageDiagnosisReportRequest struct{}"
	}

	return strings.Join([]string{"ShowMessageDiagnosisReportRequest", string(data)}, " ")
}
