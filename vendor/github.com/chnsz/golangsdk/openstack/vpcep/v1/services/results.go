package services

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

// Service contains the response of the VPC endpoint service
type Service struct {
	// the ID of the VPC endpoint service
	ID string `json:"id"`
	// the status of the VPC endpoint service
	Status string `json:"status"`
	// the ID for identifying the backend resource of the VPC endpoint service
	PortID string `json:"port_id"`
	// the ID of the VPC to which the backend resource of the VPC endpoint service belongs
	VpcID string `json:"vpc_id"`
	// the name of the VPC endpoint service
	ServiceName string `json:"service_name"`
	// the type of the VPC endpoint service
	ServiceType string `json:"service_type"`
	// the resource type
	ServerType string `json:"server_type"`
	// whether connection approval is required
	Approval bool `json:"approval_enabled"`
	// the ID of the virtual NIC to which the virtual IP address is bound
	VipPortID string `json:"vip_port_id"`
	// the project ID
	ProjectID string `json:"project_id"`
	// the network segment type. The value can be public or internal
	CidrType string `json:"cidr_type"`
	// Lists the port mappings opened to the VPC endpoint service
	Ports []PortMapping `json:"ports"`
	// whether the client IP address and port number or marker_id information is transmitted to the server
	TCPProxy string `json:"tcp_proxy"`
	// the resource tags
	Tags []tags.ResourceTag `json:"tags"`
	// the error message when the status of the VPC endpoint service changes to failed
	Error []ErrorInfo `json:"error"`
	// the description of the VPC endpoint service
	Description string `json:"description"`
	// the IP version of the VPC endpoint service
	IpVersion string `json:"ip_version"`
	// whether the VPC endpoint policy is enabled.
	EnablePolicy bool `json:"enable_policy"`
	// the creation time of the VPC endpoint service
	Created string `json:"created_at"`
	// the update time of the VPC endpoint service
	Updated string `json:"updated_at"`
	// the VPC endpoint service that matches the edge attribute in the filtering result.
	PublicBorderGroup string `json:"public_border_group"`
	// the number of VPC endpoints that are in the Creating or Accepted status.
	ConnectionCount int `json:"connection_count"`
	// the network ID of any subnet within the VPC used to create the VPC endpoint service
	SnatNetworkId string `json:"snat_network_id"`
	// the IPv4 address or domain name of the server in the interface type VLAN scenario
	IpAddress string `json:"ip"`
	// the dedicated cluster ID associated with the VPC endpoint service
	PoolId string `json:"pool_id"`
}

// PortMapping contains the port mappings opened to the VPC endpoint service
type PortMapping struct {
	// the protocol used in port mappings. The value can be TCP or UDP.
	Protocol string `json:"protocol"`
	// the port for accessing the VPC endpoint
	ClientPort int `json:"client_port"`
	// the port for accessing the VPC endpoint service
	ServerPort int `json:"server_port"`
}

type ErrorInfo struct {
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}

// PublicService contains the response of the public VPC endpoint service
type PublicService struct {
	// the ID of the public VPC endpoint service
	ID string `json:"id"`
	// the owner of the VPC endpoint service
	Owner string `json:"owner"`
	// the name of the VPC endpoint service
	ServiceName string `json:"service_name"`
	// the type of the VPC endpoint service: gateway or interface
	ServiceType string `json:"service_type"`
	// whether the associated VPC endpoint carries a charge: true or false
	IsChange bool `json:"is_charge"`
	// the creation time of the VPC endpoint service
	Created string `json:"created_at"`
}

type commonResult struct {
	golangsdk.Result
}

// ListResult represents the result of a list operation. Call its ExtractServices
// method to interpret it as Services.
type ListResult struct {
	commonResult
}

// ListPublicResult represents the result of a list public operation. Call its ExtractServices
// method to interpret it as PublicServices.
type ListPublicResult struct {
	commonResult
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Service.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Service.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Service.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

// Extract is a function that accepts a result and extracts a Service.
func (r commonResult) Extract() (*Service, error) {
	var s Service
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractServices is a function that accepts a result and extracts the given Services
func (r ListResult) ExtractServices() ([]Service, error) {
	var s struct {
		Services []Service `json:"endpoint_services"`
	}

	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Services, nil
}

// ExtractServices is a function that accepts a result and extracts the given PublicService
func (r ListPublicResult) ExtractServices() ([]PublicService, error) {
	var s struct {
		Services []PublicService `json:"endpoint_services"`
	}

	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Services, nil
}

// Connection contains the response of querying Connections of a VPC Endpoint Service
type Connection struct {
	// the ID of the VPC endpoint
	EndpointID string `json:"id"`
	// the packet ID of the VPC endpoint
	MarkerID int `json:"marker_id"`
	// the ID of the user's domain
	DomainID string `json:"domain_id"`
	// the connection status of the VPC endpoint
	Status string `json:"status"`
	// the creation time of the VPC endpoint
	Created string `json:"created_at"`
	// the update time of the VPC endpoint
	Updated string `json:"updated_at"`
	// the error message when the status of the VPC endpoint service changes to failed
	Error []ErrorInfo `json:"error"`
	// the description of endpoint service connection
	Description string `json:"description"`
}

// ConnectionResult represents the result of a list connections.
// Call its ExtractConnections method to interpret it as []Connection.
type ConnectionResult struct {
	commonResult
}

// ExtractConnections is a function that accepts a result and extracts the given []Connection
func (r ConnectionResult) ExtractConnections() ([]Connection, error) {
	var s struct {
		Connections []Connection `json:"connections"`
	}

	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Connections, nil
}

// Permission contains the response of querying Permissions of a VPC Endpoint Service
type Permission struct {
	// the unique ID of the permission.
	ID string `json:"id"`
	// the whitelist records.
	Permission string `json:"permission"`
	// the type of the whitelist.
	PermissionType string `json:"permission_type"`
	// the time of adding the whitelist record
	Created string `json:"created_at"`
}

type PermActionResult struct {
	commonResult
}

// ListPermResult represents the result of a list permissions. Call its ExtractPermissions
// method to interpret it as []Permission.
type ListPermResult struct {
	commonResult
}

// ExtractPermissions is a function that accepts a result and extracts the given []Permission
func (r ListPermResult) ExtractPermissions() ([]Permission, error) {
	var s struct {
		Permissions []Permission `json:"permissions"`
	}

	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Permissions, nil
}

// PublicServicePage is a single page maximum result representing a query by offset page.
type PublicServicePage struct {
	pagination.OffsetPageBase
}

// extractPublicService is a method to extract the list of tags supported PublicService.
func extractPublicService(r pagination.Page) ([]PublicService, error) {
	var s []PublicService
	err := r.(PublicServicePage).Result.ExtractIntoSlicePtr(&s, "endpoint_services")
	return s, err
}

// extractService is a method to extract the list of tags supported Services.
func extractService(r pagination.Page) ([]Service, error) {
	var s []Service
	err := r.(PublicServicePage).Result.ExtractIntoSlicePtr(&s, "endpoint_services")
	return s, err
}

// ConnectionPage is a single page maximum result representing a query by offset page.
type ConnectionPage struct {
	pagination.OffsetPageBase
}

// extractConnection is a method to extract the list of tags supported connection.
func extractConnection(r pagination.Page) ([]Connection, error) {
	var s []Connection
	err := r.(ConnectionPage).Result.ExtractIntoSlicePtr(&s, "connections")
	return s, err
}

// PermissionPage is a single page maximum result representing a query by offset page.
type PermissionPage struct {
	pagination.OffsetPageBase
}

// extractPermission is a method to extract the list of tags supported permission.
func extractPermission(r pagination.Page) ([]Permission, error) {
	var s []Permission
	err := r.(PermissionPage).Result.ExtractIntoSlicePtr(&s, "permissions")
	return s, err
}
