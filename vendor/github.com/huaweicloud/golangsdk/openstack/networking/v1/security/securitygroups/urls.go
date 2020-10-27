package securitygroups

import (
	"github.com/huaweicloud/golangsdk"
)

func CreateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("security-groups")
}

func DeleteURL(c *golangsdk.ServiceClient, securityGroupId string) string {
	return c.ServiceURL("security-groups", securityGroupId)
}

func GetURL(c *golangsdk.ServiceClient, securityGroupId string) string {
	return c.ServiceURL("security-groups", securityGroupId)
}

func ListURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("security-groups")
}
