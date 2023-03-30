package applications

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure required by the Create method to create a new application.
type CreateOpts struct {
	// Specified the application name with 2 to 64 characters long.
	// It consists of English letters, numbers, underscores (-), and underscores (_).
	// It must start with an English letter and end with an English letter or number.
	Name string `json:"name" required:"true"`
	// Specified the application description.
	// The description can contain a maximum of 96 characters.
	Description *string `json:"description,omitempty"`
	// Specified the enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// Application labels
	Labels []string `json:"labels"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a new ServiceStage application using create option.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Application, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst Application
	_, err = c.Post(rootURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &rst, err
}

// Get is a method to obtain the details of a specified ServiceStage application using its ID.
func Get(c *golangsdk.ServiceClient, appId string) (*Application, error) {
	var r Application
	_, err := c.Get(resourceURL(c, appId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// Value range: 0–100.
	// Default value: 1000, indicating that a maximum of 1000 records can be queried and all records are displayed on
	// the same page.
	Limit int `q:"limit"`
	// The offset number.
	Offset int `q:"offset"`
	// Sorting field. By default, query results are sorted by creation time.
	// The following enumerated values are supported: create_time, name, and update_time.
	OrderBy string `q:"order_by"`
	// Descending or ascending order. Default value: desc.
	Order string `q:"order"`
}

// List is a method to query the list of the environments using given opts.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Application, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := ApplicationPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractApplications(pages)
}

// Delete is a method to remove an existing application.
func Delete(c *golangsdk.ServiceClient, appId string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(resourceURL(c, appId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r
}
