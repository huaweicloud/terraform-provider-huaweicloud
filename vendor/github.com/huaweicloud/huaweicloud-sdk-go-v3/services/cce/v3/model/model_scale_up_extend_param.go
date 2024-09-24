package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ScaleUpExtendParam 节点池扩容时节点的计费配置
type ScaleUpExtendParam struct {

	// 包周期时长, 包月1-9, 包年1-3
	PeriodNum int32 `json:"periodNum"`

	// 包周期类型, year(包年), month(包月)
	PeriodType string `json:"periodType"`

	// 是否自动续费，可选参数，如果不填写，以节点池isAutoRenew属性为准。
	IsAutoRenew *bool `json:"isAutoRenew,omitempty"`

	// 是否自动付费，可选参数，如果不填写，以节点池isAutoPay属性为准。
	IsAutoPay *bool `json:"isAutoPay,omitempty"`
}

func (o ScaleUpExtendParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScaleUpExtendParam struct{}"
	}

	return strings.Join([]string{"ScaleUpExtendParam", string(data)}, " ")
}
