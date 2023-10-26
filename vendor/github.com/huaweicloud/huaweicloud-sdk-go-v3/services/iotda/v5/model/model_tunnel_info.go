package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TunnelInfo 隧道信息
type TunnelInfo struct {

	// 隧道ID
	TunnelId *string `json:"tunnel_id,omitempty"`

	// 设备ID
	DeviceId *string `json:"device_id,omitempty"`

	// 隧道创建时间。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	CreateTime *string `json:"create_time,omitempty"`

	// 隧道更新时间。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	ClosedTime *string `json:"closed_time,omitempty"`

	// 隧道状态 CLOSED | OPEN
	Status *string `json:"status,omitempty"`

	SourceConnectState *ConnectState `json:"source_connect_state,omitempty"`

	DeviceConnectState *ConnectState `json:"device_connect_state,omitempty"`
}

func (o TunnelInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TunnelInfo struct{}"
	}

	return strings.Join([]string{"TunnelInfo", string(data)}, " ")
}
