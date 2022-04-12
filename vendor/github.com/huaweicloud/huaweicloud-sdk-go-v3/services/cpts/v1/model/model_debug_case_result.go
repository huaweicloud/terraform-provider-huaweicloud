package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DebugCaseResult struct {
	// body

	Body *string `json:"body,omitempty"`
	// errorReason

	ErrorReason *string `json:"errorReason,omitempty"`

	Header *DebugCaseResultHeader `json:"header,omitempty"`
	// name

	Name *string `json:"name,omitempty"`
	// responseTime

	ResponseTime *int32 `json:"responseTime,omitempty"`
	// result

	Result *int32 `json:"result,omitempty"`
	// returnBody

	ReturnBody *string `json:"returnBody,omitempty"`

	ReturnHeader *DebugCaseReturnHeader `json:"returnHeader,omitempty"`
	// statusCode

	StatusCode *string `json:"statusCode,omitempty"`
	// url

	Url *string `json:"url,omitempty"`
}

func (o DebugCaseResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DebugCaseResult struct{}"
	}

	return strings.Join([]string{"DebugCaseResult", string(data)}, " ")
}
