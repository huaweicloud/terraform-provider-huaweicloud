package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTopicAccessPolicyRequest Request Object
type ShowTopicAccessPolicyRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// Topic名称。
	TopicName string `json:"topic_name"`
}

func (o ShowTopicAccessPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTopicAccessPolicyRequest struct{}"
	}

	return strings.Join([]string{"ShowTopicAccessPolicyRequest", string(data)}, " ")
}
