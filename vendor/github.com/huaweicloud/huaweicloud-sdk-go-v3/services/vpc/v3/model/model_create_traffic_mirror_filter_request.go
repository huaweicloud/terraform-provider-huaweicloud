package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTrafficMirrorFilterRequest Request Object
type CreateTrafficMirrorFilterRequest struct {
	Body *CreateTrafficMirrorFilterRequestBody `json:"body,omitempty"`
}

func (o CreateTrafficMirrorFilterRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrafficMirrorFilterRequest struct{}"
	}

	return strings.Join([]string{"CreateTrafficMirrorFilterRequest", string(data)}, " ")
}
