package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTrafficMirrorSessionRequest Request Object
type UpdateTrafficMirrorSessionRequest struct {

	// 流量镜像会话ID
	TrafficMirrorSessionId string `json:"traffic_mirror_session_id"`

	Body *UpdateTrafficMirrorSessionRequestBody `json:"body,omitempty"`
}

func (o UpdateTrafficMirrorSessionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTrafficMirrorSessionRequest struct{}"
	}

	return strings.Join([]string{"UpdateTrafficMirrorSessionRequest", string(data)}, " ")
}
