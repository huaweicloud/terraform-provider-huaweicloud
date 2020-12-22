package policies

import (
	"github.com/huaweicloud/golangsdk"
)

// Get retrieves details on a single policy, by ID.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToPolicyCreateMap() (map[string]interface{}, error)
}

type Policy struct {
	Version   string      `json:"Version" required:"true"`
	Statement []Statement `json:"Statement" required:"true"`
}

type Statement struct {
	Action    []string                       `json:"Action" required:"true"`
	Effect    string                         `json:"Effect" required:"true"`
	Condition map[string]map[string][]string `json:"Condition,omitempty"`
	Resource  []string                       `json:"Resource,omitempty"`
}

// CreateOpts provides options used to create a policy.
type CreateOpts struct {
	Name        string `json:"display_name" required:"true"`
	Type        string `json:"type" required:"true"`
	Description string `json:"description" required:"true"`
	Policy      Policy `json:"policy" required:"true"`
}

// ToPolicyCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToPolicyCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "role")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Create creates a new Policy.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{201},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}

// Update updates an existing Policy.
func Update(client *golangsdk.ServiceClient, roleID string, opts CreateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(resourceURL(client, roleID), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}

// Delete deletes a policy.
func Delete(client *golangsdk.ServiceClient, roleID string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, roleID), &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}
