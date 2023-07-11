package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListNotifyEventResponse Response Object
type ListNotifyEventResponse struct {

	// 事件名称
	EventName *[]string `json:"event_name,omitempty"`

	// 点播通知事件总数
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListNotifyEventResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListNotifyEventResponse struct{}"
	}

	return strings.Join([]string{"ListNotifyEventResponse", string(data)}, " ")
}
