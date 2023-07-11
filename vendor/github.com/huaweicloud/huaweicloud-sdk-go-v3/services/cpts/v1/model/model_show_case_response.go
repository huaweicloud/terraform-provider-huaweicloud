package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCaseResponse Response Object
type ShowCaseResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	TestCase       *CaseInfoDetail `json:"test_case,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ShowCaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCaseResponse struct{}"
	}

	return strings.Join([]string{"ShowCaseResponse", string(data)}, " ")
}
