package l7policies

import "github.com/huaweicloud/golangsdk"

const (
	rootPath     = "lbaas"
	resourcePath = "l7policies"
	rulePath     = "rules"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func ruleRootURL(c *golangsdk.ServiceClient, policyID string) string {
	return c.ServiceURL(rootPath, resourcePath, policyID, rulePath)
}

func ruleResourceURL(c *golangsdk.ServiceClient, policyID string, ruleID string) string {
	return c.ServiceURL(rootPath, resourcePath, policyID, rulePath, ruleID)
}
