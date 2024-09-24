package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PolicyTargetBase 策略目标
type PolicyTargetBase struct {

	// **参数说明**：策略绑定的目标类型。 **取值范围**：device|product|app，device表示设备，product表示产品，app表示整个资源空间。
	TargetType string `json:"target_type"`

	// 策略绑定的目标ID
	TargetId string `json:"target_id"`
}

func (o PolicyTargetBase) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PolicyTargetBase struct{}"
	}

	return strings.Join([]string{"PolicyTargetBase", string(data)}, " ")
}
