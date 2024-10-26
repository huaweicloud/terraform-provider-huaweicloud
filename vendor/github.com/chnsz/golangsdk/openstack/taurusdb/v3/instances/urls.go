package instances

import "github.com/chnsz/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func deleteURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID)
}

func getURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func updateURL(c *golangsdk.ServiceClient, instanceID string, update string) string {
	return c.ServiceURL("instances", instanceID, update)
}

func deleteReplicaURL(c *golangsdk.ServiceClient, instanceID, nodeID string) string {
	return c.ServiceURL("instances", instanceID, "nodes", nodeID)
}

func proxyURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "proxy")
}

func secondLevelMonitoringURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "monitor-policy")
}

func versionURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "database-version")
}

func slowLogShowOriginalSwitchURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "slowlog/query")
}

func jobURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("jobs")
}

func listDehURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("dedicated-resources")
}
