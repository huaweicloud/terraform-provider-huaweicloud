package link

import "github.com/chnsz/golangsdk"

// POST /v1.1/{project_id}/clusters/{cluster_id}/cdm/link
func createLinkURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "link")
}

// GET /v1.1/{project_id}/clusters/{cluster_id}/cdm/link/{link_name}
func showLinkURL(c *golangsdk.ServiceClient, clusterId string, linkName string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "link", linkName)
}

// PUT /v1.1/{project_id}/clusters/{cluster_id}/cdm/link/{link_name}
func updateLinkURL(c *golangsdk.ServiceClient, clusterId string, linkName string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "link", linkName)
}

// DELETE /v1.1/{project_id}/clusters/{cluster_id}/cdm/link/{link_name}
func deleteLinkURL(c *golangsdk.ServiceClient, clusterId string, linkName string) string {
	return c.ServiceURL("clusters", clusterId, "cdm", "link", linkName)
}
