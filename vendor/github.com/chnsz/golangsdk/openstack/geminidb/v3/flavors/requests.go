package flavors

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type ListFlavorOpts struct {
	EngineName string `q:"engine_name"`
}

type ListFlavorBuilder interface {
	ToFlavorListDetailQuery() (string, error)
}

func (opts ListFlavorOpts) ToFlavorListDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListFlavorBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToFlavorListDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.SinglePageBase(r)}
	})
	// Headers supplies additional HTTP headers to populate on each paged request
	pageList.Headers = map[string]string{"Content-Type": "application/json"}

	return pageList
}
