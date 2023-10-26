package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTrafficMirrorFilterResponse Response Object
type UpdateTrafficMirrorFilterResponse struct {
	TrafficMirrorFilter *TrafficMirrorFilter `json:"traffic_mirror_filter,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateTrafficMirrorFilterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTrafficMirrorFilterResponse struct{}"
	}

	return strings.Join([]string{"UpdateTrafficMirrorFilterResponse", string(data)}, " ")
}
