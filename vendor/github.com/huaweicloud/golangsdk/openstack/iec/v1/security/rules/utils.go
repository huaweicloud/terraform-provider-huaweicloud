package rules

import (
	"github.com/huaweicloud/golangsdk"
)

func CreateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("security-group-rules")
}

func DeleteURL(c *golangsdk.ServiceClient, securityGroupRuleID string) string {
	return c.ServiceURL("security-group-rules", securityGroupRuleID)
}

func GetURL(c *golangsdk.ServiceClient, securityGroupRuleID string) string {
	return c.ServiceURL("security-group-rules", securityGroupRuleID)
}
