package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDevicePolicy 更新策略请求体。
type UpdateDevicePolicy struct {

	// **参数说明**：策略名称。 **取值范围**：长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合。
	PolicyName *string `json:"policy_name,omitempty"`

	// **参数说明**：策略文档。
	Statement *[]Statement `json:"statement,omitempty"`
}

func (o UpdateDevicePolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDevicePolicy struct{}"
	}

	return strings.Join([]string{"UpdateDevicePolicy", string(data)}, " ")
}
