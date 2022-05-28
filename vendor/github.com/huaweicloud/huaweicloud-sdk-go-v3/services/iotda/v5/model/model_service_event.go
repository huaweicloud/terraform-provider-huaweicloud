package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 事件服务对象。
type ServiceEvent struct {

	// **参数说明**：设备事件类型。注：设备服务内不允许重复。 **取值范围**：长度不超过32，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	EventType string `json:"event_type"`

	// **参数说明**：设备事件的参数列表。
	Paras *[]ServiceCommandPara `json:"paras,omitempty"`
}

func (o ServiceEvent) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ServiceEvent struct{}"
	}

	return strings.Join([]string{"ServiceEvent", string(data)}, " ")
}
