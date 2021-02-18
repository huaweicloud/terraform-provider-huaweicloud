package volumes

import "github.com/huaweicloud/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudvolumes")
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("cloudvolumes", id)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("os-vendor-volumes", id)
}
