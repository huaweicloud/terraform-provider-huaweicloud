package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTrafficMirrorFilterRuleRequestBody
type UpdateTrafficMirrorFilterRuleRequestBody struct {
	TrafficMirrorFilterRule *UpdateTrafficMirrorFilterRuleOption `json:"traffic_mirror_filter_rule"`
}

func (o UpdateTrafficMirrorFilterRuleRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTrafficMirrorFilterRuleRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateTrafficMirrorFilterRuleRequestBody", string(data)}, " ")
}
