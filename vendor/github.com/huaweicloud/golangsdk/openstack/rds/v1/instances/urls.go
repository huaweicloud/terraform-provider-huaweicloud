package instances

import "github.com/huaweicloud/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "action")
}

func updatePolicyURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "backups/policy")
}
