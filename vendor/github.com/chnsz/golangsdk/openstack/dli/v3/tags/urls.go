package tags

import "github.com/chnsz/golangsdk"

func tagsURL(c *golangsdk.ServiceClient, id, resourceType, action string) string {
	return c.ServiceURL(resourceType, id, "tags", action)
}
