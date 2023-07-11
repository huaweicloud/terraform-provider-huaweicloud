package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCheckRuleDetailRequest Request Object
type ShowCheckRuleDetailRequest struct {

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 基线名称
	CheckName string `json:"check_name"`

	// 基线类型
	CheckType string `json:"check_type"`

	// 检查项ID
	CheckRuleId string `json:"check_rule_id"`

	// 标准类型，包含如下:   - cn_standard : 等保合规标准   - hw_standard : 华为标准   - qt_standard : 青腾标准
	Standard string `json:"standard"`

	// 主机ID
	HostId *string `json:"host_id,omitempty"`
}

func (o ShowCheckRuleDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCheckRuleDetailRequest struct{}"
	}

	return strings.Join([]string{"ShowCheckRuleDetailRequest", string(data)}, " ")
}
