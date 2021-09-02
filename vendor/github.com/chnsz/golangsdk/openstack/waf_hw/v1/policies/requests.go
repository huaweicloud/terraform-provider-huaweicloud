package policies

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToPolicyCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new policy.
type CreateOpts struct {
	//Policy name
	Name string `json:"name" required:"true"`
}

// ToPolicyCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToPolicyCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new policy based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPolicyCreateMap()
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
	ToPolicyUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a policy.
type UpdateOpts struct {
	Name          string        `json:"name,omitempty"`
	Action        *Action       `json:"action,omitempty"`
	Options       *PolicyOption `json:"options,omitempty"`
	Level         int           `json:"level,omitempty"`
	FullDetection *bool         `json:"full_detection,omitempty"`
}

// ToPolicyUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToPolicyUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and uses the values to update a policy.The response code from api is 200
func Update(c *golangsdk.ServiceClient, policyID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPolicyUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Patch(resourceURL(c, policyID), b, nil, reqOpt)
	return
}

// UpdateHostsOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateHostsOptsBuilder interface {
	ToUpdateHostsQuery() (string, error)
}

// UpdateHostsOpts contains all the values needed to update a policy hosts.
type UpdateHostsOpts struct {
	//Domain ID
	Hosts []string `q:"hosts" required:"true"`
}

// ListPolicyOpts
type ListPolicyOpts struct {
	Page     int `q:"page"`
	Pagesize int `q:"pagesize"`
	// policy name
	Name string `q:"name"`
}

// ToUpdateHostsQuery builds a update request query from UpdateHostsOpts.
func (opts UpdateHostsOpts) ToUpdateHostsQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// UpdateHosts accepts a UpdateHostsOpts struct and uses the values to update a policy hosts.The response code from api is 200
func UpdateHosts(c *golangsdk.ServiceClient, policyId string, opts UpdateHostsOptsBuilder) (r UpdateResult) {
	url := resourceURL(c, policyId)
	if opts != nil {
		var query string
		query, r.Err = opts.ToUpdateHostsQuery()
		if r.Err != nil {
			return
		}
		url += query
	}
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	}
	_, r.Err = c.Put(url, nil, r.Body, reqOpt)
	return
}

// Get retrieves a particular policy based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	}
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, reqOpt)
	return
}

// Delete will permanently delete a particular policy based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Delete(resourceURL(c, id), reqOpt)
	return
}

// ListPolicy retrieve waf policy by ListPolicyOpts
func ListPolicy(c *golangsdk.ServiceClient, opts ListPolicyOpts) (*ListPolicyRst, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst golangsdk.Result
	_, err = c.Get(url, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	if err == nil {
		var r ListPolicyRst
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}
