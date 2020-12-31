package firewalls

import (
	"github.com/huaweicloud/golangsdk"
)

func CreateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("firewalls")
}

func DeleteURL(c *golangsdk.ServiceClient, firewallId string) string {
	return c.ServiceURL("firewalls", firewallId)
}

func GetURL(c *golangsdk.ServiceClient, firewallId string) string {
	return c.ServiceURL("firewalls", firewallId)
}

func UpdateURL(c *golangsdk.ServiceClient, firewallId string) string {
	return c.ServiceURL("firewalls", firewallId)
}

func UpdateRuleURL(c *golangsdk.ServiceClient, firewallId string) string {
	return c.ServiceURL("firewalls", firewallId, "update_firewall_rules")
}
