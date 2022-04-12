package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListVariablesResponse struct {
	// code

	Code *string `json:"code,omitempty"`
	// message

	Message *string `json:"message,omitempty"`
	// variable_list

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
