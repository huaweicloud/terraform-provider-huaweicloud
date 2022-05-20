package templates

// CreateTemplateResp is the response from a Create operation.
type CreateTemplateResp struct {
	// 服务端返回的新添加的模板的id
	ID string `json:"id"`
}

// TemplateResponse is the response from a Get operation.
type TemplateResponse struct {
	// 模板ID
	Id string `json:"id"`
	// 模板名称
	Name string `json:"name"`
	// 是否是通用模板,只有为true才能在console中显示
	IsTemplate string `json:"is_template"`
	// Region信息
	Region string `json:"region"`
	// 项目ID
	Projectid string `json:"projectid"`
	// 目标端服务器名称
	TargetServerName string `json:"target_server_name"`
	// 可用区
	AvailabilityZone string `json:"availability_zone"`
	// 虚拟机规格
	Flavor string `json:"flavor"`
	// VPC信息
	Vpc VpcObject `json:"vpc"`
	// 网卡信息,支持多个网卡,如果是自动创建,只填一个,id使用“autoCreate”
	Nics []NicObject `json:"nics"`
	// 安全组,支持多个安全组,如果是自动创建,只填一个,id使用“autoCreate”
	SecurityGroups []SgObject `json:"security_groups"`
	// 公网IP信息
	PublicIP EipObject `json:"publicip"`
	// 磁盘信息
	Disks []DiskObject `json:"disk"`
	// 磁盘类型: "SAS", "SSD", "GPSSD", "ESSD"
	Volumetype string `json:"volumetype"`
	// 数据盘磁盘类型: "SAS", "SSD", "GPSSD", "ESSD"
	DataVolumeType string `json:"data_volume_type"`
	// 目的端服务器密码
	ServerPassword string `json:"target_password"`
}

// VpcObject vpc对象
type VpcObject struct {
	// 虚拟私有云ID,如果是自动创建,填“autoCreate”
	Id string `json:"id"`
	// 虚拟私有云名称
	Name string `json:"name"`
	// VPC的网段,默认192.168.0.0/16
	Cidr string `json:"cidr"`
}

// NicObject 网卡资源
type NicObject struct {
	// 子网ID,如果是自动创建,使用"autoCreate"
	Id string `json:"id"`
	// 子网名称
	Name string `json:"name"`
	// 子网网关/掩码
	Cidr string `json:"cidr"`
	// 虚拟机IP地址,如果没有这个字段,自动分配IP
	Ip string `json:"ip"`
}

// SgObject 安全组object
type SgObject struct {
	// 安全组ID
	Id string `json:"id"`
	// 安全组名称
	Name string `json:"name"`
}

// EipObject 公网ip
type EipObject struct {
	// 弹性公网IP类型,默认为5_bgp
	Type string `json:"type"`
	// 带宽大小,单位:Mbit/s
	BandwidthSize int `json:"bandwidth_size"`
}

// DiskObject 磁盘模板
type DiskObject struct {
	// 磁盘序号,从0开始
	Index int `json:"index"`
	// 磁盘名称
	Name string `json:"name"`
	// 磁盘类型,同volumetype字段
	Disktype string `json:"disktype"`
	// 磁盘大小,单位:GB
	Size int `json:"size"`
}
