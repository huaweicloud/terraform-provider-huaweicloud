package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowMergeCaseDetailResponse Response Object
type ShowMergeCaseDetailResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 扩展信息
	Extend *interface{} `json:"extend,omitempty"`

	Result         *CaseReportDetailResult `json:"result,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o ShowMergeCaseDetailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMergeCaseDetailResponse struct{}"
	}

	return strings.Join([]string{"ShowMergeCaseDetailResponse", string(data)}, " ")
}
