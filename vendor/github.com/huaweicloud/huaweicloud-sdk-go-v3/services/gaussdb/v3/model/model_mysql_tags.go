package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 标签列表，根据标签键值对创建实例。 - {key}表示标签键，不可以为空或重复。 - {value}表示标签值，可以为空。  如果创建实例时同时使用多个标签键值对，中间使用逗号分隔开，最多包含10组。
type MysqlTags struct {
	// 标签键。最大长度36个unicode字符。 key不能为空或者空字符串，不能为空格。 字符集：A-Z，a-z ，0-9，‘-’，‘_’，UNICODE字符（\\u4E00-\\u9FFF）。

	Key string `json:"key"`
	// 标签值。最大长度43个unicode字符。 可以为空字符串。 字符集：A-Z，a-z ，0-9，‘.’，‘-’，‘_’，UNICODE字符（\\u4E00-\\u9FFF）。

	Value string `json:"value"`
}

func (o MysqlTags) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlTags struct{}"
	}

	return strings.Join([]string{"MysqlTags", string(data)}, " ")
}
