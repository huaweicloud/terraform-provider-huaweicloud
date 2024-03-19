package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteTagsRequest Request Object
type BatchDeleteTagsRequest struct {
	Body *DeleteTagsRequestBody `json:"body,omitempty"`
}

func (o BatchDeleteTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteTagsRequest struct{}"
	}

	return strings.Join([]string{"BatchDeleteTagsRequest", string(data)}, " ")
}
