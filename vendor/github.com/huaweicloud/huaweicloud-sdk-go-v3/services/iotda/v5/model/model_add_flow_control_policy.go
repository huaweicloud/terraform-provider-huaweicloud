package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddFlowControlPolicy 新增数据流转流控策略请求结构体
type AddFlowControlPolicy struct {

	// **参数说明**：数据流转流控策略名称。 **取值范围**：长度不超过256，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	PolicyName *string `json:"policy_name,omitempty"`

	// **参数说明**：用户自定义的数据流转流控策略描述。 **取值范围**：长度不超过256，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	Description *string `json:"description,omitempty"`

	// **参数说明**：流控策略作用域. **取值范围**： - USER：租户级流控策略。 - CHANNEL：转发通道级流控策略。 - RULE：转发规则级流控策略。 - ACTION：转发动作级流控策略。
	Scope *string `json:"scope,omitempty"`

	// **参数说明**：流控策略作用域附加值。 scope取值为USER时，可不携带该字段，表示租户级流控。 scope取值为CHANNEL时，**取值范围**：HTTP_FORWARDING、DIS_FORWARDING、OBS_FORWARDING、AMQP_FORWARDING、DMS_KAFKA_FORWARDING。 scope取值为RULE时，该字段为对应的ruleId。 scope取值为ACTION时，该字段为对应的actionId。
	ScopeValue *string `json:"scope_value,omitempty"`

	// **参数说明**：数据转发流控大小。单位为tps，取值范围为1~1000的整数，默认为1000.
	Limit *int32 `json:"limit,omitempty"`
}

func (o AddFlowControlPolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddFlowControlPolicy struct{}"
	}

	return strings.Join([]string{"AddFlowControlPolicy", string(data)}, " ")
}
