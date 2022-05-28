package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 创建标签请求结构体。
type UnbindTagsDto struct {

	// **参数说明**：要绑定标签的资源类型。 **取值范围**： - device：设备。
	ResourceType string `json:"resource_type"`

	// **参数说明**：要绑定标签的资源id。例如，资源类型为device，那么对应的资源id就是device_id。 **取值范围**：长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合。
	ResourceId string `json:"resource_id"`

	// **参数说明**：指定资源要解绑的标签键列表，标签键列表中各项之间不允许重复，不能填写不存在的标签键值 **取值范围**：标签键长度不超过64，只允许中文、字母、数字、以及_.-等字符的组合。
	TagKeys []string `json:"tag_keys"`
}

func (o UnbindTagsDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnbindTagsDto struct{}"
	}

	return strings.Join([]string{"UnbindTagsDto", string(data)}, " ")
}
