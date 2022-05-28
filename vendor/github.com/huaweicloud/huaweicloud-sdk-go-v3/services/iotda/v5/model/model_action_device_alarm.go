package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 上报设备告警消息结构
type ActionDeviceAlarm struct {

	// **参数说明**：告警名称。
	Name string `json:"name"`

	// **参数说明**：告警状态。 **取值范围**： - fault：上报告警。 - recovery：恢复告警。
	AlarmStatus string `json:"alarm_status"`

	// **参数说明**：告警级别。 **取值范围**： - warning：警告。 - minor：一般。 - major：严重。 - critical：致命。
	Severity string `json:"severity"`

	// **参数说明**：告警的描述信息。
	Description *string `json:"description,omitempty"`
}

func (o ActionDeviceAlarm) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ActionDeviceAlarm struct{}"
	}

	return strings.Join([]string{"ActionDeviceAlarm", string(data)}, " ")
}
