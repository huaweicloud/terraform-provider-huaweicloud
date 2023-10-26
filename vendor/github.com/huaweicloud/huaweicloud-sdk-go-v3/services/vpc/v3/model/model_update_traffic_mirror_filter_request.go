package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTrafficMirrorFilterRequest Request Object
type UpdateTrafficMirrorFilterRequest struct {

	// 流量镜像筛选条件ID
	TrafficMirrorFilterId string `json:"traffic_mirror_filter_id"`

	Body *UpdateTrafficMirrorFilterRequestBody `json:"body,omitempty"`
}

func (o UpdateTrafficMirrorFilterRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTrafficMirrorFilterRequest struct{}"
	}

	return strings.Join([]string{"UpdateTrafficMirrorFilterRequest", string(data)}, " ")
}
