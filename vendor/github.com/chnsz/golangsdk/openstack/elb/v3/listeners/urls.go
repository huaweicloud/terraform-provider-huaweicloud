package listeners

import "github.com/chnsz/golangsdk"

const (
	rootPath     = "elb"
	resourcePath = "listeners"
	forcePath    = "force"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func resourceForceDeleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, forcePath)
}
