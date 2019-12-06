package configurations

import (
	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToConfigCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new configuration.
type CreateOpts struct {
	//Configuration Name
	Name string `json:"name" required:"true"`
	//Configuration Description
	Description string `json:"description,omitempty"`
	//Configuration Values
	Values map[string]string `json:"values,omitempty"`
	//Database Object
	DataStore DataStore `json:"datastore" required:"true"`
}

type DataStore struct {
	//DB Engine
	Type string `json:"type" required:"true"`
	//DB version
	Version string `json:"version" required:"true"`
}

// ToConfigCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToConfigCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new Config based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToConfigCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToConfigUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a Configuration.
type UpdateOpts struct {
	//Configuration Name
	Name string `json:"name,omitempty"`
	//Configuration Description
	Description string `json:"description,omitempty"`
	//Configuration Values
	Values map[string]string `json:"values,omitempty"`
}

// ToConfigUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToConfigUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and uses the values to update a Configuration.The response code from api is 200
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToConfigUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Put(resourceURL(c, id), b, nil, reqOpt)
	return
}

// Get retrieves a particular Configuration based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, &RequestOpts)
	return
}

// Delete will permanently delete a particular Configuration based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Delete(resourceURL(c, id), reqOpt)
	return
}
