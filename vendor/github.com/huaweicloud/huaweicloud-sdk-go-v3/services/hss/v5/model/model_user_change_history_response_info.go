package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UserChangeHistoryResponseInfo 账号变动历史信息
type UserChangeHistoryResponseInfo struct {

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`

	// 变更类型   - ADD ：添加   - DELETE ：删除   - MODIFY ： 修改
	ChangeType *string `json:"change_type,omitempty"`

	// 主机ID
	HostId *string `json:"host_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器私有IP
	PrivateIp *string `json:"private_ip,omitempty"`

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

	// 账号名称
	UserName *string `json:"user_name,omitempty"`

	// 到期时间，采用时间戳，默认毫秒，
	ExpireTime *int64 `json:"expire_time,omitempty"`

	// 账号增加、修改、删除等操作的变更时间
	RecentScanTime *int64 `json:"recent_scan_time,omitempty"`
}

func (o UserChangeHistoryResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserChangeHistoryResponseInfo struct{}"
	}

	return strings.Join([]string{"UserChangeHistoryResponseInfo", string(data)}, " ")
}
