package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type GaussDBResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type CreateResponse struct {
	Instance GaussDBResponse `json:"instance"`
	OrderId  string          `json:"order_id"`
}

type ChargeInfoResp struct {
	ChargeMode string `json:"charge_mode"`
}

type UpdateResponse struct {
	JobId   string `json:"job_id"`
	OrderId string `json:"order_id"`
}

type GaussDBInstance struct {
	Id                string   `json:"id"`
	Name              string   `json:"name"`
	Status            string   `json:"status"`
	Type              string   `json:"type"`
	Port              int      `json:"port"`
	VpcId             string   `json:"vpc_id"`
	SubnetId          string   `json:"subnet_id"`
	SecurityGroupId   string   `json:"security_group_id"`
	TimeZone          string   `json:"time_zone"`
	Region            string   `json:"region"`
	FlavorRef         string   `json:"flavor_ref"`
	DbUserName        string   `json:"db_user_name"`
	DiskEncryptionId  string   `json:"disk_encryption_id"`
	DsspoolId         string   `json:"dsspool_id"`
	SwitchStrategy    string   `json:"switch_strategy"`
	MaintenanceWindow string   `json:"maintenance_window"`
	PublicIps         []string `json:"public_ips"`
	PrivateIps        []string `json:"private_ips"`
	ReplicaNum        int      `json:"replica_num"`

	MysqlCompatibility MysqlCompatibility `json:"mysql_compatibility"`
	Volume             VolumeOpt          `json:"volume"`
	Ha                 HaOpt              `json:"ha"`
	Nodes              []Nodes            `json:"nodes"`
	DataStore          DataStoreOpt       `json:"datastore"`
	BackupStrategy     BackupStrategyOpt  `json:"backup_strategy"`
	ChargeInfo         ChargeInfoResp     `json:"charge_info"`

	EnterpriseProjectId string `json:"enterprise_project_id"`
}

type Nodes struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Status           string `json:"status"`
	Role             string `json:"role"`
	AvailabilityZone string `json:"availability_zone"`
}

type MysqlCompatibility struct {
	Port string `json:"port"`
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type UpdateResult struct {
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

type ListGaussDBResult struct {
	commonResult
}

type ListGaussDBResponse struct {
	Instances  []GaussDBInstance `json:"instances"`
	TotalCount int               `json:"total_count"`
}

type GaussDBPage struct {
	pagination.SinglePageBase
}

func (r GaussDBPage) IsEmpty() (bool, error) {
	data, err := ExtractGaussDBInstances(r)
	if err != nil {
		return false, err
	}
	return len(data.Instances) == 0, err
}

// ExtractGaussDBInstances is a function that takes a ListResult and returns the services' information.
func ExtractGaussDBInstances(r pagination.Page) (ListGaussDBResponse, error) {
	var s ListGaussDBResponse
	err := (r.(GaussDBPage)).ExtractInto(&s)
	return s, err
}

type RenameResponse struct {
	Instance GaussDBResponse `json:"instance"`
	JobId    string          `json:"job_id"`
}

type RenameResult struct {
	commonResult
}

func (r RenameResult) Extract() (*RenameResponse, error) {
	var response RenameResponse
	err := r.ExtractInto(&response)
	return &response, err
}
