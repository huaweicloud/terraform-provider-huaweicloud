package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NotifiedHistoriesResult 告警发送结果
type NotifiedHistoriesResult struct {

	// 告警流水号
	EventSn *string `json:"event_sn,omitempty"`

	// 通知结果
	Notifications *[]Notifications `json:"notifications,omitempty"`
}

func (o NotifiedHistoriesResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NotifiedHistoriesResult struct{}"
	}

	return strings.Join([]string{"NotifiedHistoriesResult", string(data)}, " ")
}
