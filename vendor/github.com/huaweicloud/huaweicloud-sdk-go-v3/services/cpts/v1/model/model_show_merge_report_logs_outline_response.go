package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowMergeReportLogsOutlineResponse Response Object
type ShowMergeReportLogsOutlineResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 扩展字段
	Extend *interface{} `json:"extend,omitempty"`

	Result         *ReportOutlineResult `json:"result,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o ShowMergeReportLogsOutlineResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMergeReportLogsOutlineResponse struct{}"
	}

	return strings.Join([]string{"ShowMergeReportLogsOutlineResponse", string(data)}, " ")
}
