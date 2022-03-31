package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateNotificationRequest struct {
	Body *CreateNotificationRequestBody `json:"body,omitempty"`
}

func (o CreateNotificationRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateNotificationRequest struct{}"
	}

	return strings.Join([]string{"CreateNotificationRequest", string(data)}, " ")
}
