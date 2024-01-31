package policyexecutelogs

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOpts is the structure that used to query the execution logs of scaling policy.
type ListOpts struct {
	// The scaling policy ID.
	PolicyID string `json:"scaling_policy_id" required:"true"`
	// The policy execution log ID.
	LogID string `q:"log_id"`
	// The scaling resource type. The value can be SCALING_GROUP or BANDWIDTH.
	ResourceType string `q:"scaling_resource_type"`
	// The scaling resource ID.
	ResourceID string `q:"scaling_resource_id"`
	// The policy execution type. The value can be SCHEDULED, RECURRENCE, ALARM or MANUAL.
	ExecuteType string `q:"execute_type"`
	// The starting time of query.
	StartTime string `q:"start_time"`
	// The ending time of query.
	EndTime string `q:"end_time"`
	// Start number value. The value must be a positive integer.
	StartNumber int `q:"start_number"`
	// Number of records displayed per page.
	// The value must be a positive integer.
	Limit int `q:"limit"`
}

// List is a method used to query the execution logs of scaling policy with given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]ExecuteLog, error) {
	url := listURL(client, opts.PolicyID)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := ExecuteLogPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return extractExecuteLogs(pages)
}
