package assets

import (
	"github.com/chnsz/golangsdk"
)

type AssetsOpts struct {
	AssetId      string `q:"asset_id"`
	DeployedType string `q:"deployed_type"`
	Region       string `q:"region"`
	AssetVersion string `q:"asset_version"`
	QueryAll     bool   `q:"query_all"`
}

type AssetsBuilder interface {
	ToAssetsListQuery() (string, error)
}

func (opts AssetsOpts) ToAssetsListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts AssetsBuilder) (r ListAssetsResult) {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToAssetsListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
