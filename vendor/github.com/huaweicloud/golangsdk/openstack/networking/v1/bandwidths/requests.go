package bandwidths

import (
	"github.com/huaweicloud/golangsdk"
)

//UpdateOptsBuilder is an interface by which can be able to build the request
//body
type UpdateOptsBuilder interface {
	ToBWUpdateMap() (map[string]interface{}, error)
}

//UpdateOpts is a struct which represents the request body of update method
type UpdateOpts struct {
	Size int    `json:"size,omitempty"`
	Name string `json:"name,omitempty"`
}

func (opts UpdateOpts) ToBWUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "bandwidth")
}

//Get is a method by which can get the detailed information of a bandwidth
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

//Update is a method which can be able to update the port of public ip
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToBWUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToBWListQuery() (string, error)
}

// ListOpts allows extensions to add additional parameters to the API.
type ListOpts struct {
	ShareType           string `q:"share_type"`
	EnterpriseProjectID string `q:"enterprise_project_id"`
}

// ToBWListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToBWListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List is a method by which can get the detailed information of all bandwidths
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) (r ListResult) {
	url := listURL(client)
	query, err := opts.ToBWListQuery()
	if err != nil {
		r.Err = err
		return
	}
	url += query

	_, r.Err = client.Get(url, &r.Body, nil)
	return
}
