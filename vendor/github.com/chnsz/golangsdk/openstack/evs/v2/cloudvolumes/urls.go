package cloudvolumes

import "github.com/chnsz/golangsdk"

const resourcePath = "cloudvolumes"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "action")
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "detail")
}
