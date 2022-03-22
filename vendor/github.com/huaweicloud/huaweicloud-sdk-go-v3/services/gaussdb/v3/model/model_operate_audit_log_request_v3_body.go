package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 开启/关闭审计日志参数体
type OperateAuditLogRequestV3Body struct {
	// 审计日志开关状态。取值：ON|OFF

	SwitchStatus string `json:"switch_status"`
}

func (o OperateAuditLogRequestV3Body) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OperateAuditLogRequestV3Body struct{}"
	}

	return strings.Join([]string{"OperateAuditLogRequestV3Body", string(data)}, " ")
}
