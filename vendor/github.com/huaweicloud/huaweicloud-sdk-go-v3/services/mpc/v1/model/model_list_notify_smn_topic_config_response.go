package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListNotifySmnTopicConfigResponse Response Object
type ListNotifySmnTopicConfigResponse struct {

	// 事件通知模板信息
	Notifications *[]Notification `json:"notifications,omitempty"`

	// 事件通知模板总数
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListNotifySmnTopicConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListNotifySmnTopicConfigResponse struct{}"
	}

	return strings.Join([]string{"ListNotifySmnTopicConfigResponse", string(data)}, " ")
}
