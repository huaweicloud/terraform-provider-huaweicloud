package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 键值对标签。
type TagWithKeyValue struct {

	// 标签键。最大长度36个unicode字符。 key不能为空，不能为空字符串，不能重复。字符集：A-Z，a-z ， 0-9，‘-’，‘_’，UNICODE字符（\\u4E00-\\u9FFF）。
	Key string `json:"key"`

	// 标签值。最大长度43个unicode字符。 可以为空字符串。 字符集：A-Z，a-z ， 0-9，‘.’，‘-’，‘_’，UNICODE字符（\\u4E00-\\u9FFF）。
	Value string `json:"value"`
}

func (o TagWithKeyValue) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagWithKeyValue struct{}"
	}

	return strings.Join([]string{"TagWithKeyValue", string(data)}, " ")
}
