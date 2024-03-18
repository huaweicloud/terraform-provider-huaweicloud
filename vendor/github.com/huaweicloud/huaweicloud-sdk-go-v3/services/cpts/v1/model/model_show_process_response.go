package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowProcessResponse Response Object
type ShowProcessResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	Json *UploadProcessJson `json:"json,omitempty"`

	// 扩展信息
	Extend         *string `json:"extend,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowProcessResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowProcessResponse struct{}"
	}

	return strings.Join([]string{"ShowProcessResponse", string(data)}, " ")
}
