package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTopicProducersRequest Request Object
type ListTopicProducersRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	// 主题
	Topic string `json:"topic"`

	// 偏移量，表示查询该偏移量后面的记录
	Offset *int32 `json:"offset,omitempty"`

	// 查询返回记录的数量限制
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListTopicProducersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTopicProducersRequest struct{}"
	}

	return strings.Join([]string{"ListTopicProducersRequest", string(data)}, " ")
}
