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

func updateURL(c *golangsdk.ServiceClient, instancesId string, update string) string {
	return c.ServiceURL("instances", instancesId, update)
}

func jobURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs")
}

func engineURL(c *golangsdk.ServiceClient, dbName string) string {
	return c.ServiceURL("datastores", dbName)
}

func resetRootPasswordURL(c *golangsdk.ServiceClient, instancesId string) string {
	return c.ServiceURL("instances", instancesId, "password")
}

func applyConfigurationURL(c *golangsdk.ServiceClient, configId string) string {
	return c.ServiceURL("configurations", configId, "apply")
}

func configurationsURL(c *golangsdk.ServiceClient, instancesId string) string {
	return c.ServiceURL("instances", instancesId, "configurations")
}

func actionURL(c *golangsdk.ServiceClient, instancesId string) string {
	return c.ServiceURL("instances", instancesId, "action")
}

func autoExpandURL(c *golangsdk.ServiceClient, instancesId string) string {
	return c.ServiceURL("instances", instancesId, "disk-auto-expansion")
}

func binlogRetentionHoursURL(c *golangsdk.ServiceClient, instancesId string) string {
	return c.ServiceURL("instances", instancesId, "binlog/clear-policy")
}
