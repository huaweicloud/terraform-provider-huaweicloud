package flavors

import (
	"github.com/chnsz/golangsdk"
)

// ListOpts allows the filtering flavors through the API.
type ListOpts struct {
	// Specifies the availability_zone, e.g. cn-north-1a
	AvailabilityZone string `q:"availability_zone"`
}

// List BMS flavors
func List(c *golangsdk.ServiceClient, opts ListOpts) (r ListResult) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		r.Err = err
		return
	}

	url := listURL(c) + q.String()
	_, r.Err = c.Get(url, &r.Body, nil)

	return
}
