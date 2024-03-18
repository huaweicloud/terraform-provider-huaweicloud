package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListConnectorTasksRequest Request Object
type ListConnectorTasksRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 偏移量，表示从此偏移量开始查询，offset大于等于0。
	Offset *int32 `json:"offset,omitempty"`

	// 当次查询返回的最大实例个数，默认值为10，取值范围为1~50。
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListConnectorTasksRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListConnectorTasksRequest struct{}"
	}

	return strings.Join([]string{"ListConnectorTasksRequest", string(data)}, " ")
}
