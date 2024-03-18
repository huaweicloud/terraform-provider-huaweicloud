package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowHttpInfoResponse Response Object
type ShowHttpInfoResponse struct {
	Https *HttpInfoResponseBody `json:"https,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowHttpInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowHttpInfoResponse struct{}"
	}

	return strings.Join([]string{"ShowHttpInfoResponse", string(data)}, " ")
}
