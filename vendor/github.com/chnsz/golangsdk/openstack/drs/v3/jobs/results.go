package jobs

import "github.com/chnsz/golangsdk/openstack/common/tags"

type CreateResp struct {
	Results []CreateJobResp `json:"results"`
	Count   int             `json:"count"`
}

type CreateJobResp struct {
	Id         string   `json:"id"`
	ChildIds   []string `json:"child_ids"`
	Name       string   `json:"name"`
	Status     string   `json:"status"`
	CreateTime string   `json:"create_time"`
	ErrorCode  string   `json:"error_code"`
	ErrorMsg   string   `json:"error_msg"`
}

type StatusResp struct {
	Count   int         `json:"count"`
	Results []JobStatus `json:"results"`
}

type JobStatus struct {
	Id           string `json:"id"`
	Status       string `json:"status"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type JobDetailResp struct {
	Count   int         `json:"count"`
	Results []JobDetail `json:"results"`
}

type JobDetail struct {
	Id                       string             `json:"id"`
	ParentId                 string             `json:"parent_id"`
	Name                     string             `json:"name"`
	Status                   string             `json:"status"`
	Description              string             `json:"description"`
	CreateTime               string             `json:"create_time"`
	TaskType                 string             `json:"task_type"`
	SourceEndpoint           Endpoint           `json:"source_endpoint"`
	DmqEndpoint              Endpoint           `json:"dmq_endpoint"`
	SourceSharding           []Endpoint         `json:"source_sharding"`
	TargetEndpoint           Endpoint           `json:"target_endpoint"`
	NetType                  string             `json:"net_type"`
	FailedReason             string             `json:"failed_reason"`
	InstInfo                 InstInfo           `json:"inst_info"`
	ActualStartTime          string             `json:"actual_start_time"`
	FullTransferCompleteTime string             `json:"full_transfer_complete_time"`
	UpdateTime               string             `json:"update_time"`
	JobDirection             string             `json:"job_direction"`
	OriginalJobDirection     string             `json:"original_job_direction"`
	DbUseType                string             `json:"db_use_type"`
	NeedRestart              bool               `json:"need_restart"`
	IsTargetReadonly         bool               `json:"is_target_readonly"`
	ConflictPolicy           string             `json:"conflict_policy"`
	FilterDdlPolicy          string             `json:"filter_ddl_policy"`
	SpeedLimit               []SpeedLimitInfo   `json:"speed_limit"`
	SchemaType               string             `json:"schema_type"`
	NodeNum                  int                `json:"node_num"`
	ObjectSwitch             bool               `json:"object_switch"`
	ObjectInfos              []ObjectInfo       `json:"object_infos"`
	MasterAz                 string             `json:"master_az"`
	MasterJobId              string             `json:"master_job_id"`
	SlaveAz                  string             `json:"slave_az"`
	FullMode                 string             `json:"full_mode"`
	StructTrans              bool               `json:"struct_trans"`
	IndexTrans               bool               `json:"index_trans"`
	ReplaceDefiner           bool               `json:"replace_definer"`
	MigrateUser              bool               `json:"migrate_user"`
	SyncDatabase             bool               `json:"sync_database"`
	ErrorCode                string             `json:"error_code"`
	ErrorMessage             string             `json:"error_message"`
	TargetRootDb             DefaultRootDb      `json:"target_root_db"`
	AzCode                   string             `json:"az_code"`
	VpcId                    string             `json:"vpc_id"`
	SubnetId                 string             `json:"subnet_id"`
	SecurityGroupId          string             `json:"security_group_id"`
	MultiWrite               bool               `json:"multi_write"`
	SupportIpV6              bool               `json:"support_ip_v6"`
	InheritId                string             `json:"inherit_id"`
	Gtid                     string             `json:"gtid"`
	AlarmNotify              AlarmNotifyInfo    `json:"alarm_notify"`
	IncreStartPosition       string             `json:"incre_start_position"`
	Tags                     []tags.ResourceTag `json:"tags"`
	PeriodOrder              OrderInfo          `json:"period_order"`
	IsOpenFastClean          bool               `json:"is_open_fast_clean"`
}

type InstInfo struct {
	EngineType string `json:"engine_type"`
	InstType   string `json:"inst_type"`
	Ip         string `json:"ip"`
	PublicIp   string `json:"public_ip"`
	StartTime  string `json:"start_time"`
	Status     string `json:"status"`
	VolumeSize int    `json:"volume_size"`
}

type ObjectInfo struct {
	Id        string `json:"id"`
	ParentId  string `json:"parent_id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	AliasName string `json:"alias_name"`
}

type DefaultRootDb struct {
	DbName     string `json:"db_name"`
	DbEncoding string `json:"db_encoding"`
}

type OrderInfo struct {
	Status       string `json:"status"`
	OrderId      string `json:"order_id"`
	ChargingMode int    `json:"charging_mode"`
	PeriodType   int    `json:"period_type"`
	PeriodNum    int    `json:"period_num"`
	IsAutoRenew  int    `json:"is_auto_renew"`
	EffTime      string `json:"eff_time"`
	ExpTime      string `json:"exp_time"`
}

type ActionResp struct {
	Results []ActionResult `json:"results"`
	Count   int            `json:"count"`
}

type ActionResult struct {
	Id string `json:"id"`
	// success
	// failed
	Status    string `json:"status"`
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	// only valid in TestConnection
	Success bool `json:"success"`
}

type PrecheckResp struct {
	Results []PreCheckDetail `json:"results"`
	Count   int              `json:"count"`
}

type PreCheckDetail struct {
	Id         string `json:"id"`
	PrecheckId string `json:"precheck_id"`
	Status     string `json:"status"`
	ErrorCode  string `json:"error_code"`
	ErrorMsg   string `json:"error_msg"`
}

type PrecheckResultResp struct {
	Count   int              `json:"count"`
	Results []PrecheckResult `json:"results"`
}

type PrecheckResult struct {
	PrecheckId      string      `json:"precheck_id"`
	Result          bool        `json:"result"`
	Process         string      `json:"process"`
	TotalPassedRate string      `json:"total_passed_rate"`
	RdsInstanceId   string      `json:"rds_instance_id"`
	JobDirection    string      `json:"job_direction"`
	PrecheckResult  []CheckItem `json:"precheck_result"`
	ErrorMsg        string      `json:"error_msg"`
	ErrorCode       string      `json:"error_code"`
}

type CheckItem struct {
	Item          string                 `json:"item"`
	Result        string                 `json:"result"`
	FailedReason  string                 `json:"failed_reason"`
	Data          string                 `json:"data"`
	RawErrorMsg   string                 `json:"raw_error_msg"`
	Group         string                 `json:"group"`
	FailedSubJobs []PrecheckFailSubJobVO `json:"failed_sub_jobs"`
}

type PrecheckFailSubJobVO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CheckResult string `json:"check_result"`
}

type JobsListResp struct {
	TotalRecord int       `json:"total_record"`
	Jobs        []JobInfo `json:"jobs"`
}

type JobInfo struct {
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	Status           string            `json:"status"`
	Description      string            `json:"description"`
	CreateTime       string            `json:"create_time"`
	EngineType       string            `json:"engine_type"`
	NetType          string            `json:"net_type"`
	BillingTag       bool              `json:"billing_tag"`
	JobDirection     string            `json:"job_direction"`
	DbUseType        string            `json:"db_use_type"`
	TaskType         string            `json:"task_type"`
	Children         []ChildrenJobInfo `json:"children"`
	NodeNewFramework bool              `json:"node_newFramework"`
	JobAction        JobAction         `json:"job_action"`
}

type ChildrenJobInfo struct {
	BillingTag       bool   `json:"billing_tag"`
	CreateTime       string `json:"create_time"`
	DbUseType        string `json:"db_use_type"`
	Description      string `json:"description"`
	EngineType       string `json:"engine_type"`
	ErrorMsg         string `json:"error_msg"`
	Id               string `json:"id"`
	JobDirection     string `json:"job_direction"`
	Name             string `json:"name"`
	NetType          string `json:"net_type"`
	NodeNewFramework bool   `json:"node_newFramework"`
	Status           string `json:"status"`
	TaskType         string `json:"task_type"`
}

type JobAction struct {
	AvailableActions   []string `json:"available_actions"`
	UnavailableActions []string `json:"unavailable_actions"`
	CurrentAction      string   `json:"current_action"`
}

type ProgressResp struct {
	Count   int      `json:"count"`
	Results []Result `json:"results"`
}

type Result struct {
	JobId                 string                  `json:"job_id"`
	Progress              string                  `json:"progress"`
	IncreTransDelay       string                  `json:"incre_trans_delay"`
	IncreTransDelayMillis string                  `json:"incre_trans_delay_millis"`
	TaskMode              string                  `json:"task_mode"`
	TransferStatus        string                  `json:"transfer_status"`
	ProcessTime           string                  `json:"process_time"`
	RemainingTime         string                  `json:"remaining_time"`
	ProgressMap           map[string]ProgressInfo `json:"progress_map"`
	ErrorCode             string                  `json:"error_code"`
	ErrorMsg              string                  `json:"error_msg"`
}

type ProgressInfo struct {
	Completed     string `json:"completed"`
	RemainingTime string `json:"remaining_time"`
}
