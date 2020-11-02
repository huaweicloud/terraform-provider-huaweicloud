package enterpriseprojects

import (
	"github.com/huaweicloud/golangsdk"
)

type ListOpts struct {
	Name   string `q:"name"`
	ID     string `q:"id"`
	Status int    `q:"status"`
}

func (opts ListOpts) ToEnterpriseProjectListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

type ListOptsBuilder interface {
	ToEnterpriseProjectListQuery() (string, error)
}

func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) (r ListResult) {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToEnterpriseProjectListQuery()
		if err != nil {
			r.Err = err
		}
		url += query
	}

	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
