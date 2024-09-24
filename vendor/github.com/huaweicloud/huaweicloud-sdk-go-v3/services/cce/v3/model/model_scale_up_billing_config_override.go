package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ScaleUpBillingConfigOverride 节点池扩容时覆盖节点的默认计费模式配置
type ScaleUpBillingConfigOverride struct {

	// 节点计费类型，0(按需)，1(包周期)
	BillingMode int32 `json:"billingMode"`

	ExtendParam *ScaleUpExtendParam `json:"extendParam,omitempty"`
}

func (o ScaleUpBillingConfigOverride) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScaleUpBillingConfigOverride struct{}"
	}

	return strings.Join([]string{"ScaleUpBillingConfigOverride", string(data)}, " ")
}
