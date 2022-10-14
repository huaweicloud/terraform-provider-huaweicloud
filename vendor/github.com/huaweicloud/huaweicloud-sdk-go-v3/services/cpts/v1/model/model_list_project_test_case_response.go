package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListProjectTestCaseResponse struct {

	// 状态码
	Code *string `json:"code,omitempty"`

	// 描述
	Message *string `json:"message,omitempty"`

	// 用例集
	Directory      *[]ProjectDirectory `json:"directory,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ListProjectTestCaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProjectTestCaseResponse struct{}"
	}

	return strings.Join([]string{"ListProjectTestCaseResponse", string(data)}, " ")
}
