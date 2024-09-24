package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CountPortsByTagsRequest Request Object
type CountPortsByTagsRequest struct {
	Body *CountPortsByTagsRequestBody `json:"body,omitempty"`
}

func (o CountPortsByTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CountPortsByTagsRequest struct{}"
	}

	return strings.Join([]string{"CountPortsByTagsRequest", string(data)}, " ")
}
