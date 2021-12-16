package clusters

import "github.com/chnsz/golangsdk"

// POST /v1.1/{project_id}/clusters
func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("clusters")
}

// DELETE /v1.1/{project_id}/clusters/{cluster_id}
func deleteURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId)
}

// GET /v1.1/{project_id}/clusters
func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("clusters")
}

// POST /v1.1/{project_id}/clusters/{cluster_id}/action
func restartURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId, "action")
}

// GET /v1.1/{project_id}/clusters/{cluster_id}
func getURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId)
}

// POST /v1.1/{project_id}/clusters/{cluster_id}/action
func actionURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId, "action")
}
