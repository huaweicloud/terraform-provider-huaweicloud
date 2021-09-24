package namespaces

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToNamespaceCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new network
type CreateOpts struct {
	Namespace string `json:"namespace" required:"true"`
}

// ToNamespaceCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToNamespaceCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new namespace.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToNamespaceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular network based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// Delete will permanently delete a particular network based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// CreateAccessOptsBuilder allows extensions to add additional parameters to the create request.
type CreateAccessOptsBuilder interface {
	ToAccessCreateMap() (map[string]interface{}, error)
}

// CreateAccessOpts contains all the values needed to create access of a namespace
type CreateAccessOpts struct {
	Users []User
}

// Access information of a user
type User struct {
	// ID of the user
	UserID string `json:"user_id" required:"true"`
	// Name of the user
	UserName string `json:"user_name" required:"true"`
	// Permission of the user, 7: Manage. 3: Write. 1: Read
	Auth int `json:"auth" required:"true"`
}

// ToAccessCreateMap builds a create request body from CreateAccessOpts.
func (opts CreateAccessOpts) ToAccessCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// CreateAccess accepts a CreateAccessOpts struct and uses the values to create access of a namespace.
func CreateAccess(c *golangsdk.ServiceClient, opts CreateAccessOptsBuilder, namespace string) (r CreateAccessResult) {
	b, err := opts.ToAccessCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(accessURL(c, namespace), b["Users"], &r.Body, nil)
	return
}

// Get retrieves the access of a namespace.
func GetAccess(c *golangsdk.ServiceClient, namespace string) (r GetAccessResult) {
	_, r.Err = c.Get(accessURL(c, namespace), &r.Body, nil)
	return
}

// Delete will permanently delete the access of a namespace.
func DeleteAccess(c *golangsdk.ServiceClient, userIDs []string, namespace string) (r DeleteAccessResult) {
	reqOpt := &golangsdk.RequestOpts{JSONBody: userIDs}
	_, r.Err = c.Delete(accessURL(c, namespace), reqOpt)
	return
}
