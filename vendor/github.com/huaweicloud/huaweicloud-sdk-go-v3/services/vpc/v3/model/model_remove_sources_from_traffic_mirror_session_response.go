package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RemoveSourcesFromTrafficMirrorSessionResponse Response Object
type RemoveSourcesFromTrafficMirrorSessionResponse struct {
	TrafficMirrorSession *TrafficMirrorSession `json:"traffic_mirror_session,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o RemoveSourcesFromTrafficMirrorSessionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveSourcesFromTrafficMirrorSessionResponse struct{}"
	}

	return strings.Join([]string{"RemoveSourcesFromTrafficMirrorSessionResponse", string(data)}, " ")
}
