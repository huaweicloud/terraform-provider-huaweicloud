package services

import "github.com/chnsz/golangsdk"

// CreateOpts is the structure required by the Create method to subscribe Workspace service.
type CreateOpts struct {
	// Configuration of domain.
	AdDomain *Domain `json:"ad_domains" required:"true"`
	// VPC ID.
	VpcId string `json:"vpc_id" required:"true"`
	// The network IDs of the service subnet. The subnet cannot conflict with 172.16.0.0/12.
	Subnets []Subnet `json:"subnet_ids" required:"true"`
	// Access mode.
	// + INTERNET: Indicates Internet access.
	// + DEDICATED: Indicates dedicated line access.
	// + BOTH: Indicates that both access methods are supported.
	AccessMode string `json:"access_mode" required:"true"`
	// Enterprise ID.
	// The enterprise ID is the unique identification in the workspace service.
	// If omited, the system will automatically generate an enterprise ID.
	// The ID can contain `1` to `32` characters, only letters, digits, hyphens (-) and underscores (_) are allowed.
	EnterpriseId string `json:"enterprise_id,omitempty"`
	// The CIDR of management subnet.
	// It cannot conflict with 172.16.0.0/12 and the CIDRs of service subnet.
	ManagementSubnetCidr string `json:"manage_subnet_cidr,omitempty"`
	// Dedicated subnet list.
	DedicatedSubnets string `json:"dedicated_subnets,omitempty"`
}

// Domain is an object to specified the configuration of AD domain.
type Domain struct {
	// Domain type.
	// + LITE_AS: Local authentication.
	// + LOCAL_AD: Local AD.
	// When the domain type is "LOCAL_AD", make sure that the selected VPC network and the network to which AD
	//   belongs can be connected.
	Type string `json:"domain_type" required:"domain_type"`
	// Domain name. It needs to be configured when the domain type is LOCAL_AD.
	// The domain name must be an existing domain name on the AD server, and the length should not exceed 55.
	Name string `json:"domain_name,omitempty"`
	// Domain administrator account. It needs to be configured when the domain type is "LOCAL_AD".
	// It must be an existing domain administrator account on the AD server.
	AdminAccount string `json:"domain_admin_account,omitempty"`
	// Domain administrator account password. It needs to be configured when the domain type is "LOCAL_AD".
	Password string `json:"domain_password,omitempty"`
	// Primary domain controller IP address. It needs to be configured when the domain type is LOCAL_AD.
	ActiveDomainIp string `json:"active_domain_ip,omitempty"`
	// Primary domain controller name. It needs to be configured when the domain type is LOCAL_AD.
	AcitveDomainName string `json:"active_domain_name,omitempty"`
	// The IP address of the standby domain controller.
	// It needs to be configured when the domain type is LOCAL_AD and the standby node is configured.
	StandyDomainIp string `json:"standby_domain_ip,omitempty"`
	// The name of the standby domain controller.
	// It needs to be configured when the domain type is LOCAL_AD and the standby node is configured.
	StandyDomainName string `json:"standby_domain_name,omitempty"`
	// Primary DNS IP address.
	// It needs to be configured when the domain type is LOCAL_AD.
	ActiveDnsIp string `json:"active_dns_ip,omitempty"`
	// Backup DNS IP address.
	// It needs to be configured when the domain type is LOCAL_AD and the standby node is configured.
	StandyDnsIp string `json:"standby_dns_ip,omitempty"`
	// Whether to delete the corresponding computer object on AD while deleting the desktop.
	// + 0 means not delete
	// + 1 means delete.
	DeleteComputerObject *int `json:"delete_computer_object,omitempty"`
	// Whether to enable LDAPS.
	UseLdaps bool `json:"use_idaps,omitempty"`
	// The configuration of TLS.
	TlsConfig *TlsConfig `json:"tls_config,omitempty"`
}

// TlsConfig is an object to specified the configuration TLS (SLL) certificate.
type TlsConfig struct {
	// The pem content, used to update or upload. The query will not return.
	CertPem string `json:"cert_pem,omitempty"`
	// The valid start time of the certificate, please refer to the example "2022-01-25T09:24:27".
	CertStartTime string `json:"cert_start_time,omitempty"`
	// The valid end time of the certificate, please refer to the example 2022-01-25T09:24:27.
	CertEndTime string `json:"cert_end_time,omitempty"`
}

// Subnet is an object to specified the network configuration of VPC subnet to which the service and desktops belongs.
type Subnet struct {
	// The network ID of subnet.
	NetworkId string `json:"subnet_id" required:"true"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to subscribe Workspace service using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*CreateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r CreateResp
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to obtain the Workspace serivce details.
func Get(c *golangsdk.ServiceClient) (*Service, error) {
	var r Service
	_, err := c.Get(rootURL(c), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// UpdateOpts is the structure required by the Update method to change servie configuration.
type UpdateOpts struct {
	// Configuration of domain.
	AdDomain *Domain `json:"ad_domains,omitempty"`
	// Access mode.
	// + INTERNET: Indicates Internet access.
	// + DEDICATED: Indicates dedicated line access.
	// + BOTH: Indicates that both access methods are supported.
	AccessMode string `json:"access_mode,omitempty"`
	// Dedicated subnet list.
	DedicatedSubnets string `json:"dedicated_subnets,omitempty"`
	// Service subnet which to specify the returned network ID to order desktops.
	Subnets []string `json:"subnet_ids,omitempty"`
	// Internet access port.
	InternetAccessPort string `json:"internet_access_port,omitempty"`
	// Enterprise ID.
	// The enterprise ID is the unique identification in the workspace service.
	// If omited, the system will automatically generate an enterprise ID.
	// The ID can contain `1` to `32` characters, only letters, digits, hyphens (-) and underscores (_) are allowed.
	EnterpriseId string `json:"enterprise_id,omitempty"`
}

// Update is a method to change service configuration using givin parameters.
func Update(c *golangsdk.ServiceClient, opts UpdateOpts) (*UpdateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r UpdateResp
	_, err = c.Put(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to unregistry the Workspace service using given parameters.
func Delete(c *golangsdk.ServiceClient) (*DeleteResp, error) {
	var r DeleteResp
	_, err := c.DeleteWithResponse(rootURL(c), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// GetAuthConfig is the method that used to query the configuration information of secondary authentication.
func GetAuthConfig(c *golangsdk.ServiceClient) (*OtpConfigResp, error) {
	var r OtpConfigResp
	_, err := c.Get(authConfigURL(c), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// UpdateAuthConfigOpts is the structure by the UpdateAuthConfig method to change auxiliary authentication.configuration.
type UpdateAuthConfigOpts struct {
	// Authentication type.
	// + OTP: Indicates OTP assist authentication.
	AuthType string `json:"auth_type" required:"true"`
	// The OTP auxiliary authentication method configuration.
	OptConfigInfo *OtpConfigInfo `json:"otp_config_info" required:"true"`
}

// OtpConfigInfo is the structure to specified the OTP auxiliary authentication configuration infomation.
type OtpConfigInfo struct {
	// Whether to enable OTP authentication mode.
	Enable *bool `json:"enable" required:"true"`
	// Verification code receiving mode.
	// + VMFA Indicates virtual MFA device.
	// + HMFA Indicates hardware MFA device.
	ReceiveMode string `json:"receive_mode" required:"true"`
	// Auxiliary authentication server address.
	AuthUrl string `json:"auth_url,omitempty"`
	// Auxiliary authentication service access account.
	AppId string `json:"app_id,omitempty"`
	// Auxiliary authentication service access password.
	AppSecrte string `json:"app_secret,omitempty"`
	// Auxiliary authentication service access mode.
	// + INTERNET: Indicates Internet access.
	// + DEDICATED: Indicates dedicated line access.
	// + SYSTEM_DEFAULT：Indicates system default.
	AuthServerAccessMode string `json:"auth_server_access_mode,omitempty"`
	// PEM format certificate content.
	CertContent string `json:"cert_content,omitempty"`
	// Authentication application object information. If null, it means it is effective for all application objects.
	ApplyRule *ApplyRule `json:"apply_rule,omitempty"`
}

// ApplyRule is the object to specified the OTP auxiliary authentication configuration infomation.
type ApplyRule struct {
	// Authentication application object type.
	// + ACCESS_MODE: Indicates access type.
	RuleType string `json:"rule_type,omitempty"`
	// Authentication application object.
	// + INTERNET: Indicates Internet access. Optional only when rule_type is "ACCESS_MODE".
	// + PRIVATE: Indicates dedicated line access. Optional only when rule_type is "ACCESS_MODE".
	Rule string `json:"rule,omitempty"`
}

// UpdateAssistAuthConfig is the method that used to modify the configuration information of auxiliary authentication
func UpdateAssistAuthConfig(c *golangsdk.ServiceClient, opts UpdateAuthConfigOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Put(authConfigURL(c), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
		OkCodes:     []int{204},
	})
	return err
}

// GetLockStatus is the method used to get whether the Workspace service is locked detail.
func GetLockStatus(c *golangsdk.ServiceClient) (*LockStatusResp, error) {
	var r LockStatusResp
	_, err := c.Get(lockStatusURL(c), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return &r, err
}

// UnlockOpts is the object to specified the Workspace service unlock action type.
type UnlockOpts struct {
	// Unlock action type.
	// + unlock: Indicates unlock the Workspace service.
	OperateType string `json:"operate_type" required:"true"`
}

// UnlockService is the method that used to unlock the Workspace service.
func UnlockService(c *golangsdk.ServiceClient, opts UnlockOpts) (*UnlockResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r UnlockResp
	_, err = c.Put(lockStatusURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}
