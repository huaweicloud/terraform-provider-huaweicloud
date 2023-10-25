package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTrafficMirrorFilterResponse Response Object
type ShowTrafficMirrorFilterResponse struct {
	TrafficMirrorFilter *TrafficMirrorFilter `json:"traffic_mirror_filter,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowTrafficMirrorFilterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTrafficMirrorFilterResponse struct{}"
	}

	return strings.Join([]string{"ShowTrafficMirrorFilterResponse", string(data)}, " ")
}
