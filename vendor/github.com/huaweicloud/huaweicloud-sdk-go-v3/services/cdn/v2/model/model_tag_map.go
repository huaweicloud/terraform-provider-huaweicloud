package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TagMap struct {

	// 标签键。长度1-128个字符, 可用 UTF-8 格式表示的字母(包含中文)、数字和空格，以及以下字符： _ . : = + - @
	Key string `json:"key"`

	// 标签值。长度0-255个字符,  可用 UTF-8 格式表示的字母(包含中文)、数字和空格，以及以下字符： _ . : / = + - @
	Value *string `json:"value,omitempty"`
}

func (o TagMap) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagMap struct{}"
	}

	return strings.Join([]string{"TagMap", string(data)}, " ")
}
