package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTrafficMirrorFilterRequestBody
type CreateTrafficMirrorFilterRequestBody struct {
	TrafficMirrorFilter *CreateTrafficMirrorFilterOption `json:"traffic_mirror_filter"`
}

func (o CreateTrafficMirrorFilterRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrafficMirrorFilterRequestBody struct{}"
	}

	return strings.Join([]string{"CreateTrafficMirrorFilterRequestBody", string(data)}, " ")
}
