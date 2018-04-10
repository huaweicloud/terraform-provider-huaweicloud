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
