package acls

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "acls")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, policyId string) string {
	return c.ServiceURL("instances", instanceId, "acls", policyId)
}
