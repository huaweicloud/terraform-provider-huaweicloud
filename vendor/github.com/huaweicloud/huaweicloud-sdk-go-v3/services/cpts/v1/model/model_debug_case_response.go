package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DebugCaseResponse struct {
	// code

	Code *string `json:"code,omitempty"`
	// message

	Message *string `json:"message,omitempty"`
	// extend

	Extend *string `json:"extend,omitempty"`
	// result

	Result         *[]DebugCaseResult `json:"result,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o DebugCaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DebugCaseResponse struct{}"
	}

	return strings.Join([]string{"DebugCaseResponse", string(data)}, " ")
}
