package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateVariableRequestBody UpdateVariableRequestBody
type UpdateVariableRequestBody struct {

	// 变量id
	Id int32 `json:"id"`

	// 变量名称
	Name string `json:"name"`

	// 变量类型（1：整数；2：枚举；3：文件[；5：文本](tag:hws,hws_hk)
	VariableType int32 `json:"variable_type"`

	// 变量值
	Variable []interface{} `json:"variable"`
}

func (o UpdateVariableRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateVariableRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateVariableRequestBody", string(data)}, " ")
}
