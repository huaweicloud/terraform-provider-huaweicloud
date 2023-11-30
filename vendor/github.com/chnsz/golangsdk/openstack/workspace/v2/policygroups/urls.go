package policygroups

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("policy-groups")
}

func resourceURL(c *golangsdk.ServiceClient, groupId string) string {
	return c.ServiceURL("policy-groups", groupId)
}
