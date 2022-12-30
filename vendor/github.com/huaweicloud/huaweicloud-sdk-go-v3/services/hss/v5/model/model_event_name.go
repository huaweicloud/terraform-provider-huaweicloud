package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 事件名称
type EventName struct {
}

func (o EventName) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventName struct{}"
	}

	return strings.Join([]string{"EventName", string(data)}, " ")
}
