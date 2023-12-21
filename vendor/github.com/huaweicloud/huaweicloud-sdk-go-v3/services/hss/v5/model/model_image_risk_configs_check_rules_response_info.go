package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ImageRiskConfigsCheckRulesResponseInfo 检查项风险信息
type ImageRiskConfigsCheckRulesResponseInfo struct {

	// 风险等级，包含如下:   - Security : 安全   - Low : 低危   - Medium : 中危   - High : 高危
	Severity *string `json:"severity,omitempty"`

	// 基线名称
	CheckName *string `json:"check_name,omitempty"`

	// 基线类型
	CheckType *string `json:"check_type,omitempty"`

	// 标准类型，包含如下:   - cn_standard : 等保合规标准   - hw_standard : 华为标准   - qt_standard : 青腾标准
	Standard *string `json:"standard,omitempty"`

	// 检查项
	CheckRuleName *string `json:"check_rule_name,omitempty"`

	// 检查项ID
	CheckRuleId *string `json:"check_rule_id,omitempty"`

	// 检测结果，包含如下：   - pass    通过   - failed  未通过
	ScanResult *string `json:"scan_result,omitempty"`
}

func (o ImageRiskConfigsCheckRulesResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ImageRiskConfigsCheckRulesResponseInfo struct{}"
	}

	return strings.Join([]string{"ImageRiskConfigsCheckRulesResponseInfo", string(data)}, " ")
}
