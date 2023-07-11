package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteInstanceTopicRequest Request Object
type BatchDeleteInstanceTopicRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *BatchDeleteInstanceTopicReq `json:"body,omitempty"`
}

func (o BatchDeleteInstanceTopicRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteInstanceTopicRequest struct{}"
	}

	return strings.Join([]string{"BatchDeleteInstanceTopicRequest", string(data)}, " ")
}
