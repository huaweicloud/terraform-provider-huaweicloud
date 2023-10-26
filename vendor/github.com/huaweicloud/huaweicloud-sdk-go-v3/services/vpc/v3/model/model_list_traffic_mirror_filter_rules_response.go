package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTrafficMirrorFilterRulesResponse Response Object
type ListTrafficMirrorFilterRulesResponse struct {

	// 流量镜像筛选条件规则对象
	TrafficMirrorFilterRules *[]TrafficMirrorFilterRule `json:"traffic_mirror_filter_rules,omitempty"`

	PageInfo *PageInfo `json:"page_info,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListTrafficMirrorFilterRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTrafficMirrorFilterRulesResponse struct{}"
	}

	return strings.Join([]string{"ListTrafficMirrorFilterRulesResponse", string(data)}, " ")
}
