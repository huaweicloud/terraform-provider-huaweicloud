package services

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
)

// PostOptsBuilder allows extensions to add parameters to the
// Post request.
type PostOptsBuilder interface {
	ToServicePostMap() (map[string]interface{}, error)
}

// CreateOpts contains the options for create a VPC endpoint service.
// This object is passed to Create().
type CreateOpts struct {
	// Specifies the ID of the VPC to which the backend resource of the VPC endpoint service belongs.
	VpcID string `json:"vpc_id" required:"true"`
	// Specifies the ID for identifying the backend resource of the VPC endpoint service.
	PortID string `json:"port_id" required:"true"`
	// Specifies the resource type.
	ServerType string `json:"server_type" required:"true"`
	// Lists the port mappings opened to the VPC endpoint service.
	Ports []PortOpts `json:"ports" required:"true"`

	// Specifies the name of the VPC endpoint service.
	// The value contains a maximum of 16 characters, including letters, digits, underscores (_), and hyphens (-).
	ServiceName string `json:"service_name,omitempty"`
	// Specifies the type of the VPC endpoint service, only interface is valid.
	ServiceType string `json:"service_type,omitempty"`
	// Specifies whether connection approval is required.
	Approval *bool `json:"approval_enabled,omitempty"`
	// Specifies the ID of the virtual NIC to which the virtual IP address is bound.
	VipPortID string `json:"vip_port_id,omitempty"`
	// Specifies whether the client IP address and port number or marker_id information is transmitted to the server.
	TCPProxy string `json:"tcp_proxy,omitempty"`
	// Specifies the resource tags in key/value format
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

// PortOpts contains the port mappings opened to the VPC endpoint service.
type PortOpts struct {
	// Specifies the protocol used in port mappings. The value can be TCP or UDP. The default value is TCP.
	Protocol string `json:"protocol,omitempty"`
	// Specifies the port for accessing the VPC endpoint.
	ClientPort int `json:"client_port,omitempty"`
	// Specifies the port for accessing the VPC endpoint service.
	ServerPort int `json:"server_port,omitempty"`
}

// ToServicePostMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToServicePostMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// VPC endpoint service.
func Create(c *golangsdk.ServiceClient, opts PostOptsBuilder) (r CreateResult) {
	b, err := opts.ToServicePostMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular nodes based on its unique ID and cluster ID.
func Get(c *golangsdk.ServiceClient, serviceID string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, serviceID), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToServiceUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a VPC endpoint service
type UpdateOpts struct {
	// Specifies the name of the VPC endpoint service.
	ServiceName string `json:"service_name,omitempty"`
	// Specifies whether connection approval is required.
	Approval *bool `json:"approval_enabled,omitempty"`
	// Specifies the ID for identifying the backend resource of the VPC endpoint service.
	PortID string `json:"port_id,omitempty"`
	// Lists the port mappings opened to the VPC endpoint service.
	Ports []PortOpts `json:"ports,omitempty"`
	// Specifies the ID of the virtual NIC to which the virtual IP address is bound.
	VipPortID string `json:"vip_port_id,omitempty"`
}

// ToServiceUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOpts) ToServiceUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update allows a VPC endpoint service to be updated.
func Update(c *golangsdk.ServiceClient, serviceID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToServiceUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, serviceID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular node based on its unique ID and cluster ID.
func Delete(c *golangsdk.ServiceClient, serviceID string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, serviceID), nil)
	return
}

// ListOptsBuilder allows extensions to add parameters to the
// List request.
type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

// ListOpts allows the filtering of list data using given parameters.
type ListOpts struct {
	ServiceName string `q:"endpoint_service_name"`
	ID          string `q:"id"`
	// Status is not supported for ListPublic
	Status string `q:"status"`
}

// ToListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list VPC endpoint services.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) ([]Service, error) {
	var r ListResult
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}
	_, r.Err = client.Get(url, &r.Body, nil)
	if r.Err != nil {
		return nil, r.Err
	}

	allNodes, err := r.ExtractServices()
	if err != nil {
		return nil, err
	}

	return allNodes, nil
}

// ListPublic makes a request against the API to list public VPC endpoint services.
func ListPublic(client *golangsdk.ServiceClient, opts ListOptsBuilder) ([]PublicService, error) {
	var r ListPublicResult
	url := publicResourceURL(client)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}
	_, r.Err = client.Get(url, &r.Body, nil)
	if r.Err != nil {
		return nil, r.Err
	}

	allNodes, err := r.ExtractServices()
	if err != nil {
		return nil, err
	}

	return allNodes, nil
}

// ConnActionOpts used to receive or reject a VPC endpoint for a VPC endpoint service.
type ConnActionOpts struct {
	// Specifies whether to receive or reject a VPC endpoint for a VPC endpoint service.
	Action string `json:"action" required:"true"`
	// Lists the VPC endpoints.
	Endpoints []string `json:"endpoints" required:"true"`
}

// ToServicePostMap assembles a request body based on the contents of a ConnActionOpts.
func (opts ConnActionOpts) ToServicePostMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// ConnAction accepts a ConnActionOpts struct and uses the values to receive or reject
// a VPC endpoint for a VPC endpoint service.
func ConnAction(c *golangsdk.ServiceClient, serviceID string, opts PostOptsBuilder) (r ConnectionResult) {
	b, err := opts.ToServicePostMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(connectionsActionURL(c, serviceID), b, &r.Body, reqOpt)
	return
}

// ListConnOpts used to query connections of a VPC endpoint service.
type ListConnOpts struct {
	// Specifies the unique ID of the VPC endpoint
	EndpointID string `q:"id"`
	// Specifies the packet ID of the VPC endpoint
	MarkerID string `q:"marker_id"`
}

// ToListQuery formats a ListConnOpts into a query string.
func (opts ListConnOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// ListConnections makes a request against the API to list connections of a VPC endpoint service.
func ListConnections(client *golangsdk.ServiceClient, serviceID string, opts ListOptsBuilder) ([]Connection, error) {
	var r ConnectionResult
	url := connectionsURL(client, serviceID)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}
	_, r.Err = client.Get(url, &r.Body, nil)
	if r.Err != nil {
		return nil, r.Err
	}

	allConnections, err := r.ExtractConnections()
	if err != nil {
		return nil, err
	}

	return allConnections, nil
}

// PermActionOpts used to add to or delete whitelist records from a VPC endpoint service.
type PermActionOpts struct {
	// Specifies the operation to be performed: dd or remove.
	Action string `json:"action" required:"true"`
	// Lists the whitelist records. The record is in the iam:domain::domain_id format.
	Permissions []string `json:"permissions" required:"true"`
}

// ToServicePostMap assembles a request body based on the contents of a PermActionOpts.
func (opts PermActionOpts) ToServicePostMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// PermAction accepts a PermActionOpts struct and uses the values toadd to or delete
// whitelist records from a VPC endpoint service.
func PermAction(c *golangsdk.ServiceClient, serviceID string, opts PostOptsBuilder) (r PermActionResult) {
	b, err := opts.ToServicePostMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(permissionsActionURL(c, serviceID), b, &r.Body, reqOpt)
	return
}

// ListPermissions makes a request against the API to query the whitelist records of
// a VPC endpoint service.
func ListPermissions(client *golangsdk.ServiceClient, serviceID string) ([]Permission, error) {
	var r ListPermResult
	url := permissionsURL(client, serviceID)

	_, r.Err = client.Get(url, &r.Body, nil)
	if r.Err != nil {
		return nil, r.Err
	}

	allPermissions, err := r.ExtractPermissions()
	if err != nil {
		return nil, err
	}

	return allPermissions, nil
}
