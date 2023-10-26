package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTrafficMirrorFilterRuleResponse Response Object
type CreateTrafficMirrorFilterRuleResponse struct {
	TrafficMirrorFilterRule *TrafficMirrorFilterRule `json:"traffic_mirror_filter_rule,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateTrafficMirrorFilterRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrafficMirrorFilterRuleResponse struct{}"
	}

	return strings.Join([]string{"CreateTrafficMirrorFilterRuleResponse", string(data)}, " ")
}
