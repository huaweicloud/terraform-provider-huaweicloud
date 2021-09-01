package blockchains

import "github.com/chnsz/golangsdk"

const resourcePath = "blockchains"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL(resourcePath, instanceID)
}

func extraURL(c *golangsdk.ServiceClient, instanceID, extra string) string {
	return c.ServiceURL(resourcePath, instanceID, extra)
}
