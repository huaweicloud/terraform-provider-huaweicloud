package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 事件告警上报实体。
type EventList struct {

	// 事件或者告警详情。
	Events []EventModel `json:"events"`
}

func (o EventList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventList struct{}"
	}

	return strings.Join([]string{"EventList", string(data)}, " ")
}
