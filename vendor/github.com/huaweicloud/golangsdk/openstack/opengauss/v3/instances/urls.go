package instances

import "github.com/huaweicloud/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func deleteURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID)
}

func updateURL(c *golangsdk.ServiceClient, instanceID, updata string) string {
	return c.ServiceURL("instances", instanceID, updata)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}
