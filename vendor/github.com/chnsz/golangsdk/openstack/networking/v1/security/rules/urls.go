package rules

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "security-group-rules")
}

func resourceURL(c *golangsdk.ServiceClient, ruleId string) string {
	return c.ServiceURL(c.ProjectID, "security-group-rules", ruleId)
}
