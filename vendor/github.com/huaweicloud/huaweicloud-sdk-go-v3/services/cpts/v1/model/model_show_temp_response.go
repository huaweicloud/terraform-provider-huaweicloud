package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTempResponse Response Object
type ShowTempResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	TempInfo       *TempInfo `json:"temp_info,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ShowTempResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTempResponse struct{}"
	}

	return strings.Join([]string{"ShowTempResponse", string(data)}, " ")
}
