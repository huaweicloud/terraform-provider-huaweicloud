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
}

func (o VariableDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VariableDetail struct{}"
	}

	return strings.Join([]string{"VariableDetail", string(data)}, " ")
}
