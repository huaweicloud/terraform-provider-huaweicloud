package protectiongroups

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToGroupCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new group.
type CreateOpts struct {
	//Group Name
	Name string `json:"name" required:"true"`
	//Group Description
	Description string `json:"description,omitempty"`
	//The source AZ of a protection group
	SourceAZ string `json:"source_availability_zone" required:"true"`
	//The target AZ of a protection group
	TargetAZ string `json:"target_availability_zone" required:"true"`
	//An active-active domain
	DomainID string `json:"domain_id" required:"true"`
	//ID of the source VPC
	SourceVpcID string `json:"source_vpc_id" required:"true"`
	//Deployment model
	DrType string `json:"dr_type,omitempty"`
}

// ToGroupCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToGroupCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "server_group")
}

// Create will create a new Group based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToGroupCreateMap()
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
	ToGroupUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a Group.
type UpdateOpts struct {
	//Group name
	Name string `json:"name" required:"true"`
}

// ToGroupUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToGroupUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "server_group")
}

// Update accepts a UpdateOpts struct and uses the values to update a Group.The response code from api is 200
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, id), b, nil, reqOpt)
	return
}

// Get retrieves a particular Group based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// Delete will permanently delete a particular Group based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r JobResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.DeleteWithResponse(resourceURL(c, id), &r.Body, reqOpt)
	return
}

// EnableOpts contains all the values needed to enable protection for a Group.
type EnableOpts struct {
	//Empty
}

// ToGroupEnableMap builds a create request body from EnableOpts.
func (opts EnableOpts) ToGroupEnableMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "start-server-group")
}

// Enable will enable protection for a protection Group.
func Enable(c *golangsdk.ServiceClient, id string) (r JobResult) {
	opts := EnableOpts{}
	b, err := opts.ToGroupEnableMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(actionURL(c, id), b, &r.Body, reqOpt)
	return
}

// DisableOpts contains all the values needed to disable protection for a Group.
type DisableOpts struct {
	//Empty
}

// ToGroupDisableMap builds a create request body from DisableOpts.
func (opts DisableOpts) ToGroupDisableMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "stop-server-group")
}

// Disable will disable protection for a protection Group.
func Disable(c *golangsdk.ServiceClient, id string) (r JobResult) {
	opts := DisableOpts{}
	b, err := opts.ToGroupDisableMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(actionURL(c, id), b, &r.Body, reqOpt)
	return
}
