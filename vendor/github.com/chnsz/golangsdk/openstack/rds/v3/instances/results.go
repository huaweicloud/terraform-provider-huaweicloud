package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type RenameResult struct {
	commonResult
}

type ModifyAliasResult struct {
	commonResult
}

type ModifyMaintainWindowResult struct {
	commonResult
}

type ModifyReplicationModeResult struct {
	commonResult
}

type ModifyBinlogRetentionHoursResult struct {
	commonResult
}

type ModifySwitchStrategyResult struct {
	commonResult
}

type ResizeFlavorResult struct {
	commonResult
}

type EnlargeVolumeResult struct {
	commonResult
}

type ApplyConfigurationOptsResult struct {
	commonResult
}

type ModifyConfigurationResult struct {
	commonResult
}

type ModifySecondLevelMonitoringResult struct {
	commonResult
}

type GetConfigurationResult struct {
	commonResult
}

type GetBinlogRetentionHoursResult struct {
	commonResult
}

type GetTdeStatusResult struct {
	commonResult
}

type GetSecondLevelMonitoringResult struct {
	commonResult
}

type ModifySlowLogShowOriginalStatusResult struct {
	commonResult
}

type JobResult struct {
	commonResult
}

type Instance struct {
	Id                  string         `json:"id"`
	Name                string         `json:"name"`
	Status              string         `json:"status"`
	Datastore           Datastore      `json:"datastore"`
	Ha                  Ha             `json:"ha"`
	ConfigurationId     string         `json:"configuration_id"`
	Port                string         `json:"port"`
	BackupStrategy      BackupStrategy `json:"backup_strategy"`
	EnterpriseProjectId string         `json:"enterprise_project_id"`
	DiskEncryptionId    string         `json:"disk_encryption_id"`
	FlavorRef           string         `json:"flavor_ref"`
	Volume              Volume         `json:"volume"`
	Region              string         `json:"region"`
	AvailabilityZone    string         `json:"availability_zone"`
	VpcId               string         `json:"vpc_id"`
	SubnetId            string         `json:"subnet_id"`
	SecurityGroupId     string         `json:"security_group_id"`
	ChargeInfo          ChargeResponse `json:"charge_info"`
}

type ChargeResponse struct {
	ChargeMode string `json:"charge_mode"`
}

type CreateResponse struct {
	Instance Instance `json:"instance"`
	JobId    string   `json:"job_id"`
	OrderId  string   `json:"order_id"`
}

func (r CreateResult) Extract() (*CreateResponse, error) {
	var response CreateResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type JobResponse struct {
	JobId     string `json:"job_id"`
	HumpJobId string `json:"jobId"`
}

func (r JobResult) Extract() (*JobResponse, error) {
	var response JobResponse
	err := r.ExtractInto(&response)
	return &response, err
}

func (job JobResponse) GetJobId() string {
	if job.JobId != "" {
		return job.JobId
	}
	return job.HumpJobId
}

type ResizeFlavor struct {
	JobId   string `json:"job_id"`
	OrderId string `json:"order_id"`
}

func (r ResizeFlavorResult) Extract() (*ResizeFlavor, error) {
	var response ResizeFlavor
	err := r.ExtractInto(&response)
	return &response, err
}

type EnlargeVolumeResp struct {
	JobId   string `json:"job_id"`
	OrderId string `json:"order_id"`
}

func (r EnlargeVolumeResult) Extract() (*EnlargeVolumeResp, error) {
	var response EnlargeVolumeResp
	err := r.ExtractInto(&response)
	return &response, err
}

type ApplyConfigurationResp struct {
	ConfigurationId   string `json:"configuration_id"`
	ConfigurationName string `json:"configuration_name"`
	Success           bool   `json:"success"`
	JobId             string `json:"job_id"`
}

func (r ApplyConfigurationOptsResult) Extract() (*ApplyConfigurationResp, error) {
	var response ApplyConfigurationResp
	err := r.ExtractInto(&response)
	return &response, err
}

type ModifyConfigurationResp struct {
	JobId   string `json:"job_id"`
	Restart bool   `json:"restart_required"`
}

func (r ModifyConfigurationResult) Extract() (*ModifyConfigurationResp, error) {
	var response ModifyConfigurationResp
	err := r.ExtractInto(&response)
	return &response, err
}

type ConfigParams struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Restart     bool   `json:"restart_required"`
	ReadOnly    bool   `json:"readonly"`
	ValueRange  string `json:"value_range"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type GetConfigurationResp struct {
	DatastoreVersion string         `json:"datastore_version_name"`
	DatastoreName    string         `json:"datastore_name"`
	Parameters       []ConfigParams `json:"configuration_parameters"`
}

func (r GetConfigurationResult) Extract() (*GetConfigurationResp, error) {
	var response GetConfigurationResp
	err := r.ExtractInto(&response)
	return &response, err
}

type GetTdeStatusResp struct {
	InstanceId string `json:"instance_id"`
	TdeStatus  string `json:"tde_status"`
}

func (r GetTdeStatusResult) Extract() (*GetTdeStatusResp, error) {
	var response GetTdeStatusResp
	err := r.ExtractInto(&response)
	return &response, err
}

type GetSecondLevelMonitoringResp struct {
	SwitchOption bool `json:"switch_option"`
	Interval     int  `json:"interval"`
}

func (r GetSecondLevelMonitoringResult) Extract() (*GetSecondLevelMonitoringResp, error) {
	var response GetSecondLevelMonitoringResp
	err := r.ExtractInto(&response)
	return &response, err
}

type GetBinlogRetentionHoursResp struct {
	BinlogRetentionHours int    `json:"binlog_retention_hours"`
	BinlogClearType      string `json:"binlog_clear_type"`
}

func (r GetBinlogRetentionHoursResult) Extract() (*GetBinlogRetentionHoursResp, error) {
	var response GetBinlogRetentionHoursResp
	err := r.ExtractInto(&response)
	return &response, err
}

type ListMsdtcHostsResponse struct {
	Hosts      []RdsMsdtcHosts `json:"hosts"`
	TotalCount int             `json:"total_count"`
}

type RdsMsdtcHosts struct {
	Id       string `json:"id"`
	Host     string `json:"host"`
	HostName string `json:"host_name"`
}

type MsdtcHostsPage struct {
	pagination.OffsetPageBase
}

func (r MsdtcHostsPage) IsEmpty() (bool, error) {
	data, err := ExtractRdsMsdtcHosts(r)
	if err != nil {
		return false, err
	}
	return len(data.Hosts) == 0, err
}

// ExtractRdsMsdtcHosts is a function that takes a ListResult and returns the msdct hosts' information.
func ExtractRdsMsdtcHosts(r pagination.Page) (ListMsdtcHostsResponse, error) {
	var s ListMsdtcHostsResponse
	err := (r.(MsdtcHostsPage)).ExtractInto(&s)
	return s, err
}

type ReplicationMode struct {
	WorkflowId      string `json:"workflowId"`
	InstanceId      string `json:"instanceId"`
	ReplicationMode string `json:"replicationMode"`
}

func (r ModifyReplicationModeResult) Extract() (*ReplicationMode, error) {
	var response ReplicationMode
	err := r.ExtractInto(&response)
	return &response, err
}

type BinlogRetentionHoursResp struct {
	Resp string `json:"resp"`
}

func (r ModifyBinlogRetentionHoursResult) Extract() (*BinlogRetentionHoursResp, error) {
	var response BinlogRetentionHoursResp
	err := r.ExtractInto(&response)
	return &response, err
}

type ListRdsResult struct {
	commonResult
}

type ListRdsResponse struct {
	Instances  []RdsInstanceResponse `json:"instances"`
	TotalCount int                   `json:"total_count"`
}

type RdsInstanceResponse struct {
	Id                  string             `json:"id"`
	Name                string             `json:"name"`
	Status              string             `json:"status"`
	PrivateIps          []string           `json:"private_ips"`
	PublicIps           []string           `json:"public_ips"`
	Port                int                `json:"port"`
	Type                string             `json:"type"`
	Ha                  Ha                 `json:"ha"`
	Region              string             `json:"region"`
	DataStore           Datastore          `json:"datastore"`
	Created             string             `json:"created"`
	Updated             string             `json:"updated"`
	DbUserName          string             `json:"db_user_name"`
	VpcId               string             `json:"vpc_id"`
	SubnetId            string             `json:"subnet_id"`
	SecurityGroupId     string             `json:"security_group_id"`
	FlavorRef           string             `json:"flavor_ref"`
	Volume              Volume             `json:"volume"`
	SwitchStrategy      string             `json:"switch_strategy"`
	BackupStrategy      BackupStrategy     `json:"backup_strategy"`
	ChargeInfo          ChargeResponse     `json:"charge_info"`
	MaintenanceWindow   string             `json:"maintenance_window"`
	Collation           string             `json:"collation"`
	Nodes               []Nodes            `json:"nodes"`
	RelatedInstance     []RelatedInstance  `json:"related_instance"`
	DiskEncryptionId    string             `json:"disk_encryption_id"`
	EnterpriseProjectId string             `json:"enterprise_project_id"`
	TimeZone            string             `json:"time_zone"`
	Alias               string             `json:"alias"`
	AssociatedWithDdm   bool               `json:"associated_with_ddm"`
	BackupUsedSpace     float64            `json:"backup_used_space"`
	Cpu                 string             `json:"cpu"`
	EnableSsl           bool               `json:"enable_ssl"`
	ExpirationTime      string             `json:"expiration_time"`
	MaxIops             int                `json:"max_iops"`
	Mem                 string             `json:"mem"`
	PrivateDnsNames     []string           `json:"private_dns_names"`
	ReadOnlyByUser      bool               `json:"read_only_by_user"`
	StorageUsedSpace    float64            `json:"storage_used_space"`
	Tags                []tags.ResourceTag `json:"tags"`
}

type Nodes struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Role             string `json:"role"`
	Status           string `json:"status"`
	AvailabilityZone string `json:"availability_zone"`
}

type RelatedInstance struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type RdsPage struct {
	pagination.SinglePageBase
}

type ErrorResponse struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (r RdsPage) IsEmpty() (bool, error) {
	data, err := ExtractRdsInstances(r)
	if err != nil {
		return false, err
	}
	return len(data.Instances) == 0, err
}

// ExtractRdsInstances is a function that takes a ListResult and returns the instances' information.
func ExtractRdsInstances(r pagination.Page) (ListRdsResponse, error) {
	var s ListRdsResponse
	err := (r.(RdsPage)).ExtractInto(&s)
	return s, err
}

type ErrorLogResult struct {
	golangsdk.Result
}

type ErrorLogResp struct {
	ErrorLogList []Errorlog `json:"error_log_list"`
	TotalRecord  int        `json:"total_record"`
}

type Errorlog struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Content string `json:"content"`
}

type ErrorLogPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no services.
func (r ErrorLogPage) IsEmpty() (bool, error) {
	data, err := ExtractErrorLog(r)
	if err != nil {
		return false, err
	}
	return len(data.ErrorLogList) == 0, err
}

func ExtractErrorLog(r pagination.Page) (ErrorLogResp, error) {
	var s ErrorLogResp
	err := (r.(ErrorLogPage)).ExtractInto(&s)
	return s, err
}

type SlowLogResp struct {
	Slowloglist []Slowloglist `json:"slow_log_list"`
	TotalRecord int           `json:"total_record"`
}

type Slowloglist struct {
	Count        string `json:"count"`
	Time         string `json:"time"`
	Locktime     string `json:"lock_time"`
	Rowssent     string `json:"rows_sent"`
	Rowsexamined string `json:"rows_examined"`
	Database     string `json:"database"`
	Users        string `json:"users"`
	QuerySample  string `json:"query_sample"`
	Type         string `json:"type"`
}

type SlowLogPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no services.
func (r SlowLogPage) IsEmpty() (bool, error) {
	data, err := ExtractSlowLog(r)
	if err != nil {
		return false, err
	}
	return len(data.Slowloglist) == 0, err
}

// ExtractCloudServers is a function that takes a ListResult and returns the services' information.
func ExtractSlowLog(r pagination.Page) (SlowLogResp, error) {
	var s SlowLogResp
	err := (r.(SlowLogPage)).ExtractInto(&s)
	return s, err
}

type RDSJobResult struct {
	commonResult
}

type ListJob struct {
	Job Job `json:"job"`
}

type Job struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Created    string `json:"created"`
	Ended      string `json:"ended"`
	Process    string `json:"process"`
	FailReason string `json:"fail_reason"`
}

func (r RDSJobResult) Extract() (ListJob, error) {
	var s ListJob
	err := r.ExtractInto(&s)
	return s, err
}

type Engine struct {
	Versions []VersionInfo `json:"dataStores"`
}

type VersionInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AutoExpansion is an object that represents the automatic expansion configuration.
type AutoExpansion struct {
	// Whether the automatic expansion is enabled.
	SwitchOption bool `json:"switch_option"`
	// The upper limit of automatic expansion of storage, in GB.
	// This parameter is mandatory when switch_option is set to true.
	// The value ranges from 40 GB to 4,000 GB and must be no less than the current storage of the instance.
	LimitSize int `json:"limit_size"`
	// The threshold to trigger automatic expansion.
	// If the available storage drops to this threshold or 10 GB, the automatic expansion is triggered.
	// This parameter is mandatory when switch_option is set to true.
	// The valid values are as follows:
	// + 10
	// + 15
	// + 20
	TriggerThreshold int `json:"trigger_threshold"`
}
