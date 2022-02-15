package resources

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type ListOpts struct {
	Region              string `q:"region_id"`
	EnterpriseProjectId string `q:"ep_id"`
	// provider.type
	Type string `q:"type"`
	// min：1
	// max：200
	Limit int `q:"limit"`
	// min：4
	// max：400
	Marker string `q:"marker"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func List(c *golangsdk.ServiceClient, opts ListOpts) pagination.Pager {
	url := listURL(c)

	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	url += query.String()

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := ResourcePage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}
