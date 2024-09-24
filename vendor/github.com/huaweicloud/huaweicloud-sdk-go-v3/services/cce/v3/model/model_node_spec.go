package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NodeSpec struct {

	// 节点的规格，CCE支持的节点规格请参考[节点规格说明](cce_02_0368.xml)获取。
	Flavor string `json:"flavor"`

	// 待创建节点所在的可用区，需要指定可用区（AZ）的名称，不填或者填random选择随机可用区。 [CCE支持的可用区请参考[地区和终端节点](https://developer.huaweicloud.com/endpoint?CCE)](tag:hws) [CCE支持的可用区请参考[地区和终端节点](https://developer.huaweicloud.com/intl/zh-cn/endpoint?CCE)](tag:hws_hk)
	Az string `json:"az"`

	// 节点的操作系统类型。具体支持的操作系统请参见[节点操作系统说明](node-os.xml)。 > - 系统会根据集群版本自动选择支持的系统版本。当前集群版本不支持该系统类型，则会报错。 > - 若在创建节点时指定了extendParam中的alpha.cce/NodeImageID参数，可以不填写此参数。 > - 创建节点池时，该参数为必选。 > - 若创建节点时使用共享磁盘空间，即磁盘初始化配置管理参数使用storage，且StorageGroups中virtualSpaces的name字段指定为share，该参数为必选。
	Os *string `json:"os,omitempty"`

	Login *Login `json:"login"`

	RootVolume *Volume `json:"rootVolume"`

	// 节点的数据盘参数（目前已支持通过控制台为CCE节点添加第二块数据盘）。 如果数据盘正供容器运行时和Kubelet组件使用，则不可被卸载，否则将导致节点不可用。 针对专属云节点，参数解释与rootVolume一致
	DataVolumes []Volume `json:"dataVolumes"`

	Storage *Storage `json:"storage,omitempty"`

	PublicIP *NodePublicIp `json:"publicIP,omitempty"`

	NodeNicSpec *NodeNicSpec `json:"nodeNicSpec,omitempty"`

	// 批量创建时节点的个数，必须为大于等于1，小于等于最大限额的正整数。作用于节点池时该项可以不填写。
	Count *int32 `json:"count,omitempty"`

	// 节点的计费模式： -  0: 按需付费 [- 1: 包周期](tag:hws,hws_hk) [- 2: 已废弃：自动付费包周期](tag:hws,hws_hk)
	BillingMode *int32 `json:"billingMode,omitempty"`

	// 支持给创建出来的节点加Taints来设置反亲和性，taints配置不超过20条。每条Taints包含以下3个参数：  - Key：必须以字母或数字开头，可以包含字母、数字、连字符、下划线和点，最长63个字符；另外可以使用DNS子域作为前缀。 - Value：必须以字符或数字开头，可以包含字母、数字、连字符、下划线和点，最长63个字符。 - Effect：只可选NoSchedule，PreferNoSchedule或NoExecute。 字段使用场景：在节点创建场景下，支持指定初始值，查询时不返回该字段；在节点池场景下，其中节点模板中支持指定初始值，查询时支持返回该字段；在其余场景下，查询时都不会返回该字段。  示例：  ``` \"taints\": [{   \"key\": \"status\",   \"value\": \"unavailable\",   \"effect\": \"NoSchedule\" }, {   \"key\": \"looks\",   \"value\": \"bad\",   \"effect\": \"NoSchedule\" }] ```
	Taints *[]Taint `json:"taints,omitempty"`

	// 格式为key/value键值对。键值对个数不超过20条。 - Key：必须以字母或数字开头，可以包含字母、数字、连字符、下划线和点，最长63个字符；另外可以使用DNS子域作为前缀，例如example.com/my-key，DNS子域最长253个字符。 - Value：可以为空或者非空字符串，非空字符串必须以字符或数字开头，可以包含字母、数字、连字符、下划线和点，最长63个字符。 字段使用场景：在节点创建场景下，支持指定初始值，查询时不返回该字段；在节点池场景下，其中节点模板中支持指定初始值，查询时支持返回该字段；在其余场景下，查询时都不会返回该字段。   示例： ``` \"k8sTags\": {   \"key\": \"value\" } ```
	K8sTags map[string]string `json:"k8sTags,omitempty"`

	// 云服务器组ID，若指定，将节点创建在该云服务器组下 > 创建节点池时该配置不会生效，若要保持节点池中的节点都在同一个云服务器组内，请在节点池 nodeManagement 字段中配置
	EcsGroupId *string `json:"ecsGroupId,omitempty"`

	// 云服务器故障域，将节点创建在指定故障域下。  >必须同时指定故障域策略的云服务器ID，且需要开启故障域特性开关
	FaultDomain *string `json:"faultDomain,omitempty"`

	// 指定DeH主机的ID，将节点调度到自己的DeH上。 >创建节点池添加节点时不支持该参数。
	DedicatedHostId *string `json:"dedicatedHostId,omitempty"`

	// 是否CCE Turbo集群节点 >创建节点池添加节点时不支持该参数。
	OffloadNode *bool `json:"offloadNode,omitempty"`

	// 节点来源是否为纳管节点
	IsStatic *bool `json:"isStatic,omitempty"`

	// 云服务器标签，键必须唯一，CCE支持的最大用户自定义标签数量依region而定，自定义标签数上限为8个。 字段使用场景：在节点创建场景下，支持指定初始值，查询时不返回该字段；在节点池场景下，其中节点模板中支持指定初始值，查询时支持返回该字段；在其余场景下，查询时都不会返回该字段。 > 标签键只能包含大写字母.小写字母、数字和特殊字符(-_)以及Unicode字符，长度不超过36个字符。
	UserTags *[]UserTag `json:"userTags,omitempty"`

	Runtime *Runtime `json:"runtime,omitempty"`

	// 自定义初始化标记。  CCE节点在初始化完成之前，会打上初始化未完成污点（node.cloudprovider.kubernetes.io/uninitialized）防止pod调度到节点上。  cce支持自定义初始化标记，在接收到initializedConditions参数后，会将参数值转换成节点标签，随节点下发，例如：cloudprovider.openvessel.io/inject-initialized-conditions=CCEInitial_CustomedInitial。  当节点上设置了此标签，会轮询节点的status.Conditions，查看conditions的type是否存在标记名，如CCEInitial、CustomedInitial标记，如果存在所有传入的标记，且状态为True，认为节点初始化完成，则移除初始化污点。  - 必须以字母、数字组成，长度范围1-20位。 - 标记数量不超过2个
	InitializedConditions *[]string `json:"initializedConditions,omitempty"`

	ExtendParam *NodeExtendParam `json:"extendParam,omitempty"`

	HostnameConfig *HostnameConfig `json:"hostnameConfig,omitempty"`

	// 服务器企业项目ID。CCE服务不实现EPS相关特性，该字段仅用于同步服务器企业项目ID。 创建节点/节点池场景：可指定已存在企业项目，当取值为空时，该字段继承集群企业项目属性。 更新节点池场景：配置修改后仅会对新增节点的服务器生效，存量节点需前往EPS界面迁移。 如果更新时不指定值，不会更新该字段。 当该字段为空时，返回集群企业项目。
	ServerEnterpriseProjectID *string `json:"serverEnterpriseProjectID,omitempty"`
}

func (o NodeSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeSpec struct{}"
	}

	return strings.Join([]string{"NodeSpec", string(data)}, " ")
}
