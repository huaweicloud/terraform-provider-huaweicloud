package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTrafficMirrorFilterRuleRequest Request Object
type UpdateTrafficMirrorFilterRuleRequest struct {

	// 流量镜像筛选条件规则ID
	TrafficMirrorFilterRuleId string `json:"traffic_mirror_filter_rule_id"`

	Body *UpdateTrafficMirrorFilterRuleRequestBody `json:"body,omitempty"`
}

func (o UpdateTrafficMirrorFilterRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTrafficMirrorFilterRuleRequest struct{}"
	}

	return strings.Join([]string{"UpdateTrafficMirrorFilterRuleRequest", string(data)}, " ")
}
