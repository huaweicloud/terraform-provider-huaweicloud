package dependencies

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Dependency type, which support public, private, and all, default to all.
	//   public
	//   private
	//   all
	DependencyType string `q:"dependency_type"`
	// Runtime of function.
	Runtime string `q:"runtime"`
	// Name of the dependency.
	Name string `q:"name"`
	// Final record queried last time. Default value: 0.
	Marker string `q:"marker"`
	// Maximum number of dependencies that can be obtained in a query, default to 400.
	Limit string `q:"limit"`
}

// ListOptsBuilder is an interface which to support request query build of
// the dependent package search.
type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

// ToListQuery is a method which to build a request query by the ListOpts.
func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List is a method to obtain an array of one or more dependent packages according to the query parameters.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := DependencyPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}
