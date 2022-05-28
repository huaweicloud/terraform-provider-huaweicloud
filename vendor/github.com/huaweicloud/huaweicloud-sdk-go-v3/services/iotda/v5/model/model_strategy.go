package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 配置条件处理策略，用于确定规则是否判断上次数据是否满足条件。当rule_type类型为 “DEVICE_LINKAGE”时，为必选参数。默认为pulse触发类型，事件有效性为永久有效
type Strategy struct {

	// **参数说明**：规则条件触发的判断策略，默认为pulse。 **取值范围**： - pulse：设备上报的数据满足条件则触发，不判断上一次上报的数据。 - reverse：设备上一次上报的数据不满足条件，本次上报的数据满足条件则触发。
	Trigger *string `json:"trigger,omitempty"`

	// **参数说明**：设备数据的有效时间，单位为秒，设备数据的产生时间以上报数据中的eventTime为基准。
	EventValidTime *int32 `json:"event_valid_time,omitempty"`
}

func (o Strategy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Strategy struct{}"
	}

	return strings.Join([]string{"Strategy", string(data)}, " ")
}
