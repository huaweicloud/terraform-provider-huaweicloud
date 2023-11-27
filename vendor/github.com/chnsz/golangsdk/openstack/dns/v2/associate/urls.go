package associate

import "github.com/chnsz/golangsdk"

func associateURL(c *golangsdk.ServiceClient, resolverRuleID string) string {
	return c.ServiceURL("resolverrules", resolverRuleID, "associaterouter")
}

func disAssociateURL(c *golangsdk.ServiceClient, resolverRuleID string) string {
	return c.ServiceURL("resolverrules", resolverRuleID, "disassociaterouter")
}
