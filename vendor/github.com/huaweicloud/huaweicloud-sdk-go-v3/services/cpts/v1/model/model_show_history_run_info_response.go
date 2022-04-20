package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowHistoryRunInfoResponse struct {
	// code

	Code *string `json:"code,omitempty"`
	// message

	Message *string `json:"message,omitempty"`
	// log_list

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
