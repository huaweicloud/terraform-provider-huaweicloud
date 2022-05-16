package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteNotificationRequest struct {

	// 标识关键操作通知id。 批量删除请使用逗号隔开，notification_id=\"xxx1,cccc2\"
	NotificationId string `json:"notification_id"`
}

func (o DeleteNotificationRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteNotificationRequest struct{}"
	}

	return strings.Join([]string{"DeleteNotificationRequest", string(data)}, " ")
}
