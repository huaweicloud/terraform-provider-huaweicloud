package jobs

import "github.com/chnsz/golangsdk"

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Get is a method to obtain the job detail of the ECS API request.
func Get(c *golangsdk.ServiceClient, jobId string) (*Job, error) {
	var rst golangsdk.Result
	_, err := c.Get(rootURL(c, jobId), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	if err == nil {
		var r Job
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}
