package instances

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
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

type ListDcsResponse struct {
	Instances  []Instance `json:"instances"`
	TotalCount int        `json:"instance_num"`
}

// Instance response
type Instance struct {
	Name                  string               `json:"name"`
	Engine                string               `json:"engine"`
	Capacity              int                  `json:"capacity"`
	CapacityMinor         string               `json:"capacity_minor"`
	IP                    string               `json:"ip"`
	Port                  int                  `json:"port"`
	Status                string               `json:"status"`
	Description           string               `json:"description"`
	InstanceID            string               `json:"instance_id"`
	ResourceSpecCode      string               `json:"resource_spec_code"`
	EngineVersion         string               `json:"engine_version"`
	InternalVersion       string               `json:"internal_version"`
	ChargingMode          int                  `json:"charging_mode"`
	VPCID                 string               `json:"vpc_id"`
	VPCName               string               `json:"vpc_name"`
	CreatedAt             string               `json:"created_at"`
	ErrorCode             string               `json:"error_code"`
	ProductID             string               `json:"product_id"`
	SecurityGroupID       string               `json:"security_group_id"`
	SecurityGroupName     string               `json:"security_group_name"`
	SubnetID              string               `json:"subnet_id"`
	SubnetName            string               `json:"subnet_name"`
	SubnetCIDR            string               `json:"subnet_cidr"`
	AvailableZones        []string             `json:"available_zones"`
	MaxMemory             int                  `json:"max_memory"`
	UsedMemory            int                  `json:"used_memory"`
	InstanceBackupPolicy  InstanceBackupPolicy `json:"instance_backup_policy"`
	UserID                string               `json:"user_id"`
	UserName              string               `json:"user_name"`
	OrderID               string               `json:"order_id"`
	MaintainBegin         string               `json:"maintain_begin"`
	MaintainEnd           string               `json:"maintain_end"`
	NoPasswordAccess      string               `json:"no_password_access"`
	AccessUser            string               `json:"access_user"`
	EnterpriseProjectID   string               `json:"enterprise_project_id"`
	EnterpriseProjectName string               `json:"enterprise_project_name"`
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

// Password response
type Password struct {
	// Whether the password is successfully changed:
	// Values:
	// Success: The password is successfully changed.
	// passwordFailed: The old password is incorrect.
	// Locked: This account has been locked.
	// Failed: Failed to change the password.
	Result         string `json:"result"`
	Message        string `json:"message"`
	RetryTimesLeft string `json:"retry_times_left"`
	LockTime       string `json:"lock_time"`
	LockTimesLeft  string `json:"lock_time_left"`
}

// UpdatePasswordResult is a struct from which can get the result of update password method
type UpdatePasswordResult struct {
	golangsdk.Result
}

// Extract from UpdatePasswordResult
func (r UpdatePasswordResult) Extract() (*Password, error) {
	var s Password
	err := r.Result.ExtractInto(&s)
	return &s, err
}

// ExtendResult is a struct from which can get the result of extend method
type ExtendResult struct {
	golangsdk.Result
}

type DcsPage struct {
	pagination.SinglePageBase
}

func (r DcsPage) IsEmpty() (bool, error) {
	data, err := ExtractDcsInstances(r)
	if err != nil {
		return false, err
	}
	return len(data.Instances) == 0, err
}

// ExtractCloudServers is a function that takes a ListResult and returns the services' information.
func ExtractDcsInstances(r pagination.Page) (ListDcsResponse, error) {
	var s ListDcsResponse
	err := (r.(DcsPage)).ExtractInto(&s)
	return s, err
}
