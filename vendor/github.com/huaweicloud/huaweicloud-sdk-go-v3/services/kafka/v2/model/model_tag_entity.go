package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TagEntity struct {

	// 键。  key不能为空，长度1~128个字符（中文也可以输入128个字符）。  可用UTF-8格式表示的字母、数字和空格，以及以下字符： _ . : = + - @  key两头不能有空格字符。
	Key *string `json:"key,omitempty"`

	// 值。  长度0~255个字符（中文也可以输入255个字符）。  可用UTF-8格式表示的字母、数字和空格，以及以下字符： _ . : / = + - @。  value可以为空字符串。
	Value *string `json:"value,omitempty"`
}

func (o TagEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagEntity struct{}"
	}

	return strings.Join([]string{"TagEntity", string(data)}, " ")
}
