package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EventResourceResponseInfo 资源信息
type EventResourceResponseInfo struct {

	// 租户账号ID
	DomainId *string `json:"domain_id,omitempty"`

	// 项目ID
	ProjectId *string `json:"project_id,omitempty"`

	// 企业项目ID
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// Region名称
	RegionName *string `json:"region_name,omitempty"`

	// VPC ID
	VpcId *string `json:"vpc_id,omitempty"`

	// 云主机ID
	CloudId *string `json:"cloud_id,omitempty"`

	// 虚拟机名称
	VmName *string `json:"vm_name,omitempty"`

	// 虚拟机UUID，即主机ID
	VmUuid *string `json:"vm_uuid,omitempty"`

	// 容器ID
	ContainerId *string `json:"container_id,omitempty"`

	// 容器状态
	ContainerStatus *string `json:"container_status,omitempty"`

	// pod uid
	PodUid *string `json:"pod_uid,omitempty"`

	// pod name
	PodName *string `json:"pod_name,omitempty"`

	// namespace
	Namespace *string `json:"namespace,omitempty"`

	// 集群id
	ClusterId *string `json:"cluster_id,omitempty"`

	// 集群名称
	ClusterName *string `json:"cluster_name,omitempty"`

	// 镜像ID
	ImageId *string `json:"image_id,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 主机属性
	HostAttr *string `json:"host_attr,omitempty"`

	// 业务服务
	Service *string `json:"service,omitempty"`

	// 微服务
	MicroService *string `json:"micro_service,omitempty"`

	// 系统CPU架构
	SysArch *string `json:"sys_arch,omitempty"`

	// 操作系统位数
	OsBit *string `json:"os_bit,omitempty"`

	// 操作系统类型
	OsType *string `json:"os_type,omitempty"`

	// 操作系统名称
	OsName *string `json:"os_name,omitempty"`

	// 操作系统版本
	OsVersion *string `json:"os_version,omitempty"`
}

func (o EventResourceResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventResourceResponseInfo struct{}"
	}

	return strings.Join([]string{"EventResourceResponseInfo", string(data)}, " ")
}
