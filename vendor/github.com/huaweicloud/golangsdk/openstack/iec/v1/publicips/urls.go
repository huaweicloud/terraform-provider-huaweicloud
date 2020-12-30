package publicips

import (
	"github.com/huaweicloud/golangsdk"
)

func CreateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("publicips")
}

func DeleteURL(c *golangsdk.ServiceClient, publicipId string) string {
	return c.ServiceURL("publicips", publicipId)
}

func GetURL(c *golangsdk.ServiceClient, publicipId string) string {
	return c.ServiceURL("publicips", publicipId)
}

func UpdateURL(c *golangsdk.ServiceClient, publicipId string) string {
	return c.ServiceURL("publicips", publicipId)
}
