package sites

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

//ListSiteOptsBuilder list sites builder
type ListSiteOptsBuilder interface {
	ToListSiteQuery() (string, error)
}

// ListSiteOpts to list site
type ListSiteOpts struct {
	//Limit query limit
	Limit string `q:"limit"`

	//Offset query begin index
	Offset string `q:"offset"`

	//id query by id
	ID string `q:"id"`

	//Area query by area
	Area string `q:"area"`

	//Province query by province
	Province string `q:"province"`

	//City query by city
	City string `q:"city"`

	//Operator query by operator
	Operator string `q:"operator"`
}

// ToListSiteQuery converts ListSiteOpts structures to query string
func (opts ListSiteOpts) ToListSiteQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// sites. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(client *golangsdk.ServiceClient, opts ListSiteOptsBuilder) pagination.Pager {
	url := ListURL(client)
	if opts != nil {
		query, err := opts.ToListSiteQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SitePage{pagination.LinkedPageBase{PageResult: r}}
	})
}
