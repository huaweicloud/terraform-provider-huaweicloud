package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// http响应头设置
type HttpResponseHeader struct {

	// 设置HTTP响应头参数。取值：\"Content-Disposition\", \"Content-Language\", \"Access-Control-Allow-Origin\",\"Access-Control-Allow-Methods\", \"Access-Control-Max-Age\", \"Access-Control-Expose-Headers\"或自定义头部。格式要求：长度1~100，以字母开头，可以使用字母、数字和短横杠。
	Name string `json:"name"`

	// 设置HTTP响应头参数的值。自定义HTTP响应头参数长度范围1~256，支持字母、数字和特定字符（.-_*#!%&+|^~'\"/:;,=@?）。
	Value *string `json:"value,omitempty"`

	// 设置http响应头操作类型，取值“set/delete”。set代表设置，delete代表删除。
	Action string `json:"action"`
}

func (o HttpResponseHeader) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HttpResponseHeader struct{}"
	}

	return strings.Join([]string{"HttpResponseHeader", string(data)}, " ")
}
