package tags

import (
	"github.com/huaweicloud/golangsdk"
)

func createURL(c *golangsdk.ServiceClient, resource_type, resource_id string) string {
	return c.ServiceURL("os-vendor-tags", resource_type, resource_id)
}

func getURL(c *golangsdk.ServiceClient, resource_type, resource_id string) string {
	return c.ServiceURL("os-vendor-tags", resource_type, resource_id)
}
