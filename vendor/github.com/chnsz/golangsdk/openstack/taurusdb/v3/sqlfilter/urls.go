package sqlfilter

import "github.com/chnsz/golangsdk"

func updateURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "sql-filter", "switch")
}

func getURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "sql-filter", "switch")
}
