package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateResourceTagRequest Request Object
type CreateResourceTagRequest struct {
	Body *ReqCreateTag `json:"body,omitempty"`
}

func (o CreateResourceTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateResourceTagRequest struct{}"
	}

	return strings.Join([]string{"CreateResourceTagRequest", string(data)}, " ")
}
