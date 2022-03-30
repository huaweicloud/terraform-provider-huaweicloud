package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateAuditLogRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID

	InstanceId string `json:"instance_id"`

	Body *OperateAuditLogRequestV3Body `json:"body,omitempty"`
}

func (o UpdateAuditLogRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAuditLogRequest struct{}"
	}

	return strings.Join([]string{"UpdateAuditLogRequest", string(data)}, " ")
}
