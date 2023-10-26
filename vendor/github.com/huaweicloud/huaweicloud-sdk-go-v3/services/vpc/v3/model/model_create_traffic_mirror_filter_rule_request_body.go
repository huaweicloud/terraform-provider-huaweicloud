package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTrafficMirrorFilterRuleRequestBody
type CreateTrafficMirrorFilterRuleRequestBody struct {
	TrafficMirrorFilterRule *CreateTrafficMirrorFilterRuleOption `json:"traffic_mirror_filter_rule"`
}

func (o CreateTrafficMirrorFilterRuleRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrafficMirrorFilterRuleRequestBody struct{}"
	}

	return strings.Join([]string{"CreateTrafficMirrorFilterRuleRequestBody", string(data)}, " ")
}
