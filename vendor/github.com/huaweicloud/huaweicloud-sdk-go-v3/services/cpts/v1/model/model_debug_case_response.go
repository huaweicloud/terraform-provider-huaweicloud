package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DebugCaseResponse Response Object
type DebugCaseResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 扩展信息
	Extend *string `json:"extend,omitempty"`

	// 结果
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
