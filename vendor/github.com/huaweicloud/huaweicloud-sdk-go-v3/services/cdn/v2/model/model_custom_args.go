package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CustomArgs 鉴权参数
type CustomArgs struct {

	// 参数类型，custom_var：自定义，nginx_preset_var：预置的变量。
	Type string `json:"type"`

	// 参数,长度支持1-256，由数字0-9、字符a-z、A-Z，及特殊字符._-*#%|+^@?=组成。
	Key string `json:"key"`

	// 取值,长度支持1-256，由数字0-9、字符a-z、A-Z，及特殊字符._-*#%|+^@?=组成。
	Value string `json:"value"`
}

func (o CustomArgs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CustomArgs struct{}"
	}

	return strings.Join([]string{"CustomArgs", string(data)}, " ")
}
