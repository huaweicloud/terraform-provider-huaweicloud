package instances

import "github.com/huaweicloud/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func deleteURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID)
}

func extendURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "extend-volume")
}

func enlargeNodeURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "enlarge-node")
}

func reduceNodeURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "reduce-node")
}

func updateNameURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "name")
}

func updatePassURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "password")
}

func resizeURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "resize")
}

func updateSgURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "security-group")
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}
