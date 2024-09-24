package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SwitchLogReplayRequest Request Object
type SwitchLogReplayRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage string `json:"X-Language"`

	Body *SwitchLogReplayRequestBody `json:"body,omitempty"`
}

func (o SwitchLogReplayRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SwitchLogReplayRequest struct{}"
	}

	return strings.Join([]string{"SwitchLogReplayRequest", string(data)}, " ")
}
