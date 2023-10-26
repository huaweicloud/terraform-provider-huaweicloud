package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddSourcesToTrafficMirrorSessionRequestBody
type AddSourcesToTrafficMirrorSessionRequestBody struct {
	TrafficMirrorSession *TrafficMirrorSourcesOption `json:"traffic_mirror_session"`
}

func (o AddSourcesToTrafficMirrorSessionRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddSourcesToTrafficMirrorSessionRequestBody struct{}"
	}

	return strings.Join([]string{"AddSourcesToTrafficMirrorSessionRequestBody", string(data)}, " ")
}
