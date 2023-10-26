package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTrafficMirrorFilterRequestBody
type UpdateTrafficMirrorFilterRequestBody struct {
	TrafficMirrorFilter *UpdateTrafficMirrorFilterOption `json:"traffic_mirror_filter"`
}

func (o UpdateTrafficMirrorFilterRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTrafficMirrorFilterRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateTrafficMirrorFilterRequestBody", string(data)}, " ")
}
