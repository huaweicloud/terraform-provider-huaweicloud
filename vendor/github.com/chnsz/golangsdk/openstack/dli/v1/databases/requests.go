package databases

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// CreateOpts is a structure which allows to create a new database using given parameters.
type CreateOpts struct {
	// Name of the created database.
	// NOTE: The default database is a built-in database. You cannot create a database named default.
	Name string `json:"database_name" required:"true"`
	// Information about the created database.
	Description string `json:"description,omitempty"`
	// Enterprise project ID. The value 0 indicates the default enterprise project.
	// NOTE: Users who have enabled Enterprise Management can set this parameter to bind a specified project.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// Database tag.
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

// Create is a method to create a new database by CreateOpts.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*RequestResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c), b, &rst.Body, nil)
	if err == nil {
		var r RequestResp
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// ListOpts is a structure which allows to obtain databases using given parameters.
type ListOpts struct {
	// Specifies whether to display the permission information. The value can be true or false.
	// The default value is false.
	IsDesplay bool `q:"with-priv"`
	// The value should be no less than 0. The default value is 0.
	Offset int `q:"offset"`
	// Number of returned data records. The value must be greater than or equal to 0.
	// By default, all data records are returned.
	Limit int `q:"limit"`
	// Database name filtering keyword. Fuzzy match is used to obtain all databases whose names contain the keyword.
	Keyword string `q:"keyword"`
	// Database tag. The format is key=value, for example:
	// GET /v1.0/{project_id}/databases?offset=0&limit=10&with-priv=true&tags=foo%3Dbar
	// In the preceding information, = needs to be escaped to %3D, foo indicates the tag key, and bar indicates the tag
	// value.
	Tags string `q:"tags"`
}

// List is a method to obtain a list of databases by ListOpts.
func List(c *golangsdk.ServiceClient, opts ListOpts) (*ListResp, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst golangsdk.Result
	_, err = c.Get(url, &rst.Body, nil)
	if err == nil {
		var r ListResp
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// UpdateDBOwnerOpts is a structure which allows to update database owner using given name.
type UpdateDBOwnerOpts struct {
	// Name of the new owner. The new user must be a sub-user of the current tenant.
	NewOwner string `json:"new_owner" required:"true"`
}

// UpdateDBOwner is a method to update database owner by UpdateDBOwnerOpts.
func UpdateDBOwner(c *golangsdk.ServiceClient, dbName string, opts UpdateDBOwnerOpts) (*RequestResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(userURL(c, dbName), b, &rst.Body, nil)
	if err == nil {
		var r RequestResp
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// Delete is a method to remove the exist database by database name.
func Delete(c *golangsdk.ServiceClient, dbName string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(resourceURL(c, dbName), nil)
	return &r
}
