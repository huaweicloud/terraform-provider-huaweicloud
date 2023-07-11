package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TagMultyValueEntity struct {

	// 键。  key不能为空，长度1~128个字符（中文也可以输入128个字符）。  可用UTF-8格式表示的字母、数字和空格，以及以下字符： _ . : = + - @  key两头不能有空格字符。
	Key *string `json:"key,omitempty"`

	// 值列表。  值长度0~255个字符（中文也可以输入255个字符）。  值可用UTF-8格式表示的字母、数字和空格，以及以下字符： _ . : / = + - @。  值可以为空字符串。
	Values *[]string `json:"values,omitempty"`
}

func (o TagMultyValueEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagMultyValueEntity struct{}"
	}

	return strings.Join([]string{"TagMultyValueEntity", string(data)}, " ")
}
