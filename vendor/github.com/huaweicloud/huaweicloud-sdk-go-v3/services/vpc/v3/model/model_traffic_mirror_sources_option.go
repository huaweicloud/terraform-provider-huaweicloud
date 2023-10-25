package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TrafficMirrorSourcesOption
type TrafficMirrorSourcesOption struct {

	// 功能说明：镜像源ID列表
	TrafficMirrorSources []string `json:"traffic_mirror_sources"`
}

func (o TrafficMirrorSourcesOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TrafficMirrorSourcesOption struct{}"
	}

	return strings.Join([]string{"TrafficMirrorSourcesOption", string(data)}, " ")
}
