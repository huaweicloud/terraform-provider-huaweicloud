package instances

import "github.com/chnsz/golangsdk"

const resourcePath = "instances"

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "onekey-purchase")
}
