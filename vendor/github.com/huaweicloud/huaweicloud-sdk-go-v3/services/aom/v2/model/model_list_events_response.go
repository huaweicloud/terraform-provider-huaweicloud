package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListEventsResponse Response Object
type ListEventsResponse struct {

	// 事件或者告警详情。
	Events *[]ListEventModel `json:"events,omitempty"`

	PageInfo       *PageInfo `json:"page_info,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ListEventsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEventsResponse struct{}"
	}

	return strings.Join([]string{"ListEventsResponse", string(data)}, " ")
}
