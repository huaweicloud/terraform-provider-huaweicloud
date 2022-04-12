package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowTaskSetResponse struct {
	// code

	Code *string `json:"code,omitempty"`
	// extend

	Extend *[]string `json:"extend,omitempty"`
	// message

	Message *string `json:"message,omitempty"`
	// 工程集详细信息

	Tasks          *[]Task `json:"tasks,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowTaskSetResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTaskSetResponse struct{}"
	}

	return strings.Join([]string{"ShowTaskSetResponse", string(data)}, " ")
}
