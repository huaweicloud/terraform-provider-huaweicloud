package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 配置检测检查项参数信息
type CheckRuleFixParamInfo struct {

	// 检查项参数ID
	RuleParamId *int32 `json:"rule_param_id,omitempty"`

	// 检查项参数描述
	RuleDesc *string `json:"rule_desc,omitempty"`

	// 检查项参数默认值
	DefaultValue *int32 `json:"default_value,omitempty"`

	// 检查项参数可取最小值
	RangeMin *int32 `json:"range_min,omitempty"`

	// 检查项参数可取最大值
	RangeMax *int32 `json:"range_max,omitempty"`
}

func (o CheckRuleFixParamInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckRuleFixParamInfo struct{}"
	}

	return strings.Join([]string{"CheckRuleFixParamInfo", string(data)}, " ")
}
