package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTrafficMirrorSessionResponse Response Object
type CreateTrafficMirrorSessionResponse struct {
	TrafficMirrorSession *TrafficMirrorSession `json:"traffic_mirror_session,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateTrafficMirrorSessionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrafficMirrorSessionResponse struct{}"
	}

	return strings.Join([]string{"CreateTrafficMirrorSessionResponse", string(data)}, " ")
}
