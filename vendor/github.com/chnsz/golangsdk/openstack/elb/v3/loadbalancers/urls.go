package loadbalancers

import "github.com/chnsz/golangsdk"

const (
	rootPath     = "elb"
	resourcePath = "loadbalancers"
	statusPath   = "statuses"
	forcePath    = "force-elb"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func statusRootURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, statusPath)
}

func resourceForceDeleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, forcePath)
}
