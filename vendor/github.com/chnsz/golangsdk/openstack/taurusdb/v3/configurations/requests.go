package configurations

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json"},
}

func List(client *golangsdk.ServiceClient) (r ListResult) {
	_, r.Err = client.Get(listURL(client), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

type ApplyOpts struct {
	InstanceIds []string `json:"instance_ids" required:"true"`
}

func Apply(c *golangsdk.ServiceClient, configID string, opts ApplyOpts) (r JobResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Put(applyURL(c, configID), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}
