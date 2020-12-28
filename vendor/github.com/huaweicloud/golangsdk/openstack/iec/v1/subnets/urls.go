package subnets

import "github.com/huaweicloud/golangsdk"

func CreateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("subnets")
}

func DeleteURL(c *golangsdk.ServiceClient, subnetId string) string {
	return c.ServiceURL("subnets", subnetId)
}

func GetURL(c *golangsdk.ServiceClient, subnetId string) string {
	return c.ServiceURL("subnets", subnetId)
}

func UpdateURL(c *golangsdk.ServiceClient, subnetId string) string {
	return c.ServiceURL("subnets", subnetId)
}
