package flavors

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type DbFlavorsOpts struct {
	Versionname string `q:"version_name"`
}

type DbFlavorsBuilder interface {
	ToDbFlavorsListQuery() (string, error)
}

func (opts DbFlavorsOpts) ToDbFlavorsListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts DbFlavorsBuilder, databasename string) pagination.Pager {
	url := listURL(client, databasename)
	if opts != nil {
		query, err := opts.ToDbFlavorsListQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageRdsList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return DbFlavorsPage{pagination.SinglePageBase(r)}
	})

	rdsheader := map[string]string{"Content-Type": "application/json"}
	pageRdsList.Headers = rdsheader
	return pageRdsList
}
