package instances

import "github.com/huaweicloud/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func deleteURL(c *golangsdk.ServiceClient, serverID string) string {
	return c.ServiceURL("instances", serverID)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func modifyURL(c *golangsdk.ServiceClient, serverID, action string) string {
	return c.ServiceURL("instances", serverID, action)
}
