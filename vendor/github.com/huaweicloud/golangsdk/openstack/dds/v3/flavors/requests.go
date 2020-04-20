package flavors

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOptsBuilder allows extensions to add parameters to the List request.
type ListOptsBuilder interface {
	ToFlavorListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListOpts struct {
	Region     string `q:"region"`
	EngineName string `q:"engine_name"`
}

// ToFlavorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List implements a flavor List request.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := baseURL(client)
	if opts != nil {
		query, err := opts.ToFlavorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
