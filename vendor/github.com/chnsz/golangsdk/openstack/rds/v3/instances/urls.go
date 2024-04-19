package instances

import "github.com/chnsz/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func deleteURL(c *golangsdk.ServiceClient, serverID string) string {
	return c.ServiceURL("instances", serverID)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func updateURL(c *golangsdk.ServiceClient, instanceId string, update string) string {
	return c.ServiceURL("instances", instanceId, update)
}

func jobURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs")
}

func getURL(c *golangsdk.ServiceClient, instanceId, getContent string) string {
	return c.ServiceURL("instances", instanceId, getContent)
}

func engineURL(c *golangsdk.ServiceClient, dbName string) string {
	return c.ServiceURL("datastores", dbName)
}

func applyConfigurationURL(c *golangsdk.ServiceClient, configId string) string {
	return c.ServiceURL("configurations", configId, "apply")
}
