package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

// InstanceCreate response
type InstanceCreate struct {
	InstanceID string `json:"instance_id"`
}

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

// Extract from CreateResult
func (r CreateResult) Extract() (*InstanceCreate, error) {
	var s InstanceCreate
	err := r.Result.ExtractInto(&s)
	return &s, err
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.ErrResult
}

type ListResponse struct {
	Instances  []Instance `json:"instances"`
	TotalCount int        `json:"instance_num"`
}

// Instance response
type Instance struct {
	Name                     string             `json:"name"`
	Description              string             `json:"description"`
	Engine                   string             `json:"engine"`
	EngineVersion            string             `json:"engine_version"`
	Specification            string             `json:"specification"`
	StorageSpace             int                `json:"storage_space"`
	UsedStorageSpace         int                `json:"used_storage_space"`
	ConnectAddress           string             `json:"connect_address"`
	Port                     int                `json:"port"`
	Status                   string             `json:"status"`
	InstanceID               string             `json:"instance_id"`
	ResourceSpecCode         string             `json:"resource_spec_code"`
	BrokerNum                int                `json:"broker_num"`
	ChargingMode             int                `json:"charging_mode"`
	VPCID                    string             `json:"vpc_id"`
	VPCName                  string             `json:"vpc_name"`
	CreatedAt                string             `json:"created_at"`
	UserID                   string             `json:"user_id"`
	UserName                 string             `json:"user_name"`
	OrderID                  string             `json:"order_id"`
	MaintainBegin            string             `json:"maintain_begin"`
	MaintainEnd              string             `json:"maintain_end"`
	EnablePublicIP           bool               `json:"enable_publicip"`
	PublicIPAddress          string             `json:"publicip_address"`
	PublicIPID               string             `json:"publicip_id"`
	ManagementConnectAddress string             `json:"management_connect_address"`
	SslEnable                bool               `json:"ssl_enable"`
	EnterpriseProjectID      string             `json:"enterprise_project_id"`
	IsLogicalVolume          bool               `json:"is_logical_volume"`
	ExtendTimes              int                `json:"extend_times"`
	Type                     string             `json:"type"`
	ProductID                string             `json:"product_id"`
	SecurityGroupID          string             `json:"security_group_id"`
	SecurityGroupName        string             `json:"security_group_name"`
	SubnetID                 string             `json:"subnet_id"`
	AvailableZones           []string           `json:"available_zones"`
	TotalStorageSpace        int                `json:"total_storage_space"`
	StorageResourceID        string             `json:"storage_resource_id"`
	StorageSpecCode          string             `json:"storage_spec_code"`
	Ipv6Enable               bool               `json:"ipv6_enable"`
	Ipv6ConnectAddresses     []string           `json:"ipv6_connect_addresses"`
	ServiceType              string             `json:"service_type"`
	StorageType              string             `json:"storage_type"`
	PublicBoundWidth         int                `json:"public_boundwidth"`
	EnableLogCollect         bool               `json:"enable_log_collection"`
	ConnectorEnalbe          bool               `json:"connector_enable"`
	ConnectorID              string             `json:"connector_id"`
	RestEnable               bool               `json:"rest_enable"`
	RestConnectAddress       string             `json:"rest_connect_address"`
	MessageQueryInstEnable   bool               `json:"message_query_inst_enable"`
	VpcClientPlain           bool               `json:"vpc_client_plain"`
	SupportFeatures          string             `json:"support_features"`
	TraceEnable              bool               `json:"trace_enable"`
	DiskEncrypted            bool               `json:"disk_encrypted"`
	CesVersion               string             `json:"ces_version"`
	AccessUser               string             `json:"access_user"`
	Task                     Task               `json:"task"`
	Tags                     []tags.ResourceTag `json:"tags"`
	EnableAcl                bool               `json:"enable_acl"`
}

type Task struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// UpdateResult is a struct from which can get the result of update method
type UpdateResult struct {
	golangsdk.Result
}

// GetResult contains the body of getting detailed
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*Instance, error) {
	var s Instance
	err := r.Result.ExtractInto(&s)
	return &s, err
}

type Page struct {
	pagination.SinglePageBase
}

func (r Page) IsEmpty() (bool, error) {
	data, err := ExtractInstances(r)
	if err != nil {
		return false, err
	}
	return len(data.Instances) == 0, err
}

// ExtractCloudServers is a function that takes a ListResult and returns the services' information.
func ExtractInstances(r pagination.Page) (ListResponse, error) {
	var s ListResponse
	err := (r.(Page)).ExtractInto(&s)
	return s, err
}

// ResetPasswordResult is a struct that contains all the return parameters of ResetPassword function
type ResetPasswordResult struct {
	golangsdk.Result
}
