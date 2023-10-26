package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTrafficMirrorFilterRuleResponse Response Object
type UpdateTrafficMirrorFilterRuleResponse struct {
	TrafficMirrorFilterRule *TrafficMirrorFilterRule `json:"traffic_mirror_filter_rule,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateTrafficMirrorFilterRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTrafficMirrorFilterRuleResponse struct{}"
	}

	return strings.Join([]string{"UpdateTrafficMirrorFilterRuleResponse", string(data)}, " ")
}
