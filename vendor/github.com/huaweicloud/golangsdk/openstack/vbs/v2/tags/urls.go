package tags

import "github.com/huaweicloud/golangsdk"

const rootPath = "backuppolicy"

func commonURL(c *golangsdk.ServiceClient, policyID string) string {
	return c.ServiceURL(rootPath, policyID, "tags")
}

func deleteURL(c *golangsdk.ServiceClient, policyID string, key string) string {
	return c.ServiceURL(rootPath, policyID, "tags", key)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, "resource_instances", "action")
}

func actionURL(c *golangsdk.ServiceClient, policyID string) string {
	return c.ServiceURL(rootPath, policyID, "tags", "action")
}
