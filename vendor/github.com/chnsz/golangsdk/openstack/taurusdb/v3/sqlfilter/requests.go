package sqlfilter

import (
	"github.com/chnsz/golangsdk"
)

type UpdateBuilder interface {
	ToUpdateMap() (map[string]interface{}, error)
}

type UpdateSqlFilterOpts struct {
	SwitchStatus string `json:"switch_status" required:"true"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func (opts UpdateSqlFilterOpts) ToUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Update(c *golangsdk.ServiceClient, instanceId string, opts UpdateBuilder) (r JobResult) {
	b, err := opts.ToUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(updateURL(c, instanceId), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}

func Get(c *golangsdk.ServiceClient, instanceId string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, instanceId), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}
