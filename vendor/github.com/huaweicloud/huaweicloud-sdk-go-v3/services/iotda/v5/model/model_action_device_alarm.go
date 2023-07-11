package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ActionDeviceAlarm 上报设备告警消息结构
type ActionDeviceAlarm struct {

	// **参数说明**：告警名称。
	Name string `json:"name"`

	// **参数说明**：告警状态。 **取值范围**： - fault：上报告警。 - recovery：恢复告警。
	AlarmStatus string `json:"alarm_status"`

	// **参数说明**：告警级别。 **取值范围**：warning（警告）、minor（一般）、major（严重）和critical（致命）。
	Severity string `json:"severity"`

	// **参数说明**：告警维度，与告警名称和告警级别组合起来共同标识一条告警，默认不携带该字段为用户维度告警，支持设备维度和资源空间维度告警。 **取值范围**： - device：设备维度。 - app：资源空间维度。
	Dimension *string `json:"dimension,omitempty"`

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
