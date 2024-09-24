package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// AutopilotClusterSpec 集群参数定义。
type AutopilotClusterSpec struct {

	// 集群类别。
	Category *AutopilotClusterSpecCategory `json:"category,omitempty"`

	// 集群Master节点架构：  - VirtualMachine：Master节点为x86架构服务器
	Type *AutopilotClusterSpecType `json:"type,omitempty"`

	// 集群规格，cce.autopilot.cluster
	Flavor string `json:"flavor"`

	// 集群版本，与Kubernetes社区基线版本保持一致，建议选择最新版本。  在CCE控制台支持创建两种最新版本的集群。可登录CCE控制台创建集群，在“版本”处获取到集群版本。 其它集群版本，当前仍可通过api创建，但后续会逐渐下线，具体下线策略请关注CCE官方公告。  >    - 若不配置，默认创建最新版本的集群。
	Version *string `json:"version,omitempty"`

	// CCE集群平台版本号，表示集群版本(version)下的内部版本。用于跟踪某一集群版本内的迭代，集群版本内唯一，跨集群版本重新计数。不支持用户指定，集群创建时自动选择对应集群版本的最新平台版本。  platformVersion格式为：cce.X.Y - X: 表示内部特性版本。集群版本中特性或者补丁修复，或者OS支持等变更场景。其值从1开始单调递增。 - Y: 表示内部特性版本的补丁版本。仅用于特性版本上线后的软件包更新，不涉及其他修改。其值从0开始单调递增。
	PlatformVersion *string `json:"platformVersion,omitempty"`

	// 集群描述，对于集群使用目的的描述，可根据实际情况自定义，默认为空。集群创建成功后可通过接口[更新指定的集群](cce_02_0240.xml)来做出修改，也可在CCE控制台中对应集群的“集群详情”下的“描述”处进行修改。仅支持utf-8编码。
	Description *string `json:"description,omitempty"`

	// 集群的API Server服务端证书中的自定义SAN（Subject Alternative Name）字段，遵从SSL标准X509定义的格式规范。  1. 不允许出现同名重复。 2. 格式符合IP和域名格式。  示例: ``` SAN 1: DNS Name=example.com SAN 2: DNS Name=www.example.com SAN 3: DNS Name=example.net SAN 4: IP Address=93.184.216.34 ```
	CustomSan *[]string `json:"customSan,omitempty"`

	// 集群是否配置SNAT。开启后您的集群可以通过NAT网关访问公网，默认使用所选的VPC中已有的NAT网关，否则系统将会为您自动创建一个默认规格的NAT网关并绑定弹性公网IP，自动配置SNAT规则。
	EnableSnat *bool `json:"enableSnat,omitempty"`

	// 集群是否配置镜像访问。为确保您的集群节点可以从容器镜像服务中拉取镜像，默认使用所选VPC中已有的SWR和OBS终端节点，否则将会为您自动新建SWR和OBS终端节点。
	EnableSWRImageAccess *bool `json:"enableSWRImageAccess,omitempty"`

	// 是否为Autopilot集群。
	EnableAutopilot *bool `json:"enableAutopilot,omitempty"`

	// 集群是否使用IPv6模式。
	Ipv6enable *bool `json:"ipv6enable,omitempty"`

	HostNetwork *AutopilotHostNetwork `json:"hostNetwork"`

	ContainerNetwork *AutopilotContainerNetwork `json:"containerNetwork"`

	EniNetwork *AutopilotEniNetwork `json:"eniNetwork,omitempty"`

	ServiceNetwork *AutopilotServiceNetwork `json:"serviceNetwork,omitempty"`

	Authentication *AutopilotAuthentication `json:"authentication,omitempty"`

	// 集群的计费方式。 - 0: 按需计费  默认为“按需计费”。
	BillingMode *int32 `json:"billingMode,omitempty"`

	// 服务网段参数，kubernetes clusterIP取值范围。创建集群时如若未传参，默认为\"10.247.0.0/16\"。该参数废弃中，推荐使用新字段serviceNetwork，包含IPv4服务网段。
	KubernetesSvcIpRange *string `json:"kubernetesSvcIpRange,omitempty"`

	// 集群资源标签
	ClusterTags *[]AutopilotResourceTag `json:"clusterTags,omitempty"`

	// 服务转发模式：  - iptables：社区传统的kube-proxy模式，完全以iptables规则的方式来实现service负载均衡。该方式最主要的问题是在服务多的时候产生太多的iptables规则，非增量式更新会引入一定的时延，大规模情况下有明显的性能问题。  > 默认使用iptables转发模式。
	KubeProxyMode *AutopilotClusterSpecKubeProxyMode `json:"kubeProxyMode,omitempty"`

	// 可用区（仅查询返回字段）。  [CCE支持的可用区请参考[地区和终端节点](https://developer.huaweicloud.com/endpoint?CCE)](tag:hws)  [CCE支持的可用区请参考[地区和终端节点](https://developer.huaweicloud.com/intl/zh-cn/endpoint?CCE)](tag:hws_hk)
	Az *string `json:"az,omitempty"`

	ExtendParam *AutopilotClusterExtendParam `json:"extendParam,omitempty"`

	// 覆盖集群默认组件配置  若指定了不支持的组件或组件不支持的参数，该配置项将被忽略。  当前支持的可配置组件及其参数详见 [[配置管理](https://support.huaweicloud.com/usermanual-cce/cce_10_0213.html)](tag:hws) [[配置管理](https://support.huaweicloud.com/intl/zh-cn/usermanual-cce/cce_10_0213.html)](tag:hws_hk)
	ConfigurationsOverride *[]AutopilotPackageConfiguration `json:"configurationsOverride,omitempty"`
}

func (o AutopilotClusterSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotClusterSpec struct{}"
	}

	return strings.Join([]string{"AutopilotClusterSpec", string(data)}, " ")
}

type AutopilotClusterSpecCategory struct {
	value string
}

type AutopilotClusterSpecCategoryEnum struct {
	TURBO AutopilotClusterSpecCategory
}

func GetAutopilotClusterSpecCategoryEnum() AutopilotClusterSpecCategoryEnum {
	return AutopilotClusterSpecCategoryEnum{
		TURBO: AutopilotClusterSpecCategory{
			value: "Turbo",
		},
	}
}

func (c AutopilotClusterSpecCategory) Value() string {
	return c.value
}

func (c AutopilotClusterSpecCategory) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AutopilotClusterSpecCategory) UnmarshalJSON(b []byte) error {
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

type AutopilotClusterSpecType struct {
	value string
}

type AutopilotClusterSpecTypeEnum struct {
	VIRTUAL_MACHINE AutopilotClusterSpecType
}

func GetAutopilotClusterSpecTypeEnum() AutopilotClusterSpecTypeEnum {
	return AutopilotClusterSpecTypeEnum{
		VIRTUAL_MACHINE: AutopilotClusterSpecType{
			value: "VirtualMachine",
		},
	}
}

func (c AutopilotClusterSpecType) Value() string {
	return c.value
}

func (c AutopilotClusterSpecType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AutopilotClusterSpecType) UnmarshalJSON(b []byte) error {
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

type AutopilotClusterSpecKubeProxyMode struct {
	value string
}

type AutopilotClusterSpecKubeProxyModeEnum struct {
	IPTABLES AutopilotClusterSpecKubeProxyMode
}

func GetAutopilotClusterSpecKubeProxyModeEnum() AutopilotClusterSpecKubeProxyModeEnum {
	return AutopilotClusterSpecKubeProxyModeEnum{
		IPTABLES: AutopilotClusterSpecKubeProxyMode{
			value: "iptables",
		},
	}
}

func (c AutopilotClusterSpecKubeProxyMode) Value() string {
	return c.value
}

func (c AutopilotClusterSpecKubeProxyMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AutopilotClusterSpecKubeProxyMode) UnmarshalJSON(b []byte) error {
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
