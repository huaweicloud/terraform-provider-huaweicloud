package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddSourcesToTrafficMirrorSessionRequest Request Object
type AddSourcesToTrafficMirrorSessionRequest struct {

	// 流量镜像会话ID
	TrafficMirrorSessionId string `json:"traffic_mirror_session_id"`

	Body *AddSourcesToTrafficMirrorSessionRequestBody `json:"body,omitempty"`
}

func (o AddSourcesToTrafficMirrorSessionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddSourcesToTrafficMirrorSessionRequest struct{}"
	}

	return strings.Join([]string{"AddSourcesToTrafficMirrorSessionRequest", string(data)}, " ")
}
