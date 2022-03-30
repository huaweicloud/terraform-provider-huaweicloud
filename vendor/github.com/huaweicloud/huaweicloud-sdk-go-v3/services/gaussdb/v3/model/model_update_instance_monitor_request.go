package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateInstanceMonitorRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID

	InstanceId string `json:"instance_id"`

	Body *TaurusModifyInstanceMonitorRequestBody `json:"body,omitempty"`
}

func (o UpdateInstanceMonitorRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceMonitorRequest struct{}"
	}

	return strings.Join([]string{"UpdateInstanceMonitorRequest", string(data)}, " ")
}
