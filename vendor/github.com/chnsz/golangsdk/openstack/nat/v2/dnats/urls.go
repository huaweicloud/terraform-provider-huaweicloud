package dnats

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("dnat_rules")
}

func resourceURL(c *golangsdk.ServiceClient, ruleId string) string {
	return c.ServiceURL("dnat_rules", ruleId)
}

func deleteURL(c *golangsdk.ServiceClient, gatewayId, ruleId string) string {
	return c.ServiceURL("nat_gateways", gatewayId, "dnat_rules", ruleId)
}
