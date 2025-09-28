package cloudservers

import (
	"strconv"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type cloudServerResult struct {
	golangsdk.Result
}

type PasswordResult struct {
	golangsdk.ErrResult
}

type UpdateResult struct {
	golangsdk.ErrResult
}
type Flavor struct {
	Disk  string `json:"disk"`
	Vcpus string `json:"vcpus"`
	RAM   string `json:"ram"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

// Image defines a image struct in details of a server.
type Image struct {
	ID string `json:"id"`
}

type SysTags struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OsSchedulerHints struct {
	Group           []string `json:"group"`
	Tenancy         []string `json:"tenancy"`
	DedicatedHostID []string `json:"dedicated_host_id"`
	FaultDomain     string   `json:"fault_domain,omitempty"`
}

// Metadata is only used for method that requests details on a single server, by ID.
// Because metadata struct must be a map.
type Metadata struct {
	ChargingMode      string `json:"charging_mode"`
	OrderID           string `json:"metering.order_id"`
	ProductID         string `json:"metering.product_id"`
	VpcID             string `json:"vpc_id"`
	EcmResStatus      string `json:"EcmResStatus"`
	ImageID           string `json:"metering.image_id"`
	Imagetype         string `json:"metering.imagetype"`
	Resourcespeccode  string `json:"metering.resourcespeccode"`
	ImageName         string `json:"image_name"`
	OsBit             string `json:"os_bit"`
	LockCheckEndpoint string `json:"lock_check_endpoint"`
	LockSource        string `json:"lock_source"`
	LockSourceID      string `json:"lock_source_id"`
	LockScene         string `json:"lock_scene"`
	VirtualEnvType    string `json:"virtual_env_type"`
	AgencyName        string `json:"agency_name"`
	AgentList         string `json:"__support_agent_list"`
}

type Address struct {
	Version string `json:"version"`
	Addr    string `json:"addr"`
	MacAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
	PortID  string `json:"OS-EXT-IPS:port_id"`
	Type    string `json:"OS-EXT-IPS:type"`
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
	Status              string               `json:"status"`
	Updated             time.Time            `json:"updated"`
	HostID              string               `json:"hostId"`
	Addresses           map[string][]Address `json:"addresses"`
	ID                  string               `json:"id"`
	Name                string               `json:"name"`
	AccessIPv4          string               `json:"accessIPv4"`
	AccessIPv6          string               `json:"accessIPv6"`
	Created             time.Time            `json:"created"`
	Tags                []string             `json:"tags"`
	Description         string               `json:"description"`
	Locked              *bool                `json:"locked"`
	ConfigDrive         string               `json:"config_drive"`
	TenantID            string               `json:"tenant_id"`
	UserID              string               `json:"user_id"`
	HostStatus          string               `json:"host_status"`
	EnterpriseProjectID string               `json:"enterprise_project_id"`
	SysTags             []SysTags            `json:"sys_tags"`
	Flavor              Flavor               `json:"flavor"`
	Metadata            Metadata             `json:"metadata"`
	SecurityGroups      []SecurityGroups     `json:"security_groups"`
	KeyName             string               `json:"key_name"`
	Image               Image                `json:"image"`
	Progress            *int                 `json:"progress"`
	PowerState          *int                 `json:"OS-EXT-STS:power_state"`
	VMState             string               `json:"OS-EXT-STS:vm_state"`
	TaskState           string               `json:"OS-EXT-STS:task_state"`
	DiskConfig          string               `json:"OS-DCF:diskConfig"`
	AvailabilityZone    string               `json:"OS-EXT-AZ:availability_zone"`
	LaunchedAt          string               `json:"OS-SRV-USG:launched_at"`
	TerminatedAt        string               `json:"OS-SRV-USG:terminated_at"`
	RootDeviceName      string               `json:"OS-EXT-SRV-ATTR:root_device_name"`
	RamdiskID           string               `json:"OS-EXT-SRV-ATTR:ramdisk_id"`
	KernelID            string               `json:"OS-EXT-SRV-ATTR:kernel_id"`
	LaunchIndex         *int                 `json:"OS-EXT-SRV-ATTR:launch_index"`
	ReservationID       string               `json:"OS-EXT-SRV-ATTR:reservation_id"`
	Hostname            string               `json:"OS-EXT-SRV-ATTR:hostname"`
	UserData            string               `json:"OS-EXT-SRV-ATTR:user_data"`
	Host                string               `json:"OS-EXT-SRV-ATTR:host"`
	InstanceName        string               `json:"OS-EXT-SRV-ATTR:instance_name"`
	HypervisorHostname  string               `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
	VolumeAttached      []VolumeAttached     `json:"os-extended-volumes:volumes_attached"`
	OsSchedulerHints    OsSchedulerHints     `json:"os:scheduler_hints"`
	Fault               Fault                `json:"fault"`
	AutoTerminateTime   string               `json:"auto_terminate_time"`
	EnclaveOptions      *EnclaveOptions      `json:"enclave_options"`
}

// ECS fault causes
type Fault struct {
	Code    int       `json:"code"`
	Created time.Time `json:"created"`
	Details string    `json:"details"`
	Message string    `json:"message"`
}

// NewCloudServer defines the response from details on a single server, by ID.
type NewCloudServer struct {
	CloudServer
	Metadata map[string]string `json:"metadata"`
}

// GetResult is the response from a Get operation. Call its Extract
// method to interpret it as a Server.
type GetResult struct {
	cloudServerResult
}

func (r GetResult) Extract() (*CloudServer, error) {
	var s struct {
		Server *CloudServer `json:"server"`
	}
	err := r.ExtractInto(&s)
	return s.Server, err
}

// ServerPage abstracts the raw results of making a List() request against
// the API.
type ServerPage struct {
	pagination.PageSizeBase
}

// IsEmpty returns true if a page contains no Server results.
func (r ServerPage) IsEmpty() (bool, error) {
	s, err := ExtractServers(r)
	return len(s) == 0, err
}

// NextPageURL returns the next page with offset and limit.
func (r ServerPage) NextPageURL() (string, error) {
	pageName := "offset"
	currentURL := r.URL

	q := currentURL.Query()
	pageNum := q.Get(pageName)
	if pageNum == "" {
		pageNum = "1"
	}

	sizeVal, err := strconv.ParseInt(pageNum, 10, 32)
	if err != nil {
		return "", err
	}

	pageNum = strconv.Itoa(int(sizeVal + 1))
	q.Set(pageName, pageNum)
	currentURL.RawQuery = q.Encode()
	return currentURL.String(), nil
}

// ExtractServers interprets the results of a single page from a List() call,
// producing a slice of CloudServer entities.
func ExtractServers(r pagination.Page) ([]CloudServer, error) {
	var s struct {
		Servers []CloudServer `json:"servers"`
	}

	err := (r.(ServerPage)).ExtractInto(&s)
	return s.Servers, err
}

// UpdateMetadataResult contains the result of an UpdateMetadata operation.
// Call its Extract method to interpret it as a map[string]interface{}.
type UpdateMetadataResult struct {
	golangsdk.Result
}

// DeleteMetadatItemResult contains the result of a DeleteMetadatItem operation.
// Call its ExtractErr method to determine if the call succeeded or failed.
type DeleteMetadatItemResult struct {
	golangsdk.ErrResult
}

func (r UpdateMetadataResult) Extract() (map[string]interface{}, error) {
	var s struct {
		Metadata map[string]interface{} `json:"metadata"`
	}
	err := r.ExtractInto(&s)
	return s.Metadata, err
}
