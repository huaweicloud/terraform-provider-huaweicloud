package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateMessageDiagnosisTaskResponse Response Object
type CreateMessageDiagnosisTaskResponse struct {

	// 诊断报告ID。
	ReportId       *string `json:"report_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateMessageDiagnosisTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMessageDiagnosisTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateMessageDiagnosisTaskResponse", string(data)}, " ")
}
