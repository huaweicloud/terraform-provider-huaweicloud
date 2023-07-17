package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TagVo  标签对象。
type TagVo struct {

	//   键。 最大长度36个字符。 字符集：A-Z，a-z ， 0-9，‘-’，‘_’，UNICODE字符（\\u4E00-\\u9FFF）。
	Key string `json:"key"`

	// 值。 最大长度43个字符，可以为空字符串。 字符集：A-Z，a-z ， 0-9，‘.’，‘-’，‘_’，UNICODE字符（\\u4E00-\\u9FFF）。
	Value *string `json:"value,omitempty"`
}

func (o TagVo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagVo struct{}"
	}

	return strings.Join([]string{"TagVo", string(data)}, " ")
}
