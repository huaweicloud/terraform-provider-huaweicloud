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
	StandyDomainIp string `json:"standy_domain_ip,omitempty"`
	// The name of the standby domain controller.
	// It needs to be configured when the domain type is LOCAL_AD and the standby node is configured.
	StandyDomainName string `json:"standy_domain_name,omitempty"`
	// Primary DNS IP address.
	// It needs to be configured when the domain type is LOCAL_AD.
	ActiveDnsIp string `json:"active_dns_ip,omitempty"`
	// Backup DNS IP address.
	// It needs to be configured when the domain type is LOCAL_AD and the standby node is configured.
	StandyDnsIp string `json:"standy_dns_ip,omitempty"`
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
