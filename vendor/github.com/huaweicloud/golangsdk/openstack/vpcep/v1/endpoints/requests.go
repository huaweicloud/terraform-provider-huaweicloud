package endpoints

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
)

// CreateOptsBuilder allows extensions to add parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToEndpointCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains the options for create a VPC endpoint
// This object is passed to Create().
type CreateOpts struct {
	// Specifies the ID of the VPC endpoint service
	ServiceID string `json:"endpoint_service_id" required:"true"`
	// Specifies the ID of the VPC where the VPC endpoint is to be created
	VpcID string `json:"vpc_id" required:"true"`

	// Specifies the network ID of the subnet created in the VPC specified by vpc_id
	// The parameter is mandatory to create an interface VPC endpoint
	SubnetID string `json:"subnet_id,omitempty"`
	// Specifies the IP address for accessing the associated VPC endpoint service
	PortIP string `json:"port_ip,omitempty"`
	// Specifies whether to create a private domain name
	EnableDNS *bool `json:"enable_dns,omitempty"`
	// Specifies whether to enable access control
	EnableWhitelist *bool `json:"enable_whitelist,omitempty"`
	// Specifies the whitelist for controlling access to the VPC endpoint
	Whitelist []string `json:"whitelist,omitempty"`
	// Specifies the IDs of route tables
	RouteTables []string `json:"routeTables,omitempty"`
	// Specifies the resource tags in key/value format
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

// ToEndpointCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToEndpointCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// VPC endpoint
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToEndpointCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular VPC endpoint based on its unique ID
func Get(c *golangsdk.ServiceClient, endpointID string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, endpointID), &r.Body, nil)
	return
}

// Delete will permanently delete a particular VPC endpoint based on its unique ID
func Delete(c *golangsdk.ServiceClient, endpointID string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, endpointID), nil)
	return
}

// ListOptsBuilder allows extensions to add parameters to the
// List request.
type ListOptsBuilder interface {
	ToEndpointListQuery() (string, error)
}

// ListOpts allows the filtering of list data using given parameters.
type ListOpts struct {
	ServiceName string `q:"endpoint_service_name"`
	VPCID       string `q:"vpc_id"`
	ID          string `q:"id"`
}

// ToEndpointListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToEndpointListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list VPC endpoints.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) ([]Endpoint, error) {
	var r ListResult
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToEndpointListQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}
	_, r.Err = client.Get(url, &r.Body, nil)
	if r.Err != nil {
		return nil, r.Err
	}

	allEndpoints, err := r.ExtractEndpoints()
	if err != nil {
		return nil, err
	}

	return allEndpoints, nil
}
