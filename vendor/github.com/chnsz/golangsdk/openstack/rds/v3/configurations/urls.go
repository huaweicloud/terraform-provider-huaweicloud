package configurations

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("configurations")
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("configurations", id)
}
