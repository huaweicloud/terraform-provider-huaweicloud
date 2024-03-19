package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAuditlogPolicyResponse Response Object
type ShowAuditlogPolicyResponse struct {

	// 审计日志保存天数，取值范围0~732。0表示关闭审计日志策略。
	KeepDays *int32 `json:"keep_days,omitempty"`

	// 审计记录的操作类型，动态范围。空表示不过滤任何操作类型。
	AuditTypes *[]string `json:"audit_types,omitempty"`

	// 审计记录的所有操作类型。
	AllAuditLogAction *string `json:"all_audit_log_action,omitempty"`
	HttpStatusCode    int     `json:"-"`
}

func (o ShowAuditlogPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAuditlogPolicyResponse struct{}"
	}

	return strings.Join([]string{"ShowAuditlogPolicyResponse", string(data)}, " ")
}
