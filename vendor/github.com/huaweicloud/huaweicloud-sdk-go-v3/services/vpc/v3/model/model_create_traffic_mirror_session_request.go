package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTrafficMirrorSessionRequest Request Object
type CreateTrafficMirrorSessionRequest struct {
	Body *CreateTrafficMirrorSessionRequestBody `json:"body,omitempty"`
}

func (o CreateTrafficMirrorSessionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrafficMirrorSessionRequest struct{}"
	}

	return strings.Join([]string{"CreateTrafficMirrorSessionRequest", string(data)}, " ")
}
