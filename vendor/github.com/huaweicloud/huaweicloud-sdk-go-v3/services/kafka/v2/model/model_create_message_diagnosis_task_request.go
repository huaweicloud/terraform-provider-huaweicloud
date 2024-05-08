package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateMessageDiagnosisTaskRequest Request Object
type CreateMessageDiagnosisTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *CreateMessageDiagnosisTaskReq `json:"body,omitempty"`
}

func (o CreateMessageDiagnosisTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMessageDiagnosisTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateMessageDiagnosisTaskRequest", string(data)}, " ")
}
