package services

import (
	"github.com/huaweicloud/golangsdk"
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
	Tags []ResourceTags `json:"tags"`
	// the error message when the status of the VPC endpoint service changes to failed
	Error []ErrorInfo `json:"error"`
	// the creation time of the VPC endpoint service
	Created string `json:"created_at"`
	// the update time of the VPC endpoint service
	Updated string `json:"updated_at"`
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

type ResourceTags struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ErrorInfo struct {
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}

type commonResult struct {
	golangsdk.Result
}

// ListResult represents the result of a list operation. Call its ExtractServices
// method to interpret it as Services.
type ListResult struct {
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
