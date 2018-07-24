package peerings

import "github.com/huaweicloud/golangsdk"

const (
	resourcePath = "peerings"
	rootpath     = "vpc"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootpath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootpath, resourcePath, id)
}

func acceptURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootpath, resourcePath, id, "accept")
}

func rejectURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootpath, resourcePath, id, "reject")
}
