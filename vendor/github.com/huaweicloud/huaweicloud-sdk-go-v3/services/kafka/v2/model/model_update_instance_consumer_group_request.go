package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateInstanceConsumerGroupRequest Request Object
type UpdateInstanceConsumerGroupRequest struct {

	// 消息引擎的类型。
	Engine string `json:"engine"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 消费者组。
	Group string `json:"group"`

	Body *CreateGroupReq `json:"body,omitempty"`
}

func (o UpdateInstanceConsumerGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceConsumerGroupRequest struct{}"
	}

	return strings.Join([]string{"UpdateInstanceConsumerGroupRequest", string(data)}, " ")
}
