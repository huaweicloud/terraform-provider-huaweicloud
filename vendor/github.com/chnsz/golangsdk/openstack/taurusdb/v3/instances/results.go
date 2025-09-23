package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/structs"
	"github.com/chnsz/golangsdk/pagination"
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
	Id                string   `json:"id"`
	Name              string   `json:"name"`
	Alias             string   `json:"alias"`
	Status            string   `json:"status"`
	Type              string   `json:"type"`
	Port              string   `json:"port"`
	NodeCount         int      `json:"node_count"`
	VpcId             string   `json:"vpc_id"`
	SubnetId          string   `json:"subnet_id"`
	SecurityGroupId   string   `json:"security_group_id"`
	ConfigurationId   string   `json:"configuration_id"`
	AZMode            string   `json:"az_mode"`
	MasterAZ          string   `json:"master_az_code"`
	TimeZone          string   `json:"time_zone"`
	ProjectId         string   `json:"project_id"`
	DbUserName        string   `json:"db_user_name"`
	PublicIps         string   `json:"public_ips"`
	PrivateDnsNames   []string `json:"private_dns_names"`
	PrivateIps        []string `json:"private_write_ips"`
	Created           string   `json:"created"`
	Updated           string   `json:"updated"`
	MaintenanceWindow string   `json:"maintenance_window"`

	Volume         Volume         `json:"volume"`
	Nodes          []Nodes        `json:"nodes"`
	DataStore      DataStore      `json:"datastore"`
	BackupStrategy BackupStrategy `json:"backup_strategy"`

	EnterpriseProjectId string `json:"enterprise_project_id"`
	DedicatedResourceId string `json:"dedicated_resource_id"`
}

type ListTaurusDBInstance struct {
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
	PublicIps       []string `json:"public_ips"`
	PrivateIps      []string `json:"private_write_ips"`
	Created         string   `json:"-"`
	Updated         string   `json:"-"`

	Volume         Volume         `json:"volume"`
	Nodes          []Nodes        `json:"nodes"`
	DataStore      DataStore      `json:"datastore"`
	BackupStrategy BackupStrategy `json:"backup_strategy"`

	EnterpriseProjectId string `json:"enterprise_project_id"`
	DedicatedResourceId string `json:"dedicated_resource_id"`
}

type Volume struct {
	Type string `json:"type"`
	Used string `json:"used"`
}

type NodeVolume struct {
	Size int `json:"size"`
}

type Nodes struct {
	Id               string     `json:"id"`
	Name             string     `json:"name"`
	Type             string     `json:"type"`
	Status           string     `json:"status"`
	PrivateIps       []string   `json:"private_read_ips"`
	Port             int        `json:"port"`
	Flavor           string     `json:"flavor_ref"`
	Region           string     `json:"region_code"`
	AvailabilityZone string     `json:"az_code"`
	Volume           NodeVolume `json:"volume"`
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

type ExtendResult struct {
	golangsdk.ErrResult
}

type ListTaurusDBResult struct {
	commonResult
}

type ListTaurusDBResponse struct {
	Instances  []ListTaurusDBInstance `json:"instances"`
	TotalCount int                    `json:"total_count"`
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

type Proxy struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	ElbVip  string `json:"elb_vip"`
	Eip     string `json:"eip"`
	NodeNum int    `json:"node_num"`
	Flavor  string `json:"flavor_ref"`
}

type GetProxyResult struct {
	commonResult
}

func (r GetProxyResult) Extract() (*Proxy, error) {
	var proxy Proxy
	err := r.ExtractIntoStructPtr(&proxy, "proxy")
	return &proxy, err
}

type DehResource struct {
	Id               string   `json:"id"`
	ResourceName     string   `json:"resource_name"`
	EngineName       string   `json:"engine_name"`
	AvailabilityZone []string `json:"availability_zone"`
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

type UpdateAliasResponse struct {
}

type UpdateAliasResult struct {
	commonResult
}

func (r UpdateAliasResult) ExtractUpdateAliasResponse() (*UpdateAliasResponse, error) {
	job := new(UpdateAliasResponse)
	err := r.ExtractInto(job)
	return job, err
}

type UpdateMaintenanceWindowResponse struct {
}

type UpdateMaintenanceWindowResult struct {
	commonResult
}

func (r UpdateMaintenanceWindowResult) ExtractUpdateMaintenanceWindowResponse() (*UpdateMaintenanceWindowResponse, error) {
	job := new(UpdateMaintenanceWindowResponse)
	err := r.ExtractInto(job)
	return job, err
}

type SecondLevelMonitoring struct {
	MonitorSwitch bool `json:"monitor_switch"`
	Period        int  `json:"period"`
}

type GetSecondLevelMonitoringResult struct {
	commonResult
}

func (r GetSecondLevelMonitoringResult) Extract() (*SecondLevelMonitoring, error) {
	var secondLevelMonitoring SecondLevelMonitoring
	err := r.ExtractInto(&secondLevelMonitoring)
	return &secondLevelMonitoring, err
}

type Version struct {
	UpgradeFlag bool      `json:"upgrade_flag"`
	Datastore   Datastore `json:"datastore"`
}

type Datastore struct {
	CurrentVersion       string `json:"current_version"`
	CurrentKernelVersion string `json:"current_kernel_version"`
	LatestVersion        string `json:"latest_version"`
	LatestKernelVersion  string `json:"latest_kernel_version"`
}

type GetVersionResult struct {
	commonResult
}

func (r GetVersionResult) Extract() (*Version, error) {
	var version Version
	err := r.ExtractInto(&version)
	return &version, err
}

type UpdateSlowLogShowOriginalSwitchResponse struct {
	Response string `json:"response"`
}

type UpdateSlowLogShowOriginalSwitchResult struct {
	commonResult
}

func (r UpdateSlowLogShowOriginalSwitchResult) ExtractUpdateSlowLogShowOriginalSwitchResponse() (*UpdateSlowLogShowOriginalSwitchResponse, error) {
	res := new(UpdateSlowLogShowOriginalSwitchResponse)
	err := r.ExtractInto(res)
	return res, err
}

type SlowLogShowOriginalSwitch struct {
	OpenSlowLogSwitch string `json:"open_slow_log_switch"`
}

type GetSlowLogShowOriginalSwitchResult struct {
	commonResult
}

func (r GetSlowLogShowOriginalSwitchResult) Extract() (*SlowLogShowOriginalSwitch, error) {
	var slowLogShowOriginalSwitch SlowLogShowOriginalSwitch
	err := r.ExtractInto(&slowLogShowOriginalSwitch)
	return &slowLogShowOriginalSwitch, err
}
