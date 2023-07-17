package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSnapshotConfigsRequest Request Object
type ListSnapshotConfigsRequest struct {

	// 域名
	Domain string `json:"domain"`

	// 应用名称
	AppName *string `json:"app_name,omitempty"`

	// 每页记录数  取值范围[1,100]  默认值：10
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量，表示从此偏移量开始查询，offset大于等于0
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListSnapshotConfigsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSnapshotConfigsRequest struct{}"
	}

	return strings.Join([]string{"ListSnapshotConfigsRequest", string(data)}, " ")
}
