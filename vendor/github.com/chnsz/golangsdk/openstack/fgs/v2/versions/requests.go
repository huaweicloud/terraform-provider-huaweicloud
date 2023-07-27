package versions

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOpts is the structure that used to query function version list.
type ListOpts struct {
	// Function URN.
	FunctionUrn string `json:"-" required:"true"`
	// The current query index.
	Marker int `q:"marker"`
	// Maximum number of functions to obtain in a request.
	MaxItems int `q:"maxitems"`
}

// List is a method to query the list of the function versions using given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Version, error) {
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url := rootURL(client, opts.FunctionUrn) + query.String()
	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := VersionPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	pageInfo, err := extractPageInfo(pages)
	if err != nil {
		return nil, err
	}
	return pageInfo.Versions, nil
}
