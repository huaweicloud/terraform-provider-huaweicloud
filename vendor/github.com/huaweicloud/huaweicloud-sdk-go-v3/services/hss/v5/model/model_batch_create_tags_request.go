package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateTagsRequest Request Object
type BatchCreateTagsRequest struct {

	// 资源类别，hss
	ResourceType string `json:"resource_type"`

	// 资源ID
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
