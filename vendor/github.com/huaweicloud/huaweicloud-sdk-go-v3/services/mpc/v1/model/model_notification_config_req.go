package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NotificationConfigReq struct {

	// 事件通知模板信息
	Notifications []Notification `json:"notifications"`
}

func (o NotificationConfigReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NotificationConfigReq struct{}"
	}

	return strings.Join([]string{"NotificationConfigReq", string(data)}, " ")
}
