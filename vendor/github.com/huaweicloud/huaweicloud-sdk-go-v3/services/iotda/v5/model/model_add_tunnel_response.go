package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddTunnelResponse Response Object
type AddTunnelResponse struct {

	// 隧道ID
	TunnelId *string `json:"tunnel_id,omitempty"`

	// 鉴权信息
	TunnelAccessToken *string `json:"tunnel_access_token,omitempty"`

	// 鉴权信息的过期时间, 单位：秒
	ExpiresIn *int32 `json:"expires_in,omitempty"`

	// websocket接入地址
	TunnelUri      *string `json:"tunnel_uri,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AddTunnelResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddTunnelResponse struct{}"
	}

	return strings.Join([]string{"AddTunnelResponse", string(data)}, " ")
}
