package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteTrafficMirrorFilterRequest Request Object
type DeleteTrafficMirrorFilterRequest struct {

	// 流量镜像筛选条件ID
	TrafficMirrorFilterId string `json:"traffic_mirror_filter_id"`
}

func (o DeleteTrafficMirrorFilterRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTrafficMirrorFilterRequest struct{}"
	}

	return strings.Join([]string{"DeleteTrafficMirrorFilterRequest", string(data)}, " ")
}
