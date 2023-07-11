package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateInstanceTopicRequest Request Object
type CreateInstanceTopicRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *CreateInstanceTopicReq `json:"body,omitempty"`
}

func (o CreateInstanceTopicRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateInstanceTopicRequest struct{}"
	}

	return strings.Join([]string{"CreateInstanceTopicRequest", string(data)}, " ")
}
