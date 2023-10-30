package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTrafficMirrorSessionResponse Response Object
type UpdateTrafficMirrorSessionResponse struct {
	TrafficMirrorSession *TrafficMirrorSession `json:"traffic_mirror_session,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateTrafficMirrorSessionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTrafficMirrorSessionResponse struct{}"
	}

	return strings.Join([]string{"UpdateTrafficMirrorSessionResponse", string(data)}, " ")
}
