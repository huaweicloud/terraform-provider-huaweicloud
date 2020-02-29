package groups

import (
	"github.com/huaweicloud/golangsdk"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToGroupCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a API group. This object is passed to
// the API groups Create function.
type CreateOpts struct {
	// Name of the API group
	Name string `json:"name" required:"true"`
	// Description of the API group
	Remark string `json:"remark,omitempty"`
}

// ToGroupCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToGroupCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new API group based on the values in CreateOpts.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToGroupUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contain options for updating an existing API group.
type UpdateOpts struct {
	// Name of the API group
	Name string `json:"name" required:"true"`
	// Description of the API group
	Remark string `json:"remark,omitempty"`
}

// ToGroupUpdateMap assembles a request body based on the contents of an
// UpdateOpts.
func (opts UpdateOpts) ToGroupUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update will update the API group with provided information. To extract the updated
// API group from the response, call the Extract method on the UpdateResult.
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(groupURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will delete the existing group with the provided ID.
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(groupURL(client, id), nil)
	return
}

// Get retrieves the group with the provided ID. To extract the Group object
// from the response, call the Extract method on the GetResult.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(groupURL(client, id), &r.Body, nil)
	return
}
