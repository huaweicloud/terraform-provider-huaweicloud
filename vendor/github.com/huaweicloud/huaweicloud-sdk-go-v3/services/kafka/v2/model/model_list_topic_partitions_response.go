package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTopicPartitionsResponse Response Object
type ListTopicPartitionsResponse struct {

	// 总条数
	Total *int32 `json:"total,omitempty"`

	// 分区数组
	Partitions     *[]KafkaTopicPartitionResponsePartitions `json:"partitions,omitempty"`
	HttpStatusCode int                                      `json:"-"`
}

func (o ListTopicPartitionsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTopicPartitionsResponse struct{}"
	}

	return strings.Join([]string{"ListTopicPartitionsResponse", string(data)}, " ")
}
