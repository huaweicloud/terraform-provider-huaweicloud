package endpoints

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
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
	// Specifies whether to create a private domain name
	EnableDNS *bool `json:"enable_dns,omitempty"`
	// Specifies the resource tags in key/value format
	Tags []tags.ResourceTag `json:"tags,omitempty"`
	// Specifies the IDs of route tables
	// The parameter is mandatory to create a gateway type VPC endpoint
	// If the parameter is not set, will be associated with the default route table
	RouteTables []string `json:"routetables,omitempty"`
	// Specifies the IP address for accessing the associated VPC endpoint service
	PortIP string `json:"port_ip,omitempty"`
	// Specifies the whitelist for controlling access to the VPC endpoint
	Whitelist []string `json:"whitelist,omitempty"`
	// Specifies whether to enable access control
	EnableWhitelist *bool `json:"enable_whitelist,omitempty"`
	// Specifies the description of the VPC endpoint service
	Description string `json:"description,omitempty"`
	// Specifies the endpoint policy information for the gateway type
	PolicyStatement []PolicyStatement `json:"policy_statement,omitempty"`
	// Specifies the IP version
	IPVersion string `json:"ip_version,omitempty"`
	// Specifies the IPv6 address
	IPv6Address string `json:"ipv6_address,omitempty"`
}

// PolicyStatement represents the Statement of the gateway
type PolicyStatement struct {
	Effect    string                 `json:"Effect" required:"true"`
	Action    []string               `json:"Action" required:"true"`
	Resource  []string               `json:"Resource" required:"true"`
	Condition map[string]interface{} `json:"Condition,omitempty"`
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

// UpdateOptsBuilder allows extensions to add parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToEndpointUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the options for update a VPC endpoint
// This object is passed to Update().
type UpdateOpts struct {
	// Specifies whether to enable access control
	EnableWhitelist *bool `json:"enable_whitelist,omitempty"`
	// Specifies the whitelist for controlling access to the VPC endpoint.
	// If the value is [], means delete all white list
	Whitelist []string `json:"whitelist"`
}

// ToEndpointUpdateMap assembles a request body based on the contents of a UpdateOpts.
func (opts UpdateOpts) ToEndpointUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and uses the values to update a VPC endpoint
func Update(c *golangsdk.ServiceClient, opts UpdateOptsBuilder, endpointID string) (r UpdateResult) {
	b, err := opts.ToEndpointUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Put(resourceURL(c, endpointID), b, &r.Body, nil)
	return
}

// UpdatePolicyOptsBuilder using to update endpoint policy
type UpdatePolicyOptsBuilder interface {
	UpdatePolicyOptsMap() (map[string]interface{}, error)
}

// UpdatePolicyOpts using to pass to UpdatePolicy()
type UpdatePolicyOpts struct {
	// Specifies the endpoint policy information for the gateway type
	PolicyStatement []PolicyStatement `json:"policy_statement,omitempty"`
}

// UpdatePolicyOptsMap assembles a request body based on the contents of a UpdatePolicyOpts.
func (opts UpdatePolicyOpts) UpdatePolicyOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// UpdatePolicy accepts a UpdatePolicyOpts struct and uses the values to update endpoint policy
func UpdatePolicy(c *golangsdk.ServiceClient, opts UpdatePolicyOptsBuilder, endpointID string) (r UpdateResult) {
	b, err := opts.UpdatePolicyOptsMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Put(updatePolicyURL(c, endpointID), b, &r.Body, nil)
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
