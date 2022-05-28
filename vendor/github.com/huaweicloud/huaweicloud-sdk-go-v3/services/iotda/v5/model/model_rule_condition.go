package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 规则条件结构体
type RuleCondition struct {

	// **参数说明**：规则条件的类型。 **取值范围**： - DEVICE_DATA：设备数据类型条件。 - SIMPLE_TIMER：简单定时类型条件。 - DAILY_TIMER：每日定时类型条件。 - DEVICE_STATUS：设备状态类型条件。 - DEVICE_LIFE_CYCLE：设备生命周期类型条件。 - DEVICE_MESSAGE：设备消息条件。 - MESSAGE_RESULTS：下行消息结果条件。
	Type string `json:"type"`

	DevicePropertyCondition *DeviceDataCondition `json:"device_property_condition,omitempty"`

	SimpleTimerCondition *SimpleTimerType `json:"simple_timer_condition,omitempty"`

	DailyTimerCondition *DailyTimerType `json:"daily_timer_condition,omitempty"`
}

func (o RuleCondition) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RuleCondition struct{}"
	}

	return strings.Join([]string{"RuleCondition", string(data)}, " ")
}
