package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ReinstallNodeSpec 节点重装配置参数
type ReinstallNodeSpec struct {

	// 操作系统。指定自定义镜像场景将以IMS镜像的实际操作系统版本为准。请选择当前集群支持的操作系统版本，例如EulerOS 2.5、CentOS 7.6、EulerOS 2.8。
	Os string `json:"os"`

	Login *Login `json:"login"`

	// 节点名称 > 重装时指定将修改节点名称，且服务器名称会同步修改。默认以服务器当前名称作为节点名称。 > 命名规则：以小写字母开头，由小写字母、数字、中划线(-)组成，长度范围1-56位。
	Name *string `json:"name,omitempty"`

	ServerConfig *ReinstallServerConfig `json:"serverConfig,omitempty"`

	VolumeConfig *ReinstallVolumeConfig `json:"volumeConfig,omitempty"`

	RuntimeConfig *ReinstallRuntimeConfig `json:"runtimeConfig,omitempty"`

	K8sOptions *ReinstallK8sOptionsConfig `json:"k8sOptions,omitempty"`

	Lifecycle *NodeLifecycleConfig `json:"lifecycle,omitempty"`

	// 自定义初始化标记。  CCE节点在初始化完成之前，会打上初始化未完成污点（node.cloudprovider.kubernetes.io/uninitialized）防止pod调度到节点上。  cce支持自定义初始化标记，在接收到initializedConditions参数后，会将参数值转换成节点标签，随节点下发，例如：cloudprovider.openvessel.io/inject-initialized-conditions=CCEInitial_CustomedInitial。  当节点上设置了此标签，会轮询节点的status.Conditions，查看conditions的type是否存在标记名，如CCEInitial、CustomedInitial标记，如果存在所有传入的标记，且状态为True，认为节点初始化完成，则移除初始化污点。  - 必须以字母、数字组成，长度范围1-20位。 - 标记数量不超过2个
	InitializedConditions *[]string `json:"initializedConditions,omitempty"`

	ExtendParam *ReinstallExtendParam `json:"extendParam,omitempty"`

	HostnameConfig *HostnameConfig `json:"hostnameConfig,omitempty"`
}

func (o ReinstallNodeSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReinstallNodeSpec struct{}"
	}

	return strings.Join([]string{"ReinstallNodeSpec", string(data)}, " ")
}
