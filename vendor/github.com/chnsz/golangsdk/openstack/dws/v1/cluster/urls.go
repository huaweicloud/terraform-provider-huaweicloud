package cluster

import "github.com/chnsz/golangsdk"

const (
	resourcePath = "clusters"
)

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

// resetPasswordURL /v1.0/{project_id}/clusters/{cluster_id}/reset-password
func resetPasswordURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL(resourcePath, clusterId, "reset-password")
}

// resizeURL /v1.0/{project_id}/clusters/{cluster_id}/resize
func resizeURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL(resourcePath, clusterId, "resize")
}
