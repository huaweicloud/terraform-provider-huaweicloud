package networks

import "github.com/huaweicloud/golangsdk"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("networks", id)
}

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("networks")
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *golangsdk.ServiceClient) string {
	return rootURL(c)
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return resourceURL(c, id)
}
