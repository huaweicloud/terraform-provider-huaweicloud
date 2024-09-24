package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// JarPackageHostInfo 服务器列表
type JarPackageHostInfo struct {

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器ip
	HostIp *string `json:"host_ip,omitempty"`

	// Jar包名称
	FileName *string `json:"file_name,omitempty"`

	// Jar包名称(不带后缀)
	Name *string `json:"name,omitempty"`

	// Jar包类型
	Catalogue *string `json:"catalogue,omitempty"`

	// Jar包后缀
	FileType *string `json:"file_type,omitempty"`

	// Jar包版本
	Version *string `json:"version,omitempty"`

	// Jar包路径
	Path *string `json:"path,omitempty"`

	// Jar包hash
	Hash *string `json:"hash,omitempty"`

	// Jar包大小
	Size *int32 `json:"size,omitempty"`

	// uid
	Uid *int32 `json:"uid,omitempty"`

	// gid
	Gid *int32 `json:"gid,omitempty"`

	// 文件权限
	Mode *string `json:"mode,omitempty"`

	// 进程id
	Pid *int32 `json:"pid,omitempty"`

	// 进程可执行文件路径
	ProcPath *string `json:"proc_path,omitempty"`

	// 容器实例id
	ContainerId *string `json:"container_id,omitempty"`

	// 容器名称
	ContainerName *string `json:"container_name,omitempty"`

	// 包路径
	PackagePath *string `json:"package_path,omitempty"`

	// 显示的是否是嵌套包
	IsEmbedded *int32 `json:"is_embedded,omitempty"`

	// 扫描时间
	RecordTime *int64 `json:"record_time,omitempty"`
}

func (o JarPackageHostInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "JarPackageHostInfo struct{}"
	}

	return strings.Join([]string{"JarPackageHostInfo", string(data)}, " ")
}
