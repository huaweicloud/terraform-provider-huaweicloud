package flavors

import (
	"github.com/chnsz/golangsdk"
)

// ListOptsBuilder allows extensions to add parameters to the List request.
type ListOptsBuilder interface {
	ToFlavorListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned.
type ListOpts struct {
	CacheMode     string `q:"cache_mode"`
	Engine        string `q:"engine"`
	EngineVersion string `q:"engine_version"`
	Capacity      string `q:"capacity"`
	SpecCode      string `q:"spec_code"`
	CPUType       string `q:"cpu_type"`
}

// ToFlavorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List implements a flavor List request.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) (r ListResult) {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToFlavorListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}
