package job

import "github.com/chnsz/golangsdk"

// POST /v1.1/{project_id}/clusters/{cluster_id}/cdm/job
func createURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "job")
}

// DELETE /v1.1/{project_id}/clusters/{cluster_id}/cdm/job/{job_name}
func deleteURL(c *golangsdk.ServiceClient, clusterId string, jobName string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "job", jobName)
}

// PUT /v1.1/{project_id}/clusters/{cluster_id}/cdm/job/{job_name}
func updateURL(c *golangsdk.ServiceClient, clusterId string, jobName string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "job", jobName)
}

// GET /v1.1/{project_id}/clusters/{cluster_id}/cdm/job/{job_name}
func getURL(c *golangsdk.ServiceClient, clusterId string, jobName string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "job", jobName)
}

// PUT /v1.1/{project_id}/clusters/{cluster_id}/cdm/job/{job_name}/start
func startURL(c *golangsdk.ServiceClient, clusterId string, jobName string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "job", jobName, "start")
}

// PUT /v1.1/{project_id}/clusters/{cluster_id}/cdm/job/{job_name}/stop
func stopURL(c *golangsdk.ServiceClient, clusterId string, jobName string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "job", jobName, "stop")
}

// GET /v1.1/{project_id}/clusters/{cluster_id}/cdm/job/{job_name}/status
func getStatusURL(c *golangsdk.ServiceClient, clusterId string, jobName string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "job", jobName, "status")
}

// GET /v1.1/{project_id}/clusters/{cluster_id}/cdm/submissions
func ListJobSubmissionsURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "submissions")
}
