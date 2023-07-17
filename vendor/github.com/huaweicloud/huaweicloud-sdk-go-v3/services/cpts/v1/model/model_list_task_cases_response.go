package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTaskCasesResponse Response Object
type ListTaskCasesResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 用例列表
	TestCases      *[]TestCaseBriefInfo `json:"test_cases,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o ListTaskCasesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTaskCasesResponse struct{}"
	}

	return strings.Join([]string{"ListTaskCasesResponse", string(data)}, " ")
}
