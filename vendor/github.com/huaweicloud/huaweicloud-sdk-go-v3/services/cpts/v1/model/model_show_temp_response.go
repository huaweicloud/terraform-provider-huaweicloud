package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowTempResponse struct {
	// code

	Code *string `json:"code,omitempty"`
	// message

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
