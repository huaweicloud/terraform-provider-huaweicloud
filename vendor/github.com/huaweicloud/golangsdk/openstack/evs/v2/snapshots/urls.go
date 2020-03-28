package snapshots

import "github.com/huaweicloud/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudsnapshots")
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudsnapshots/detail")
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("cloudsnapshots", id)
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return deleteURL(c, id)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return deleteURL(c, id)
}
