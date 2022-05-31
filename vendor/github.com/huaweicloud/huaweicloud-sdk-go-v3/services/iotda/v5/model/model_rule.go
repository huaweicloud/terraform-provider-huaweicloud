package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 规则请求结构体
type Rule struct {

	// **参数说明**：规则名称。
	Name string `json:"name"`

	// **参数说明**：规则的描述信息。
	Description *string `json:"description,omitempty"`

	ConditionGroup *ConditionGroup `json:"condition_group"`

	// **参数说明**：规则的动作列表，单个规则最多支持设置10个动作。
	Actions []RuleAction `json:"actions"`

	// **参数说明**：规则的类型。 **取值范围**： - DEVICE_LINKAGE：设备联动。 - DATA_FORWARDING：数据转发。 - EDGE：边缘侧。
	RuleType string `json:"rule_type"`

	// **参数说明**：规则的状态，默认值：active。 **取值范围**： - active：激活。 - inactive：未激活。
	Status *string `json:"status,omitempty"`

	// **参数说明**：资源空间ID。此参数为非必选参数，存在多资源空间的用户需要使用该接口时，建议携带该参数指定创建的规则归属到哪个资源空间下，否则创建的规则将会归属到[默认资源空间](https://support.huaweicloud.com/usermanual-iothub/iot_01_0006.html#section0)下。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	AppId *string `json:"app_id,omitempty"`
}

func (o Rule) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Rule struct{}"
	}

	return strings.Join([]string{"Rule", string(data)}, " ")
}
