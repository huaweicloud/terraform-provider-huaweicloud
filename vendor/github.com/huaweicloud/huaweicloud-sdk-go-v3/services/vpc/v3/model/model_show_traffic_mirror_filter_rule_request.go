package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTrafficMirrorFilterRuleRequest Request Object
type ShowTrafficMirrorFilterRuleRequest struct {

	// 流量镜像筛选规则ID
	TrafficMirrorFilterRuleId string `json:"traffic_mirror_filter_rule_id"`
}

func (o ShowTrafficMirrorFilterRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTrafficMirrorFilterRuleRequest struct{}"
	}

	return strings.Join([]string{"ShowTrafficMirrorFilterRuleRequest", string(data)}, " ")
}
