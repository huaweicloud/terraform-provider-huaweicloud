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

func updateURL(c *golangsdk.ServiceClient, instancesId string, updata string) string {
	return c.ServiceURL("instances", instancesId, updata)
}

func jobURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs")
}

func engineURL(c *golangsdk.ServiceClient, dbName string) string {
	return c.ServiceURL("datastores", dbName)
}
