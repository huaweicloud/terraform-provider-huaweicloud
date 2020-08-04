package tags

import (
	"github.com/huaweicloud/golangsdk"
)

// CreateOptsBuilder describes struct types that can be accepted by the Create call.
// The CreateOpts struct in this package does.
type CreateOptsBuilder interface {
	// Returns value that can be passed to json.Marshal
	ToTagsCreateMap() (map[string]interface{}, error)
}

// CreateOpts implements CreateOptsBuilder
type CreateOpts struct {
	// Tags is a set of tags.
	Tags map[string]string `json:"tags" required:"true"`
}

// ToImageCreateMap assembles a request body based on the contents of
// a CreateOpts.
func (opts CreateOpts) ToTagsCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Create implements create image request
func Create(client *golangsdk.ServiceClient, resource_type, resource_id string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTagsCreateMap()
	if err != nil {
		r.Err = err
		return r
	}
	_, r.Err = client.Put(createURL(client, resource_type, resource_id), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Get implements tags get request
func Get(client *golangsdk.ServiceClient, resource_type, resource_id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, resource_type, resource_id), &r.Body, nil)
	return
}

// Delete implements image delete request by creating empty tag map
func Delete(client *golangsdk.ServiceClient, resource_type, resource_id string) (r DeleteResult) {
	createOpts := CreateOpts{
		Tags: map[string]string{},
	}
	_, r.Err = Create(client, resource_type, resource_id, createOpts).Extract()
	return
}
