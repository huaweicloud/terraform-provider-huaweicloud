package tags

import (
	"github.com/huaweicloud/golangsdk"
)

const tagPath = "scaling_group_tag"

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(tagPath, id, "tags/action")
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(tagPath, id, "tags")
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(tagPath, "tags")
}
