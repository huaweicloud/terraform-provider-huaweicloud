package jobs

import (
	"github.com/chnsz/golangsdk"
)

var requestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// QueryOpts allows to filter list data using given parameters.
type QueryOpts struct {
	// Job ID.
	JobId string `q:"id"`
}

// Get is a method to retrieves a particular job based on its unique ID.
func Get(c *golangsdk.ServiceClient, jobId string) (*Job, error) {
	opts := QueryOpts{
		JobId: jobId,
	}
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r struct {
		Job Job `json:"job"`
	}
	_, err = c.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Job, err
}
