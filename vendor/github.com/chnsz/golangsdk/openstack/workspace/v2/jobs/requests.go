package jobs

import "github.com/chnsz/golangsdk"

// ListOpts is the structure required by the List method to query job list.
type ListOpts struct {
	// Job status.
	Status string `q:"status"`
	// Job ID.
	JobId string `q:"job_id"`
	//Job type.
	JobType string `q:"job_type"`
	// Number of records to be queried.
	// Value range: 0–1000.
	// Default value: 1000, indicating that a maximum of 1000 records can be queried and all records are displayed on
	// the same page.
	Limit int `q:"limit"`
	// The offset number, start with 0.
	Offset int `q:"offset"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// List is a method to query the job details using given parameters.
func List(c *golangsdk.ServiceClient, opts ListOpts) (*QueryResp, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r QueryResp
	_, err = c.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}
