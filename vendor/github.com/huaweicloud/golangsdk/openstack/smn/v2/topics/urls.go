package topics

import "github.com/huaweicloud/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("topics")
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("topics", id)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("topics", id)
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("topics", id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("topics?offset=0&limit=100")
}
