package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTopicReplicaRequest Request Object
type UpdateTopicReplicaRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// Topic名称。
	Topic string `json:"topic"`

	Body *ResetReplicaReq `json:"body,omitempty"`
}

func (o UpdateTopicReplicaRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTopicReplicaRequest struct{}"
	}

	return strings.Join([]string{"UpdateTopicReplicaRequest", string(data)}, " ")
}
