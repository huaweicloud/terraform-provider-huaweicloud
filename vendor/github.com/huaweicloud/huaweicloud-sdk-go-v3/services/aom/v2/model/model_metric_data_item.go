package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 指标
type MetricDataItem struct {

	// 数据收集时间支持过去1天和未来半小时范围内的数据上报。数据收集时间需要满足：  当前UTC时间减去collect_time小于等于24小时或者collect_time减去当前UTC时间小于等于30分钟。  若数据上报时间早于当天8点，则指标监控页面只显示当天8点后的数据。 取值范围： UNIX时间戳，单位毫秒。
	CollectTime int64 `json:"collect_time"`

	Metric *MetricItemInfo `json:"metric"`

	// 指标数据的值。
	Values []ValueData `json:"values"`
}

func (o MetricDataItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MetricDataItem struct{}"
	}

	return strings.Join([]string{"MetricDataItem", string(data)}, " ")
}
