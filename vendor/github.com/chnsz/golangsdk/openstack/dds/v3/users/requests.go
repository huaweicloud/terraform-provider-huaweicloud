package users

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dds/v3/roles"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOpts is the structure required by the Create method to create a new database user.
type CreateOpts struct {
	// User name.
	// The length is 1~64 bits and can contain letters, numbers, hyphens, underscores and dots.
	Name string `json:"user_name" required:"true"`
	// Database user password.
	// The length is 8~32 digits, and must be uppercase letters (A~Z), lowercase letters (a~z), numbers (0~9), special
	// characters ~!@#%^*-_=+? The combination.
	// It is recommended that you enter a strong password to improve security and prevent security risks such as
	// password cracking by brute force.
	Password string `json:"user_pwd" required:"true"`
	// List of roles inherited by the newly created role.
	Roles []roles.Role `json:"roles" required:"true"`
	// The name of the database where the user is located.
	// The length is 1~64 bits and can contain letters, numbers and underscores.
	DbName string `json:"db_name,omitempty"`
}

// Create is a method to create a new database user using given parameters.
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
	// User name.
	// The length is 1~64 bits and can contain letters, numbers, hyphens, underscores and dots.
	Name string `q:"user_name"`
	// The name of the database where the user is located.
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

// List is a method to query the list of the users using given opts.
func List(c *golangsdk.ServiceClient, instanceId string, opts ListOpts) ([]UserResp, error) {
	url := resourceURL(c, instanceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var result []UserResp
	err = pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := UserPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).EachPage(func(page pagination.Page) (bool, error) {
		resp, err := ExtractUsers(page)
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

// PwdResetOpts is the structure required by the ResetPassword method to reset the database user password.
type PwdResetOpts struct {
	// New user password for reset.
	// The length is 8~32 digits, and must be uppercase letters (A~Z), lowercase letters (a~z), numbers (0~9), special
	// characters ~!@#%^*-_=+? The combination.
	// It is recommended that you enter a strong password to improve security and prevent security risks such as
	// password cracking by brute force.
	Password string `json:"user_pwd" required:"true"`
	// User name. Defaults to "rwuser".
	// The length is 1~64 bits and can contain letters, numbers, hyphens, underscores and dots.
	Name string `json:"user_name,omitempty"`
	// The name of the database where the user is located.
	// The length is 1~64 bits and can contain letters, numbers and underscores.
	DbName string `json:"db_name,omitempty"`
}

// ResetPassword is a method to reset the database user password.
func ResetPassword(c *golangsdk.ServiceClient, instanceId string, opts PwdResetOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Put(pwdResetURL(c, instanceId), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// DeleteOpts is the structure required by the Delete method to remove an existing database user.
type DeleteOpts struct {
	// User name.
	// The length is 1~64 bits and can contain letters, numbers, hyphens, underscores and dots.
	Name string `json:"user_name" required:"true"`
	// The name of the database where the user is located.
	// The length is 1~64 bits and can contain letters, numbers and underscores.
	DbName string `json:"db_name" required:"true"`
}

// Delete is a method to remove an existing database user.
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
