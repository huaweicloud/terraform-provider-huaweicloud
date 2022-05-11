package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListLogItemsRequest struct {

	// 日志接口调用方式,当值为\"querylogs\"时接口功能为查询日志内容。
	Type string `json:"type"`

	Body *QueryBodyParam `json:"body,omitempty"`
}

func (o ListLogItemsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListLogItemsRequest struct{}"
	}

	return strings.Join([]string{"ListLogItemsRequest", string(data)}, " ")
}
