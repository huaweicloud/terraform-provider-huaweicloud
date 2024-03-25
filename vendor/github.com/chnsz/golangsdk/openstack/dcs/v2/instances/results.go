package instances

type CreateResponse struct {
	OrderId   string           `json:"order_id"`
	Instances []SimpleInstance `json:"instances"`
}

type SimpleInstance struct {
	InstanceId   string `json:"instance_id,omitempty"`
	InstanceName string `json:"instance_name,omitempty"`
}

type DcsInstance struct {
	VpcName                   string               `json:"vpc_name"`
	ChargingMode              int                  `json:"charging_mode"`
	VpcId                     string               `json:"vpc_id"`
	UserName                  string               `json:"user_name"`
	CreatedAt                 string               `json:"created_at"`
	LaunchedAt                string               `json:"launched_at"`
	Description               string               `json:"description"`
	SecurityGroupId           string               `json:"security_group_id"`
	SecurityGroupName         string               `json:"security_group_name"`
	MaxMemory                 int                  `json:"max_memory"`
	UsedMemory                int                  `json:"used_memory"`
	Capacity                  float64              `json:"capacity"`
	CapacityMinor             string               `json:"capacity_minor"`
	MaintainBegin             string               `json:"maintain_begin"`
	MaintainEnd               string               `json:"maintain_end"`
	Engine                    string               `json:"engine"`
	NoPasswordAccess          string               `json:"no_password_access"`
	Ip                        string               `json:"ip"`
	BackupPolicy              InstanceBackupPolicy `json:"instance_backup_policy"`
	AzCodes                   []string             `json:"az_codes"`
	AccessUser                string               `json:"access_user"`
	InstanceId                string               `json:"instance_id"`
	Port                      int                  `json:"port"`
	UserId                    string               `json:"user_id"`
	Name                      string               `json:"name"`
	SpecCode                  string               `json:"spec_code"`
	SubnetId                  string               `json:"subnet_id"`
	SubnetName                string               `json:"subnet_name"`
	SubnetCidr                string               `json:"subnet_cidr"`
	EngineVersion             string               `json:"engine_version"`
	OrderId                   string               `json:"order_id"`
	Status                    string               `json:"status"`
	DomainName                string               `json:"domain_name"`
	EnablePublicIp            bool                 `json:"enable_publicip"`
	PublicIpId                string               `json:"publicip_id"`
	PublicIpAddress           string               `json:"publicip_address"`
	EnableSsl                 bool                 `json:"enable_ssl"`
	ServiceUpgrade            bool                 `json:"service_upgrade"`
	ServiceTaskId             string               `json:"service_task_id"`
	EnterpriseProjectId       string               `json:"enterprise_project_id"`
	BackendAddress            string               `json:"backend_addrs"`
	BandWidthDetail           BandWidthInfo        `json:"bandwidth_info"`
	CacheMode                 string               `json:"cache_mode"`
	CpuType                   string               `json:"cpu_type"`
	ReplicaCount              int                  `json:"replica_count"`
	ReadOnlyDomainName        string               `json:"readonly_domain_name"`
	TransparentClientIpEnable bool                 `json:"transparent_client_ip_enable"`
	ShardingCount             int                  `json:"sharding_count"`
	ProductType               string               `json:"product_type"`
}

type InstanceBackupPolicy struct {
	BackupPolicyId string          `json:"backup_policy_id"`
	Policy         DcsBackupPolicy `json:"policy"`
}

type DcsBackupPolicy struct {
	BackupType           string     `json:"backup_type"`
	SaveDays             int        `json:"save_days"`
	PeriodicalBackupPlan BackupPlan `json:"periodical_backup_plan"`
}

type ResizeResponse struct {
	OrderId string `json:"order_id"`
}

type ResizePrePaidResponse struct {
	OrderId string `json:"order_id"`
}

type RestartResponse struct {
	Results []RestartResultResponse `json:"results"`
}

type RestartResultResponse struct {
	Result   string `json:"result"`
	Instance string `json:"instance"`
}

type Configuration struct {
	ConfigTime   string        `json:"config_time"`
	InstanceId   int           `json:"instance_id"`
	RedisConfig  []RedisConfig `json:"redis_config"`
	ConfigStatus string        `json:"config_status"`
	Status       string        `json:"status"`
}

type RedisConfig struct {
	Proxy              bool   `json:"proxy"`
	ParamId            string `json:"param_id"`
	ParamName          string `json:"param_name"`
	DefaultValue       string `json:"default_value"`
	ValueRange         string `json:"value_range"`
	ValueType          string `json:"value_type"`
	Description        string `json:"description"`
	ParamValue         string `json:"param_value"`
	NodeRole           string `json:"node_role"`
	ParamType          string `json:"param_type"`
	UserPermission     string `json:"user_permission"`
	NeedRestart        bool   `json:"need_restart"`
	SupportedVersion   string `json:"supported_version"`
	InternationalKey   string `json:"international_key"`
	NewVersionOnly     bool   `json:"new_version_only"`
	SupportDataVersion string `json:"support_data_version"`
	Customized         bool   `json:"customized"`
}

type GetSslResponse struct {
	Enable       bool   `json:"enabled"`
	Ip           string `json:"ip"`
	SslValidated bool   `json:"ssl_validated"`
	Port         string `json:"port"`
	DomainName   string `json:"domain_name"`
	SslExpiredAt string `json:"ssl_expired_at"`
}

type BandWidthInfo struct {
	BandWidth          int  `json:"bandwidth"`
	BeginTime          int  `json:"begin_time"`
	CurrentTime        int  `json:"current_time"`
	EndTime            int  `json:"end_time"`
	ExpandCount        int  `json:"expand_count"`
	ExpandEffectTime   int  `json:"expand_effect_time"`
	ExpandIntervalTime int  `json:"expand_interval_time"`
	MaxExpandCount     int  `json:"max_expand_count"`
	NextExpandTime     int  `json:"next_expand_time"`
	TaskRunning        bool `json:"task_running"`
}
