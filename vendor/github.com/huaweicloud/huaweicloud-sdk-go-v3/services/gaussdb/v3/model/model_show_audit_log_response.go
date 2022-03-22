package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowAuditLogResponse struct {
	// 审计日志开关状态。取值：ON|OFF

	SwitchStatus   *string `json:"switch_status,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowAuditLogResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAuditLogResponse struct{}"
	}

	return strings.Join([]string{"ShowAuditLogResponse", string(data)}, " ")
}
