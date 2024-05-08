package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteMessageDiagnosisReportsRequest Request Object
type BatchDeleteMessageDiagnosisReportsRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *BatchDeleteMessageDiagnosisReportsReq `json:"body,omitempty"`
}

func (o BatchDeleteMessageDiagnosisReportsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteMessageDiagnosisReportsRequest struct{}"
	}

	return strings.Join([]string{"BatchDeleteMessageDiagnosisReportsRequest", string(data)}, " ")
}
