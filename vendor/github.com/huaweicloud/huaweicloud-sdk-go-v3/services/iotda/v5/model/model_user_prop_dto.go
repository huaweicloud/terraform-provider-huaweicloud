package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UserPropDto 用户自定义属性
type UserPropDto struct {

	// **参数说明**：用户自定义属性键。 **取值范围**：长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合。
	PropKey *string `json:"prop_key,omitempty"`

	// **参数说明**：用户自定义属性值。 **取值范围**：长度不超过128，只允许中文、字母、数字、以及_? '#().,&%@!-等字符的组合。
	PropValue *string `json:"prop_value,omitempty"`
}

func (o UserPropDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserPropDto struct{}"
	}

	return strings.Join([]string{"UserPropDto", string(data)}, " ")
}
