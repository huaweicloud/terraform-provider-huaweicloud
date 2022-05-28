package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 规则状态
type RuleStatus struct {

	// **参数说明**：规则的激活状态。 **取值范围**： - active：激活。 - inactive：未激活。
	Status string `json:"status"`
}

func (o RuleStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RuleStatus struct{}"
	}

	return strings.Join([]string{"RuleStatus", string(data)}, " ")
}
