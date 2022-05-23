package templates

import (
	"github.com/chnsz/golangsdk"
)

// TemplateOpts 创建/更新虚拟机模板
type TemplateOpts struct {
	// 模板名称
	Name string `json:"name" required:"true"`
	// 是否是通用模板,只有为true才能在console中显示
	IsTemplate *bool `json:"is_template" required:"true"`
	// Region信息,如 cn-north-4
	Region string `json:"region" required:"true"`
	// 项目ID
	ProjectID string `json:"projectid" required:"true"`
	// 可用区,如 cn-north-4b
	AvailabilityZone string `json:"availability_zone,omitempty"`
	// 目标端服务器名称
	TargetServerName string `json:"target_server_name,omitempty"`
	// 虚拟机规格
	Flavor string `json:"flavor,omitempty"`
	// VPC信息
	Vpc *VpcRequest `json:"vpc,omitempty"`
	// 网卡信息,支持多个网卡,如果是自动创建,只填一个,id使用“autoCreate”
	Nics []NicRequest `json:"nics,omitempty"`
	// 安全组,支持多个安全组,如果是自动创建,只填一个,id使用“autoCreate”
	SecurityGroups []SgRequest `json:"security_groups,omitempty"`
	// 公网IP信息
	PublicIP *EipRequest `json:"publicip,omitempty"`
	// 磁盘信息
	Disks []DiskRequest `json:"disk,omitempty"`
	// 磁盘类型: "SAS", "SSD", "GPSSD", "ESSD"
	VolumeType string `json:"volumetype,omitempty"`
	// 数据盘磁盘类型: "SAS", "SSD", "GPSSD", "ESSD"
	DataVolumeType string `json:"data_volume_type,omitempty"`
	// 目的端服务器密码
	ServerPassword string `json:"target_password,omitempty"`
}

// VpcRequest vpc对象
type VpcRequest struct {
	// 虚拟私有云ID,如果是自动创建,填“autoCreate”
	Id string `json:"id" required:"true"`
	// 虚拟私有云名称
	Name string `json:"name" required:"true"`
	// VPC的网段,默认192.168.0.0/16
	Cidr string `json:"cidr,omitempty"`
}

// NicRequest 网卡资源
type NicRequest struct {
	// 子网ID,如果是自动创建,使用"autoCreate"
	Id string `json:"id" required:"true"`
	// 子网名称
	Name string `json:"name" required:"true"`
	// 子网网关/掩码
	Cidr string `json:"cidr,omitempty"`
	// 虚拟机IP地址,如果没有这个字段,自动分配IP
	Ip string `json:"ip,omitempty"`
}

// SgRequest 安全组object
type SgRequest struct {
	// 安全组ID
	Id string `json:"id" required:"true"`
	// 安全组名称
	Name string `json:"name" required:"true"`
}

// EipRequest 公网ip
type EipRequest struct {
	// 弹性公网IP类型,默认为5_bgp
	Type string `json:"type" required:"true"`
	// 带宽大小,单位:Mbit/s
	BandwidthSize int `json:"bandwidth_size" required:"true"`
}

// DiskRequest 磁盘模板
type DiskRequest struct {
	// 磁盘序号,从0开始
	Index int `json:"index" required:"true"`
	// 磁盘名称
	Name string `json:"name" required:"true"`
	// 磁盘类型,同volumetype字段
	Disktype string `json:"disktype" required:"true"`
	// 磁盘大小,单位:GB
	Size int `json:"size" required:"true"`
}

// Create 新增模板信息
func Create(c *golangsdk.ServiceClient, opts *TemplateOpts) (string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "template")
	if err != nil {
		return "", err
	}

	var rst golangsdk.Result
	_, rst.Err = c.Post(rootURL(c), b, &rst.Body, nil)

	var r CreateTemplateResp
	if err := rst.ExtractInto(&r); err != nil {
		return "", err
	}

	return r.ID, nil
}

// Get 查询指定ID模板信息
func Get(c *golangsdk.ServiceClient, id string) (*TemplateResponse, error) {
	var rst golangsdk.Result
	_, rst.Err = c.Get(templateURL(c, id), &rst.Body, nil)

	var r struct {
		Template TemplateResponse `json:"template"`
	}

	err := rst.ExtractInto(&r)
	if err != nil {
		return nil, err
	}

	return &r.Template, nil
}

// Update 修改模板信息
func Update(c *golangsdk.ServiceClient, id string, opts *TemplateOpts) *golangsdk.ErrResult {
	b, err := golangsdk.BuildRequestBody(opts, "template")
	if err != nil {
		return nil
	}

	var r golangsdk.ErrResult
	_, r.Err = c.Put(templateURL(c, id), b, nil, nil)
	return &r
}

// Delete 删除指定ID的模板
func Delete(c *golangsdk.ServiceClient, id string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(templateURL(c, id), nil)
	return &r
}
