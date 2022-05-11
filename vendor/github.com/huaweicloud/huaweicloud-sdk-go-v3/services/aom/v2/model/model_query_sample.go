package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 查询参数集。
type QuerySample struct {

	// 时间序列命名空间。 取值范围： PAAS.CONTAINER、PAAS.NODE、PAAS.SLA、PAAS.AGGR、CUSTOMMETRICS等； PAAS.CONTAINER：应用时间序列； PAAS.NODE：节点时间序列； PAAS.SLA：SLA时间序列； PAAS.AGGR：集群时间序列； CUSTOMMETRICS：自定义时间序列。
	Namespace string `json:"namespace"`

	// 时间序列维度列表 可通过/v2/{project_id}/series接口中namespace+metric_name， 查询当前监控的时间序列名称的时间序列维度列表。
	Dimensions []DimensionSeries `json:"dimensions"`

	// 时间序列名称。名称长度取值范围为1~255个字符。 取值范围： AOM提供的基础时间序列名称，cpuUsage、cpuCoreUsed等， cpuUage：cpu使用率； cpuCoreUsed：cpu内核占用； 用户上报的自定义时间序列名称。
	MetricName string `json:"metric_name"`
}

func (o QuerySample) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QuerySample struct{}"
	}

	return strings.Join([]string{"QuerySample", string(data)}, " ")
}
