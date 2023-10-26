package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTrafficMirrorFilterResponse Response Object
type CreateTrafficMirrorFilterResponse struct {
	TrafficMirrorFilter *TrafficMirrorFilter `json:"traffic_mirror_filter,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateTrafficMirrorFilterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrafficMirrorFilterResponse struct{}"
	}

	return strings.Join([]string{"CreateTrafficMirrorFilterResponse", string(data)}, " ")
}
