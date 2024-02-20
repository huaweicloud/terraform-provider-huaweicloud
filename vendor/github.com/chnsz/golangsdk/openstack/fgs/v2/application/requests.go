package application

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOpts is the structure that used to query application template list.
type ListOpts struct {
	// Template runtime.
	Runtime string `q:"runtime"`
	// Template category.
	Category string `q:"category"`
	// The current query index.
	Marker int `q:"marker"`
	// Maximum number of templates to obtain in a request.
	MaxItems int `q:"maxitems"`
}

// ListTemplates is a method to query the list of the application templates using given parameters.
func ListTemplates(client *golangsdk.ServiceClient, opts ListOpts) ([]Template, error) {
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url := listURL(client) + query.String()
	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := TemplatePage{pagination.MarkerPageBase{PageResult: r}}
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
	return pageInfo.Templates, nil
}
