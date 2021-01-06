package ports

import (
	"github.com/huaweicloud/golangsdk"
)

func CreateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("ports")
}

func DeleteURL(c *golangsdk.ServiceClient, portId string) string {
	return c.ServiceURL("ports", portId)
}

func GetURL(c *golangsdk.ServiceClient, portId string) string {
	return c.ServiceURL("ports", portId)
}

func UpdateURL(c *golangsdk.ServiceClient, portId string) string {
	return c.ServiceURL("ports", portId)
}
