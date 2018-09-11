package floatingips

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// List returns a Pager that allows you to iterate over a collection of FloatingIPs.
func List(client *golangsdk.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return FloatingIPPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToFloatingIPCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies a Floating IP allocation request.
type CreateOpts struct {
	// Pool is the pool of Floating IPs to allocate one from.
	Pool string `json:"pool" required:"true"`
}

// ToFloatingIPCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToFloatingIPCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create requests the creation of a new Floating IP.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToFloatingIPCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get returns data about a previously created Floating IP.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// Delete requests the deletion of a previous allocated Floating IP.
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// AssociateOptsBuilder allows extensions to add additional parameters to the
// Associate request.
type AssociateOptsBuilder interface {
	ToFloatingIPAssociateMap() (map[string]interface{}, error)
}

// AssociateOpts specifies the required information to associate a Floating IP with an instance
type AssociateOpts struct {
	// FloatingIP is the Floating IP to associate with an instance.
	FloatingIP string `json:"address" required:"true"`

	// FixedIP is an optional fixed IP address of the server.
	FixedIP string `json:"fixed_address,omitempty"`
}

// ToFloatingIPAssociateMap constructs a request body from AssociateOpts.
func (opts AssociateOpts) ToFloatingIPAssociateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "addFloatingIp")
}

// AssociateInstance pairs an allocated Floating IP with a server.
func AssociateInstance(client *golangsdk.ServiceClient, serverID string, opts AssociateOptsBuilder) (r AssociateResult) {
	b, err := opts.ToFloatingIPAssociateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(associateURL(client, serverID), b, nil, nil)
	return
}

// DisassociateOptsBuilder allows extensions to add additional parameters to
// the Disassociate request.
type DisassociateOptsBuilder interface {
	ToFloatingIPDisassociateMap() (map[string]interface{}, error)
}

// DisassociateOpts specifies the required information to disassociate a
// Floating IP with a server.
type DisassociateOpts struct {
	FloatingIP string `json:"address" required:"true"`
}

// ToFloatingIPDisassociateMap constructs a request body from DisassociateOpts.
func (opts DisassociateOpts) ToFloatingIPDisassociateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "removeFloatingIp")
}

// DisassociateInstance decouples an allocated Floating IP from an instance
func DisassociateInstance(client *golangsdk.ServiceClient, serverID string, opts DisassociateOptsBuilder) (r DisassociateResult) {
	b, err := opts.ToFloatingIPDisassociateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(disassociateURL(client, serverID), b, nil, nil)
	return
}
