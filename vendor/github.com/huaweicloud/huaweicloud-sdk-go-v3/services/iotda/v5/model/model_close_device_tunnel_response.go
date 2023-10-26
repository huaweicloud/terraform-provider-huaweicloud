package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CloseDeviceTunnelResponse Response Object
type CloseDeviceTunnelResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CloseDeviceTunnelResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CloseDeviceTunnelResponse struct{}"
	}

	return strings.Join([]string{"CloseDeviceTunnelResponse", string(data)}, " ")
}
