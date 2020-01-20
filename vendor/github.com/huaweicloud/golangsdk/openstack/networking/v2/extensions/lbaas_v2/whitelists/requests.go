package whitelists

import (
	"github.com/huaweicloud/golangsdk"
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToWhitelistCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	TenantId        string `json:"tenant_id,omitempty"`
	ListenerId      string `json:"listener_id" required:"true"`
	EnableWhitelist *bool  `json:"enable_whitelist,omitempty"`
	Whitelist       string `json:"whitelist,omitempty"`
}

// ToWhitelistCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToWhitelistCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "whitelist")
}

// Create is an operation which provisions a new whitelist based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create loadbalancers on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToWhitelistCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// Get retrieves a particular Whitelist based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToWhitelistUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	EnableWhitelist *bool  `json:"enable_whitelist,omitempty"`
	Whitelist       string `json:"whitelist,omitempty"`
}

// ToWhitelistUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToWhitelistUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "whitelist")
}

// Update is an operation which modifies the attributes of the specified Whitelist.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToWhitelistUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular Whitelist based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
