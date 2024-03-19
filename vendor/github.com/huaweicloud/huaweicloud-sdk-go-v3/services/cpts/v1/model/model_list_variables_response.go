package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListVariablesResponse Response Object
type ListVariablesResponse struct {

	// 响应吗
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 变量列表
	VariableList   *[]VariableDetail `json:"variable_list,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o ListVariablesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListVariablesResponse struct{}"
	}

	return strings.Join([]string{"ListVariablesResponse", string(data)}, " ")
}
