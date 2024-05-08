package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ProcessesHostResponseInfo 进程主机统计信息
type ProcessesHostResponseInfo struct {

	// path对应的sha256值
	Hash *string `json:"hash,omitempty"`

	// 主机ip
	HostIp *string `json:"host_ip,omitempty"`

	// 主机名称
	HostName *string `json:"host_name,omitempty"`

	// 启动参数
	LaunchParams *string `json:"launch_params,omitempty"`

	// 启动时间
	LaunchTime *int64 `json:"launch_time,omitempty"`

	// 进程可执行文件路径
	ProcessPath *string `json:"process_path,omitempty"`

	// 进程pid
	ProcessPid *int32 `json:"process_pid,omitempty"`

	// 文件权限
	RunPermission *string `json:"run_permission,omitempty"`

	// 容器id
	ContainerId *string `json:"container_id,omitempty"`

	// 容器名称
	ContainerName *string `json:"container_name,omitempty"`
}

func (o ProcessesHostResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProcessesHostResponseInfo struct{}"
	}

	return strings.Join([]string{"ProcessesHostResponseInfo", string(data)}, " ")
}
