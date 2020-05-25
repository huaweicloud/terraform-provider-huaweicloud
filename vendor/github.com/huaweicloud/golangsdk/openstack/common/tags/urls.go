package tags

import (
	"github.com/huaweicloud/golangsdk"
)

// supported resourceType: "vpcs", "subnets", "publicips"
// "DNS-public_zone", "DNS-private_zone", "DNS-ptr_record"
// "DNS-public_recordset", "DNS-private_recordset"
func actionURL(c *golangsdk.ServiceClient, resourceType, id string) string {
	return c.ServiceURL(c.ProjectID, resourceType, id, "tags/action")
}

func getURL(c *golangsdk.ServiceClient, resourceType, id string) string {
	return c.ServiceURL(c.ProjectID, resourceType, id, "tags")
}

func listURL(c *golangsdk.ServiceClient, resourceType string) string {
	return c.ServiceURL(c.ProjectID, resourceType, "tags")
}
