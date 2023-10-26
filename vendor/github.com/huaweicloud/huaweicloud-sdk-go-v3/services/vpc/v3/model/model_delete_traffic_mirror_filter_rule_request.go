package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteTrafficMirrorFilterRuleRequest Request Object
type DeleteTrafficMirrorFilterRuleRequest struct {

	// 流量镜像筛选条件规则ID
	TrafficMirrorFilterRuleId string `json:"traffic_mirror_filter_rule_id"`
}

func (o DeleteTrafficMirrorFilterRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTrafficMirrorFilterRuleRequest struct{}"
	}

	return strings.Join([]string{"DeleteTrafficMirrorFilterRuleRequest", string(data)}, " ")
}
