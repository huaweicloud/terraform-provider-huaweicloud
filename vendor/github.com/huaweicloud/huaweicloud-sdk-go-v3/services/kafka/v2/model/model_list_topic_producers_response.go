package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTopicProducersResponse Response Object
type ListTopicProducersResponse struct {

	// 总条数
	Total *int32 `json:"total,omitempty"`

	// 生产者列表
	Producers      *[]KafkaTopicProducerResponseProducers `json:"producers,omitempty"`
	HttpStatusCode int                                    `json:"-"`
}

func (o ListTopicProducersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTopicProducersResponse struct{}"
	}

	return strings.Join([]string{"ListTopicProducersResponse", string(data)}, " ")
}
