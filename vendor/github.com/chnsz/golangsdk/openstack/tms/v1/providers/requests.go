package providers

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOpts allows to filter supported provider list.
type ListOpts struct {
	// Specifies the display language, defaults to 'zh-cn'.
	Locale string `q:"locale"`
	// Number of records to be queried.
	// The valid value is range from 1 to 200, defaults to 200.
	Limit int `q:"limit"`
	// Specifies the index position, which starts from the next data record specified by offset.
	// The value must be a number and connot be a negative number, defaults to 0.
	Offset string `q:"offset"`
	// Specifies the cloud service name.
	Provider string `q:"provider"`
}

// List is a method to query the supported provider list using given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Provider, error) {
	url := queryURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := ProviderPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()
	if err != nil {
		return nil, err
	}
	return extractProviders(pages)
}
