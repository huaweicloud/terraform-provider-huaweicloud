package flavors

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToFlavorListMap() (string, error)
}

//ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	// Specifies the AZ name.
	AvailabilityZone string `q:"availability_zone"`
}

// ToFlavorListMap formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListMap() (string, error) {
	s, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return s.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// flavors.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		queryString, err := opts.ToFlavorListMap()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += queryString
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FlavorsPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
