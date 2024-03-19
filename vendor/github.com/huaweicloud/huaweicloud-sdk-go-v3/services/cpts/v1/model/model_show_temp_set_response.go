package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTempSetResponse Response Object
type ShowTempSetResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// temps
	Temps          *[]TempDetailInfo `json:"temps,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o ShowTempSetResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTempSetResponse struct{}"
	}

	return strings.Join([]string{"ShowTempSetResponse", string(data)}, " ")
}
