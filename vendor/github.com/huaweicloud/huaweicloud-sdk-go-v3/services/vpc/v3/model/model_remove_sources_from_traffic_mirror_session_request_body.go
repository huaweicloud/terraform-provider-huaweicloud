package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RemoveSourcesFromTrafficMirrorSessionRequestBody
type RemoveSourcesFromTrafficMirrorSessionRequestBody struct {
	TrafficMirrorSession *TrafficMirrorSourcesOption `json:"traffic_mirror_session"`
}

func (o RemoveSourcesFromTrafficMirrorSessionRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveSourcesFromTrafficMirrorSessionRequestBody struct{}"
	}

	return strings.Join([]string{"RemoveSourcesFromTrafficMirrorSessionRequestBody", string(data)}, " ")
}
