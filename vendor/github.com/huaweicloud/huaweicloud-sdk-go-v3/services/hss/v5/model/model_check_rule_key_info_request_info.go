package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CheckRuleKeyInfoRequestInfo 检查项key
type CheckRuleKeyInfoRequestInfo struct {

	// 配置检查（基线）的名称，例如SSH、CentOS 7、Windows
	CheckName *string `json:"check_name,omitempty"`

	// 检查项ID，值可以通过这个接口的返回数据获得：/v5/{project_id}/baseline/risk-config/{check_name}/check-rules
	CheckRuleId *string `json:"check_rule_id,omitempty"`

	// 基线标准, 类别包含如下：   - cn_standard#等保合规标准   - hw_standard#云安全实践标准
	Standard *string `json:"standard,omitempty"`

	// 用户键入的检查项修复参数数组
	FixValues *[]CheckRuleFixValuesInfo `json:"fix_values,omitempty"`
}

func (o CheckRuleKeyInfoRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckRuleKeyInfoRequestInfo struct{}"
	}

	return strings.Join([]string{"CheckRuleKeyInfoRequestInfo", string(data)}, " ")
}
