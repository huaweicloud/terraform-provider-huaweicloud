package desktops

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure required by the Create method to create a new desktop.
type CreateOpts struct {
	// Configuration of system volume.
	RootVolume *Volume `json:"root_volume" required:"true"`
	// Configuration of desktops.
	Desktops []DesktopConfig `json:"desktops" required:"true"`
	// Desktop type.
	// + DEDICATED: dedicated desktop.
	DesktopType string `json:"desktop_type" required:"true"`
	// Product ID of desktop.
	ProductId string `json:"product_id" required:"true"`
	// The availability zone where the desktop is located.
	AvailabilityZone string `json:"availability_zone,omitempty"`
	// Image type, the default value is private.
	// + market
	// + gold
	// + private
	ImageType string `json:"image_type,omitempty"`
	// Image ID.
	ImageId string `json:"image_id,omitempty"`
	// Vpc ID, first creation time must be specified.
	VpcId string `json:"vpc_id,omitempty"`
	// Whether to send emails to user mailbox during important operations.
	EmailNotification *bool `json:"email_notification,omitempty"`
	// Configuration of data volumes.
	DataVolumes []Volume `json:"data_volumes,omitempty"`
	// NIC information corresponding to the desktop.
	Nics []Nic `json:"nics,omitempty"`
	// Configuration of security groups, the default security group (WorkspaceUserSecurityGroup) must be specified.
	SecurityGroups []SecurityGroup `json:"security_groups,omitempty"`
	// Specifies the key/value pairs of the desktop.
	Tags []tags.ResourceTag `json:"tags,omitempty"`
	// EnterpriseProject ID of desktop
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
}

// DesktopConfig is an object to specified the basic configuration of desktop.
type DesktopConfig struct {
	// User name.
	UserName string `json:"user_name" required:"true"`
	// User email.
	UserEmail string `json:"user_email" required:"true"`
	// User group.
	UserGroup string `json:"user_group,omitempty"`
	// Desktop name.
	DesktopName string `json:"computer_name,omitempty"`
	// Name prefix of desktop.
	DesktopNamePrefix string `json:"desktop_name_prefix"`
}

// Volume is an object to specified the disk configuration of root volume or data volume.
type Volume struct {
	// Volume type.
	// + **SAS**: High I/O disk type.
	// + **SSD**: Ultra-high I/O disk type.
	Type string `json:"type" required:"true"`
	// Volume size.
	// For root volume, the valid value is range from 80 to 1020.
	// For data volume, the valid value is range from 10 to 8200.
	Size int `json:"size" required:"true"`
}

// Nic is an object to specified the NIC information corresponding to the desktop.
type Nic struct {
	// Network ID.
	NetworkId string `json:"subnet_id" required:"true"`
}

// SecurityGroup is an object to specified the security group to which the desktop belongs.
type SecurityGroup struct {
	ID string `json:"id" required:"true"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a desktop using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*CreateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r CreateResp
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to obtain the desktop detail by its ID.
func Get(c *golangsdk.ServiceClient, desktopId string) (*Desktop, error) {
	var r GetResp
	_, err := c.Get(resourceURL(c, desktopId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Desktop, err
}

// ProductUpdateOpts is the structure required by the UpdateProduct method to change the desktop product.
type ProductUpdateOpts struct {
	// Batch create configuration of desktop list.
	Desktops []DesktopUpdateConfig `json:"desktops" required:"true"`
	// Product ID.
	ProductId string `json:"product_id" required:"true"`
	// Whether the product ID can be changed when the desktop is powered on.
	Mode string `json:"mode" required:"true"`
}

// DesktopUpdateConfig is an object to specified the update configuration of desktop.
type DesktopUpdateConfig struct {
	// Desktop ID.
	DesktopId string `json:"desktop_id"`
}

// UpdateProduct is a method to create a desktop using given parameters.
func UpdateProduct(c *golangsdk.ServiceClient, opts ProductUpdateOpts) ([]Job, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r UpdateResp
	_, err = c.Post(productURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.Jobs, err
}

// DeleteOpts is the structure required by the Delete method to delete an existing desktop.
type DeleteOpts struct {
	// Whether to delete user associated with this desktop after deleting it.
	DeleteUser bool `q:"delete_users"`
	// Whether to send emails to user mailbox during delete operation.
	EmailNotification bool `q:"email_notification"`
}

// Delete is a method to remove an existing desktop using given parameters, if the user does not have any desktop under
// it, the user can delete it together with the desktop.
func Delete(c *golangsdk.ServiceClient, desktopId string, opts DeleteOpts) error {
	url := resourceURL(c, desktopId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return err
	}
	url += query.String()

	_, err = c.Delete(url, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// NewVolumeOpts is the structure required by the NewVolumes method to add some volumes to the desktop.
type NewVolumeOpts struct {
	// New volumes parameters.
	VolumeConfigs []NewVolumeConfig `json:"addDesktopVolumesReq,omitempty"`
}

// NewVolumeConfig is an object to specified the volume configuration.
type NewVolumeConfig struct {
	// The desktop ID to which the volume belongs.
	DesktopId string `json:"desktop_id,omitempty"`
	//Configuration of data volumes.
	Volumes []Volume `json:"volumes,omitempty"`
}

// NewVolumes is a method to add some new volumes to the desktop.
func NewVolumes(c *golangsdk.ServiceClient, opts NewVolumeOpts) (*NewVolumesResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r NewVolumesResp
	_, err = c.Post(volumeURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// VolumeExpandOpts is the structure required by the ExpandVolumes method to batch expand volumes size.
type VolumeExpandOpts struct {
	// Volumes expansion parameters.
	VolumeConfigs []ExpandVolumeConfig `json:"expandVolumesReq,omitempty"`
}

// ExpandVolumeConfig is an object to specified the volume configuration.
type ExpandVolumeConfig struct {
	// The desktop ID to which the volume belongs.
	DesktopId string `json:"desktop_id,omitempty"`
	// Volume ID.
	VolumeId string `json:"volume_id,omitempty"`
	// The size of the disk after resizing, in GB.
	// For root volume, the valid value is range from 80 to 1020.
	// For data volume, the valid value is range from 10 to 8200.
	NewSize int `json:"new_size,omitempty"`
}

// ExpandVolumes is a method to batch expand the desktop volumes size.
func ExpandVolumes(c *golangsdk.ServiceClient, opts VolumeExpandOpts) (*ExpandVolumesResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r ExpandVolumesResp
	_, err = c.Post(volumeExpandURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// RebuildOpts is the structure that used to modify desktop image and os.
type RebuildOpts struct {
	// ID list of workspace desktops that wants to rebuild.
	DesktopIds []string `json:"desktop_ids" required:"true"`
	// New image type.
	ImageType string `json:"image_type" required:"true"`
	// New image ID.
	ImageId string `json:"image_id" required:"true"`
	// New OS type.
	OsType string `json:"os_type,omitempty"`
	// Delay time.
	DelayTime string `json:"delay_time,omitempty"`
	// Rebuild message send to the users.
	Message string `json:"message,omitempty"`
	// Enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
}

// Rebuild is the method that used to modify desktop using given parameters.
func Rebuild(c *golangsdk.ServiceClient, opts RebuildOpts) (*RebuildResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r RebuildResp
	_, err = c.Post(rebuildURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListEipOpts is the structure that used to query the bound desktop and unbound EIPs.
type ListEipOpts struct {
	// EnterpriseProject ID of desktop.
	EnterpriseProjectId string `q:"enterprise_project_id"`
	// Desktop ID.
	DesktopId string `q:"desktop_id"`
	// Desktop name.
	DesktopName string `q:"desktop_name"`
	// User name.
	UserName string `q:"user_name"`
	// EIP address.
	Address string `q:"address"`
	// Offset from which the query starts.
	// The starting record sequence number of the query starts from 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page.
	// If not specified, all matching records are returned.
	Limit int `q:"limit"`
	// EIP binding status.
	// + bind: binded EIP
	// + unbind: unbinded EIP
	State string `q:"state"`
}

// ListEips is the method that used to query the EIPs in which bound desktop and unbound desktop.
func ListEips(c *golangsdk.ServiceClient, desktopId string) ([]EipResp, error) {
	url := eipsURL(c)
	opts := ListEipOpts{
		DesktopId: desktopId,
	}
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pager := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := EipPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	})
	pager.Headers = requestOpts.MoreHeaders
	pages, err := pager.AllPages()

	if err != nil {
		return nil, err
	}

	return ExtractEips(pages)
}

// BindEipOpts is the structure that used to bind EIP to desktop.
type BindEipOpts struct {
	// ID list of workspace desktops that wants to apply EIP.
	DesktopId string `json:"desktop_id" required:"true"`
	// EIP ID of worksoaces that wants to apply.
	ID string `json:"eip_id" required:"true"`
}

// BindEip is the method that used to bind EIP to desktop.
func BindEip(c *golangsdk.ServiceClient, opts BindEipOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")

	if err != nil {
		return err
	}

	_, err = c.Post(bindEipURL(c, "binding"), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// UnbindEipOpt is the structure that used to unbind EIP from desktop.
type UnbindEipOpt struct {
	// Desktop ID of the EIP want to be unbind.
	DesktopIds []string `json:"desktop_ids" required:"true"`
}

// UnbindEip is the method that used to unbind EIP from desktop.
func UnbindEip(c *golangsdk.ServiceClient, opts UnbindEipOpt) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Post(bindEipURL(c, "unbinding"), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// GetNetwork is the method that used to query desktop network infomation.
func GetNetwork(c *golangsdk.ServiceClient, desktopId string) ([]NetworkInfos, error) {
	var r NetworkResp
	_, err := c.Get(networkURL(c, desktopId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.Network, err
}

// UpdateNetworkOpts is the structure that used to modify desktop network.
type UpdateNetworkOpts struct {
	DesktopId string `json:"-"`
	// The ID of the vpc to be change.
	VpcId string `json:"vpc_id" required:"true"`
	// The ID of the subnet to be change.
	SubnetId string `json:"subnet_id" required:"true"`
	// ID list of security group.
	SecurityGroupIds []string `json:"security_group_ids" required:"true"`
	// Specifies a private ID address.
	PrivateId string `json:"private_ip,omitempty"`
}

// UpdateNetwork is the method that used to modify desktop network infomation using given parameters.
func UpdateNetwork(c *golangsdk.ServiceClient, opts UpdateNetworkOpts) (*UpdateNetworkResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r UpdateNetworkResp
	_, err = c.Put(networkURL(c, opts.DesktopId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ActionOpts is the structure required by the DoAction method operate the power state of the desktop.
type ActionOpts struct {
	// ID list of workspace desktops that wants to operate.
	DesktopIds []string `json:"desktop_ids" required:"true"`
	// The power type of the desktop. The valid values are as follows:
	// + os-start
	// + os-stop
	// + reboot
	// + os-hibernate
	OpType string `json:"op_type" requires:"true"`
	// The operation type. The valid values are as follows:
	// + SOFT: Normal operation.
	// + HARD: Forced operation.
	Type string `json:"type,omitempty"`
}

// DoAction is a method that used to operate the power state of the desktop.
func DoAction(client *golangsdk.ServiceClient, opts ActionOpts) (*ActionResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r ActionResp
	_, err = client.Post(actionURL(client), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}
