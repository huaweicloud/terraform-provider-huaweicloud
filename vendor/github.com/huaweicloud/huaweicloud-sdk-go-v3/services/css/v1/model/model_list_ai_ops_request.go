package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAiOpsRequest Request Object
type ListAiOpsRequest struct {

	// 指定待查询的集群ID。
	ClusterId string `json:"cluster_id"`

	// 分页参数，列表当前分页的数量限制。
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量。 偏移量为一个大于0小于终端节点服务总个数的整数， 表示从偏移量后面的终端节点服务开始查询。
	Start *int32 `json:"start,omitempty"`
}

func (o ListAiOpsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAiOpsRequest struct{}"
	}

	return strings.Join([]string{"ListAiOpsRequest", string(data)}, " ")
}
