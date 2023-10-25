package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTrafficMirrorSessionResponse Response Object
type ShowTrafficMirrorSessionResponse struct {
	TrafficMirrorSession *TrafficMirrorSession `json:"traffic_mirror_session,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowTrafficMirrorSessionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTrafficMirrorSessionResponse struct{}"
	}

	return strings.Join([]string{"ShowTrafficMirrorSessionResponse", string(data)}, " ")
}
