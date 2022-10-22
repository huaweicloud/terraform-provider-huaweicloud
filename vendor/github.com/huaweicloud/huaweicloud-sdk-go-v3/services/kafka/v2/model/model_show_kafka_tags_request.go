package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowKafkaTagsRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o ShowKafkaTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowKafkaTagsRequest struct{}"
	}

	return strings.Join([]string{"ShowKafkaTagsRequest", string(data)}, " ")
}
