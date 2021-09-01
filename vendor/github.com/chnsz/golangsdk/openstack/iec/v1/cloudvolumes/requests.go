package cloudvolumes

import (
	"net/http"

	"github.com/chnsz/golangsdk"
)

func Get(client *golangsdk.ServiceClient, CloudVolumeID string) (r GetResult) {
	url := GetURL(client, CloudVolumeID)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

func ListVolumeType(client *golangsdk.ServiceClient) (r GetResult) {
	url := ListVolumeTypeURL(client)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

// ListOpts holds options for listing Voulmes. It is passed to the List
// function.
type ListOpts struct {
	Limit  int    `q:"limit"`
	Offset int    `q:"offset"`
	Name   string `q:"name"`
	Status string `q:"status"`
}

type ListOptsBuilder interface {
	ToListVolumesQuery() (string, error)
}

func (opts ListOpts) ToListVolumesQuery() (string, error) {
	b, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) (r ListResult) {
	listURL := rootURL(client)
	if opts != nil {
		query, err := opts.ToListVolumesQuery()
		if err != nil {
			r.Err = err
			return r
		}
		listURL += query
	}

	_, r.Err = client.Get(listURL, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}
