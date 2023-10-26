package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTrafficMirrorSessionRequestBody
type UpdateTrafficMirrorSessionRequestBody struct {
	TrafficMirrorSession *UpdateTrafficMirrorSessionOption `json:"traffic_mirror_session"`
}

func (o UpdateTrafficMirrorSessionRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTrafficMirrorSessionRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateTrafficMirrorSessionRequestBody", string(data)}, " ")
}
