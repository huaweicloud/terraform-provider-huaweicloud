package parameters

import "github.com/chnsz/golangsdk"

func updateURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "configurations")
}

func listURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "configurations")
}
