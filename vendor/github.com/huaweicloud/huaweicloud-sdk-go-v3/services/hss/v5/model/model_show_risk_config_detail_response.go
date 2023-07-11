package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowRiskConfigDetailResponse Response Object
type ShowRiskConfigDetailResponse struct {

	// 风险等级，包含如下:   - Low : 低危   - Medium : 中危   - High : 高危
	Severity *string `json:"severity,omitempty"`

	// 基线类型
	CheckType *string `json:"check_type,omitempty"`

	// 基线描述
	CheckTypeDesc *string `json:"check_type_desc,omitempty"`

	// 检查项总数量
	CheckRuleNum *int32 `json:"check_rule_num,omitempty"`

	// 未通过的检查项数量
	FailedRuleNum *int32 `json:"failed_rule_num,omitempty"`

	// 已通过的检查项数量
	PassedRuleNum *int32 `json:"passed_rule_num,omitempty"`

	// 已忽略的检查项数量
	IgnoredRuleNum *int32 `json:"ignored_rule_num,omitempty"`

	// 受影响的服务器的数量
	HostNum        *int64 `json:"host_num,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ShowRiskConfigDetailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRiskConfigDetailResponse struct{}"
	}

	return strings.Join([]string{"ShowRiskConfigDetailResponse", string(data)}, " ")
}
