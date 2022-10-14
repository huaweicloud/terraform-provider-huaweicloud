package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListInstanceConsumerGroupsResponse struct {

	// 所有的消费组ID
	GroupIds *[]string `json:"group_ids,omitempty"`

	// 所有的消费组总数
	Total *int32 `json:"total,omitempty"`

	// 下一个消费组序号
	NextOffset *int32 `json:"next_offset,omitempty"`

	// 上一个消费组序号
	PreviousOffset *int32 `json:"previous_offset,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListInstanceConsumerGroupsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstanceConsumerGroupsResponse struct{}"
	}

	return strings.Join([]string{"ListInstanceConsumerGroupsResponse", string(data)}, " ")
}
