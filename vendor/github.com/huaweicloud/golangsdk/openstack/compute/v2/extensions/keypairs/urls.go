package keypairs

import "github.com/huaweicloud/golangsdk"

const resourcePath = "os-keypairs"

func resourceURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *golangsdk.ServiceClient) string {
	return resourceURL(c)
}

func createURL(c *golangsdk.ServiceClient) string {
	return resourceURL(c)
}

func getURL(c *golangsdk.ServiceClient, name string) string {
	return c.ServiceURL(resourcePath, name)
}

func deleteURL(c *golangsdk.ServiceClient, name string) string {
	return getURL(c, name)
}
