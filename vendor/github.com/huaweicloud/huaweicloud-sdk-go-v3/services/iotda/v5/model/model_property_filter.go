package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PropertyFilter 设备属性过滤信息，自定义结构。
type PropertyFilter struct {

	// **参数说明**：设备属性的路径信息，格式：service_id/DataProperty，例如门磁状态为“DoorWindow/status”。
	Path string `json:"path"`

	// **参数说明**：数据比较的操作符。 **取值范围**：当前支持的操作符有：>，<，>=，<=，=，in:表示在指定值中匹配和between:表示数值区间。
	Operator string `json:"operator"`

	// **参数说明**：数据比较表达式的右值。与数据比较操作符between联用时，右值表示最小值和最大值，用逗号隔开，如“20,30”表示大于等于20小于30。
	Value *string `json:"value,omitempty"`

	// **参数说明**：当operator为in时该字段必填，使用该字段传递比较表达式右值，上限为20个。
	InValues *[]string `json:"in_values,omitempty"`

	Strategy *Strategy `json:"strategy,omitempty"`
}

func (o PropertyFilter) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PropertyFilter struct{}"
	}

	return strings.Join([]string{"PropertyFilter", string(data)}, " ")
}
