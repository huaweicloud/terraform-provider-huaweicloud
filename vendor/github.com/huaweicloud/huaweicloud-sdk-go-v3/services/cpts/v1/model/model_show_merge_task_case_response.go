package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowMergeTaskCaseResponse Response Object
type ShowMergeTaskCaseResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 扩展信息
	Extend *interface{} `json:"extend,omitempty"`

	Result         *CaseReportSummary `json:"result,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ShowMergeTaskCaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMergeTaskCaseResponse struct{}"
	}

	return strings.Join([]string{"ShowMergeTaskCaseResponse", string(data)}, " ")
}
