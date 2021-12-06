package bandwidths

import (
	"net/http"

	"github.com/chnsz/golangsdk"
)

func Get(client *golangsdk.ServiceClient, bandwidthId string) (r GetResult) {
	url := GetURL(client, bandwidthId)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

type UpdateOpts struct {
	// Specifies the bandwidth name. The value is a string of 1 to 64
	// characters that can contain letters, digits, underscores (_), and hyphens (-).
	Name string `json:"name,omitempty"`

	// Specifies the bandwidth size. The value ranges from 1 Mbit/s to
	// 300 Mbit/s.
	Size int `json:"size,omitempty"`
}

type UpdateOptsBuilder interface {
	ToBandwidthsUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToBandwidthsUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "bandwidth")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Update(client *golangsdk.ServiceClient, bandwidthId string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToBandwidthsUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(UpdateURL(client, bandwidthId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

type ListOpts struct {
	Limit  int    `q:"limit"`
	Offset int    `q:"offset"`
	SiteID string `q:"site_id"`
}

type ListBandwidthsOptsBuilder interface {
	ToListBandwidthsQuery() (string, error)
}

func (opts ListOpts) ToListBandwidthsQuery() (string, error) {
	b, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func List(client *golangsdk.ServiceClient, opts ListBandwidthsOptsBuilder) (r ListResult) {
	listURL := listURL(client)
	if opts != nil {
		query, err := opts.ToListBandwidthsQuery()
		if err != nil {
			r.Err = err
			return r
		}
		listURL += query
	}

	_, r.Err = client.Get(listURL, &r.Body, nil)
	return
}
