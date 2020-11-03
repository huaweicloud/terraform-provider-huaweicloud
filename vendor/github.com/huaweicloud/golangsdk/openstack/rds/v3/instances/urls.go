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

func restartURL(c *golangsdk.ServiceClient, instancesId string) string {
	return c.ServiceURL("instances", instancesId, "action")
}

func singletohaURL(c *golangsdk.ServiceClient, instancesId string) string {
	return c.ServiceURL("instances", instancesId, "action")
}

func resizeURL(c *golangsdk.ServiceClient, instancesId string) string {
	return c.ServiceURL("instances", instancesId, "action")
}

func enlargeURL(c *golangsdk.ServiceClient, instancesId string) string {
	return c.ServiceURL("instances", instancesId, "action")
}

func listerrorlogURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "errorlog")
}

func listslowlogURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "slowlog")
}
