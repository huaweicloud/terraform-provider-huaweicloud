package environments

import "github.com/chnsz/golangsdk"

const rootPath = "instances"

func rootURL(c *golangsdk.ServiceClient, instanceId, path string) string {
	return c.ServiceURL(rootPath, instanceId, path)
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, path, id string) string {
	return c.ServiceURL(rootPath, instanceId, path, id)
}
