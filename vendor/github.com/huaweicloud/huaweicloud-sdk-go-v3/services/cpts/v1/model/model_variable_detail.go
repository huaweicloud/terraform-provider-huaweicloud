package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VariableDetail struct {

	// file_size
	FileSize *int32 `json:"file_size,omitempty"`

	// id
	Id *int32 `json:"id,omitempty"`

	// 是否被引用
	IsQuoted *bool `json:"is_quoted,omitempty"`

	// name
	Name *string `json:"name,omitempty"`

	// variable
	Variable *[]interface{} `json:"variable,omitempty"`

	// variable_type
	VariableType *int32 `json:"variable_type,omitempty"`

	// 变量读取模式，0：顺序模式；1：随机模式
	VariableMode *int32 `json:"variable_mode,omitempty"`

	// 变量共享模式，0：用例模式；1：并发模式
	ShareMode *int32 `json:"share_mode,omitempty"`
}

func (o VariableDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VariableDetail struct{}"
	}

	return strings.Join([]string{"VariableDetail", string(data)}, " ")
}
