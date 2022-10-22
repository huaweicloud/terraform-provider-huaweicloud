package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListInstanceConsumerGroupsRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 偏移量，表示从此偏移量开始查询， offset大于等于0。
	Offset *string `json:"offset,omitempty"`

	// 当次查询返回的最大消费组ID个数，默认值为10，取值范围为1~50。
	Limit *string `json:"limit,omitempty"`

	// 消费组名过滤查询，过滤方式为字段包含过滤。
	Group *string `json:"group,omitempty"`
}

func (o ListInstanceConsumerGroupsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstanceConsumerGroupsRequest struct{}"
	}

	return strings.Join([]string{"ListInstanceConsumerGroupsRequest", string(data)}, " ")
}
