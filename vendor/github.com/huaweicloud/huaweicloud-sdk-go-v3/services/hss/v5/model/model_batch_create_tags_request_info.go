package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateTagsRequestInfo 批量添加标签的请求体
type BatchCreateTagsRequestInfo struct {

	// 标签对象列表
	Tags []ResourceTagInfo `json:"tags"`
}

func (o BatchCreateTagsRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateTagsRequestInfo struct{}"
	}

	return strings.Join([]string{"BatchCreateTagsRequestInfo", string(data)}, " ")
}
