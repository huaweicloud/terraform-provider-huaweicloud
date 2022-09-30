package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type BatchCreateOrDeleteKafkaTagRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *BatchCreateOrDeleteTagReq `json:"body,omitempty"`
}

func (o BatchCreateOrDeleteKafkaTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateOrDeleteKafkaTagRequest struct{}"
	}

	return strings.Join([]string{"BatchCreateOrDeleteKafkaTagRequest", string(data)}, " ")
}
