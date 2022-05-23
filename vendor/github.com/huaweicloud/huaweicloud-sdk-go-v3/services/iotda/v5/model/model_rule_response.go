package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 规则响应结构体
type RuleResponse struct {

	// 规则id。
	RuleId *string `json:"rule_id,omitempty"`

	// 规则名称。
	Name string `json:"name"`

	// 规则的描述信息。
	Description *string `json:"description,omitempty"`

	ConditionGroup *ConditionGroup `json:"condition_group"`

	// 规则的动作列表，单个规则最多支持设置10个动作。
	Actions []RuleAction `json:"actions"`

	// 规则的类型 - DEVICE_LINKAGE：设备联动。 - DATA_FORWARDING：数据转发。 - EDGE：边缘侧。
	RuleType string `json:"rule_type"`

	// 规则的状态，默认值：active。 - active：激活。 - inactive：未激活。
	Status *string `json:"status,omitempty"`

	// 资源空间ID。此参数为非必选参数，存在多资源空间的用户需要使用该接口时，建议携带该参数指定创建的规则归属到哪个资源空间下，否则创建的规则将会归属到[默认资源空间](https://support.huaweicloud.com/usermanual-iothub/iot_01_0006.html#section0)下。
	AppId *string `json:"app_id,omitempty"`

	// 归属边缘侧节点设备ID列表。
	EdgeNodeIds *[]string `json:"edge_node_ids,omitempty"`

	// 规则最后更新时间，使用UTC时区，格式：yyyyMMdd'T'HHmmss'Z'。
	LastUpdateTime *string `json:"last_update_time,omitempty"`
}

func (o RuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RuleResponse struct{}"
	}

	return strings.Join([]string{"RuleResponse", string(data)}, " ")
}
