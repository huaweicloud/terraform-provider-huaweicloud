package certificates

import (
	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToCertCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new certificate.
type CreateOpts struct {
	// Certificate name
	Name string `json:"name" required:"true"`
	// Certificate content
	Content string `json:"content" required:"true"`
	// Private Key
	Key string `json:"key" required:"true"`
}

// ToCertCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToCertCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new certificate based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCertCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToCertUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a certificate.
type UpdateOpts struct {
	//Certificate name
	Name string `json:"name,omitempty"`
}

// ToCertUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToCertUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and uses the values to update a certificate.The response code from api is 200
func Update(c *golangsdk.ServiceClient, certID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToCertUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, certID), b, nil, reqOpt)
	return
}

// Get retrieves a particular certificate based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	}
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, reqOpt)
	return
}

// Delete will permanently delete a particular certificate based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200, 204},
		MoreHeaders: RequestOpts.MoreHeaders,
	}
	_, r.Err = c.Delete(resourceURL(c, id), reqOpt)
	return
}
