package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type NotifySmnTopicConfigRequest struct {
	Body *NotificationConfigReq `json:"body,omitempty"`
}

func (o NotifySmnTopicConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NotifySmnTopicConfigRequest struct{}"
	}

	return strings.Join([]string{"NotifySmnTopicConfigRequest", string(data)}, " ")
}
