package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTopicPartitionsRequest Request Object
type ListTopicPartitionsRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	// 主题
	Topic string `json:"topic"`

	// 偏移量，表示查询该偏移量后面的记录
	Offset *int32 `json:"offset,omitempty"`

	// 查询返回记录的数量限制
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListTopicPartitionsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTopicPartitionsRequest struct{}"
	}

	return strings.Join([]string{"ListTopicPartitionsRequest", string(data)}, " ")
}
