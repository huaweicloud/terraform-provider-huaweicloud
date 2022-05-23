package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteRoutingRuleRequest struct {

	// **参数说明**：实例ID。物理多租下各实例的唯一标识，一般华为云租户无需携带该参数，仅在物理多租场景下从管理面访问API时需要携带该参数。
	InstanceId *string `json:"Instance-Id,omitempty"`

	// **参数说明**：规则条件ID。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	RuleId string `json:"rule_id"`
}

func (o DeleteRoutingRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteRoutingRuleRequest struct{}"
	}

	return strings.Join([]string{"DeleteRoutingRuleRequest", string(data)}, " ")
}
