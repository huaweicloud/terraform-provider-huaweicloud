package servergroups

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// List returns a Pager that allows you to iterate over a collection of
// ServerGroups.
func List(client *golangsdk.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, rootURL(client), func(r pagination.PageResult) pagination.Page {
		return ServerGroupPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToServerGroupCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies Server Group creation parameters.
type CreateOpts struct {
	// Name is the name of the server group
	Name string `json:"name" required:"true"`

	// Policies are the server group policies
	Policies []string `json:"policies" required:"true"`
}

// ToServerGroupCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToServerGroupCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "server_group")
}

// Create requests the creation of a new Server Group.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToServerGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get returns data about a previously created ServerGroup.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

// Delete requests the deletion of a previously allocated ServerGroup.
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return
}

type MemberOptsBuilder interface {
	ToServerGroupUpdateMemberMap(string) (map[string]interface{}, error)
}

type MemberOpts struct {
	InstanceID string `json:"instance_uuid" required:"true"`
}

func (opts MemberOpts) ToServerGroupUpdateMemberMap(optsType string) (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, optsType)
}

// UpdateMember is used to add and delete members from Server Group.
// (params)optsType: The opts type is title of block in request body.
//                   Add options is "add_memebr" and remove options is "remove_member".
func UpdateMember(client *golangsdk.ServiceClient, opts MemberOptsBuilder, optsType, id string) (r MemberResult) {
	b, err := opts.ToServerGroupUpdateMemberMap(optsType)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}
