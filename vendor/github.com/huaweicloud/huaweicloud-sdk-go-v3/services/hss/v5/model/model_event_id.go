package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EventId 事件编号
type EventId struct {
}

func (o EventId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventId struct{}"
	}

	return strings.Join([]string{"EventId", string(data)}, " ")
}
