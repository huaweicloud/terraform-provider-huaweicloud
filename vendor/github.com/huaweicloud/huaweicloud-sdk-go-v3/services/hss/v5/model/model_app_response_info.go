package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AppResponseInfo 软件信息
type AppResponseInfo struct {

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器ip
	HostIp *string `json:"host_ip,omitempty"`

	// 软件名称
	AppName *string `json:"app_name,omitempty"`

	// 版本号
	Version *string `json:"version,omitempty"`

	// 更新时间，最近一次更新的时间，用毫秒表示
	UpdateTime *int64 `json:"update_time,omitempty"`

	// 最近扫描时间，用毫秒表示
	RecentScanTime *int64 `json:"recent_scan_time,omitempty"`

	// 容器id
	ContainerId *string `json:"container_id,omitempty"`

	// 容器名称
	ContainerName *string `json:"container_name,omitempty"`
}

func (o AppResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AppResponseInfo struct{}"
	}

	return strings.Join([]string{"AppResponseInfo", string(data)}, " ")
}
