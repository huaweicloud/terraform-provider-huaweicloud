package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PortResponseInfo 端口信息
type PortResponseInfo struct {

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 监听ip
	Laddr *string `json:"laddr,omitempty"`

	// port status, normal, danger or unknow   - \"normal\" : 正常   - \"danger\" : 危险   - \"unknow\" : 未知
	Status *string `json:"status,omitempty"`

	// 端口号
	Port *int32 `json:"port,omitempty"`

	// 类型
	Type *string `json:"type,omitempty"`

	// 进程ID
	Pid *int32 `json:"pid,omitempty"`

	// 程序文件
	Path *string `json:"path,omitempty"`

	// agent id
	AgentId *string `json:"agent_id,omitempty"`

	// 容器id
	ContainerId *string `json:"container_id,omitempty"`
}

func (o PortResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PortResponseInfo struct{}"
	}

	return strings.Join([]string{"PortResponseInfo", string(data)}, " ")
}
