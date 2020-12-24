package endpoints

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
)

// Endpoint contains the response of the VPC endpoint
type Endpoint struct {
	// the ID of the VPC endpoint
	ID string `json:"id"`
	// the connection status of the VPC endpoint
	Status string `json:"status"`
	// the account status: frozen or active
	ActiveStatus []string `json:"active_status"`
	// the type of the VPC endpoint service that is associated with the VPC endpoint
	ServiceType string `json:"service_type"`
	// the name of the VPC endpoint service
	ServiceName string `json:"endpoint_service_name"`
	// the ID of the VPC endpoint service
	ServiceID string `json:"endpoint_service_id"`
	// the ID of the VPC where the VPC endpoint is to be created
	VpcID string `json:"vpc_id"`
	// the network ID of the subnet in the VPC specified by vpc_id
	SubnetID string `json:"subnet_id"`
	// the IP address for accessing the associated VPC endpoint service
	IPAddr string `json:"ip"`
	// the packet ID of the VPC endpoint
	MarkerID int `json:"marker_id"`
	// whether to create a private domain name
	EnableDNS bool `json:"enable_dns"`
	// the domain name for accessing the associated VPC endpoint service
	DNSNames []string `json:"dns_names"`
	// whether to enable access control
	EnableWhitelist bool `json:"enable_whitelist"`
	// the whitelist for controlling access to the VPC endpoint
	Whitelist []string `json:"whitelist"`
	// the IDs of route tables
	RouteTables []string `json:"routetables"`
	// the resource tags
	Tags []tags.ResourceTag `json:"tags"`
	// the project ID
	ProjectID string `json:"project_id"`
	// the creation time of the VPC endpoint
	Created string `json:"created_at"`
	// the update time of the VPC endpoint
	Updated string `json:"updated_at"`
}

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Endpoint.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Endpoint.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

// ListResult represents the result of a list operation. Call its ExtractEndpoints
// method to interpret it as Endpoints.
type ListResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a Endpoint
func (r commonResult) Extract() (*Endpoint, error) {
	var ep Endpoint
	err := r.ExtractInto(&ep)
	return &ep, err
}

// ExtractEndpoints is a function that accepts a result and extracts the given Endpoints
func (r ListResult) ExtractEndpoints() ([]Endpoint, error) {
	var s struct {
		Endpoints []Endpoint `json:"endpoints"`
	}

	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Endpoints, nil
}
