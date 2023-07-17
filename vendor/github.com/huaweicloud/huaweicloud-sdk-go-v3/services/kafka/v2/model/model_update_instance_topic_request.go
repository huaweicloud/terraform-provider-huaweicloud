package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateInstanceTopicRequest Request Object
type UpdateInstanceTopicRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *UpdateInstanceTopicReq `json:"body,omitempty"`
}

func (o UpdateInstanceTopicRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceTopicRequest struct{}"
	}

	return strings.Join([]string{"UpdateInstanceTopicRequest", string(data)}, " ")
}
