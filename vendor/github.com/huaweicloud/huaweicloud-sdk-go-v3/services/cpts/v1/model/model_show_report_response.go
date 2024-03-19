package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowReportResponse Response Object
type ShowReportResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 扩展信息
	Extend *string `json:"extend,omitempty"`

	Result         *ReportInfo `json:"result,omitempty"`
	HttpStatusCode int         `json:"-"`
}

func (o ShowReportResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowReportResponse struct{}"
	}

	return strings.Join([]string{"ShowReportResponse", string(data)}, " ")
}
