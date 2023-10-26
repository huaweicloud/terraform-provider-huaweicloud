package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDeviceTunnelResponse Response Object
type DeleteDeviceTunnelResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteDeviceTunnelResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDeviceTunnelResponse struct{}"
	}

	return strings.Join([]string{"DeleteDeviceTunnelResponse", string(data)}, " ")
}
