package rules

import "github.com/chnsz/golangsdk"

const rootPath = "aad/instances"

func batchCreateURL(c *golangsdk.ServiceClient, instanceId, ip string) string {
	return c.ServiceURL(rootPath, instanceId, ip, "rules/batch-create")
}

func listURL(c *golangsdk.ServiceClient, instanceId, ip string) string {
	return c.ServiceURL(rootPath, instanceId, ip, "rules")
}

func updateURL(c *golangsdk.ServiceClient, instanceId, ip, ruleId string) string {
	return c.ServiceURL(rootPath, instanceId, ip, "rules", ruleId)
}

func batchDeleteURL(c *golangsdk.ServiceClient, instanceId, ip string) string {
	return c.ServiceURL(rootPath, instanceId, ip, "rules/batch-delete")
}
