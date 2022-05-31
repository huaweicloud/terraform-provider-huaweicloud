package jobs

import (
	"github.com/chnsz/golangsdk"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Get is a method to obtain the details of a specified deployment job using its ID.
func Get(c *golangsdk.ServiceClient, jobId string) (*JobResp, error) {
	var r JobResp
	_, err := c.Get(rootURL(c, jobId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}
