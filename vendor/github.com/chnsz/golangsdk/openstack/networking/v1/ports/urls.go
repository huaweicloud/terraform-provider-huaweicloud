package ports

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("ports")
}

func resourceURL(c *golangsdk.ServiceClient, portId string) string {
	return c.ServiceURL("ports", portId)
}
