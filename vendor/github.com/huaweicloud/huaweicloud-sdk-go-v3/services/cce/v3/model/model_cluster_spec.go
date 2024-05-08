package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ClusterSpec 集群参数定义。
type ClusterSpec struct {

	// 集群类别： - CCE：CCE集群   CCE集群支持虚拟机与裸金属服务器混合、GPU、NPU等异构节点的混合部署，基于高性能网络模型提供全方位、多场景、安全稳定的容器运行环境。 [- Turbo: CCE Turbo集群。   全面基于云原生基础设施构建的云原生2.0的容器引擎服务，具备软硬协同、网络无损、安全可靠、调度智能的优势，为用户提供一站式、高性价比的全新容器服务体验。](tag:hws,hws_hk,dt,hcs,g42,sbc)
	Category *ClusterSpecCategory `json:"category,omitempty"`

	// 集群Master节点架构：  - VirtualMachine：Master节点为x86架构服务器 [- ARM64: Master节点为鲲鹏（ARM架构）服务器](tag:hws,hws_hk,hcs)
	Type *ClusterSpecType `json:"type,omitempty"`

	// 集群规格，当集群为v1.15及以上版本时支持创建后变更，详情请参见[变更集群规格](ResizeCluster.xml)。请按实际业务需求进行选择： - cce.s1.small: 小规模单控制节点CCE集群（最大50节点） - cce.s1.medium: 中等规模单控制节点CCE集群（最大200节点） - cce.s2.small: 小规模多控制节点CCE集群（最大50节点） - cce.s2.medium: 中等规模多控制节点CCE集群（最大200节点） - cce.s2.large: 大规模多控制节点CCE集群（最大1000节点） - cce.s2.xlarge: 超大规模多控制节点CCE集群（最大2000节点）  >    关于规格参数中的字段说明如下： >    - s1：单控制节点的集群，控制节点数为1。单控制节点故障后，集群将不可用，但已运行工作负载不受影响。 >    - s2：多控制节点的集群，即高可用集群，控制节点数为3。当某个控制节点故障时，集群仍然可用。 >    [- dec：表示专属云的CCE集群规格。例如cce.dec.s1.small表示小规模单控制节点的专属云CCE集群（最大50节点）。](tag:hws,hws_hk) >    - small：表示集群支持管理的最大节点规模为50节点。 >    - medium：表示集群支持管理的最大节点规模为200节点。 >    - large：表示集群支持管理的最大节点规模为1000节点。 >    - xlarge：表示集群支持管理的最大节点规模为2000节点。
	Flavor string `json:"flavor"`

	// 集群版本，与Kubernetes社区基线版本保持一致，建议选择最新版本。  在CCE控制台支持创建两种最新版本的集群。可登录CCE控制台创建集群，在“版本”处获取到集群版本。 其它集群版本，当前仍可通过api创建，但后续会逐渐下线，具体下线策略请关注CCE官方公告。  >    - 若不配置，默认创建最新版本的集群。 >    - 若指定集群基线版本但是不指定具体r版本，则系统默认选择对应集群版本的最新r版本。建议不指定具体r版本由系统选择最新版本。 [>    - Turbo集群支持1.19及以上版本商用。](tag:hws,hws_hk,dt) [>    - Turbo集群支持1.23及以上版本商用。](tag:hcs,g42,sbc)
	Version *string `json:"version,omitempty"`

	// CCE集群平台版本号，表示集群版本(version)下的内部版本。用于跟踪某一集群版本内的迭代，集群版本内唯一，跨集群版本重新计数。不支持用户指定，集群创建时自动选择对应集群版本的最新平台版本。  platformVersion格式为：cce.X.Y - X: 表示内部特性版本。集群版本中特性或者补丁修复，或者OS支持等变更场景。其值从1开始单调递增。 - Y: 表示内部特性版本的补丁版本。仅用于特性版本上线后的软件包更新，不涉及其他修改。其值从0开始单调递增。
	PlatformVersion *string `json:"platformVersion,omitempty"`

	// 集群描述，对于集群使用目的的描述，可根据实际情况自定义，默认为空。集群创建成功后可通过接口[更新指定的集群](cce_02_0240.xml)来做出修改，也可在CCE控制台中对应集群的“集群详情”下的“描述”处进行修改。仅支持utf-8编码。
	Description *string `json:"description,omitempty"`

	// 集群的API Server服务端证书中的自定义SAN（Subject Alternative Name）字段，遵从SSL标准X509定义的格式规范。  1. 不允许出现同名重复。 2. 格式符合IP和域名格式。  示例: ``` SAN 1: DNS Name=example.com SAN 2: DNS Name=www.example.com SAN 3: DNS Name=example.net SAN 4: IP Address=93.184.216.34 ```
	CustomSan *[]string `json:"customSan,omitempty"`

	// 集群是否使用IPv6模式，1.15版本及以上支持。
	Ipv6enable *bool `json:"ipv6enable,omitempty"`

	// CCE Turbo集群
	OffloadCluster *bool `json:"offloadCluster,omitempty"`

	HostNetwork *HostNetwork `json:"hostNetwork"`

	ContainerNetwork *ContainerNetwork `json:"containerNetwork"`

	EniNetwork *EniNetwork `json:"eniNetwork,omitempty"`

	ServiceNetwork *ServiceNetwork `json:"serviceNetwork,omitempty"`

	Authentication *Authentication `json:"authentication,omitempty"`

	// 集群的计费方式。 - 0: 按需计费 [- 1: 包周期](tag:hws,hws_hk)  默认为“按需计费”。
	BillingMode *int32 `json:"billingMode,omitempty"`

	// 控制节点的高级配置
	Masters *[]MasterSpec `json:"masters,omitempty"`

	// 服务网段参数，kubernetes clusterIP取值范围，1.11.7版本及以上支持。创建集群时如若未传参，默认为\"10.247.0.0/16\"。该参数废弃中，推荐使用新字段serviceNetwork，包含IPv4服务网段。
	KubernetesSvcIpRange *string `json:"kubernetesSvcIpRange,omitempty"`

	// 集群资源标签
	ClusterTags *[]ResourceTag `json:"clusterTags,omitempty"`

	// 服务转发模式，支持以下两种实现：  - iptables：社区传统的kube-proxy模式，完全以iptables规则的方式来实现service负载均衡。该方式最主要的问题是在服务多的时候产生太多的iptables规则，非增量式更新会引入一定的时延，大规模情况下有明显的性能问题。 - ipvs：主导开发并在社区获得广泛支持的kube-proxy模式，采用增量式更新，吞吐更高，速度更快，并可以保证service更新期间连接保持不断开，适用于大规模场景。  > 默认使用iptables转发模式。
	KubeProxyMode *ClusterSpecKubeProxyMode `json:"kubeProxyMode,omitempty"`

	// 可用区（仅查询返回字段）。  [CCE支持的可用区请参考[地区和终端节点](https://developer.huaweicloud.com/endpoint?CCE)](tag:hws)  [CCE支持的可用区请参考[地区和终端节点](https://developer.huaweicloud.com/intl/zh-cn/endpoint?CCE)](tag:hws_hk)
	Az *string `json:"az,omitempty"`

	ExtendParam *ClusterExtendParam `json:"extendParam,omitempty"`

	// 支持Istio
	SupportIstio *bool `json:"supportIstio,omitempty"`

	// 集群控制节点系统盘、数据盘加密。默认使用AES_256加密算法。CCE、Turbo集群1.25及以上版本开始支持。集群创建后不支持修改。开启后存在一定的磁盘读写性能损耗。
	EnableMasterVolumeEncryption *bool `json:"enableMasterVolumeEncryption,omitempty"`

	// 覆盖集群默认组件配置  若指定了不支持的组件或组件不支持的参数，该配置项将被忽略。  当前支持的可配置组件及其参数详见 [[配置管理](https://support.huaweicloud.com/usermanual-cce/cce_10_0213.html)](tag:hws) [[配置管理](https://support.huaweicloud.com/intl/zh-cn/usermanual-cce/cce_10_0213.html)](tag:hws_hk)
	ConfigurationsOverride *[]PackageConfiguration `json:"configurationsOverride,omitempty"`
}

func (o ClusterSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterSpec struct{}"
	}

	return strings.Join([]string{"ClusterSpec", string(data)}, " ")
}

type ClusterSpecCategory struct {
	value string
}

type ClusterSpecCategoryEnum struct {
	CCE   ClusterSpecCategory
	TURBO ClusterSpecCategory
}

func GetClusterSpecCategoryEnum() ClusterSpecCategoryEnum {
	return ClusterSpecCategoryEnum{
		CCE: ClusterSpecCategory{
			value: "CCE",
		},
		TURBO: ClusterSpecCategory{
			value: "Turbo",
		},
	}
}

func (c ClusterSpecCategory) Value() string {
	return c.value
}

func (c ClusterSpecCategory) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ClusterSpecCategory) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type ClusterSpecType struct {
	value string
}

type ClusterSpecTypeEnum struct {
	VIRTUAL_MACHINE ClusterSpecType
	ARM64           ClusterSpecType
}

func GetClusterSpecTypeEnum() ClusterSpecTypeEnum {
	return ClusterSpecTypeEnum{
		VIRTUAL_MACHINE: ClusterSpecType{
			value: "VirtualMachine",
		},
		ARM64: ClusterSpecType{
			value: "ARM64",
		},
	}
}

func (c ClusterSpecType) Value() string {
	return c.value
}

func (c ClusterSpecType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ClusterSpecType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type ClusterSpecKubeProxyMode struct {
	value string
}

type ClusterSpecKubeProxyModeEnum struct {
	IPTABLES ClusterSpecKubeProxyMode
	IPVS     ClusterSpecKubeProxyMode
}

func GetClusterSpecKubeProxyModeEnum() ClusterSpecKubeProxyModeEnum {
	return ClusterSpecKubeProxyModeEnum{
		IPTABLES: ClusterSpecKubeProxyMode{
			value: "iptables",
		},
		IPVS: ClusterSpecKubeProxyMode{
			value: "ipvs",
		},
	}
}

func (c ClusterSpecKubeProxyMode) Value() string {
	return c.value
}

func (c ClusterSpecKubeProxyMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ClusterSpecKubeProxyMode) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
