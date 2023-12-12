package connections

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("data-connections")
}

func resourceURL(c *golangsdk.ServiceClient, connectionId string) string {
	return c.ServiceURL("data-connections", connectionId)
}

func validateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("data-connections/validation")
}
