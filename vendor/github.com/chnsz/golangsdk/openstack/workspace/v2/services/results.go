package services

// RequestResp is the structure that represents the API response of service methods request.
type RequestResp struct {
	JobId string `json:"job_id"`
}

// CreateResp is the structure that represents the API response of Create method request.
type CreateResp struct {
	RequestResp
}

// UpdateResp is the structure that represents the API response of Update method request.
type UpdateResp struct {
	RequestResp
	// Enterprise ID.
	EnterpriseId string `json:"enterprise_id"`
}

// DeleteResp is the structure that represents the API response of Delete method request.
type DeleteResp struct {
	RequestResp
}

// Service is the structure that represents the Workspace service details.
type Service struct {
	// Workspace service ID.
	ID string `json:"id"`
	// Domain information.
	AdDomain DomainResp `json:"ad_domains"`
	// VPC ID.
	VpcId string `json:"vpc_id"`
	// VPC name.
	VpcName string `json:"vpc_name"`
	// Access mode.
	// + INTERNET: Indicates Internet access.
	// + DEDICATED: Indicates dedicated line access.
	// + BOTH: Indicates that both access methods are supported.
	AccessMode string `json:"access_mode"`
	// Dedicated access network segment.
	// This parameter is returned only when the access_mode is "DEDICATED" or "BOTH".
	DedicatedSubnets string `json:"dedicated_subnets"`
	// Dedicated access address
	// This parameter is returned only when the access_mode is "DEDICATED" or "BOTH".
	DedicatedAccessAddress string `json:"dedicated_access_address"`
	// Internet access address.
	// This parameter is returned only when the access_mode is "INTERNET" or "BOTH".
	InternetAccessAddress string `json:"internet_access_address"`
	// Internet access port.
	InternetAccessPort string `json:"internet_access_port"`
	// Status of cloud office services.
	// + PREPARING: ready to open.
	// + SUBSCRIBING: Subscription is in progress.
	// + SUBSCRIBED: subscribed.
	// + SUBSCRIPTION_FAILED: Subscription failed.
	// + DEREGISTERING: The account is being sold out.
	// + DEREGISTRATION_FAILED: Account cancellation failed.
	// + CLOSED: The account that has been closed has not been opened.
	Status string `json:"status"`
	// The status of the Internet and private line switching tasks.
	// + init: Initialization - the initial state after the service is activated.
	// + available: Available - The normal state where the task was executed and returned after success.
	// + internetOpening: On - Internet Access is on.
	// + dedicatedOpening: Open - Dedicated line access is open.
	// + internetOpenFailed: Failed to open - Open internet access failed to open.
	// + dedicatedOpenFailed: Failed to open - Dedicated line access failed to open.
	// + openSuccess: Open successfully - Internet access is successfully opened.
	// + internetClosing: Closing - Closing Internet access is closing.
	// + dedicatedClosing: Closing - Dedicated line access is closed.
	// + internetCloseFailed: Failed to close - Failed to close the internet access method.
	// + dedicatedCloseFailed: Failed to close - Failed to close dedicated line access.
	// + closeSuccess: Close Success - Close the access method successfully.
	// + internetAccessPortModifying: The internet access port is being modified.
	// + internetAccessPortModifyFailed: Port modification failed.
	AccessStatus string `json:"access_status"`
	// Service subnet which to specify the returned network ID to order desktops.
	SubnetIds []Subnet `json:"subnet_ids"`
	// The subnet segment for the management component.
	ManagementSubentCidr string `json:"management_subnet_cidr"`
	// The management component security group automatically created under the specified VPC after the service is enabled.
	InfrastructureSecurityGroup SecurityGroup `json:"infrastructure_security_group"`
	// The desktop security group automatically created under the specified VPC after the service is enabled.
	DesktopSecurityGroup SecurityGroup `json:"desktop_security_group"`
	// Whether the service can be canceled.
	Closable bool `json:"closable"`
	// Configuration status.
	// + 0: The service is successfully activated and the connection to AD is successful.
	// + 1: The service is successfully activated, but the AD configuration fails.
	// + 2: The service is successfully activated, but there are other errors after the AD configuration fails.
	// + 3: The service is successfully activated, but the AD connection is not enabled.
	ConfigStatus string `json:"config_status"`
	// The progress of service activation or deregistration, in percentage format, for example: 100%.
	Progress string `json:"progress"`
	// The job ID of service activation or deregistration.
	JobId string `json:"job_id"`
	// Failure error code.
	FailCode int `json:"fail_code"`
	// Failure reason.
	FailReason string `json:"fail_reason"`
	// Enterprise ID.
	EnterpriseId string `json:"enterprise_id"`
}

// DomainResp is an object to specified the configuration details of AD domain.
type DomainResp struct {
	// Domain type.
	// + LITE_AS: Local authentication.
	// + LOCAL_AD: Local AD.
	// When the domain type is "LOCAL_AD", make sure that the selected VPC network and the network to which AD
	//   belongs can be connected.
	Type string `json:"domain_type"`
	// Domain name. It needs to be configured when the domain type is LOCAL_AD.
	// The domain name must be an existing domain name on the AD server, and the length should not exceed 55.
	Name string `json:"domain_name"`
	// Domain administrator account. It needs to be configured when the domain type is "LOCAL_AD".
	// It must be an existing domain administrator account on the AD server.
	AdminAccount string `json:"domain_admin_account"`
	// Domain administrator account password. It needs to be configured when the domain type is "LOCAL_AD".
	Password string `json:"domain_password"`
	// Primary domain controller IP address. It needs to be configured when the domain type is LOCAL_AD.
	ActiveDomainIp string `json:"active_domain_ip"`
	// Primary domain controller name. It needs to be configured when the domain type is LOCAL_AD.
	AcitveDomainName string `json:"active_domain_name"`
	// The IP address of the standby domain controller.
	// It needs to be configured when the domain type is LOCAL_AD and the standby node is configured.
	StandyDomainIp string `json:"standby_domain_ip"`
	// The name of the standby domain controller.
	// It needs to be configured when the domain type is LOCAL_AD and the standby node is configured.
	StandyDomainName string `json:"standby_domain_name"`
	// Primary DNS IP address.
	// It needs to be configured when the domain type is LOCAL_AD.
	ActiveDnsIp string `json:"active_dns_ip"`
	// Backup DNS IP address.
	// It needs to be configured when the domain type is LOCAL_AD and the standby node is configured.
	StandyDnsIp string `json:"standby_dns_ip"`
	// Whether to delete the corresponding computer object on AD while deleting the desktop.
	// + 0 means not delete
	// + 1 means delete.
	DeleteComputerObject string `json:"delete_computer_object"`
	// Whether to enable LDAPS.
	UseLdaps bool `json:"use_idaps"`
	// The configuration of TLS.
	TlsConfig TlsConfig `json:"tls_config"`
}

// SecurityGroup is an object to specified the security group that service have.
type SecurityGroup struct {
	// Security Group ID.
	ID string `json:"id"`
	// Security Group name.
	Name string `json:"name"`
}

// OtpConfigResp represents that the GetAuthConfig method response details.
type OtpConfigResp struct {
	OptConfigInfo OtpConfig `json:"otp_config_info"`
}

// OtpConfig represents that the Auxiliary authentication configuration details.
type OtpConfig struct {
	// Whether to enable OTP authentication mode.
	Enable bool `json:"enable"`
	// Verification code receiving mode.
	// + VMFA Indicates virtual MFA device.
	// + HMFA Indicates hardware MFA device.
	ReceiveMode string `json:"receive_mode"`
	// Auxiliary authentication server address.
	AuthUrl string `json:"auth_url"`
	// Auxiliary authentication service access account.
	AppId string `json:"app_id"`
	// Auxiliary authentication service access password.
	AppSecrte string `json:"app_secret"`
	// Auxiliary authentication service access mode.
	// + INTERNET: Indicates Internet access.
	// + DEDICATED: Indicates dedicated line access.
	// + SYSTEM_DEFAULT：Indicates system default.
	AuthServerAccessMode string `json:"auth_server_access_mode"`
	// PEM format certificate content.
	CertContent string `json:"cert_content"`
	// Authentication application object information. If null, it means it is effective for all application objects.
	ApplyRule ApplyRuleInfo `json:"apply_rule"`
}

// ApplyRuleInfo is an object specified the detail of Authentication application object.
type ApplyRuleInfo struct {
	// Authentication application object type.
	// + ACCESS_MODE: Indicates access type.
	RuleType string `json:"rule_type"`
	// Authentication application object.
	// + INTERNET: Indicates Internet access. Optional only when rule_type is "ACCESS_MODE".
	// + PRIVATE: Indicates dedicated line access. Optional only when rule_type is "ACCESS_MODE".
	Rule string `json:"rule"`
}

// LockStatusResp is the structure that represents the API response of GetLockStatus method request.
type LockStatusResp struct {
	// Whether the Workspace service is locked.
	// + 0: Indicates not locked.
	// + 1: Indicates locked.
	IsLocked int `json:"is_locked"`
	// The lock time of the Workspace service.
	LockTime string `json:"lock_time"`
	// The reason of the Workspace service is locked.
	LockReason string `json:"lock_reason"`
}

type UnlockResp struct {
	RequestResp
}
