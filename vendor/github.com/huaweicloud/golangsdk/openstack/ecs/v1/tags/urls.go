package tags

import "github.com/huaweicloud/golangsdk"

const (
	rootPath     = "servers"
	resourcePath = "tags"
	actionPath   = "tags/action"
)

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, actionPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, resourcePath)
}
