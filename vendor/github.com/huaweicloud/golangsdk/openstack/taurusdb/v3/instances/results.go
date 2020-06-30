package instances

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/structs"
	"github.com/huaweicloud/golangsdk/pagination"
)

type DataStore struct {
	Type    string `json:"type" required:"true"`
	Version string `json:"version" required:"true"`
}

type BackupStrategy struct {
	StartTime string `json:"start_time"`
	KeepDays  string `json:"keep_days"`
}

type TaurusDBResponse struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	Region          string `json:"region"`
	Mode            string `json:"mode"`
	Port            string `json:"port"`
	VpcId           string `json:"vpc_id"`
	SubnetId        string `json:"subnet_id"`
	SecurityGroupId string `json:"security_group_id"`
	ConfigurationId string `json:"configuration_id"`
	AZMode          string `json:"availability_zone_mode"`
	MasterAZ        string `json:"master_availability_zone"`
	SlaveCount      int    `json:"slave_count"`

	DataStore      DataStore          `json:"datastore"`
	BackupStrategy BackupStrategy     `json:"backup_strategy"`
	ChargeInfo     structs.ChargeInfo `json:"charge_info"`

	EnterpriseProjectId string `json:"enterprise_project_id"`
}

type CreateResponse struct {
	Instance TaurusDBResponse `json:"instance"`
	JobId    string           `json:"job_id"`
	OrderId  string           `json:"order_id"`
}

type TaurusDBInstance struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	Type            string   `json:"type"`
	Port            string   `json:"port"`
	NodeCount       int      `json:"node_count"`
	VpcId           string   `json:"vpc_id"`
	SubnetId        string   `json:"subnet_id"`
	SecurityGroupId string   `json:"security_group_id"`
	ConfigurationId string   `json:"configuration_id"`
	AZMode          string   `json:"az_mode"`
	MasterAZ        string   `json:"master_az_code"`
	TimeZone        string   `json:"time_zone"`
	ProjectId       string   `json:"project_id"`
	DbUserName      string   `json:"db_user_name"`
	PublicIps       string   `json:"public_ips"`
	PrivateIps      []string `json:"private_write_ips"`
	Created         string   `json:"-"`
	Updated         string   `json:"-"`

	Volume         Volume         `json:"volume"`
	Nodes          []Nodes        `json:"nodes"`
	DataStore      DataStore      `json:"datastore"`
	BackupStrategy BackupStrategy `json:"backup_strategy"`

	EnterpriseProjectId string `json:"enterprise_project_id"`
}

type Volume struct {
	Type string `json:"type"`
	Used string `json:"used"`
}

type Nodes struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	Type             string   `json:"type"`
	Status           string   `json:"status"`
	PrivateIps       []string `json:"private_read_ips"`
	Port             int      `json:"port"`
	Flavor           string   `json:"flavor_ref"`
	Region           string   `json:"region_code"`
	AvailabilityZone string   `json:"az_code"`
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

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*TaurusDBInstance, error) {
	var instance TaurusDBInstance
	err := r.ExtractIntoStructPtr(&instance, "instance")
	return &instance, err
}

type ListTaurusDBResult struct {
	commonResult
}

type ListTaurusDBResponse struct {
	Instances  []TaurusDBInstance `json:"instances"`
	TotalCount int                `json:"total_count"`
}

type TaurusDBPage struct {
	pagination.SinglePageBase
}

func (r TaurusDBPage) IsEmpty() (bool, error) {
	data, err := ExtractTaurusDBInstances(r)
	if err != nil {
		return false, err
	}
	return len(data.Instances) == 0, err
}

// ExtractTaurusDBInstances is a function that takes a ListResult and returns the services' information.
func ExtractTaurusDBInstances(r pagination.Page) (ListTaurusDBResponse, error) {
	var s ListTaurusDBResponse
	err := (r.(TaurusDBPage)).ExtractInto(&s)
	return s, err
}
