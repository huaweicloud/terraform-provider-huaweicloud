package orders

import (
	"github.com/huaweicloud/golangsdk"
)

type UnsubscribeOpts struct {
	ResourceIds     []string `json:"resource_ids" required:"true"`
	UnsubscribeType int      `json:"unsubscribe_type" required:"true"`
}

type UnsubscribeOptsBuilder interface {
	ToOrderUnsubscribeMap() (map[string]interface{}, error)
}

func (opts UnsubscribeOpts) ToOrderUnsubscribeMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Unsubscribe(client *golangsdk.ServiceClient, opts UnsubscribeOptsBuilder) (r UnsubscribeResult) {
	reqBody, err := opts.ToOrderUnsubscribeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(unsubscribeURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}
