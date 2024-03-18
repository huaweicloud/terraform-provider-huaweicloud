package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateVariableRequestBody CreateVariableRequestBody
type CreateVariableRequestBody struct {

	// 变量id
	Id int32 `json:"id"`

	// 变量名称
	Name string `json:"name"`

	// 变量类型（1：整数；2：枚举；3：文件[；5：文本](tag:hws,hws_hk)
	VariableType int32 `json:"variable_type"`

	// 变量值
	Variable []interface{} `json:"variable"`

	// 是否被引用
	IsQuoted bool `json:"is_quoted"`
}

func (o CreateVariableRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateVariableRequestBody struct{}"
	}

	return strings.Join([]string{"CreateVariableRequestBody", string(data)}, " ")
}
