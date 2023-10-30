package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTrafficMirrorSessionRequestBody
type CreateTrafficMirrorSessionRequestBody struct {
	TrafficMirrorSession *CreateTrafficMirrorSessionOption `json:"traffic_mirror_session"`
}

func (o CreateTrafficMirrorSessionRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrafficMirrorSessionRequestBody struct{}"
	}

	return strings.Join([]string{"CreateTrafficMirrorSessionRequestBody", string(data)}, " ")
}
