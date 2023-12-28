package activitylogs

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOpts is the structure that used to query the activity logs of scaling group.
type ListOpts struct {
	// Starting time of query.
	StartTime string `q:"start_time"`
	// Ending time of query.
	EndTime string `q:"end_time"`
	// Number of records displayed per page.
	// The value must be a positive integer.
	Limit int `q:"limit"`
	// Start number value. The value must be a positive integer.
	StartNumber int `q:"start_number"`
}

// List is a method used to query the activity logs of scaling group with given parameters.
func List(client *golangsdk.ServiceClient, groupID string, opts ListOpts) ([]ActivityLog, error) {
	url := listURL(client, groupID)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := ActivityLogPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractActivityLogs(pages)
}
