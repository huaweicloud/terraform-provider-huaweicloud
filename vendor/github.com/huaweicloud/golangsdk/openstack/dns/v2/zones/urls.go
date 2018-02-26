package zones

import "github.com/huaweicloud/golangsdk"

func baseURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("zones")
}

func zoneURL(c *golangsdk.ServiceClient, zoneID string) string {
	return c.ServiceURL("zones", zoneID)
}
