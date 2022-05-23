package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 标签结构体。
type TagV5Dto struct {

	// **参数说明**：标签键，在同一资源下标签键唯一。绑定资源时，如果设置的键已存在，则将覆盖之前的标签值。如果设置的键值不存在，则新增标签。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_.-等字符的组合。
	TagKey string `json:"tag_key"`

	// **参数说明**：标签值。 **取值范围**：长度不超过128，只允许中文、字母、数字、以及_.-等字符的组合。
	TagValue *string `json:"tag_value,omitempty"`
}

func (o TagV5Dto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagV5Dto struct{}"
	}

	return strings.Join([]string{"TagV5Dto", string(data)}, " ")
}
