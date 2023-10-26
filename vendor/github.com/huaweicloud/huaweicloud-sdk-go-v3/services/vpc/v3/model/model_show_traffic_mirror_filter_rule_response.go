package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTrafficMirrorFilterRuleResponse Response Object
type ShowTrafficMirrorFilterRuleResponse struct {
	TrafficMirrorFilterRule *TrafficMirrorFilterRule `json:"traffic_mirror_filter_rule,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowTrafficMirrorFilterRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTrafficMirrorFilterRuleResponse struct{}"
	}

	return strings.Join([]string{"ShowTrafficMirrorFilterRuleResponse", string(data)}, " ")
}
