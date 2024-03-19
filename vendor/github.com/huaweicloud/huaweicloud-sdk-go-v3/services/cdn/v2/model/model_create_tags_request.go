package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTagsRequest Request Object
type CreateTagsRequest struct {
	Body *CreateTagsRequestBody `json:"body,omitempty"`
}

func (o CreateTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTagsRequest struct{}"
	}

	return strings.Join([]string{"CreateTagsRequest", string(data)}, " ")
}
