package resources

import "github.com/chnsz/golangsdk"

// GET /v1/resource-manager/domains/{domain_id}/all-resources
func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("resource-manager", "domains", c.DomainID, "all-resources")
}
