package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UnBindDevicePolicy 解绑策略请求体。
type UnBindDevicePolicy struct {

	// **参数说明**：解除策略绑定的目标类型。 **取值范围**：device|product|app，device表示设备，product表示产品，app表示整个资源空间。
	TargetType string `json:"target_type"`

	// 策略解绑的目标ID列表
	TargetIds []string `json:"target_ids"`
}

func (o UnBindDevicePolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnBindDevicePolicy struct{}"
	}

	return strings.Join([]string{"UnBindDevicePolicy", string(data)}, " ")
}
