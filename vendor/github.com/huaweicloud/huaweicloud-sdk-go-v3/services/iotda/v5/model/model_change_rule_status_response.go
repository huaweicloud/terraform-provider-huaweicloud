package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ChangeRuleStatusResponse struct {

	// **参数说明**：规则的激活状态。 **取值范围**： - active：激活。 - inactive：未激活。
	Status         *string `json:"status,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ChangeRuleStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeRuleStatusResponse struct{}"
	}

	return strings.Join([]string{"ChangeRuleStatusResponse", string(data)}, " ")
}
