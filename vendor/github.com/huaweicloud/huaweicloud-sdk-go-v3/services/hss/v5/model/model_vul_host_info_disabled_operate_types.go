package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VulHostInfoDisabledOperateTypes struct {

	// 操作类型 - ignore : 忽略 - not_ignore : 取消忽略 - immediate_repair : 修复 - manual_repair: 人工修复 - verify : 验证 - add_to_whitelist : 加入白名单
	OperateType *string `json:"operate_type,omitempty"`

	// 不可进行操作的原因
	Reason *string `json:"reason,omitempty"`
}

func (o VulHostInfoDisabledOperateTypes) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VulHostInfoDisabledOperateTypes struct{}"
	}

	return strings.Join([]string{"VulHostInfoDisabledOperateTypes", string(data)}, " ")
}
