package ipgroups

import (
	"github.com/chnsz/golangsdk"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToIpGroupCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options for creating a ipgroup.
type CreateOpts struct {
	// Human-readable name for the IpGroup. Does not have to be unique.
	Name string `json:"name,omitempty"`

	// Human-readable description for the IpGroup.
	Description *string `json:"description,omitempty"`

	// A list of IP addresses.
	IpList *[]IpListOpt `json:"ip_list" required:"true"`

	// Specifies the enterprise project id.
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
}

type IpListOpt struct {
	Ip          string `json:"ip" required:"true"`
	Description string `json:"description,omitempty"`
}

// ToIpGroupCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToIpGroupCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "ipgroup")
}

// Create is an operation which provisions a new IpGroups based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create IpGroups on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToIpGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// Get retrieves a particular Ipgroups based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToIpGroupUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options for updating a IpGroup.
type UpdateOpts struct {
	// Human-readable name for the IpGroup. Does not have to be unique.
	Name string `json:"name,omitempty"`

	// Human-readable description for the IpGroup.
	Description *string `json:"description,omitempty"`

	// A list of IP addresses.
	IpList *[]IpListOpt `json:"ip_list,omitempty"`
}

// ToIpGroupUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToIpGroupUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "ipgroup")
}

// Update is an operation which modifies the attributes of the specified
// IpGroup.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToIpGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

// Delete will permanently delete a particular Ipgroups based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
