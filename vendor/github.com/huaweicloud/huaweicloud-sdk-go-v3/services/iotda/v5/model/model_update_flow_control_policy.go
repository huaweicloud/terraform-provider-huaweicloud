package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateFlowControlPolicy 修改数据流转流控策略请求结构体
type UpdateFlowControlPolicy struct {

	// **参数说明**：数据流转流控策略名称。 **取值范围**：长度不超过256，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	PolicyName *string `json:"policy_name,omitempty"`

	// **参数说明**：用户自定义的数据流转流控策略描述。 **取值范围**：长度不超过256，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	Description *string `json:"description,omitempty"`

	// **参数说明**：数据转发流控大小。单位为tps，取值范围为1~1000的整数，默认为1000.
	Limit *int32 `json:"limit,omitempty"`
}

func (o UpdateFlowControlPolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateFlowControlPolicy struct{}"
	}

	return strings.Join([]string{"UpdateFlowControlPolicy", string(data)}, " ")
}
