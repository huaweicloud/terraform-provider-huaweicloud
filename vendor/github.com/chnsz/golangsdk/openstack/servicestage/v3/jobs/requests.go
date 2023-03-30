package jobs

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Get is a method to obtain the details of a specified deployment job using its ID.
//
// Deprecated: Please use the List method to query task details because of the Get method can only obtain maximum of 20
// results of task (The first page).
func Get(c *golangsdk.ServiceClient, jobId string) (*JobResp, error) {
	var r JobResp
	_, err := c.Get(rootURL(c, jobId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Instance ID of the component.
	InstanceId string `q:"instance_id"`
	// Number of records to be queried.
	// Default value: 20
	Limit int `q:"limit"`
	// The offset number.
	Offset int `q:"offset"`
	// Descending or ascending order.
	Desc string `q:"desc"`
}

// List is a method to query the list of the deployment task using given ID and opts.
func List(c *golangsdk.ServiceClient, jobId string, opts ListOpts) ([]Task, error) {
	url := rootURL(c, jobId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := TaskPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractTasks(pages)
}
