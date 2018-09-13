package secgroups

import "github.com/huaweicloud/golangsdk"

const (
	secgrouppath = "os-security-groups"
	rulepath     = "os-security-group-rules"
)

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(secgrouppath, id)
}

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(secgrouppath)
}

func listByServerURL(c *golangsdk.ServiceClient, serverID string) string {
	return c.ServiceURL("servers", serverID, secgrouppath)
}

func rootRuleURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rulepath)
}

func resourceRuleURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rulepath, id)
}

func serverActionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("servers", id, "action")
}
