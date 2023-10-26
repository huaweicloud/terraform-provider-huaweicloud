package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteTrafficMirrorSessionRequest Request Object
type DeleteTrafficMirrorSessionRequest struct {

	// 流量镜像会话ID
	TrafficMirrorSessionId string `json:"traffic_mirror_session_id"`
}

func (o DeleteTrafficMirrorSessionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTrafficMirrorSessionRequest struct{}"
	}

	return strings.Join([]string{"DeleteTrafficMirrorSessionRequest", string(data)}, " ")
}
