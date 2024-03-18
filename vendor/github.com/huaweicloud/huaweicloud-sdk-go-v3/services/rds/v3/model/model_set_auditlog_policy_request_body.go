package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SetAuditlogPolicyRequestBody struct {

	// 审计日志保存天数，取值范围0~732。0表示关闭审计日志策略。
	KeepDays int32 `json:"keep_days"`

	// 仅关闭审计日志策略时有效。  - true（默认），表示关闭审计日志策略的同时，延迟删除已有的历史审计日志。 - false，表示关闭审计日志策略的同时，删除已有的历史审计日志。
	ReserveAuditlogs *bool `json:"reserve_auditlogs,omitempty"`

	// 审计记录的操作类型，动态范围。空表示不过滤任何操作类型。
	AuditTypes *[]string `json:"audit_types,omitempty"`
}

func (o SetAuditlogPolicyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetAuditlogPolicyRequestBody struct{}"
	}

	return strings.Join([]string{"SetAuditlogPolicyRequestBody", string(data)}, " ")
}
