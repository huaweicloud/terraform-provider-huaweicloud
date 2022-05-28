package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 更新设备组请求结构体
type UpdateDeviceGroupDto struct {

	// **参数说明**：设备组名称，单个资源空间下不可重复。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_? '#().,&%@!-等字符的组合。
	Name *string `json:"name,omitempty"`

	// **参数说明**：设备组描述。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_? '#().,&%@!-等字符的组合。
	Description *string `json:"description,omitempty"`
}

func (o UpdateDeviceGroupDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDeviceGroupDto struct{}"
	}

	return strings.Join([]string{"UpdateDeviceGroupDto", string(data)}, " ")
}
