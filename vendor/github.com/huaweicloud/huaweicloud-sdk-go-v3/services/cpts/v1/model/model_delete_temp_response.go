package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteTempResponse Response Object
type DeleteTempResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteTempResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTempResponse struct{}"
	}

	return strings.Join([]string{"DeleteTempResponse", string(data)}, " ")
}
