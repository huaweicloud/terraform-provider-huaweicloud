package services

import "github.com/huaweicloud/golangsdk"

const (
	rootPath = "vpc-endpoint-services"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, serviceID string) string {
	return c.ServiceURL(rootPath, serviceID)
}

func publicResourceURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, "public")
}

func connectionsURL(c *golangsdk.ServiceClient, serviceID string) string {
	return c.ServiceURL(rootPath, serviceID, "connections")
}

func connectionsActionURL(c *golangsdk.ServiceClient, serviceID string) string {
	return c.ServiceURL(rootPath, serviceID, "connections/action")
}

func permissionsURL(c *golangsdk.ServiceClient, serviceID string) string {
	return c.ServiceURL(rootPath, serviceID, "permissions")
}

func permissionsActionURL(c *golangsdk.ServiceClient, serviceID string) string {
	return c.ServiceURL(rootPath, serviceID, "permissions/action")
}
