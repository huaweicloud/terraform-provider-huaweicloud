package roles

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOpts is the structure required by the Create method to create a new database role.
type CreateOpts struct {
	// Role name.
	// The length is 1~64 bits and can contain letters, numbers, hyphens, underscores and dots.
	Name string `json:"role_name" required:"true"`
	// The name of the database where the role is located.
	// The length is 1~64 bits and can contain letters, numbers and underscores.
	DbName string `json:"db_name,omitempty"`
	// List of roles inherited by the newly created role.
	Roles []Role `json:"roles,omitempty"`
}

// Role is the object that represent the role details.
type Role struct {
	// Role name.
	// The length is 1~64 bits and can contain letters, numbers, hyphens, underscores and dots.
	Name string `json:"role_name" required:"true"`
	// The name of the database where the role is located.
	// The length is 1~64 bits and can contain letters, numbers and underscores.
	DbName string `json:"role_db_name" required:"true"`
}

// Create is a method to create a new database role using given parameters.
func Create(c *golangsdk.ServiceClient, instanceId string, opts CreateOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Post(rootURL(c, instanceId), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Role name.
	// The length is 1~64 bits and can contain letters, numbers, hyphens, underscores and dots.
	Name string `q:"role_name"`
	// The name of the database where the role is located.
	// The length is 1~64 bits and can contain letters, numbers and underscores.
	DbName string `q:"db_name"`
	// The offset number.
	// Default value: 0.
	Offset int `q:"offset"`
	// Number of records to be queried.
	// Value range: 0–100.
	// Default value: 100, indicating that a maximum of 1000 records can be queried and all records are displayed on
	// the same page.
	Limit int `q:"limit"`
}

// List is a method to query the list of the database roles using given opts.
func List(c *golangsdk.ServiceClient, instanceId string, opts ListOpts) ([]RoleResp, error) {
	url := resourceURL(c, instanceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var result []RoleResp
	err = pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := RolePage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).EachPage(func(page pagination.Page) (bool, error) {
		resp, err := ExtractRoles(page)
		if err != nil {
			return false, err
		}
		if len(resp) == 0 {
			return false, nil
		}
		result = append(result, resp...)
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, err
}

// DeleteOpts is the structure required by the Delete method to remove an existing database role.
type DeleteOpts struct {
	// Role name.
	// The length is 1~64 bits and can contain letters, numbers, hyphens, underscores and dots.
	Name string `json:"role_name" required:"true"`
	// The name of the database where the role is located.
	// The length is 1~64 bits and can contain letters, numbers and underscores.
	DbName string `json:"db_name" required:"true"`
}

// Delete is a method to remove an existing database role.
func Delete(c *golangsdk.ServiceClient, instanceId string, opts DeleteOpts) error {
	url := rootURL(c, instanceId)
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.DeleteWithBody(url, b, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
