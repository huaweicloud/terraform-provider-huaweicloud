package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowInstanceMonitorExtendRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID

	InstanceId string `json:"instance_id"`
}

func (o ShowInstanceMonitorExtendRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceMonitorExtendRequest struct{}"
	}

	return strings.Join([]string{"ShowInstanceMonitorExtendRequest", string(data)}, " ")
}
