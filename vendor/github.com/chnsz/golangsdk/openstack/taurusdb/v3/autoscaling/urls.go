package autoscaling

import "github.com/chnsz/golangsdk"

func updateURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "auto-scaling/policy")
}

func getURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "auto-scaling/policy")
}
