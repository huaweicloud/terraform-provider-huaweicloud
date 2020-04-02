package hw_snatrules

import "github.com/huaweicloud/golangsdk"

const resourcePath = "snat_rules"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}
