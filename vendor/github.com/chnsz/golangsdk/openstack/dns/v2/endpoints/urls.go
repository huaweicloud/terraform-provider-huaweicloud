package endpoints

import "github.com/chnsz/golangsdk"

func baseUrl(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("endpoints")
}

func resourceUrl(c *golangsdk.ServiceClient, endpointID string) string {
	return c.ServiceURL("endpoints", endpointID)
}
