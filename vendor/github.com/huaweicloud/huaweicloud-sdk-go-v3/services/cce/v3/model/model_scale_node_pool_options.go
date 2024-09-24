package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ScaleNodePoolOptions 节点池伸缩选项配置
type ScaleNodePoolOptions struct {

	// 扩容状态检查策略: instant(同步检查), async(异步检查)。默认同步检查instant
	ScalableChecking *string `json:"scalableChecking,omitempty"`

	BillingConfigOverride *ScaleUpBillingConfigOverride `json:"billingConfigOverride,omitempty"`
}

func (o ScaleNodePoolOptions) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScaleNodePoolOptions struct{}"
	}

	return strings.Join([]string{"ScaleNodePoolOptions", string(data)}, " ")
}
