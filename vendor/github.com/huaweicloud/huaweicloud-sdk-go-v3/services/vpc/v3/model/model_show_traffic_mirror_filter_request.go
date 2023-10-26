package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTrafficMirrorFilterRequest Request Object
type ShowTrafficMirrorFilterRequest struct {

	// 流量镜像筛选条件ID
	TrafficMirrorFilterId string `json:"traffic_mirror_filter_id"`
}

func (o ShowTrafficMirrorFilterRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTrafficMirrorFilterRequest struct{}"
	}

	return strings.Join([]string{"ShowTrafficMirrorFilterRequest", string(data)}, " ")
}
