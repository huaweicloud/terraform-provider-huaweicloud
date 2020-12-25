package common

import (
	"time"
)

// Operator 运营商
type Operator struct {
	// 运营商的唯一uuid
	ID string `json:"id"`

	// 运营商的名称
	Name string `json:"name,omitempty"`

	// 运营商的国际化名称
	I18nName string `json:"i18n_name,omitempty"`

	// 运营商的简写
	Sa string `json:"sa"`
}

type PublicIP struct {
	// Specifies the ID of the elastic IP address, which uniquely
	// identifies the elastic IP address.
	ID string `json:"id"`

	// Specifies the status of the elastic IP address.
	Status string `json:"status"`

	// Specifies the obtained elastic IP address.
	PublicIpAddress string `json:"public_ip_address"`

	// Value range: 4, 6, respectively, to create ipv4 and ipv6, when not created ipv4 by default
	IPVersion int `json:"ip_version"`

	// Specifies the private IP address bound to the elastic IP
	// address.
	PrivateIpAddress string `json:"private_ip_address"`

	// Specifies the port ID.
	PortID string `json:"port_id"`

	// Specifies the time for applying for the elastic IP address.
	CreateTime string `json:"create_time"`

	// Specifies the bandwidth ID of the elastic IP address.
	BandwidthID string `json:"bandwidth_id"`

	// Specifies the bandwidth size.
	BandwidthSize int `json:"bandwidth_size"`

	// Specifies whether the bandwidth is shared or exclusive.
	BandwidthShareType string `json:"bandwidth_share_type"`

	// Specifies the bandwidth name.
	BandwidthName string `json:"bandwidth_name"`

	//Operator information
	Operator Operator `json:"operator"`

	// Specifies the Siteid.
	SiteID string `json:"site_id"`

	// SiteInfo
	SiteInfo string `json:"site_info"`

	Region string `json:"region,omitempty"`
}

// GeoLocation 地理位置
type GeoLocation struct {
	// ID 标志
	ID string `json:"id"`

	// City 城市
	City string `json:"city,omitempty"`

	// I18nCity 城市的国家化名称
	I18nCity string `json:"i18n_city,omitempty"`

	// Province 省份
	Province string `json:"province,omitempty"`

	// I18nProvince 省份的国际化名称
	I18nProvince string `json:"i18n_province,omitempty"`

	// Area 区域
	Area string `json:"area,omitempty"`

	// I18nArea 区域国际化名称
	I18nArea string `json:"i18n_area,omitempty"`

	// Country 国家
	Country string `json:"country"`

	// I18nCountry 国家的国际化名称
	I18nCountry string `json:"i18n_country,omitempty"`
}

// Subnet represents a subnet. See package documentation for a top-level
// description of what this is.
type Subnet struct {
	// Specifies a resource ID in UUID format.
	ID string `json:"id"`

	// Specifies the subnet name. The value is a string of 1 to 64
	// characters that can contain letters, digits, underscores (_), and hyphens (-).
	Name string `json:"name"`

	// Specifies the network segment on which the subnet resides. The
	// value must be in CIDR format. The value must be within the CIDR block of the VPC. The
	// subnet mask cannot be greater than 28.
	Cidr string `json:"cidr"`

	// Specifies the gateway of the subnet. The value must be a valid
	// IP address. The value must be an IP address in the subnet segment.
	GatewayIP string `json:"gateway_ip"`

	// Specifies whether the DHCP function is enabled for the subnet.
	// The value can be true or false. If this parameter is left blank, it is set to true by
	// default.
	DhcpEnable bool `json:"dhcp_enable,omitempty"`

	// Specifies the IP address of DNS server 1 on the subnet. The
	// value must be a valid IP address.
	PrimaryDNS string `json:"primary_dns,omitempty"`

	// Specifies the IP address of DNS server 2 on the subnet. The
	// value must be a valid IP address.
	SecondaryDNS string `json:"secondary_dns,omitempty"`

	// Specifies the DNS server address list of a subnet. This field
	// is required if you need to use more than two DNS servers. This parameter value is the
	// superset of both DNS server address 1 and DNS server address 2.
	DNSList []string `json:"dnsList,omitempty"`

	// Specifies the ID of the VPC to which the subnet belongs.
	VpcID string `json:"vpc_id"`

	// Specifies the status of the subnet. The value can be ACTIVE,
	// DOWN, UNKNOWN, or ERROR.
	Status string `json:"status"`

	// Specifies the network (Native OpenStack API) ID.
	NeutronNetworkID string `json:"neutron_network_id"`

	// Specifies the subnet (Native OpenStack API) ID.
	NeutronSubnetID string `json:"neutron_subnet_id"`

	// SiteID
	SiteID string `json:"site_id,omitempty"`

	//SiteInfo
	SiteInfo string `json:"site_info,omitempty"`
}

type VPC struct {
	// Specifies a resource ID in UUID format.
	ID string `json:"id"`

	// Specifies the name of the VPC. The name must be unique for a
	// tenant. The value is a string of no more than 64 characters and can contain digits,
	// letters, underscores (_), and hyphens (-).
	Name string `json:"name,omitempty"`

	// Specifies the range of available subnets in the VPC. The value
	// must be in CIDR format, for example, 192.168.0.0/16. The value ranges from 10.0.0.0/8
	// to 10.255.255.0/24, 172.16.0.0/12 to 172.31.255.0/24, or 192.168.0.0/16 to
	// 192.168.255.0/24.
	Cidr string `json:"cidr,omitempty"`

	// SubnetNum
	SubnetNum int64 `json:"subnet_num"`

	Mode string `json:"mode,omitempty"`
}

// SiteAttribute 站点扩展属性
type SiteAttribute struct {
	// 站点属性的唯一uuid
	ID string `json:"id"`

	// 站点相对应的属性的key
	Key string `json:"site_attr"`

	// 站点相对应属性的value
	Value string `json:"site_attr_value"`
}

// SiteBase 站点的基本信息
type SiteBase struct {
	// 站点所在的城市
	City string `json:"city,omitempty"`

	// 城市的国家化名称
	I18nCity string `json:"i18n_city,omitempty"`

	// 站点所在的省份
	Province string `json:"province,omitempty"`

	// 省份的国际化名称
	I18nProvince string `json:"i18n_province,omitempty"`

	// 站点所在的区域
	Area string `json:"area,omitempty"`

	// 区域国际化名称
	I18nArea string `json:"i18n_area,omitempty"`

	// 站点所在的国家
	Country string `json:"country,omitempty"`

	// 国家的国际化名称
	I18nCountry string `json:"i18n_country,omitempty"`

	Operator *Operator `json:"operator,omitempty"`
}

// Site 站点信息
type Site struct {
	// 站点的唯一uuid
	ID string `json:"id"`

	// 站点的名称，最好按照一定的规则命名，比如:IEG-国家-区域-省-市-运营商
	Name string `json:"name"`

	SiteBase

	//站点的状态,
	Status string `json:"status"`
}

type Port struct {
	// Specifies the port ID, which uniquely identifies the port.
	ID string `json:"id"`

	// Specifies the port name. The value can contain no more than 255
	// characters. This parameter is left blank by default.
	Name string `json:"name"`

	// Specifies the ID of the network to which the port belongs. The
	// network ID must be a real one in the network environment.
	NetworkID string `json:"network_id"`

	// Specifies the administrative state of the port. The value can
	// only be?true, and the default value is?true.
	AdminStateUp bool `json:"admin_state_up"`

	// Specifies the port MAC address. The system automatically sets
	// this parameter, and you are not allowed to configure the parameter value.
	MacAddress string `json:"mac_address"`

	// Specifies the port IP address. A port supports only one fixed
	// IP address that cannot be changed.
	FixedIPs []FixedIp `json:"fixed_ips"`

	// Specifies the ID of the device to which the port belongs. The
	// system automatically sets this parameter, and you are not allowed to configure or
	// change the parameter value.
	DeviceID string `json:"device_id"`

	// Specifies the belonged device, which can be the DHCP server,
	// router, load balancers, or Nova. The system automatically sets this parameter, and
	// you are not allowed to configure or change the parameter value.
	DeviceOwner string `json:"device_owner"`

	// Specifies the status of the port. The value can
	// be?ACTIVE,?BUILD, or?DOWN.
	Status string `json:"status"`

	// Specifies the UUID of the security group. This attribute is
	// extended.
	SecurityGroups []string `json:"security_groups"`

	// 1. Specifies a set of zero or more allowed address pairs. An
	// address pair consists of an IP address and MAC address. This attribute is extended.
	// For details, see parameter?allow_address_pair. 2. The IP address cannot be?0.0.0.0.
	// 3. Configure an independent security group for the port if a large CIDR block (subnet
	// mask less than 24) is configured for parameter?allowed_address_pairs.
	AllowedAddressPairs []AllowedAddressPair `json:"allowed_address_pairs"`

	// Specifies a set of zero or more extra DHCP option pairs. An
	// option pair consists of an option value and name. This attribute is extended.
	ExtraDhcpOpts []ExtraDHCPOpt `json:"extra_dhcp_opts"`

	// Specifies the type of the bound vNIC. The value can
	// be?normal?or?direct. Parameter?normal?indicates software switching.
	// Parameter?direct?indicates SR-IOV PCIe passthrough, which is not supported.
	BindingvnicType string `json:"binding:vnic_type"`

	// Default private domain name of the main NIC
	DnsAssignment []DnsAssignment `json:"dns_assignment"`

	// Default private DNS name of the main NIC
	DnsName string `json:"dns_name"`

	// site id
	SiteID string `json:"site_id"`
}

type FixedIp struct {
	// Specifies the subnet ID. You cannot change the parameter
	// value.
	SubnetId string `json:"subnet_id,omitempty"`

	// Specifies the port IP address. You cannot change the parameter
	// value.
	IpAddress string `json:"ip_address,omitempty"`
}

type DnsAssignment struct {
	// 功能说明：fqdn
	Fqdn string `json:"fqdn,omitempty"`

	// 功能说明：hostname
	HostName string `json:"hostname,omitempty"`

	// 功能说明：ip_address
	IpAddress string `json:"ip_address,omitempty"`
}

type ExtraDHCPOpt struct {
	// 功能说明：Option名称
	OptName string `json:"opt_name,omitempty"`

	// 功能说明：Option值
	OptValue string `json:"opt_value,omitempty"`
}

type AllowedAddressPair struct {
	// Specifies the IP address. You cannot set it to 0.0.0.0.
	// Configure an independent security group for the port if a large CIDR block (subnet
	// mask less than 24) is configured for parameter allowed_address_pairs.
	IpAddress string `json:"ip_address,omitempty"`

	// Specifies the MAC address.
	MacAddress string `json:"mac_address,omitempty"`
}

type Bandwidth struct {
	// Specifies the bandwidth name. The value is a string of 1 to 64
	// characters that can contain letters, digits, underscores (_), and hyphens (-).
	Name string `json:"name"`

	// Specifies the bandwidth size. The value ranges from 1 Mbit/s to
	// 300 Mbit/s.
	Size int `json:"size"`

	// Specifies the bandwidth ID, which uniquely identifies the
	// bandwidth.
	ID string `json:"id"`

	// Specifies whether the bandwidth is shared or exclusive. The
	// value can be PER or WHOLE.
	ShareType string `json:"share_type"`

	// Specifies the elastic IP address of the bandwidth.  The
	// bandwidth, whose type is set to WHOLE, supports up to 20 elastic IP addresses. The
	// bandwidth, whose type is set to PER, supports only one elastic IP address.
	PublicipInfo []PublicIpinfo `json:"publicip_info"`

	// Specifies the tenant ID of the user.
	TenantId string `json:"tenant_id"`

	// Specifies the bandwidth type.
	BandwidthType string `json:"bandwidth_type"`

	// Specifies the charging mode (by traffic or by bandwidth).
	ChargeMode string `json:"charge_mode"`

	// Specifies the status of bandwidth
	Status string `json:"status"`

	SiteID string `json:"site_id,omitempty"`

	CreateTime time.Time `json:"create_time,omitempty"`

	SiteInfo string `json:"site_info,omitempty"`

	Operator Operator `json:"operator,omitempty"`

	UpdateTime time.Time `json:"update_time,omitempty"`
}

type PublicIpinfo struct {
	// Specifies the tenant ID of the user.
	PublicipId string `json:"publicip_id"`

	// Specifies the elastic IP address.
	PublicipAddress string `json:"publicip_address"`

	// Specifies the elastic IP version.
	IPVersion int `json:"ip_version"`

	// Specifies the elastic IP address type. The value can be
	// 5_telcom, 5_union, or 5_bgp.
	PublicipType string `json:"publicip_type"`
}

// Volume contains all the information associated with an OpenStack Volume.
type Volume struct {
	// Unique identifier for the volume.
	ID string `json:"id"`
	// Current status of the volume.
	Status string `json:"status"`
	// Size of the volume in GB.
	Size int `json:"size"`
	// AvailabilityZone is which availability zone the volume is in.
	AvailabilityZone string `json:"availability_zone"`
	// The date when this volume was created.
	CreatedAt time.Time `json:"-"`
	// The date when this volume was last updated
	UpdatedAt time.Time `json:"-"`
	// Instances onto which the volume is attached.
	Attachments []Attachment `json:"attachments"`
	// Human-readable display name for the volume.
	Name string `json:"name"`
	// Human-readable description for the volume.
	Description string `json:"description"`
	// The type of volume to create, either SATA or SSD.
	VolumeType string `json:"volume_type"`
	// The ID of the snapshot from which the volume was created
	SnapshotID string `json:"snapshot_id"`
	// The ID of another block storage volume from which the current volume was created
	SourceVolID string `json:"source_volid"`
	// Arbitrary key-value pairs defined by the user.
	Metadata map[string]string `json:"metadata"`
	// UserID is the id of the user who created the volume.
	UserID string `json:"user_id"`
	// Indicates whether this is a bootable volume.
	Bootable string `json:"bootable"`
	// Encrypted denotes if the volume is encrypted.
	Encrypted bool `json:"encrypted"`
	// ReplicationStatus is the status of replication.
	ReplicationStatus string `json:"replication_status"`
	// ConsistencyGroupID is the consistency group ID.
	ConsistencyGroupID string `json:"consistencygroup_id"`
	// Multiattach denotes if the volume is multi-attach capable.
	Multiattach bool `json:"multiattach"`

	//Cloud hard disk uri self-description information.
	Links []map[string]string `json:"links"`

	//Whether it is a shared cloud drive.
	//Shareable bool `json:"shareable"`
	//Volume image metadata
	VolumeImageMetadata map[string]string `json:"volume_image_metadata"`

	//The tenant ID to which the cloud drive belongs.
	TenantAttr string `json:"os-vol-tenant-attr:tenant_id"`

	//The host name to which the cloud drive belongs.
	HostAttr string `json:"os-vol-host-attr:host"`
	//Reserved attribute
	RepAttrDriverData string `json:"os-volume-replication:driver_data"`
	//Reserved attribute
	RepAttrExtendedStatus string `json:"os-volume-replication:extended_status"`
	//Reserved attribute
	MigAttrStat string `json:"os-vol-mig-status-attr:migstat"`
	//Reserved attribute
	MigAttrNameID string `json:"os-vol-mig-status-attr:name_id"`
}

type Attachment struct {
	AttachedAt   time.Time `json:"-"`
	AttachmentID string    `json:"attachment_id"`
	Device       string    `json:"device"`
	HostName     string    `json:"host_name"`
	ID           string    `json:"id"`
	ServerID     string    `json:"server_id"`
	VolumeID     string    `json:"volume_id"`
}

// VolumeType 卷类型
type VolumeType struct {
	// Unique identifier for the volume type.
	ID string `json:"id"`
	// Human-readable display name for the volume type.
	Name string `json:"name"`
}

type Flavor struct {
	// Specifies the ID of ECS specifications.
	ID string `json:"id"`

	// Specifies the name of the ECS specifications.
	Name string `json:"name"`

	// Specifies the number of CPU cores in the ECS specifications.
	Vcpus string `json:"vcpus"`

	// Specifies the memory size (MB) in the ECS specifications.
	Ram int64 `json:"ram"`

	// Specifies the system disk size in the ECS specifications.
	// The value 0 indicates that the disk size is not limited.
	Disk string `json:"disk"`

	// Specifies shortcut links for ECS flavors.
	Links []Link `json:"links"`

	// Specifies extended ECS specifications.
	OsExtraSpecs OsExtraSpecs `json:"os_extra_specs"`

	// Reserved
	Swap string `json:"swap"`

	// Reserved
	FlvEphemeral int64 `json:"OS-FLV-EXT-DATA:ephemeral"`

	// Reserved
	FlvDisabled bool `json:"OS-FLV-DISABLED:disabled"`

	// Reserved
	RxtxFactor int64 `json:"rxtx_factor"`

	// Reserved
	RxtxQuota string `json:"rxtx_quota"`

	// Reserved
	RxtxCap string `json:"rxtx_cap"`

	// Reserved
	AccessIsPublic bool `json:"os-flavor-access:is_public"`
}

type Link struct {
	// Specifies the shortcut link marker name.
	Rel string `json:"rel"`

	// Provides the corresponding shortcut link.
	Href string `json:"href"`

	// Specifies the shortcut link type.
	Type string `json:"type"`
}

type OsExtraSpecs struct {
	// Specifies the ECS specifications types
	PerformanceType string `json:"ecs:performancetype"`

	// Specifies the resource type.
	ResourceType string `json:"resource_type"`

	// Specifies the generation of an ECS type
	Generation string `json:"ecs:generation"`

	// Specifies a virtualization type
	VirtualizationEnvTypes string `json:"ecs:virtualization_env_types"`

	// Indicates whether the GPU is passthrough.
	PciPassthroughEnableGpu string `json:"pci_passthrough:enable_gpu"`

	// Indicates the technology used on the G1 and G2 ECSs,
	// including GPU virtualization and GPU passthrough.
	PciPassthroughGpuSpecs string `json:"pci_passthrough:gpu_specs"`

	// Indicates the model and quantity of passthrough-enabled GPUs on P1 ECSs.
	PciPassthroughAlias string `json:"pci_passthrough:alias"`

	// gpu info.wuzilin add
	InfoGPUName string `json:"info:gpu:name,omitempty"`

	// cpu
	InfoCpuName string `json:"info:cpu:name,omitempty"`

	CondOperationStatus string `json:"cond:operation:status"`

	CondOperationAz string `json:"cond:operation:az"`

	CondCompute string `json:"cond:compute"`

	CondImage string `json:"cond:image"`

	VifMaxNum string `json:"quota:vif_max_num"`

	PhysicsMaxRate string `json:"quota:physics_max_rate"`

	VifMultiqueueNum string `json:"quota:vif_multiqueue_num"`

	MinRate string `json:"quota:min_rate"`

	MaxRate string `json:"quota:max_rate"`

	MaxPps string `json:"quota:max_pps"`

	CPUSockets string `json:"hw:cpu_sockets"`

	NumaNodes string `json:"hw:numa_nodes"`

	CPUThreads string `json:"hw:cpu_threads"`

	MemPageSize string `json:"hw:mem_page_size"`

	ConnLimitTotal string `json:"quota:conn_limit_total"`

	CPUCores string `json:"hw:cpu_cores"`

	SpotExtraSpecs
}

// SpotExtraSpecs 增加spot属性
type SpotExtraSpecs struct {
	CondSpotBlockOperationAz     string `json:"cond:spot_block:operation:az"`
	CondSpotBlockLdh             string `json:"cond:spot_block:operation:longest_duration_hours"`
	CondSpotBlockLdc             string `json:"cond:spot_block:operation:longest_duration_count"`
	CondSpotBlockInterruptPolicy string `json:"cond:spot_block:operation:interrupt_policy"`
	CondSpotOperationAz          string `json:"cond:spot:operation:az"`
	CondSpotOperationStatus      string `json:"cond:spot:operation:status"`
}

// EdgeImageInfo 边缘镜像基本字段
type EdgeImageInfo struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Status                string `json:"status"`
	DiskFormat            string `json:"disk_format"`
	MinDiskGigabytes      int    `json:"min_disk"`
	MinRAMMegabytes       int    `json:"min_ram"`
	Owner                 string `json:"owner"`
	Protected             bool   `json:"protected"`
	Visibility            string `json:"visibility"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`
	Self                  string `json:"self"`
	Deleted               bool   `json:"deleted"`
	VirtualEnvType        string `json:"virtual_env_type"`
	DeletedAt             string `json:"deleted_at"`
	RelatedJobID          string `json:"related_job_id"`
	ImageType             string `json:"__imagetype"`
	Platform              string `json:"__platform"`
	OsType                string `json:"__os_type"`
	OsVersion             string `json:"__os_version"`
	IsRegistered          bool   `json:"__isregistered"`
	SupportKvm            string `json:"__support_kvm,omitempty"`
	SupportKvmGpuType     string `json:"__support_kvm_gpu_type,omitempty"`
	SupportKvmAscend310   string `json:"__support_kvm_ascend_310,omitempty"`
	SupportKvmHi1822Hiovs string `json:"__support_kvm_hi1822_hiovs,omitempty"`
	SupportArm            string `json:"__support_arm,omitempty"`
	HwFirmwareType        string `json:"hw_firmware_type,omitempty"`
}

// Coverage :Edge Coverage Rule
type Coverage struct {
	CoveragePolicy string         `json:"coverage_policy" required:"true"`
	CoverageLevel  string         `json:"coverage_level" required:"true"`
	CoverageSites  []CoverageSite `json:"coverage_sites,omitempty"`
}

// CoverageSite :Edge service coverage site
type CoverageSite struct {
	Site    string   `json:"site"`
	Demands []Demand `json:"demands"`
}

// Demand
type Demand struct {
	Operator string `json:"operator" required:"true"`
	Count    int    `json:"demand_count" required:"true"`
}

type ResourceOpts struct {
	//Name is the name to assign to the newly launched server.
	Name string `json:"name" required:"true"`

	// ImageRef [optional; required if ImageName is not provided] is the ID or
	// full URL to the image that contains the server's OS and initial state.
	// Also optional if using the boot-from-volume extension.
	ImageRef string `json:"image_ref" required:"true"`

	// FlavorRef [optional; required if FlavorName is not provided] is the ID or
	// full URL to the flavor that describes the server's specs.
	FlavorRef string `json:"flavor_ref" required:"true"`

	// UserData contains configuration information or scripts to use upon launch.
	// Create will base64-encode it for you, if it isn't already.
	UserData string `json:"user_data"`

	// AdminPass sets the root user password. If not set, a randomly-generated
	// password will be created and returned in the response.
	AdminPass string `json:"admin_pass,omitempty"`

	//secret for logging in server
	KeyName string `json:"key_name,omitempty"`

	// Networks dictates how this server will be attached to available networks.
	// By default, the server will be attached to all isolated networks for the
	// tenant.
	NetConfig NetConfig `json:"net_config" required:"true"`

	//Specifies the EIP bandwidth. If this parameter does not exist, no EIP is bound.
	//If this parameter exists, an EIP is bound.
	BandWidth *BandWidth `json:"bandwidth,omitempty"`

	//the number of servers created
	Count int `json:"count"`

	//System disk configuration of the ECS
	RootVolume RootVolume `json:"root_volume"`

	//Specifies the data disk configuration of the ECS. Each data structure indicates
	//a data disk to be created.
	DataVolumes []DataVolume `json:"data_volumes,omitempty"`

	//Specifies the security group of the ECS.
	SecurityGroups []SecurityGroup `json:"security_groups,omitempty"`

	// 边缘场景，待使用
	EdgeScenes string `json:"edge_scenes,omitempty"`
}

// NetConfig
type NetConfig struct {
	VpcID   string     `json:"vpc_id" required:"true"`
	Subnets []SubnetID `json:"subnets"`
	NicNum  int        `json:"nic_num"`
}

// Subnet
type SubnetID struct {
	ID string `json:"id" required:"true"`
}

type BandWidth struct {
	//BandWidth（Mbit/s）[1,300]。
	Size int `json:"size,omitempty"`

	//ShareTypde PER indicates exclusive, and WHOLE indicates shared.
	ShareType string `json:"sharetype" required:"true"`

	//ChargeMode
	ChargeMode string `json:"chargemode,omitempty"`

	//BandWidthID，When creating an elastic IP address for a bandwidth of the
	//WHOLE type, you can specify the original shared bandwidth.
	Id string `json:"id,omitempty"`
}

type RootVolume struct {
	//the disk type of the ECS system disk. The disk type must match the disk
	//type provided by the system.
	VolumeType string `json:"volume_type" required:"true"`

	//the system disk size. The unit is GB. The value ranges from 1 to 1024.
	Size int `json:"size,omitempty"`
}

type DataVolume struct {
	//the disk type of the ECS data disk. The disk type must match the disk
	//type provided by the system.
	VolumeType string `json:"volume_type" required:"true"`

	//the data disk size in GB. The value ranges from 10 to 32768.
	Size int `json:"size" required:"true"`
}

type SecurityGroup struct {
	//云服务器组ID，UUID格式。
	ID string `json:"id" required:"true"`
}

// KeyPair is an SSH key known to the OpenStack Cloud that is available to be
// injected into servers.
type KeyPair struct {
	// Name is used to refer to this keypair from other services within this
	// region.
	Name string `json:"name"`

	// Fingerprint is a short sequence of bytes that can be used to authenticate
	// or validate a longer public key.
	Fingerprint string `json:"fingerprint"`

	// PublicKey is the public key from this pair, in OpenSSH format.
	// "ssh-rsa AAAAB3Nz..."
	PublicKey string `json:"public_key"`

	// PrivateKey is the private key from this pair, in PEM format.
	// "-----BEGIN RSA PRIVATE KEY-----\nMIICXA..."
	// It is only present if this KeyPair was just returned from a Create call.
	PrivateKey string `json:"private_key,omitempty"`

	// UserID is the user who owns this KeyPair.
	UserID string `json:"user_id,omitempty"`
}

// ReqSecurityGroupRuleEntity 创建安全组的规则的结构体
type ReqSecurityGroupRuleEntity struct {
	Description     string      `json:"description,omitempty"`
	SecurityGroupID string      `json:"security_group_id"`
	Direction       string      `json:"direction"`
	EtherType       string      `json:"ethertype,omitempty"`
	Protocol        string      `json:"protocol,omitempty"`
	PortRangeMin    interface{} `json:"port_range_min"`
	PortRangeMax    interface{} `json:"port_range_max"`
	RemoteIPPrefix  string      `json:"remote_ip_prefix,omitempty"`
	RemoteGroupID   string      `json:"remote_group_id,omitempty"`
}

// RegionSecurityGroupItem region级别的安全组信息
type RegionSecurityGroupItem struct {
	RegionID              string `json:"region_id,omitempty"`
	RegionSecurityGroupID string `json:"region_security_group_id,omitempty"`
}

// RespSecurityGroupRuleEntity 获取安全组安全组的规则的结构体
type RespSecurityGroupRuleEntity struct {
	ID              string      `json:"id"`
	Description     string      `json:"description"`
	SecurityGroupID string      `json:"security_group_id"`
	Direction       string      `json:"direction"`
	EtherType       string      `json:"ethertype"`
	Protocol        string      `json:"protocol"`
	PortRangeMin    interface{} `json:"port_range_min"`
	PortRangeMax    interface{} `json:"port_range_max"`
	RemoteIPPrefix  string      `json:"remote_ip_prefix"`
	RemoteGroupID   string      `json:"remote_group_id"`
}
