package groups

import "github.com/huaweicloud/golangsdk"

const rootPath = "security-groups"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id)
}
