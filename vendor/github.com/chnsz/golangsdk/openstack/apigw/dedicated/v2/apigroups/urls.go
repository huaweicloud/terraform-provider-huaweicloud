package apigroups

import "github.com/chnsz/golangsdk"

const rootPath = "instances"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "api-groups")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, groupId string) string {
	return c.ServiceURL(rootPath, instanceId, "api-groups", groupId)
}

func associateDomainURL(c *golangsdk.ServiceClient, instanceId string, groupId string) string {
	return c.ServiceURL(rootPath, instanceId, "api-groups", groupId, "domains")
}
func disAssociateDomainURL(c *golangsdk.ServiceClient, instanceId string, groupId string, domainId string) string {
	return c.ServiceURL(rootPath, instanceId, "api-groups", groupId, "domains", domainId)
}

func domainAccessEnabledURL(c *golangsdk.ServiceClient, instanceId, groupId string) string {
	return c.ServiceURL(rootPath, instanceId, "api-groups", groupId, "sl-domain-access-settings")
}
