package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateBandwidthPolicyReq struct {

	// 配置流量控制策略。数组中一个元素对应一个时段的最大带宽，最多允许5个时段，且时段不能重叠。
	BandwidthPolicy []BandwidthPolicyDto `json:"bandwidth_policy"`
}

func (o UpdateBandwidthPolicyReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBandwidthPolicyReq struct{}"
	}

	return strings.Join([]string{"UpdateBandwidthPolicyReq", string(data)}, " ")
}
