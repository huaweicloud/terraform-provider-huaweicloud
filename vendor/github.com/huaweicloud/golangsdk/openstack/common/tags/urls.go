package tags

import (
	"strings"

	"github.com/huaweicloud/golangsdk"
)

// supported resourceType: "vpcs", "subnets", "publicips"
// "DNS-public_zone", "DNS-private_zone", "DNS-ptr_record"
// "DNS-public_recordset", "DNS-private_recordset"
func actionURL(c *golangsdk.ServiceClient, resourceType, id string) string {
	if hasProjectID(c) {
		return c.ServiceURL(resourceType, id, "tags/action")
	}
	return c.ServiceURL(c.ProjectID, resourceType, id, "tags/action")
}

func getURL(c *golangsdk.ServiceClient, resourceType, id string) string {
	if hasProjectID(c) {
		return c.ServiceURL(resourceType, id, "tags")
	}
	return c.ServiceURL(c.ProjectID, resourceType, id, "tags")
}

func listURL(c *golangsdk.ServiceClient, resourceType string) string {
	if hasProjectID(c) {
		return c.ServiceURL(resourceType, "tags")
	}
	return c.ServiceURL(c.ProjectID, resourceType, "tags")
}

func hasProjectID(c *golangsdk.ServiceClient) bool {
	url := c.ResourceBaseURL()
	array := strings.Split(url, "/")

	// the baseURL must be end with "/"
	if array[len(array)-2] == c.ProjectID {
		return true
	}
	return false
}
