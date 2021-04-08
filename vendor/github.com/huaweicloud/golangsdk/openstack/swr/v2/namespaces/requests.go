package namespaces

import (
	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToNamespaceCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new network
type CreateOpts struct {
	Namespace string `json:"namespace" required:"true"`
}

// ToNamespaceCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToNamespaceCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new namespace.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToNamespaceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular network based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// Delete will permanently delete a particular network based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}
