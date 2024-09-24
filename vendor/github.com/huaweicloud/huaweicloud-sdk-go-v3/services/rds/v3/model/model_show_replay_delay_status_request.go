package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowReplayDelayStatusRequest Request Object
type ShowReplayDelayStatusRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ShowReplayDelayStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowReplayDelayStatusRequest struct{}"
	}

	return strings.Join([]string{"ShowReplayDelayStatusRequest", string(data)}, " ")
}
