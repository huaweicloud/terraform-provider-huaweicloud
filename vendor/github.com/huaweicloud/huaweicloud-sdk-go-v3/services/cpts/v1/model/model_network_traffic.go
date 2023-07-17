package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NetworkTraffic struct {

	// 平均网络流量
	AvgNetworkTraffic *float32 `json:"avg_network_traffic,omitempty"`

	// 最大下行带宽
	MaxDownStream *int32 `json:"max_down_stream,omitempty"`

	// 最大网络流量（流量峰值）
	MaxNetworkTraffic *int32 `json:"max_network_traffic,omitempty"`

	// 最大上行带宽
	MaxUpstream *int32 `json:"max_upstream,omitempty"`

	// 最小网络流量
	MinNetworkTraffic *int32 `json:"min_network_traffic,omitempty"`
}

func (o NetworkTraffic) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NetworkTraffic struct{}"
	}

	return strings.Join([]string{"NetworkTraffic", string(data)}, " ")
}
