package snats

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("snat_rules")
}

func resourceURL(c *golangsdk.ServiceClient, ruleId string) string {
	return c.ServiceURL("snat_rules", ruleId)
}

func deleteURL(c *golangsdk.ServiceClient, gatewayId, ruleId string) string {
	return c.ServiceURL("nat_gateways", gatewayId, "snat_rules", ruleId)
}
