package parameters

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type UpdateParametersOpts struct {
	ParameterValues map[string]string `json:"parameter_values" required:"true"`
}

func Update(c *golangsdk.ServiceClient, instanceId string, opts UpdateParametersOpts) (r JobResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Put(updateURL(c, instanceId), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}

func List(client *golangsdk.ServiceClient, instanceId string) ([]ParameterValue, error) {
	url := listURL(client, instanceId)
	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ParameterPage{pagination.OffsetPageBase{PageResult: r}}
	}).AllPages()
	if err != nil {
		return nil, err
	}
	res, err := ExtractParameters(pages)
	if err != nil {
		return nil, err
	}
	return res.ParameterValues, nil
}
