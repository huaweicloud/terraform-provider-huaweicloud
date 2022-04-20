package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowProcessResponse struct {
	// code

	Code *string `json:"code,omitempty"`
	// message

	Message *string `json:"message,omitempty"`

	Json *UploadProcessJson `json:"json,omitempty"`
	// extend

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
