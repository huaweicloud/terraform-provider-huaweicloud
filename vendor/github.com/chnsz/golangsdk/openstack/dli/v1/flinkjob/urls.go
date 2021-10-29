package flinkjob

import (
	"strconv"

	"github.com/chnsz/golangsdk"
)

const (
	rootPath   = "streaming"
	sqlJobPath = "sql-jobs"
	jobsPath   = "jobs"
)

// POST /v1.0/{project_id}/streaming/sql-jobs
func createFlinkSqlUrl(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, sqlJobPath)
}

// PUT /v1.0/{project_id}/streaming/sql-jobs/{job_id}
func updateFlinkSqlURL(c *golangsdk.ServiceClient, jobId int) string {
	return c.ServiceURL(rootPath, sqlJobPath, strconv.Itoa(jobId))
}

// POST /v1.0/{project_id}/streaming/jobs/run
func runFlinkJobURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, jobsPath, "run")
}

// GET /v1.0/{project_id}/streaming/jobs/{job_id}
func getURL(c *golangsdk.ServiceClient, jobId int) string {
	return c.ServiceURL(rootPath, jobsPath, strconv.Itoa(jobId))
}

// GET /v1.0/{project_id}/streaming/jobs
func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, jobsPath)
}

// DELETE /v1.0/{project_id}/streaming/jobs/{job_id}
func deleteURL(c *golangsdk.ServiceClient, jobId int) string {
	return c.ServiceURL(rootPath, jobsPath, strconv.Itoa(jobId))
}

// POST /v1.0/{project_id}/dli/obs-authorize
func authorizeBucketURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("dli", "obs-authorize")
}

// POST /v1.0/{project_id}/streaming/flink-jobs
func createJarJobURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, "flink-jobs")
}

// PUT /v1.0/{project_id}/streaming/flink-jobs/{job_id}
func updateJarJobURL(c *golangsdk.ServiceClient, jobId int) string {
	return c.ServiceURL(rootPath, "flink-jobs", strconv.Itoa(jobId))
}

// POST /v1.0/{project_id}/streaming/jobs/stop
func stopJobURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, jobsPath, "stop")
}
