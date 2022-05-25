package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListLiveStreamsOnlineRequest struct {

	// 推流域名
	PublishDomain string `json:"publish_domain"`

	// 应用名
	App *string `json:"app,omitempty"`

	// 偏移量，表示从此偏移量开始查询， offset大于等于0
	Offset *int32 `json:"offset,omitempty"`

	// 每页记录数，取值范围[1,100]，默认值10
	Limit *int32 `json:"limit,omitempty"`

	// 流名，用于单流查询，携带stream参数时app不能缺省
	Stream *string `json:"stream,omitempty"`
}

func (o ListLiveStreamsOnlineRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListLiveStreamsOnlineRequest struct{}"
	}

	return strings.Join([]string{"ListLiveStreamsOnlineRequest", string(data)}, " ")
}
