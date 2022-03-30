package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowAuditLogRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID

	InstanceId string `json:"instance_id"`
}

func (o ShowAuditLogRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAuditLogRequest struct{}"
	}

	return strings.Join([]string{"ShowAuditLogRequest", string(data)}, " ")
}
