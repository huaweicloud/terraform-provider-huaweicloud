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
	VpcName             string               `json:"vpc_name"`
	ChargingMode        int                  `json:"charging_mode"`
	VpcId               string               `json:"vpc_id"`
	UserName            string               `json:"user_name"`
	CreatedAt           string               `json:"created_at"`
	Description         string               `json:"description"`
	SecurityGroupId     string               `json:"security_group_id"`
	SecurityGroupName   string               `json:"security_group_name"`
	MaxMemory           int                  `json:"max_memory"`
	UsedMemory          int                  `json:"used_memory"`
	Capacity            float64              `json:"capacity"`
	CapacityMinor       string               `json:"capacity_minor"`
	MaintainBegin       string               `json:"maintain_begin"`
	MaintainEnd         string               `json:"maintain_end"`
	Engine              string               `json:"engine"`
	NoPasswordAccess    string               `json:"no_password_access"`
	Ip                  string               `json:"ip"`
	BackupPolicy        InstanceBackupPolicy `json:"instance_backup_policy"`
	AzCodes             []string             `json:"az_codes"`
	AccessUser          string               `json:"access_user"`
	InstanceId          string               `json:"instance_id"`
	Port                int                  `json:"port"`
	UserId              string               `json:"user_id"`
	Name                string               `json:"name"`
	SpecCode            string               `json:"spec_code"`
	SubnetId            string               `json:"subnet_id"`
	SubnetName          string               `json:"subnet_name"`
	SubnetCidr          string               `json:"subnet_cidr"`
	EngineVersion       string               `json:"engine_version"`
	OrderId             string               `json:"order_id"`
	Status              string               `json:"status"`
	DomainName          string               `json:"domain_name"`
	EnablePublicIp      bool                 `json:"enable_publicip"`
	PublicIpId          string               `json:"publicip_id"`
	PublicIpAddress     string               `json:"publicip_address"`
	EnableSsl           bool                 `json:"enable_ssl"`
	ServiceUpgrade      bool                 `json:"service_upgrade"`
	ServiceTaskId       string               `json:"service_task_id"`
	EnterpriseProjectId string               `json:"enterprise_project_id"`
	BackendAddress      string               `json:"backend_addrs"`
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
