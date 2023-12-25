package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodeExtendParam 创建节点时的扩展参数。
type NodeExtendParam struct {

	// 云服务器规格的分类。响应中会返回此字段。
	Ecsperformancetype *string `json:"ecs:performancetype,omitempty"`

	// 订单ID，节点付费类型为自动付费包周期类型时，响应中会返回此字段(仅创建场景涉及)。
	OrderID *string `json:"orderID,omitempty"`

	// 产品ID，节点付费类型为自动付费包周期类型时，响应中会返回此字段。
	ProductID *string `json:"productID,omitempty"`

	// 节点最大允许创建的实例数(Pod)，该数量包含系统默认实例，取值范围为16~256。  该设置的目的为防止节点因管理过多实例而负载过重，请根据您的业务需要进行设置。  节点可以创建多少个Pod，受多个参数影响，具体请参见[节点最多可以创建多少Pod](maxPods.xml)。
	MaxPods *int32 `json:"maxPods,omitempty"`

	// - month：月 - year：年 > 作为请求参数，billingMode为1（包周期）或2（已废弃：自动付费包周期）时生效，且为必选。 > 作为响应参数，仅在创建包周期节点时返回。
	PeriodType *string `json:"periodType,omitempty"`

	// 订购周期数，取值范围： - periodType=month（周期类型为月）时，取值为[1-9]。 - periodType=year（周期类型为年）时，取值为1。 > 作为请求参数，billingMode为1或2（已废弃）时生效，且为必选。 > 作为响应参数，仅在创建包周期节点时返回。
	PeriodNum *int32 `json:"periodNum,omitempty"`

	// 是否自动续订 - “true”：自动续订 - “false”：不自动续订 > billingMode为1或2（已废弃）时生效，不填写此参数时默认不会自动续费。
	IsAutoRenew *string `json:"isAutoRenew,omitempty"`

	// 是否自动扣款  - “true”：自动扣款 - “false”：不自动扣款  > billingMode为1或2（已废弃）时生效，billingMode为1时不填写此参数时默认不会自动扣款。（已废弃：billingMode为2时不填写此参数时默认会自动扣款）
	IsAutoPay *string `json:"isAutoPay,omitempty"`

	// Docker数据盘配置项（已废弃，请使用storage字段）。默认配置示例如下：  ``` \"DockerLVMConfigOverride\":\"dockerThinpool=vgpaas/90%VG;kubernetesLV=vgpaas/10%VG;diskType=evs;lvType=linear\" ```  默认配置在无VD类型磁盘时，会由于数据盘查找失败而出错，请根据真实盘符类型填写diskType。 包含如下字段：   - userLV（可选）：用户空间的大小，示例格式：vgpaas/20%VG   - userPath（可选）：用户空间挂载路径，示例格式：/home/wqt-test   - diskType：磁盘类型，目前只有evs、hdd和ssd三种格式   - lvType：逻辑卷的类型，目前支持linear和striped两种，示例格式：striped   - dockerThinpool：Docker盘的空间大小，示例格式：vgpaas/60%VG   - kubernetesLV：Kubelet空间大小，示例格式：vgpaas/20%VG
	DockerLVMConfigOverride *string `json:"DockerLVMConfigOverride,omitempty"`

	// 节点上单容器的可用磁盘空间大小，单位G。  不配置该值或值为0时将使用默认值，Devicemapper模式下默认值为10；OverlayFS模式默认不限制单容器可用空间大小，且dockerBaseSize设置仅在新版本集群的EulerOS节点上生效。  CCE节点容器运行时空间配置请参考[数据盘空间分配说明](cce_01_0341.xml)。  Devicemapper模式下建议dockerBaseSize配置不超过80G，设置过大时可能会导致容器运行时初始化时间过长而启动失败，若对容器磁盘大小有特殊要求，可考虑使用挂载外部或本地存储方式代替。
	DockerBaseSize *int32 `json:"dockerBaseSize,omitempty"`

	// 是否为CCE Turbo集群节点。
	OffloadNode *string `json:"offloadNode,omitempty"`

	// 节点的公钥。
	PublicKey *string `json:"publicKey,omitempty"`

	// 安装前执行脚本 > 输入的值需要经过Base64编码，方法为echo -n \"待编码内容\" | base64
	AlphaCcePreInstall *string `json:"alpha.cce/preInstall,omitempty"`

	// 安装后执行脚本 > 输入的值需要经过Base64编码，方法为echo -n \"待编码内容\" | base64。
	AlphaCcePostInstall *string `json:"alpha.cce/postInstall,omitempty"`

	// 如果创建裸金属节点，需要使用自定义镜像时用此参数。
	AlphaCceNodeImageID *string `json:"alpha.cce/NodeImageID,omitempty"`

	// - 弹性网卡队列数配置，默认配置示例如下：  ``` \"[{\\\"queue\\\":4}]\" ```  包含如下字段： - queue: 弹性网卡队列数。 - 仅在turbo集群的BMS节点时，该字段才可配置。 - 当前支持可配置队列数以及弹性网卡数：{\"1\":128, \"2\":92, \"4\":92, \"8\":32, \"16\":16,\"28\":9}, 既1弹性网卡队列可绑定128张弹性网卡，2队列弹性网卡可绑定92张，以此类推。 - 弹性网卡队列数越多，性能越强，但可绑定弹性网卡数越少，请根据您的需求进行配置（创建后不可修改）。
	NicMultiqueue *string `json:"nicMultiqueue,omitempty"`

	// - 弹性网卡预绑定比例配置，默认配置示例如下： ``` \"0.3:0.6\" ```   - 第一位小数：预绑定低水位，弹性网卡预绑定的最低比例（最小预绑定弹性网卡数 = ⌊节点的总弹性网卡数 * 预绑定低水位⌋）   - 第二位小数：预绑定高水位，弹性网卡预绑定的最高比例（最大预绑定弹性网卡数 = ⌊节点的总弹性网卡数 * 预绑定高水位⌋）   - BMS节点上绑定的弹性网卡数：Pod正在使用的弹性网卡数 + 最小预绑定弹性网卡数 < BMS节点上绑定的弹性网卡数 < Pod正在使用的弹性网卡数 + 最大预绑定弹性网卡数   - BMS节点上当预绑定弹性网卡数 < 最小预绑定弹性网卡数时：会绑定弹性网卡，使得预绑定弹性网卡数 = 最小预绑定弹性网卡数   - BMS节点上当预绑定弹性网卡数 > 最大预绑定弹性网卡数时：会定时解绑弹性网卡（约2分钟一次），直到预绑定弹性网卡数 = 最大预绑定弹性网卡数   - 取值范围：[0.0, 1.0]; 一位小数; 低水位 <= 高水位 - 仅在turbo集群的BMS节点时，该字段才可配置。 - 弹性网卡预绑定能加快工作负载的创建，但会占用IP，请根据您的需求进行配置。
	NicThreshold *string `json:"nicThreshold,omitempty"`

	// 节点的计费模式。已废弃，请使用NodeSpec中的billingMode字段。
	ChargingMode *int32 `json:"chargingMode,omitempty"`

	// 委托的名称。  委托是由租户管理员在统一身份认证服务（Identity and Access Management，IAM）上创建的，可以为CCE节点提供访问云服务器的临时凭证。 作为响应参数仅在创建节点传入时返回该字段。
	AgencyName *string `json:"agency_name,omitempty"`

	// 节点内存预留，Kubernetes相关组件预留值。[随节点规格变动，具体请参见[节点预留资源策略说明](https://support.huaweicloud.com/usermanual-cce/cce_10_0178.html)。](tag:hws)
	KubeReservedMem *int32 `json:"kubeReservedMem,omitempty"`

	// 节点内存预留，系统组件预留值。[随节点规格变动，具体请参见[节点预留资源策略说明](https://support.huaweicloud.com/usermanual-cce/cce_10_0178.html)。](tag:hws)
	SystemReservedMem *int32 `json:"systemReservedMem,omitempty"`

	// 节点密码，作为响应参数时，固定展示星号。
	InitNodePassword *string `json:"init-node-password,omitempty"`
}

func (o NodeExtendParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeExtendParam struct{}"
	}

	return strings.Join([]string{"NodeExtendParam", string(data)}, " ")
}
