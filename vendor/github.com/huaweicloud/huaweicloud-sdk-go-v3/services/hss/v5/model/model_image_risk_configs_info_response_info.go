package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ImageRiskConfigsInfoResponseInfo 配置检测结果信息
type ImageRiskConfigsInfoResponseInfo struct {

	// 风险等级，包含如下:   - Security : 安全   - Low : 低危   - Medium : 中危   - High : 高危
	Severity *string `json:"severity,omitempty"`

	// 基线名称
	CheckName *string `json:"check_name,omitempty"`

	// 基线类型
	CheckType *string `json:"check_type,omitempty"`

	// 标准类型，包含如下:   - cn_standard : 等保合规标准   - hw_standard : 华为标准   - qt_standard : 青腾标准
	Standard *string `json:"standard,omitempty"`

	// 检查项数量
	CheckRuleNum *int32 `json:"check_rule_num,omitempty"`

	// 风险项数量
	FailedRuleNum *int32 `json:"failed_rule_num,omitempty"`

	// 基线描述信息
	CheckTypeDesc *string `json:"check_type_desc,omitempty"`
}

func (o ImageRiskConfigsInfoResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ImageRiskConfigsInfoResponseInfo struct{}"
	}

	return strings.Join([]string{"ImageRiskConfigsInfoResponseInfo", string(data)}, " ")
}
