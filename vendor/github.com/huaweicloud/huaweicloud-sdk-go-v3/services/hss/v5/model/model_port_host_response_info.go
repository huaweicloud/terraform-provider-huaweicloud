package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PortHostResponseInfo 开放端口的机器统计信息
type PortHostResponseInfo struct {

	// 镜像id
	ContainerId *string `json:"container_id,omitempty"`

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 主机ip
	HostIp *string `json:"host_ip,omitempty"`

	// 主机名称
	HostName *string `json:"host_name,omitempty"`

	// 监听ip
	Laddr *string `json:"laddr,omitempty"`

	// 进程可执行文件路径
	Path *string `json:"path,omitempty"`

	// pid
	Pid *int32 `json:"pid,omitempty"`

	// 端口
	Port *int32 `json:"port,omitempty"`

	// 状态
	Status *string `json:"status,omitempty"`

	// 端口类型：目前包括TCP，UDP两种
	Type *string `json:"type,omitempty"`

	// 容器名称
	ContainerName *string `json:"container_name,omitempty"`

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`
}

func (o PortHostResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PortHostResponseInfo struct{}"
	}

	return strings.Join([]string{"PortHostResponseInfo", string(data)}, " ")
}
