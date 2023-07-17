package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListResourceRequest Request Object
type ListResourceRequest struct {
	Body *ResqTagResource `json:"body,omitempty"`
}

func (o ListResourceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListResourceRequest struct{}"
	}

	return strings.Join([]string{"ListResourceRequest", string(data)}, " ")
}
