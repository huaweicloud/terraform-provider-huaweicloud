package servergroups

import "github.com/huaweicloud/golangsdk"

const resourcePath = "os-server-groups"

func resourceURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *golangsdk.ServiceClient) string {
	return resourceURL(c)
}

func createURL(c *golangsdk.ServiceClient) string {
	return resourceURL(c)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return getURL(c, id)
}
