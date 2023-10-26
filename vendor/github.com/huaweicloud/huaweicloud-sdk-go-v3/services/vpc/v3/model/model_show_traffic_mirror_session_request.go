package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTrafficMirrorSessionRequest Request Object
type ShowTrafficMirrorSessionRequest struct {

	// 流量镜像会话ID
	TrafficMirrorSessionId string `json:"traffic_mirror_session_id"`
}

func (o ShowTrafficMirrorSessionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTrafficMirrorSessionRequest struct{}"
	}

	return strings.Join([]string{"ShowTrafficMirrorSessionRequest", string(data)}, " ")
}
