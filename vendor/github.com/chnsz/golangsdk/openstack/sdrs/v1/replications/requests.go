package replications

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToReplicationCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new replication.
type CreateOpts struct {
	//Protection Group ID
	GroupID string `json:"server_group_id" required:"true"`
	//Volume ID
	VolumeID string `json:"volume_id" required:"true"`
	//Replication Name
	Name string `json:"name" required:"true"`
	//Replication Description
	Description string `json:"description,omitempty"`
}

// ToReplicationCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToReplicationCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "replication")
}

// Create will create a new Replication based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToReplicationCreateMap()
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
	ToReplicationUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a Replication.
type UpdateOpts struct {
	//Replication name
	Name string `json:"name" required:"true"`
}

// ToReplicationUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToReplicationUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "replication")
}

// Update accepts a UpdateOpts struct and uses the values to update a Replication.The response code from api is 200
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToReplicationUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, id), b, nil, reqOpt)
	return
}

// Get retrieves a particular Replication based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// DeleteOptsBuilder allows extensions to add additional parameters to the
// Delete request.
type DeleteOptsBuilder interface {
	ToReplicationDeleteMap() (map[string]interface{}, error)
}

// DeleteOpts contains all the values needed to delete a Replication.
type DeleteOpts struct {
	//Group ID
	GroupID string `json:"server_group_id,omitempty"`
	//Delete Target Volume
	DeleteVolume bool `json:"delete_target_volume,omitempty"`
}

// ToReplicationDeleteMap builds a update request body from DeleteOpts.
func (opts DeleteOpts) ToReplicationDeleteMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "replication")
}

// Delete will permanently delete a particular Replication based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string, opts DeleteOptsBuilder) (r JobResult) {
	b, err := opts.ToReplicationDeleteMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.DeleteWithBodyResp(resourceURL(c, id), b, &r.Body, reqOpt)
	return
}
