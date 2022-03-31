package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateNotificationRequest struct {
	Body *UpdateNotificationRequestBody `json:"body,omitempty"`
}

func (o UpdateNotificationRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateNotificationRequest struct{}"
	}

	return strings.Join([]string{"UpdateNotificationRequest", string(data)}, " ")
}
