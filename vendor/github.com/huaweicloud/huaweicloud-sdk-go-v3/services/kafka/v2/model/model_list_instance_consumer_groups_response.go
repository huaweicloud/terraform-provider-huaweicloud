package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListInstanceConsumerGroupsResponse Response Object
type ListInstanceConsumerGroupsResponse struct {

	// 所有的消费组。
	Groups *[]GroupInfoSimple `json:"groups,omitempty"`

	// 所有的消费组总数。
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListInstanceConsumerGroupsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstanceConsumerGroupsResponse struct{}"
	}

	return strings.Join([]string{"ListInstanceConsumerGroupsResponse", string(data)}, " ")
}
