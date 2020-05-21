package tags

import (
	"github.com/huaweicloud/golangsdk"
)

// supported resourceType: "vpcs", "subnets", "publicips"
func actionURL(c *golangsdk.ServiceClient, resourceType, id string) string {
	return c.ServiceURL(c.ProjectID, resourceType, id, "tags/action")
}

func getURL(c *golangsdk.ServiceClient, resourceType, id string) string {
	return c.ServiceURL(c.ProjectID, resourceType, id, "tags")
}

func listURL(c *golangsdk.ServiceClient, resourceType string) string {
	return c.ServiceURL(c.ProjectID, resourceType, "tags")
}
