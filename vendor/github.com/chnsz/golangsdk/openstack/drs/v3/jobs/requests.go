package jobs

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

const (
	DeleteTypeTerminate      = "terminate"
	DeleteTypeForceTerminate = "force_terminate"
	DeleteTypeDelete         = "delete"
)

type BatchCreateJobReq struct {
	Jobs []CreateJobReq `json:"jobs" required:"true"`
}

type CreateJobReq struct {
	BindEip          *bool              `json:"bind_eip,omitempty"`
	DbUseType        string             `json:"db_use_type" required:"true"`
	Name             string             `json:"name" required:"true"`
	Description      string             `json:"description,omitempty"`
	EngineType       string             `json:"engine_type" required:"true"`
	IsTargetReadonly *bool              `json:"is_target_readonly,omitempty"`
	JobDirection     string             `json:"job_direction" required:"true"`
	NetType          string             `json:"net_type" required:"true"`
	NodeNum          *int               `json:"node_num,omitempty"`
	NodeType         string             `json:"node_type" required:"true"`
	SourceEndpoint   Endpoint           `json:"source_endpoint" required:"true"`
	TargetEndpoint   Endpoint           `json:"target_endpoint" required:"true"`
	TaskType         string             `json:"task_type" required:"true"`
	SubnetId         string             `json:"customize_sutnet_id" required:"true"`
	ProductId        string             `json:"product_id,omitempty"`
	ExpiredDays      string             `json:"expired_days,omitempty"`
	MultiWrite       *bool              `json:"multi_write,omitempty"`
	Tags             []tags.ResourceTag `json:"tags,omitempty"`
	SysTags          []tags.ResourceTag `json:"sys_tags,omitempty"`
	MasterAz         string             `json:"master_az,omitempty"`
	SlaveAz          string             `json:"slave_az,omitempty"`
	ChargingMode     string             `json:"charging_mode,omitempty"`
	PeriodOrder      *PeriodOrder       `json:"period_order,omitempty"`
	PublciIpList     []PublciIpList     `json:"public_ip_list,omitempty"`
	IsOpenFastClean  bool               `json:"is_open_fast_clean,omitempty"`
}

type Endpoint struct {
	DbType          string `json:"db_type" required:"true"`
	AzCode          string `json:"az_code,omitempty"`
	Region          string `json:"region,omitempty"`
	InstanceId      string `json:"inst_id,omitempty"`
	InstanceName    string `json:"inst_name,omitempty"`
	VpcId           string `json:"vpc_id,omitempty"`
	SubnetId        string `json:"subnet_id,omitempty"`
	SecurityGroupId string `json:"security_group_id,omitempty"`
	ProjectId       string `json:"project_id,omitempty"`

	DbName     string `json:"db_name,omitempty"`
	Ip         string `json:"ip,omitempty"`
	DbPort     *int   `json:"db_port,omitempty"`
	DbUser     string `json:"db_user,omitempty"`
	DbPassword string `json:"db_password,omitempty"`

	SslCertPassword string `json:"ssl_cert_password,omitempty"`
	SslCertCheckSum string `json:"ssl_cert_check_sum,omitempty"`
	SslCertKey      string `json:"ssl_cert_key,omitempty"`
	SslCertName     string `json:"ssl_cert_name,omitempty"`
	SslLink         *bool  `json:"ssl_link,omitempty"`

	SafeMode    *int   `json:"safe_mode,omitempty"`
	MongoHaMode string `json:"mongo_ha_mode,omitempty"`
	Topic       string `json:"topic,omitempty"`
	ClusterMode string `json:"cluster_mode,omitempty"`

	KafkaSecurityConfig *KafkaSecurityConfig `json:"kafka_security_config,omitempty"`
}

type PeriodOrder struct {
	PeriodType  int `json:"period_type,omitempty"`
	PeriodNum   int `json:"period_num,omitempty"`
	IsAutoRenew int `json:"is_auto_renew,omitempty"`
}

type PublciIpList struct {
	Id       string `json:"id" required:"true"`
	PublicIp string `json:"public_ip" required:"true"`
	Type     string `json:"type" required:"true"`
}

type KafkaSecurityConfig struct {
	Type                  string `json:"type,omitempty"`
	SaslMechanism         string `json:"sasl_mechanism,omitempty"`
	TrustStoreKeyName     string `json:"trust_store_key_name,omitempty"`
	TrustStoreKey         string `json:"trust_store_key,omitempty"`
	TrustStorePassword    string `json:"trust_store_password,omitempty"`
	EndpointAlgorithm     string `json:"endpoint_algorithm,omitempty"`
	DelegationTokens      bool   `json:"delegation_tokens,omitempty"`
	EnableKeyStore        bool   `json:"enable_key_store,omitempty"`
	KeyStoreKeyName       string `json:"key_store_key_name,omitempty"`
	KeyStoreKey           string `json:"key_store_key,omitempty"`
	KeyStorePassword      string `json:"key_store_password,omitempty"`
	SetPrivateKeyPassword bool   `json:"set_private_key_password,omitempty"`
	KeyPassword           string `json:"key_password,omitempty"`
}

type QueryJobReq struct {
	Jobs    []string `json:"jobs" required:"true"`
	PageReq PageReq  `json:"page_req,omitempty"`
}

type PageReq struct {
	CurPage int `json:"cur_page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
}

type StartJobReq struct {
	Jobs []StartInfo `json:"jobs" required:"true"`
}

type StartInfo struct {
	JobId     string `json:"job_id" required:"true"`
	StartTime string `json:"start_time,omitempty"`
}

type TestConnectionsReq struct {
	Jobs []TestEndPoint `json:"jobs" required:"true"`
}

type TestEndPoint struct {
	JobId   string `json:"id" required:"true"`
	NetType string `json:"net_type" required:"true"`
	// source DB:so, target DB:ta
	EndPointType string `json:"end_point_type" required:"true"`

	DbType     string `json:"db_type" required:"true"`
	Ip         string `json:"ip" required:"true"`
	DbUser     string `json:"db_user,omitempty"`
	DbPassword string `json:"db_password,omitempty"`
	//when type is Mongo、DDS, must be `0`
	DbPort *int   `json:"db_port,omitempty"`
	DbName string `json:"db_name,omitempty"`

	Region    string `json:"region,omitempty"`
	ProjectId string `json:"project_id,omitempty"`
	InstId    string `json:"inst_id,omitempty"`
	VpcId     string `json:"vpc_id,omitempty"`
	SubnetId  string `json:"subnet_id,omitempty"`

	SslLink         *bool  `json:"ssl_link,omitempty"`
	SslCertKey      string `json:"ssl_cert_key,omitempty"`
	SslCertName     string `json:"ssl_cert_name,omitempty"`
	SslCertCheckSum string `json:"ssl_cert_check_sum,omitempty"`
	SslCertPassword string `json:"ssl_cert_password,omitempty"`

	KafkaSecurityConfig *KafkaSecurityConfig `json:"kafka_security_config,omitempty"`
}

type TestClusterConnectionsReq struct {
	Jobs []TestJob `json:"jobs" required:"true"`
}

type TestJob struct {
	Action   string `json:"action" required:"true"`
	JobId    string `json:"job_id" required:"true"`
	Property string `json:"property" required:"true"`
}

type PropertyParam struct {
	NetType             string               `json:"nettype" required:"true"`
	EndPointType        string               `json:"endpointtype" required:"true"`
	DbType              string               `json:"dbtype" required:"true"`
	Ip                  string               `json:"ip" required:"true"`
	DbUser              string               `json:"dbuser,omitempty"`
	DbPassword          string               `json:"dbpassword,omitempty"`
	DbPort              *int                 `json:"dbport,omitempty"`
	DbName              string               `json:"dbName,omitempty"`
	Region              string               `json:"region,omitempty"`
	ProjectId           string               `json:"projectId,omitempty"`
	InstId              string               `json:"instid,omitempty"`
	VpcId               string               `json:"vpcId,omitempty"`
	SubnetId            string               `json:"subnetId,omitempty"`
	SslLink             *bool                `json:"ssllink,omitempty"`
	SslCertKey          string               `json:"sslcertkey,omitempty"`
	SslCertName         string               `json:"sslcertname,omitempty"`
	SslCertCheckSum     string               `json:"sslcertchecksum,omitempty"`
	KafkaSecurityConfig *KafkaSecurityConfig `json:"kafka_security_config,omitempty"`
}

type BatchDeleteJobReq struct {
	Jobs []DeleteJobReq `json:"jobs" required:"true"`
}

type DeleteJobReq struct {
	DeleteType string `json:"delete_type" required:"true"`
	JobId      string `json:"job_id" required:"true"`
}

type UpdateReq struct {
	Jobs []UpdateJobReq `json:"jobs" required:"true"`
}

type UpdateJobReq struct {
	JobId          string    `json:"job_id" required:"true"`
	Name           string    `json:"name,omitempty"`
	SourceEndpoint *Endpoint `json:"source_endpoint,omitempty"`
	TargetEndpoint *Endpoint `json:"target_endpoint,omitempty"`
	NodeType       string    `json:"node_type,omitempty"`
	EngineType     string    `json:"engine_type,omitempty"`
	NetType        string    `json:"net_type,omitempty"`
	StoreDbInfo    bool      `json:"store_db_info,omitempty"`

	IsRecreate       *bool  `json:"is_recreate,omitempty"`
	Description      string `json:"description,omitempty"`
	TaskType         string `json:"task_type,omitempty"`
	DbUseType        string `json:"db_use_type,omitempty"`
	JobDirection     string `json:"job_direction,omitempty"`
	IsTargetReadonly *bool  `json:"is_target_readonly,omitempty"`
	ProductId        string `json:"product_id,omitempty"`

	ReplaceDefiner *bool              `json:"replace_definer,omitempty"`
	Tags           []tags.ResourceTag `json:"tags,omitempty"`
	AlarmNotify    *AlarmNotifyInfo   `json:"alarm_notify,omitempty"`
}

type AlarmNotifyInfo struct {
	DelayTime     *int               `json:"delay_time,omitempty"`
	RtoDelay      *int               `json:"rto_delay,omitempty"`
	RpoDelay      *int               `json:"rpo_delay,omitempty"`
	AlarmToUser   bool               `json:"alarm_to_user" required:"true"`
	Subscriptions []SubscriptionInfo `json:"subscriptions,omitempty"`
}

type SubscriptionInfo struct {
	Endpoints []string `json:"endpoints" required:"true"`
	// sms; email
	Protocol string `json:"protocol" required:"true"`
}

type BatchLimitSpeedReq struct {
	SpeedLimits []LimitSpeedReq `json:"speed_limits" required:"true"`
}

type LimitSpeedReq struct {
	JobId      string           `json:"job_id" required:"true"`
	SpeedLimit []SpeedLimitInfo `json:"speed_limit" required:"true"`
}

type SpeedLimitInfo struct {
	// format: hh:mm
	Begin string `json:"begin" required:"true"`
	// format: hh:mm
	End string `json:"end" required:"true"`
	// range: 1~9999 , unit: MB/s
	Speed string `json:"speed" required:"true"`
}

type ListJobsReq struct {
	CurPage             int    `json:"cur_page" required:"true"`
	PerPage             int    `json:"per_page" required:"true"`
	DbUseType           string `json:"db_use_type" required:"true"`
	EngineType          string `json:"engine_type,omitempty"`
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// name  or id
	Name        string `json:"name,omitempty"`
	NetType     string `json:"net_type,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	Status      string `json:"status,omitempty"`

	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

type BatchPrecheckReq struct {
	Jobs []PreCheckInfo `json:"jobs" required:"true"`
}

type PreCheckInfo struct {
	JobId        string `json:"job_id" required:"true"`
	PrecheckMode string `json:"precheck_mode" required:"true"`
}

type QueryPrecheckResultReq struct {
	Jobs []string `json:"jobs" required:"true"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func Create(c *golangsdk.ServiceClient, opts BatchCreateJobReq) (*CreateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r CreateResp
	_, err = c.Post(createURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r, err
}

func Get(c *golangsdk.ServiceClient, opts QueryJobReq) (*JobDetailResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst JobDetailResp
	_, err = c.Post(detailsURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Start(c *golangsdk.ServiceClient, opts StartJobReq) (*ActionResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r ActionResp
	_, err = c.Post(startURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r, err
}

func TestConnections(c *golangsdk.ServiceClient, opts TestConnectionsReq) (*ActionResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst ActionResp
	_, err = c.Post(testConnectionsURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func TestClusterConnections(c *golangsdk.ServiceClient, opts TestClusterConnectionsReq) (*ActionResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst ActionResp
	_, err = c.Post(testClusterConnectionsURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Update(c *golangsdk.ServiceClient, opts UpdateReq) (*golangsdk.Result, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r golangsdk.Result
	_, err = c.Put(updateJobURL(c), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r, err
}

func Delete(c *golangsdk.ServiceClient, opts BatchDeleteJobReq) *golangsdk.ErrResult {
	var r golangsdk.ErrResult

	url := deleteURL(c)
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return &r
	}

	_, r.Err = c.DeleteWithBody(url, b, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r
}

func LimitSpeed(c *golangsdk.ServiceClient, opts BatchLimitSpeedReq) (*ActionResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst ActionResp
	_, err = c.Put(limitSpeedURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func PreCheckJobs(c *golangsdk.ServiceClient, opts BatchPrecheckReq) (*PrecheckResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst PrecheckResp
	_, err = c.Post(preCheckURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func CheckResults(c *golangsdk.ServiceClient, opts QueryPrecheckResultReq) (*PrecheckResultResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst PrecheckResultResp
	_, err = c.Post(batchCheckResultsURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func List(c *golangsdk.ServiceClient, opts ListJobsReq) (*JobsListResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst JobsListResp
	_, err = c.Post(listURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Status(c *golangsdk.ServiceClient, opts QueryJobReq) (*StatusResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst StatusResp
	_, err = c.Post(statusURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Progress(c *golangsdk.ServiceClient, opts QueryJobReq) (*ProgressResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst ProgressResp
	_, err = c.Post(progressURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}
