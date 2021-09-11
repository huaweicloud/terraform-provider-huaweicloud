package instances

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id)
}

func resizeResourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "resize")
}

func updatePasswordURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "password")
}
