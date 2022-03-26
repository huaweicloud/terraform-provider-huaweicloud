package ports

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "ports")
}

func resourceURL(c *golangsdk.ServiceClient, portId string) string {
	return c.ServiceURL(c.ProjectID, "ports", portId)
}
