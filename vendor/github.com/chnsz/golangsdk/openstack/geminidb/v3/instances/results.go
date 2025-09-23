package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type DataStore struct {
	Type          string `json:"type" required:"true"`
	Version       string `json:"version"`
	StorageEngine string `json:"storage_engine" required:"true"`
}

type Flavor struct {
	Num      string `json:"num"`
	Size     string `json:"size"`
	Storage  string `json:"storage"`
	SpecCode string `json:"spec_code"`
}

type BackupStrategy struct {
	StartTime string `json:"start_time"`
	KeepDays  int    `json:"keep_days"`
}

type GeminiDBBase struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	Region          string    `json:"region"`
	Mode            string    `json:"mode"`
	Created         string    `json:"-"`
	VpcId           string    `json:"vpc_id"`
	SubnetId        string    `json:"subnet_id"`
	SecurityGroupId string    `json:"security_group_id"`
	DataStore       DataStore `json:"datastore"`

	EnterpriseProjectId string `json:"enterprise_project_id"`
}

type CreateResponse struct {
	GeminiDBBase
	JobId            string            `json:"job_id"`
	OrderId          string            `json:"order_id"`
	AvailabilityZone string            `json:"availability_zone"`
	Flavor           []Flavor          `json:"flavor"`
	BackupStrategy   BackupStrategyOpt `json:"backup_strategy"`
}

type GeminiDBInstance struct {
	GeminiDBBase
	Port              string         `json:"port"`
	Engine            string         `json:"engine"`
	Updated           string         `json:"-"`
	DbUserName        string         `json:"db_user_name"`
	PayMode           string         `json:"pay_mode"`
	TimeZone          string         `json:"time_zone"`
	MaintenanceWindow string         `json:"maintenance_window"`
	LbIpAddress       string         `json:"lb_ip_address"`
	LbPort            string         `json:"lb_port"`
	Actions           []string       `json:"actions"`
	Groups            []Groups       `json:"groups"`
	BackupStrategy    BackupStrategy `json:"backup_strategy"`

	DedicatedResourceId string `json:"dedicated_resource_id"`
}

type Groups struct {
	Id     string  `json:"id"`
	Status string  `json:"status"`
	Volume Volume  `json:"volume"`
	Nodes  []Nodes `json:"nodes"`
}

type Volume struct {
	Size string `json:"size"`
	Used string `json:"used"`
}

type Nodes struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Status           string `json:"status"`
	PrivateIp        string `json:"private_ip"`
	PublicIp         string `json:"public_ip"`
	SpecCode         string `json:"spec_code"`
	AvailabilityZone string `json:"availability_zone"`
	SupportReduce    bool   `json:"support_reduce"`
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*CreateResponse, error) {
	var response CreateResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type DeleteResult struct {
	commonResult
}

type DeleteResponse struct {
	JobId string `json:"job_id"`
}

func (r DeleteResult) Extract() (*DeleteResponse, error) {
	var response DeleteResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type ExtendResult struct {
	commonResult
}

type ExtendResponse struct {
	JobId   string `json:"job_id"`
	OrderId string `json:"order_id"`
}

func (r ExtendResult) Extract() (*ExtendResponse, error) {
	var response ExtendResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type UpdateResult struct {
	golangsdk.ErrResult
}

type ListGeminiDBResult struct {
	commonResult
}

type ListGeminiDBResponse struct {
	Instances  []GeminiDBInstance `json:"instances"`
	TotalCount int                `json:"total_count"`
}

type GeminiDBPage struct {
	pagination.SinglePageBase
}

func (r GeminiDBPage) IsEmpty() (bool, error) {
	data, err := ExtractGeminiDBInstances(r)
	if err != nil {
		return false, err
	}
	return len(data.Instances) == 0, err
}

// ExtractGeminiDBInstances is a function that takes a ListResult and returns the services' information.
func ExtractGeminiDBInstances(r pagination.Page) (ListGeminiDBResponse, error) {
	var s ListGeminiDBResponse
	err := (r.(GeminiDBPage)).ExtractInto(&s)
	return s, err
}

type DehResource struct {
	Id               string   `json:"id"`
	ResourceName     string   `json:"resource_name"`
	EngineName       string   `json:"engine_name"`
	AvailabilityZone string   `json:"availability_zone"`
	Architecture     string   `json:"architecture"`
	Status           string   `json:"status"`
	Capacity         Capacity `json:"capacity"`
}

type Capacity struct {
	Vcpus  int   `json:"vcpus"`
	Ram    int   `json:"ram"`
	Volume int64 `json:"volume"`
}

type ListDehResponse struct {
	Resources  []DehResource `json:"resources"`
	TotalCount int           `json:"total_count"`
}

type DehResourcePage struct {
	pagination.SinglePageBase
}

func ExtractDehResources(r pagination.Page) (ListDehResponse, error) {
	var s ListDehResponse
	err := (r.(DehResourcePage)).ExtractInto(&s)
	return s, err
}

type SslResult struct {
	commonResult
}

type SslResponse struct {
	JobId string `json:"job_id"`
}

func (r SslResult) Extract() (*SslResponse, error) {
	var response SslResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type PublicIpResult struct {
	commonResult
}

type PublicIpResponse struct {
	JobId string `json:"job_id"`
}

func (r PublicIpResult) Extract() (*PublicIpResponse, error) {
	var response PublicIpResponse
	err := r.ExtractInto(&response)
	return &response, err
}
