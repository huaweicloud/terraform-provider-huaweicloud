package groups

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("groups")
}

func resourceURL(c *golangsdk.ServiceClient, groupId string) string {
	return c.ServiceURL("groups", groupId)
}

func userURL(c *golangsdk.ServiceClient, groupId string) string {
	return c.ServiceURL("groups", groupId, "users")
}

func actionURL(c *golangsdk.ServiceClient, groupId string) string {
	return c.ServiceURL("groups", groupId, "actions")
}
