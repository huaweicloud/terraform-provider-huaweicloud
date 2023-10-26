package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListInstancesResourceMetricsResponse Response Object
type ListInstancesResourceMetricsResponse struct {

	// 总记录数
	TotalCount *int32 `json:"total_count,omitempty"`

	// 资源监控信息
	ResourceMonitoringInfos *[]ResourceMonitoringInfo `json:"resource_monitoring_infos,omitempty"`
	HttpStatusCode          int                       `json:"-"`
}

func (o ListInstancesResourceMetricsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstancesResourceMetricsResponse struct{}"
	}

	return strings.Join([]string{"ListInstancesResourceMetricsResponse", string(data)}, " ")
}
