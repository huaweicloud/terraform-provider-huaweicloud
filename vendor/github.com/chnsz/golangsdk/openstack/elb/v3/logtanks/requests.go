package logtanks

import (
	"github.com/chnsz/golangsdk"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToLogTanksCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// The LoadBalancer on which the log will be associated with.
	LoadbalancerID string `json:"loadbalancer_id" required:"true"`

	// The log group on which the log will be associated with.
	LogGroupId string `json:"log_group_id" required:"true"`

	// The topic on which the log will subscribe.
	LogTopicId string `json:"log_topic_id" required:"true"`
}

// ToLogTanksCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToLogTanksCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "logtank")
}

// Create is an operation which provisions a new Logtanks based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToLogTanksCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// Get retrieves a particular Logtanks based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToLogTanksUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// The log group on which the log will be associated with.
	LogGroupId string `json:"log_group_id,omitempty"`

	// The topic on which the log will subscribe.
	LogTopicId string `json:"log_topic_id,omitempty"`
}

// ToLogTanksUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToLogTanksUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "logtank")
}

// Update is an operation which modifies the attributes of the specified
// Logtank.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToLogTanksUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular Logtank based on its
// unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
