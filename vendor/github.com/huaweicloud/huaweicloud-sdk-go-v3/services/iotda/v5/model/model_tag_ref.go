package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TagRef struct {

	// **参数说明**：标签键名称，可以是一个明确的静态字符串，也可以是动态的模板参数引用 - 明确的静态字符串：\"myTagKey\"。**取值范围**：长度不超过64，只允许中文、字母、数字、以及_.-等字符的组合 - 参数引用: {\"ref\" : \"iotda::certificate::country\"}
	TagKey *interface{} `json:"tag_key,omitempty"`

	// **参数说明**：标签值，可以是一个明确的静态字符串，也可以是动态的模板参数引用 - 明确的静态字符串：\"myTagValue\"。**取值范围**：长度不超过128，只允许中文、字母、数字、以及_.-等字符的组合。 - 参数引用: {\"ref\" : \"iotda::certificate::country\"}
	TagValue *interface{} `json:"tag_value,omitempty"`
}

func (o TagRef) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagRef struct{}"
	}

	return strings.Join([]string{"TagRef", string(data)}, " ")
}
