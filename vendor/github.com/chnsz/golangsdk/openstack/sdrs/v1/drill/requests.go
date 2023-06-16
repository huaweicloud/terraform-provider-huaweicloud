package drill

import (
	"github.com/chnsz/golangsdk"
)

var requestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToDrillCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new dr-drill.
type CreateOpts struct {
	// Protection Group ID
	GroupID string `json:"server_group_id" required:"true"`
	//DR-Drill Name
	Name string `json:"name" required:"true"`
	// Drill vpc id
	DrillVpcID string `json:"drill_vpc_id,omitempty"`
}

// ToDrillCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToDrillCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "disaster_recovery_drill")
}

// Create will create a new DR-Drill based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToDrillCreateMap()
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
	ToDrillUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a dr-drill.
type UpdateOpts struct {
	// DR-Drill name
	Name string `json:"name" required:"true"`
}

// ToDrillUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToDrillUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "disaster_recovery_drill")
}

// Update accepts a UpdateOpts struct and uses the values to update a dr-drill.
// The response code from api is 200
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToDrillUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, id), b, nil, reqOpt)
	return
}

// Get retrieves a particular dr-drill based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}

// Delete will permanently delete a particular dr-drill based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r JobResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.DeleteWithResponse(resourceURL(c, id), &r.Body, reqOpt)
	return
}
