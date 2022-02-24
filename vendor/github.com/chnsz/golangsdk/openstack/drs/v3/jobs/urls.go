package jobs

import "github.com/chnsz/golangsdk"

// POST /v3/{project_id}/jobs/batch-creation
func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs", "batch-creation")
}

// POST /v3/{project_id}/jobs/batch-status
func statusURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs", "batch-status")
}

// POST /v3/{project_id}/jobs/batch-detail
func detailsURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs", "batch-detail")
}

// DELETE  /v3/{project_id}/jobs/batch-jobs
func deleteURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs", "batch-jobs")
}

// POST /v3/{project_id}/jobs/batch-starting
func startURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs", "batch-starting")
}

// POST /v3/{project_id}/jobs/batch-connection
func testConnectionsURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs", "batch-connection")
}

// PUT /v3/{project_id}/jobs/batch-modification
func updateJobURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs", "batch-modification")
}

// PUT /v3/{project_id}/jobs/batch-limit-speed
func limitSpeedURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs", "batch-limit-speed")
}

// POST /v3/{project_id}/jobs/batch-precheck
func preCheckURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs", "batch-precheck")
}

// POST /v3/{project_id}/jobs/batch-precheck-result
func batchCheckResultsURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs", "batch-precheck-result")
}

// POST /v3/{project_id}/jobs
func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs")
}
