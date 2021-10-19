package resources

import "github.com/chnsz/golangsdk"

const rootPath = "resources"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, subPath string) string {
	return c.ServiceURL(rootPath, subPath)
}
