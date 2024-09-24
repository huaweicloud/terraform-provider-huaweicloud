package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SwitchLogReplayResponse Response Object
type SwitchLogReplayResponse struct {

	// 提示信息
	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SwitchLogReplayResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SwitchLogReplayResponse struct{}"
	}

	return strings.Join([]string{"SwitchLogReplayResponse", string(data)}, " ")
}
