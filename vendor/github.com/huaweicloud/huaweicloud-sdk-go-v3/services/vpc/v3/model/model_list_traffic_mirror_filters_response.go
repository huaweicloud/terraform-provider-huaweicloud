package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTrafficMirrorFiltersResponse Response Object
type ListTrafficMirrorFiltersResponse struct {

	// 流量镜像筛选条件对象列表
	TrafficMirrorFilters *[]TrafficMirrorFilter `json:"traffic_mirror_filters,omitempty"`

	PageInfo *PageInfo `json:"page_info,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListTrafficMirrorFiltersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTrafficMirrorFiltersResponse struct{}"
	}

	return strings.Join([]string{"ListTrafficMirrorFiltersResponse", string(data)}, " ")
}
