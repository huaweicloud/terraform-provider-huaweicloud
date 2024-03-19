package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowHistoryRunInfoResponse Response Object
type ShowHistoryRunInfoResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 报告列表
	LogList        *[]HistoryRunInfo `json:"log_list,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o ShowHistoryRunInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowHistoryRunInfoResponse struct{}"
	}

	return strings.Join([]string{"ShowHistoryRunInfoResponse", string(data)}, " ")
}
