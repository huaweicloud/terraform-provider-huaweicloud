package networkipavailabilities

import "github.com/huaweicloud/golangsdk"

const resourcePath = "network-ip-availabilities"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, networkIPAvailabilityID string) string {
	return c.ServiceURL(resourcePath, networkIPAvailabilityID)
}

func listURL(c *golangsdk.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *golangsdk.ServiceClient, networkIPAvailabilityID string) string {
	return resourceURL(c, networkIPAvailabilityID)
}
