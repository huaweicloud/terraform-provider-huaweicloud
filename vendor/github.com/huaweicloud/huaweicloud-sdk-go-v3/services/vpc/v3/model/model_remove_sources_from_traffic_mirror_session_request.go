package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RemoveSourcesFromTrafficMirrorSessionRequest Request Object
type RemoveSourcesFromTrafficMirrorSessionRequest struct {

	// 流量镜像会话ID
	TrafficMirrorSessionId string `json:"traffic_mirror_session_id"`

	Body *RemoveSourcesFromTrafficMirrorSessionRequestBody `json:"body,omitempty"`
}

func (o RemoveSourcesFromTrafficMirrorSessionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveSourcesFromTrafficMirrorSessionRequest struct{}"
	}

	return strings.Join([]string{"RemoveSourcesFromTrafficMirrorSessionRequest", string(data)}, " ")
}
