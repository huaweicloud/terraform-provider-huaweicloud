package volumeattach

import "github.com/huaweicloud/golangsdk"

const resourcePath = "os-volume_attachments"

func resourceURL(c *golangsdk.ServiceClient, serverID string) string {
	return c.ServiceURL("servers", serverID, resourcePath)
}

func listURL(c *golangsdk.ServiceClient, serverID string) string {
	return resourceURL(c, serverID)
}

func createURL(c *golangsdk.ServiceClient, serverID string) string {
	return resourceURL(c, serverID)
}

func getURL(c *golangsdk.ServiceClient, serverID, aID string) string {
	return c.ServiceURL("servers", serverID, resourcePath, aID)
}

func deleteURL(c *golangsdk.ServiceClient, serverID, aID string) string {
	return getURL(c, serverID, aID)
}
