package sqljob

import "github.com/chnsz/golangsdk"

const (
	resourcePath = "jobs"
)

// POST /v1.0/{project_id}/jobs/submit-job
func submitURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "submit-job")
}

// DELETE /v1.0/{project_id}/jobs/{job_id}
func resourceURL(c *golangsdk.ServiceClient, jobId string) string {
	return c.ServiceURL(resourcePath, jobId)
}

// GET /v1.0/{project_id}/jobs
func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

// GET /v1.0/{project_id}/jobs/{job_id}/status
func queryStatusURL(c *golangsdk.ServiceClient, jobId string) string {
	return c.ServiceURL(resourcePath, jobId, "status")
}

// GET/v1.0/{project_id}/jobs/{job_id}/detail
func detailURL(c *golangsdk.ServiceClient, jobId string) string {
	return c.ServiceURL(resourcePath, jobId, "detail")
}

// POST /v1.0/{project_id}/jobs/check-sql
func checkSqlURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "check-sql")
}

// POST /v1.0/{project_id}/jobs/{job_id}/export-result
func exportResultURL(c *golangsdk.ServiceClient, jobId string) string {
	return c.ServiceURL(resourcePath, jobId, "export-result")
}

// GET /v1/{project_id}/jobs/{job_id}/progress
func progressURL(c *golangsdk.ServiceClient, jobId string) string {
	return c.ServiceURL(resourcePath, jobId, "progress")
}

// POST /v1.0/{project_id}/jobs/import-table
func importTableURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "import-table")
}

// POST /v1.0/{project_id}/jobs/export-table
func exportTableURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "export-table")
}
