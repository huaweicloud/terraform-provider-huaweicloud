package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowInstanceTopicDetailRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// Topic名称。
	Topic string `json:"topic"`
}

func (o ShowInstanceTopicDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceTopicDetailRequest struct{}"
	}

	return strings.Join([]string{"ShowInstanceTopicDetailRequest", string(data)}, " ")
}
