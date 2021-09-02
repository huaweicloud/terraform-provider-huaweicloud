package apigroups

import "github.com/chnsz/golangsdk"

const rootPath = "instances"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "api-groups")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, groupId string) string {
	return c.ServiceURL(rootPath, instanceId, "api-groups", groupId)
}
