package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DebugCaseResult struct {

	// 响应内容
	Body *string `json:"body,omitempty"`

	// 错误原因
	ErrorReason *string `json:"errorReason,omitempty"`

	Header *DebugCaseResultHeader `json:"header,omitempty"`

	// 请求名称
	Name *string `json:"name,omitempty"`

	// 响应时间
	ResponseTime *int32 `json:"responseTime,omitempty"`

	// 调试结果（1：成功；）
	Result *int32 `json:"result,omitempty"`

	// 响应正文
	ReturnBody *string `json:"returnBody,omitempty"`

	ReturnHeader *DebugCaseReturnHeader `json:"returnHeader,omitempty"`

	// 响应状态码
	StatusCode *string `json:"statusCode,omitempty"`

	// 请求地址
	Url *string `json:"url,omitempty"`
}

func (o DebugCaseResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DebugCaseResult struct{}"
	}

	return strings.Join([]string{"DebugCaseResult", string(data)}, " ")
}
