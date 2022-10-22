package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateInstanceAutoCreateTopicReq struct {

	// 是否开启自动创建topic功能。
	EnableAutoTopic bool `json:"enable_auto_topic"`
}

func (o UpdateInstanceAutoCreateTopicReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceAutoCreateTopicReq struct{}"
	}

	return strings.Join([]string{"UpdateInstanceAutoCreateTopicReq", string(data)}, " ")
}
