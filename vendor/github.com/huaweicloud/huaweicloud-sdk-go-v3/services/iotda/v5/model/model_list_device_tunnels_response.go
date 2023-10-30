package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDeviceTunnelsResponse Response Object
type ListDeviceTunnelsResponse struct {

	// 隧道信息列表。
	Tunnels        *[]TunnelInfo `json:"tunnels,omitempty"`
	HttpStatusCode int           `json:"-"`
}

func (o ListDeviceTunnelsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDeviceTunnelsResponse struct{}"
	}

	return strings.Join([]string{"ListDeviceTunnelsResponse", string(data)}, " ")
}
