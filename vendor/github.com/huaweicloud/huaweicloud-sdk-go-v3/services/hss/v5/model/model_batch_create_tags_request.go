package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateTagsRequest Request Object
type BatchCreateTagsRequest struct {

	// 缺省值:application/json; charset=utf-8
	ContentType *string `json:"Content-Type,omitempty"`

	// 由标签管理服务定义的资源类别，企业主机安全服务调用此接口时资源类别为hss
	ResourceType string `json:"resource_type"`

	// 由标签管理服务定义的资源id，企业主机安全服务调用此接口时资源id为配额ID
	ResourceId string `json:"resource_id"`

	Body *BatchCreateTagsRequestInfo `json:"body,omitempty"`
}

func (o BatchCreateTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateTagsRequest struct{}"
	}

	return strings.Join([]string{"BatchCreateTagsRequest", string(data)}, " ")
}
