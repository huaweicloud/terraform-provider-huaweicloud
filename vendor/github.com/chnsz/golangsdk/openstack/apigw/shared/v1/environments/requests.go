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
}

// List is a method to obtain an array of one or more environments according to the query parameters.
// Note: The list returned by the function only contains the environment of the first page. This is because the return
//       body does not contain page number information, so the page number of the next page cannot be obtained.
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
