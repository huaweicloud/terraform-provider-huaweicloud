package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 修改规则触发条件请求结构体
type UpdateRuleReq struct {

	// **参数说明**：规则名称。 **取值范围**：长度不超过256，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合
	RuleName *string `json:"rule_name,omitempty"`

	// **参数说明**：用户自定义的规则描述。
	Description *string `json:"description,omitempty"`

	// **参数说明**：用户自定义sql select语句，最大长度500，更新sql时，select跟where必须同时传参，如果需要清除该参数的值，输入空字符串，该参数仅供标准版和企业版用户使用。
	Select *string `json:"select,omitempty"`

	// **参数说明**：用户自定义sql where语句，最大长度500，更新操作时，select跟where必须同时传参，如果需要清除该参数的值，输入空字符串，该参数仅供标准版和企业版用户使用。
	Where *string `json:"where,omitempty"`

	// **参数说明**：修改规则条件的状态是否为激活。
	Active *bool `json:"active,omitempty"`
}

func (o UpdateRuleReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateRuleReq struct{}"
	}

	return strings.Join([]string{"UpdateRuleReq", string(data)}, " ")
}
