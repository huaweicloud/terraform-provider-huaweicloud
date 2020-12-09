package tags

import "github.com/huaweicloud/golangsdk"

const (
	rootPath     = "images"
	resourcePath = "tags"
	actionPath   = "tags/action"
)

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, rootPath, id, actionPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, rootPath, id, resourcePath)
}
