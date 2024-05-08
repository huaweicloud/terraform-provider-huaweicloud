package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UserResponseInfo 账号信息
type UserResponseInfo struct {

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`

	// 主机ID
	HostId *string `json:"host_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器ip
	HostIp *string `json:"host_ip,omitempty"`

	// 用户名
	UserName *string `json:"user_name,omitempty"`

	// 是否有登录权限
	LoginPermission *bool `json:"login_permission,omitempty"`

	// 是否有root权限
	RootPermission *bool `json:"root_permission,omitempty"`

	// 用户组
	UserGroupName *string `json:"user_group_name,omitempty"`

	// 用户目录
	UserHomeDir *string `json:"user_home_dir,omitempty"`

	// 用户启动shell
	Shell *string `json:"shell,omitempty"`

	// 最近扫描时间
	RecentScanTime *int64 `json:"recent_scan_time,omitempty"`

	// 容器id
	ContainerId *string `json:"container_id,omitempty"`

	// 容器名称
	ContainerName *string `json:"container_name,omitempty"`
}

func (o UserResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserResponseInfo struct{}"
	}

	return strings.Join([]string{"UserResponseInfo", string(data)}, " ")
}
