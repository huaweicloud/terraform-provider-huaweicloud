package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Port 1、功能说明：端口对象 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A
type Port struct {

	// 1、功能描述：管理状态 2、取值范围：true/false 3、约束：N/A 4、默认值：true 5、权限：N/A
	AdminStateUp bool `json:"admin_state_up"`

	// 1、功能描述：主机ID 2、取值范围：N/A 3、约束：管理员权限，普通租户不可见 4、默认值：N/A 5、权限：N/A
	BindinghostId string `json:"binding:host_id"`

	// 1、功能描述：提供用户设置自定义信息 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	Bindingprofile *interface{} `json:"binding:profile"`

	// 1、功能描述：vif的详细信息， \"ovs_hybrid_plug\": 是否为ovs/bridge混合模式 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	BindingvifDetails *interface{} `json:"binding:vif_details"`

	// 1、功能描述：端口的接口类型 (ovs/hw_veb等)(扩展属性) 2、取值范围：N/A 3、约束：管理员权限，普通租户不可见 4、默认值：N/A 5、权限：N/A
	BindingvifType string `json:"binding:vif_type"`

	// 1、功能描述：绑定的vNIC类型normal: 软交换direct: SRIOV硬直通（不支持） 2、取值范围：normal或者direct 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	BindingvnicType string `json:"binding:vnic_type"`

	// 1、功能说明：创建时间 2、取值范围：格式 \"UTC时间 格式: yyyy-MM-ddTHH:mm:ss\"  3、约束：N/A 4、默认值：N/A 5、权限：N/A
	CreatedAt *sdktime.SdkTime `json:"created_at"`

	// 1、功能说明：创建时间 2、取值范围：格式 \"UTC时间 格式: yyyy-MM-ddTHH:mm:ss\"  3、约束：N/A 4、默认值：N/A 5、权限：N/A
	UpdatedAt *sdktime.SdkTime `json:"updated_at"`

	// 1、功能说明：端口描述 2、取值范围：0-255个字符，不能包含“<”和“>” 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	Description string `json:"description"`

	// 1、功能描述：端口所属设备ID 2、取值范围：标准UUID 3、约束：不支持设置和更新，由系统自动维护 4、默认值：N/A 5、权限：N/A
	DeviceId string `json:"device_id"`

	// 1、功能描述：设备所属（DHCP/Router/ lb/Nova） 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	DeviceOwner string `json:"device_owner"`

	// 1、功能描述：标识这个端口所属虚拟机的flavor 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	EcsFlavor string `json:"ecs_flavor"`

	// 1、功能描述：端口唯一标识 2、取值范围：标准UUID 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	Id string `json:"id"`

	// 1、功能描述：端口所属实例ID，例如RDS实例ID 2、取值范围：N/A 3、约束：不支持设置和更新，由系统自动维护 4、默认值：N/A 5、权限：N/A
	InstanceId string `json:"instance_id"`

	// 1、功能描述：端口所属实例类型，例如“RDS” 2、取值范围：N/A 3、约束：不支持设置和更新，由系统自动维护 4、默认值：N/A 5、权限：N/A
	InstanceType string `json:"instance_type"`

	// 1、功能描述：MAC地址 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	MacAddress string `json:"mac_address"`

	// 1、功能描述：端口名称 2、取值范围：默认为空，最大长度不超过255 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	Name string `json:"name"`

	// 1、功能描述：端口安全使能标记，如果不使能则安全组和dhcp防欺骗不生效 2、取值范围：true/false 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	PortSecurityEnabled bool `json:"port_security_enabled"`

	// 1、功能描述：port的私有IP地址 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	PrivateIps []PrivateIpInfo `json:"private_ips"`

	// 1、功能描述：项目ID 2、取值范围：UUID 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	ProjectId string `json:"project_id"`

	// 1、功能描述：安全组 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	SecurityGroups []string `json:"security_groups"`

	// 1、功能描述：端口状态 2、取值范围：ACTIVE，BUILD，DOWN 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	Status string `json:"status"`

	// 1、功能描述：租户ID 2、取值范围：UUID 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	TenantId string `json:"tenant_id"`

	// 1、功能描述：所属网络ID 2、取值范围：标准UUID 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	VirsubnetId string `json:"virsubnet_id"`

	// 1、功能描述：VPC的ID 2、取值范围：标准UUID 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	VpcId string `json:"vpc_id"`

	// 1、功能描述：VPC_租户ID 2、取值范围：UUID 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	VpcTenantId string `json:"vpc_tenant_id"`

	// 1、功能描述：本地IP 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	VtepIp string `json:"vtep_ip"`

	// 1、功能描述：是否使能efi，使能则表示端口支持vRoCE能力 2、取值范围：true or false 3、约束：N/A 4、默认值：false 5、权限：N/A
	EnableEfi bool `json:"enable_efi"`

	// 1、功能描述：作用域 2、取值范围：center，表示作用域为中心；{azId}，表示作用域为具体的可用区 3、约束：N/A 4、默认值：center 5、权限：N/A
	Scope string `json:"scope"`

	// 1、功能描述：端口所属的可用分区 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	ZoneId string `json:"zone_id"`

	// 1、功能描述：迁移目的节点信息，包括目的节点的binding:vif_details和binding:vif_type 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	BindingmigrationInfo *interface{} `json:"binding:migration_info"`

	// 功能说明：DHCP的扩展属性
	ExtraDhcpOpts []ExtraDhcpOpt `json:"extra_dhcp_opts"`

	// 1、功能描述：边缘场景位置类型 2、取值范围：N/A 3、约束：N/A 4、默认值：center 5、权限：N/A
	PositionType string `json:"position_type"`

	// 1、功能描述：端口绑定实例信息 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	InstanceInfo *interface{} `json:"instance_info"`

	// 1、功能描述：端口标签 2、取值范围：N/A 3、约束：N/A 4、默认值：N/A 5、权限：N/A
	Tags []string `json:"tags"`

	// 1、功能描述：IP/Mac对列表 2、取值范围：N/A 3、约束： - IP地址不允许为 “0.0.0.0/0” - 如果allowed_address_pairs配置地址池较大的CIDR（掩码小于24位），建议为该port配置一个单独的安全组。 - 如果allowed_address_pairs的IP地址为“1.1.1.1/0”，表示关闭源目地址检查开关。 - 被绑定的云服务器网卡allowed_address_pairs的IP地址填“1.1.1.1/0”。 4、默认值：N/A 5、权限：N/A
	AllowedAddressPairs []AllowedAddressPair `json:"allowed_address_pairs"`
}

func (o Port) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Port struct{}"
	}

	return strings.Join([]string{"Port", string(data)}, " ")
}
