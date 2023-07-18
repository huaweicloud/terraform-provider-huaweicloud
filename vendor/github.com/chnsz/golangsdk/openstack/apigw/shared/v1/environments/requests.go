package environments

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type ListOpts struct {
	// Environment name.
	EnvName string `q:"name"`
	// Number of records displayed on each page. The default value is 20.
	PageSize int `q:"page_size"`
	// Page number. The default value is 1.
	PageNum int `q:"page_no"`
	// Parameter name for exact matching.
	PreciseSearch string `q:"precise_search"`
	// Limit number of records displayed on each page. The default value is 20.
	Limit int `q:"limit"`
	// Offset Page number. The default value is 1.
	Offset int `q:"offset"`
}

// List is a method to obtain an array of one or more environments according to the query parameters.
// Note: The list returned by the function only contains the environment of the first page. This is because the return
// body does not contain page number information, so the page number of the next page cannot be obtained.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Environment, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return EnvironmentPage{pagination.SinglePageBase(r)}
	}).AllPages()
	if err != nil {
		return nil, err
	}

	var s []Environment
	err = pages.(EnvironmentPage).Result.ExtractIntoSlicePtr(&s, "envs")
	return s, err
}

// CreateOpts allows to create a new APIG environment using given parameters
type CreateOpts struct {
	// Environment name, which can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name" required:"true"`
	// Description of the environment, which can contain a maximum of 255 characters,
	// and the angle brackets (< and >) are not allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description string `json:"remark,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a new APIG shared environment using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Environment, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Environment
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// UpdateOpts allows to update an existing APIG environment using given parameters
type UpdateOpts struct {
	// Name of the APIG shared environment.
	Name string `json:"name" required:"true"`
	// Description of the APIG shared environment.
	Description *string `json:"remark,omitempty"`
}

// Update is a method to update a APIG shared environment using given parameters.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (*Environment, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Environment
	_, err = c.Put(environmentURL(c, id), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to delete an existing APIG shared environment.
func Delete(c *golangsdk.ServiceClient, id string) error {
	_, err := c.Delete(environmentURL(c, id), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
