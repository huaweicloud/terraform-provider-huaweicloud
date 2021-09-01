package baremetalservers

import (
	"github.com/chnsz/golangsdk"
)

type cloudServerResult struct {
	golangsdk.Result
}

type Flavor struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Disk  string `json:"disk"`
	Vcpus string `json:"vcpus"`
	RAM   string `json:"ram"`
}

// Image defines a image struct in details of a server.
type Image struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"__os_type"`
}

type SysTags struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OsSchedulerHints struct {
	DecBaremetal []string `json:"dec_baremetal"`
}

// Metadata is only used for method that requests details on a single server, by ID.
// Because metadata struct must be a map.
type Metadata struct {
	ChargingMode     string `json:"charging_mode"`
	OrderID          string `json:"metering.order_id"`
	ProductID        string `json:"metering.product_id"`
	VpcID            string `json:"vpc_id"`
	ImageID          string `json:"metering.image_id"`
	Imagetype        string `json:"metering.imagetype"`
	PortList         string `json:"baremetalPortIDList"`
	Resourcespeccode string `json:"metering.resourcespeccode"`
	ResourceType     string `json:"metering.resourcetype"`
	ImageName        string `json:"image_name"`
	OpSvcUserId      string `json:"op_svc_userid"`
	OsType           string `json:"os_type"`
	BmsSupportEvs    string `json:"__bms_support_evs"`
	OsBit            string `json:"os_bit"`
}

type Address struct {
	Version string `json:"version"`
	Addr    string `json:"addr"`
	MacAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
	PortID  string `json:"OS-EXT-IPS:port_id"`
	Type    string `json:"OS-EXT-IPS:type"`
}

type Fault struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Details string `json:"details"`
	Created string `json:"created"`
}
type VolumeAttached struct {
	ID                  string `json:"id"`
	DeleteOnTermination string `json:"delete_on_termination"`
	BootIndex           string `json:"bootIndex"`
	Device              string `json:"device"`
}

type SecurityGroups struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// CloudServer is only used for method that requests details on a single server, by ID.
// Because metadata struct must be a map.
type CloudServer struct {
	ID                  string               `json:"id"`
	UserID              string               `json:"user_id"`
	Name                string               `json:"name"`
	TenantID            string               `json:"tenant_id"`
	HostID              string               `json:"hostId"`
	Addresses           map[string][]Address `json:"addresses"`
	KeyName             string               `json:"key_name"`
	Image               Image                `json:"image"`
	Flavor              Flavor               `json:"flavor"`
	SecurityGroups      []SecurityGroups     `json:"security_groups"`
	AccessIPv4          string               `json:"accessIPv4"`
	AccessIPv6          string               `json:"accessIPv6"`
	Status              string               `json:"status"`
	Progress            *int                 `json:"progress"`
	ConfigDrive         string               `json:"config_drive"`
	Metadata            Metadata             `json:"metadata"`
	TaskState           string               `json:"OS-EXT-STS:task_state"`
	VMState             string               `json:"OS-EXT-STS:vm_state"`
	Host                string               `json:"OS-EXT-SRV-ATTR:host"`
	InstanceName        string               `json:"OS-EXT-SRV-ATTR:instance_name"`
	PowerState          *int                 `json:"OS-EXT-STS:power_state"`
	HypervisorHostname  string               `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
	AvailabilityZone    string               `json:"OS-EXT-AZ:availability_zone"`
	DiskConfig          string               `json:"OS-DCF:diskConfig"`
	Fault               Fault                `json:"fault"`
	VolumeAttached      []VolumeAttached     `json:"os-extended-volumes:volumes_attached"`
	Description         string               `json:"description"`
	HostStatus          string               `json:"host_status"`
	Hostname            string               `json:"OS-EXT-SRV-ATTR:hostname"`
	ReservationID       string               `json:"OS-EXT-SRV-ATTR:reservation_id"`
	LaunchIndex         *int                 `json:"OS-EXT-SRV-ATTR:launch_index"`
	KernelID            string               `json:"OS-EXT-SRV-ATTR:kernel_id"`
	RamdiskID           string               `json:"OS-EXT-SRV-ATTR:ramdisk_id"`
	RootDeviceName      string               `json:"OS-EXT-SRV-ATTR:root_device_name"`
	UserData            string               `json:"OS-EXT-SRV-ATTR:user_data"`
	Locked              *bool                `json:"locked"`
	Tags                []string             `json:"tags"`
	OsSchedulerHints    OsSchedulerHints     `json:"os:scheduler_hints"`
	EnterpriseProjectID string               `json:"enterprise_project_id"`
	SysTags             []SysTags            `json:"sys_tags"`
}

// GetResult is the response from a Get operation. Call its Extract
// method to interpret it as a Server.
type GetResult struct {
	cloudServerResult
}

type UpdateResult struct {
	cloudServerResult
}

func (r cloudServerResult) Extract() (*CloudServer, error) {
	var s struct {
		Server *CloudServer `json:"server"`
	}
	err := r.ExtractInto(&s)
	return s.Server, err
}
