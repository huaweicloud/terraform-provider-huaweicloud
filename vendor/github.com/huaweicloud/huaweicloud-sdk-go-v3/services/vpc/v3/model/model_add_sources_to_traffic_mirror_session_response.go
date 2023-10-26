package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddSourcesToTrafficMirrorSessionResponse Response Object
type AddSourcesToTrafficMirrorSessionResponse struct {
	TrafficMirrorSession *TrafficMirrorSession `json:"traffic_mirror_session,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AddSourcesToTrafficMirrorSessionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddSourcesToTrafficMirrorSessionResponse struct{}"
	}

	return strings.Join([]string{"AddSourcesToTrafficMirrorSessionResponse", string(data)}, " ")
}
