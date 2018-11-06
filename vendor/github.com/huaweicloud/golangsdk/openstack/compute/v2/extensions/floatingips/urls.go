package floatingips

import "github.com/huaweicloud/golangsdk"

const resourcePath = "os-floating-ips"

func resourceURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *golangsdk.ServiceClient) string {
	return resourceURL(c)
}

func createURL(c *golangsdk.ServiceClient) string {
	return resourceURL(c)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return getURL(c, id)
}

func serverURL(c *golangsdk.ServiceClient, serverID string) string {
	return c.ServiceURL("servers/" + serverID + "/action")
}

func associateURL(c *golangsdk.ServiceClient, serverID string) string {
	return serverURL(c, serverID)
}

func disassociateURL(c *golangsdk.ServiceClient, serverID string) string {
	return serverURL(c, serverID)
}
